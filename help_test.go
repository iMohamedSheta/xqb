package xqb_test

import (
	"testing"

	dialects "github.com/iMohamedSheta/xqb/grammar"
)

// forEachDialect iterates over all supported dialects and runs the test function
func forEachDialect(t *testing.T, test func(t *testing.T, dialect dialects.Driver)) {
	t.Helper()
	for _, dialect := range []dialects.Driver{dialects.DriverMySQL, dialects.DriverPostgres} {
		t.Run(string(dialect), func(t *testing.T) {
			test(t, dialect)
		})
	}
}
