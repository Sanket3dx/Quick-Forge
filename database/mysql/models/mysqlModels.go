package mysql_models

import (
	"fmt"
	mysql_configer "quick_forge/database/mysql"
	"quick_forge/utils"
	"strings"
)

func GetAllData(tableName string) ([]map[string]interface{}, error) {
	db := mysql_configer.InitDB()
	defer db.Close()
	query := fmt.Sprintf("SELECT * FROM %s", tableName)
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Get column names
	columnNames, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var results []map[string]interface{}

	for rows.Next() {
		// Create a slice to hold the column values
		columns := make([]interface{}, len(columnNames))
		columnPointers := make([]interface{}, len(columnNames))

		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		err := rows.Scan(columnPointers...)
		if err != nil {
			return nil, err
		}

		rowData := make(map[string]interface{})

		for i, colName := range columnNames {
			val := columnPointers[i].(*interface{})
			if bytesVal, ok := (*val).([]byte); ok {
				rowData[colName] = string(bytesVal)
			} else {
				rowData[colName] = *val
			}
		}
		results = append(results, rowData)
	}
	return results, nil
}

func GetData(route utils.Route, arg string) ([]map[string]interface{}, error) {
	db := mysql_configer.InitDB()
	defer db.Close()

	// Ensure that the primary key column exists in the struct
	primaryKey, primaryKeyExists := route.DBTableStruct["primary_key"]
	if !primaryKeyExists {
		return nil, fmt.Errorf("primary key not found in DBTableStruct")
	}

	query := fmt.Sprintf("SELECT * FROM %s WHERE %s = ?", route.DBTableName, primaryKey)

	// Execute the query with a prepared statement
	rows, err := db.Query(query, arg)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Get column names
	columnNames, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var results []map[string]interface{}

	for rows.Next() {
		// Create a slice to hold the column values
		columns := make([]interface{}, len(columnNames))
		columnPointers := make([]interface{}, len(columnNames))

		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		err := rows.Scan(columnPointers...)
		if err != nil {
			return nil, err
		}

		rowData := make(map[string]interface{})

		for i, colName := range columnNames {
			val := columnPointers[i].(*interface{})
			if bytesVal, ok := (*val).([]byte); ok {
				rowData[colName] = string(bytesVal)
			} else {
				rowData[colName] = *val
			}
		}
		results = append(results, rowData)
	}
	return results, nil
}

func InsertData(tableName string, data map[string]interface{}) (int64, error) {
	db := mysql_configer.InitDB()
	defer db.Close()

	// Create placeholders for the column names and values
	var columns []string
	var placeholders []string
	var values []interface{}

	for colName, colValue := range data {
		columns = append(columns, colName)
		placeholders = append(placeholders, "?")
		values = append(values, colValue)
	}

	// Build the SQL query
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		tableName,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "))

	// Execute the insert query with prepared statement
	result, err := db.Exec(query, values...)
	if err != nil {
		return 0, err
	}

	// Get the ID of the newly inserted row
	insertID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return insertID, nil
}
