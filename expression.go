package xqb

import "github.com/iMohamedSheta/xqb/shared/types"

// Raw creates a new raw Sql expression
func Raw(sql string, bindings ...any) *types.Expression {
	return &types.Expression{
		Sql:      sql,
		Bindings: bindings,
	}
}

func RawDialect(defaultDialect string, dialects map[string]*types.Expression) types.DialectExpression {
	return types.DialectExpression{Default: defaultDialect, Dialects: dialects}
}
