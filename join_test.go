package xqb_test

import (
	"testing"

	"github.com/iMohamedSheta/xqb"
	xqbErr "github.com/iMohamedSheta/xqb/shared/errors"
	"github.com/iMohamedSheta/xqb/shared/types"
	"github.com/stretchr/testify/assert"
)

func Test_Join_String_Table(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("users").SetDialect(dialect).Join("posts", "users.id = posts.user_id")
		sql, bindings, err := qb.ToSql()
		expectedSql := map[types.Driver]string{
			types.DriverMySql:    "SELECT * FROM `users` JOIN `posts` ON users.id = posts.user_id",
			types.DriverPostgres: `SELECT * FROM "users" JOIN "posts" ON users.id = posts.user_id`,
		}
		assert.Equal(t, expectedSql[dialect], sql)
		assert.Empty(t, bindings)
		assert.NoError(t, err)
	})
}

func Test_Join_With_Bindings(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("users").SetDialect(dialect).Join("posts", "users.id = posts.user_id AND posts.status = ?", "active")
		sql, bindings, err := qb.ToSql()
		expectedSql := map[types.Driver]string{
			types.DriverMySql:    "SELECT * FROM `users` JOIN `posts` ON users.id = posts.user_id AND posts.status = ?",
			types.DriverPostgres: `SELECT * FROM "users" JOIN "posts" ON users.id = posts.user_id AND posts.status = $1`,
		}
		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, []any{"active"}, bindings)
		assert.NoError(t, err)
	})
}

func Test_LeftJoin(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("users").SetDialect(dialect).LeftJoin("comments", "users.id = comments.user_id")
		sql, bindings, err := qb.ToSql()
		expectedSql := map[types.Driver]string{
			types.DriverMySql:    "SELECT * FROM `users` LEFT JOIN `comments` ON users.id = comments.user_id",
			types.DriverPostgres: `SELECT * FROM "users" LEFT JOIN "comments" ON users.id = comments.user_id`,
		}
		assert.Equal(t, expectedSql[dialect], sql)
		assert.NoError(t, err)
		assert.Empty(t, bindings)
	})
}

func Test_RightJoin(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("users").SetDialect(dialect).RightJoin("logins", "users.id = logins.user_id")
		sql, bindings, err := qb.ToSql()
		expectedSql := map[types.Driver]string{
			types.DriverMySql:    "SELECT * FROM `users` RIGHT JOIN `logins` ON users.id = logins.user_id",
			types.DriverPostgres: `SELECT * FROM "users" RIGHT JOIN "logins" ON users.id = logins.user_id`,
		}
		assert.Equal(t, expectedSql[dialect], sql)
		assert.NoError(t, err)
		assert.Empty(t, bindings)
	})
}

func Test_FullJoin(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("users").SetDialect(dialect).FullJoin("sessions", "users.id = sessions.user_id")
		sql, bindings, err := qb.ToSql()
		expectedSql := map[types.Driver]string{
			types.DriverMySql:    "",
			types.DriverPostgres: `SELECT * FROM "users" FULL JOIN "sessions" ON users.id = sessions.user_id`,
		}
		expectedErr := map[types.Driver]error{
			types.DriverMySql:    xqbErr.ErrUnsupportedFeature, // MySql does not support FULL JOIN
			types.DriverPostgres: nil,
		}

		assert.Equal(t, expectedSql[dialect], sql)
		if expectedErr[dialect] != nil {
			assert.ErrorIs(t, err, expectedErr[dialect])
		} else {
			assert.NoError(t, err)
		}

		assert.Empty(t, bindings)
	})
}

func Test_CrossJoin(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("users").SetDialect(dialect).CrossJoin("roles")
		sql, bindings, err := qb.ToSql()
		expectedSql := map[types.Driver]string{
			types.DriverMySql:    "SELECT * FROM `users` CROSS JOIN `roles`",
			types.DriverPostgres: `SELECT * FROM "users" CROSS JOIN "roles"`,
		}
		assert.Equal(t, expectedSql[dialect], sql)
		assert.NoError(t, err)
		assert.Empty(t, bindings)
	})
}

func Test_CrossJoinSub(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		sub := xqb.Table("plans").SetDialect(dialect).Where("expired", "=", false)
		qb := xqb.Table("users").SetDialect(dialect).CrossJoinSub(sub, "p")
		sql, bindings, err := qb.ToSql()
		expectedSql := map[types.Driver]string{
			types.DriverMySql:    "SELECT * FROM `users` CROSS JOIN (SELECT * FROM `plans` WHERE `expired` = ?) AS `p`",
			types.DriverPostgres: `SELECT * FROM "users" CROSS JOIN (SELECT * FROM "plans" WHERE "expired" = $1) AS "p"`,
		}
		assert.Equal(t, expectedSql[dialect], sql)
		assert.NoError(t, err)
		assert.Equal(t, []any{false}, bindings)
	})
}

func Test_CrossJoin_With_Expr(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		raw := xqb.Raw("(SELECT * FROM regions WHERE active = ?) AS r", true)
		qb := xqb.Table("users").SetDialect(dialect).CrossJoinExpr(raw)
		sql, bindings, err := qb.ToSql()
		expectedSql := map[types.Driver]string{
			types.DriverMySql:    "SELECT * FROM `users` CROSS JOIN (SELECT * FROM regions WHERE active = ?) AS `r`",
			types.DriverPostgres: `SELECT * FROM "users" CROSS JOIN (SELECT * FROM regions WHERE active = $1) AS "r"`,
		}
		assert.Equal(t, expectedSql[dialect], sql)
		assert.NoError(t, err)
		assert.Equal(t, []any{true}, bindings)
	})
}

func Test_Join_SubQuery_DefaultAlias(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		sub := xqb.Table("posts").SetDialect(dialect).Where("published", "=", true)
		qb := xqb.Table("users").SetDialect(dialect).JoinSub(sub, "sub", "users.id = sub.user_id")
		sql, bindings, err := qb.ToSql()
		expectedSql := map[types.Driver]string{
			types.DriverMySql:    "SELECT * FROM `users` JOIN (SELECT * FROM `posts` WHERE `published` = ?) AS `sub` ON users.id = sub.user_id",
			types.DriverPostgres: `SELECT * FROM "users" JOIN (SELECT * FROM "posts" WHERE "published" = $1) AS "sub" ON users.id = sub.user_id`,
		}
		assert.Equal(t, expectedSql[dialect], sql)
		assert.NoError(t, err)
		assert.Equal(t, []any{true}, bindings)
	})
}

func Test_Join_SubQuery_With_Alias(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		sub := xqb.Table("posts").SetDialect(dialect).Where("published", "=", true)
		qb := xqb.Table("users").SetDialect(dialect).JoinSub(sub, "p", "users.id = p.user_id")
		sql, bindings, err := qb.ToSql()
		expectedSql := map[types.Driver]string{
			types.DriverMySql:    "SELECT * FROM `users` JOIN (SELECT * FROM `posts` WHERE `published` = ?) AS `p` ON users.id = p.user_id",
			types.DriverPostgres: `SELECT * FROM "users" JOIN (SELECT * FROM "posts" WHERE "published" = $1) AS "p" ON users.id = p.user_id`,
		}
		assert.Equal(t, expectedSql[dialect], sql)
		assert.NoError(t, err)
		assert.Equal(t, []any{true}, bindings)
	})
}

func Test_LeftJoin_SubQuery(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		sub := xqb.Table("comments").SetDialect(dialect).Where("active", "=", true)
		qb := xqb.Table("users").SetDialect(dialect).LeftJoinSub(sub, "c", "users.id = c.user_id")
		sql, bindings, err := qb.ToSql()
		expectedSql := map[types.Driver]string{
			types.DriverMySql:    "SELECT * FROM `users` LEFT JOIN (SELECT * FROM `comments` WHERE `active` = ?) AS `c` ON users.id = c.user_id",
			types.DriverPostgres: `SELECT * FROM "users" LEFT JOIN (SELECT * FROM "comments" WHERE "active" = $1) AS "c" ON users.id = c.user_id`,
		}
		assert.Equal(t, expectedSql[dialect], sql)
		assert.NoError(t, err)
		assert.Equal(t, []any{true}, bindings)
	})
}

func Test_RightJoin_SubQuery(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		sub := xqb.Table("orders").SetDialect(dialect).Where("status", "=", "paid")
		qb := xqb.Table("users").SetDialect(dialect).RightJoinSub(sub, "o", "users.id = o.user_id")
		sql, bindings, err := qb.ToSql()
		expectedSql := map[types.Driver]string{
			types.DriverMySql:    "SELECT * FROM `users` RIGHT JOIN (SELECT * FROM `orders` WHERE `status` = ?) AS `o` ON users.id = o.user_id",
			types.DriverPostgres: `SELECT * FROM "users" RIGHT JOIN (SELECT * FROM "orders" WHERE "status" = $1) AS "o" ON users.id = o.user_id`,
		}
		assert.Equal(t, expectedSql[dialect], sql)
		assert.NoError(t, err)
		assert.Equal(t, []any{"paid"}, bindings)
	})
}

func Test_Join_With_Condition_Expression(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("users").SetDialect(dialect).Join("posts", "users.id = posts.user_id AND posts.status = ?", "active")
		sql, bindings, err := qb.ToSql()
		expectedSql := map[types.Driver]string{
			types.DriverMySql:    "SELECT * FROM `users` JOIN `posts` ON users.id = posts.user_id AND posts.status = ?",
			types.DriverPostgres: `SELECT * FROM "users" JOIN "posts" ON users.id = posts.user_id AND posts.status = $1`,
		}
		assert.Equal(t, expectedSql[dialect], sql)
		assert.NoError(t, err)
		assert.Equal(t, []any{"active"}, bindings)
	})
}

func Test_Join_With_Expression_Table(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		table := xqb.Raw("(SELECT * FROM posts WHERE published = ?) AS p", true)
		qb := xqb.Table("users").SetDialect(dialect).JoinExpr(table, "users.id = p.user_id")
		sql, bindings, err := qb.ToSql()
		expectedSql := map[types.Driver]string{
			types.DriverMySql:    "SELECT * FROM `users` JOIN (SELECT * FROM posts WHERE published = ?) AS `p` ON users.id = p.user_id",
			types.DriverPostgres: `SELECT * FROM "users" JOIN (SELECT * FROM posts WHERE published = $1) AS "p" ON users.id = p.user_id`,
		}
		assert.Equal(t, expectedSql[dialect], sql)
		assert.NoError(t, err)
		assert.Equal(t, []any{true}, bindings)
	})
}

func Test_FullJoinExpr(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		expr := xqb.Raw("(SELECT * FROM stats WHERE active = ?) AS s", true)
		qb := xqb.Table("users").SetDialect(dialect).FullJoinExpr(expr, "users.id = s.user_id")
		sql, bindings, err := qb.ToSql()

		expectedSql := map[types.Driver]string{
			types.DriverMySql:    "",
			types.DriverPostgres: `SELECT * FROM "users" FULL JOIN (SELECT * FROM stats WHERE active = $1) AS "s" ON users.id = s.user_id`,
		}
		expectedErr := map[types.Driver]error{
			types.DriverMySql:    xqbErr.ErrUnsupportedFeature, // MySql does not support FULL JOIN
			types.DriverPostgres: nil,
		}

		assert.Equal(t, expectedSql[dialect], sql)
		if expectedErr[dialect] != nil {
			assert.ErrorIs(t, err, expectedErr[dialect])
			assert.Empty(t, bindings)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, []any{true}, bindings)
		}

	})
}

func Test_JoinSub_ErrorOnMissingSubQueryAlias(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		sub := xqb.Table("posts").SetDialect(dialect).Where("published", "=", true)
		qb := xqb.Table("users").SetDialect(dialect).JoinSub(sub, "", "users.id = sub.user_id")
		sql, bindings, err := qb.ToSql()

		assert.Equal(t, "", sql)
		assert.ErrorIs(t, err, xqbErr.ErrInvalidQuery) // missing subquery alias should result in an error
		assert.Empty(t, bindings)
	})
}

func Test_JoinExpr_With_Expression_Condition(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		table := xqb.Raw("(SELECT * FROM payments WHERE confirmed = ?) AS p", true)
		cond := xqb.Raw("users.id = p.user_id AND p.status = ?", "success")
		qb := xqb.Table("users").SetDialect(dialect).JoinExpr(table, cond)
		sql, bindings, err := qb.ToSql()
		expectedSql := map[types.Driver]string{
			types.DriverMySql:    "SELECT * FROM `users` JOIN (SELECT * FROM payments WHERE confirmed = ?) AS `p` ON users.id = p.user_id AND p.status = ?",
			types.DriverPostgres: `SELECT * FROM "users" JOIN (SELECT * FROM payments WHERE confirmed = $1) AS "p" ON users.id = p.user_id AND p.status = $2`,
		}
		assert.Equal(t, expectedSql[dialect], sql)
		assert.NoError(t, err)
		assert.Equal(t, []any{true, "success"}, bindings)
	})
}

func Test_CrossJoinSub_With_Alias(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		sub := xqb.Table("plans").SetDialect(dialect).Where("expired", "=", false)
		qb := xqb.Table("users").SetDialect(dialect).CrossJoinSub(sub, "sub")
		sql, bindings, err := qb.ToSql()
		expectedSql := map[types.Driver]string{
			types.DriverMySql:    "SELECT * FROM `users` CROSS JOIN (SELECT * FROM `plans` WHERE `expired` = ?) AS `sub`",
			types.DriverPostgres: `SELECT * FROM "users" CROSS JOIN (SELECT * FROM "plans" WHERE "expired" = $1) AS "sub"`,
		}
		assert.Equal(t, expectedSql[dialect], sql)
		assert.NoError(t, err)
		assert.Equal(t, []any{false}, bindings)
	})
}

func Test_Multiple_Joins_Mixed_Types(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		sub := xqb.Table("orders").SetDialect(dialect).Where("status", "=", "shipped")
		expr := xqb.Raw("(SELECT * FROM invoices WHERE paid = ?) AS inv", true)
		qb := xqb.Table("users").SetDialect(dialect).
			Join("addresses", "users.id = addresses.user_id AND addresses.city = ?", "Cairo").
			LeftJoinSub(sub, "o", "users.id = o.user_id").
			RightJoinExpr(expr, "users.id = inv.user_id AND inv.total > ?", 1000)

		sql, bindings, err := qb.ToSql()
		expectedSql := map[types.Driver]string{
			types.DriverMySql: "SELECT * FROM `users` JOIN `addresses` ON users.id = addresses.user_id AND addresses.city = ?" +
				" LEFT JOIN (SELECT * FROM `orders` WHERE `status` = ?) AS `o` ON users.id = o.user_id" +
				" RIGHT JOIN (SELECT * FROM invoices WHERE paid = ?) AS `inv` ON users.id = inv.user_id AND inv.total > ?",
			types.DriverPostgres: `SELECT * FROM "users" JOIN "addresses" ON users.id = addresses.user_id AND addresses.city = $1` +
				` LEFT JOIN (SELECT * FROM "orders" WHERE "status" = $2) AS "o" ON users.id = o.user_id` +
				` RIGHT JOIN (SELECT * FROM invoices WHERE paid = $3) AS "inv" ON users.id = inv.user_id AND inv.total > $4`,
		}

		assert.Equal(t, expectedSql[dialect], sql)
		assert.NoError(t, err)
		assert.Equal(t, []any{"Cairo", "shipped", true, 1000}, bindings)
	})
}

func Test_Join_With_SubQuery_That_Has_Join(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		innerSub := xqb.Table("payments").Where("amount", ">", 500)
		sub := xqb.Table("orders").SetDialect(dialect).
			JoinSub(innerSub, "pay", "orders.payment_id = pay.id").
			Where("orders.status", "=", "completed")

		qb := xqb.Table("users").SetDialect(dialect).JoinSub(sub, "o", "users.id = o.user_id")
		sql, bindings, err := qb.ToSql()

		expectedSql := map[types.Driver]string{
			types.DriverMySql: "SELECT * FROM `users` JOIN (SELECT * FROM `orders` JOIN (SELECT * FROM `payments` WHERE `amount` > ?) AS `pay`" +
				" ON orders.payment_id = pay.id WHERE `orders`.`status` = ?) AS `o` ON users.id = o.user_id",
			types.DriverPostgres: `SELECT * FROM "users" JOIN (SELECT * FROM "orders" JOIN (SELECT * FROM "payments" WHERE "amount" > $1) AS "pay"` +
				` ON orders.payment_id = pay.id WHERE "orders"."status" = $2) AS "o" ON users.id = o.user_id`,
		}
		assert.Equal(t, expectedSql[dialect], sql)
		assert.NoError(t, err)
		assert.Equal(t, []any{500, "completed"}, bindings)
	})
}

func Test_CrossJoin_Combined_With_Other_Joins(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		qb := xqb.Table("users").SetDialect(dialect).
			Join("posts", "users.id = posts.user_id").
			CrossJoin("countries")
		sql, bindings, err := qb.ToSql()

		expectedSql := map[types.Driver]string{
			types.DriverMySql:    "SELECT * FROM `users` JOIN `posts` ON users.id = posts.user_id CROSS JOIN `countries`",
			types.DriverPostgres: `SELECT * FROM "users" JOIN "posts" ON users.id = posts.user_id CROSS JOIN "countries"`,
		}

		assert.Equal(t, expectedSql[dialect], sql)
		assert.NoError(t, err)
		assert.Empty(t, bindings)
	})
}

func Test_Stores_With_Orders_SubQuery(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Driver) {
		subQuery := xqb.Table("orders").SetDialect(dialect).
			Select("store_id", xqb.Raw("COUNT(*) as total_orders")).
			WhereNull("cancelled_at").
			WhereNotNull("confirmed_at").
			Where("status", "!=", "failed").
			GroupBy("store_id")

		sql, bindings, err := xqb.Table("stores").SetDialect(dialect).
			Select(
				"managers.fullname",
				"managers.email",
				"stores.id",
				"order_stats.total_orders",
			).
			AddSelectRaw("locations.city location_city").
			AddSelectRaw("locations.zip_code location_zip").
			AddSelectRaw("managers.id manager_id").
			LeftJoinSub(subQuery, "order_stats", "stores.id = order_stats.store_id").
			Join("managers", "stores.manager_id = managers.id").
			Join("locations", "stores.location_id = locations.id").
			OrderBy("stores.id", "ASC").
			Where("stores.region_id", "=", 22).
			Limit(5).
			ToSql()

		expectedSql := map[types.Driver]string{
			types.DriverMySql: "SELECT `managers`.`fullname`, `managers`.`email`, `stores`.`id`, `order_stats`.`total_orders`, locations.city location_city, locations.zip_code location_zip, managers.id manager_id " +
				"FROM `stores` " +
				"LEFT JOIN (SELECT `store_id`, COUNT(*) AS total_orders FROM `orders` WHERE `cancelled_at` IS NULL AND `confirmed_at` IS NOT NULL AND `status` != ? GROUP BY `store_id`) AS `order_stats` ON stores.id = order_stats.store_id " +
				"JOIN `managers` ON stores.manager_id = managers.id " +
				"JOIN `locations` ON stores.location_id = locations.id " +
				"WHERE `stores`.`region_id` = ? " +
				"ORDER BY `stores`.`id` ASC " +
				"LIMIT 5",
			types.DriverPostgres: `SELECT "managers"."fullname", "managers"."email", "stores"."id", "order_stats"."total_orders", locations.city location_city, locations.zip_code location_zip, managers.id manager_id ` +
				`FROM "stores" ` +
				`LEFT JOIN (SELECT "store_id", COUNT(*) AS total_orders FROM "orders" WHERE "cancelled_at" IS NULL AND "confirmed_at" IS NOT NULL AND "status" != $1 GROUP BY "store_id") AS "order_stats" ON stores.id = order_stats.store_id ` +
				`JOIN "managers" ON stores.manager_id = managers.id ` +
				`JOIN "locations" ON stores.location_id = locations.id ` +
				`WHERE "stores"."region_id" = $2 ` +
				`ORDER BY "stores"."id" ASC ` +
				`LIMIT 5`,
		}

		assert.Equal(t, expectedSql[dialect], sql)
		assert.NoError(t, err)
		assert.Equal(t, []any{"failed", 22}, bindings)
	})
}
