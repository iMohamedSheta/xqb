package integration

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Seed: Alice(30, 95.5, active), Bob(25, 80.0, active),
// Carol(35, 70.25, inactive), Dave(NULL, NULL, NULL).

func TestAggregate_Count(t *testing.T) {
	forEachDB(t, func(t *testing.T, conn string) {
		resetTables(t, conn)

		n, err := QB(conn, "xqb_users").Count("*")
		require.NoError(t, err)
		assert.Equal(t, int64(4), n)

		n, err = QB(conn, "xqb_users").Where("status", "=", "active").Count("*")
		require.NoError(t, err)
		assert.Equal(t, int64(2), n)

		// COUNT(column) skips NULLs.
		n, err = QB(conn, "xqb_users").Count("age")
		require.NoError(t, err)
		assert.Equal(t, int64(3), n)
	})
}

func TestAggregate_SumAvgMinMax(t *testing.T) {
	forEachDB(t, func(t *testing.T, conn string) {
		resetTables(t, conn)

		sum, err := QB(conn, "xqb_users").Sum("age")
		require.NoError(t, err)
		assert.InDelta(t, 90.0, sum, 0.001, "[%s] sum(age)", conn)

		avg, err := QB(conn, "xqb_users").Avg("age")
		require.NoError(t, err)
		assert.InDelta(t, 30.0, avg, 0.001, "[%s] avg(age)", conn)

		mn, err := QB(conn, "xqb_users").Min("age")
		require.NoError(t, err)
		assert.InDelta(t, 25.0, mn, 0.001)

		mx, err := QB(conn, "xqb_users").Max("age")
		require.NoError(t, err)
		assert.InDelta(t, 35.0, mx, 0.001)

		// Score column exercises float handling on every driver.
		sumScore, err := QB(conn, "xqb_users").Sum("score")
		require.NoError(t, err)
		assert.InDelta(t, 245.75, sumScore, 0.001, "[%s] sum(score)", conn)
	})
}

func TestAggregate_WithWhere(t *testing.T) {
	forEachDB(t, func(t *testing.T, conn string) {
		resetTables(t, conn)

		n, err := QB(conn, "xqb_users").Where("age", ">", 26).Count("*")
		require.NoError(t, err)
		assert.Equal(t, int64(2), n)

		mx, err := QB(conn, "xqb_users").Where("status", "=", "active").Max("score")
		require.NoError(t, err)
		assert.InDelta(t, 95.5, mx, 0.001)
	})
}
