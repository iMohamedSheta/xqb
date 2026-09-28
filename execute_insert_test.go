package xqb_test

import (
	"testing"

	"github.com/iMohamedSheta/xqb"
	xqbErr "github.com/iMohamedSheta/xqb/shared/errors"
	"github.com/iMohamedSheta/xqb/shared/types"
	"github.com/stretchr/testify/assert"
)

func Test_InsertSql_ConsistentOrder(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		values := []map[string]any{
			{
				"email":    "mohamed@gmail.com",
				"name":     "mohamed",
				"age":      20,
				"password": "hashed_password",
			},
			{
				"email":    "ali@gmail.com",
				"age":      21,
				"password": "hashed_password",
				"name":     "ali",
			},
			{
				"age":      22,
				"password": "hashed_password",
				"name":     "ahmed",
				"email":    "ahmed@gmail.com",
			},
		}
		qb := xqb.Table("users").SetDialect(dialect)

		sql, bindings, err := qb.InsertSql(values)

		expectedSQL := map[types.Dialect]string{
			types.DialectMySql:     "INSERT INTO `users` (`age`, `email`, `name`, `password`) VALUES (?, ?, ?, ?), (?, ?, ?, ?), (?, ?, ?, ?)",
			types.DialectMariaDB:   "INSERT INTO `users` (`age`, `email`, `name`, `password`) VALUES (?, ?, ?, ?), (?, ?, ?, ?), (?, ?, ?, ?)",
			types.DialectPostgres:  `INSERT INTO "users" ("age", "email", "name", "password") VALUES ($1, $2, $3, $4), ($5, $6, $7, $8), ($9, $10, $11, $12)`,
			types.DialectSQLite:    `INSERT INTO "users" ("age", "email", "name", "password") VALUES (?, ?, ?, ?), (?, ?, ?, ?), (?, ?, ?, ?)`,
			types.DialectSQLServer: `INSERT INTO [users] ([age], [email], [name], [password]) VALUES (@p1, @p2, @p3, @p4), (@p5, @p6, @p7, @p8), (@p9, @p10, @p11, @p12)`,
		}
		expectedBindings := []any{
			20, "mohamed@gmail.com", "mohamed", "hashed_password",
			21, "ali@gmail.com", "ali", "hashed_password",
			22, "ahmed@gmail.com", "ahmed", "hashed_password",
		}
		assert.Equal(t, expectedSQL[dialect], sql)
		assert.Equal(t, expectedBindings, bindings)
		assert.NoError(t, err)
	})
}

func Test_InsertSql_TakesInsertedColumnsFromFirstRow(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		qb := xqb.Table("users").SetDialect(dialect)
		sql, bindings, err := qb.InsertSql([]map[string]any{
			{
				"name": "mohamed",
			},
			{
				"name": "ali",
				"age":  20,
			},
		})

		expectedSQL := map[types.Dialect]string{
			types.DialectMySql:     "INSERT INTO `users` (`name`) VALUES (?), (?)",
			types.DialectMariaDB:   "INSERT INTO `users` (`name`) VALUES (?), (?)",
			types.DialectPostgres:  `INSERT INTO "users" ("name") VALUES ($1), ($2)`,
			types.DialectSQLite:    `INSERT INTO "users" ("name") VALUES (?), (?)`,
			types.DialectSQLServer: `INSERT INTO [users] ([name]) VALUES (@p1), (@p2)`,
		}
		expectedBindings := []any{
			"mohamed",
			"ali",
		}
		assert.Equal(t, expectedSQL[dialect], sql)
		assert.Equal(t, expectedBindings, bindings)
		assert.Empty(t, err)
	})
}

func Test_InsertSql_NullableColumns(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		qb := xqb.Table("users").SetDialect(dialect)
		sql, bindings, err := qb.InsertSql([]map[string]any{
			{
				// Take the inserted columns from the first row
				"name":  "mohamed",
				"email": nil,
				"age":   20,
			},
			{
				"name":  "ali",
				"email": nil,
				// "age":   20, // Non existing column are inserted as null
			},
		})

		expectedSQL := map[types.Dialect]string{
			types.DialectMySql:     "INSERT INTO `users` (`age`, `email`, `name`) VALUES (?, ?, ?), (?, ?, ?)",
			types.DialectMariaDB:   "INSERT INTO `users` (`age`, `email`, `name`) VALUES (?, ?, ?), (?, ?, ?)",
			types.DialectPostgres:  `INSERT INTO "users" ("age", "email", "name") VALUES ($1, $2, $3), ($4, $5, $6)`,
			types.DialectSQLite:    `INSERT INTO "users" ("age", "email", "name") VALUES (?, ?, ?), (?, ?, ?)`,
			types.DialectSQLServer: `INSERT INTO [users] ([age], [email], [name]) VALUES (@p1, @p2, @p3), (@p4, @p5, @p6)`,
		}

		expectedBindings := []any{
			20, nil, "mohamed",
			nil, nil, "ali",
		}

		assert.Equal(t, expectedSQL[dialect], sql)
		assert.Equal(t, expectedBindings, bindings)
		assert.NoError(t, err)
	})
}

func Test_UpsertSql_WithTwoUniqueByColumns(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		qb := xqb.Table("users").SetDialect(dialect)

		insertedValues := []map[string]any{
			{
				"email":    "mohamed@gmail.com",
				"name":     "mohamed",
				"age":      20,
				"password": "hashed_password",
			},
			{
				"email":    "ali@gmail.com",
				"age":      21,
				"password": "hashed_password",
				"name":     "ali",
			},
			{
				"age":      22,
				"password": "hashed_password",
				"name":     "ahmed",
				"email":    "ahmed@gmail.com",
			},
		}

		sql, bindings, err := qb.UpsertSql(insertedValues, []string{"email", "name"}, []string{"age", "email", "name", "password"})

		expectedSQL := map[types.Dialect]string{
			types.DialectMySql: "INSERT INTO `users` (`age`, `email`, `name`, `password`) VALUES (?, ?, ?, ?), (?, ?, ?, ?), (?, ?, ?, ?) " +
				"ON DUPLICATE KEY UPDATE `age` = VALUES(`age`), `password` = VALUES(`password`)",
			types.DialectMariaDB: "INSERT INTO `users` (`age`, `email`, `name`, `password`) VALUES (?, ?, ?, ?), (?, ?, ?, ?), (?, ?, ?, ?) " +
				"ON DUPLICATE KEY UPDATE `age` = VALUES(`age`), `password` = VALUES(`password`)",
			types.DialectPostgres: `INSERT INTO "users" ("age", "email", "name", "password") VALUES ($1, $2, $3, $4), ($5, $6, $7, $8), ($9, $10, $11, $12) ` +
				`ON CONFLICT ("email", "name") DO UPDATE SET "age" = EXCLUDED."age", "password" = EXCLUDED."password"`,
			types.DialectSQLite: `INSERT INTO "users" ("age", "email", "name", "password") VALUES (?, ?, ?, ?), (?, ?, ?, ?), (?, ?, ?, ?) ` +
				`ON CONFLICT ("email", "name") DO UPDATE SET "age" = EXCLUDED."age", "password" = EXCLUDED."password"`,
			types.DialectSQLServer: "",
		}
		expectedErr := map[types.Dialect]error{
			types.DialectMySql:     nil,
			types.DialectMariaDB:   nil,
			types.DialectPostgres:  nil,
			types.DialectSQLite:    nil,
			types.DialectSQLServer: xqbErr.ErrUnsupportedFeature,
		}
		expectedBindings := []any{
			20, "mohamed@gmail.com", "mohamed", "hashed_password",
			21, "ali@gmail.com", "ali", "hashed_password",
			22, "ahmed@gmail.com", "ahmed", "hashed_password",
		}
		assert.Equal(t, expectedSQL[dialect], sql)
		if expectedErr[dialect] != nil {
			assert.Empty(t, bindings)
			assert.ErrorIs(t, err, expectedErr[dialect])
		} else {
			assert.Equal(t, expectedBindings, bindings)
			assert.NoError(t, err)
		}
	})
}
func Test_UpsertSql_ErrorOnMissingUpdateColumns(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		qb := xqb.Table("users").SetDialect(dialect)

		insertedValues := []map[string]any{
			{
				"email": "khaled@gmail.com",
				"name":  "khaled",
			},
		}

		sql, bindings, err := qb.UpsertSql(insertedValues, []string{"email"}, []string{"name", "password"})

		assert.Empty(t, sql)
		assert.Empty(t, bindings)
		assert.ErrorIs(t, err, xqbErr.ErrInvalidQuery) // cannot update column "password" because it is not part of inserted values
	})
}

func Test_UpsertSql_EmptyValues(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		qb := xqb.Table("users").SetDialect(dialect)

		_, _, err := qb.UpsertSql([]map[string]any{}, []string{"email"}, []string{"name"})
		assert.Error(t, err)
		assert.ErrorIs(t, err, xqbErr.ErrInvalidQuery) // values cannot be empty
	})
}

func Test_UpsertSql_EmptyUpdateColumns(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		qb := xqb.Table("users").SetDialect(dialect)

		_, _, err := qb.UpsertSql([]map[string]any{
			{"email": "mohamed@gmail.com", "name": "mohamed"},
		}, []string{"email"}, []string{})
		assert.Error(t, err)
		assert.ErrorIs(t, err, xqbErr.ErrInvalidQuery) // updateColumns cannot be empty
	})
}

func Test_UpsertSql_UpdateColumnNotInInsert(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		qb := xqb.Table("users").SetDialect(dialect)

		sql, bindings, err := qb.UpsertSql([]map[string]any{
			{"email": "mohamed@gmail.com", "name": "mohamed"},
		}, []string{"email"}, []string{"non_existing_col"})

		assert.Empty(t, sql)
		assert.Empty(t, bindings)
		assert.Error(t, err)
		assert.ErrorIs(t, err, xqbErr.ErrInvalidQuery) // cannot update column "non_existing_col" because it is not part of inserted values
	})
}

func Test_UpsertSql_SkipUniqueByInUpdate(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		qb := xqb.Table("users").SetDialect(dialect)

		values := []map[string]any{
			{
				"email": "mohamed@gmail.com",
				"name":  "mohamed",
				"age":   30,
			},
		}
		sql, bindings, err := qb.UpsertSql(values, []string{"email"}, []string{"email", "age"})

		expectedSQL := map[types.Dialect]string{
			types.DialectMySql:     "INSERT INTO `users` (`age`, `email`, `name`) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE `age` = VALUES(`age`)",
			types.DialectMariaDB:   "INSERT INTO `users` (`age`, `email`, `name`) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE `age` = VALUES(`age`)",
			types.DialectPostgres:  `INSERT INTO "users" ("age", "email", "name") VALUES ($1, $2, $3) ON CONFLICT ("email") DO UPDATE SET "age" = EXCLUDED."age"`,
			types.DialectSQLite:    `INSERT INTO "users" ("age", "email", "name") VALUES (?, ?, ?) ON CONFLICT ("email") DO UPDATE SET "age" = EXCLUDED."age"`,
			types.DialectSQLServer: "",
		}
		expectedErr := map[types.Dialect]error{
			types.DialectMySql:     nil,
			types.DialectMariaDB:   nil,
			types.DialectPostgres:  nil,
			types.DialectSQLite:    nil,
			types.DialectSQLServer: xqbErr.ErrUnsupportedFeature,
		}
		expectedBindings := []any{
			30, "mohamed@gmail.com", "mohamed",
		}
		assert.Equal(t, expectedSQL[dialect], sql)
		if expectedErr[dialect] != nil {
			assert.Empty(t, bindings)
			assert.ErrorIs(t, err, expectedErr[dialect])
		} else {
			assert.Equal(t, expectedBindings, bindings)
			assert.NoError(t, err)
		}
	})
}

func Test_Insert_GetId(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		qb := xqb.Table("users").SetDialect(dialect)
		sql, bindings, err := qb.InsertGetIdSql([]map[string]any{
			{
				"email": "mohamed@gmail.com",
				"name":  "mohamed",
				"age":   20,
			},
		})

		expectedSQL := map[types.Dialect]string{
			types.DialectMySql:     "INSERT INTO `users` (`age`, `email`, `name`) VALUES (?, ?, ?)",
			types.DialectMariaDB:   "INSERT INTO `users` (`age`, `email`, `name`) VALUES (?, ?, ?)",
			types.DialectPostgres:  `INSERT INTO "users" ("age", "email", "name") VALUES ($1, $2, $3) RETURNING id`,
			types.DialectSQLite:    `INSERT INTO "users" ("age", "email", "name") VALUES (?, ?, ?) RETURNING id`,
			types.DialectSQLServer: `INSERT INTO [users] ([age], [email], [name]) OUTPUT INSERTED.[id] VALUES (@p1, @p2, @p3)`,
		}
		expectedBindings := []any{
			20, "mohamed@gmail.com", "mohamed",
		}
		assert.Equal(t, expectedSQL[dialect], sql)
		assert.Equal(t, expectedBindings, bindings)
		assert.NoError(t, err)
	})
}
