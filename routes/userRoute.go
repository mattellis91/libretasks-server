package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mattellis91/libretasks-server/controllers"
)

func UserRoute(router *gin.Engine) {
	router.POST("/user", controllers.CreateUser())	
	router.GET("/user/:userId", controllers.GetAUser())
}