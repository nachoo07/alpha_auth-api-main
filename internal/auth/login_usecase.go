package auth

import (
	"context"
	"fmt"

	"github.com/AlphaCodinggroup/alpha_auth-api/internal/auth/session"
	"github.com/AlphaCodinggroup/alpha_auth-api/internal/auth/user"
	authcontext "github.com/AlphaCodinggroup/alpha_auth-api/pkg/context"
)

type (
	UserTokenManager interface {
		UpdateRefreshTokens(ctx context.Context, user user.User) error
		DeleteRefreshToken(ctx context.Context, id uint, token string) error
	}

	TokenGenerator interface {
		CreateAccessToken(tokenInfo session.AccessTokenInfo) (string, error)
		CreateRefreshToken(tokenInfo session.RefreshTokenInfo) (string, error)
	}

	LoginCommand interface {
		CreateLoginInfo(ctx context.Context, userInfo *user.Login) error
	}

	loginUseCase struct {
		userQuery      UserFinder
		userCommand    UserTokenManager
		tokenGenerator TokenGenerator
		loginCommand   LoginCommand
	}
)

func NewLoginUseCase(userQuery UserFinder, userCommand UserTokenManager,
	tokenGenerator TokenGenerator, loginCommand LoginCommand) *loginUseCase {
	return &loginUseCase{
		userQuery:      userQuery,
		userCommand:    userCommand,
		tokenGenerator: tokenGenerator,
		loginCommand:   loginCommand,
	}
}

func (uc *loginUseCase) Login(ctx context.Context, userLogin *user.Login) (string, string, error) {
	logger := authcontext.Logger(ctx)
	logger.Info("Entering UserService: Login()")

	userFound, err := uc.userQuery.GetUserByUsername(ctx, userLogin.Username)
	if err != nil {
		return "", "", fmt.Errorf("error getting user: %w", err)
	}

	if userFound == nil {
		return "", "", ErrInvalidUser
	}

	if !userFound.ComparePasswords(userLogin.Password) {
		return "", "", user.InvalidCredentialsError{
			Message: "invalid credentials",
		}
	}

	accessToken, err := uc.tokenGenerator.CreateAccessToken(session.AccessTokenInfo{
		ID:       userFound.ID,
		Username: userFound.Username,
		Rol:      userFound.IDRol,
		Hash:     userFound.TokenHash,
	})
	if err != nil {
		return "", "", err
	}

	refreshToken, err := uc.tokenGenerator.CreateRefreshToken(session.RefreshTokenInfo{
		ID:   userFound.ID,
		Hash: userFound.TokenHash,
	})
	if err != nil {
		return "", "", err
	}

	userFound.RefreshTokens = append(userFound.RefreshTokens, refreshToken)

	err = uc.userCommand.UpdateRefreshTokens(ctx, *userFound)
	if err != nil {
		return "", "", err
	}

	go func() {
		userLogin.UserID = uint64(userFound.ID)
		if err := uc.loginCommand.CreateLoginInfo(context.Background(), userLogin); err != nil {
			logger.Error("Failed to register user login", "error", err, "userID", userLogin.UserID)
		}
	}()

	return accessToken, refreshToken, nil
}

func (uc *loginUseCase) Logout(ctx context.Context, userID uint, refreshToken string) error {
	err := uc.userCommand.DeleteRefreshToken(ctx, userID, refreshToken)
	if err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}

	return nil
}

func (uc *loginUseCase) GenerateNewAccessToken(username, hash string, userID, rolID uint) (tokenString string, err error) {
	tokenString, err = uc.tokenGenerator.CreateAccessToken(session.AccessTokenInfo{
		ID:       userID,
		Username: username,
		Rol:      rolID,
		Hash:     hash,
	})
	if err != nil {
		return tokenString, fmt.Errorf("error signing token: %w", err)
	}

	return
}
