package user

import (
	"fmt"
	"time"

	"github.com/AlphaCodinggroup/alpha_auth-api/internal/platform/bcrypt"
)

type User struct {
	ID            uint
	IDRol         uint
	Email         string
	Username      string
	Password      string
	TokenHash     string
	RefreshTokens []string
	IsVerified    bool
	Active        bool
	CreatedAt     *time.Time
	UpdatedAt     *time.Time
	CreatedBy     int
	UpdatedBy     int
}

type Login struct {
	UserID          uint64
	Username        string
	Password        string
	LoginAt         time.Time
	IPAddress       string
	DeviceInfo      string
	LogoutAt        *time.Time
	SessionDuration *time.Duration
}

func (user *User) HashPassword() error {
	pass, err := bcrypt.HashPassword(user.Password)
	if err != nil {
		return fmt.Errorf("could not hash password %w", err)
	}

	user.Password = pass
	return nil
}

func (user *User) ComparePasswords(providedPassword string) bool {
	return bcrypt.ComparePasswords(user.Password, providedPassword)
}
