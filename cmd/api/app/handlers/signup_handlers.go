package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/AlphaCodinggroup/alpha_auth-api/cmd/api/app/handlers/presenter"
	"github.com/AlphaCodinggroup/alpha_auth-api/internal/auth"
	"github.com/AlphaCodinggroup/alpha_auth-api/internal/auth/user"
	"github.com/gin-gonic/gin"
)

type (
	RegisterUseCase interface {
		GetUsers(ctx context.Context) ([]user.User, error)
		GetUserByID(ctx context.Context, id uint) (*user.User, error)
		Signup(ctx context.Context, user *user.User) error
		UpdateUser(ctx context.Context, user *user.User) error
		DeleteUser(ctx context.Context, id uint, deletedBy int) error
	}

	signupHandler struct {
		registerUsecase RegisterUseCase
	}

	userDTO struct {
		Username        string `validate:"required" binding:"required"`
		Password        string `validate:"required"`
		PasswordConfirm string `json:"password_confirm" validate:"required"`
	}
)

func NewSignupHandler(registerUsecase RegisterUseCase) *signupHandler {
	return &signupHandler{
		registerUsecase: registerUsecase,
	}
}

func (handler *signupHandler) GetUsers(ctx *gin.Context) {
	users, err := handler.registerUsecase.GetUsers(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	var usersResponse []any
	for _, user := range users {
		usersResponse = append(usersResponse, presenter.User(&user))
	}

	ctx.JSON(http.StatusOK, usersResponse)
}

func (handler *signupHandler) GetUserByID(ctx *gin.Context) {
	userID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	user, err := handler.registerUsecase.GetUserByID(ctx, uint(userID))
	if err != nil {
		if errors.Is(err, auth.ErrInvalidUser) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": auth.ErrInvalidUser.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, presenter.User(user))
}

func (handler *signupHandler) Signup(ctx *gin.Context) {
	body := new(userDTO)
	if err := decodeAndValidate(ctx.Request, body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := getUserInfo(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if body.Password != body.PasswordConfirm {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Passwords does not match"})
		return
	}

	user := &user.User{
		Username:  body.Username,
		Password:  body.Password,
		CreatedBy: userID,
	}

	if err := handler.registerUsecase.Signup(ctx, user); err != nil {
		if errors.Is(err, auth.ErrUserAlreadyExists) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": auth.ErrUserAlreadyExists.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": fmt.Sprintf("user %s has been created successfully", body.Username)})
}

func (handler *signupHandler) UpdateUser(ctx *gin.Context) {
	userID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	updatedBy, err := getUserInfo(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	body := new(userDTO)
	if err := decodeAndValidate(ctx.Request, body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if body.Password != body.PasswordConfirm {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Passwords does not match"})
		return
	}

	user := &user.User{
		ID:        uint(userID),
		Username:  body.Username,
		Password:  body.Password,
		UpdatedBy: updatedBy,
	}

	if err := handler.registerUsecase.UpdateUser(ctx, user); err != nil {
		if errors.Is(err, auth.ErrUserAlreadyExists) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": auth.ErrUserAlreadyExists.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("user %s has been created successfully", body.Username)})
}

func (handler *signupHandler) DeleteUser(ctx *gin.Context) {
	userID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	deletedBy, err := getUserInfo(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := handler.registerUsecase.DeleteUser(ctx, uint(userID), deletedBy); err != nil {
		if errors.Is(err, auth.ErrUserAlreadyExists) {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": auth.ErrUserAlreadyExists.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("user %d has been deleted successfully", userID)})
}

func getUserInfo(ctx *gin.Context) (int, error) {
	userIDHeader := ctx.Request.Header.Get("rolID")
	userID, err := strconv.Atoi(userIDHeader)
	if err != nil {
		return userID, fmt.Errorf("context error: invalid userID: %w", err)
	}

	return userID, nil
}
