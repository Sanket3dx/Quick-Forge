package mysql_configer

import (
	"database/sql"
	"fmt"
	"os"
	"quick_forge/utils"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func DataBaseString(DbConfig utils.Database) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", DbConfig.Username, DbConfig.Password, DbConfig.Host, DbConfig.DatabaseName)
}

func InitDB(dbConfig utils.Database) *sql.DB {
	var DB *sql.DB
	DB, _ = sql.Open("mysql", DataBaseString(dbConfig))
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

func CreateTableIfNotExist(dbConfig utils.Database, route utils.Route) {
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
	db := InitDB(dbConfig)

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

func InitDBInternal() *sql.DB {
	var DB *sql.DB
	DbConfig := utils.Database{
		Username:     os.Getenv("DB_USERNAME"),
		Password:     os.Getenv("DB_PASS"),
		Host:         os.Getenv("DB_HOST"),
		DatabaseName: os.Getenv("DB_NAME"),
	}
	dbString := fmt.Sprintf("%s:%s@tcp(%s)/%s", DbConfig.Username, DbConfig.Password, DbConfig.Host, DbConfig.DatabaseName)
	DB, _ = sql.Open("mysql", dbString)
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
