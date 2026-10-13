package utils

import (
	"database/sql"
	"fmt"
)

func ScanResultDB(rows *sql.Rows) ([]map[string]any, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to get columns: %v", err)
	}

	payloads := []map[string]any{}
	for rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("failed to get rows: %v", err)
		}
		values := make([]any, len(columns))
		pointers := make([]any, len(columns))
		for i := range values {
			pointers[i] = &values[i]
		}

		if err := rows.Scan(pointers...); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v", err)
		}

		payload := map[string]any{}
		for i, value := range values {
			stringValue := asString(value)
			if stringValue != nil {
				payload[columns[i]] = *stringValue
			} else {
				payload[columns[i]] = value
			}
		}

		payloads = append(payloads, payload)
	}

	return payloads, nil
}

func asString(v interface{}) *string {
	switch x := v.(type) {
	case nil:
		value := ""
		return &value
	case []byte:
		value := string(x)
		return &value
	case string:
		return &x
	default:
		return nil
	}
}
