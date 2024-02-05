package mysql_configer

import (
	"database/sql"
	"fmt"
	"quick_forge/utils"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func DataBaseString() string {
	DbConfig := utils.GetProjectConfig().DB
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", DbConfig.Username, DbConfig.Password, DbConfig.Host, DbConfig.DatabaseName)
}

func InitDB() *sql.DB {
	var DB *sql.DB
	DB, _ = sql.Open("mysql", DataBaseString())
	// set max connection
	DB.SetConnMaxLifetime(50)
	// // set max idle connections
	// DB.SetMaxIdleConns(10)
	// // verify the connection

	DB.SetMaxOpenConns(10000)

	if err := DB.Ping(); err != nil {
		fmt.Println("database connection fail")
		panic(err)
	}
	return DB
}

func CreateTableIfNotExist(route utils.Route) {
	primaryKey, primaryKeyExists := route.DBTableStruct["primary_key"]
	if !primaryKeyExists {
		return
	}

	columnDefinitions := make([]string, 0)

	for key, value := range route.DBTableStruct {
		if key == "primary_key" {
			continue
		}
		columnName := key

		valueMap, ok := value.(map[string]interface{})
		if !ok {
			fmt.Printf("Unexpected type for column '%s'. Skipping.\n", columnName)
			continue
		}

		dType := getDbStructDataType(valueMap)
		columnDefinitions = append(columnDefinitions, fmt.Sprintf("%s %s", columnName, dType))
	}

	primaryKeyColumn := fmt.Sprintf("%s INT AUTO_INCREMENT PRIMARY KEY", primaryKey)
	columnDefinitions = append(columnDefinitions, primaryKeyColumn)

	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s)",
		route.DBTableName, strings.Join(columnDefinitions, ", "))
	println(query)
	db := InitDB()

	db.Exec(query)
	db.Close()
}

func getDbStructDataType(value map[string]interface{}) string {
	if val, ok := value["type"].(string); ok {
		switch val {
		case "int":
			return "INT"
		case "float":
			return "DOUBLE"
		case "string":
			return "VARCHAR(255)"
		default:
			return "TEXT"
		}
	}
	return "TEXT"
}
