package xqb_test

import (
	"testing"

	"github.com/iMohamedSheta/xqb"
	"github.com/iMohamedSheta/xqb/shared/types"
	"github.com/stretchr/testify/assert"
)

func Test_Join_Closure_On(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		qb := xqb.Table("users").SetDialect(dialect).Join("orders", func(j *xqb.JoinClause) {
			j.On("users.id", "=", "orders.user_id")
		})
		sql, bindings, err := qb.ToSql()
		expectedSQL := map[types.Dialect]string{
			types.DialectMySql:    "SELECT * FROM `users` JOIN `orders` ON users.id = orders.user_id",
			types.DialectPostgres: `SELECT * FROM "users" JOIN "orders" ON users.id = orders.user_id`,
		}
		assert.Equal(t, expectedSQL[dialect], sql)
		assert.NoError(t, err)
		assert.Empty(t, bindings)
	})
}

func Test_Join_Closure_On_OrOn(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		qb := xqb.Table("users").SetDialect(dialect).Join("orders", func(j *xqb.JoinClause) {
			j.On("users.id", "=", "orders.user_id").
				OrOn("users.alt_id", "=", "orders.user_id")
		})
		sql, bindings, err := qb.ToSql()
		expectedSQL := map[types.Dialect]string{
			types.DialectMySql:    "SELECT * FROM `users` JOIN `orders` ON users.id = orders.user_id OR users.alt_id = orders.user_id",
			types.DialectPostgres: `SELECT * FROM "users" JOIN "orders" ON users.id = orders.user_id OR users.alt_id = orders.user_id`,
		}
		assert.Equal(t, expectedSQL[dialect], sql)
		assert.NoError(t, err)
		assert.Empty(t, bindings)
	})
}

func Test_Join_Closure_Where(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		qb := xqb.Table("users").SetDialect(dialect).Join("orders", func(j *xqb.JoinClause) {
			j.On("users.id", "=", "orders.user_id").
				Where("orders.type", "=", "invoice")
		})
		sql, bindings, err := qb.ToSql()
		expectedSQL := map[types.Dialect]string{
			types.DialectMySql:    "SELECT * FROM `users` JOIN `orders` ON users.id = orders.user_id AND orders.type = ?",
			types.DialectPostgres: `SELECT * FROM "users" JOIN "orders" ON users.id = orders.user_id AND orders.type = $1`,
		}
		assert.Equal(t, expectedSQL[dialect], sql)
		assert.NoError(t, err)
		assert.Equal(t, []any{"invoice"}, bindings)
	})
}

func Test_Join_Closure_OrWhere(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		qb := xqb.Table("users").SetDialect(dialect).Join("orders", func(j *xqb.JoinClause) {
			j.On("users.id", "=", "orders.user_id").
				OrWhere("orders.type", "=", "receipt")
		})
		sql, bindings, err := qb.ToSql()
		expectedSQL := map[types.Dialect]string{
			types.DialectMySql:    "SELECT * FROM `users` JOIN `orders` ON users.id = orders.user_id OR orders.type = ?",
			types.DialectPostgres: `SELECT * FROM "users" JOIN "orders" ON users.id = orders.user_id OR orders.type = $1`,
		}
		assert.Equal(t, expectedSQL[dialect], sql)
		assert.NoError(t, err)
		assert.Equal(t, []any{"receipt"}, bindings)
	})
}

func Test_Join_Closure_WhereNull(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		qb := xqb.Table("users").SetDialect(dialect).Join("orders", func(j *xqb.JoinClause) {
			j.On("users.id", "=", "orders.user_id").
				WhereNull("orders.deleted_at")
		})
		sql, bindings, err := qb.ToSql()
		expectedSQL := map[types.Dialect]string{
			types.DialectMySql:    "SELECT * FROM `users` JOIN `orders` ON users.id = orders.user_id AND orders.deleted_at IS NULL",
			types.DialectPostgres: `SELECT * FROM "users" JOIN "orders" ON users.id = orders.user_id AND orders.deleted_at IS NULL`,
		}
		assert.Equal(t, expectedSQL[dialect], sql)
		assert.NoError(t, err)
		assert.Empty(t, bindings)
	})
}

func Test_Join_Closure_OrWhereNull(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		qb := xqb.Table("users").SetDialect(dialect).Join("orders", func(j *xqb.JoinClause) {
			j.On("users.id", "=", "orders.user_id").
				OrWhereNull("orders.cancelled_at")
		})
		sql, bindings, err := qb.ToSql()
		expectedSQL := map[types.Dialect]string{
			types.DialectMySql:    "SELECT * FROM `users` JOIN `orders` ON users.id = orders.user_id OR orders.cancelled_at IS NULL",
			types.DialectPostgres: `SELECT * FROM "users" JOIN "orders" ON users.id = orders.user_id OR orders.cancelled_at IS NULL`,
		}
		assert.Equal(t, expectedSQL[dialect], sql)
		assert.NoError(t, err)
		assert.Empty(t, bindings)
	})
}

func Test_Join_Closure_WhereNotNull(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		qb := xqb.Table("users").SetDialect(dialect).Join("orders", func(j *xqb.JoinClause) {
			j.On("users.id", "=", "orders.user_id").
				WhereNotNull("orders.confirmed_at")
		})
		sql, bindings, err := qb.ToSql()
		expectedSQL := map[types.Dialect]string{
			types.DialectMySql:    "SELECT * FROM `users` JOIN `orders` ON users.id = orders.user_id AND orders.confirmed_at IS NOT NULL",
			types.DialectPostgres: `SELECT * FROM "users" JOIN "orders" ON users.id = orders.user_id AND orders.confirmed_at IS NOT NULL`,
		}
		assert.Equal(t, expectedSQL[dialect], sql)
		assert.NoError(t, err)
		assert.Empty(t, bindings)
	})
}

func Test_Join_Closure_OrWhereNotNull(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		qb := xqb.Table("users").SetDialect(dialect).Join("orders", func(j *xqb.JoinClause) {
			j.On("users.id", "=", "orders.user_id").
				OrWhereNotNull("orders.confirmed_at")
		})
		sql, bindings, err := qb.ToSql()
		expectedSQL := map[types.Dialect]string{
			types.DialectMySql:    "SELECT * FROM `users` JOIN `orders` ON users.id = orders.user_id OR orders.confirmed_at IS NOT NULL",
			types.DialectPostgres: `SELECT * FROM "users" JOIN "orders" ON users.id = orders.user_id OR orders.confirmed_at IS NOT NULL`,
		}
		assert.Equal(t, expectedSQL[dialect], sql)
		assert.NoError(t, err)
		assert.Empty(t, bindings)
	})
}

func Test_Join_Closure_OnRaw(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		qb := xqb.Table("users").SetDialect(dialect).Join("orders", func(j *xqb.JoinClause) {
			j.OnRaw("users.id = orders.user_id AND orders.active = ?", 1)
		})
		sql, bindings, err := qb.ToSql()
		expectedSQL := map[types.Dialect]string{
			types.DialectMySql:    "SELECT * FROM `users` JOIN `orders` ON users.id = orders.user_id AND orders.active = ?",
			types.DialectPostgres: `SELECT * FROM "users" JOIN "orders" ON users.id = orders.user_id AND orders.active = $1`,
		}
		assert.Equal(t, expectedSQL[dialect], sql)
		assert.NoError(t, err)
		assert.Equal(t, []any{1}, bindings)
	})
}

func Test_Join_Closure_OrOnRaw(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		qb := xqb.Table("users").SetDialect(dialect).Join("orders", func(j *xqb.JoinClause) {
			j.On("users.id", "=", "orders.user_id").
				OrOnRaw("orders.guest_id = ? AND orders.active = ?", 99, 1)
		})
		sql, bindings, err := qb.ToSql()
		expectedSQL := map[types.Dialect]string{
			types.DialectMySql:    "SELECT * FROM `users` JOIN `orders` ON users.id = orders.user_id OR orders.guest_id = ? AND orders.active = ?",
			types.DialectPostgres: `SELECT * FROM "users" JOIN "orders" ON users.id = orders.user_id OR orders.guest_id = $1 AND orders.active = $2`,
		}
		assert.Equal(t, expectedSQL[dialect], sql)
		assert.NoError(t, err)
		assert.Equal(t, []any{99, 1}, bindings)
	})
}

func Test_Join_Closure_OnGroup(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		qb := xqb.Table("users").SetDialect(dialect).Join("orders", func(j *xqb.JoinClause) {
			j.On("users.id", "=", "orders.user_id").
				OnGroup(func(g *xqb.JoinClause) {
					g.On("orders.type", "=", "orders.default_type").
						OrOn("orders.type", "=", "orders.fallback_type")
				})
		})
		sql, bindings, err := qb.ToSql()
		expectedSQL := map[types.Dialect]string{
			types.DialectMySql:    "SELECT * FROM `users` JOIN `orders` ON users.id = orders.user_id AND (orders.type = orders.default_type OR orders.type = orders.fallback_type)",
			types.DialectPostgres: `SELECT * FROM "users" JOIN "orders" ON users.id = orders.user_id AND (orders.type = orders.default_type OR orders.type = orders.fallback_type)`,
		}
		assert.Equal(t, expectedSQL[dialect], sql)
		assert.NoError(t, err)
		assert.Empty(t, bindings)
	})
}

func Test_Join_Closure_OrOnGroup(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		qb := xqb.Table("users").SetDialect(dialect).Join("orders", func(j *xqb.JoinClause) {
			j.On("users.id", "=", "orders.user_id").
				OrOnGroup(func(g *xqb.JoinClause) {
					g.On("orders.type", "=", "orders.default_type").
						Where("orders.status", "=", "active")
				})
		})
		sql, bindings, err := qb.ToSql()
		expectedSQL := map[types.Dialect]string{
			types.DialectMySql:    "SELECT * FROM `users` JOIN `orders` ON users.id = orders.user_id OR (orders.type = orders.default_type AND orders.status = ?)",
			types.DialectPostgres: `SELECT * FROM "users" JOIN "orders" ON users.id = orders.user_id OR (orders.type = orders.default_type AND orders.status = $1)`,
		}
		assert.Equal(t, expectedSQL[dialect], sql)
		assert.NoError(t, err)
		assert.Equal(t, []any{"active"}, bindings)
	})
}

func Test_Join_Closure_EmptyGroup_IsIgnored(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		qb := xqb.Table("users").SetDialect(dialect).Join("orders", func(j *xqb.JoinClause) {
			j.On("users.id", "=", "orders.user_id").
				OnGroup(func(g *xqb.JoinClause) {
					// empty group — should be silently ignored
				})
		})
		sql, bindings, err := qb.ToSql()
		expectedSQL := map[types.Dialect]string{
			types.DialectMySql:    "SELECT * FROM `users` JOIN `orders` ON users.id = orders.user_id",
			types.DialectPostgres: `SELECT * FROM "users" JOIN "orders" ON users.id = orders.user_id`,
		}
		assert.Equal(t, expectedSQL[dialect], sql)
		assert.NoError(t, err)
		assert.Empty(t, bindings)
	})
}

func Test_LeftJoin_Closure(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		qb := xqb.Table("users").SetDialect(dialect).LeftJoin("comments", func(j *xqb.JoinClause) {
			j.On("users.id", "=", "comments.user_id").
				Where("comments.approved", "=", true)
		})
		sql, bindings, err := qb.ToSql()
		expectedSQL := map[types.Dialect]string{
			types.DialectMySql:    "SELECT * FROM `users` LEFT JOIN `comments` ON users.id = comments.user_id AND comments.approved = ?",
			types.DialectPostgres: `SELECT * FROM "users" LEFT JOIN "comments" ON users.id = comments.user_id AND comments.approved = $1`,
		}
		assert.Equal(t, expectedSQL[dialect], sql)
		assert.NoError(t, err)
		assert.Equal(t, []any{true}, bindings)
	})
}

func Test_Join_Closure_Complex_Mixed_Conditions(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		qb := xqb.Table("users").SetDialect(dialect).
			Join("orders", func(j *xqb.JoinClause) {
				j.On("users.id", "=", "orders.user_id").
					Where("orders.status", "=", "active").
					WhereNotNull("orders.confirmed_at").
					OnGroup(func(g *xqb.JoinClause) {
						g.On("orders.type", "=", "orders.primary_type").
							OrWhere("orders.priority", "=", 1)
					})
			})
		sql, bindings, err := qb.ToSql()
		expectedSQL := map[types.Dialect]string{
			types.DialectMySql: "SELECT * FROM `users` JOIN `orders` ON users.id = orders.user_id" +
				" AND orders.status = ?" +
				" AND orders.confirmed_at IS NOT NULL" +
				" AND (orders.type = orders.primary_type OR orders.priority = ?)",
			types.DialectPostgres: `SELECT * FROM "users" JOIN "orders" ON users.id = orders.user_id` +
				` AND orders.status = $1` +
				` AND orders.confirmed_at IS NOT NULL` +
				` AND (orders.type = orders.primary_type OR orders.priority = $2)`,
		}
		assert.Equal(t, expectedSQL[dialect], sql)
		assert.NoError(t, err)
		assert.Equal(t, []any{"active", 1}, bindings)
	})
}

func Test_Join_Closure_OnRaw_Multiple_Bindings(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		qb := xqb.Table("users").SetDialect(dialect).Join("orders", func(j *xqb.JoinClause) {
			j.On("users.id", "=", "orders.user_id").
				OnRaw("orders.region = ? AND orders.priority > ?", "EU", 2)
		})
		sql, bindings, err := qb.ToSql()
		expectedSQL := map[types.Dialect]string{
			types.DialectMySql:    "SELECT * FROM `users` JOIN `orders` ON users.id = orders.user_id AND orders.region = ? AND orders.priority > ?",
			types.DialectPostgres: `SELECT * FROM "users" JOIN "orders" ON users.id = orders.user_id AND orders.region = $1 AND orders.priority > $2`,
		}
		assert.Equal(t, expectedSQL[dialect], sql)
		assert.NoError(t, err)
		assert.Equal(t, []any{"EU", 2}, bindings)
	})
}
