package auth

import (
	"context"
	"errors"

	"github.com/AlphaCodinggroup/alpha_auth-api/internal/auth/user"
)

var ErrInvalidUser = errors.New("user do not exists with the given username")
var ErrUserAlreadyExists = errors.New("user already exists with the given username")

type UserFinder interface {
	GetUsers(ctx context.Context) ([]user.User, error)
	GetUserByUsername(ctx context.Context, username string) (*user.User, error)
	GetUserByID(ctx context.Context, id uint) (*user.User, error)
}
