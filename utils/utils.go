package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

type Database struct {
	DBType       string `json:"db_type"`
	Host         string `json:"host"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	DatabaseName string `json:"database_name"`
	Port         string `json:"port"`
	MongoString  string `json:"mongo_string"`
}

type RouteMethod struct {
	Get    *bool `json:"get"`
	Post   *bool `json:"post"`
	Put    *bool `json:"put"`
	Delete *bool `json:"delete"`
}

type Route struct {
	Endpoint      string                 `json:"endpoint" binding:"required"`
	DBTableName   string                 `json:"db_table_name" binding:"required"`
	Methods       RouteMethod            `json:"methods" binding:"required"`
	DBTableStruct map[string]interface{} `json:"db_table_struct" binding:"required"`
	Secure        bool                   `json:"secure" binding:"required"`
}

type ProjectConfig struct {
	ProjectName string   `json:"project_name"`
	DB          Database `json:"db"`
	Routes      []Route  `json:"routes"`
	Auth        Auth     `json:"auth"`
	Port        string   `json:"port"`
}

type Auth struct {
	Secure    bool   `json:"secure"`
	SecretKey string `json:"secret_key"`
}
type WebSocketMessage struct {
	Message string `json:"message"`
}

func GetProjectConfig() ProjectConfig {
	var config ProjectConfig
	mydir, _ := os.Getwd()
	filePath := filepath.Join(mydir, "/config", "project_config.json")
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Error opening JSON file: %v", err)
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file) // Read the entire file content
	if err != nil {
		log.Fatalf("Error reading JSON file: %v", err)
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}
	fmt.Println(config.DB)
	return config
}

func IsEndpointAllowed(config ProjectConfig, endpoint string, method string) bool {
	for _, route := range config.Routes {
		if route.Endpoint == endpoint {
			switch strings.ToUpper(method) {
			case "GET":
				return route.Methods.Get != nil && *route.Methods.Get
			case "POST":
				return route.Methods.Post != nil && *route.Methods.Post
			case "PUT":
				return route.Methods.Put != nil && *route.Methods.Put
			case "DELETE":
				return route.Methods.Delete != nil && *route.Methods.Delete
			}
		}
	}
	return false
}

func AddRoute(newRoute Route) error {
	config := GetProjectConfig()

	for _, route := range config.Routes {
		if route.Endpoint == newRoute.Endpoint {
			return fmt.Errorf("Endpoint '%s' already exists", newRoute.Endpoint)
		}
	}

	config.Routes = append(config.Routes, newRoute)

	err := saveProjectConfig(config)
	if err != nil {
		return fmt.Errorf("Error saving project configuration: %v", err)
	}

	return nil
}

func UpdateRoute(updatedRoute Route) error {
	config := GetProjectConfig()

	// Find the index of the route to update
	var indexToUpdate int = -1
	for i, route := range config.Routes {
		if route.Endpoint == updatedRoute.Endpoint {
			indexToUpdate = i
			break
		}
	}

	if indexToUpdate == -1 {
		return fmt.Errorf("Route '%s' not found", updatedRoute.Endpoint)
	}

	// Update the route
	config.Routes[indexToUpdate] = updatedRoute

	err := saveProjectConfig(config)
	if err != nil {
		return fmt.Errorf("Error saving project configuration: %v", err)
	}

	return nil
}

func DeleteRoute(endpoint string) error {
	config := GetProjectConfig()
	var indexToDelete int = -1
	for i, route := range config.Routes {
		if route.Endpoint == endpoint {
			indexToDelete = i
			break
		}
	}
	if indexToDelete == -1 {
		return fmt.Errorf("Route '%s' not found", endpoint)
	}
	config.Routes = append(config.Routes[:indexToDelete], config.Routes[indexToDelete+1:]...)

	// Save the updated configuration
	err := saveProjectConfig(config)
	if err != nil {
		return fmt.Errorf("Error saving project configuration: %v", err)
	}

	return nil
}

func saveProjectConfig(config ProjectConfig) error {
	mydir, _ := os.Getwd()
	filePath := filepath.Join(mydir, "/config", "project_config.json")

	jsonData, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("Error encoding JSON: %v", err)
	}

	err = ioutil.WriteFile(filePath, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("Error writing to JSON file: %v", err)
	}

	fmt.Println("Project configuration updated and saved.")
	return nil
}

func HandleError(ctx *gin.Context, statusCode int, errorMessage string) {
	ctx.JSON(statusCode, gin.H{
		"error": errorMessage,
	})
	ctx.Abort()
}
