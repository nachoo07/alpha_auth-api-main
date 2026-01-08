package handlers

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/AlphaCodinggroup/alpha_auth-api/internal/auth"
	"github.com/AlphaCodinggroup/alpha_auth-api/internal/auth/user"
	authcontext "github.com/AlphaCodinggroup/alpha_auth-api/pkg/context"
	"github.com/gin-gonic/gin"
)

var PgDuplicateKeyMsg = "duplicate key value violates unique constraint"
var ErrUserNotFound = "No user account exists with given email. Please sign in first"

type (
	LoginUseCase interface {
		Login(ctx context.Context, userLogin *user.Login) (string, string, error)
		Logout(ctx context.Context, userID uint, refreshToken string) error
		GenerateNewAccessToken(username, hash string, userID, rolID uint) (tokenString string, err error)
	}

	loginHandler struct {
		loginUseCase LoginUseCase
	}

	loginDTO struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}
)

func NewLoginHandler(loginUseCase LoginUseCase) *loginHandler {
	return &loginHandler{
		loginUseCase: loginUseCase,
	}
}

func (h *loginHandler) Login(ctx *gin.Context) {
	authCtx := authcontext.New(ctx.Request)
	logger := authcontext.Logger(authCtx)
	logger.Info("Entering UserHandler: Login()")

	body := new(loginDTO)
	if err := decodeAndValidate(ctx.Request, body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	loginInfo := &user.Login{
		Username:   body.Username,
		Password:   body.Password,
		IPAddress:  extractIPAddress(ctx),
		DeviceInfo: extractDeviceInfo(ctx),
		LoginAt:    time.Now(),
	}

	accessToken, refreshToken, err := h.loginUseCase.Login(authCtx, loginInfo)
	if err != nil {
		if ok := errors.Is(err, auth.ErrInvalidUser); ok {
			ctx.JSON(http.StatusNotFound, gin.H{"error": auth.ErrInvalidUser.Error()})
			return
		} else {
			userNotVerifyError := user.UserNotVerifyError{}
			if ok := errors.As(err, &userNotVerifyError); ok {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": userNotVerifyError.Error()})
				return
			}
			if ok := errors.As(err, &user.InvalidCredentialsError{}); ok {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}

			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"login_at":      loginInfo.LoginAt,
	})
}

type logoutDTO struct {
	RefreshToken string `json:"refresh_token" validate:"required" binding:"required"`
}

func (h *loginHandler) Logout(ctx *gin.Context) {
	userIDHeader := ctx.GetHeader("userID")

	userID, err := strconv.Atoi(userIDHeader)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	body := new(logoutDTO)
	if err := decodeAndValidate(ctx.Request, body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.loginUseCase.Logout(ctx, uint(userID), body.RefreshToken); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (h *loginHandler) GenerateAccessToken(ctx *gin.Context) {
	userIDHeader := ctx.GetHeader("userID")
	rolIDHeader := ctx.GetHeader("rolID")
	username := ctx.GetHeader("username")
	hash := ctx.GetHeader("hash")

	userID, err := strconv.Atoi(userIDHeader)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rolID, err := strconv.Atoi(rolIDHeader)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accessToken, err := h.loginUseCase.GenerateNewAccessToken(username, hash, uint(userID), uint(rolID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
	})
}

func extractIPAddress(ginCtx *gin.Context) string {
	clientIP := ginCtx.ClientIP()
	if clientIP == "" {
		clientIP = ginCtx.GetHeader("X-Forwarded-For")
		if clientIP == "" {
			clientIP = ginCtx.GetHeader("X-Real-IP")
			if clientIP == "" {
				clientIP = ginCtx.Request.RemoteAddr
			}
		}
	}

	return clientIP
}

func extractDeviceInfo(ginCtx *gin.Context) string {
	userAgent := ginCtx.GetHeader("User-Agent")
	return userAgent
}
