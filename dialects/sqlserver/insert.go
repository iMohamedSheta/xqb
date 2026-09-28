package sqlserver

import (
	"fmt"
	"sort"
	"strings"

	xqbErr "github.com/iMohamedSheta/xqb/shared/errors"
	"github.com/iMohamedSheta/xqb/shared/types"
)

func (d *SqlServerDialect) CompileInsert(qb *types.QueryBuilderData) (string, []any, error) {
	if qb == nil {
		return "", nil, fmt.Errorf("%w: query builder data is nil", xqbErr.ErrInvalidQuery)
	}

	tableName, _, err := d.resolveTable(qb, "insert", false)
	if err != nil {
		return "", nil, err
	}

	if len(qb.InsertedValues) == 0 {
		return fmt.Sprintf("INSERT INTO %s DEFAULT VALUES", d.Wrap(tableName)), nil, nil
	}

	// SQL Server upsert uses MERGE - reject it explicitly.
	if isUpsert, ok := qb.GetBoolOption(types.OptionIsUpsert); ok && isUpsert {
		return "", nil, fmt.Errorf("%w: upsert is not supported by the SQLServer dialect in this version (use raw MERGE)", xqbErr.ErrUnsupportedFeature)
	}

	columns := getSortedColumns(qb.InsertedValues[0])
	columnStr := wrapColumns(columns, d.Wrap)

	valueStrings, bindings := buildValuePlaceholders(qb.InsertedValues, columns)

	sql := fmt.Sprintf("INSERT INTO %s (%s)", d.Wrap(tableName), columnStr)

	// OUTPUT INSERTED for InsertGetId (RETURNING equivalent on SQL Server)
	if returningId, ok := qb.GetBoolOption(types.OptionReturningId); ok && returningId {
		sql += " OUTPUT INSERTED.[id]"
	}

	sql += fmt.Sprintf(" VALUES %s", strings.Join(valueStrings, ", "))

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
