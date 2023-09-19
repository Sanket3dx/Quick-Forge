package mysql_models

import (
	"fmt"
	mysql_configer "quick_forge/database/mysql"
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
