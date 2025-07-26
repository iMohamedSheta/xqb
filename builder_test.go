package xqb_test

import (
	"testing"

	"github.com/iMohamedSheta/xqb"
)

// TestOnBeforeQueryHook tests the OnBeforeQuery and OnAfterQuery hooks
func TestOnBeforeQueryHook(t *testing.T) {
	beforeCalled := false
	afterCalled := false

	xqb.DefaultSettings().OnBeforeQuery(func(qb *xqb.QueryBuilder) {
		beforeCalled = true
	})

	xqb.DefaultSettings().OnAfterQuery(func(query *xqb.QueryExecuted) {
		afterCalled = true
	})

	Setup()

	qb := xqb.Query().Table("users").Where("id", "=", 1)

	_, _, err := qb.ToSql()
	if err != nil {
		t.Fatalf("ToSql failed: %v", err)
	}

	if !beforeCalled {
		t.Errorf("OnBeforeQuery hook was not called")
	}

	if !afterCalled {
		t.Errorf("OnAfterQuery hook was not called")
	}
}
