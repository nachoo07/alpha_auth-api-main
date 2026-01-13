package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type (
	LoginHandler interface {
		Login(ctx *gin.Context)
		Logout(ctx *gin.Context)
		GenerateAccessToken(ctx *gin.Context)
	}

	SignupHandler interface {
		GetUsers(ctx *gin.Context)
		GetUserByID(ctx *gin.Context)
		Signup(ctx *gin.Context)
		UpdateUser(ctx *gin.Context)
		DeleteUser(ctx *gin.Context)
	}

	AuthMiddleware interface {
		TokenAuthMiddleware() gin.HandlerFunc
		ValidateRefreshToken() gin.HandlerFunc
		ValidateRolMiddleware() gin.HandlerFunc
	}
)

func RegisterRouter(router *gin.Engine, md AuthMiddleware, loginHdl LoginHandler, singUpHdl SignupHandler) {
	basePath := "/api/v1/auth"
	publicRouter := router.Group(basePath)

	publicRouter.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "Pong!"})
	})

	newLoginRouter(loginHdl, md, publicRouter)

	protectedRouter := router.Group(basePath)
	// Middleware to verify AccessToken
	protectedRouter.Use(md.TokenAuthMiddleware())
	protectedRouter.Use(md.ValidateRolMiddleware())
	newSignupRouter(singUpHdl, protectedRouter)

	//protectedRouter.GET("/ping", srv.loginController.Ping)
}
