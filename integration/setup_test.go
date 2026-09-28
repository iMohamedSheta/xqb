package integration

// Integration harness for xqb.
//
// Runs the SAME test suite against every supported database
// (db matrix): mysql, mariadb, postgres, sqlite, sqlserver.
//
// How it works:
//   - Each DB is configured via env (with CI-friendly defaults).
//   - TestMain tries to connect to every DB. Unreachable DBs are SKIPPED,
//     reachable ones run the full suite as subtests (t.Run(dialect)).
//   - Each test calls resetTables(t, conn) to get a clean schema + seed,
//     so tests are isolated and order-independent.
//
// Env overrides:
//   XQB_MYSQL_DSN      default test:test@tcp(127.0.0.1:3306)/test_xqb_db?parseTime=true&multiStatements=true&charset=utf8mb4
//   XQB_MARIADB_DSN    default test:test@tcp(127.0.0.1:3307)/test_xqb_db?parseTime=true&multiStatements=true&charset=utf8mb4
//   XQB_POSTGRES_DSN   default postgres://test:test@127.0.0.1:5432/test_xqb_db?sslmode=disable
//   XQB_SQLITE_DSN     default file:/tmp/xqb_integration_test.db?cache=shared&mode=rwc&_journal_mode=WAL
//   XQB_SQLSERVER_DSN  default sqlserver://sa:Strong!Pass123@127.0.0.1:1433?database=test_xqb_db
//   XQB_DIALECTS       comma list to limit matrix, e.g. "sqlite,mysql" (default all)

import (
	"database/sql"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/iMohamedSheta/xqb"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/microsoft/go-mssqldb"
)

// connName is the xqb connection name (= dialect key).
type dbHandle struct {
	connName string
	dialect  xqb.Dialect
	db       *sql.DB
}

var testHandles = map[string]*dbHandle{}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func wantedDialects() map[string]bool {
	raw := strings.ToLower(strings.TrimSpace(os.Getenv("XQB_DIALECTS")))
	if raw == "" || raw == "all" {
		return nil // all
	}
	m := map[string]bool{}
	for _, p := range strings.Split(raw, ",") {
		p = strings.TrimSpace(p)
		if p != "" {
			m[p] = true
		}
	}
	return m
}

func TestMain(m *testing.M) {
	configs := []struct {
		connName string
		dialect  xqb.Dialect
		driver   string
		dsn      string
	}{
		{
			connName: "mysql",
			dialect:  xqb.DialectMySql,
			driver:   "mysql",
			dsn:      envOr("XQB_MYSQL_DSN", "test:test@tcp(127.0.0.1:3306)/test_xqb_db?parseTime=true&multiStatements=true&charset=utf8mb4"),
		},
		{
			connName: "mariadb",
			dialect:  xqb.DialectMariaDB,
			driver:   "mysql",
			dsn:      envOr("XQB_MARIADB_DSN", "test:test@tcp(127.0.0.1:3307)/test_xqb_db?parseTime=true&multiStatements=true&charset=utf8mb4"),
		},
		{
			connName: "postgres",
			dialect:  xqb.DialectPostgres,
			driver:   "postgres",
			dsn:      envOr("XQB_POSTGRES_DSN", "postgres://test:test@127.0.0.1:5432/test_xqb_db?sslmode=disable"),
		},
		{
			connName: "sqlite",
			dialect:  xqb.DialectSQLite,
			driver:   "sqlite3",
			dsn:      envOr("XQB_SQLITE_DSN", ""),
		},
		{
			connName: "sqlserver",
			dialect:  xqb.DialectSQLServer,
			driver:   "sqlserver",
			dsn:      envOr("XQB_SQLSERVER_DSN", "sqlserver://sa:Strong!Pass123@127.0.0.1:1433?database=test_xqb_db"),
		},
	}

	// sqlite always works locally: fall back to a temp file when no DSN given.
	sqlitePath := ""
	for i, c := range configs {
		if c.connName == "sqlite" && c.dsn == "" {
			sqlitePath = fmt.Sprintf("%s/xqb_integration_test.db", os.TempDir())
			_ = os.Remove(sqlitePath)
			configs[i].dsn = fmt.Sprintf("file:%s?cache=shared&mode=rwc&_journal_mode=WAL", sqlitePath)
		}
	}

	filter := wantedDialects()
	// In CI set XQB_REQUIRE_ALL=1 so a missing database FAILS instead of
	// being skipped (db matrix: every dialect must pass).
	requireAll := strings.TrimSpace(os.Getenv("XQB_REQUIRE_ALL")) == "1"
	var missing []string

	for _, c := range configs {
		if filter != nil && !filter[c.connName] {
			fmt.Printf("⏭️  skip %s (filtered by XQB_DIALECTS)\n", c.connName)
			continue
		}
		db, err := sql.Open(c.driver, c.dsn)
		if err != nil {
			fmt.Printf("⏭️  skip %s: open failed: %v\n", c.connName, err)
			missing = append(missing, c.connName)
			continue
		}
		// Keep sqlite to a single pooled connection so :memory:/file state is stable.
		if c.connName == "sqlite" {
			db.SetMaxOpenConns(1)
		}
		// Retry ping for slow CI services (up to ~60s).
		var pingErr error
		for attempt := 0; attempt < 60; attempt++ {
			pingErr = db.Ping()
			if pingErr == nil {
				break
			}
			time.Sleep(time.Second)
		}
		if pingErr != nil {
			fmt.Printf("⏭️  skip %s: ping failed: %v\n", c.connName, pingErr)
			_ = db.Close()
			missing = append(missing, c.connName)
			continue
		}
		if err := xqb.AddConnection(&xqb.Connection{Name: c.connName, Dialect: c.dialect, DB: db}); err != nil {
			fmt.Printf("⏭️  skip %s: add connection failed: %v\n", c.connName, err)
			_ = db.Close()
			missing = append(missing, c.connName)
			continue
		}
		testHandles[c.connName] = &dbHandle{connName: c.connName, dialect: c.dialect, db: db}
		fmt.Printf("✅ integration DB ready: %s (%s)\n", c.connName, c.dialect)
	}

	if requireAll && len(missing) > 0 {
		fmt.Printf("❌ XQB_REQUIRE_ALL=1 but these dialects are unreachable: %s\n", strings.Join(missing, ", "))
		os.Exit(1)
	}

	if len(testHandles) == 0 {
		fmt.Println("❌ No integration databases reachable. Start at least one (sqlite works with no setup).")
		os.Exit(1)
	}

	// Default connection = first available so Table() without .Connection() still works.
	names := sortedConnNames()
	_ = xqb.DBManager().SetDefaultConnection(names[0])

	code := m.Run()

	// Teardown: drop test tables, close, remove sqlite file.
	for _, name := range names {
		h := testHandles[name]
		_ = dropTables(h.db, h.dialect)
		_ = h.db.Close()
	}
	if sqlitePath != "" {
		_ = os.Remove(sqlitePath)
	}

	os.Exit(code)
}

func sortedConnNames() []string {
	names := make([]string, 0, len(testHandles))
	for n := range testHandles {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}

// forEachDB runs fn once per reachable database as a named subtest.
func forEachDB(t *testing.T, fn func(t *testing.T, conn string)) {
	t.Helper()
	for _, name := range sortedConnNames() {
		name := name
		t.Run(name, func(t *testing.T) {
			fn(t, name)
		})
	}
}

// QB is a shortcut for a builder bound to a named connection.
func QB(conn, table string) *xqb.QueryBuilder {
	return xqb.Table(table).Connection(conn)
}

// ---------------------------------------------------------------------------
// Schema
// ---------------------------------------------------------------------------

func dropTables(db *sql.DB, dialect xqb.Dialect) error {
	for _, tbl := range []string{"xqb_posts", "xqb_users"} {
		q := fmt.Sprintf("DROP TABLE IF EXISTS %s", tbl)
		// SQL Server 2016+ supports DROP TABLE IF EXISTS as well.
		if _, err := db.Exec(q); err != nil {
			// Best effort: old SQL Server fallback.
			_, _ = db.Exec(fmt.Sprintf("IF OBJECT_ID('%s', 'U') IS NOT NULL DROP TABLE %s", tbl, tbl))
		}
	}
	return nil
}

func createTables(t *testing.T, conn string) {
	t.Helper()
	h := testHandles[conn]
	db := h.db
	if err := dropTables(db, h.dialect); err != nil {
		t.Fatalf("[%s] drop tables: %v", conn, err)
	}

	var usersDDL, postsDDL string
	switch h.dialect {
	case xqb.DialectPostgres:
		usersDDL = `CREATE TABLE xqb_users (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			email VARCHAR(255) UNIQUE,
			age INT NULL,
			score DOUBLE PRECISION NULL,
			status VARCHAR(50) NULL
		)`
		postsDDL = `CREATE TABLE xqb_posts (
			id SERIAL PRIMARY KEY,
			user_id INT NOT NULL,
			title VARCHAR(255) NOT NULL
		)`
	case xqb.DialectSQLite:
		usersDDL = `CREATE TABLE xqb_users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			email TEXT UNIQUE,
			age INTEGER NULL,
			score REAL NULL,
			status TEXT NULL
		)`
		postsDDL = `CREATE TABLE xqb_posts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			title TEXT NOT NULL
		)`
	case xqb.DialectSQLServer:
		usersDDL = `CREATE TABLE xqb_users (
			id INT IDENTITY(1,1) PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			email VARCHAR(255) UNIQUE,
			age INT NULL,
			score FLOAT NULL,
			status VARCHAR(50) NULL
		)`
		postsDDL = `CREATE TABLE xqb_posts (
			id INT IDENTITY(1,1) PRIMARY KEY,
			user_id INT NOT NULL,
			title VARCHAR(255) NOT NULL
		)`
	default: // mysql + mariadb
		usersDDL = "CREATE TABLE xqb_users (" +
			"id INT AUTO_INCREMENT PRIMARY KEY," +
			"name VARCHAR(255) NOT NULL," +
			"email VARCHAR(255) UNIQUE," +
			"age INT NULL," +
			"score DOUBLE NULL," +
			"status VARCHAR(50) NULL" +
			")"
		postsDDL = "CREATE TABLE xqb_posts (" +
			"id INT AUTO_INCREMENT PRIMARY KEY," +
			"user_id INT NOT NULL," +
			"title VARCHAR(255) NOT NULL" +
			")"
	}

	if _, err := db.Exec(usersDDL); err != nil {
		t.Fatalf("[%s] create xqb_users: %v", conn, err)
	}
	if _, err := db.Exec(postsDDL); err != nil {
		t.Fatalf("[%s] create xqb_posts: %v", conn, err)
	}
}

// resetTables drops, recreates and seeds both tables. Call at the start of
// every test so tests stay isolated and order-independent.
func resetTables(t *testing.T, conn string) {
	t.Helper()
	createTables(t, conn)
	seedUsers(t, conn)
	seedPosts(t, conn)
}

// resetUsersOnly is for insert-focused tests that need an empty table.
func resetUsersEmpty(t *testing.T, conn string) {
	t.Helper()
	createTables(t, conn)
}

// seedUsers inserts 4 deterministic rows via raw SQL (no placeholders, so it
// works unchanged on every dialect).
func seedUsers(t *testing.T, conn string) {
	t.Helper()
	db := testHandles[conn].db
	rows := []string{
		`('Alice', 'alice@example.com', 30, 95.5, 'active')`,
		`('Bob', 'bob@example.com', 25, 80.0, 'active')`,
		`('Carol', 'carol@example.com', 35, 70.25, 'inactive')`,
		`('Dave', NULL, NULL, NULL, NULL)`,
	}
	for _, r := range rows {
		q := `INSERT INTO xqb_users (name, email, age, score, status) VALUES ` + r
		if _, err := db.Exec(q); err != nil {
			t.Fatalf("[%s] seed users %s: %v", conn, r, err)
		}
	}
}

func seedPosts(t *testing.T, conn string) {
	t.Helper()
	db := testHandles[conn].db
	rows := []string{
		`(1, 'Alice first post')`,
		`(1, 'Alice second post')`,
		`(2, 'Bob post')`,
	}
	for _, r := range rows {
		q := `INSERT INTO xqb_posts (user_id, title) VALUES ` + r
		if _, err := db.Exec(q); err != nil {
			t.Fatalf("[%s] seed posts %s: %v", conn, r, err)
		}
	}
}

func countUsers(t *testing.T, conn string) int64 {
	t.Helper()
	n, err := QB(conn, "xqb_users").Count("*")
	if err != nil {
		t.Fatalf("[%s] count users: %v", conn, err)
	}
	return n
}

// ---------------------------------------------------------------------------
// Value normalizers (drivers return different Go types per dialect)
// ---------------------------------------------------------------------------

func asString(v any) string {
	if v == nil {
		return ""
	}
	switch x := v.(type) {
	case string:
		return x
	case []byte:
		return string(x)
	case fmt.Stringer:
		return x.String()
	default:
		return fmt.Sprintf("%v", x)
	}
}

func asInt64(v any) int64 {
	if v == nil {
		return 0
	}
	switch x := v.(type) {
	case int64:
		return x
	case int32:
		return int64(x)
	case int:
		return int64(x)
	case int16:
		return int64(x)
	case int8:
		return int64(x)
	case uint64:
		return int64(x)
	case uint32:
		return int64(x)
	case uint:
		return int64(x)
	case float64:
		return int64(x)
	case float32:
		return int64(x)
	case []byte:
		n, _ := strconv.ParseInt(strings.TrimSpace(string(x)), 10, 64)
		return n
	case string:
		n, _ := strconv.ParseInt(strings.TrimSpace(x), 10, 64)
		return n
	case bool:
		if x {
			return 1
		}
		return 0
	default:
		n, _ := strconv.ParseInt(fmt.Sprintf("%v", x), 10, 64)
		return n
	}
}

func asFloat64(v any) float64 {
	if v == nil {
		return 0
	}
	switch x := v.(type) {
	case float64:
		return x
	case float32:
		return float64(x)
	case int64:
		return float64(x)
	case int:
		return float64(x)
	case int32:
		return float64(x)
	case []byte:
		f, _ := strconv.ParseFloat(strings.TrimSpace(string(x)), 64)
		return f
	case string:
		f, _ := strconv.ParseFloat(strings.TrimSpace(x), 64)
		return f
	default:
		f, _ := strconv.ParseFloat(fmt.Sprintf("%v", x), 64)
		return f
	}
}
