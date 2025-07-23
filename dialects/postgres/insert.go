package postgres

import (
	"errors"
	"fmt"
	"sort"
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
		return fmt.Sprintf("INSERT INTO %s DEFAULT VALUES", pg.Wrap(tableName)), nil, nil
	}

	columns := getSortedColumns(qb.InsertedValues[0])
	columnStr := wrapColumns(columns, pg.Wrap)

	valueStrings, bindings := buildValuePlaceholders(qb.InsertedValues, columns)

	sql := fmt.Sprintf("INSERT INTO %s (%s) VALUES %s",
		pg.Wrap(tableName),
		columnStr,
		strings.Join(valueStrings, ", "),
	)

	if isUpsert, ok := qb.GetOption(types.OptionIsUpsert); ok && isUpsert.(bool) {
		upsertClause, err := buildUpsertClause(qb, columns, pg.Wrap)
		if err != nil {
			return "", nil, err
		}
		if upsertClause != "" {
			sql += " " + upsertClause
		}
	}

	if len(qb.Errors) > 0 {
		errs := errors.Join(qb.Errors...)
		return "", nil, fmt.Errorf("%w: %s", xqbErr.ErrInvalidQuery, errs)
	}

	return sql, bindings, nil
}

func getSortedColumns(row map[string]any) []string {
	columns := make([]string, 0, len(row))
	for col := range row {
		columns = append(columns, col)
	}
	sort.Strings(columns)
	return columns
}

func wrapColumns(columns []string, wrapFn func(string) string) string {
	wrapped := make([]string, len(columns))
	for i, col := range columns {
		wrapped[i] = wrapFn(col)
	}
	return strings.Join(wrapped, ", ")
}

func buildValuePlaceholders(rows []map[string]any, columns []string) ([]string, []any) {
	var (
		values   = make([]string, len(rows))
		bindings = make([]any, 0, len(rows)*len(columns))
	)

	for i, row := range rows {
		placeholders := make([]string, len(columns))
		for j, col := range columns {
			placeholders[j] = "?"
			bindings = append(bindings, row[col])
		}
		values[i] = "(" + strings.Join(placeholders, ", ") + ")"
	}
	return values, bindings
}

func buildUpsertClause(qb *types.QueryBuilderData, allCols []string, wrapFn func(string) string) (string, error) {
	uniqueVal, ok := qb.GetOption(types.OptionUpsertUniqueBy)
	if !ok {
		return "", fmt.Errorf("%w: you must set the unique by column for the upsert operation", xqbErr.ErrInvalidQuery)
	}

	uniqueBy, ok := uniqueVal.([]string)
	if !ok || len(uniqueBy) == 0 {
		return "", fmt.Errorf("%w: unique by value must be a non-empty []string", xqbErr.ErrInvalidQuery)
	}

	uniqueCols := make(map[string]struct{}, len(uniqueBy))
	for _, col := range uniqueBy {
		uniqueCols[col] = struct{}{}
	}

	updatedVal, ok := qb.GetOption(types.OptionUpsertUpdatedCols)
	if !ok {
		return "", nil
	}
	updatedCols, ok := updatedVal.([]string)
	if !ok || len(updatedCols) == 0 {
		return "", nil
	}

	sort.Strings(updatedCols)

	updates := make([]string, 0, len(updatedCols))
	for _, col := range updatedCols {
		if _, isUnique := uniqueCols[col]; isUnique {
			continue
		}
		wrappedCol := wrapFn(col)
		updates = append(updates, fmt.Sprintf("%s = EXCLUDED.%s", wrappedCol, wrappedCol))
	}

	if len(updates) == 0 {
		return "", nil
	}

	wrappedUnique := make([]string, len(uniqueBy))
	for i, col := range uniqueBy {
		wrappedUnique[i] = wrapFn(col)
	}

	return fmt.Sprintf("ON CONFLICT (%s) DO UPDATE SET %s",
		strings.Join(wrappedUnique, ", "),
		strings.Join(updates, ", "),
	), nil
}
