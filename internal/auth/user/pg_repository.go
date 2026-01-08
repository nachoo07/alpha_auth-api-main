package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/AlphaCodinggroup/alpha_auth-api/internal/platform/repository/pg"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) (*userRepository, error) {
	if db == nil {
		return nil, fmt.Errorf("database can not be nil")
	}

	return &userRepository{db: db}, nil
}

func (repo *userRepository) Save(ctx context.Context, user *User) error {
	userModel := entityToModel(user)
	userModel.CreatedBy = user.CreatedBy
	userModel.UpdatedBy = user.CreatedBy

	if err := repo.db.WithContext(ctx).Create(&userModel).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return fmt.Errorf("user already exist")
		}
		return fmt.Errorf("error saving User: %w", err)
	}

	user.ID = userModel.ID
	user.CreatedAt = &userModel.CreatedAt
	user.UpdatedAt = &userModel.UpdatedAt

	return nil
}

func (repo *userRepository) CreateLoginInfo(ctx context.Context, loginInfo *Login) error {
	userLogin := &pg.UserLogin{
		UserID:          loginInfo.UserID,
		IPAddress:       loginInfo.IPAddress,
		DeviceInfo:      loginInfo.DeviceInfo,
		Success:         true,
		SessionDuration: loginInfo.SessionDuration,
	}

	if !loginInfo.LoginAt.IsZero() {
		userLogin.LoginAt = loginInfo.LoginAt
	} else {
		userLogin.LoginAt = time.Now()
	}

	result := repo.db.WithContext(ctx).Create(userLogin)
	if result.Error != nil {
		return fmt.Errorf("error creating login info: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("no rows affected when creating login info")
	}

	return nil
}

func (repo *userRepository) UpdatePassword(ctx context.Context, user *User) error {
	result := repo.db.WithContext(ctx).Model(&pg.User{}).
		Where("email = ?", user.Email).
		Update("password", user.Password).
		Update("token_hash", user.TokenHash).
		Update("refresh_tokens", []string{})

	if result.Error != nil || result.RowsAffected == 0 {
		return fmt.Errorf("error updating User")
	}

	return nil
}

func (repo *userRepository) Update(ctx context.Context, user *User) error {
	result := repo.db.WithContext(ctx).Model(&pg.User{}).
		Where("id = ?", user.ID).
		Updates(map[string]any{
			"password":       user.Password,
			"token_hash":     user.TokenHash,
			"updated_by":     user.UpdatedBy,
			"refresh_tokens": []string{},
		})

	if result.Error != nil || result.RowsAffected == 0 {
		return fmt.Errorf("error updating User")
	}

	return nil
}

func (repo *userRepository) GetUsers(ctx context.Context) ([]User, error) {
	var err error
	var usersModel []pg.User

	result := repo.db.WithContext(ctx).Where("id_rol != 1").Find(&usersModel)
	if err = result.Error; err != nil {
		return nil, fmt.Errorf("error getting all users: %w", err)
	}

	var users []User

	for _, user := range usersModel {
		users = append(users, *modelToEntity(&user))
	}

	return users, nil
}

func (repo *userRepository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var err error
	userModel := new(pg.User)

	if err = repo.db.WithContext(ctx).Where("email = ?", email).Last(userModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, fmt.Errorf("error getting User with email: %s", email)
	}

	return modelToEntity(userModel), nil
}

func (repo *userRepository) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	var err error
	userModel := new(pg.User)

	if err = repo.db.WithContext(ctx).Where("username = ?", username).Last(userModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, fmt.Errorf("error getting User with username: %s", username)
	}

	return modelToEntity(userModel), nil
}

func (repo *userRepository) GetUserByID(ctx context.Context, id uint) (*User, error) {
	var err error
	userModel := new(pg.User)

	err = repo.db.WithContext(ctx).Where("id = ?", id).First(userModel).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			msg := fmt.Sprintf("No User found with id: %d", id)

			return nil, ResourceNotFoundError{Message: msg}
		}

		return nil, fmt.Errorf("error getting User with id: %d", id)
	}

	return modelToEntity(userModel), nil
}

func (repo *userRepository) UpdateIsVerifiedField(ctx context.Context, email string) error {
	result := repo.db.WithContext(ctx).Model(&pg.User{}).
		Where("email = ?", email).
		Update("is_verified", true)

	if result.Error != nil || result.RowsAffected == 0 {
		return fmt.Errorf("error updating User")
	}

	return nil
}

func (repo *userRepository) UpdateRefreshTokens(ctx context.Context, user User) error {
	err := repo.db.WithContext(ctx).Model(&pg.User{}).
		Where("id = ?", user.ID).
		Updates(pg.User{RefreshTokens: user.RefreshTokens}).Error
	if err != nil {
		return fmt.Errorf("error updating User: %w", err)
	}

	return nil
}

func (repo *userRepository) DeleteRefreshToken(ctx context.Context, id uint, token string) error {
	var user pg.User
	err := repo.db.WithContext(ctx).First(&user, id).Error
	if err != nil {
		return fmt.Errorf("error finding User: %w", err)
	}

	var updatedTokens []string
	for _, t := range user.RefreshTokens {
		if t != token {
			updatedTokens = append(updatedTokens, t)
		}
	}

	if len(updatedTokens) == 0 {
		err = repo.db.WithContext(ctx).Model(&pg.User{}).
			Where("id = ?", id).
			Update("refresh_tokens", []string{}).Error
	} else {
		err = repo.db.WithContext(ctx).Model(&pg.User{}).
			Where("id = ?", id).
			Updates(pg.User{RefreshTokens: updatedTokens}).Error
	}

	if err != nil {
		return fmt.Errorf("error deleting refresh token: %w", err)
	}

	return nil
}

func (repo *userRepository) Delete(ctx context.Context, userID uint, deletedBy int) error {
	result := repo.db.WithContext(ctx).Model(&pg.User{}).
		Where("id = ?", userID).
		Updates(map[string]any{
			"deleted_by": deletedBy,
			"deleted_at": gorm.DeletedAt{Time: time.Now(), Valid: true},
		})

	if result.Error != nil || result.RowsAffected == 0 {
		return fmt.Errorf("error deleting user: %w", result.Error)
	}

	return nil
}

func modelToEntity(u *pg.User) *User {
	return &User{
		ID:            u.ID,
		IDRol:         u.IDRol,
		Email:         u.Email,
		Username:      u.Username,
		Password:      u.Password,
		TokenHash:     u.TokenHash,
		RefreshTokens: u.RefreshTokens,
		IsVerified:    u.IsVerified,
		Active:        u.Active,
		CreatedAt:     &u.CreatedAt,
		UpdatedAt:     &u.UpdatedAt,
	}
}

func entityToModel(u *User) pg.User {
	return pg.User{
		IDRol:         u.IDRol,
		Email:         u.Email,
		Username:      u.Username,
		Password:      u.Password,
		TokenHash:     u.TokenHash,
		RefreshTokens: u.RefreshTokens,
		IsVerified:    u.IsVerified,
		Active:        u.Active,
	}
}
