# XQB Query Builder

XQB is a powerful and flexible SQL query builder for Go that provides a fluent interface for building complex SQL queries. It supports MySQL syntax and can be extended to support other database dialects.

## Features

- Fluent interface for building SQL queries
- Support for complex queries including:
  - SELECT, INSERT, UPDATE, DELETE operations
  - JOINs (INNER, LEFT, RIGHT, CROSS)
  - WHERE conditions with various operators
  - GROUP BY and HAVING clauses
  - ORDER BY with multiple columns
  - LIMIT and OFFSET
  - Common Table Expressions (CTEs)
  - Subqueries
  - Unions
  - JSON operations
  - Aggregate functions
  - String and Date functions
  - Mathematical expressions
  - Conditional expressions
  - Index hints
  - Locking clauses
- Parameter binding for safe SQL execution
- Extensible grammar system for different database dialects
- Type-safe query building

## Installation

```bash
go get github.com/iMohamedSheta/xqb
```

## Database Connection

XQB provides a simple database connection manager (`DBManager`) that acts as a bridge between your database connection and the query builder. This design choice follows the principle of separation of concerns:

- The `DBManager` is responsible only for managing the connection instance for the query builder
- Connection pool configuration and database-specific settings are handled at the database driver level
- This gives you full control over your database connection while keeping the query builder focused on its primary responsibility

Here's how to set up and use it:

```go
package main

import (
    "database/sql"
    "time"
    _ "github.com/go-sql-driver/mysql" // Import your database driver
    "github.com/iMohamedSheta/xqb"
)

func main() {
    // Create a database connection
    db, err := sql.Open("mysql", "user:password@tcp(localhost:3306)/database_name?parseTime=true")
    if err != nil {
        panic(err)
    }

    // Configure connection pool settings
    // These settings are managed by the database driver, not XQB
    db.SetMaxOpenConns(25)                // Maximum number of open connections
    db.SetMaxIdleConns(5)                 // Maximum number of idle connections
    db.SetConnMaxLifetime(5 * time.Minute) // Maximum lifetime of a connection
    db.SetConnMaxIdleTime(10 * time.Minute) // Maximum idle time of a connection

    // Set up the connection in XQB
    // DBManager only manages the connection instance for the query builder
    dbManager := xqb.GetDBManager()
    dbManager.SetDB(db)

    // Verify the connection
    if !dbManager.IsDBConnected() {
        panic("Database connection failed")
    }

    // Don't forget to close the connection when done
    defer dbManager.CloseDB()
}
```

### Using Transactions

```go
dbManager := xqb.GetDBManager()

// Simple transaction
err := dbManager.Transaction(func(tx *sql.Tx) error {
    // Create a query builder
    qb := xqb.Table("users")

    // Execute your queries within the transaction
    // If any query fails, the entire transaction will be rolled back
    return nil
})

// Manual transaction control
tx, err := dbManager.BeginTransaction()
if err != nil {
    panic(err)
}

// Use the transaction
// ...

// Commit or rollback
if err != nil {
    tx.Rollback()
} else {
    tx.Commit()
}
```

## Executing Queries

Once you have set up the database connection, you can execute various types of queries. Here are examples of common operations:

### Basic Query Execution

```go
dbManager := xqb.GetDBManager()

// SELECT query
qb := xqb.Table("users").
    Select("id", "name", "email").
    Where("active", "=", true)

// Execute the query and get results as map
results, err := qb.Get() // Returns []map[string]interface{}
if err != nil {
    panic(err)
}

// Process the results
for _, user := range results {
    // Access fields using map keys
    id := user["id"]
    name := user["name"]
    email := user["email"]
    // Use the data
}

// Get first row
user, err := qb.First() // Returns map[string]interface{}
if err != nil {
    if err == sql.ErrNoRows {
        // No user found
    } else {
        panic(err)
    }
}

// Get count
count, err := qb.Count("*", nil)
if err != nil {
    panic(err)
}
```

### Insert, Update, and Delete Operations

```go
// Insert a new record
affected, err := xqb.Table("users").
    Insert([]map[string]interface{}{
        {
            "name":  "John Doe",
            "email": "john@example.com",
        },
    })

// Get last inserted ID
lastID, err := xqb.Table("users").
    InsertGetId([]map[string]interface{}{
        {
            "name":  "John Doe",
            "email": "john@example.com",
        },
    })

// Update records
affected, err = xqb.Table("users").
    Update(map[string]interface{}{
        "name": "Jane Doe",
    }).
    Where("id", "=", 1)

// Delete records
affected, err = xqb.Table("users").
    Delete().
    Where("id", "=", 1)
```

### Pagination

XQB provides built-in support for pagination:

```go
// Using the Paginate method
results, meta, err := xqb.Table("users").
    Select("id", "name", "email").
    Where("active", "=", true).
    OrderBy("id", "DESC").
    Paginate(10, 1, true) // perPage, page, withCount

if err != nil {
    panic(err)
}

// Access paginated results
for _, user := range results {
    // Process each user
}

// Access pagination metadata
totalCount := meta["total_count"]
currentPage := meta["current_page"]
lastPage := meta["last_page"]
nextPage := meta["next_page"]
prevPage := meta["prev_page"]
```

### Advanced Query Execution

```go
// Chunking large result sets
err = xqb.Table("users").
    Where("active", "=", true).
    Chunk(100, func(rows []map[string]any) error {
        // Process 100 records at a time
        for _, row  := range rows {
            // Process each row
        }
        return nil
    })

// Get a single value from first record
value, err := xqb.Table("users").
    Where("id", "=", 1).
    Value("username")

// Check if record exists
exists, err := xqb.Table("users").
    Where("email", "=", "john@example.com").
    Exists()

// Pluck specific columns
names, err := xqb.Table("users").
    Select("name").
    Pluck("name", "id") // Returns map[string]any
```

### Error Handling

```go
// Proper error handling for queries
qb := xqb.Table("users").
    Select("id", "name").
    Where("email", "=", "john@example.com")

user, err := qb.First()
if err != nil {
    switch {
    case err == sql.ErrNoRows:
        // Handle no results found
        fmt.Println("User not found")
    case err != nil:
        // Handle other errors
        fmt.Printf("Error executing query: %v\n", err)
    }
    return
}

// Use the user data
fmt.Printf("Found user: %+v\n", user)
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/iMohamedSheta/xqb"
)

func main() {
    // Create a new query builder instance
    qb := xqb.Table("users")

    // Build a simple SELECT query
    sql, bindings, err := qb.
        Select("id", "name", "email").
        Where("active", "=", true).
        OrderBy("name", "ASC").
        Limit(10).
        ToSQL()

    if err != nil {
        panic(err)
    }

    fmt.Println("SQL:", sql)
    fmt.Println("Bindings:", bindings)
}
```

## Basic Usage

### SELECT Queries

```go
// Simple SELECT
qb := xqb.Table("users").
    Select("id", "name", "email").
    Where("active", "=", true)

// SELECT with JOIN
qb := xqb.Table("orders").
    Select("orders.id", "users.name").
    Join("users", "orders.user_id = users.id").
    Where("orders.status", "=", "pending")

// SELECT with GROUP BY and HAVING
qb := xqb.Table("orders").
    Select("user_id", "COUNT(*) as order_count").
    GroupBy("user_id").
    Having("order_count", ">", 5)

// SELECT with subquery
subquery := xqb.Table("orders").
    Select("user_id").
    Where("status", "=", "completed")

qb := xqb.Table("users").
    Select("id", "name").
    WhereIn("id", subquery)
```

### INSERT Queries

```go
qb := xqb.Table("users").
    Insert(map[string]interface{}{
        "name":  "John Doe",
        "email": "john@example.com",
    })
```

### UPDATE Queries

```go
qb := xqb.Table("users").
    Where("id", "=", 1).
    Update(map[string]interface{}{
        "name": "Jane Doe",
    })
```

### DELETE Queries

```go
qb := xqb.Table("users").
    Where("id", "=", 1).
    Delete()
```

## Advanced Features

### Common Table Expressions (CTEs)

```go
cte := xqb.Table("orders").
    Select("user_id", "SUM(amount) as total").
    GroupBy("user_id")

qb := xqb.Table("users").
    With("user_totals", cte).
    Select("users.*", "user_totals.total").
    Join("user_totals", "users.id = user_totals.user_id")
```

### JSON Operations

```go
qb := xqb.Table("users").
    Select("id", "JSON_EXTRACT(data, '$.name') as name").
    Where("JSON_EXTRACT(data, '$.age')", ">", 18)
```

### Complex WHERE Conditions

```go
qb := xqb.Table("users").
    Where("age", ">", 18).
    Where("status", "=", "active").
    WhereIn("id", []interface{}{1, 2, 3}).
    WhereBetween("created_at", "2023-01-01", "2023-12-31")
```

## Adding New Grammar

To add support for a new database dialect, you need to implement the `GrammarInterface`:

1. Create a new grammar struct that embeds `BaseGrammar`:

```go
type PostgreSQLGrammar struct {
    BaseGrammar
}
```

2. Implement the required interface methods:

```go
func (pg *PostgreSQLGrammar) CompileSelect(qb *types.QueryBuilderData) (string, []interface{}, error) {
    // Implement PostgreSQL-specific SELECT compilation
}

func (pg *PostgreSQLGrammar) CompileInsert(qb *types.QueryBuilderData) (string, []interface{}, error) {
    // Implement PostgreSQL-specific INSERT compilation
}

// ... implement other required methods
```

3. Register your grammar:

```go
func GetGrammar(driverName string) GrammarInterface {
    switch driverName {
    case "postgres":
        return &PostgreSQLGrammar{}
    case "mysql":
        return &MySQLGrammar{}
    default:
        return &MySQLGrammar{} // Default to MySQL
    }
}
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
