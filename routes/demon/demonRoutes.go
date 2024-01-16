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
	demonRouter.Use(middlewares.MethodAllowedMiddleware(config))
	demonRouter.GET("/:endpoint", demonController.GetAllhandler)
	demonRouter.GET("/:endpoint/:arg", demonController.Gethandler)
	demonRouter.POST("/:endpoint", middlewares.ValidateRequest(), demonController.Posthandler)
	demonRouter.PUT("/:endpoint/:arg", middlewares.ValidateRequest(), demonController.Puthandler)
	demonRouter.DELETE("/:endpoint/:arg", demonController.Deletehandler)
	fmt.Println("🤘 Routes Loaded ...✅ ")
	fmt.Println("😎 Router started on Port : " + config.Port + " ...✅ ")
	demonRouter.Run(config.Port)
}
