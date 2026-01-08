package pg

import (
	"fmt"
	"os"
	"time"

	"github.com/AlphaCodinggroup/alpha_auth-api/internal/platform/bcrypt"
	"github.com/AlphaCodinggroup/alpha_auth-api/internal/platform/strbuilder"
	"github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB //nolint:gochecknoglobals

type User struct {
	gorm.Model
	Email         string
	Username      string `gorm:"unique;not null"`
	Password      string
	TokenHash     string
	RefreshTokens pq.StringArray `gorm:"type:text[]"`
	IDRol         uint           `gorm:"column:id_rol"`
	IsVerified    bool
	Active        bool
	CreatedBy     int  `gorm:"column:created_by"`
	UpdatedBy     int  `gorm:"column:updated_by"`
	DeletedBy     *int `gorm:"column:deleted_by"`
}

type UserLogin struct {
	ID              uint64         `gorm:"primaryKey;column:id"`
	UserID          uint64         `gorm:"column:user_id;not null"`
	LoginAt         time.Time      `gorm:"column:login_at;default:CURRENT_TIMESTAMP"`
	IPAddress       string         `gorm:"column:ip_address"`
	DeviceInfo      string         `gorm:"column:device_info"`
	Success         bool           `gorm:"column:success;default:true"`
	LogoutAt        *time.Time     `gorm:"column:logout_at"`
	SessionDuration *time.Duration `gorm:"column:session_duration"`
}

func (UserLogin) TableName() string {
	return "user_logins"
}

func GetDB() (*gorm.DB, error) {
	if db == nil {
		dialector, err := getDialector()
		if err != nil {
			return nil, err
		}

		db, err = gorm.Open(dialector, &gorm.Config{})
		if err != nil {
			fmt.Printf("# DB ERROR: %s", err.Error())
			return nil, err
		}

		if err := createAdminUser(db); err != nil {
			fmt.Printf("# MIGRATION ERROR: %s", err.Error())
			return nil, err
		}
	}

	return db, nil
}

func createAdminUser(db *gorm.DB) error {
	if err := db.AutoMigrate(&User{}, &UserLogin{}); err != nil {
		return err
	}

	const adminEmail = "admin@example.com"
	adminUsername := os.Getenv("ADMIN_USER")
	adminPass := os.Getenv("ADMIN_PASS")

	if adminUsername == "" || adminPass == "" {
		return fmt.Errorf("username and password cannot be nil")
	}

	hashedPassword, err := bcrypt.HashPassword(adminPass)
	if err != nil {
		return fmt.Errorf("could not hash password %w", err)
	}

	const adminRoleID = 1

	var user User
	result := db.Where("username = ?", adminUsername).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			admin := User{
				Email:      adminEmail,
				Username:   adminUsername,
				Password:   hashedPassword,
				IDRol:      adminRoleID,
				TokenHash:  strbuilder.GenerateRandomString(),
				IsVerified: true,
				Active:     true,
				CreatedBy:  1,
				UpdatedBy:  1,
			}

			if err := db.Create(&admin).Error; err != nil {
				return fmt.Errorf("error al crear el usuario admin: %w", err)
			}
		} else {
			return fmt.Errorf("error finding user admin: %w", result.Error)
		}
	}

	return nil
}

func getDialector() (gorm.Dialector, error) {
	//if os.Getenv("K_SERVICE") != "" {
	//	return connectWithConnectorIAMAuthN()
	//}

	dbUser := os.Getenv("DB_USER")
	dbName := os.Getenv("DB_NAME")
	dbHost := os.Getenv("DB_HOST")
	sslmode := os.Getenv("SSL_MODE")
	dbPass := os.Getenv("DB_PASS")
	dbPort := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPass, dbName, sslmode)

	return postgres.Open(dsn), nil
}

/* func connectWithConnectorIAMAuthN() (gorm.Dialector, error) {
	mustGetenv := func(k string) string {
		v := os.Getenv(k)
		if v == "" {
			log.Fatalf("Warning: %s environment variable not set.", k)
		}
		return v
	}

	var (
		dbUser                 = mustGetenv("DB_USER")                  // e.g. 'service-account-name@project-id.iam'
		dbName                 = mustGetenv("DB_NAME")                  // e.g. 'my-database'
		instanceConnectionName = mustGetenv("INSTANCE_CONNECTION_NAME") // e.g. 'project:region:instance'
		usePrivate             = os.Getenv("PRIVATE_IP")
	)

	d, err := cloudsqlconn.NewDialer(
		context.Background(),
		cloudsqlconn.WithIAMAuthN(),
		cloudsqlconn.WithLazyRefresh(),
	)
	if err != nil {
		return nil, fmt.Errorf("cloudsqlconn.NewDialer: %w", err)
	}
	var opts []cloudsqlconn.DialOption
	if usePrivate != "" {
		opts = append(opts, cloudsqlconn.WithPrivateIP())
	}

	dsn := fmt.Sprintf("user=%s database=%s", dbUser, dbName)
	config, err := pgx.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	config.DialFunc = func(ctx context.Context, network, instance string) (net.Conn, error) {
		return d.Dial(ctx, instanceConnectionName, opts...)
	}
	dbURI := stdlib.RegisterConnConfig(config)
	sqlDB, err := sql.Open("pgx", dbURI)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
	}

	return postgres.New(postgres.Config{Conn: sqlDB}), nil
}
*/
