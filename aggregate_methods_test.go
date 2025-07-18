package xqb_test

import (
	"testing"

	"github.com/iMohamedSheta/xqb"
	"github.com/stretchr/testify/assert"
)

func Test_Aggregate(t *testing.T) {
	qb := xqb.Table("users")
	sql, binding, _ := qb.Select(
		xqb.Sum("price", "total_price"),
		"username",
		"email",
	).Where("id", "=", 15).
		Limit(1).
		ToSQL()

	assert.Equal(t, "SELECT SUM(price) AS total_price, username, email FROM users WHERE id = ? LIMIT 1", sql)
	assert.Equal(t, []any{15}, binding)
}

func Test_DialectExpr(t *testing.T) {
	qb := xqb.Table("users")
	sql, binding, _ := qb.Select(
		xqb.Sum("price", "total_price"),
		xqb.DateFormat("created_at", "%Y-%m-%d", "created_at"),
		"username",
		"email",
	).Where("id", "=", 15).
		Limit(1).
		ToSQL()

	assert.Equal(t, "SELECT SUM(price) AS total_price, DATE_FORMAT(created_at, '%Y-%m-%d') AS created_at, username, email FROM users WHERE id = ? LIMIT 1", sql)
	assert.Equal(t, []any{15}, binding)
}

func Test_CountExpression(t *testing.T) {
	expr := xqb.Count("id", "total_users")
	assert.Equal(t, "COUNT(id) AS total_users", expr.SQL)
}

func Test_JsonExtract_MySQL(t *testing.T) {
	expr := xqb.JsonExtract("data", "user.name", "username")
	sql := expr.Dialects["mysql"].SQL
	assert.Equal(t, "JSON_EXTRACT(data, '$.user.name') AS username", sql)
}

func Test_JsonExtract_Postgres(t *testing.T) {
	expr := xqb.JsonExtract("data", "user.name", "username")
	sql := expr.Dialects["postgres"].SQL
	assert.Equal(t, "data->'user'->>'name' AS username", sql)
}

func Test_DateFunctions(t *testing.T) {
	assert.Equal(t, "DATE(created_at) AS date_only", xqb.Date("created_at", "date_only").SQL)
	assert.Equal(t, "DATEDIFF(end_date, start_date) AS diff", xqb.DateDiff("end_date", "start_date", "diff").SQL)
}

func Test_StringFunctions(t *testing.T) {
	assert.Equal(t, "LOWER(name) AS lower_name", xqb.Lower("name", "lower_name").SQL)
	assert.Equal(t, "REPLACE(title, 'foo', 'bar') AS updated", xqb.Replace("title", "'foo'", "'bar'", "updated").SQL)
}

func Test_ConcatWithBindings(t *testing.T) {
	expr := xqb.Concat([]string{"first_name", "' '", "last_name"}, "full_name")
	assert.Equal(t, "CONCAT(first_name, ' ', last_name) AS full_name", expr.SQL)
}

func Test_MathExpression(t *testing.T) {
	expr := xqb.Math("price * quantity", "total")
	assert.Equal(t, "price * quantity AS total", expr.SQL)
}

func Test_Coalesce(t *testing.T) {
	expr := xqb.Coalesce([]string{"middle_name", "'N/A'"}, "coalesced_name")
	assert.Equal(t, "COALESCE(middle_name, 'N/A') AS coalesced_name", expr.SQL)
}

func Test_QueryBuilder_Locks(t *testing.T) {
	qb := xqb.Table("users").LockForUpdate()
	assert.True(t, qb.GetData().IsLockedForUpdate)

	qb = xqb.Table("users").SharedLock()
	assert.True(t, qb.GetData().IsInSharedLock)
}
