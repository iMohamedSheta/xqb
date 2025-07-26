package xqb_test

import (
	"testing"

	"github.com/iMohamedSheta/xqb/shared/types"
)

// forEachDialect iterates over all supported dialects and runs the test function
func forEachDialect(t *testing.T, test func(t *testing.T, dialect types.Driver)) {
	t.Helper()
	for _, dialect := range []types.Driver{types.DriverMySql, types.DriverPostgres} {
		t.Run(string(dialect), func(t *testing.T) {
			test(t, dialect)
		})
	}
}
