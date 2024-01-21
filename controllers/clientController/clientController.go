package clientcontroller

import (
	"fmt"
	"net/http"
	"time"

	mysql_models "quick_forge/database/mysql/models"

	"quick_forge/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type LoginResponse struct {
	Message string                 `json:"message"`
	User    map[string]interface{} `json:"user,omitempty"`
	Token   string                 `json:"token,omitempty"`
}

func LoginHandler(c *gin.Context) {
	var loginRequest LoginRequest

	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, LoginResponse{Message: "Invalid request"})
		return
	}

	// Authenticate user
	user, err := mysql_models.AuthenticateUser(loginRequest.Username, loginRequest.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, LoginResponse{Message: "Authentication failed"})
		return
	}
	if user["user_id"] != 0 && user["user_id"] != nil {
		claims := jwt.MapClaims{
			"user_id": user["user_id"],
			"exp":     time.Now().Add(time.Hour * 1).Unix(), // Token will expire in 1 hour
			"nbf":     time.Now().Unix(),
			"iat":     time.Now().Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		secretKey := []byte("kaHFVKAJHFVKAJHFOIWGYFOIAUSFOAKSbe8127467824678fsd78gtasd68trgasdrgsrds7df8")
		// Sign the token with your secret key
		tokenString, err := token.SignedString(secretKey)
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusUnauthorized, LoginResponse{Message: "Authentication failed"})
			return
		}

		c.JSON(http.StatusOK, LoginResponse{
			Message: "Login successful",
			Token:   tokenString,
		})
	} else {
		c.JSON(http.StatusUnauthorized, LoginResponse{Message: "Incorcet Username or Password."})
		return
	}

}

func GetAllRoutes(c *gin.Context) {
	config := utils.GetProjectConfig()
	c.JSON(http.StatusOK, gin.H{
		"error": false,
		"data":  config.Routes,
	})
}
