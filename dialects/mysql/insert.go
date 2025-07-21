package mysql

import (
	"fmt"
	"strings"

	"github.com/iMohamedSheta/xqb/shared/types"
)

func (mg *MySQLGrammar) CompileInsert(qb *types.QueryBuilderData) (string, []any, error) {
	tableName, _, err := mg.resolveTable(qb, "insert", false)
	if err != nil {
		return "", nil, err
	}

	if len(qb.InsertedValues) == 0 {
		// Supported in all versions of MySQL for +8.0 use of INSERT INTO {TABLE} DEFAULT VALUES
		return fmt.Sprintf("INSERT INTO %s () VALUES ()", tableName), nil, nil
	}

	var bindings []any

	// Get columns from the first row of values
	columns := make([]string, 0, len(qb.InsertedValues[0]))
	for col := range qb.InsertedValues[0] {
		columns = append(columns, col)
	}

	// Build column names string
	columnStr := strings.Join(columns, ", ")

	// Build values strings
	valueStrings := make([]string, len(qb.InsertedValues))
	for i, row := range qb.InsertedValues {
		placeholders := make([]string, len(columns))
		for j, col := range columns {
			value := row[col]
			if value == nil {
				placeholders[j] = "NULL"
			} else {
				placeholders[j] = "?"
				bindings = append(bindings, value)
			}
		}
		valueStrings[i] = "(" + strings.Join(placeholders, ", ") + ")"
	}

	// Build the final SQL
	sql := fmt.Sprintf("INSERT INTO %s (%s) VALUES %s",
		tableName,
		columnStr,
		strings.Join(valueStrings, ", "),
	)

	return sql, bindings, nil
}
