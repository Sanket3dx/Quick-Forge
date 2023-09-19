package middlewares

import (
	"fmt"
	"net/http"
	"quick_forge/utils"

	"github.com/gin-gonic/gin"
)

func MethodAllowedMiddleware(config utils.ProjectConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
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
			c.Next()
		} else {
			c.JSON(http.StatusMethodNotAllowed, gin.H{
				"error":   "Method Not Allowed",
				"message": requestedEndpoint + "is not defined in project",
			})
			c.Abort()
		}
	}
}

func ValidateRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		routeConfig, exists := c.Get("routeInfo")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Route info not found in context",
			})
			c.Abort()
			return
		}
		route, _ := routeConfig.(utils.Route)
		var requestBody map[string]interface{}
		err := c.ShouldBindJSON(&requestBody)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			c.Abort()
			return
		}

		// Check if the POST request body contains the required parameters
		for paramName, paramConfig := range route.DBTableStruct {
			if paramConfigMap, ok := paramConfig.(map[string]interface{}); ok {
				required, requiredExists := paramConfigMap["required"]
				if requiredExists && required.(bool) {
					_, exists := requestBody[paramName]
					if !exists {
						c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Required parameter '%s' is missing in the request body", paramName)})
						c.Abort()
						return
					}
				}
			}
		}

		c.Next()
	}
}
