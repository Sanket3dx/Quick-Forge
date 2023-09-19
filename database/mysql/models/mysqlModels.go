package mysql_models

import (
	"encoding/json"
	"fmt"
	mysql_configer "quick_forge/database/mysql"
)

func GetAllDataAsJSON(tableName string) (string, error) {
	db := mysql_configer.InitDB()
	defer db.Close()
	query := fmt.Sprintf("SELECT * FROM %s", tableName)
	rows, err := db.Query(query)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	// Create a slice to hold the results
	var results []map[string]interface{}

	// Get column names
	columnNames, err := rows.Columns()
	if err != nil {
		return "", err
	}

	// Iterate over the rows
	for rows.Next() {

		columns := make([]interface{}, len(columnNames))
		columnPointers := make([]interface{}, len(columnNames))

		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		err := rows.Scan(columnPointers...)
		if err != nil {
			return "", err
		}

		rowData := make(map[string]interface{})

		for i, colName := range columnNames {
			val := columnPointers[i].(*interface{})
			rowData[colName] = *val
		}

		results = append(results, rowData)
	}
	fmt.Println(results)
	// Marshal the results to JSON format
	jsonData, err := json.Marshal(results)
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}
