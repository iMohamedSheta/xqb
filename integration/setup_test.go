package integration

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/iMohamedSheta/xqb"
)

var dbManager *xqb.DBManager

// TestMain runs before all tests
func TestMain(m *testing.M) {
	TestingSetup()
	code := m.Run()
	TestingTeardown()
	os.Exit(code)
}

// TestingSetup runs before all tests
func TestingSetup() {
	dsn := "test:test@tcp(127.0.0.1:3306)/test_xqb_db"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		fmt.Println("❌ Failed to connect to MySql:", err)
		os.Exit(1)
	}

	if err = db.Ping(); err != nil {
		fmt.Println("❌ Failed to ping MySql:", err)
		os.Exit(1)
	}

	dbManager = xqb.GetDBManager()
	dbManager.SetDB(db)
}

// TestingTeardown runs after all tests
func TestingTeardown() {
	if dbManager != nil {
		db, err := dbManager.GetDB()
		if err != nil {
			fmt.Println("⚠️ Failed to get DB for teardown:", err)
		} else {
			dropAllTables(db)
		}

		if err := dbManager.CloseDB(); err != nil {
			fmt.Println("⚠️ Failed to close DB:", err)
		}
	}
}

func dropAllTables(db *sql.DB) {
	rows, err := db.Query("SHOW TABLES")
	if err != nil {
		fmt.Println("⚠️ Failed to list tables:", err)
		return
	}
	defer rows.Close()

	var table string
	for rows.Next() {
		if err := rows.Scan(&table); err == nil {
			if strings.HasPrefix(table, "test_") {
				_, err := db.Exec("DROP TABLE IF EXISTS `" + table + "`")
				if err != nil {
					fmt.Printf("⚠️ Failed to drop table %s: %v\n", table, err)
				} else {
					fmt.Printf("✅ Dropped table: %s\n", table)
				}
			}
		}
	}
}
