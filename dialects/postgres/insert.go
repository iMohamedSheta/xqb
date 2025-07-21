package postgres

import (
	"errors"
	"fmt"
	"strings"

	xqbErr "github.com/iMohamedSheta/xqb/shared/errors"
	"github.com/iMohamedSheta/xqb/shared/types"
)

func (pg *PostgresDialect) CompileInsert(qb *types.QueryBuilderData) (string, []any, error) {
	tableName, _, err := pg.resolveTable(qb, "insert", false)
	if err != nil {
		return "", nil, err
	}

	if len(qb.InsertedValues) == 0 {
		return fmt.Sprintf("INSERT INTO %s DEFAULT VALUES", tableName), nil, nil
	}

	var bindings []any

	// Get columns from the first row of values
	columns := make([]string, 0, len(qb.InsertedValues[0]))
	for col := range qb.InsertedValues[0] {
		columns = append(columns, col)
	}

	// Build column names string
	columnStr := strings.Join(columns, ", ")

	// Build values strings ($1, $2, $3), ($4, $5, $6)
	placeholderNum := 1
	valueStrings := make([]string, len(qb.InsertedValues))
	for i, row := range qb.InsertedValues {
		placeholders := make([]string, len(columns))
		for j, col := range columns {
			value := row[col]
			if value == nil {
				placeholders[j] = "NULL"
			} else {
				placeholders[j] = fmt.Sprintf("$%d", placeholderNum)
				placeholderNum++
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

	// Check if there are any errors in building the query
	if len(qb.Errors) > 0 {
		errs := errors.Join(qb.Errors...)
		return "", nil, fmt.Errorf("%w: %s", xqbErr.ErrInvalidQuery, errs)
	}

	return sql, bindings, nil
}
