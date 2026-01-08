package auth

import (
	"context"
	"fmt"

	"github.com/AlphaCodinggroup/alpha_auth-api/internal/auth/user"
	"github.com/AlphaCodinggroup/alpha_auth-api/internal/platform/strbuilder"
	authcontext "github.com/AlphaCodinggroup/alpha_auth-api/pkg/context"
)

type (
	UserAccountCommand interface {
		Save(ctx context.Context, user *user.User) error
		Update(ctx context.Context, user *user.User) error
		Delete(ctx context.Context, userID uint, deletedBy int) error
	}

	registerUseCase struct {
		userQuery   UserFinder
		userCommand UserAccountCommand
	}
)

func NewRegisterUseCase(userQuery UserFinder, userCommand UserAccountCommand) *registerUseCase {
	return &registerUseCase{
		userQuery:   userQuery,
		userCommand: userCommand,
	}
}

func (uc *registerUseCase) GetUsers(ctx context.Context) ([]user.User, error) {
	logger := authcontext.Logger(ctx)

	result, err := uc.userQuery.GetUsers(ctx)
	if err != nil {
		logger.Error(err)
		return nil, fmt.Errorf("something wrong happened, try again later: %w", err)
	}

	return result, nil
}

func (uc *registerUseCase) Signup(ctx context.Context, userData *user.User) error {
	logger := authcontext.Logger(ctx)

	result, err := uc.userQuery.GetUserByUsername(ctx, userData.Username)
	if err != nil {
		return fmt.Errorf("something wrong happened, try again later: %w", err)
	}

	if result != nil {
		return ErrUserAlreadyExists
	}

	if err := userData.HashPassword(); err != nil {
		logger.Error(err)
		return fmt.Errorf("something wrong happened, try again later: %w", err)
	}

	userData.TokenHash = strbuilder.GenerateRandomString()
	userData.IDRol = 2
	userData.IsVerified = true
	userData.Active = true

	if err := uc.userCommand.Save(ctx, userData); err != nil {
		logger.Error(err)
		return fmt.Errorf("error saving in repository: %w", err)
	}

	return nil
}

func (uc *registerUseCase) GetUserByID(ctx context.Context, id uint) (*user.User, error) {
	logger := authcontext.Logger(ctx)

	result, err := uc.userQuery.GetUserByID(ctx, id)
	if err != nil {
		logger.Error(err)
		return nil, fmt.Errorf("something wrong happened, try again later: %w", err)
	}

	return result, nil
}

func (uc *registerUseCase) UpdateUser(ctx context.Context, userData *user.User) error {
	logger := authcontext.Logger(ctx)

	result, err := uc.userQuery.GetUserByID(ctx, userData.ID)
	if err != nil {
		logger.Error(err)
		return fmt.Errorf("something wrong happened, try again later: %w", err)
	}

	if result == nil {
		return ErrInvalidUser
	}

	if err := userData.HashPassword(); err != nil {
		logger.Error(err)
		return fmt.Errorf("something wrong happened, try again later: %w", err)
	}

	userData.TokenHash = strbuilder.GenerateRandomString()

	if err := uc.userCommand.Update(ctx, userData); err != nil {
		logger.Error(err)
		return fmt.Errorf("error saving in repository: %w", err)
	}

	return nil
}

func (uc *registerUseCase) DeleteUser(ctx context.Context, id uint, deletedBy int) error {
	logger := authcontext.Logger(ctx)

	result, err := uc.userQuery.GetUserByID(ctx, id)
	if err != nil {
		logger.Error(err)
		return fmt.Errorf("something wrong happened, try again later: %w", err)
	}

	if result == nil {
		return ErrInvalidUser
	}

	if err := uc.userCommand.Delete(ctx, id, deletedBy); err != nil {
		logger.Error(err)
		return fmt.Errorf("error deleting in repository: %w", err)
	}

	return nil
}
