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
		demonRouter.Use(middlewares.AuthDemonMiddleware())
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
		clientRouter.GET("/routes/:endpoint", middlewares.AuthMiddleware(), clientcontroller.GetRoute)
		clientRouter.POST("/routes", middlewares.AuthMiddleware(), clientcontroller.AddRouteHandler)
		clientRouter.PUT("/routes", middlewares.AuthMiddleware(), clientcontroller.UpdateRouteHandler)
		clientRouter.DELETE("/routes/:endpoint", middlewares.AuthMiddleware(), clientcontroller.DeleteRouteHandler)
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
		clientRouter.GET("/console", func(c *gin.Context) {
			// Serve the HTML content directly
			htmlContent := `
				<!DOCTYPE html>
				<html lang="en">
				<head>
				<meta charset="UTF-8">
				<meta name="viewport" content="width=device-width, initial-scale=1.0">
				<title>WebSocket Terminal</title>
				<style>
					body {
						font-family: monospace;
						background-color: #000;
						color: #0f0;
						padding: 20px;
					}
					#terminal {
						background-color: #000;
						color: #0f0;
						padding: 10px;
						border: 1px solid #0f0;
						overflow-y: auto;
						height: 300px;
					}
					input[type="text"] {
						background-color: transparent;
						border: none;
						color: #0f0;
						width: calc(100% - 20px);
						padding: 5px;
						font-family: monospace;
					}
					input[type="text"]:focus {
						outline: none;
					}
				</style>
				</head>
				<body>
				<div id="terminal"></div>
			
				<script>
					const terminal = document.getElementById('terminal');
					// const input = document.getElementById('input');
					const socket = new WebSocket('wss://quick-forge.onrender.com/x-client/ws');
				
					// Connection opened
					socket.addEventListener('open', (event) => {
						writeToTerminal('WebSocket connection opened');
					});
				
					// Listen for messages
					socket.addEventListener('message', (event) => {
						writeToTerminal('Received message: ' + event.data);
						// Handle the received message as needed
					});
				
					// Connection closed
					socket.addEventListener('close', (event) => {
						writeToTerminal('WebSocket connection closed');
					});
				
					// Handle errors
					socket.addEventListener('error', (event) => {
						writeToTerminal('WebSocket error: ' + event.message);
					});
				
					// Send a message
					function sendMessage() {
						const message = input.value;
						socket.send(message);
						writeToTerminal('Sent message: ' + message);
						input.value = '';
					}
				
					// Function to write to terminal
					function writeToTerminal(message) {
						const line = document.createElement('div');
						line.textContent = message;
						terminal.appendChild(line);
						terminal.scrollTop = terminal.scrollHeight;
					}
				
					// Listen for Enter key press to send message
					input.addEventListener('keyup', (event) => {
						if (event.key === 'Enter') {
							sendMessage();
						}
					});
				</script>
				</body>
				</html>
				`
			c.Header("Content-Type", "text/html")
			c.String(http.StatusOK, htmlContent)
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
