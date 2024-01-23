package demonRoutes

import (
	"fmt"
	clientcontroller "quick_forge/controllers/clientController"
	"quick_forge/controllers/demonController"
	"quick_forge/middlewares"
	"quick_forge/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitDemonRouter() {
	config := utils.GetProjectConfig()
	projectName := config.ProjectName
	fmt.Println("👊 Config files are Loaded ...✅ ")
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Allow all origins
		AllowMethods:     []string{"*"}, // Allow all HTTP methods
		AllowHeaders:     []string{"*"}, // Allow all headers
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	demonRouter := router.Group("/" + projectName)
	{
		demonRouter.Use(middlewares.MethodAllowedMiddleware())
		demonRouter.GET("/:endpoint", demonController.GetAllhandler)
		demonRouter.GET("/:endpoint/:arg", demonController.Gethandler)
		demonRouter.POST("/:endpoint", middlewares.ValidateRequest(), demonController.Posthandler)
		demonRouter.PUT("/:endpoint/:arg", middlewares.ValidateRequest(), demonController.Puthandler)
		demonRouter.DELETE("/:endpoint/:arg", demonController.Deletehandler)
	}
	clientRouter := router.Group("/x-client")
	{
		clientRouter.POST("/login", clientcontroller.LoginHandler)
		clientRouter.GET("/routes", middlewares.AuthMiddleware(), clientcontroller.GetAllRoutes)
		clientRouter.POST("/routes", middlewares.AuthMiddleware(), clientcontroller.AddRouteHandler)
	}

	fmt.Println("🤘 Routes Loaded ...✅ ")
	fmt.Println("😎 Router started on Port : " + config.Port + " ...✅ ")
	router.Run(config.Port)
}
