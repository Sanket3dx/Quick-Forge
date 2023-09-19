package demonRoutes

import (
	"fmt"
	"quick_forge/controllers/demonController"
	"quick_forge/middlewares"
	"quick_forge/utils"

	"github.com/gin-gonic/gin"
)

func InitDemonRouter() {
	config := utils.GetProjectConfig()
	fmt.Println("👊 Config files are Loaded ...✅ ")
	demonRouter := gin.Default()

	demonRouter.GET("/:endpoint", middlewares.MethodAllowedMiddleware(config), demonController.GetAllhandler)
	demonRouter.GET("/:endpoint/:arg", middlewares.MethodAllowedMiddleware(config), demonController.Gethandler)
	demonRouter.POST("/:endpoint", middlewares.MethodAllowedMiddleware(config), demonController.Posthandler)
	demonRouter.PUT("/:endpoint/:arg", middlewares.MethodAllowedMiddleware(config), demonController.Puthandler)
	demonRouter.DELETE("/:endpoint/:arg", middlewares.MethodAllowedMiddleware(config), demonController.Deletehandler)
	fmt.Println("🤘 Routes Loaded ...✅ ")
	fmt.Println("😎 Router started on Port : " + config.Port + " ...✅ ")
	demonRouter.Run(config.Port)
}
