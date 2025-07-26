package xqb_test

import (
	"fmt"
	"testing"

	"github.com/iMohamedSheta/xqb"
)

func TestMain(m *testing.M) {
	Setup()
	m.Run()
}

func Setup() {
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
}
