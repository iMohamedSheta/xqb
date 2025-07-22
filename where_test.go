package xqb_test

import (
	"testing"

	"github.com/iMohamedSheta/xqb"
	"github.com/iMohamedSheta/xqb/shared/types"
	"github.com/stretchr/testify/assert"
)

// .Where - Test cases for the Where clause in QueryBuilder

func Test_Where_Subquery_1(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("users").SetDialect(dialect)
		subQuery := xqb.Table("orders").Select("user_id").Where("status", "=", "active")
		sql, bindings, err := qb.Where("id", "IN", subQuery).ToSQL()

		expectedSql := map[types.Driver]string{
			types.DriverMySQL:    "SELECT * FROM users WHERE id IN (SELECT user_id FROM orders WHERE status = ?)",
			types.DriverPostgres: "SELECT * FROM users WHERE id IN (SELECT user_id FROM orders WHERE status = $1)",
		}
		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, []any{"active"}, bindings)
		assert.NoError(t, err)
	})
}

func Test_Where_Subquery_2(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("users").SetDialect(dialect)
		subQuery := xqb.Table("admins").Select("user_id").Where("role", "=", "superadmin").AddSelect("id").Latest("id")

		sql, bindings, err := qb.Where("id", "IN", subQuery).ToSQL()
		expectedSql := map[types.Driver]string{
			types.DriverMySQL:    "SELECT * FROM users WHERE id IN (SELECT user_id, id FROM admins WHERE role = ? ORDER BY id DESC)",
			types.DriverPostgres: "SELECT * FROM users WHERE id IN (SELECT user_id, id FROM admins WHERE role = $1 ORDER BY id DESC)",
		}

		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, []any{"superadmin"}, bindings)
		assert.NoError(t, err)
	})
}

func Test_Where_Subquery_3(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("admins").SetDialect(dialect)
		subQuery := xqb.Table("users").
			Join("orders", "users.id = orders.user_id").
			Where("orders.status", "=", "paid").
			Select("users.id", "users.name", "orders.id AS order_id")

		sql, bindings, err := qb.Where("id", "IN", subQuery).ToSQL()
		expectedSql := map[types.Driver]string{
			types.DriverMySQL:    "SELECT * FROM admins WHERE id IN (SELECT users.id, users.name, orders.id AS order_id FROM users JOIN orders ON users.id = orders.user_id WHERE orders.status = ?)",
			types.DriverPostgres: "SELECT * FROM admins WHERE id IN (SELECT users.id, users.name, orders.id AS order_id FROM users JOIN orders ON users.id = orders.user_id WHERE orders.status = $1)",
		}
		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, []any{"paid"}, bindings)
		assert.NoError(t, err)
	})
}

func Test_Where_WithRaw_CaseExpression(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("users").SetDialect(dialect)
		sql, bindings, err := qb.Where(xqb.Raw("CASE WHEN status = 'active' THEN 1 ELSE 0 END"), "=", 1).ToSQL()

		expectedSql := map[types.Driver]string{
			types.DriverMySQL:    "SELECT * FROM users WHERE CASE WHEN status = 'active' THEN 1 ELSE 0 END = ?",
			types.DriverPostgres: "SELECT * FROM users WHERE CASE WHEN status = 'active' THEN 1 ELSE 0 END = $1",
		}
		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, []any{1}, bindings)
		assert.NoError(t, err)
	})
}

func Test_Where_WithRaw_1(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("users").SetDialect(dialect)
		sql, bindings, err := qb.Join("orders", "users.id = orders.user_id").Where(xqb.Raw("CASE WHEN status = 'active' THEN 1 ELSE 0 END"), "=", 1).ToSQL()

		expectedSql := map[types.Driver]string{
			types.DriverMySQL:    "SELECT * FROM users JOIN orders ON users.id = orders.user_id WHERE CASE WHEN status = 'active' THEN 1 ELSE 0 END = ?",
			types.DriverPostgres: "SELECT * FROM users JOIN orders ON users.id = orders.user_id WHERE CASE WHEN status = 'active' THEN 1 ELSE 0 END = $1",
		}
		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, []any{1}, bindings)
		assert.NoError(t, err)
	})
}

func Test_OrWhere_SubQuery_1(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("users").SetDialect(dialect)
		subQuery := xqb.Table("orders").Select("user_id").Where("status", "=", "active").Latest("id")
		sql, bindings, err := qb.OrWhere("id", "IN", subQuery).ToSQL()

		expectedSql := map[types.Driver]string{
			types.DriverMySQL:    "SELECT * FROM users WHERE id IN (SELECT user_id FROM orders WHERE status = ? ORDER BY id DESC)",
			types.DriverPostgres: "SELECT * FROM users WHERE id IN (SELECT user_id FROM orders WHERE status = $1 ORDER BY id DESC)",
		}
		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, []any{"active"}, bindings)
		assert.NoError(t, err)
	})
}

func Test_OrWhere_SubQuery_2(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("users").SetDialect(dialect)
		subQuery := xqb.Table("orders").Join("admins", "users.id = admins.user_id").Select("user_id").Where("role", "=", "superadmin").Latest("id")
		sql, bindings, err := qb.OrWhere("id", "IN", subQuery).ToSQL()

		expectedSql := map[types.Driver]string{
			types.DriverMySQL:    "SELECT * FROM users WHERE id IN (SELECT user_id FROM orders JOIN admins ON users.id = admins.user_id WHERE role = ? ORDER BY id DESC)",
			types.DriverPostgres: "SELECT * FROM users WHERE id IN (SELECT user_id FROM orders JOIN admins ON users.id = admins.user_id WHERE role = $1 ORDER BY id DESC)",
		}
		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, []any{"superadmin"}, bindings)
		assert.NoError(t, err)
	})
}

func Test_OrWhere_Raw_1(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("users").SetDialect(dialect)
		sql, bindings, err := qb.OrWhere(xqb.Raw("CASE WHEN status IN ('active', 'pending') THEN 1 ELSE 0 END"), "=", 1).ToSQL()

		expectedSql := map[types.Driver]string{
			types.DriverMySQL:    "SELECT * FROM users WHERE CASE WHEN status IN ('active', 'pending') THEN 1 ELSE 0 END = ?",
			types.DriverPostgres: "SELECT * FROM users WHERE CASE WHEN status IN ('active', 'pending') THEN 1 ELSE 0 END = $1",
		}
		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, []any{1}, bindings)
		assert.NoError(t, err)
	})

}

func Test_OrWhere_Raw_2(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("users").SetDialect(dialect)
		sql, bindings, err := qb.OrWhere(xqb.Raw("CASE WHEN status IN ('active', 'pending') THEN 1 ELSE 0 END"), "=", 1).Join("orders", "users.id = orders.user_id").ToSQL()

		expectedSql := map[types.Driver]string{
			types.DriverMySQL:    "SELECT * FROM users JOIN orders ON users.id = orders.user_id WHERE CASE WHEN status IN ('active', 'pending') THEN 1 ELSE 0 END = ?",
			types.DriverPostgres: "SELECT * FROM users JOIN orders ON users.id = orders.user_id WHERE CASE WHEN status IN ('active', 'pending') THEN 1 ELSE 0 END = $1",
		}
		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, []any{1}, bindings)
		assert.NoError(t, err)
	})
}

func Test_WhereNull_With_OrWhereNotNull(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("users").SetDialect(dialect)
		sql, bindings, err := qb.Where("id", "=", 1).WhereGroup(func(qb *xqb.QueryBuilder) {
			qb.OrWhereNull("deleted_at").OrWhereNotNull("disabled_at")
		}).ToSQL()

		expectedSql := map[types.Driver]string{
			types.DriverMySQL:    "SELECT * FROM users WHERE id = ? AND (deleted_at IS NULL OR disabled_at IS NOT NULL)",
			types.DriverPostgres: "SELECT * FROM users WHERE id = $1 AND (deleted_at IS NULL OR disabled_at IS NOT NULL)",
		}
		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, 1, len(bindings))
		assert.Equal(t, []any{1}, bindings)
		assert.NoError(t, err)
	})
}

func Test_WhereNull_With_Grouping(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("users").SetDialect(dialect)
		sql, bindings, err := qb.Where("id", "=", 1).WhereGroup(func(qb *xqb.QueryBuilder) {
			qb.WhereNull("deleted_at").OrWhereNull("disabled_at")
		}).ToSQL()

		expectedSql := map[types.Driver]string{
			types.DriverMySQL:    "SELECT * FROM users WHERE id = ? AND (deleted_at IS NULL OR disabled_at IS NULL)",
			types.DriverPostgres: "SELECT * FROM users WHERE id = $1 AND (deleted_at IS NULL OR disabled_at IS NULL)",
		}
		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, 1, len(bindings))
		assert.Equal(t, []any{1}, bindings)
		assert.NoError(t, err)
	})
}

func Test_WhereIn_normal(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("users").SetDialect(dialect)
		sql, bindings, err := qb.WhereIn("id", []any{1, 2, 3}).ToSQL()

		expectedSql := map[types.Driver]string{
			types.DriverMySQL:    "SELECT * FROM users WHERE id IN (?, ?, ?)",
			types.DriverPostgres: "SELECT * FROM users WHERE id IN ($1, $2, $3)",
		}

		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, []any{1, 2, 3}, bindings)
		assert.NoError(t, err)
	})
}

func Test_WhereIn_With_Raw(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("users").SetDialect(dialect)
		sql, bindings, err := qb.WhereIn("id", []any{xqb.Raw("? UNION ?", 1, 2)}).ToSQL()

		expectedSql := map[types.Driver]string{
			types.DriverMySQL:    "SELECT * FROM users WHERE id IN (? UNION ?)",
			types.DriverPostgres: "SELECT * FROM users WHERE id IN ($1 UNION $2)",
		}

		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, []any{1, 2}, bindings)
		assert.NoError(t, err)
	})

}

func Test_WhereIn_With_Raw_2(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("users").SetDialect(dialect)
		sql, bindings, err := qb.WhereIn("id", []any{xqb.Raw("? UNION ?", 1, 2)}).ToSQL()

		expectedSql := map[types.Driver]string{
			types.DriverMySQL:    "SELECT * FROM users WHERE id IN (? UNION ?)",
			types.DriverPostgres: "SELECT * FROM users WHERE id IN ($1 UNION $2)",
		}

		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, []any{1, 2}, bindings)
		assert.NoError(t, err)
	})
}

func Test_WhereIn_With_Query(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("customers").SetDialect(dialect)
		subQuery := xqb.Table("users").Select("id").Where("type", "=", "active")
		sql, bindings, err := qb.WhereIn("user_id", []any{subQuery}).ToSQL()

		expectedSql := map[types.Driver]string{
			types.DriverMySQL:    "SELECT * FROM customers WHERE user_id IN (SELECT id FROM users WHERE type = ?)",
			types.DriverPostgres: "SELECT * FROM customers WHERE user_id IN (SELECT id FROM users WHERE type = $1)",
		}
		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, 1, len(bindings))
		assert.Equal(t, []any{"active"}, bindings)
		assert.NoError(t, err)
	})
}

func Test_WhereIn_With_Query_Assert_If_There_Is_SubQuery_Use_It_Only(t *testing.T) { // TODO: need to be handled differently
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("customers").SetDialect(dialect)
		subQuery := xqb.Table("users").Select("id").Where("type", "=", "active")
		sql, bindings, err := qb.WhereIn("user_id", []any{15, 20, subQuery}).ToSQL()

		expectedSql := map[types.Driver]string{
			types.DriverMySQL:    "SELECT * FROM customers WHERE user_id IN (SELECT id FROM users WHERE type = ?)",
			types.DriverPostgres: "SELECT * FROM customers WHERE user_id IN (SELECT id FROM users WHERE type = $1)",
		}
		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, 1, len(bindings))
		assert.Equal(t, []any{"active"}, bindings)
		assert.NoError(t, err)
	})
}

func Test_WhereInQuery(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("customers").SetDialect(dialect)
		subQuery := xqb.Table("users").Select("id").Where("type", "=", "active")
		sql, bindings, err := qb.WhereInQuery("user_id", subQuery).ToSQL()

		expectedSql := map[types.Driver]string{
			types.DriverMySQL:    "SELECT * FROM customers WHERE user_id IN (SELECT id FROM users WHERE type = ?)",
			types.DriverPostgres: "SELECT * FROM customers WHERE user_id IN (SELECT id FROM users WHERE type = $1)",
		}
		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, 1, len(bindings))
		assert.Equal(t, []any{"active"}, bindings)
		assert.NoError(t, err)
	})
}

func Test_WhereExists_With_SubQuery_1(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("users").SetDialect(dialect)
		subQuery := xqb.Table("admins").Select("user_id").Where("role", "IN", []any{"superadmin", "admin"}).Latest("id")
		sql, bindings, err := qb.Select("1").WhereExists(subQuery).ToSQL()

		expectedSql := map[types.Driver]string{
			types.DriverMySQL:    "SELECT 1 FROM users WHERE EXISTS (SELECT user_id FROM admins WHERE role IN (?, ?) ORDER BY id DESC)",
			types.DriverPostgres: "SELECT 1 FROM users WHERE EXISTS (SELECT user_id FROM admins WHERE role IN ($1, $2) ORDER BY id DESC)",
		}
		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, []any{"superadmin", "admin"}, bindings)
		assert.NoError(t, err)
	})
}

func Test_WhereExists_With_SubQuery_2(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("customers").SetDialect(dialect)
		subQuery := xqb.Table("users").Select("id").Where("type", "=", "active")
		sql, bindings, err := qb.Select("1").WhereExists(subQuery).ToSQL()

		expectedSql := map[types.Driver]string{
			types.DriverMySQL:    "SELECT 1 FROM customers WHERE EXISTS (SELECT id FROM users WHERE type = ?)",
			types.DriverPostgres: "SELECT 1 FROM customers WHERE EXISTS (SELECT id FROM users WHERE type = $1)",
		}
		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, 1, len(bindings))
		assert.Equal(t, []any{"active"}, bindings)
		assert.NoError(t, err)
	})
}

func Test_WhereExists_With_Raw(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("orders").SetDialect(dialect)
		raw := xqb.Raw("SELECT user_id FROM users WHERE type = ?", "active")
		sql, bindings, err := qb.Select("1").WhereExists(raw).ToSQL()

		expectedSql := map[types.Driver]string{
			types.DriverMySQL:    "SELECT 1 FROM orders WHERE EXISTS (SELECT user_id FROM users WHERE type = ?)",
			types.DriverPostgres: "SELECT 1 FROM orders WHERE EXISTS (SELECT user_id FROM users WHERE type = $1)",
		}
		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, 1, len(bindings))
		assert.Equal(t, []any{"active"}, bindings)
		assert.NoError(t, err)
	})
}

func Test_WhereNotExists_With_SubQuery_1(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("users").SetDialect(dialect)
		subQuery := xqb.Table("admins").Select("user_id").Where("role", "IN", []any{"superadmin", "admin"}).Latest("id")
		sql, bindings, err := qb.Select("1").WhereNotExists(subQuery).ToSQL()

		expectedSql := map[types.Driver]string{
			types.DriverMySQL:    "SELECT 1 FROM users WHERE NOT EXISTS (SELECT user_id FROM admins WHERE role IN (?, ?) ORDER BY id DESC)",
			types.DriverPostgres: "SELECT 1 FROM users WHERE NOT EXISTS (SELECT user_id FROM admins WHERE role IN ($1, $2) ORDER BY id DESC)",
		}
		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, []any{"superadmin", "admin"}, bindings)
		assert.NoError(t, err)
	})
}

func Test_OrWhereExists_WithSubQuery(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("users").SetDialect(dialect)
		subQuery := xqb.Table("admins").Select("user_id").Where("role", "IN", []any{"superadmin", "admin"}).Latest("id")
		sql, bindings, err := qb.Select("1").Where("id", "=", 15).OrWhereExists(subQuery).ToSQL()

		expectedSql := map[types.Driver]string{
			types.DriverMySQL:    "SELECT 1 FROM users WHERE id = ? OR EXISTS (SELECT user_id FROM admins WHERE role IN (?, ?) ORDER BY id DESC)",
			types.DriverPostgres: "SELECT 1 FROM users WHERE id = $1 OR EXISTS (SELECT user_id FROM admins WHERE role IN ($2, $3) ORDER BY id DESC)",
		}

		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, []any{15, "superadmin", "admin"}, bindings)
		assert.NoError(t, err)
	})
}

func Test_WhereValue(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("users").SetDialect(dialect)
		sql, bindings, err := qb.WhereValue("age", ">", 18).ToSQL()

		expectedSql := map[types.Driver]string{
			types.DriverMySQL:    "SELECT * FROM users WHERE age > ?",
			types.DriverPostgres: "SELECT * FROM users WHERE age > $1",
		}

		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, []any{18}, bindings)
		assert.NoError(t, err)
	})
}

func Test_OrWhereValue(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("users").SetDialect(dialect)
		sql, bindings, err := qb.Where("name", "=", "admin").OrWhereValue("role", "=", "guest").ToSQL()

		expectedSql := map[types.Driver]string{
			types.DriverMySQL:    "SELECT * FROM users WHERE name = ? OR role = ?",
			types.DriverPostgres: "SELECT * FROM users WHERE name = $1 OR role = $2",
		}

		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, []any{"admin", "guest"}, bindings)
		assert.NoError(t, err)
	})
}

func Test_WhereExpr(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		expr := xqb.Raw("LOWER(name)")
		qb := xqb.Table("users").SetDialect(dialect)
		sql, bindings, err := qb.WhereExpr("LOWER(name)", "=", expr).ToSQL()

		expectedSql := map[types.Driver]string{
			types.DriverMySQL:    "SELECT * FROM users WHERE LOWER(name) = (LOWER(name))",
			types.DriverPostgres: "SELECT * FROM users WHERE LOWER(name) = (LOWER(name))",
		}
		assert.Equal(t, expectedSql[dialect], sql)
		assert.Empty(t, bindings)
		assert.NoError(t, err)
	})
}

func Test_OrWhereExpr(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		expr := xqb.Raw("LOWER(role)")
		qb := xqb.Table("users").SetDialect(dialect)
		sql, bindings, err := qb.Where("name", "=", "mohamed").OrWhereExpr("LOWER(role)", "=", expr).ToSQL()

		expectedSql := map[types.Driver]string{
			types.DriverMySQL:    "SELECT * FROM users WHERE name = ? OR LOWER(role) = (LOWER(role))",
			types.DriverPostgres: "SELECT * FROM users WHERE name = $1 OR LOWER(role) = (LOWER(role))",
		}
		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, []any{"mohamed"}, bindings)
		assert.NoError(t, err)
	})
}

func Test_WhereSub(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		sub := xqb.Table("admins").Select("id").Where("active", "=", true)
		qb := xqb.Table("users").SetDialect(dialect).WhereSub("admin_id", "IN", sub)
		sql, bindings, err := qb.ToSQL()

		expectedSql := map[types.Driver]string{
			types.DriverMySQL:    "SELECT * FROM users WHERE admin_id IN (SELECT id FROM admins WHERE active = ?)",
			types.DriverPostgres: "SELECT * FROM users WHERE admin_id IN (SELECT id FROM admins WHERE active = $1)",
		}

		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, []any{true}, bindings)
		assert.NoError(t, err)
	})
}

func Test_OrWhereSub(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		sub := xqb.Table("admins").Select("id").Where("active", "=", true)
		qb := xqb.Table("users").SetDialect(dialect).Where("role", "=", "staff").OrWhereSub("admin_id", "IN", sub)
		sql, bindings, err := qb.ToSQL()

		expectedSql := map[types.Driver]string{
			types.DriverMySQL:    "SELECT * FROM users WHERE role = ? OR admin_id IN (SELECT id FROM admins WHERE active = ?)",
			types.DriverPostgres: "SELECT * FROM users WHERE role = $1 OR admin_id IN (SELECT id FROM admins WHERE active = $2)",
		}

		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, []any{"staff", true}, bindings)
		assert.NoError(t, err)
	})
}

func Test_WhereNotInQuery(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		sub := xqb.Table("banned_users").Select("id")
		qb := xqb.Table("users").SetDialect(dialect).WhereNotInQuery("id", sub)
		sql, bindings, err := qb.ToSQL()

		expectedSql := map[types.Driver]string{
			types.DriverMySQL:    "SELECT * FROM users WHERE id NOT IN (SELECT id FROM banned_users)",
			types.DriverPostgres: "SELECT * FROM users WHERE id NOT IN (SELECT id FROM banned_users)",
		}

		assert.Equal(t, expectedSql[dialect], sql)
		assert.Empty(t, bindings)
		assert.NoError(t, err)
	})
}

func Test_OrWhereNotInQuery(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		sub := xqb.Table("banned_users").Select("id")
		qb := xqb.Table("users").SetDialect(dialect).Where("role", "=", "staff").OrWhereNotInQuery("id", sub)
		sql, bindings, err := qb.ToSQL()

		expectedSql := map[types.Driver]string{
			types.DriverMySQL:    "SELECT * FROM users WHERE role = ? OR id NOT IN (SELECT id FROM banned_users)",
			types.DriverPostgres: "SELECT * FROM users WHERE role = $1 OR id NOT IN (SELECT id FROM banned_users)",
		}

		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, []any{"staff"}, bindings)
		assert.NoError(t, err)
	})
}

func Test_WhereNotBetween(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("users").SetDialect(dialect).WhereNotBetween("age", 18, 60)
		sql, bindings, err := qb.ToSQL()

		expectedSql := map[types.Driver]string{
			types.DriverMySQL:    "SELECT * FROM users WHERE age NOT BETWEEN ? AND ?",
			types.DriverPostgres: "SELECT * FROM users WHERE age NOT BETWEEN $1 AND $2",
		}

		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, []any{18, 60}, bindings)
		assert.NoError(t, err)
	})
}

func Test_OrWhereNotBetween(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("users").SetDialect(dialect).Where("role", "=", "guest").OrWhereNotBetween("age", 10, 20)
		sql, bindings, err := qb.ToSQL()

		expectedSql := map[types.Driver]string{
			types.DriverMySQL:    "SELECT * FROM users WHERE role = ? OR age NOT BETWEEN ? AND ?",
			types.DriverPostgres: "SELECT * FROM users WHERE role = $1 OR age NOT BETWEEN $2 AND $3",
		}

		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, []any{"guest", 10, 20}, bindings)
		assert.NoError(t, err)
	})
}

func Test_WhereGroup_MultipleLevels(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("users").SetDialect(dialect).WhereGroup(func(q1 *xqb.QueryBuilder) {
			q1.Where("status", "=", "active").OrWhereGroup(func(q2 *xqb.QueryBuilder) {
				q2.Where("email_verified", "=", false).Where("banned", "=", false)
			})
		})

		sql, bindings, err := qb.ToSQL()

		expectedSql := map[types.Driver]string{
			types.DriverMySQL:    "SELECT * FROM users WHERE (status = ? OR (email_verified = ? AND banned = ?))",
			types.DriverPostgres: "SELECT * FROM users WHERE (status = $1 OR (email_verified = $2 AND banned = $3))",
		}

		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, []any{"active", false, false}, bindings)
		assert.NoError(t, err)

	})
}

func Test_WhereRaw_WithBindings(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("logs").SetDialect(dialect).WhereRaw("created_at > ?", "2024-01-01")
		sql, bindings, err := qb.ToSQL()

		expectedSql := map[types.Driver]string{
			types.DriverMySQL:    "SELECT * FROM logs WHERE created_at > ?",
			types.DriverPostgres: "SELECT * FROM logs WHERE created_at > $1",
		}
		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, []any{"2024-01-01"}, bindings)
		assert.NoError(t, err)
	})
}

func Test_OrWhereRaw_WithBindings(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("logs").SetDialect(dialect).
			Where("type", "=", "info").
			OrWhereRaw("created_at > ?", "2024-01-01")
		sql, bindings, err := qb.ToSQL()

		expectedSql := map[types.Driver]string{
			types.DriverMySQL:    "SELECT * FROM logs WHERE type = ? OR created_at > ?",
			types.DriverPostgres: "SELECT * FROM logs WHERE type = $1 OR created_at > $2",
		}
		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, []any{"info", "2024-01-01"}, bindings)
		assert.NoError(t, err)
	})
}

func Test_WhereIn_Empty(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("users").SetDialect(dialect).WhereIn("id", []any{})
		sql, bindings, err := qb.ToSQL()

		expectedSql := map[types.Driver]string{
			types.DriverMySQL:    "SELECT * FROM users",
			types.DriverPostgres: "SELECT * FROM users",
		}

		assert.Equal(t, expectedSql[dialect], sql)
		assert.Empty(t, bindings)
		assert.NoError(t, err)
	})
}

func Test_WhereNotIn_Empty(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("users").SetDialect(dialect).WhereNotIn("id", []any{})
		sql, bindings, err := qb.ToSQL()

		expectedSql := map[types.Driver]string{
			types.DriverMySQL:    "SELECT * FROM users",
			types.DriverPostgres: "SELECT * FROM users",
		}

		assert.Equal(t, expectedSql[dialect], sql)
		assert.Empty(t, bindings)
		assert.NoError(t, err)
	})
}

func Test_WhereBetween_WithExpr(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		min := xqb.Raw("NOW() - INTERVAL 1 DAY")
		max := xqb.Raw("NOW()")
		qb := xqb.Table("logs").SetDialect(dialect).WhereBetween("created_at", min, max)
		sql, bindings, err := qb.ToSQL()

		expectedSql := map[types.Driver]string{
			types.DriverMySQL:    "SELECT * FROM logs WHERE created_at BETWEEN NOW() - INTERVAL 1 DAY AND NOW()",
			types.DriverPostgres: "SELECT * FROM logs WHERE created_at BETWEEN NOW() - INTERVAL 1 DAY AND NOW()",
		}
		assert.Equal(t, expectedSql[dialect], sql)
		assert.Empty(t, bindings)
		assert.NoError(t, err)
	})
}

func Test_WhereExists_Chained(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		sub := xqb.Table("admins").Select("id").Where("active", "=", true)
		qb := xqb.Table("users").SetDialect(dialect).Where("status", "=", "staff").WhereExists(sub)
		sql, bindings, err := qb.ToSQL()

		expectedSql := map[types.Driver]string{
			types.DriverMySQL:    "SELECT * FROM users WHERE status = ? AND EXISTS (SELECT id FROM admins WHERE active = ?)",
			types.DriverPostgres: "SELECT * FROM users WHERE status = $1 AND EXISTS (SELECT id FROM admins WHERE active = $2)",
		}

		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, []any{"staff", true}, bindings)
		assert.NoError(t, err)
	})
}

func Test_WhereNotExists_Chained(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		sub := xqb.Table("admins").Select("id").Where("active", "=", false)
		qb := xqb.Table("users").SetDialect(dialect).Where("status", "=", "guest").WhereNotExists(sub)
		sql, bindings, err := qb.ToSQL()

		expectedSql := map[types.Driver]string{
			types.DriverMySQL:    "SELECT * FROM users WHERE status = ? AND NOT EXISTS (SELECT id FROM admins WHERE active = ?)",
			types.DriverPostgres: "SELECT * FROM users WHERE status = $1 AND NOT EXISTS (SELECT id FROM admins WHERE active = $2)",
		}

		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, []any{"guest", false}, bindings)
		assert.NoError(t, err)
	})
}

func Test_Mixed_WhereRaw_And_Normal(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		jsonSQL := xqb.JsonExtract("meta", "age", "").Dialects[dialect.String()].SQL

		qb := xqb.Table("users").SetDialect(dialect).
			WhereRaw(jsonSQL+" > ?", 18).
			Where("active", "=", true)

		sql, bindings, err := qb.ToSQL()

		expectedSql := map[types.Driver]string{
			types.DriverMySQL:    "SELECT * FROM users WHERE JSON_EXTRACT(meta, '$.age') > ? AND active = ?",
			types.DriverPostgres: "SELECT * FROM users WHERE meta->>'age' > $1 AND active = $2",
		}

		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, []any{18, true}, bindings)
		assert.NoError(t, err)
	})
}

func Test_OrWhereGroup_Complex(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("products").SetDialect(dialect).
			Where("stock", ">", 0).
			OrWhereGroup(func(q *xqb.QueryBuilder) {
				q.Where("archived", "=", false).
					Where("discontinued", "=", false)
			})
		sql, bindings, err := qb.ToSQL()

		expectedSql := map[types.Driver]string{
			types.DriverMySQL:    "SELECT * FROM products WHERE stock > ? OR (archived = ? AND discontinued = ?)",
			types.DriverPostgres: "SELECT * FROM products WHERE stock > $1 OR (archived = $2 AND discontinued = $3)",
		}
		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, []any{0, false, false}, bindings)
		assert.NoError(t, err)
	})
}

func Test_WhereExpr_ComplexBothSides(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		left := xqb.Raw("LOWER(username)")
		right := xqb.Raw("LOWER(?)", "Mohamed")
		qb := xqb.Table("users").SetDialect(dialect).Where(left, "=", right)
		sql, bindings, err := qb.ToSQL()

		expectedSql := map[types.Driver]string{
			types.DriverMySQL:    "SELECT * FROM users WHERE (LOWER(username)) = (LOWER(?))",
			types.DriverPostgres: "SELECT * FROM users WHERE (LOWER(username)) = (LOWER($1))",
		}
		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, []any{"Mohamed"}, bindings)
		assert.NoError(t, err)
	})
}

func TestWhereWithRawExpressions(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("users").SetDialect(dialect)
		sql, bindings, err := qb.Where(xqb.Raw("LOWER(name)"), "=", "john").ToSQL()

		expectedSql := map[types.Driver]string{
			types.DriverMySQL:    "SELECT * FROM users WHERE LOWER(name) = ?",
			types.DriverPostgres: "SELECT * FROM users WHERE LOWER(name) = $1",
		}
		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, []any{"john"}, bindings)
		assert.NoError(t, err)
	})
}

func TestWhereRaw(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("users").SetDialect(dialect)
		sql, bindings, err := qb.WhereRaw("LOWER(name) = ? OR LOWER(email) = ?", "john", "john@example.com").ToSQL()

		expectedSql := map[types.Driver]string{
			types.DriverMySQL:    "SELECT * FROM users WHERE LOWER(name) = ? OR LOWER(email) = ?",
			types.DriverPostgres: "SELECT * FROM users WHERE LOWER(name) = $1 OR LOWER(email) = $2",
		}
		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, []any{"john", "john@example.com"}, bindings)
		assert.NoError(t, err)
	})
}

func TestWhereNull(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("users").SetDialect(dialect)
		sql, bindings, err := qb.WhereNull("deleted_at").ToSQL()

		expectedSql := map[types.Driver]string{
			types.DriverMySQL:    "SELECT * FROM users WHERE deleted_at IS NULL",
			types.DriverPostgres: "SELECT * FROM users WHERE deleted_at IS NULL",
		}
		assert.Equal(t, expectedSql[dialect], sql)
		assert.Empty(t, bindings)
		assert.NoError(t, err)
	})

}

func TestWhereNotNull(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("users").SetDialect(dialect)
		sql, bindings, err := qb.WhereNotNull("email").ToSQL()

		expectedSql := map[types.Driver]string{
			types.DriverMySQL:    "SELECT * FROM users WHERE email IS NOT NULL",
			types.DriverPostgres: "SELECT * FROM users WHERE email IS NOT NULL",
		}
		assert.Equal(t, expectedSql[dialect], sql)
		assert.Empty(t, bindings)
		assert.NoError(t, err)
	})
}

func TestWhereNullWithSelect(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("users").SetDialect(dialect)
		sql, bindings, err := qb.Select("id", "name").Where("name", "LIKE", "%mohamedsheta%").WhereNull("deleted_at").ToSQL()

		expectedSql := map[types.Driver]string{
			types.DriverMySQL:    "SELECT id, name FROM users WHERE name LIKE ? AND deleted_at IS NULL",
			types.DriverPostgres: "SELECT id, name FROM users WHERE name LIKE $1 AND deleted_at IS NULL",
		}

		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, []any{"%mohamedsheta%"}, bindings)
		assert.NoError(t, err)
	})
}

func TestWhereNotNullWithSelect(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("users").SetDialect(dialect)
		sql, bindings, err := qb.Select("id", "name").Where("name", "LIKE", "%mohamedsheta%").WhereNotNull("email").ToSQL()

		expectedSql := map[types.Driver]string{
			types.DriverMySQL:    "SELECT id, name FROM users WHERE name LIKE ? AND email IS NOT NULL",
			types.DriverPostgres: "SELECT id, name FROM users WHERE name LIKE $1 AND email IS NOT NULL",
		}

		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, []any{"%mohamedsheta%"}, bindings)
		assert.NoError(t, err)
	})
}

func TestWhereIn(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("users").SetDialect(dialect)
		sql, bindings, err := qb.WhereIn("id", []any{1, 2, 3}).ToSQL()

		expectedSql := map[types.Driver]string{
			types.DriverMySQL:    "SELECT * FROM users WHERE id IN (?, ?, ?)",
			types.DriverPostgres: "SELECT * FROM users WHERE id IN ($1, $2, $3)",
		}

		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, []any{1, 2, 3}, bindings)
		assert.NoError(t, err)
	})
}

func TestWhereNotIn(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("users").SetDialect(dialect)
		sql, bindings, err := qb.WhereNotIn("id", []any{1, 2, 3}).ToSQL()

		expectedSql := map[types.Driver]string{
			types.DriverMySQL:    "SELECT * FROM users WHERE id NOT IN (?, ?, ?)",
			types.DriverPostgres: "SELECT * FROM users WHERE id NOT IN ($1, $2, $3)",
		}

		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, []any{1, 2, 3}, bindings)
		assert.NoError(t, err)
	})
}

func TestWhereInWithSubquery(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("users").SetDialect(dialect)
		subQuery := xqb.Table("orders").Select("user_id").Where("status", "=", "active")
		sql, bindings, err := qb.WhereIn("id", []any{subQuery}).ToSQL()

		expectedSql := map[types.Driver]string{
			types.DriverMySQL:    "SELECT * FROM users WHERE id IN (SELECT user_id FROM orders WHERE status = ?)",
			types.DriverPostgres: "SELECT * FROM users WHERE id IN (SELECT user_id FROM orders WHERE status = $1)",
		}
		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, []any{"active"}, bindings)
		assert.NoError(t, err)
	})
}

func TestWhereBetween(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("users").SetDialect(dialect)
		sql, bindings, err := qb.WhereBetween("age", 18, 30).ToSQL()

		expectedSql := map[types.Driver]string{
			types.DriverMySQL:    "SELECT * FROM users WHERE age BETWEEN ? AND ?",
			types.DriverPostgres: "SELECT * FROM users WHERE age BETWEEN $1 AND $2",
		}

		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, []any{18, 30}, bindings)
		assert.NoError(t, err)
	})
}

func TestWhereRawWithSubqueryRaw(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("users").SetDialect(dialect)
		sql, bindings, err := qb.WhereRaw("EXISTS (SELECT 1 FROM orders WHERE orders.user_id = users.id AND amount > ?)", 1000).ToSQL()

		expectedSql := map[types.Driver]string{
			types.DriverMySQL:    "SELECT * FROM users WHERE EXISTS (SELECT 1 FROM orders WHERE orders.user_id = users.id AND amount > ?)",
			types.DriverPostgres: "SELECT * FROM users WHERE EXISTS (SELECT 1 FROM orders WHERE orders.user_id = users.id AND amount > $1)",
		}
		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, []any{1000}, bindings)
		assert.NoError(t, err)
	})
}

func Test_WhereGroup(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("orders").SetDialect(dialect)
		sql, bindings, err := qb.WhereGroup(func(qb *xqb.QueryBuilder) {
			qb.Where("email", "=", "mohamed@mail.com").
				OrWhere("username", "=", "mohamed")
		}).WhereGroup(func(qb *xqb.QueryBuilder) {
			qb.Where("uuid", "=", "bbee7431-454d-4a8a-9435-961d191de2a7").OrWhere("user_id", "=", 4)
		}).OrWhereGroup(func(qb *xqb.QueryBuilder) {
			qb.Where("username", "=", "ahmed").Where("user_id", "=", 6)
		}).ToSQL()

		expectedSql := map[types.Driver]string{
			types.DriverMySQL:    "SELECT * FROM orders WHERE (email = ? OR username = ?) AND (uuid = ? OR user_id = ?) OR (username = ? AND user_id = ?)",
			types.DriverPostgres: "SELECT * FROM orders WHERE (email = $1 OR username = $2) AND (uuid = $3 OR user_id = $4) OR (username = $5 AND user_id = $6)",
		}
		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, []any{"mohamed@mail.com", "mohamed", "bbee7431-454d-4a8a-9435-961d191de2a7", 4, "ahmed", 6}, bindings)
		assert.NoError(t, err)
	})

}

func Test_Where_Is_Null(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("users").SetDialect(dialect)
		sql, bindings, err := qb.Where("id", "=", 1).Where("deleted_at", "IS NULL", nil).ToSQL()

		expectedSql := map[types.Driver]string{
			types.DriverMySQL:    "SELECT * FROM users WHERE id = ? AND deleted_at IS NULL",
			types.DriverPostgres: "SELECT * FROM users WHERE id = $1 AND deleted_at IS NULL",
		}
		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, 1, len(bindings))
		assert.Equal(t, []any{1}, bindings)
		assert.NoError(t, err)
	})
}
