package xqb_test

import (
	"testing"

	"github.com/iMohamedSheta/xqb"
	"github.com/iMohamedSheta/xqb/dialects"
	"github.com/iMohamedSheta/xqb/shared/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDialects_SelectBasic(t *testing.T) {
	tests := []struct {
		dialect  types.Dialect
		expected string
		binding  string // placeholder style check
	}{
		{types.DialectSQLite, `SELECT "id", "name" FROM "users" WHERE "id" = ?`, "?"},
		{types.DialectMariaDB, "SELECT `id`, `name` FROM `users` WHERE `id` = ?", "?"},
		{types.DialectSQLServer, `SELECT [id], [name] FROM [users] WHERE [id] = @p1`, "@p1"},
		{types.DialectMySql, "SELECT `id`, `name` FROM `users` WHERE `id` = ?", "?"},
		{types.DialectPostgres, `SELECT "id", "name" FROM "users" WHERE "id" = $1`, "$1"},
	}

	for _, tt := range tests {
		t.Run(string(tt.dialect), func(t *testing.T) {
			sql, bindings, err := xqb.Table("users").SetDialect(tt.dialect).
				Select("id", "name").Where("id", "=", 1).ToSql()
			require.NoError(t, err)
			assert.Equal(t, tt.expected, sql)
			assert.Equal(t, []any{1}, bindings)
			assert.Contains(t, sql, tt.binding)
		})
	}
}

func TestNewDialects_LimitOffset(t *testing.T) {
	// SQLite: LIMIT/OFFSET like MySQL
	sql, _, err := xqb.Table("users").SetDialect(types.DialectSQLite).
		Select("id").Limit(10).Offset(20).ToSql()
	require.NoError(t, err)
	assert.Equal(t, `SELECT "id" FROM "users" LIMIT 10 OFFSET 20`, sql)

	// MariaDB: same as MySQL
	sql, _, err = xqb.Table("users").SetDialect(types.DialectMariaDB).
		Select("id").Limit(10).Offset(20).ToSql()
	require.NoError(t, err)
	assert.Equal(t, "SELECT `id` FROM `users` LIMIT 10 OFFSET 20", sql)

	// SQLServer: TOP when limit without offset
	sql, _, err = xqb.Table("users").SetDialect(types.DialectSQLServer).
		Select("id").Limit(10).ToSql()
	require.NoError(t, err)
	assert.Equal(t, `SELECT TOP 10 [id] FROM [users]`, sql)

	// SQLServer: OFFSET/FETCH with ORDER BY
	sql, bindings, err := xqb.Table("users").SetDialect(types.DialectSQLServer).
		Select("id").OrderBy("name", "ASC").Limit(10).Offset(20).ToSql()
	require.NoError(t, err)
	assert.Equal(t, `SELECT [id] FROM [users] ORDER BY [name] ASC OFFSET 20 ROWS FETCH NEXT 10 ROWS ONLY`, sql)
	assert.Empty(t, bindings)

	// SQLServer: OFFSET without ORDER BY injects ORDER BY (SELECT 0)
	sql, _, err = xqb.Table("users").SetDialect(types.DialectSQLServer).
		Select("id").Offset(5).ToSql()
	require.NoError(t, err)
	assert.Equal(t, `SELECT [id] FROM [users] ORDER BY (SELECT 0) OFFSET 5 ROWS`, sql)
}

func TestNewDialects_Locks(t *testing.T) {
	// SQLite: locks are ignored
	sql, _, err := xqb.Table("users").SetDialect(types.DialectSQLite).LockForUpdate().ToSql()
	require.NoError(t, err)
	assert.Equal(t, `SELECT * FROM "users"`, sql)

	// SQLServer: WITH hints on FROM
	sql, _, err = xqb.Table("users").SetDialect(types.DialectSQLServer).LockForUpdate().ToSql()
	require.NoError(t, err)
	assert.Equal(t, `SELECT * FROM [users] WITH (ROWLOCK,UPDLOCK,HOLDLOCK)`, sql)

	sql, _, err = xqb.Table("users").SetDialect(types.DialectSQLServer).SharedLock().ToSql()
	require.NoError(t, err)
	assert.Equal(t, `SELECT * FROM [users] WITH (ROWLOCK,HOLDLOCK)`, sql)

	// MariaDB: MySQL-style locks
	sql, _, err = xqb.Table("users").SetDialect(types.DialectMariaDB).LockForUpdate().ToSql()
	require.NoError(t, err)
	assert.Equal(t, "SELECT * FROM `users` FOR UPDATE", sql)

	sql, _, err = xqb.Table("users").SetDialect(types.DialectMariaDB).SharedLock().ToSql()
	require.NoError(t, err)
	assert.Equal(t, "SELECT * FROM `users` LOCK IN SHARE MODE", sql)
}

func TestNewDialects_Insert(t *testing.T) {
	vals := []map[string]any{{"name": "John", "email": "john@example.com"}}

	// SQLite basic
	sql, bindings, err := xqb.Table("users").SetDialect(types.DialectSQLite).InsertSql(vals)
	require.NoError(t, err)
	assert.Equal(t, `INSERT INTO "users" ("email", "name") VALUES (?, ?)`, sql)
	assert.Equal(t, []any{"john@example.com", "John"}, bindings)

	// SQLite upsert (ON CONFLICT)
	sql, _, err = xqb.Table("users").SetDialect(types.DialectSQLite).
		UpsertSql(vals, []string{"email"}, []string{"name"})
	require.NoError(t, err)
	assert.Equal(t, `INSERT INTO "users" ("email", "name") VALUES (?, ?) ON CONFLICT ("email") DO UPDATE SET "name" = EXCLUDED."name"`, sql)

	// SQLite returning
	sql, _, err = xqb.Table("users").SetDialect(types.DialectSQLite).InsertGetIdSql(vals)
	require.NoError(t, err)
	assert.Contains(t, sql, "RETURNING id")

	// SQLServer basic (@p placeholders)
	sql, bindings, err = xqb.Table("users").SetDialect(types.DialectSQLServer).InsertSql(vals)
	require.NoError(t, err)
	assert.Equal(t, `INSERT INTO [users] ([email], [name]) VALUES (@p1, @p2)`, sql)
	assert.Equal(t, []any{"john@example.com", "John"}, bindings)

	// SQLServer returning (OUTPUT INSERTED)
	sql, _, err = xqb.Table("users").SetDialect(types.DialectSQLServer).InsertGetIdSql(vals)
	require.NoError(t, err)
	assert.Equal(t, `INSERT INTO [users] ([email], [name]) OUTPUT INSERTED.[id] VALUES (@p1, @p2)`, sql)

	// SQLServer upsert unsupported
	_, _, err = xqb.Table("users").SetDialect(types.DialectSQLServer).
		UpsertSql(vals, []string{"email"}, []string{"name"})
	require.Error(t, err)

	// MariaDB basic + upsert (ON DUPLICATE KEY)
	sql, _, err = xqb.Table("users").SetDialect(types.DialectMariaDB).InsertSql(vals)
	require.NoError(t, err)
	assert.Equal(t, "INSERT INTO `users` (`email`, `name`) VALUES (?, ?)", sql)

	sql, _, err = xqb.Table("users").SetDialect(types.DialectMariaDB).
		UpsertSql(vals, []string{"email"}, []string{"name"})
	require.NoError(t, err)
	assert.Equal(t, "INSERT INTO `users` (`email`, `name`) VALUES (?, ?) ON DUPLICATE KEY UPDATE `name` = VALUES(`name`)", sql)
}

func TestNewDialects_UpdateDelete(t *testing.T) {
	// SQLite update
	sql, bindings, err := xqb.Table("users").SetDialect(types.DialectSQLite).
		Where("id", "=", 1).UpdateSql(map[string]any{"name": "Jane"})
	require.NoError(t, err)
	assert.Equal(t, `UPDATE "users" SET "name" = ? WHERE "id" = ?`, sql)
	assert.Equal(t, []any{"Jane", 1}, bindings)

	// SQLite update with LIMIT rejected
	_, _, err = xqb.Table("users").SetDialect(types.DialectSQLite).
		Where("id", "=", 1).Limit(1).UpdateSql(map[string]any{"name": "Jane"})
	require.Error(t, err)

	// SQLServer update
	sql, bindings, err = xqb.Table("users").SetDialect(types.DialectSQLServer).
		Where("id", "=", 1).UpdateSql(map[string]any{"name": "Jane"})
	require.NoError(t, err)
	assert.Equal(t, `UPDATE [users] SET [name] = @p1 WHERE [id] = @p2`, sql)
	assert.Equal(t, []any{"Jane", 1}, bindings)

	// SQLServer update with JOIN uses FROM
	sql, _, err = xqb.Table("users").SetDialect(types.DialectSQLServer).
		Join("orders", "users.id = orders.user_id").
		Where("users.id", "=", 1).
		UpdateSql(map[string]any{"status": "x"})
	require.NoError(t, err)
	assert.Contains(t, sql, "FROM [users]")
	assert.Contains(t, sql, "JOIN [orders]")

	// MariaDB update with JOIN + LIMIT allowed (MySQL-style)
	sql, _, err = xqb.Table("users").SetDialect(types.DialectMariaDB).
		Where("status", "=", "pending").Limit(10).
		UpdateSql(map[string]any{"status": "verified"})
	require.NoError(t, err)
	assert.Equal(t, "UPDATE `users` SET `status` = ? WHERE `status` = ? LIMIT 10", sql)

	// SQLite delete
	sql, bindings, err = xqb.Table("users").SetDialect(types.DialectSQLite).
		Where("id", "=", 1).DeleteSql()
	require.NoError(t, err)
	assert.Equal(t, `DELETE FROM "users" WHERE "id" = ?`, sql)
	assert.Equal(t, []any{1}, bindings)

	// SQLite delete with LIMIT rejected
	_, _, err = xqb.Table("users").SetDialect(types.DialectSQLite).
		Where("id", "=", 1).Limit(1).DeleteSql()
	require.Error(t, err)

	// SQLServer delete with TOP
	sql, _, err = xqb.Table("users").SetDialect(types.DialectSQLServer).
		Where("status", "=", "pending").Limit(10).DeleteSql()
	require.NoError(t, err)
	assert.Equal(t, `DELETE TOP (10) FROM [users] WHERE [status] = @p1`, sql)

	// MariaDB delete with LIMIT allowed
	sql, _, err = xqb.Table("users").SetDialect(types.DialectMariaDB).
		Where("status", "=", "pending").Limit(10).DeleteSql()
	require.NoError(t, err)
	assert.Equal(t, "DELETE FROM `users` WHERE `status` = ? LIMIT 10", sql)
}

func TestNewDialects_UnionAndJoins(t *testing.T) {
	// INTERSECT allowed on sqlite / mariadb / sqlserver (unlike mysql)
	for _, d := range []types.Dialect{types.DialectSQLite, types.DialectMariaDB, types.DialectSQLServer} {
		sub := xqb.Table("admins").SetDialect(d).Select("id")
		_, _, err := xqb.Table("users").SetDialect(d).Select("id").IntersectUnion(sub).ToSql()
		require.NoError(t, err, "dialect %s should support INTERSECT", d)

		sub2 := xqb.Table("banned").SetDialect(d).Select("id")
		_, _, err = xqb.Table("users").SetDialect(d).Select("id").ExceptUnion(sub2).ToSql()
		require.NoError(t, err, "dialect %s should support EXCEPT", d)
	}

	// FULL JOIN: allowed on sqlserver, rejected on sqlite/mariadb
	_, _, err := xqb.Table("users").SetDialect(types.DialectSQLServer).
		FullJoin("sessions", "users.id = sessions.user_id").ToSql()
	require.NoError(t, err)

	_, _, err = xqb.Table("users").SetDialect(types.DialectSQLite).
		FullJoin("sessions", "users.id = sessions.user_id").ToSql()
	require.Error(t, err)

	_, _, err = xqb.Table("users").SetDialect(types.DialectMariaDB).
		FullJoin("sessions", "users.id = sessions.user_id").ToSql()
	require.Error(t, err)
}

func TestNewDialects_InjectBindings(t *testing.T) {
	// SQLite uses ?
	out, err := xqb.InjectBindings(types.DialectSQLite, `SELECT * FROM "users" WHERE "id" = ?`, []any{1})
	require.NoError(t, err)
	assert.Equal(t, `SELECT * FROM "users" WHERE "id" = 1`, out)

	// MariaDB uses ?
	out, err = xqb.InjectBindings(types.DialectMariaDB, "SELECT * FROM `users` WHERE `id` = ?", []any{1})
	require.NoError(t, err)
	assert.Equal(t, "SELECT * FROM `users` WHERE `id` = 1", out)

	// SQLServer uses @pN
	out, err = xqb.InjectBindings(types.DialectSQLServer, `SELECT * FROM [users] WHERE [id] = @p1 AND [age] > @p2`, []any{1, 18})
	require.NoError(t, err)
	assert.Equal(t, `SELECT * FROM [users] WHERE [id] = 1 AND [age] > 18`, out)

	// ToSqlView end-to-end per dialect
	view, err := xqb.Table("users").SetDialect(types.DialectSQLServer).Where("id", "=", 1).ToSqlView()
	require.NoError(t, err)
	assert.Equal(t, `SELECT * FROM [users] WHERE [id] = 1`, view)

	view, err = xqb.Table("users").SetDialect(types.DialectSQLite).Where("id", "=", 1).ToSqlView()
	require.NoError(t, err)
	assert.Equal(t, `SELECT * FROM "users" WHERE "id" = 1`, view)
}

func TestNewDialects_AliasesAndRegistry(t *testing.T) {
	// Normalize aliases
	assert.Equal(t, types.DialectPostgres, types.Dialect("pgsql").Normalize())
	assert.Equal(t, types.DialectPostgres, types.Dialect("PostgreSQL").Normalize())
	assert.Equal(t, types.DialectSQLite, types.Dialect("sqlite3").Normalize())
	assert.Equal(t, types.DialectSQLServer, types.Dialect("sqlsrv").Normalize())
	assert.Equal(t, types.DialectSQLServer, types.Dialect("mssql").Normalize())
	assert.Equal(t, types.DialectMySql, types.Dialect("mysql").Normalize())
	assert.Equal(t, types.DialectMariaDB, types.Dialect("mariadb").Normalize())

	// Connection-level mapping
	assert.Equal(t, types.DialectPostgres, xqb.DialectPgsql.MappedDialect())
	assert.Equal(t, types.DialectSQLServer, xqb.DialectSqlSrv.MappedDialect())
	assert.Equal(t, types.DialectSQLite, xqb.DialectSQLite.MappedDialect())
	assert.Equal(t, types.DialectMariaDB, xqb.DialectMariaDB.MappedDialect())

	// GetDialect registry returns the right implementation
	assert.Equal(t, types.DialectSQLite, dialects.GetDialect("sqlite").Getdialect())
	assert.Equal(t, types.DialectSQLite, dialects.GetDialect("sqlite3").Getdialect())
	assert.Equal(t, types.DialectSQLServer, dialects.GetDialect("sqlsrv").Getdialect())
	assert.Equal(t, types.DialectSQLServer, dialects.GetDialect("mssql").Getdialect())
	assert.Equal(t, types.DialectPostgres, dialects.GetDialect("pgsql").Getdialect())
	assert.Equal(t, types.DialectMariaDB, dialects.GetDialect("mariadb").Getdialect())
	assert.Equal(t, types.DialectMySql, dialects.GetDialect("mysql").Getdialect())

	// SetDialect accepts aliases
	sql, _, err := xqb.Table("users").SetDialect(types.Dialect("pgsql")).Select("id").ToSql()
	require.NoError(t, err)
	assert.Contains(t, sql, `"id"`)

	sql, _, err = xqb.Table("users").SetDialect(types.Dialect("sqlsrv")).Select("id").ToSql()
	require.NoError(t, err)
	assert.Contains(t, sql, `[id]`)
}

func TestNewDialects_WrapAndHelpers(t *testing.T) {
	// SQLServer bracket wrapping
	d := dialects.GetDialect(types.DialectSQLServer)
	w, ok := d.(interface{ Wrap(string) string })
	require.True(t, ok)
	assert.Equal(t, "[users].[id]", w.Wrap("users.id"))
	assert.Equal(t, "[users] AS [u]", w.Wrap("users as u"))
	assert.Equal(t, "*", w.Wrap("*"))

	// Dialect-aware helpers resolve per new dialect
	assert.Equal(t, "json_extract(data, '$.a.b') AS x", mustToSQL(t, xqb.JsonExtract("data", "a.b", "x"), "sqlite"))
	assert.Equal(t, "JSON_VALUE(data, '$.a.b') AS x", mustToSQL(t, xqb.JsonExtract("data", "a.b", "x"), "sqlserver"))
	assert.Equal(t, "JSON_EXTRACT(data, '$.a.b') AS x", mustToSQL(t, xqb.JsonExtract("data", "a.b", "x"), "mariadb"))
	assert.Equal(t, "strftime('%Y-%m-%d', created_at) AS d", mustToSQL(t, xqb.DateFormat("created_at", "%Y-%m-%d", "d"), "sqlite"))
	assert.Equal(t, "FORMAT(created_at, '%Y-%m-%d') AS d", mustToSQL(t, xqb.DateFormat("created_at", "%Y-%m-%d", "d"), "sqlserver"))
}

func mustToSQL(t *testing.T, e types.ExpressionInterface, dialect string) string {
	t.Helper()
	sql, _, err := e.ToSql(dialect)
	require.NoError(t, err)
	return sql
}
