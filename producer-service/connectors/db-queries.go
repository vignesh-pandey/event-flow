package connectors

import (
	"producer-service/logs"
	"context"
	"fmt"
	"regexp"
	"strings"
)

// GetFilteredUsers retrieves users from the database based on dynamic filters.
func (p *Postgres) GetFilteredUsers(filters map[string][]string) ([]map[string]interface{}, error) {
	// Validate input filters to prevent SQL injection
	for key, value := range filters {
		if !isValidColumnName(key) {
			logs.Log.Errorf("invalid filter key: %s\n", key)
			return nil, fmt.Errorf("invalid filter key: %s", key)
		}
		if len(value) != 1 {
			logs.Log.Errorf("multiple values are not supported for filter key: %s\n", key)
			return nil, fmt.Errorf("multiple values are not supported for filter key: %s", key)
		}
	}

	// Base query
	baseQuery := "SELECT * FROM users WHERE 1=1"
	var conditions []string
	var args []interface{}

	// Build dynamic query based on filters
	for key, value := range filters {
		conditions = append(conditions, fmt.Sprintf("%s = $%d", key, len(args)+1))
		args = append(args, value[0]) // Use only the first value for each key
	}

	if len(conditions) > 0 {
		baseQuery += " AND " + strings.Join(conditions, " AND ")
	}

	// Execute the query
	rows, err := p.Conn.Query(context.Background(), baseQuery, args...)
	if err != nil {
		logs.Log.Errorln("Failed to execute query:", err)
		return nil, err
	}
	defer rows.Close()

	// Retrieve column names
	columns := rows.FieldDescriptions()

	// Prepare the result
	var users []map[string]interface{}
	for rows.Next() {
		// Create a slice to hold each column value
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))

		for i := range values {
			valuePtrs[i] = &values[i]
		}

		// Scan the row into the slice
		if err := rows.Scan(valuePtrs...); err != nil {
			logs.Log.Errorln("Failed to scan row: ", err)
			return nil, err
		}

		// Create a map for the row
		user := make(map[string]interface{})
		for i, col := range columns {
			user[string(col.Name)] = values[i]
		}

		// Append to the result
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		logs.Log.Errorln("Failed to retrieve rows:", err)
		return nil, err
	}

	return users, nil
}

// isValidColumnName validates that the column name is alphanumeric with underscores
func isValidColumnName(column string) bool {
	re := regexp.MustCompile("^[a-zA-Z0-9_]+$") // Allows letters, numbers, and underscores
	return re.MatchString(column)
}
