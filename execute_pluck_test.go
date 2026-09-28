package xqb_test

import (
	"testing"

	"github.com/iMohamedSheta/xqb"
	xqbErr "github.com/iMohamedSheta/xqb/shared/errors"
	"github.com/iMohamedSheta/xqb/shared/types"
	"github.com/stretchr/testify/assert"
)

// PluckSliceSql Tests

func Test_PluckSliceSql_WithValueField(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		qb := xqb.Table("users").SetDialect(dialect).Where("name", "LIKE", "%mohamed%")
		sql, bindings, err := qb.PluckSliceSql("name")

		expectedSQL := map[types.Dialect]string{
			types.DialectMySql:     "SELECT `name` FROM `users` WHERE `name` LIKE ?",
			types.DialectMariaDB:   "SELECT `name` FROM `users` WHERE `name` LIKE ?",
			types.DialectPostgres:  `SELECT "name" FROM "users" WHERE "name" LIKE $1`,
			types.DialectSQLite:    `SELECT "name" FROM "users" WHERE "name" LIKE ?`,
			types.DialectSQLServer: "SELECT [name] FROM [users] WHERE [name] LIKE @p1",
		}

		assert.Equal(t, expectedSQL[dialect], sql)
		assert.Equal(t, []any{"%mohamed%"}, bindings)
		assert.NoError(t, err)
	})
}

func Test_PluckSliceSql_WithComplexQuery(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		qb := xqb.Table("users").SetDialect(dialect)
		qb.Where("age", ">", 18).OrderBy("created_at", "DESC").Limit(10)
		sql, bindings, err := qb.PluckSliceSql("email")

		expectedSQL := map[types.Dialect]string{
			types.DialectMySql:     "SELECT `email` FROM `users` WHERE `age` > ? ORDER BY `created_at` DESC LIMIT 10",
			types.DialectMariaDB:   "SELECT `email` FROM `users` WHERE `age` > ? ORDER BY `created_at` DESC LIMIT 10",
			types.DialectPostgres:  `SELECT "email" FROM "users" WHERE "age" > $1 ORDER BY "created_at" DESC LIMIT 10`,
			types.DialectSQLite:    `SELECT "email" FROM "users" WHERE "age" > ? ORDER BY "created_at" DESC LIMIT 10`,
			types.DialectSQLServer: "SELECT TOP 10 [email] FROM [users] WHERE [age] > @p1 ORDER BY [created_at] DESC",
		}

		assert.Equal(t, expectedSQL[dialect], sql)
		assert.Equal(t, []any{18}, bindings)
		assert.NoError(t, err)
	})
}

func Test_PluckSliceSql_EmptyValueField_ReturnsError(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		qb := xqb.Table("users").SetDialect(dialect)
		sql, bindings, err := qb.PluckSliceSql("")

		assert.Error(t, err)
		assert.ErrorIs(t, err, xqbErr.ErrInvalidQuery)
		assert.Empty(t, bindings)
		assert.Equal(t, "", sql)
	})
}

func Test_PluckSliceSql_OverridesExistingSelect(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		qb := xqb.Table("users").SetDialect(dialect)
		qb.Select("id", "name", "email").Where("active", "=", true)
		sql, bindings, err := qb.PluckSliceSql("name")

		expectedSQL := map[types.Dialect]string{
			types.DialectMySql:     "SELECT `name` FROM `users` WHERE `active` = ?",
			types.DialectMariaDB:   "SELECT `name` FROM `users` WHERE `active` = ?",
			types.DialectPostgres:  `SELECT "name" FROM "users" WHERE "active" = $1`,
			types.DialectSQLite:    `SELECT "name" FROM "users" WHERE "active" = ?`,
			types.DialectSQLServer: "SELECT [name] FROM [users] WHERE [active] = @p1",
		}

		assert.Equal(t, expectedSQL[dialect], sql)
		assert.Equal(t, []any{true}, bindings)
		assert.NoError(t, err)
	})
}

// PluckMapSql Tests

func Test_PluckMapSql_WithValueAndKeyFields(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		qb := xqb.Table("users").SetDialect(dialect).Where("status", "=", "active")
		sql, bindings, err := qb.PluckMapSql("name", "id")

		expectedSQL := map[types.Dialect]string{
			types.DialectMySql:     "SELECT `name`, `id` FROM `users` WHERE `status` = ?",
			types.DialectMariaDB:   "SELECT `name`, `id` FROM `users` WHERE `status` = ?",
			types.DialectPostgres:  `SELECT "name", "id" FROM "users" WHERE "status" = $1`,
			types.DialectSQLite:    `SELECT "name", "id" FROM "users" WHERE "status" = ?`,
			types.DialectSQLServer: "SELECT [name], [id] FROM [users] WHERE [status] = @p1",
		}

		assert.Equal(t, expectedSQL[dialect], sql)
		assert.Equal(t, []any{"active"}, bindings)
		assert.NoError(t, err)
	})
}

func Test_PluckMapSql_WithComplexQuery(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		qb := xqb.Table("products").SetDialect(dialect)
		qb.Where("category", "=", "electronics").Where("price", "<", 1000).OrderBy("price", "ASC")
		sql, bindings, err := qb.PluckMapSql("title", "sku")

		expectedSQL := map[types.Dialect]string{
			types.DialectMySql:     "SELECT `title`, `sku` FROM `products` WHERE `category` = ? AND `price` < ? ORDER BY `price` ASC",
			types.DialectMariaDB:   "SELECT `title`, `sku` FROM `products` WHERE `category` = ? AND `price` < ? ORDER BY `price` ASC",
			types.DialectPostgres:  `SELECT "title", "sku" FROM "products" WHERE "category" = $1 AND "price" < $2 ORDER BY "price" ASC`,
			types.DialectSQLite:    `SELECT "title", "sku" FROM "products" WHERE "category" = ? AND "price" < ? ORDER BY "price" ASC`,
			types.DialectSQLServer: "SELECT [title], [sku] FROM [products] WHERE [category] = @p1 AND [price] < @p2 ORDER BY [price] ASC",
		}

		assert.Equal(t, expectedSQL[dialect], sql)
		assert.Equal(t, []any{"electronics", 1000}, bindings)
		assert.NoError(t, err)
	})
}

func Test_PluckMapSql_EmptyValueField_ReturnsError(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		qb := xqb.Table("users").SetDialect(dialect)
		sql, bindings, err := qb.PluckMapSql("", "id")

		assert.Error(t, err)
		assert.ErrorIs(t, err, xqbErr.ErrInvalidQuery)
		assert.Empty(t, bindings)
		assert.Equal(t, "", sql)
	})
}

func Test_PluckMapSql_EmptyKeyField_ReturnsError(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		qb := xqb.Table("users").SetDialect(dialect)
		sql, bindings, err := qb.PluckMapSql("name", "")

		assert.Error(t, err)
		assert.ErrorIs(t, err, xqbErr.ErrInvalidQuery)
		assert.Empty(t, bindings)
		assert.Equal(t, "", sql)
	})
}

func Test_PluckMapSql_BothFieldsEmpty_ReturnsError(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		qb := xqb.Table("users").SetDialect(dialect)
		sql, bindings, err := qb.PluckMapSql("", "")

		assert.Error(t, err)
		assert.ErrorIs(t, err, xqbErr.ErrInvalidQuery)
		assert.Empty(t, bindings)
		assert.Equal(t, "", sql)
	})
}

func Test_PluckMapSql_OverridesExistingSelect(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		qb := xqb.Table("users").SetDialect(dialect)
		qb.Select("id", "name", "email", "phone").Where("country", "=", "US")
		sql, bindings, err := qb.PluckMapSql("email", "id")

		expectedSQL := map[types.Dialect]string{
			types.DialectMySql:     "SELECT `email`, `id` FROM `users` WHERE `country` = ?",
			types.DialectMariaDB:   "SELECT `email`, `id` FROM `users` WHERE `country` = ?",
			types.DialectPostgres:  `SELECT "email", "id" FROM "users" WHERE "country" = $1`,
			types.DialectSQLite:    `SELECT "email", "id" FROM "users" WHERE "country" = ?`,
			types.DialectSQLServer: "SELECT [email], [id] FROM [users] WHERE [country] = @p1",
		}

		assert.Equal(t, expectedSQL[dialect], sql)
		assert.Equal(t, []any{"US"}, bindings)
		assert.NoError(t, err)
	})
}

func Test_PluckSlice_EmptyValueField_ReturnsError(t *testing.T) {
	qb := xqb.Table("users")
	values, err := qb.PluckSlice("")

	assert.Error(t, err)
	assert.ErrorIs(t, err, xqbErr.ErrInvalidQuery)
	assert.Nil(t, values)
}

func Test_PluckMap_EmptyValueField_ReturnsError(t *testing.T) {
	qb := xqb.Table("users")
	mappedValues, err := qb.PluckMap("", "id")

	assert.Error(t, err)
	assert.ErrorIs(t, err, xqbErr.ErrInvalidQuery)
	assert.Nil(t, mappedValues)
}

func Test_PluckMap_EmptyKeyField_ReturnsError(t *testing.T) {
	qb := xqb.Table("users")
	mappedValues, err := qb.PluckMap("name", "")

	assert.Error(t, err)
	assert.ErrorIs(t, err, xqbErr.ErrInvalidQuery)
	assert.Nil(t, mappedValues)
}

func Test_PluckMap_BothFieldsEmpty_ReturnsError(t *testing.T) {
	qb := xqb.Table("users")
	mappedValues, err := qb.PluckMap("", "")

	assert.Error(t, err)
	assert.ErrorIs(t, err, xqbErr.ErrInvalidQuery)
	assert.Nil(t, mappedValues)
}

func Test_PluckSliceSql_WithJoins(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		qb := xqb.Table("users").SetDialect(dialect)
		qb.Join("posts", "users.id = posts.user_id")
		qb.Where("posts.published", "=", true)
		sql, bindings, err := qb.PluckSliceSql("users.name")

		expectedSQL := map[types.Dialect]string{
			types.DialectMySql:     "SELECT `users`.`name` FROM `users` JOIN `posts` ON users.id = posts.user_id WHERE `posts`.`published` = ?",
			types.DialectMariaDB:   "SELECT `users`.`name` FROM `users` JOIN `posts` ON users.id = posts.user_id WHERE `posts`.`published` = ?",
			types.DialectPostgres:  `SELECT "users"."name" FROM "users" JOIN "posts" ON users.id = posts.user_id WHERE "posts"."published" = $1`,
			types.DialectSQLite:    `SELECT "users"."name" FROM "users" JOIN "posts" ON users.id = posts.user_id WHERE "posts"."published" = ?`,
			types.DialectSQLServer: "SELECT [users].[name] FROM [users] JOIN [posts] ON users.id = posts.user_id WHERE [posts].[published] = @p1",
		}

		assert.Equal(t, expectedSQL[dialect], sql)
		assert.Equal(t, []any{true}, bindings)
		assert.NoError(t, err)
	})
}

func Test_PluckMapSql_WithGroupBy(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		qb := xqb.Table("orders").SetDialect(dialect)
		qb.GroupBy("customer_id").Having("total", ">", 1000)
		sql, bindings, err := qb.PluckMapSql("total", "customer_id")

		expectedSQL := map[types.Dialect]string{
			types.DialectMySql:     "SELECT `total`, `customer_id` FROM `orders` GROUP BY `customer_id` HAVING `total` > ?",
			types.DialectMariaDB:   "SELECT `total`, `customer_id` FROM `orders` GROUP BY `customer_id` HAVING `total` > ?",
			types.DialectPostgres:  `SELECT "total", "customer_id" FROM "orders" GROUP BY "customer_id" HAVING "total" > $1`,
			types.DialectSQLite:    `SELECT "total", "customer_id" FROM "orders" GROUP BY "customer_id" HAVING "total" > ?`,
			types.DialectSQLServer: "SELECT [total], [customer_id] FROM [orders] GROUP BY [customer_id] HAVING [total] > @p1",
		}

		assert.Equal(t, expectedSQL[dialect], sql)
		assert.Equal(t, []any{1000}, bindings)
		assert.NoError(t, err)
	})
}
