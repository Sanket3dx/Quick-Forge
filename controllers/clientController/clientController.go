package clientcontroller

import (
	"fmt"
	"net/http"
	"os"
	mongooDbModel "quick_forge/database/mongooDb/models"
	mysql_models "quick_forge/database/mysql/models"
	"time"

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

	dbType := os.Getenv("DB_TYPE")
	secretKey := []byte("kaHFVKAJHFVKAJHFOIWGYFOIAUSFOAKSbe8127467824678fsd78gtasd68trgasdrgsrds7df8")
	fmt.Println(dbType)
	switch dbType {

	case "mysql":
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
	case "mongo":
		// Authenticate user
		user, err := mongooDbModel.AuthenticateUser(loginRequest.Username, loginRequest.Password)

		if err != nil {
			c.JSON(http.StatusUnauthorized, LoginResponse{Message: err.Error()})
			return
		}

		if user != nil {
			claims := jwt.MapClaims{
				"user_id": user["_id"],
				"exp":     time.Now().Add(time.Hour * 1).Unix(), // Token will expire in 1 hour
				"nbf":     time.Now().Unix(),
				"iat":     time.Now().Unix(),
			}
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
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
			c.JSON(http.StatusUnauthorized, LoginResponse{Message: "Incorrect Username or Password."})
			return
		}

	}

}

func GetAllRoutes(c *gin.Context) {
	projectName := utils.GetProjectNameFormGinContext(c)
	configs := utils.Projects{}
	if err := configs.LoadFromJSON("config/project_config.json"); err != nil {
		utils.HandleError(c, http.StatusInternalServerError, err.Error())
		return
	}
	config := configs.FindProjectByName(projectName)
	c.JSON(http.StatusOK, gin.H{
		"error": false,
		"data":  config.Routes,
	})
}

func GetRoute(c *gin.Context) {
	projectName := utils.GetProjectNameFormGinContext(c)
	endpoint := c.Param("endpoint")

	configs := utils.Projects{}
	if err := configs.LoadFromJSON("config/project_config.json"); err != nil {
		utils.HandleError(c, http.StatusInternalServerError, err.Error())
		return
	}
	config := configs.FindProjectByName(projectName)
	for _, route := range config.Routes {
		if route.Endpoint == endpoint {
			c.JSON(http.StatusOK, gin.H{"error": false, "data": route})
			return
		}
	}

	utils.HandleError(c, http.StatusNotFound, "Route not found")
}

func AddRouteHandler(c *gin.Context) {
	projectName := utils.GetProjectNameFormGinContext(c)
	var newRoute utils.Route
	if err := c.ShouldBindJSON(&newRoute); err != nil {
		utils.HandleError(c, http.StatusBadRequest, err.Error())
		return
	}

	config := utils.Projects{}
	if err := config.LoadFromJSON("config/project_config.json"); err != nil {
		utils.HandleError(c, http.StatusInternalServerError, err.Error())
		return
	}

	if err := config.AddRoute(projectName, newRoute); err != nil {
		utils.HandleError(c, http.StatusBadRequest, err.Error())
		return
	}
	//mysql_configer.CreateTableIfNotExist(config.FindProjectByName(projectName).DB, newRoute)

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "New route added successfully"})
}

func UpdateRouteHandler(c *gin.Context) {
	projectName := utils.GetProjectNameFormGinContext(c)
	var updatedRoute utils.Route
	if err := c.ShouldBindJSON(&updatedRoute); err != nil {
		utils.HandleError(c, http.StatusBadRequest, err.Error())
		return
	}

	config := utils.Projects{}
	if err := config.LoadFromJSON("config/project_config.json"); err != nil {
		utils.HandleError(c, http.StatusInternalServerError, err.Error())
		return
	}

	if err := config.UpdateRoute(projectName, updatedRoute); err != nil {
		utils.HandleError(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Route updated successfully"})
}

func DeleteRouteHandler(c *gin.Context) {
	projectName := utils.GetProjectNameFormGinContext(c)
	endpoint := c.Param("endpoint")

	config := utils.Projects{}
	if err := config.LoadFromJSON("config/project_config.json"); err != nil {
		utils.HandleError(c, http.StatusInternalServerError, err.Error())
		return
	}

	if err := config.DeleteRoute(projectName, endpoint); err != nil {
		utils.HandleError(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"error": false, "message": "Route deleted successfully"})
}

func GetAllProjects(c *gin.Context) {

	config := utils.Projects{}
	if err := config.LoadFromJSON("config/project_config.json"); err != nil {
		utils.HandleError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"error": false,
		"data":  config,
	})
}

func GetAllProjectNames(c *gin.Context) {
	config := utils.Projects{}
	if err := config.LoadFromJSON("config/project_config.json"); err != nil {
		utils.HandleError(c, http.StatusInternalServerError, err.Error())
		return
	}

	projectNames := make([]string, len(config.Project))
	for i, project := range config.Project {
		projectNames[i] = project.ProjectName
	}

	c.JSON(http.StatusOK, gin.H{
		"error": false,
		"data":  projectNames,
	})
}

func GetProject(c *gin.Context) {
	project := c.Param("project")

	configs := utils.Projects{}
	if err := configs.LoadFromJSON("config/project_config.json"); err != nil {
		utils.HandleError(c, http.StatusInternalServerError, err.Error())
		return
	}
	config := configs.FindProjectByName(project)
	c.JSON(http.StatusOK, gin.H{
		"error": false,
		"data":  config,
	})
}
