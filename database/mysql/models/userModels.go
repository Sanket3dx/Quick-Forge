package mysql_models

import mysql_configer "quick_forge/database/mysql"

func AuthenticateUser(username, password string) (map[string]interface{}, error) {
	db := mysql_configer.InitDBInternal()
	defer db.Close()

	query := "SELECT * FROM app_users WHERE username = ? AND password = ?"

	rows, err := db.Query(query, username, password)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Get column names
	columnNames, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}

	for rows.Next() {
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
		result = rowData
	}
	return result, nil
}
