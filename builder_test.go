package xqb_test

import (
	"fmt"
	"testing"

	"github.com/iMohamedSheta/xqb"
	"github.com/iMohamedSheta/xqb/shared/types"
)

func TestOnBeforeQueryHook(t *testing.T) {
	xqb.DefaultSettings().OnAfterQuery(func(query *xqb.QueryExecuted) {
		elapsed := query.Time.String()
		sql := query.Sql
		bindings := query.Bindings

		// inject bindings into SQL for debugging
		boundSql, err := xqb.InjectBindings(query.Dialect, sql, bindings)
		if err != nil {
			xqb.Dump(err)
			return
		}

		xqb.Dump(fmt.Sprintf("[%s] %s", elapsed, boundSql))
	})

	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("users").SetDialect(dialect)
		qb.Select("id", "name").Where("id", "=", 1).ToSql()
	})
}
