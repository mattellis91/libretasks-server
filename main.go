package main

import (
	"github.com/gin-gonic/gin"

	"github.com/mattellis91/libretasks-server/config"
	"github.com/mattellis91/libretasks-server/routes"
)

func main() {

	router := gin.Default()

	config.ConnectDB()
	
	routes.UserRoute(router)
	
	router.Run("localhost:6000")
	
}
