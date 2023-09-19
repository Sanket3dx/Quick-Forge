package middlewares

import (
	"net/http"
	"quick_forge/utils"

	"github.com/gin-gonic/gin"
)

func MethodAllowedMiddleware(config utils.ProjectConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the request method
		requestMethod := c.Request.Method

		// Get the requested endpoint
		requestedEndpoint := c.Param("endpoint")
		endpointAllowed := false
		for _, route := range config.Routes {
			if route.Endpoint == requestedEndpoint {
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
			c.Next()
		} else {
			c.JSON(http.StatusMethodNotAllowed, gin.H{
				"error": "Method Not Allowed",
			})
			c.Abort()
		}
	}
}
