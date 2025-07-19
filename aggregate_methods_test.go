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

func Test_Upper(t *testing.T) {
	expr := xqb.Upper("name", "upper_name")
	assert.Equal(t, "UPPER(name) AS upper_name", expr.SQL)
}

func Test_Length(t *testing.T) {
	expr := xqb.Length("description", "desc_length")
	assert.Equal(t, "LENGTH(description) AS desc_length", expr.SQL)
}

func Test_Trim(t *testing.T) {
	expr := xqb.Trim("username", "trimmed_username")
	assert.Equal(t, "TRIM(username) AS trimmed_username", expr.SQL)
}

func Test_DateAdd(t *testing.T) {
	expr := xqb.DateAdd("created_at", "7", "DAY", "next_week")
	assert.Equal(t, "DATE_ADD(created_at, INTERVAL 7 DAY) AS next_week", expr.SQL)
}

func Test_DateSub(t *testing.T) {
	expr := xqb.DateSub("created_at", "1", "MONTH", "last_month")
	assert.Equal(t, "DATE_SUB(created_at, INTERVAL 1 MONTH) AS last_month", expr.SQL)
}

func Test_JSONFunc(t *testing.T) {
	expr := xqb.JSONFunc("JSON_UNQUOTE", []string{"data", "'$.email'"}, "email")
	assert.Equal(t, "JSON_UNQUOTE(data, '$.email') AS email", expr.SQL)
}

func Test_Substring(t *testing.T) {
	expr := xqb.Substring("bio", 1, 10, "short_bio")
	assert.Equal(t, "SUBSTRING(bio, 1, 10) AS short_bio", expr.SQL)
}

func Test_QueryBuilder_Upper_Length_Trim(t *testing.T) {
	qb := xqb.Table("users")
	sql, bindings, _ := qb.Select(
		xqb.Upper("name", "upper_name"),
		xqb.Length("bio", "bio_length"),
		xqb.Trim("username", "trimmed_username"),
	).Where("active", "=", true).ToSQL()
	assert.Equal(t, "SELECT UPPER(name) AS upper_name, LENGTH(bio) AS bio_length, TRIM(username) AS trimmed_username FROM users WHERE active = ?", sql)
	assert.Equal(t, []any{true}, bindings)
}

func Test_QueryBuilder_DateAdd_DateSub(t *testing.T) {
	qb := xqb.Table("events")
	sql, bindings, _ := qb.Select(
		xqb.DateAdd("event_date", "1", "DAY", "tomorrow"),
		xqb.DateSub("event_date", "7", "DAY", "last_week"),
	).Where("status", "=", "open").ToSQL()
	assert.Equal(t, "SELECT DATE_ADD(event_date, INTERVAL 1 DAY) AS tomorrow, DATE_SUB(event_date, INTERVAL 7 DAY) AS last_week FROM events WHERE status = ?", sql)
	assert.Equal(t, []any{"open"}, bindings)
}

func Test_QueryBuilder_JSONFunc_Substring(t *testing.T) {
	qb := xqb.Table("profiles")
	sql, bindings, _ := qb.Select(
		xqb.JSONFunc("JSON_UNQUOTE", []string{"data", "'$.email'"}, "email"),
		xqb.Substring("bio", 1, 20, "short_bio"),
	).Where("id", "=", 42).ToSQL()
	assert.Equal(t, "SELECT JSON_UNQUOTE(data, '$.email') AS email, SUBSTRING(bio, 1, 20) AS short_bio FROM profiles WHERE id = ?", sql)
	assert.Equal(t, []any{42}, bindings)
}

func Test_QueryBuilder_AggregateMethods(t *testing.T) {
	qb := xqb.Table("test_table")
	sql, bindings, _ := qb.Select(
		xqb.Count("id", "total_count"),
		xqb.Sum("amount", "total_amount"),
		xqb.Avg("score", "avg_score"),
		xqb.Min("age", "min_age"),
		xqb.Max("salary", "max_salary"),
		xqb.JsonExtract("data", "user.email", "user_email"),
		xqb.JSONFunc("JSON_UNQUOTE", []string{"data", "'$.phone'"}, "phone"),
		xqb.Math("price * quantity", "total_price"),
		xqb.Date("created_at", "created_date"),
		xqb.DateDiff("end_date", "start_date", "days_between"),
		xqb.DateAdd("created_at", "1", "DAY", "next_day"),
		xqb.DateSub("created_at", "1", "MONTH", "prev_month"),
		xqb.DateFormat("created_at", "%Y-%m-%d", "formatted_date"),
		xqb.Coalesce([]string{"middle_name", "'N/A'"}, "coalesced_name"),
		xqb.Concat([]string{"first_name", "' '", "last_name"}, "full_name"),
		xqb.Lower("email", "lower_email"),
		xqb.Upper("username", "upper_username"),
		xqb.Length("bio", "bio_length"),
		xqb.Trim("nickname", "trimmed_nickname"),
		xqb.Replace("title", "'foo'", "'bar'", "replaced_title"),
		xqb.Substring("description", 1, 10, "short_desc"),
	).Where("active", "=", true).ToSQL()

	expected := "SELECT COUNT(id) AS total_count, SUM(amount) AS total_amount, AVG(score) AS avg_score, MIN(age) AS min_age, MAX(salary) AS max_salary, JSON_EXTRACT(data, '$.user.email') AS user_email, JSON_UNQUOTE(data, '$.phone') AS phone, price * quantity AS total_price, DATE(created_at) AS created_date, DATEDIFF(end_date, start_date) AS days_between, DATE_ADD(created_at, INTERVAL 1 DAY) AS next_day, DATE_SUB(created_at, INTERVAL 1 MONTH) AS prev_month, DATE_FORMAT(created_at, '%Y-%m-%d') AS formatted_date, COALESCE(middle_name, 'N/A') AS coalesced_name, CONCAT(first_name, ' ', last_name) AS full_name, LOWER(email) AS lower_email, UPPER(username) AS upper_username, LENGTH(bio) AS bio_length, TRIM(nickname) AS trimmed_nickname, REPLACE(title, 'foo', 'bar') AS replaced_title, SUBSTRING(description, 1, 10) AS short_desc FROM test_table WHERE active = ?"
	assert.Equal(t, expected, sql)
	assert.Equal(t, []any{true}, bindings)

	// Also check Postgres SQL for DateFormat and JsonExtract
	dialectExprDateFormat := xqb.DateFormat("created_at", "%Y-%m-%d", "formatted_date")
	assert.Equal(t, "TO_CHAR(created_at, '%Y-%m-%d') AS formatted_date", dialectExprDateFormat.Dialects["postgres"].SQL)
	dialectExprJsonExtract := xqb.JsonExtract("data", "user.email", "user_email")
	assert.Equal(t, "data->'user'->>'email' AS user_email", dialectExprJsonExtract.Dialects["postgres"].SQL)
}

func Test_QueryBuilder_AggregateMethods_2(t *testing.T) {
	qb := xqb.Table("coverage_table")
	sql, bindings, _ := qb.Select(
		xqb.Count("*", ""),
		xqb.Count("id", "cnt"),
		xqb.Sum("amount", ""),
		xqb.Sum("amount", "total_amount"),
		xqb.Avg("score", ""),
		xqb.Avg("score", "avg_score"),
		xqb.Min("age", ""),
		xqb.Min("age", "min_age"),
		xqb.Max("salary", ""),
		xqb.Max("salary", "max_salary"),
		xqb.JsonExtract("data", "user.name", ""),
		xqb.JsonExtract("data", "user.name", "user_name"),
		xqb.JSONFunc("JSON_LENGTH", []string{"data", "'$.phones'"}, ""),
		xqb.JSONFunc("JSON_UNQUOTE", []string{"data", "'$.email'"}, "email"),
		xqb.Math("price * quantity", ""),
		xqb.Math("price * quantity + tax", "total_price"),
		xqb.Date("created_at", ""),
		xqb.Date("created_at", "created_date"),
		xqb.DateDiff("end_date", "start_date", ""),
		xqb.DateDiff("end_date", "start_date", "days_between"),
		xqb.DateAdd("created_at", "1", "DAY", ""),
		xqb.DateAdd("created_at", "1", "DAY", "next_day"),
		xqb.DateSub("created_at", "1", "MONTH", ""),
		xqb.DateSub("created_at", "1", "MONTH", "prev_month"),
		xqb.DateFormat("created_at", "%Y-%m-%d", ""),
		xqb.DateFormat("created_at", "%Y-%m-%d", "formatted_date"),
		xqb.Coalesce([]string{"middle_name", "'N/A'"}, ""),
		xqb.Coalesce([]string{"middle_name", "'N/A'"}, "coalesced_name"),
		xqb.Concat([]string{"first_name", "' '", "last_name"}, ""),
		xqb.Concat([]string{"first_name", "' '", "last_name"}, "full_name"),
		xqb.Lower("email", ""),
		xqb.Lower("email", "lower_email"),
		xqb.Upper("username", ""),
		xqb.Upper("username", "upper_username"),
		xqb.Length("bio", ""),
		xqb.Length("bio", "bio_length"),
		xqb.Trim("nickname", ""),
		xqb.Trim("nickname", "trimmed_nickname"),
		xqb.Replace("title", "'foo'", "'bar'", ""),
		xqb.Replace("title", "'foo'", "'bar'", "replaced_title"),
		xqb.Substring("description", 1, 10, ""),
		xqb.Substring("description", 1, 10, "short_desc"),
	).Where(
		xqb.Lower("status", ""), "=", "active",
	).GroupBy(
		xqb.Date("created_at", ""),
		xqb.Upper("region", ""),
	).Having(
		xqb.Sum("amount", ""), ">", 1000,
	).OrderBy(
		xqb.Length("bio", ""), "DESC",
	).Limit(5).Offset(10).ToSQL()

	expected := "SELECT COUNT(*), COUNT(id) AS cnt, SUM(amount), SUM(amount) AS total_amount, AVG(score), AVG(score) AS avg_score, MIN(age), MIN(age) AS min_age, MAX(salary), MAX(salary) AS max_salary, JSON_EXTRACT(data, '$.user.name'), JSON_EXTRACT(data, '$.user.name') AS user_name, JSON_LENGTH(data, '$.phones'), JSON_UNQUOTE(data, '$.email') AS email, price * quantity, price * quantity + tax AS total_price, DATE(created_at), DATE(created_at) AS created_date, DATEDIFF(end_date, start_date), DATEDIFF(end_date, start_date) AS days_between, DATE_ADD(created_at, INTERVAL 1 DAY), DATE_ADD(created_at, INTERVAL 1 DAY) AS next_day, DATE_SUB(created_at, INTERVAL 1 MONTH), DATE_SUB(created_at, INTERVAL 1 MONTH) AS prev_month, DATE_FORMAT(created_at, '%Y-%m-%d'), DATE_FORMAT(created_at, '%Y-%m-%d') AS formatted_date, COALESCE(middle_name, 'N/A'), COALESCE(middle_name, 'N/A') AS coalesced_name, CONCAT(first_name, ' ', last_name), CONCAT(first_name, ' ', last_name) AS full_name, LOWER(email), LOWER(email) AS lower_email, UPPER(username), UPPER(username) AS upper_username, LENGTH(bio), LENGTH(bio) AS bio_length, TRIM(nickname), TRIM(nickname) AS trimmed_nickname, REPLACE(title, 'foo', 'bar'), REPLACE(title, 'foo', 'bar') AS replaced_title, SUBSTRING(description, 1, 10), SUBSTRING(description, 1, 10) AS short_desc FROM coverage_table WHERE LOWER(status) = ? GROUP BY DATE(created_at), UPPER(region) HAVING SUM(amount) > ? ORDER BY LENGTH(bio) DESC LIMIT 5 OFFSET 10"
	assert.Equal(t, expected, sql)
	assert.Equal(t, []any{"active", 1000}, bindings)
}
