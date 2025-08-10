# XQB Query Builder

A powerful and flexible Sql query builder for Go with fluent interface for building complex Sql queries.

## Installation

```bash
go get github.com/iMohamedSheta/xqb
```

## Quick Start

```go
package main

import (
    "database/sql"
    _ "github.com/go-sql-dialect/mysql"
    "github.com/iMohamedSheta/xqb"
)

func main() {
    // Setup database connection
    db, _ := sql.Open("mysql", "user:password@tcp(localhost:3306)/database")

    xqb.AddConnection(&xqb.Connection{
      Name:    "default", // Default connection name
      Dialect: xqb.DialectMySql,
      DB:      db,
    })

    // Or if you want different connection name and set it as default connection
    myDefaultConnection := "my_connection"
    xqb.AddConnection(&xqb.Connection{
      Name:     myDefaultConnection,
      Dialect: xqb.DialectMySql,
      DB:      db,
    })
    // Set default connection as my_connection
    xqb.SetDefaultConnection(myDefaultConnection)

    // Build and execute query
    qb := xqb.Table("users").
        Select("id", "name", "email").
        Where("active", "=", true).
        OrderBy("name", "ASC").
        Limit(10)

    results, _ := qb.Get()
    // Process results...
}
```

## Query Hooks

XQB supports query hooks for logging, profiling, or custom behavior.

### OnBeforeQuery

Executes right before a query is built.

```go
xqb.DefaultSettings().OnBeforeQuery(func(qb *xqb.QueryBuilder) {
    fmt.Println("Before Query:", qb.GetTable().Name)
})
```

### OnAfterQuery

Executes right after a query is built.

```go
xqb.DefaultSettings().OnAfterQuery(func(query *xqb.QueryExecuted) {
    sql, _ := xqb.InjectBindings(query.Dialect, query.Sql, query.Bindings)
    fmt.Printf("[%s] %s\n", query.Time, sql)
})
```

### OnAfterQueryExecution

Execute right after a query is executed.

```go
xqb.DefaultSettings().OnAfterQueryExecution(func(ctx context.Context) {
    reqId, _ := q.Context.Value(enums.ContextKeyRequestId.String()).(string)
    fmt.Printf("request_id: %s", reqId)
})
```

### Instance-based Hooks

Hooks can be set globally or per query using `WithSettings()`.

## Raw Sql Expressions

### Raw Function

```go
// Raw Sql in SELECT
qb := xqb.Table("users").
    Select(
        xqb.Raw("COUNT(*) as total"),
        "name",
        xqb.Raw("CONCAT(first_name, ' ', last_name) as full_name"),
    )
// Sql: SELECT COUNT(*) as total, name, CONCAT(first_name, ' ', last_name) as full_name FROM users

// Raw Sql in WHERE
qb := xqb.Table("users").
    Where(xqb.Raw("LOWER(email)"), "LIKE", "%@example.com")
// Sql: SELECT * FROM users WHERE LOWER(email) LIKE ?

// Raw Sql in ORDER BY
qb := xqb.Table("orders").
    OrderBy(xqb.Raw("DATE_FORMAT(created_at, '%Y-%m')"), "ASC")
// Sql: SELECT * FROM orders ORDER BY DATE_FORMAT(created_at, '%Y-%m') ASC

// Raw Sql in GROUP BY
qb := xqb.Table("orders").
    GroupBy(xqb.Raw("DATE_FORMAT(created_at, '%Y-%m')"))
// Sql: SELECT * FROM orders GROUP BY DATE_FORMAT(created_at, '%Y-%m')

// Raw Sql in HAVING
qb := xqb.Table("orders").
    GroupBy("user_id").
    Having(xqb.Raw("SUM(amount)"), ">", 1000)
// Sql: SELECT * FROM orders GROUP BY user_id HAVING SUM(amount) > ?
```

### RawDialect Function

```go
// Database-specific expressions
expr := xqb.RawDialect("mysql", map[string]*xqb.Expression{
    "mysql":    xqb.Raw("DATE_FORMAT(created_at, '%Y-%m-%d')"),
    "postgres": xqb.Raw("TO_CHAR(created_at, 'YYYY-MM-DD')"),
})

qb := xqb.Table("users").
    Select(expr, "formatted_date")
// MySql: SELECT DATE_FORMAT(created_at, '%Y-%m-%d') AS formatted_date FROM users
// PostgreSql: SELECT TO_CHAR(created_at, 'YYYY-MM-DD') AS formatted_date FROM users

// Complex dialect-specific expressions
jsonExpr := xqb.RawDialect("mysql", map[string]*xqb.Expression{
    "mysql":    xqb.Raw("JSON_EXTRACT(data, '$.user.email')"),
    "postgres": xqb.Raw("data->'user'->>'email'"),
})

qb := xqb.Table("profiles").
    Select(jsonExpr, "user_email")
// MySql: SELECT JSON_EXTRACT(data, '$.user.email') AS user_email FROM profiles
// PostgreSql: SELECT data->'user'->>'email' AS user_email FROM profiles
```

## Database Connection

```go
// Add connection
db, _ := sql.Open("mysql", "dsn")
xqb.AddConnection("default", db)

// Use specific connection
qb := xqb.Table("users").Connection("default").Where("active", "=", true)

// Close connection
xqb.Close("default")
xqb.CloseAll()
```

## SELECT Queries

### Basic Select

```go
// Simple select
qb := xqb.Table("users").
    Select("id", "name", "email")
// Sql: SELECT id, name, email FROM users

// Select with conditions
qb := xqb.Table("users").
    Select("id", "name").
    Where("age", ">", 18)
// Sql: SELECT id, name FROM users WHERE age > ?

// Select distinct
qb := xqb.Table("users").
    Select("id", "name").
    Distinct()
// Sql: SELECT DISTINCT id, name FROM users
```

### Joins

```go
// Inner join
qb := xqb.Table("users").
    Select("users.id", "users.name", "orders.id as order_id").
    Join("orders", "users.id = orders.user_id")
// Sql: SELECT users.id, users.name, orders.id as order_id FROM users JOIN orders ON users.id = orders.user_id

// Left join
qb := xqb.Table("users").
    LeftJoin("comments", "users.id = comments.user_id")
// Sql: SELECT * FROM users LEFT JOIN comments ON users.id = comments.user_id

// Right join
qb := xqb.Table("users").
    RightJoin("logins", "users.id = logins.user_id")
// Sql: SELECT * FROM users RIGHT JOIN logins ON users.id = logins.user_id

// Full join
qb := xqb.Table("users").
    FullJoin("sessions", "users.id = sessions.user_id")
// Sql: SELECT * FROM users FULL JOIN sessions ON users.id = sessions.user_id

// Cross join
qb := xqb.Table("users").
    CrossJoin("roles")
// Sql: SELECT * FROM users CROSS JOIN roles

// Join with conditions
qb := xqb.Table("users").
    Join("posts", "users.id = posts.user_id AND posts.status = ?", "active")
// Sql: SELECT * FROM users JOIN posts ON users.id = posts.user_id AND posts.status = ?
```

### Subquery Joins

```go
// Join subquery
sub := xqb.Table("posts").Where("published", "=", true)
qb := xqb.Table("users").
    JoinSub(sub, "p", "users.id = p.user_id")
// Sql: JOIN (SELECT * FROM posts WHERE published = ?) AS p ON users.id = p.user_id

// Left join subquery
sub := xqb.Table("comments").Where("active", "=", true)
qb := xqb.Table("users").
    LeftJoinSub(sub, "c", "users.id = c.user_id")
// Sql: LEFT JOIN (SELECT * FROM comments WHERE active = ?) AS c ON users.id = c.user_id

// Cross join subquery
sub := xqb.Table("plans").Where("expired", "=", false)
qb := xqb.Table("users").
    CrossJoinSub(sub, "p")
// Sql: CROSS JOIN (SELECT * FROM plans WHERE expired = ?) AS p
```

### Where Conditions

```go
// Basic where
qb := xqb.Table("users").
    Where("age", ">", 18)
// Sql: SELECT * FROM users WHERE age > ?

// Multiple conditions
qb := xqb.Table("users").
    Where("age", ">", 18).
    Where("active", "=", true)
// Sql: SELECT * FROM users WHERE age > ? AND active = ?

// OR conditions
qb := xqb.Table("users").
    Where("id", "=", 1).
    OrWhere("email", "=", "admin@example.com")
// Sql: SELECT * FROM users WHERE id = ? OR email = ?

// IN conditions
qb := xqb.Table("users").
    WhereIn("id", []any{1, 2, 3})
// Sql: SELECT * FROM users WHERE id IN (?, ?, ?)

// BETWEEN conditions
qb := xqb.Table("users").
    WhereBetween("age", 18, 65)
// Sql: SELECT * FROM users WHERE age BETWEEN ? AND ?

// NULL conditions
qb := xqb.Table("users").
    WhereNull("deleted_at")
// Sql: SELECT * FROM users WHERE deleted_at IS NULL

// EXISTS conditions
subQuery := xqb.Table("orders").Select("user_id").Where("status", "=", "active")
qb := xqb.Table("users").
    WhereExists(subQuery)
// Sql: SELECT * FROM users WHERE EXISTS (SELECT user_id FROM orders WHERE status = ?)

// Raw where
qb := xqb.Table("users").
    WhereRaw("CASE WHEN status = 'active' THEN 1 ELSE 0 END = ?", 1)
// Sql: SELECT * FROM users WHERE CASE WHEN status = 'active' THEN 1 ELSE 0 END = ?

// Where groups
qb := xqb.Table("users").
    Where("id", "=", 1).
    WhereGroup(func(qb *xqb.QueryBuilder) {
        qb.WhereNull("deleted_at").OrWhereNull("disabled_at")
    })
// Sql: SELECT * FROM users WHERE id = ? AND (deleted_at IS NULL OR disabled_at IS NULL)
```

### Group By and Having

```go
// Group by
qb := xqb.Table("orders").
    Select("user_id", "COUNT(*) as order_count").
    GroupBy("user_id")
// Sql: SELECT user_id, COUNT(*) as order_count FROM orders GROUP BY user_id

// Having
qb := xqb.Table("orders").
    Select("user_id", "COUNT(*) as order_count").
    GroupBy("user_id").
    Having("order_count", ">", 5)
// Sql: SELECT user_id, COUNT(*) as order_count FROM orders GROUP BY user_id HAVING order_count > ?
```

### Order By

```go
// Order by
qb := xqb.Table("users").
    Select("id", "name").
    OrderBy("name", "ASC")
// Sql: SELECT id, name FROM users ORDER BY name ASC

// Multiple order by
qb := xqb.Table("users").
    OrderBy("age", "DESC").
    OrderBy("name", "ASC")
// Sql: SELECT * FROM users ORDER BY age DESC, name ASC
```

### Limit and Offset

```go
// Limit
qb := xqb.Table("users").
    Select("id", "name").
    Limit(10)
// Sql: SELECT id, name FROM users LIMIT 10

// Limit with offset
qb := xqb.Table("users").
    Select("id", "name").
    Limit(10).
    Offset(20)
// Sql: SELECT id, name FROM users LIMIT 10 OFFSET 20
```

### Common Table Expressions (CTE)

```go
// Simple CTE
qb := xqb.Table("users").
    WithRaw("user_totals", "SELECT user_id, SUM(amount) as total_spent FROM orders GROUP BY user_id").
    Select("users.id", "users.name", "user_totals.total_spent").
    Join("user_totals", "users.id = user_totals.user_id")
// Sql: WITH user_totals AS (SELECT user_id, SUM(amount) as total_spent FROM orders GROUP BY user_id) SELECT users.id, users.name, user_totals.total_spent FROM users JOIN user_totals ON users.id = user_totals.user_id

// Complex CTE
qb := xqb.Table("products").
    WithRaw("active_users",
        "WITH user_orders AS (SELECT user_id, COUNT(*) as order_count FROM orders GROUP BY user_id) "+
            "SELECT users.id, users.name, user_orders.order_count FROM users "+
            "JOIN user_orders ON users.id = user_orders.user_id").
    Select("products.id", "products.name", "active_users.name as buyer").
    Join("active_users", "products.id = active_users.id")
```

### Locking

```go
// Lock for update
qb := xqb.Table("users").
    Select("id", "name").
    LockForUpdate()
// Sql: SELECT id, name FROM users FOR UPDATE

// Shared lock
qb := xqb.Table("users").
    Select("id", "name").
    SharedLock()
// Sql: SELECT id, name FROM users LOCK IN SHARE MODE
```

## Aggregate Functions

```go
// Basic aggregates
qb := xqb.Table("orders").
    Select(
        xqb.Count("id", "order_count"),
        xqb.Sum("amount", "total_amount"),
        xqb.Avg("amount", "average_amount"),
        xqb.Min("amount", "min_amount"),
        xqb.Raw("MAX(amount) AS max_amount")
    )
// Sql: SELECT COUNT(id) AS order_count, SUM(amount) AS total_amount, AVG(amount) AS average_amount, MIN(amount) AS min_amount, MAX(amount) AS max_amount FROM orders
```

## String Functions

```go
// String operations
qb := xqb.Table("users").
    Select(
        xqb.Concat([]string{"first_name", "' '", "last_name"}, "full_name"),
        xqb.Lower("email", "lower_email"),
        xqb.Upper("username", "upper_username"),
        xqb.Length("bio", "bio_length"),
        xqb.Trim("nickname", "trimmed_nickname"),
        xqb.Replace("title", "'foo'", "'bar'", "replaced_title"),
        xqb.Substring("description", 1, 10, "short_desc"),
    )
// Sql: SELECT CONCAT(first_name, ' ', last_name) AS full_name, LOWER(email) AS lower_email, UPPER(username) AS upper_username, LENGTH(bio) AS bio_length, TRIM(nickname) AS trimmed_nickname, REPLACE(title, 'foo', 'bar') AS replaced_title, SUBSTRING(description, 1, 10) AS short_desc FROM users
```

## Date Functions

```go
// Date operations
qb := xqb.Table("events").
    Select(
        xqb.Date("created_at", "created_date"),
        xqb.DateDiff("end_date", "start_date", "days_between"),
        xqb.DateAdd("created_at", "1", "DAY", "next_day"),
        xqb.DateSub("created_at", "1", "MONTH", "prev_month"),
        xqb.DateFormat("created_at", "%Y-%m-%d", "formatted_date"),
    )
// Sql: SELECT DATE(created_at) AS created_date, DATEDIFF(end_date, start_date) AS days_between, DATE_ADD(created_at, INTERVAL 1 DAY) AS next_day, DATE_SUB(created_at, INTERVAL 1 MONTH) AS prev_month, DATE_FORMAT(created_at, '%Y-%m-%d') AS formatted_date FROM events
```

## JSON Functions

```go
// JSON operations
qb := xqb.Table("users").
    Select(
        xqb.JsonExtract("metadata", "preferences.theme", "theme"),
        xqb.JSONFunc("JSON_UNQUOTE", []string{"data", "'$.phone'"}, "phone"),
    )
// Sql: SELECT JSON_EXTRACT(metadata, '$.preferences.theme') AS theme, JSON_UNQUOTE(data, '$.phone') AS phone FROM users
```

## Math Expressions

```go
// Math operations
qb := xqb.Table("orders").
    Select(
        xqb.Math("amount * 1.1", "total_with_tax"),
        xqb.Coalesce([]string{"middle_name", "'N/A'"}, "coalesced_name"),
    )
// Sql: SELECT amount * 1.1 AS total_with_tax, COALESCE(middle_name, 'N/A') AS coalesced_name FROM orders
```

## INSERT Queries

```go
// Insert single record
affected, _ := xqb.Table("users").
    Insert([]map[string]any{
        {"name": "John Doe", "email": "john@example.com"},
    })

// Insert multiple records
affected, _ := xqb.Table("users").
    Insert([]map[string]any{
        {"name": "John Doe", "email": "john@example.com"},
        {"name": "Jane Doe", "email": "jane@example.com"},
    })

// Insert and get ID
lastId, _ := xqb.Table("users").
    InsertGetId([]map[string]any{
        {"name": "John Doe", "email": "john@example.com"},
    })
```

## UPDATE Queries

```go
// Update records
affected, _ := xqb.Table("users").
    Where("id", "=", 1).
    Update(map[string]any{
        "name": "Jane Doe",
        "email": "jane@example.com",
    })
// Sql: UPDATE users SET name = ?, email = ? WHERE id = ?

// Update with multiple conditions
affected, _ := xqb.Table("users").
    Where("active", "=", true).
    Where("age", ">", 18).
    Update(map[string]any{
        "status": "verified",
    })
```

## DELETE Queries

```go
// Delete records
affected, _ := xqb.Table("users").
    Where("id", "=", 1).
    Delete()
// Sql: DELETE FROM users WHERE id = ?

// Delete with multiple conditions
affected, _ := xqb.Table("users").
    Where("active", "=", false).
    Where("last_login", "<", "2023-01-01").
    Delete()
```

## Raw Sql

```go
// Execute raw Sql
result, _ := xqb.Sql("INSERT INTO users (name, email) VALUES (?, ?)", "John", "john@example.com").
    Connection("secondary_connection").
    Execute()

// Query raw Sql
rows, _ := xqb.Sql("SELECT * FROM users WHERE age > ?", 18).
    Query()

// Query single row
row, _ := xqb.Sql("SELECT COUNT(*) FROM users").
    QueryRow()
```

## Transactions

```go
// Simple transaction
err := xqb.Transaction(func(tx *sql.Tx) error {
    lastId, _ := xqb.Table("users").WithTx(tx).
        InsertGetId([]map[string]any{
            {"name": "John", "email": "john@example.com"},
        })

    affected, _ := xqb.Table("profiles").WithTx(tx).
        Where("user_id", "=", lastId).
        Update(map[string]any{
            "bio": "New user",
        })

    return nil
})

// Transaction on specific connection
// Note: You can use any connection in the DBManager
err := xqb.TransactionOn("connection_name", func(tx *sql.Tx) error {
  //...
})

// Manual transaction
tx, _ := xqb.BeginTx() || xqb.BeginTxOn("connection_name")

lastId, err := xqb.Table("users").WithTx(tx).
    InsertGetId([]map[string]any{
        {"name": "John", "email": "john@example.com"},
    })
if err != nil {
    tx.Rollback()
}

tx.Commit()
```

## Query Execution

```go
// Get all results
results, _ := qb.Get() // Returns []map[string]any

// Get first result
user, _ := qb.First() // Returns map[string]any

// aggregate execution
count, _ := qb.Count("id")
max, _ := qb.Max("id")
min, _ := qb.Min("id")
avg, _ := qb.Avg("id")
sum, _ := qb.Sum("id")

// Check if exists
exists, _ := qb.Exists()

// Get single value
value, _ := qb.Value("name")

// Pluck specific columns
names, _ := qb.Pluck("name", "id") // Returns map[string]any

// Chunk large results
err := qb.Chunk(100, func(rows []map[string]any) error {
    // Process 100 records at a time
    return nil
})

// Pagination
results, meta, _ := qb.Paginate(10, 1, true)
// meta contains: total_count, current_page, last_page, next_page, prev_page
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.
