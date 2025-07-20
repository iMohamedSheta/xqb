package mysql

import (
	"fmt"
	"sort"
	"strings"

	"github.com/iMohamedSheta/xqb/shared/enums"
	"github.com/iMohamedSheta/xqb/shared/types"
)

// MySQLGrammar implements MySQL-specific SQL syntax
type MySQLGrammar struct {
}

// CompileSelect generates a SELECT SQL statement for MySQL
func (mg *MySQLGrammar) CompileSelect(qb *types.QueryBuilderData) (string, []any, error) {
	// If there are no unions, just compile the base query
	if len(qb.Unions) == 0 {
		return mg.compileBaseQuery(qb)
	}

	var bindings []any
	var sql strings.Builder

	// Get base query SQL and bindings
	baseSQL, baseBindings, err := mg.compileBaseQuery(qb)
	if err != nil {
		return "", nil, err
	}

	sql.WriteString(baseSQL)
	if baseBindings != nil {
		bindings = append(bindings, baseBindings...)
	}

	// Add each union
	for _, union := range qb.Unions {
		switch union.Type {
		case types.UnionTypeUnion:
			sql.WriteString(" UNION ")
		case types.UnionTypeIntersect, types.UnionTypeExcept:
			// unsupported types in MySQL
			return "", nil, fmt.Errorf("xqb_error_mysql_union_type: union type %s is not supported in MySQL", string(union.Type))
		}

		if union.All {
			sql.WriteString("ALL ")
		}

		// Add the union query
		sql.WriteString("(")
		sql.WriteString(union.Expression.SQL)
		sql.WriteString(")")

		if len(union.Expression.Bindings) > 0 {
			bindings = append(bindings, union.Expression.Bindings...)
		}
	}

	return sql.String(), bindings, nil
}

// compileBaseQuery compiles a query without unions
func (mg *MySQLGrammar) compileBaseQuery(qb *types.QueryBuilderData) (string, []any, error) {
	var bindings []any
	var sql strings.Builder

	// Compile each part of the query in order
	cteSQL, cteBindings, _ := mg.compileCTEs(qb)
	if cteSQL != "" {
		sql.WriteString(cteSQL)
		sql.WriteString(" ")
		bindings = append(bindings, cteBindings...)
	}

	selectSQL, selectBindings, _ := mg.compileSelectClause(qb)
	sql.WriteString(selectSQL)
	if selectBindings != nil {
		bindings = append(bindings, selectBindings...)
	}

	fromSQL, fromBindings, _ := mg.compileFromClause(qb)
	if fromSQL != "" {
		sql.WriteString(fromSQL)
		if fromBindings != nil {
			bindings = append(bindings, fromBindings...)
		}
	}

	joinsSQL, joinsBindings, _ := mg.compileJoins(qb)
	if joinsSQL != "" {
		sql.WriteString(joinsSQL)
		if joinsBindings != nil {
			bindings = append(bindings, joinsBindings...)
		}
	}

	whereSQL, whereBindings, _ := mg.compileWhereClause(qb)
	if whereSQL != "" {
		sql.WriteString(whereSQL)
		if whereBindings != nil {
			bindings = append(bindings, whereBindings...)
		}
	}

	groupBySQL, groupByBindings, _ := mg.compileGroupByClause(qb)
	if groupBySQL != "" {
		sql.WriteString(groupBySQL)
		if groupByBindings != nil {
			bindings = append(bindings, groupByBindings...)
		}
	}

	havingSQL, havingBindings, _ := mg.compileHavingClause(qb)
	if havingSQL != "" {
		sql.WriteString(havingSQL)
		if havingBindings != nil {
			bindings = append(bindings, havingBindings...)
		}
	}

	orderBySQL, orderByBindings, _ := mg.compileOrderByClause(qb)
	if orderBySQL != "" {
		sql.WriteString(orderBySQL)
		if orderByBindings != nil {
			bindings = append(bindings, orderByBindings...)
		}
	}

	limitSQL, limitBindings, _ := mg.compileLimitClause(qb)
	if limitSQL != "" {
		sql.WriteString(limitSQL)
		if limitBindings != nil {
			bindings = append(bindings, limitBindings...)
		}
	}

	offsetSQL, offsetBindings, _ := mg.compileOffsetClause(qb)
	if offsetSQL != "" {
		sql.WriteString(offsetSQL)
		if offsetBindings != nil {
			bindings = append(bindings, offsetBindings...)
		}
	}

	lockSQL, lockBindings, _ := mg.compileLockingClause(qb)
	if lockSQL != "" {
		sql.WriteString(lockSQL)
		if lockBindings != nil {
			bindings = append(bindings, lockBindings...)
		}
	}

	return sql.String(), bindings, nil
}

func (mg *MySQLGrammar) CompileInsert(qb *types.QueryBuilderData) (string, []any, error) {
	if qb.Table.Name == "" {
		return "", nil, fmt.Errorf("table name is required for insert operation")
	}

	if len(qb.InsertedValues) == 0 {
		return fmt.Sprintf("INSERT INTO %s DEFAULT VALUES", qb.Table.Name), nil, nil
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
		qb.Table.Name,
		columnStr,
		strings.Join(valueStrings, ", "),
	)

	return sql, bindings, nil
}

func (mg *MySQLGrammar) CompileUpdate(qb *types.QueryBuilderData) (string, []any, error) {
	if len(qb.UpdatedBindings) == 0 {
		return "", nil, fmt.Errorf("no bindings provided for update operation")
	}

	// Sort bindings by column name for consistency
	sort.Slice(qb.Bindings, func(i, j int) bool {
		return qb.Bindings[i].Column < qb.Bindings[j].Column
	})

	var setParts []string

	var bindings []any
	var sql strings.Builder

	for _, binding := range qb.UpdatedBindings {
		setParts = append(setParts, fmt.Sprintf("%s = ?", binding.Column))
		bindings = append(bindings, binding.Value)
	}

	sql.WriteString(fmt.Sprintf("UPDATE %s SET %s", qb.Table.Name, strings.Join(setParts, ", ")))

	whereSQL, whereBindings, _ := mg.compileWhereClause(qb)

	if whereSQL != "" {
		sql.WriteString(whereSQL)
		if whereBindings != nil {
			bindings = append(bindings, whereBindings...)
		}
	}

	limitSQL, limitBindings, _ := mg.compileLimitClause(qb)
	if limitSQL != "" {
		sql.WriteString(limitSQL)
		if limitBindings != nil {
			bindings = append(bindings, limitBindings...)
		}
	}

	return sql.String(), bindings, nil
}

func (mg *MySQLGrammar) CompileDelete(qb *types.QueryBuilderData) (string, []any, error) {
	if len(qb.Bindings) == 0 {
		return "", nil, fmt.Errorf("no bindings provided for delete operation this will destroy all data")
	}

	var bindings []any
	var sql strings.Builder

	sql.WriteString(fmt.Sprintf("DELETE FROM %s", qb.Table.Name))

	whereSQL, whereBindings, _ := mg.compileWhereClause(qb)

	if whereSQL != "" {
		sql.WriteString(whereSQL)
		if whereBindings != nil {
			bindings = append(bindings, whereBindings...)
		}
	}

	limitSQL, limitBindings, _ := mg.compileLimitClause(qb)
	if limitSQL != "" {
		sql.WriteString(limitSQL)
		if limitBindings != nil {
			bindings = append(bindings, limitBindings...)
		}
	}

	return sql.String(), bindings, nil
}

func (mg *MySQLGrammar) Build(qbd *types.QueryBuilderData) (string, []any, error) {
	var sql string
	var bindings []any
	var err error

	switch qbd.QueryType {
	case enums.SELECT:
		sql, bindings, err = mg.CompileSelect(qbd)
	case enums.INSERT:
		sql, bindings, err = mg.CompileInsert(qbd)
	case enums.UPDATE:
		sql, bindings, err = mg.CompileUpdate(qbd)
	case enums.DELETE:
		sql, bindings, err = mg.CompileDelete(qbd)
	}

	if err != nil {
		return "", nil, err
	}

	if sql == "" {
		return "", nil, fmt.Errorf("xqb_error_unexpected_error: couldn't build the query sql is empty")
	}

	return sql, bindings, nil
}
