package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

type Projects struct {
	Project []ProjectConfig
}

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
	Get    bool `json:"get"`
	Post   bool `json:"post"`
	Put    bool `json:"put"`
	Delete bool `json:"delete"`
}

type Route struct {
	Endpoint      string                 `json:"endpoint" binding:"required"`
	DBTableName   string                 `json:"db_table_name" binding:"required"`
	Methods       RouteMethod            `json:"methods" binding:"required"`
	DBTableStruct map[string]interface{} `json:"db_table_struct" binding:"required"`
	IsSecure      bool                   `json:"is_secure"`
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

// Correct GetAllProjectConfig function to return pointer to Projects
func (p *Projects) LoadFromJSON(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening JSON file: %v", err)
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return fmt.Errorf("error reading JSON file: %v", err)
	}
	if err := json.Unmarshal(data, p); err != nil {
		return fmt.Errorf("error parsing JSON: %v", err)
	}
	return nil
}

func (p *Projects) SaveToJSON(filePath string) error {
	jsonData, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return fmt.Errorf("error encoding JSON: %v", err)
	}

	if err := ioutil.WriteFile(filePath, jsonData, 0644); err != nil {
		return fmt.Errorf("error writing to JSON file: %v", err)
	}

	fmt.Println("Project configuration updated and saved.")
	return nil
}

func (p *Projects) FindProjectByName(name string) ProjectConfig {

	for i, project := range p.Project {
		if strings.EqualFold(project.ProjectName, name) {
			return p.Project[i]
		}
	}
	return ProjectConfig{}
}

func (p *Projects) AddRoute(projectName string, newRoute Route) error {
	project := p.FindProjectByName(projectName)
	if project.ProjectName == "" {
		return fmt.Errorf("project '%s' not found", projectName)
	}

	for _, route := range project.Routes {
		if route.Endpoint == newRoute.Endpoint {
			return fmt.Errorf("endpoint '%s' already exists", newRoute.Endpoint)
		}
	}

	project.Routes = append(project.Routes, newRoute)

	// Save project configuration after adding a route
	if err := p.SaveToJSON("config/project_config.json"); err != nil {
		return fmt.Errorf("error saving project configuration: %v", err)
	}

	return nil
}

func (p *Projects) UpdateRoute(projectName string, updatedRoute Route) error {
	project := p.FindProjectByName(projectName)
	if project.ProjectName == "" {
		return fmt.Errorf("project '%s' not found", projectName)
	}

	found := false
	for i, route := range project.Routes {
		if route.Endpoint == updatedRoute.Endpoint {
			project.Routes[i] = updatedRoute
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("route '%s' not found", updatedRoute.Endpoint)
	}

	// Save project configuration after updating a route
	if err := p.SaveToJSON("config/project_config.json"); err != nil {
		return fmt.Errorf("error saving project configuration: %v", err)
	}

	return nil
}

func (p *Projects) DeleteRoute(projectName, endpoint string) error {
	project := p.FindProjectByName(projectName)
	if project.ProjectName == "" {
		return fmt.Errorf("project '%s' not found", projectName)
	}
	found := false
	for i, route := range project.Routes {
		if route.Endpoint == endpoint {
			project.Routes = append(project.Routes[:i], project.Routes[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("route '%s' not found", endpoint)
	}

	// Save project configuration after deleting a route
	if err := p.SaveToJSON("config/project_config.json"); err != nil {
		return fmt.Errorf("error saving project configuration: %v", err)
	}

	return nil
}

func IsEndpointAllowed(config *ProjectConfig, endpoint string, method string) bool {
	for _, route := range config.Routes {
		if route.Endpoint == endpoint {
			switch strings.ToUpper(method) {
			case "GET":
				return route.Methods.Get
			case "POST":
				return route.Methods.Post
			case "PUT":
				return route.Methods.Put
			case "DELETE":
				return route.Methods.Delete
			}
		}
	}
	return false
}

func HandleError(ctx *gin.Context, statusCode int, errorMessage string) {
	ctx.JSON(statusCode, gin.H{
		"error": errorMessage,
	})
	ctx.Abort()
}

func ConvertURLValuesToMapSlice(queryParams url.Values) []map[string]interface{} {
	var result []map[string]interface{}

	for key, values := range queryParams {
		for _, value := range values {
			param := map[string]interface{}{
				key: value,
			}
			result = append(result, param)
		}
	}

	return result
}

func GetProjectNameFormGinContext(ctx *gin.Context) string {
	projectNameRaw, exists := ctx.Get("projectName")
	if !exists {
		HandleError(ctx, http.StatusInternalServerError, "projectName not found in context")
		return ""
	}

	projectName, ok := projectNameRaw.(string)
	if !ok {
		HandleError(ctx, http.StatusInternalServerError, "projectName in context is not a string")
		return ""
	}
	return projectName
}
