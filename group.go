package xqb

import (
	"fmt"

	"github.com/iMohamedSheta/xqb/shared/types"
)

// GroupBy sets the GROUP BY clause
func (qb *QueryBuilder) GroupBy(columns ...any) *QueryBuilder {
	for _, column := range columns {
		var col string
		var bindings []any

		switch v := column.(type) {
		case string:
			col = v
		case *types.Expression:
			col = v.SQL
			bindings = v.Bindings
		default:
			col = fmt.Sprintf("%v", v)
		}

		qb.groupBy = append(qb.groupBy, col)

		if len(bindings) > 0 {
			for _, binding := range bindings {
				qb.bindings = append(qb.bindings, types.Binding{Value: binding})
			}
		}
	}
	return qb
}

// GroupByRaw adds a raw GROUP BY clause
func (qb *QueryBuilder) GroupByRaw(sql string, bindings ...any) *QueryBuilder {
	return qb.GroupBy(Raw(sql, bindings...))
}
