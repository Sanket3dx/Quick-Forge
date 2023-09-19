package demonRoutes

import (
	"quick_forge/controllers/demonController"
	"quick_forge/middlewares"
	"quick_forge/utils"

	"github.com/gin-gonic/gin"
)

func InitDemonRouter() {
	config := utils.GetProjectConfig()
	demonRouter := gin.Default()

	demonRouter.GET("/:endpoint", middlewares.MethodAllowedMiddleware(config), demonController.GetAllhandler)

	demonRouter.Run(":8091")
}
