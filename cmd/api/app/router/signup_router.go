package router

import (
	"github.com/gin-gonic/gin"
)

func newSignupRouter(ctrl SignupHandler, router *gin.RouterGroup) {
	router.GET("/users", ctrl.GetUsers)
	router.GET("/users/:id", ctrl.GetUserByID)
	router.POST("/users", ctrl.Signup)
	router.PUT("/users/:id", ctrl.UpdateUser)
	router.DELETE("/users/:id", ctrl.DeleteUser)
}
