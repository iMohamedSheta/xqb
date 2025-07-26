package postgres

import (
	"testing"

	"github.com/iMohamedSheta/xqb/shared/types"
	"github.com/stretchr/testify/assert"
)

func TestPostgresDialect_CompileSelectClause(t *testing.T) {
	tests := []struct {
		name     string
		qb       *types.QueryBuilderData
		expected string
		bindings []any
	}{
		{
			name: `Basic SELECT *`,
			qb: &types.QueryBuilderData{
				Columns: []any{},
			},
			expected: `SELECT *`,
			bindings: nil,
		},
		{
			name: `SELECT with columns`,
			qb: &types.QueryBuilderData{
				Columns: []any{`id`, `name`, `email`},
			},
			expected: `SELECT "id", "name", "email"`,
			bindings: nil,
		},
		{
			name: `SELECT DISTINCT`,
			qb: &types.QueryBuilderData{
				IsUsingDistinct: true,
				Columns:         []any{`id`, `name`},
			},
			expected: `SELECT DISTINCT "id", "name"`,
			bindings: nil,
		},
	}

	dialect := &PostgresDialect{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sql, bindings, _ := dialect.compileSelectClause(tt.qb)
			assert.Equal(t, tt.expected, sql)
			assert.Equal(t, tt.bindings, bindings)
		})
	}
}

func TestPostgresDialect_CompileFromClause(t *testing.T) {
	tests := []struct {
		name     string
		qb       *types.QueryBuilderData
		expected string
		bindings []any
	}{
		{
			name: `Basic FROM clause`,
			qb: &types.QueryBuilderData{
				Table: &types.Table{Name: `users`},
			},
			expected: ` FROM "users"`,
			bindings: nil,
		},
		{
			name: `Empty table name`,
			qb: &types.QueryBuilderData{
				Table: &types.Table{Name: ``},
			},
			expected: ``,
			bindings: nil,
		},
	}

	dialect := &PostgresDialect{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sql, bindings, _ := dialect.compileFromClause(tt.qb)
			assert.Equal(t, tt.expected, sql)
			assert.Equal(t, tt.bindings, bindings)
		})
	}
}

func TestPostgresDialect_CompileJoins(t *testing.T) {
	tests := []struct {
		name     string
		qb       *types.QueryBuilderData
		expected string
		bindings []any
	}{
		{
			name: `Single JOIN`,
			qb: &types.QueryBuilderData{
				Joins: []*types.Join{
					{Type: `JOIN`, Table: `orders`, Condition: `users.id = orders.user_id`},
				},
			},
			expected: ` JOIN "orders" ON users.id = orders.user_id`,
			bindings: nil,
		},
		{
			name: `Multiple JOINs`,
			qb: &types.QueryBuilderData{
				Joins: []*types.Join{
					{Type: `LEFT JOIN`, Table: `orders`, Condition: `users.id = orders.user_id`},
					{Type: `JOIN`, Table: `order_items`, Condition: `orders.id = order_items.order_id`},
				},
			},
			expected: ` LEFT JOIN "orders" ON users.id = orders.user_id JOIN "order_items" ON orders.id = order_items.order_id`,
			bindings: nil,
		},
		{
			name: `No joins`,
			qb: &types.QueryBuilderData{
				Joins: []*types.Join{},
			},
			expected: ``,
			bindings: nil,
		},
	}

	dialect := &PostgresDialect{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sql, bindings, _ := dialect.compileJoins(tt.qb)
			assert.Equal(t, tt.expected, sql)
			assert.Equal(t, tt.bindings, bindings)
		})
	}
}

func TestPostgresDialect_CompileWhereClause(t *testing.T) {
	tests := []struct {
		name     string
		qb       *types.QueryBuilderData
		expected string
		bindings []any
	}{
		{
			name: `Simple WHERE condition`,
			qb: &types.QueryBuilderData{
				Where: []*types.WhereCondition{
					{Column: `age`, Operator: `>`, Value: 18},
				},
			},
			expected: ` WHERE "age" > ?`,
			bindings: []any{18},
		},
		{
			name: `Multiple WHERE conditions with AND`,
			qb: &types.QueryBuilderData{
				Where: []*types.WhereCondition{
					{Column: `age`, Operator: `>`, Value: 18},
					{Connector: `AND`, Column: `active`, Operator: `=`, Value: true},
				},
			},
			expected: ` WHERE "age" > ? AND "active" = ?`,
			bindings: []any{18, true},
		},
		{
			name: `IN condition`,
			qb: &types.QueryBuilderData{
				Where: []*types.WhereCondition{
					{Column: `id`, Operator: `IN`, Value: []any{1, 2, 3}},
				},
			},
			expected: ` WHERE "id" IN (?, ?, ?)`,
			bindings: []any{1, 2, 3},
		},
		{
			name: `BETWEEN condition`,
			qb: &types.QueryBuilderData{
				Where: []*types.WhereCondition{
					{Column: `age`, Operator: `BETWEEN`, Value: []any{18, 65}},
				},
			},
			expected: ` WHERE "age" BETWEEN ? AND ?`,
			bindings: []any{18, 65},
		},
		{
			name: `Raw Sql condition`,
			qb: &types.QueryBuilderData{
				Where: []*types.WhereCondition{
					{
						Raw: &types.Expression{
							Sql:      `EXISTS (SELECT 1 FROM orders WHERE orders.user_id = users.id)`,
							Bindings: nil,
						},
					},
				},
			},
			expected: ` WHERE EXISTS (SELECT 1 FROM orders WHERE orders.user_id = users.id)`,
			bindings: nil,
		},
	}

	dialect := &PostgresDialect{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sql, bindings, _ := dialect.compileWhereClause(tt.qb)
			assert.Equal(t, tt.expected, sql)
			assert.Equal(t, tt.bindings, bindings)
		})
	}
}

func TestPostgresDialect_CompileGroupByClause(t *testing.T) {
	tests := []struct {
		name     string
		qb       *types.QueryBuilderData
		expected string
		bindings []any
	}{
		{
			name: `Single GROUP BY column`,
			qb: &types.QueryBuilderData{
				GroupBy: []string{`user_id`},
			},
			expected: ` GROUP BY "user_id"`,
			bindings: nil,
		},
		{
			name: `Multiple GROUP BY columns`,
			qb: &types.QueryBuilderData{
				GroupBy: []string{`user_id`, `status`},
			},
			expected: ` GROUP BY "user_id", "status"`,
			bindings: nil,
		},
		{
			name: `No GROUP BY`,
			qb: &types.QueryBuilderData{
				GroupBy: []string{},
			},
			expected: ``,
			bindings: nil,
		},
	}

	dialect := &PostgresDialect{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sql, bindings, _ := dialect.compileGroupByClause(tt.qb)
			assert.Equal(t, tt.expected, sql)
			assert.Equal(t, tt.bindings, bindings)
		})
	}
}

func TestPostgresDialect_CompileHavingClause(t *testing.T) {
	tests := []struct {
		name     string
		qb       *types.QueryBuilderData
		expected string
		bindings []any
	}{
		{
			name: `Single HAVING condition`,
			qb: &types.QueryBuilderData{
				Having: []*types.Having{
					{Column: `total_amount`, Operator: `>`, Value: 1000},
				},
			},
			expected: ` HAVING "total_amount" > ?`,
			bindings: []any{1000},
		},
		{
			name: `Multiple HAVING conditions`,
			qb: &types.QueryBuilderData{
				Having: []*types.Having{
					{Column: `total_amount`, Operator: `>`, Value: 1000},
					{Connector: types.AND, Column: `order_count`, Operator: `>=`, Value: 5},
				},
			},
			expected: ` HAVING "total_amount" > ? AND "order_count" >= ?`,
			bindings: []any{1000, 5},
		},
		{
			name: `No HAVING`,
			qb: &types.QueryBuilderData{
				Having: []*types.Having{},
			},
			expected: ``,
			bindings: nil,
		},
	}

	dialect := &PostgresDialect{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sql, bindings, _ := dialect.compileHavingClause(tt.qb)
			assert.Equal(t, tt.expected, sql)
			assert.Equal(t, tt.bindings, bindings)
		})
	}
}

func TestPostgresDialect_CompileOrderByClause(t *testing.T) {
	tests := []struct {
		name     string
		qb       *types.QueryBuilderData
		expected string
		bindings []any
	}{
		{
			name: `Single ORDER BY column`,
			qb: &types.QueryBuilderData{
				OrderBy: []*types.OrderBy{
					{Column: `created_at`, Direction: `DESC`},
				},
			},
			expected: ` ORDER BY "created_at" DESC`,
			bindings: nil,
		},
		{
			name: `Multiple ORDER BY columns`,
			qb: &types.QueryBuilderData{
				OrderBy: []*types.OrderBy{
					{Column: `status`, Direction: `ASC`},
					{Column: `created_at`, Direction: `DESC`},
				},
			},
			expected: ` ORDER BY "status" ASC, "created_at" DESC`,
			bindings: nil,
		},
		{
			name: `No ORDER BY`,
			qb: &types.QueryBuilderData{
				OrderBy: []*types.OrderBy{},
			},
			expected: ``,
			bindings: nil,
		},
	}

	dialect := &PostgresDialect{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sql, bindings, _ := dialect.compileOrderByClause(tt.qb)
			assert.Equal(t, tt.expected, sql)
			assert.Equal(t, tt.bindings, bindings)
		})
	}
}

func TestPostgresDialect_CompileLimitOffsetClause(t *testing.T) {
	tests := []struct {
		name     string
		qb       *types.QueryBuilderData
		expected string
		bindings []any
	}{
		{
			name: `Only LIMIT`,
			qb: &types.QueryBuilderData{
				Limit: 10,
			},
			expected: ` LIMIT 10`,
			bindings: nil,
		},
		{
			name: `Only OFFSET`,
			qb: &types.QueryBuilderData{
				Offset: 20,
			},
			expected: ` OFFSET 20`,
			bindings: nil,
		},
		{
			name: `Both LIMIT and OFFSET`,
			qb: &types.QueryBuilderData{
				Limit:  10,
				Offset: 20,
			},
			expected: ` LIMIT 10 OFFSET 20`,
			bindings: nil,
		},
		{
			name: `No LIMIT or OFFSET`,
			qb: &types.QueryBuilderData{
				Limit:  0,
				Offset: 0,
			},
			expected: ``,
			bindings: nil,
		},
	}

	dialect := &PostgresDialect{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			limitSql, limitBindings, _ := dialect.compileLimitClause(tt.qb)
			offsetSql, offsetBindings, _ := dialect.compileOffsetClause(tt.qb)

			sql := limitSql + offsetSql
			bindings := append(limitBindings, offsetBindings...)

			assert.Equal(t, tt.expected, sql)
			assert.Equal(t, tt.bindings, bindings)
		})
	}
}

func TestPostgresDialect_CompileCTEs(t *testing.T) {
	tests := []struct {
		name     string
		qb       *types.QueryBuilderData
		expected string
		bindings []any
	}{
		{
			name: `Simple CTE`,
			qb: &types.QueryBuilderData{
				WithCTEs: []*types.CTE{
					{
						Name: `user_orders`,
						Expression: &types.Expression{
							Sql:      `SELECT user_id, COUNT(*) AS order_count FROM orders GROUP BY user_id`,
							Bindings: nil,
						},
					},
				},
			},
			expected: `WITH user_orders AS (SELECT user_id, COUNT(*) AS order_count FROM orders GROUP BY user_id) `,
			bindings: nil,
		},
		{
			name: `Multiple CTEs`,
			qb: &types.QueryBuilderData{
				WithCTEs: []*types.CTE{
					{
						Name: `active_users`,
						Expression: &types.Expression{
							Sql:      `SELECT * FROM users WHERE active = ?`,
							Bindings: []any{true},
						},
					},
					{
						Name: `user_stats`,
						Expression: &types.Expression{
							Sql:      `SELECT user_id, SUM(amount) AS total FROM orders GROUP BY user_id`,
							Bindings: nil,
						},
					},
				},
			},
			expected: `WITH active_users AS (SELECT * FROM users WHERE active = ?), user_stats AS (SELECT user_id, SUM(amount) AS total FROM orders GROUP BY user_id) `,
			bindings: []any{true},
		},
		{
			name: `No CTEs`,
			qb: &types.QueryBuilderData{
				WithCTEs: []*types.CTE{},
			},
			expected: ``,
			bindings: nil,
		},
	}

	dialect := &PostgresDialect{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sql, bindings, _ := dialect.compileCTEs(tt.qb)
			assert.Equal(t, tt.expected, sql)
			assert.Equal(t, tt.bindings, bindings)
		})
	}
}

func TestPostgresDialect_CompileSelect(t *testing.T) {
	tests := []struct {
		name     string
		qb       *types.QueryBuilderData
		expected string
		bindings []any
	}{
		{
			name: `Simple SELECT`,
			qb: &types.QueryBuilderData{
				Table:   &types.Table{Name: `users`},
				Columns: []any{`id`, `name`, `email`},
			},
			expected: `SELECT "id", "name", "email" FROM "users"`,
			bindings: nil,
		},
		{
			name: `SELECT with WHERE and ORDER BY`,
			qb: &types.QueryBuilderData{
				Table:   &types.Table{Name: `users`},
				Columns: []any{`id`, `name`},
				Where: []*types.WhereCondition{
					{Column: `active`, Operator: `=`, Value: true},
				},
				OrderBy: []*types.OrderBy{
					{Column: `name`, Direction: `ASC`},
				},
			},
			expected: `SELECT "id", "name" FROM "users" WHERE "active" = ? ORDER BY "name" ASC`,
			bindings: []any{true},
		},
		{
			name: `SELECT with UNION`,
			qb: &types.QueryBuilderData{
				Table:   &types.Table{Name: `active_users`},
				Columns: []any{`id`, `name`},
				Unions: []*types.Union{
					{
						All:  true,
						Type: types.UnionTypeUnion,
						Expression: &types.Expression{
							Sql:      `SELECT id, name FROM inactive_users`,
							Bindings: nil,
						},
					},
				},
			},
			expected: `(SELECT "id", "name" FROM "active_users") UNION ALL (SELECT id, name FROM inactive_users)`,
			bindings: nil,
		},
		{
			name: `Complex SELECT with all clauses`,
			qb: &types.QueryBuilderData{
				Table:   &types.Table{Name: `orders`},
				Columns: []any{`id`, `user_id`, `amount`},
				Joins: []*types.Join{
					{Type: types.INNER_JOIN, Table: `users`, Condition: `orders.user_id = users.id`},
				},
				Where: []*types.WhereCondition{
					{Column: `status`, Operator: `=`, Value: `pending`},
				},
				GroupBy: []string{`user_id`},
				Having: []*types.Having{
					{Column: `total_amount`, Operator: `>`, Value: 1000},
				},
				OrderBy: []*types.OrderBy{
					{Column: `total_amount`, Direction: `DESC`},
				},
				Limit:  10,
				Offset: 20,
			},
			expected: `SELECT "id", "user_id", "amount" FROM "orders" JOIN "users" ON orders.user_id = users.id WHERE "status" = ? GROUP BY "user_id" HAVING "total_amount" > ? ORDER BY "total_amount" DESC LIMIT 10 OFFSET 20`,
			bindings: []any{`pending`, 1000},
		},
	}

	dialect := &PostgresDialect{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sql, bindings, _ := dialect.CompileSelect(tt.qb)
			assert.Equal(t, tt.expected, sql)
			assert.Equal(t, tt.bindings, bindings)
		})
	}
}

func TestPostgresDialect_CompileUpdate(t *testing.T) {
	tests := []struct {
		name     string
		qb       *types.QueryBuilderData
		expected string
		bindings []any
		wantErr  bool
	}{
		{
			name: `Basic update`,
			qb: &types.QueryBuilderData{
				Table: &types.Table{Name: `users`},
				UpdatedBindings: []*types.Binding{
					{Column: `name`, Value: `John Updated`},
					{Column: `email`, Value: `john.updated@example.com`},
				},
				Where: []*types.WhereCondition{
					{Column: `id`, Operator: `=`, Value: 1},
				},
			},
			expected: `UPDATE "users" SET "email" = ?, "name" = ? WHERE "id" = ?`,
			bindings: []any{`john.updated@example.com`, `John Updated`, 1},
			wantErr:  false,
		},
		{
			name: `Update with multiple conditions`,
			qb: &types.QueryBuilderData{
				Table: &types.Table{Name: `users`},
				UpdatedBindings: []*types.Binding{
					{Column: `status`, Value: `inactive`},
				},
				Where: []*types.WhereCondition{
					{Column: `status`, Operator: `=`, Value: `active`},
					{Connector: `AND`, Column: `last_login`, Operator: `<`, Value: `2024-01-01`},
				},
			},
			expected: `UPDATE "users" SET "status" = ? WHERE "status" = ? AND "last_login" < ?`,
			bindings: []any{`inactive`, `active`, `2024-01-01`},
			wantErr:  false,
		},
		{
			name: `Update with limit`,
			qb: &types.QueryBuilderData{
				Table: &types.Table{Name: `users`},
				UpdatedBindings: []*types.Binding{
					{Column: `status`, Value: `verified`},
				},
				Where: []*types.WhereCondition{
					{Column: `status`, Operator: `=`, Value: `pending`},
				},
				Limit: 10,
			},
			expected: `UPDATE "users" SET "status" = ? WHERE "status" = ? LIMIT 10`,
			bindings: []any{`verified`, `pending`},
			wantErr:  false,
		},
		{
			name: `Update with no bindings`,
			qb: &types.QueryBuilderData{
				Table: &types.Table{Name: `users`},
				Where: []*types.WhereCondition{
					{Column: `id`, Operator: `=`, Value: 1},
				},
				UpdatedBindings: nil,
			},
			expected: ``,
			bindings: nil,
			wantErr:  true,
		},
	}

	dialect := &PostgresDialect{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sql, bindings, err := dialect.CompileUpdate(tt.qb)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, sql)
			assert.Equal(t, tt.bindings, bindings)
		})
	}
}

func TestPostgresDialect_CompileDelete(t *testing.T) {
	tests := []struct {
		name     string
		qb       *types.QueryBuilderData
		expected string
		bindings []any
		wantErr  bool
	}{
		{
			name: `Basic delete`,
			qb: &types.QueryBuilderData{
				Table: &types.Table{Name: `users`},
				Bindings: []*types.Binding{
					{Column: `id`, Value: 1},
				},
				Where: []*types.WhereCondition{
					{Column: `id`, Operator: `=`, Value: 1},
				},
			},
			expected: `DELETE FROM "users" WHERE "id" = ?`,
			bindings: []any{1},
			wantErr:  false,
		},
		{
			name: `Delete with multiple conditions`,
			qb: &types.QueryBuilderData{
				Table: &types.Table{Name: `users`},
				Bindings: []*types.Binding{
					{Column: `status`, Value: `inactive`},
				},
				Where: []*types.WhereCondition{
					{Column: `status`, Operator: `=`, Value: `inactive`},
					{Connector: `AND`, Column: `last_login`, Operator: `<`, Value: `2024-01-01`},
				},
			},
			expected: `DELETE FROM "users" WHERE "status" = ? AND "last_login" < ?`,
			bindings: []any{`inactive`, `2024-01-01`},
			wantErr:  false,
		},
		{
			name: `Delete with limit`,
			qb: &types.QueryBuilderData{
				Table: &types.Table{Name: `users`},
				Bindings: []*types.Binding{
					{Column: `status`, Value: `pending`},
				},
				Where: []*types.WhereCondition{
					{Column: `status`, Operator: `=`, Value: `pending`},
				},
				Limit: 10,
			},
			expected: ``,
			bindings: nil,
			wantErr:  true,
		},
		{
			name: `Delete with no where conditions`,
			qb: &types.QueryBuilderData{
				Table: &types.Table{Name: `users`},
				Where: nil,
			},
			expected: ``,
			bindings: nil,
			wantErr:  true,
		},
	}

	dialect := &PostgresDialect{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sql, bindings, err := dialect.CompileDelete(tt.qb)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, sql)
			assert.Equal(t, tt.bindings, bindings)
		})
	}
}
