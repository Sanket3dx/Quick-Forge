package demonRoutes

import (
	"fmt"
	"net/http"
	clientcontroller "quick_forge/controllers/clientController"
	"quick_forge/controllers/demonController"
	mysql_configer "quick_forge/database/mysql"
	"quick_forge/middlewares"
	"quick_forge/utils"
	"sync"

	"github.com/gorilla/websocket"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins
		return true
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var clients = make(map[*websocket.Conn]bool) // connected clients
var clientsMutex sync.Mutex
var broadcast = make(chan utils.WebSocketMessage) // channel to broadcast messages to all clients

func InitDemonRouter() {
	config := utils.GetProjectConfig()
	projectName := config.ProjectName
	if config.DB.DBType == "mysql" {
		for _, val := range config.Routes {
			mysql_configer.CreateTableIfNotExist(val)
		}
	}
	fmt.Println("👊 Config files are Loaded ...✅ ")
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Allow all origins
		AllowMethods:     []string{"*"}, // Allow all HTTP methods
		AllowHeaders:     []string{"*"}, // Allow all headers
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	router.Use(middlewares.WebsoketBrodcaster(&broadcast))
	go handleMessages()

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
		clientRouter.GET("/ws", func(c *gin.Context) {
			w := c.Writer
			r := c.Request
			conn, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				fmt.Println(err)
				return
			}

			clientsMutex.Lock()
			clients[conn] = true
			clientsMutex.Unlock()

			defer func() {
				clientsMutex.Lock()
				delete(clients, conn)
				clientsMutex.Unlock()
				conn.Close()
			}()

			for {
				var msg utils.WebSocketMessage
				err := conn.ReadJSON(&msg)
				if err != nil {
					fmt.Println(err)
					break
				}
			}
		})

	}

	fmt.Println("🤘 Routes Loaded ...✅ ")
	fmt.Println("😎 Router started on Port : " + config.Port + " ...✅ ")
	router.Run(config.Port)
}

func handleMessages() {
	for {
		msg := <-broadcast
		clientsMutex.Lock()
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				fmt.Println(err)
				client.Close()
				delete(clients, client)
			}
		}
		clientsMutex.Unlock()
	}
}
