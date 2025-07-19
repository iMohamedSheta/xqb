# XQB Query Builder

A powerful and flexible SQL query builder for Go with fluent interface for building complex SQL queries.

## Installation

```bash
go get github.com/iMohamedSheta/xqb
```

## Quick Start

```go
package main

import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "github.com/iMohamedSheta/xqb"
)

func main() {
    // Setup database connection
    db, _ := sql.Open("mysql", "user:password@tcp(localhost:3306)/database")
    xqb.AddConnection("default", db)
    
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

## Database Connection

```go
// Add connection
db, _ := sql.Open("mysql", "dsn")
xqb.AddConnection("default", db)

// Use specific connection
qb := xqb.Table("users").WithConnection("default")

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
// SQL: SELECT id, name, email FROM users

// Select with conditions
qb := xqb.Table("users").
    Select("id", "name").
    Where("age", ">", 18)
// SQL: SELECT id, name FROM users WHERE age > ?

// Select distinct
qb := xqb.Table("users").
    Select("id", "name").
    Distinct()
// SQL: SELECT DISTINCT id, name FROM users
```

### Joins
```go
// Inner join
qb := xqb.Table("users").
    Select("users.id", "users.name", "orders.id as order_id").
    Join("orders", "users.id = orders.user_id")
// SQL: SELECT users.id, users.name, orders.id as order_id FROM users JOIN orders ON users.id = orders.user_id

// Left join
qb := xqb.Table("users").
    LeftJoin("comments", "users.id = comments.user_id")
// SQL: SELECT * FROM users LEFT JOIN comments ON users.id = comments.user_id

// Right join
qb := xqb.Table("users").
    RightJoin("logins", "users.id = logins.user_id")
// SQL: SELECT * FROM users RIGHT JOIN logins ON users.id = logins.user_id

// Full join
qb := xqb.Table("users").
    FullJoin("sessions", "users.id = sessions.user_id")
// SQL: SELECT * FROM users FULL JOIN sessions ON users.id = sessions.user_id

// Cross join
qb := xqb.Table("users").
    CrossJoin("roles")
// SQL: SELECT * FROM users CROSS JOIN roles

// Join with conditions
qb := xqb.Table("users").
    Join("posts", "users.id = posts.user_id AND posts.status = ?", "active")
// SQL: SELECT * FROM users JOIN posts ON users.id = posts.user_id AND posts.status = ?
```

### Subquery Joins
```go
// Join subquery
sub := xqb.Table("posts").Where("published", "=", true)
qb := xqb.Table("users").
    JoinSub(sub, "p", "users.id = p.user_id")
// SQL: JOIN (SELECT * FROM posts WHERE published = ?) AS p ON users.id = p.user_id

// Left join subquery
sub := xqb.Table("comments").Where("active", "=", true)
qb := xqb.Table("users").
    LeftJoinSub(sub, "c", "users.id = c.user_id")
// SQL: LEFT JOIN (SELECT * FROM comments WHERE active = ?) AS c ON users.id = c.user_id

// Cross join subquery
sub := xqb.Table("plans").Where("expired", "=", false)
qb := xqb.Table("users").
    CrossJoinSub(sub, "p")
// SQL: CROSS JOIN (SELECT * FROM plans WHERE expired = ?) AS p
```

### Where Conditions
```go
// Basic where
qb := xqb.Table("users").
    Where("age", ">", 18)
// SQL: SELECT * FROM users WHERE age > ?

// Multiple conditions
qb := xqb.Table("users").
    Where("age", ">", 18).
    Where("active", "=", true)
// SQL: SELECT * FROM users WHERE age > ? AND active = ?

// OR conditions
qb := xqb.Table("users").
    Where("id", "=", 1).
    OrWhere("email", "=", "admin@example.com")
// SQL: SELECT * FROM users WHERE id = ? OR email = ?

// IN conditions
qb := xqb.Table("users").
    WhereIn("id", []any{1, 2, 3})
// SQL: SELECT * FROM users WHERE id IN (?, ?, ?)

// BETWEEN conditions
qb := xqb.Table("users").
    WhereBetween("age", 18, 65)
// SQL: SELECT * FROM users WHERE age BETWEEN ? AND ?

// NULL conditions
qb := xqb.Table("users").
    WhereNull("deleted_at")
// SQL: SELECT * FROM users WHERE deleted_at IS NULL

// EXISTS conditions
subQuery := xqb.Table("orders").Select("user_id").Where("status", "=", "active")
qb := xqb.Table("users").
    WhereExists(subQuery)
// SQL: SELECT * FROM users WHERE EXISTS (SELECT user_id FROM orders WHERE status = ?)

// Raw where
qb := xqb.Table("users").
    WhereRaw("CASE WHEN status = 'active' THEN 1 ELSE 0 END = ?", 1)
// SQL: SELECT * FROM users WHERE CASE WHEN status = 'active' THEN 1 ELSE 0 END = ?

// Where groups
qb := xqb.Table("users").
    Where("id", "=", 1).
    WhereGroup(func(qb *xqb.QueryBuilder) {
        qb.WhereNull("deleted_at").OrWhereNull("disabled_at")
    })
// SQL: SELECT * FROM users WHERE id = ? AND (deleted_at IS NULL OR disabled_at IS NULL)
```

### Group By and Having
```go
// Group by
qb := xqb.Table("orders").
    Select("user_id", "COUNT(*) as order_count").
    GroupBy("user_id")
// SQL: SELECT user_id, COUNT(*) as order_count FROM orders GROUP BY user_id

// Having
qb := xqb.Table("orders").
    Select("user_id", "COUNT(*) as order_count").
    GroupBy("user_id").
    Having("order_count", ">", 5)
// SQL: SELECT user_id, COUNT(*) as order_count FROM orders GROUP BY user_id HAVING order_count > ?
```

### Order By
```go
// Order by
qb := xqb.Table("users").
    Select("id", "name").
    OrderBy("name", "ASC")
// SQL: SELECT id, name FROM users ORDER BY name ASC

// Multiple order by
qb := xqb.Table("users").
    OrderBy("age", "DESC").
    OrderBy("name", "ASC")
// SQL: SELECT * FROM users ORDER BY age DESC, name ASC
```

### Limit and Offset
```go
// Limit
qb := xqb.Table("users").
    Select("id", "name").
    Limit(10)
// SQL: SELECT id, name FROM users LIMIT 10

// Limit with offset
qb := xqb.Table("users").
    Select("id", "name").
    Limit(10).
    Offset(20)
// SQL: SELECT id, name FROM users LIMIT 10 OFFSET 20
```

### Common Table Expressions (CTE)
```go
// Simple CTE
qb := xqb.Table("users").
    WithRaw("user_totals", "SELECT user_id, SUM(amount) as total_spent FROM orders GROUP BY user_id").
    Select("users.id", "users.name", "user_totals.total_spent").
    Join("user_totals", "users.id = user_totals.user_id")
// SQL: WITH user_totals AS (SELECT user_id, SUM(amount) as total_spent FROM orders GROUP BY user_id) SELECT users.id, users.name, user_totals.total_spent FROM users JOIN user_totals ON users.id = user_totals.user_id

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
// SQL: SELECT id, name FROM users FOR UPDATE

// Shared lock
qb := xqb.Table("users").
    Select("id", "name").
    SharedLock()
// SQL: SELECT id, name FROM users LOCK IN SHARE MODE
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
// SQL: SELECT COUNT(id) AS order_count, SUM(amount) AS total_amount, AVG(amount) AS average_amount, MIN(amount) AS min_amount, MAX(amount) AS max_amount FROM orders
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
// SQL: SELECT CONCAT(first_name, ' ', last_name) AS full_name, LOWER(email) AS lower_email, UPPER(username) AS upper_username, LENGTH(bio) AS bio_length, TRIM(nickname) AS trimmed_nickname, REPLACE(title, 'foo', 'bar') AS replaced_title, SUBSTRING(description, 1, 10) AS short_desc FROM users
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
// SQL: SELECT DATE(created_at) AS created_date, DATEDIFF(end_date, start_date) AS days_between, DATE_ADD(created_at, INTERVAL 1 DAY) AS next_day, DATE_SUB(created_at, INTERVAL 1 MONTH) AS prev_month, DATE_FORMAT(created_at, '%Y-%m-%d') AS formatted_date FROM events
```

## JSON Functions

```go
// JSON operations
qb := xqb.Table("users").
    Select(
        xqb.JsonExtract("metadata", "preferences.theme", "theme"),
        xqb.JSONFunc("JSON_UNQUOTE", []string{"data", "'$.phone'"}, "phone"),
    )
// SQL: SELECT JSON_EXTRACT(metadata, '$.preferences.theme') AS theme, JSON_UNQUOTE(data, '$.phone') AS phone FROM users
```

## Math Expressions

```go
// Math operations
qb := xqb.Table("orders").
    Select(
        xqb.Math("amount * 1.1", "total_with_tax"),
        xqb.Coalesce([]string{"middle_name", "'N/A'"}, "coalesced_name"),
    )
// SQL: SELECT amount * 1.1 AS total_with_tax, COALESCE(middle_name, 'N/A') AS coalesced_name FROM orders
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
// SQL: UPDATE users SET name = ?, email = ? WHERE id = ?

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
// SQL: DELETE FROM users WHERE id = ?

// Delete with multiple conditions
affected, _ := xqb.Table("users").
    Where("active", "=", false).
    Where("last_login", "<", "2023-01-01").
    Delete()
```

## Raw SQL

```go
// Execute raw SQL
result, _ := xqb.Sql("INSERT INTO users (name, email) VALUES (?, ?)", "John", "john@example.com").
    Execute()

// Query raw SQL
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

// Manual transaction
tx, _ := xqb.BeginTransaction()
defer tx.Rollback()

lastId, _ := xqb.Table("users").WithTx(tx).
    InsertGetId([]map[string]any{
        {"name": "John", "email": "john@example.com"},
    })

tx.Commit()
```

## Query Execution

```go
// Get all results
results, _ := qb.Get() // Returns []map[string]any

// Get first result
user, _ := qb.First() // Returns map[string]any

// Get count
count, _ := qb.Count("id")

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
