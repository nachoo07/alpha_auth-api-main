package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"slices"

	"github.com/AlphaCodinggroup/alpha_auth-api/internal/auth/session"
	"github.com/AlphaCodinggroup/alpha_auth-api/internal/auth/user"
	"github.com/gin-gonic/gin"
)

type (
	TokenValidator interface {
		ValidateAccessToken(ctx context.Context, token string) (session.AccessTokenInfo, error)
		ValidateRefreshToken(ctx context.Context, tokenString string) (session.RefreshTokenInfo, error)
	}

	UserQuery interface {
		GetUserByID(ctx context.Context, id uint) (*user.User, error)
	}

	UserCommand interface {
		DeleteRefreshToken(ctx context.Context, id uint, token string) error
	}

	authMiddleware struct {
		tokenService TokenValidator
		userQuery    UserQuery
		userCommand  UserCommand
		// tokenCache   *cache.Cache
	}
)

func NewAuthMiddleware(jWtService TokenValidator,
	userQuery UserQuery, userCommand UserCommand) *authMiddleware {
	return &authMiddleware{
		tokenService: jWtService,
		userQuery:    userQuery,
		userCommand:  userCommand,
		// tokenCache:   cache.New(25*time.Minute, 30*time.Minute),
	}
}

func (am *authMiddleware) TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := extractToken(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userInfo, err := am.tokenService.ValidateAccessToken(c, tokenString)
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Next()
			return
		}

		user, err := am.userQuery.GetUserByID(c, userInfo.ID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid token: authentication failed"})
			return
		}

		if userInfo.Hash != user.TokenHash {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token: user credentials have changed"})
			return
		}

		c.Request.Header.Add("userID", fmt.Sprintf("%d", userInfo.ID))
		c.Request.Header.Add("rolID", fmt.Sprintf("%v", userInfo.Rol))
		c.Request.Header.Add("hash", fmt.Sprintf("%v", userInfo.Hash))
		c.Request.Header.Add("exp", fmt.Sprintf("%v", userInfo.Exp))

		c.Next()
	}
}

func (am *authMiddleware) ValidateRolMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		rolID := c.Request.Header.Get("rolID")
		if rolID == "" || rolID != "1" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbiden"})
			return
		}

		c.Next()
	}
}

func (am *authMiddleware) ValidateRefreshToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := extractToken(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		tokenInfo, err := am.tokenService.ValidateRefreshToken(c, token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, err := am.userQuery.GetUserByID(c, tokenInfo.ID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid token: authentication failed"})
			return
		}

		if tokenInfo.Hash != user.TokenHash {
			if err := am.userCommand.DeleteRefreshToken(c, user.ID, token); err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "error deleting refresh token"})
				return
			}
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "invalid token: authentication failed"})
			return
		}

		isValid := slices.Contains(user.RefreshTokens, token)
		if !isValid {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid token: authentication failed"})
			return
		}

		c.Request.Header.Add("userID", fmt.Sprintf("%v", user.ID))
		c.Request.Header.Add("username", user.Username)
		c.Request.Header.Add("hash", user.TokenHash)
		c.Request.Header.Add("rolID", fmt.Sprintf("%v", user.IDRol))
		c.Next()
	}
}

func extractToken(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	authHeaderContent := strings.Split(authHeader, " ")
	if len(authHeaderContent) != 2 {
		return "", fmt.Errorf("token not provided or malformed")
	}
	return authHeaderContent[1], nil
}
