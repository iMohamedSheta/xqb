package xqb_test

import (
	"os"
	"testing"

	"github.com/iMohamedSheta/xqb/shared/types"
)

func TestMain(m *testing.M) {
	Setup()
	os.Exit(m.Run())
}

func Setup() {
	// xqb.DefaultSettings().OnAfterQuery(func(query *xqb.QueryExecuted) {
	// 	elapsed := query.Time.String()
	// 	sql := query.Sql
	// 	bindings := query.Bindings

	// 	if sql == "" {
	// 		return
	// 	}

	// 	boundSql, err := xqb.InjectBindings(query.Dialect, sql, bindings)
	// 	if err != nil {
	// 		xqb.Dump(err)
	// 		return
	// 	}

	// 	xqb.Dump(fmt.Sprintf("[%s] %s", elapsed, boundSql))
	// })
}

// forEachDialect iterates over all supported dialects and runs the test function
func forEachDialect(t *testing.T, test func(t *testing.T, dialect types.Driver)) {
	t.Helper()
	for _, dialect := range []types.Driver{types.DriverMySql, types.DriverPostgres} {
		t.Run(string(dialect), func(t *testing.T) {
			test(t, dialect)
		})
	}
}
