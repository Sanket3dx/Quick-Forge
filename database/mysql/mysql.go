package mysql_configer

import (
	"database/sql"
	"fmt"
	"quick_forge/utils"

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
