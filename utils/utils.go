package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Database struct {
	DBType       string `json:"db_type"`
	Host         string `json:"host"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	DatabaseName string `json:"database_name"`
	Port         string `json:"port"`
}

type RouteMethod struct {
	Get    *bool `json:"get"`
	Post   *bool `json:"post"`
	Put    *bool `json:"put"`
	Delete *bool `json:"delete"`
}

type Route struct {
	Endpoint      string                 `json:"endpoint"`
	DBTableName   string                 `json:"db_table_name"`
	Methods       RouteMethod            `json:"methods"`
	DBTableStruct map[string]interface{} `json:"db_table_struct"`
}

type ProjectConfig struct {
	ProjectName string   `json:"project_name"`
	DB          Database `json:"db"`
	Routes      []Route  `json:"routes"`
	Port        string   `json:"port"`
}

func GetProjectConfig() ProjectConfig {
	var config ProjectConfig
	mydir, _ := os.Getwd()
	filePath := filepath.Join(mydir, "config", "project_config.json")
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
