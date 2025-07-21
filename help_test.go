package xqb_test

import (
	"testing"

	"github.com/iMohamedSheta/xqb/grammar"
)

// forEachDialect iterates over all supported dialects and runs the test function
func forEachDialect(t *testing.T, test func(t *testing.T, dialect grammar.Driver)) {
	t.Helper()
	for _, dialect := range []grammar.Driver{grammar.DriverMySQL, grammar.DriverPostgres} {
		t.Run(string(dialect), func(t *testing.T) {
			test(t, dialect)
		})
	}
}
