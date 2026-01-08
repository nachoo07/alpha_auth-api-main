package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func newLoginRouter(ctrl LoginHandler, md AuthMiddleware, router *gin.RouterGroup) {
	router.POST("/login", ctrl.Login)
	router.POST("/logout", md.TokenAuthMiddleware(), ctrl.Logout)
	router.GET("/validate-token", md.TokenAuthMiddleware(), func(ctx *gin.Context) {
		userID := ctx.Request.Header.Get("userID")
		rolID := ctx.Request.Header.Get("rolID")
		hash := ctx.Request.Header.Get("hash")
		exp := ctx.Request.Header.Get("exp")

		ctx.JSON(http.StatusOK, gin.H{
			"status": "valid",
			"userID": userID,
			"rolID":  rolID,
			"hash":   hash,
			"exp":    exp,
		})
	})
	router.POST("/access-token", md.ValidateRefreshToken(), ctrl.GenerateAccessToken)
}
