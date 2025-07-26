package xqb

import (
	"fmt"

	xqbErr "github.com/iMohamedSheta/xqb/shared/errors"
	"github.com/iMohamedSheta/xqb/shared/types"
)

// GroupBy sets the GROUP BY clause
func (qb *QueryBuilder) GroupBy(columns ...any) *QueryBuilder {
	if len(columns) == 0 {
		qb.appendError(fmt.Errorf("%w: GroupBy() requires at least one column", xqbErr.ErrInvalidQuery))
	}
	for _, column := range columns {
		var col string
		var bindings []any
		var err error

		switch v := column.(type) {
		case string:
			col = v
		case *types.DialectExpression:
			col, bindings, err = v.ToSql(qb.GetDialect().GetDriver().String())
			if err != nil {
				qb.appendError(fmt.Errorf("%w: GroupBy() invalid DialectExpression - %v", xqbErr.ErrInvalidQuery, err))
			}
		case *types.Expression:
			col = v.Sql
			bindings = v.Bindings
		default:
			col = fmt.Sprintf("%v", v)
		}

		if col == "" {
			qb.appendError(fmt.Errorf("%w: GroupBy() received empty or invalid column", xqbErr.ErrInvalidQuery))
			continue
		}

		qb.groupBy = append(qb.groupBy, col)

		for _, binding := range bindings {
			qb.bindings = append(qb.bindings, &types.Binding{Value: binding})
		}
	}
	return qb
}

// GroupByRaw adds a raw GROUP BY clause
func (qb *QueryBuilder) GroupByRaw(sql string, bindings ...any) *QueryBuilder {
	return qb.GroupBy(Raw(sql, bindings...))
}
