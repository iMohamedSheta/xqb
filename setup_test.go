package xqb_test

import (
	"os"
	"testing"

	"github.com/iMohamedSheta/xqb"
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

	// Default connection settings for testing it will allow creating sql but not executing it
	xqb.AddConnection(&xqb.Connection{
		Name:    "default", // Default connection name
		Dialect: types.DialectMySql,
		DB:      nil,
	})
}

// forEachDialect iterates over all supported dialects and runs the test function
func forEachDialect(t *testing.T, test func(t *testing.T, dialect types.Dialect)) {
	t.Helper()
	for _, dialect := range []types.Dialect{types.DialectMySql, types.DialectPostgres} {
		t.Run(string(dialect), func(t *testing.T) {
			test(t, dialect)
		})
	}
}
