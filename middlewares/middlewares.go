package middlewares

import (
	"fmt"
	"net/http"
	"quick_forge/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func WebsoketBrodcaster(broadcast *chan utils.WebSocketMessage) gin.HandlerFunc {
	return func(c *gin.Context) {
		message := fmt.Sprintf("Request received: Method=%s, URI=%s, Client=%s", c.Request.Method, c.Request.RequestURI, c.ClientIP())
		// Broadcast the combined message
		go brodcast(broadcast, message)
	}
}

func brodcast(broadcast *chan utils.WebSocketMessage, message string) {
	*broadcast <- utils.WebSocketMessage{Message: message}
}
func MethodAllowedMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		config := utils.GetProjectConfig()
		// Get the request method
		requestMethod := c.Request.Method
		requestedEndpoint := c.Param("endpoint")
		var routeInfo utils.Route
		endpointAllowed := false
		for _, route := range config.Routes {
			if route.Endpoint == requestedEndpoint {
				routeInfo = route
				switch requestMethod {
				case "GET":
					endpointAllowed = route.Methods.Get != nil && *route.Methods.Get
				case "POST":
					if c.Request.ContentLength == 0 {
						c.JSON(http.StatusBadRequest, gin.H{
							"error":   "Empty Request Body",
							"message": "POST request body is empty",
						})
						c.Abort()
						return
					}
					endpointAllowed = route.Methods.Post != nil && *route.Methods.Post
				case "PUT":
					endpointAllowed = route.Methods.Put != nil && *route.Methods.Put
				case "DELETE":
					endpointAllowed = route.Methods.Delete != nil && *route.Methods.Delete
				}
				break
			}
		}
		if endpointAllowed {
			c.Set("routeInfo", routeInfo)
			c.Set("auth", config.Auth)
			c.Set("dbType", config.DB.DBType)
			c.Next()
		} else {
			c.JSON(http.StatusMethodNotAllowed, gin.H{
				"error":   "Method Not Allowed",
				"message": requestedEndpoint + " is not defined in project",
			})
			c.Abort()
		}
	}
}

func ValidateRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		routeConfig, exists := c.Get("routeInfo")
		if !exists {
			utils.HandleError(c, http.StatusBadRequest, "Route info not found in context")
			return
		}
		route, _ := routeConfig.(utils.Route)
		var requestBody map[string]interface{}
		err := c.ShouldBindJSON(&requestBody)
		if err != nil {
			utils.HandleError(c, http.StatusBadRequest, "Invalid request body")
			return
		}

		// Check if the POST request body contains the required parameters
		for paramName, paramConfig := range route.DBTableStruct {
			if paramConfigMap, ok := paramConfig.(map[string]interface{}); ok {
				required, requiredExists := paramConfigMap["required"]
				if requiredExists && required.(bool) {
					_, exists := requestBody[paramName]
					if !exists {
						utils.HandleError(c, http.StatusBadRequest, fmt.Sprintf("Required parameter '%s' is missing in the request body", paramName))
						return
					}
				}
			}
		}
		c.Set("validatedRequestBody", requestBody)
		c.Next()
	}
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		secretKey := []byte("kaHFVKAJHFVKAJHFOIWGYFOIAUSFOAKSbe8127467824678fsd78gtasd68trgasdrgsrds7df8")

		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
			c.Abort()
			return

		}

		authHeaderParts := strings.Split(authHeader, " ")
		if len(authHeaderParts) != 2 || authHeaderParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Bearer token"})
			c.Abort()
			return
		}

		tokenString := authHeaderParts[1]

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		_, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func AuthDemonMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth, exists := c.Get("auth")
		if !exists {
			utils.HandleError(c, http.StatusBadRequest, "auth info not found in context")
			return
		}
		authConfig, _ := auth.(utils.Auth)
		routeConfig, exists := c.Get("routeInfo")
		if !exists {
			utils.HandleError(c, http.StatusBadRequest, "Route info not found in context")
			return
		}
		route, _ := routeConfig.(utils.Route)
		if authConfig.Secure && route.Secure {
			secretKey := []byte(authConfig.SecretKey)
			authHeader := c.GetHeader("Authorization")

			if authHeader == "" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
				c.Abort()
				return

			}

			authHeaderParts := strings.Split(authHeader, " ")
			if len(authHeaderParts) != 2 || authHeaderParts[0] != "Bearer" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Bearer token"})
				c.Abort()
				return
			}

			tokenString := authHeaderParts[1]

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				return []byte(secretKey), nil
			})

			if err != nil || !token.Valid {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
				c.Abort()
				return
			}

			_, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
				c.Abort()
				return
			}

			c.Next()
		} else {
			c.Next()
		}
	}
}
