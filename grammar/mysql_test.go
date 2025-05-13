package grammar

import (
	"testing"

	"github.com/iMohamedSheta/xqb/types"
	"github.com/stretchr/testify/assert"
)

func TestMySQLGrammar_CompileSelectClause(t *testing.T) {
	tests := []struct {
		name     string
		qb       *types.QueryBuilderData
		expected string
		bindings []interface{}
	}{
		{
			name: "Basic SELECT *",
			qb: &types.QueryBuilderData{
				Columns: []interface{}{},
			},
			expected: "SELECT *",
			bindings: nil,
		},
		{
			name: "SELECT with columns",
			qb: &types.QueryBuilderData{
				Columns: []interface{}{"id", "name", "email"},
			},
			expected: "SELECT id, name, email",
			bindings: nil,
		},
		{
			name: "SELECT DISTINCT",
			qb: &types.QueryBuilderData{
				IsUsingDistinct: true,
				Columns:         []interface{}{"id", "name"},
			},
			expected: "SELECT DISTINCT id, name",
			bindings: nil,
		},
		{
			name: "SELECT with aggregate functions",
			qb: &types.QueryBuilderData{
				AggregateFuncs: []types.AggregateExpr{
					{Function: types.COUNT, Column: "*", Alias: "total"},
					{Function: types.SUM, Column: "amount", Alias: "total_amount"},
				},
			},
			expected: "SELECT COUNT(*) AS total, SUM(amount) AS total_amount",
			bindings: nil,
		},
		{
			name: "SELECT with JSON expressions",
			qb: &types.QueryBuilderData{
				JSONExpressions: []types.JSONExpression{
					{Column: "data", Path: "$.name", Alias: "user_name"},
					{Function: "JSON_UNQUOTE", Column: "data", Path: "$.email", Alias: "user_email"},
				},
			},
			expected: "SELECT JSON_EXTRACT(data, '$.name') AS user_name, JSON_UNQUOTE(data, '$.email') AS user_email",
			bindings: nil,
		},
		{
			name: "SELECT with string functions",
			qb: &types.QueryBuilderData{
				StringFuncs: []types.StringFunction{
					{Function: "CONCAT", Column: "first_name", Params: []interface{}{" ", "last_name"}, Alias: "full_name"},
					{Function: "UPPER", Column: "name", Alias: "upper_name"},
				},
			},
			expected: "SELECT CONCAT(first_name, ?, ?) AS full_name, UPPER(name) AS upper_name",
			bindings: []interface{}{" ", "last_name"},
		},
	}

	grammar := &MySQLGrammar{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sql, bindings, _ := grammar.compileSelectClause(tt.qb)
			assert.Equal(t, tt.expected, sql)
			assert.Equal(t, tt.bindings, bindings)
		})
	}
}

func TestMySQLGrammar_CompileFromClause(t *testing.T) {
	tests := []struct {
		name     string
		qb       *types.QueryBuilderData
		expected string
		bindings []interface{}
	}{
		{
			name: "Basic FROM clause",
			qb: &types.QueryBuilderData{
				Table: "users",
			},
			expected: " FROM users",
			bindings: nil,
		},
		{
			name: "FROM with FORCE INDEX",
			qb: &types.QueryBuilderData{
				Table:      "users",
				ForceIndex: "idx_email",
			},
			expected: " FROM users FORCE INDEX (idx_email)",
			bindings: nil,
		},
		{
			name: "FROM with USE INDEX",
			qb: &types.QueryBuilderData{
				Table:    "users",
				UseIndex: "idx_name",
			},
			expected: " FROM users USE INDEX (idx_name)",
			bindings: nil,
		},
		{
			name: "Empty table name",
			qb: &types.QueryBuilderData{
				Table: "",
			},
			expected: "",
			bindings: nil,
		},
	}

	grammar := &MySQLGrammar{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sql, bindings, _ := grammar.compileFromClause(tt.qb)
			assert.Equal(t, tt.expected, sql)
			assert.Equal(t, tt.bindings, bindings)
		})
	}
}

func TestMySQLGrammar_CompileJoins(t *testing.T) {
	tests := []struct {
		name     string
		qb       *types.QueryBuilderData
		expected string
		bindings []interface{}
	}{
		{
			name: "Single INNER JOIN",
			qb: &types.QueryBuilderData{
				Joins: []types.Join{
					{Type: "INNER JOIN", Table: "orders", Condition: "users.id = orders.user_id"},
				},
			},
			expected: " INNER JOIN orders ON users.id = orders.user_id",
			bindings: nil,
		},
		{
			name: "Multiple JOINs",
			qb: &types.QueryBuilderData{
				Joins: []types.Join{
					{Type: "LEFT JOIN", Table: "orders", Condition: "users.id = orders.user_id"},
					{Type: "INNER JOIN", Table: "order_items", Condition: "orders.id = order_items.order_id"},
				},
			},
			expected: " LEFT JOIN orders ON users.id = orders.user_id INNER JOIN order_items ON orders.id = order_items.order_id",
			bindings: nil,
		},
		{
			name: "No joins",
			qb: &types.QueryBuilderData{
				Joins: []types.Join{},
			},
			expected: "",
			bindings: nil,
		},
	}

	grammar := &MySQLGrammar{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sql, bindings, _ := grammar.compileJoins(tt.qb)
			assert.Equal(t, tt.expected, sql)
			assert.Equal(t, tt.bindings, bindings)
		})
	}
}

func TestMySQLGrammar_CompileWhereClause(t *testing.T) {
	tests := []struct {
		name     string
		qb       *types.QueryBuilderData
		expected string
		bindings []interface{}
	}{
		{
			name: "Simple WHERE condition",
			qb: &types.QueryBuilderData{
				Where: []types.WhereCondition{
					{Column: "age", Operator: ">", Value: 18},
				},
			},
			expected: " WHERE age > ?",
			bindings: []interface{}{18},
		},
		{
			name: "Multiple WHERE conditions with AND",
			qb: &types.QueryBuilderData{
				Where: []types.WhereCondition{
					{Column: "age", Operator: ">", Value: 18},
					{Connector: "AND", Column: "active", Operator: "=", Value: true},
				},
			},
			expected: " WHERE age > ? AND active = ?",
			bindings: []interface{}{18, true},
		},
		{
			name: "IN condition",
			qb: &types.QueryBuilderData{
				Where: []types.WhereCondition{
					{Column: "id", Operator: "IN", Value: []interface{}{1, 2, 3}},
				},
			},
			expected: " WHERE id IN (?, ?, ?)",
			bindings: []interface{}{1, 2, 3},
		},
		{
			name: "BETWEEN condition",
			qb: &types.QueryBuilderData{
				Where: []types.WhereCondition{
					{Column: "age", Operator: "BETWEEN", Value: []interface{}{18, 65}},
				},
			},
			expected: " WHERE age BETWEEN ? AND ?",
			bindings: []interface{}{18, 65},
		},
		{
			name: "Raw SQL condition",
			qb: &types.QueryBuilderData{
				Where: []types.WhereCondition{
					{
						Raw: &types.Expression{
							SQL:      "EXISTS (SELECT 1 FROM orders WHERE orders.user_id = users.id)",
							Bindings: nil,
						},
					},
				},
			},
			expected: " WHERE EXISTS (SELECT 1 FROM orders WHERE orders.user_id = users.id)",
			bindings: nil,
		},
	}

	grammar := &MySQLGrammar{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sql, bindings, _ := grammar.compileWhereClause(tt.qb)
			assert.Equal(t, tt.expected, sql)
			assert.Equal(t, tt.bindings, bindings)
		})
	}
}

func TestMySQLGrammar_CompileGroupByClause(t *testing.T) {
	tests := []struct {
		name     string
		qb       *types.QueryBuilderData
		expected string
		bindings []interface{}
	}{
		{
			name: "Single GROUP BY column",
			qb: &types.QueryBuilderData{
				GroupBy: []string{"user_id"},
			},
			expected: " GROUP BY user_id",
			bindings: nil,
		},
		{
			name: "Multiple GROUP BY columns",
			qb: &types.QueryBuilderData{
				GroupBy: []string{"user_id", "status"},
			},
			expected: " GROUP BY user_id, status",
			bindings: nil,
		},
		{
			name: "No GROUP BY",
			qb: &types.QueryBuilderData{
				GroupBy: []string{},
			},
			expected: "",
			bindings: nil,
		},
	}

	grammar := &MySQLGrammar{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sql, bindings, _ := grammar.compileGroupByClause(tt.qb)
			assert.Equal(t, tt.expected, sql)
			assert.Equal(t, tt.bindings, bindings)
		})
	}
}

func TestMySQLGrammar_CompileHavingClause(t *testing.T) {
	tests := []struct {
		name     string
		qb       *types.QueryBuilderData
		expected string
		bindings []interface{}
	}{
		{
			name: "Single HAVING condition",
			qb: &types.QueryBuilderData{
				Having: []types.Having{
					{Column: "total_amount", Operator: ">", Value: 1000},
				},
			},
			expected: " HAVING total_amount > ?",
			bindings: []interface{}{1000},
		},
		{
			name: "Multiple HAVING conditions",
			qb: &types.QueryBuilderData{
				Having: []types.Having{
					{Column: "total_amount", Operator: ">", Value: 1000},
					{Connector: types.AND, Column: "order_count", Operator: ">=", Value: 5},
				},
			},
			expected: " HAVING total_amount > ? AND order_count >= ?",
			bindings: []interface{}{1000, 5},
		},
		{
			name: "No HAVING",
			qb: &types.QueryBuilderData{
				Having: []types.Having{},
			},
			expected: "",
			bindings: nil,
		},
	}

	grammar := &MySQLGrammar{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sql, bindings, _ := grammar.compileHavingClause(tt.qb)
			assert.Equal(t, tt.expected, sql)
			assert.Equal(t, tt.bindings, bindings)
		})
	}
}

func TestMySQLGrammar_CompileOrderByClause(t *testing.T) {
	tests := []struct {
		name     string
		qb       *types.QueryBuilderData
		expected string
		bindings []interface{}
	}{
		{
			name: "Single ORDER BY column",
			qb: &types.QueryBuilderData{
				OrderBy: []types.OrderBy{
					{Column: "created_at", Direction: "DESC"},
				},
			},
			expected: " ORDER BY created_at DESC",
			bindings: nil,
		},
		{
			name: "Multiple ORDER BY columns",
			qb: &types.QueryBuilderData{
				OrderBy: []types.OrderBy{
					{Column: "status", Direction: "ASC"},
					{Column: "created_at", Direction: "DESC"},
				},
			},
			expected: " ORDER BY status ASC, created_at DESC",
			bindings: nil,
		},
		{
			name: "No ORDER BY",
			qb: &types.QueryBuilderData{
				OrderBy: []types.OrderBy{},
			},
			expected: "",
			bindings: nil,
		},
	}

	grammar := &MySQLGrammar{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sql, bindings, _ := grammar.compileOrderByClause(tt.qb)
			assert.Equal(t, tt.expected, sql)
			assert.Equal(t, tt.bindings, bindings)
		})
	}
}

func TestMySQLGrammar_CompileLimitOffsetClause(t *testing.T) {
	tests := []struct {
		name     string
		qb       *types.QueryBuilderData
		expected string
		bindings []interface{}
	}{
		{
			name: "Only LIMIT",
			qb: &types.QueryBuilderData{
				Limit: 10,
			},
			expected: " LIMIT 10",
			bindings: nil,
		},
		{
			name: "Only OFFSET",
			qb: &types.QueryBuilderData{
				Offset: 20,
			},
			expected: " OFFSET 20",
			bindings: nil,
		},
		{
			name: "Both LIMIT and OFFSET",
			qb: &types.QueryBuilderData{
				Limit:  10,
				Offset: 20,
			},
			expected: " LIMIT 10 OFFSET 20",
			bindings: nil,
		},
		{
			name: "No LIMIT or OFFSET",
			qb: &types.QueryBuilderData{
				Limit:  0,
				Offset: 0,
			},
			expected: "",
			bindings: nil,
		},
	}

	grammar := &MySQLGrammar{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sql, bindings, _ := grammar.compileLimitOffsetClause(tt.qb)
			assert.Equal(t, tt.expected, sql)
			assert.Equal(t, tt.bindings, bindings)
		})
	}
}

func TestMySQLGrammar_CompileLockingClause(t *testing.T) {
	tests := []struct {
		name     string
		qb       *types.QueryBuilderData
		expected string
		bindings []interface{}
	}{
		{
			name: "FOR UPDATE",
			qb: &types.QueryBuilderData{
				IsForUpdate: true,
			},
			expected: " FOR UPDATE",
			bindings: nil,
		},
		{
			name: "LOCK IN SHARE MODE",
			qb: &types.QueryBuilderData{
				IsLockInShareMode: true,
			},
			expected: " LOCK IN SHARE MODE",
			bindings: nil,
		},
		{
			name: "No locking",
			qb: &types.QueryBuilderData{
				IsForUpdate:       false,
				IsLockInShareMode: false,
			},
			expected: "",
			bindings: nil,
		},
	}

	grammar := &MySQLGrammar{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sql, bindings, _ := grammar.compileLockingClause(tt.qb)
			assert.Equal(t, tt.expected, sql)
			assert.Equal(t, tt.bindings, bindings)
		})
	}
}

func TestMySQLGrammar_CompileCTEs(t *testing.T) {
	tests := []struct {
		name     string
		qb       *types.QueryBuilderData
		expected string
		bindings []interface{}
	}{
		{
			name: "Simple CTE",
			qb: &types.QueryBuilderData{
				WithCTEs: []types.CTE{
					{
						Name: "user_orders",
						Expression: &types.Expression{
							SQL:      "SELECT user_id, COUNT(*) AS order_count FROM orders GROUP BY user_id",
							Bindings: nil,
						},
					},
				},
			},
			expected: "WITH user_orders AS (SELECT user_id, COUNT(*) AS order_count FROM orders GROUP BY user_id)",
			bindings: nil,
		},
		{
			name: "Multiple CTEs",
			qb: &types.QueryBuilderData{
				WithCTEs: []types.CTE{
					{
						Name: "active_users",
						Expression: &types.Expression{
							SQL:      "SELECT * FROM users WHERE active = ?",
							Bindings: []interface{}{true},
						},
					},
					{
						Name: "user_stats",
						Expression: &types.Expression{
							SQL:      "SELECT user_id, SUM(amount) AS total FROM orders GROUP BY user_id",
							Bindings: nil,
						},
					},
				},
			},
			expected: "WITH active_users AS (SELECT * FROM users WHERE active = ?), user_stats AS (SELECT user_id, SUM(amount) AS total FROM orders GROUP BY user_id)",
			bindings: []interface{}{true},
		},
		{
			name: "No CTEs",
			qb: &types.QueryBuilderData{
				WithCTEs: []types.CTE{},
			},
			expected: "",
			bindings: nil,
		},
	}

	grammar := &MySQLGrammar{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sql, bindings, _ := grammar.compileCTEs(tt.qb)
			assert.Equal(t, tt.expected, sql)
			assert.Equal(t, tt.bindings, bindings)
		})
	}
}

func TestMySQLGrammar_CompileSelect(t *testing.T) {
	tests := []struct {
		name     string
		qb       *types.QueryBuilderData
		expected string
		bindings []interface{}
	}{
		{
			name: "Simple SELECT",
			qb: &types.QueryBuilderData{
				Table:   "users",
				Columns: []interface{}{"id", "name", "email"},
			},
			expected: "SELECT id, name, email FROM users",
			bindings: nil,
		},
		{
			name: "SELECT with WHERE and ORDER BY",
			qb: &types.QueryBuilderData{
				Table:   "users",
				Columns: []interface{}{"id", "name"},
				Where: []types.WhereCondition{
					{Column: "active", Operator: "=", Value: true},
				},
				OrderBy: []types.OrderBy{
					{Column: "name", Direction: "ASC"},
				},
			},
			expected: "SELECT id, name FROM users WHERE active = ? ORDER BY name ASC",
			bindings: []interface{}{true},
		},
		{
			name: "SELECT with UNION",
			qb: &types.QueryBuilderData{
				Table:   "active_users",
				Columns: []interface{}{"id", "name"},
				Unions: []types.Union{
					{
						All:  true,
						Type: types.UnionTypeUnion,
						Expression: &types.Expression{
							SQL:      "SELECT id, name FROM inactive_users",
							Bindings: nil,
						},
					},
				},
			},
			expected: "SELECT id, name FROM active_users UNION ALL (SELECT id, name FROM inactive_users)",
			bindings: nil,
		},
		{
			name: "Complex SELECT with all clauses",
			qb: &types.QueryBuilderData{
				Table:   "orders",
				Columns: []interface{}{"id", "user_id", "amount"},
				Joins: []types.Join{
					{Type: types.INNER_JOIN, Table: "users", Condition: "orders.user_id = users.id"},
				},
				Where: []types.WhereCondition{
					{Column: "status", Operator: "=", Value: "pending"},
				},
				GroupBy: []string{"user_id"},
				Having: []types.Having{
					{Column: "total_amount", Operator: ">", Value: 1000},
				},
				OrderBy: []types.OrderBy{
					{Column: "total_amount", Direction: "DESC"},
				},
				Limit:  10,
				Offset: 20,
			},
			expected: "SELECT id, user_id, amount FROM orders INNER JOIN users ON orders.user_id = users.id WHERE status = ? GROUP BY user_id HAVING total_amount > ? ORDER BY total_amount DESC LIMIT 10 OFFSET 20",
			bindings: []interface{}{"pending", 1000},
		},
	}

	grammar := &MySQLGrammar{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sql, bindings, _ := grammar.CompileSelect(tt.qb)
			assert.Equal(t, tt.expected, sql)
			assert.Equal(t, tt.bindings, bindings)
		})
	}
}
