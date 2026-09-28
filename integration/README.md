# xqb integration tests (db matrix)

Same suite runs against **every** supported database:
`mysql`, `mariadb`, `postgres`, `sqlite`, `sqlserver`.

## Layout

- `setup_test.go` — matrix harness (env DSNs, per-dialect DDL, seed, value normalizers)
- `insert_test.go` — `Insert`, `InsertGetId`, `Upsert`
- `select_test.go` — `Get`, `First`, `Find`, `FindOrFail`, `Value`, `Pluck*`, `Exists`, `Distinct`, `Order/Limit`, `Paginate`, `Chunks`
- `where_test.go` — `Where*` family (`In`, `Null`, `Between`, `Or`, groups, `Raw`, `Exists`, subqueries)
- `update_test.go` / `delete_test.go` — write paths + affected rows + error paths
- `aggregates_test.go` — `Count`, `Sum`, `Avg`, `Min`, `Max`
- `joins_test.go` — `Join`, `LeftJoin`, `Union`, `UnionAll`
- `advanced_test.go` — `GroupBy/Having`, CTE (`With`), `FromSubquery`, raw, `Bind` models
- `transactions_test.go` — `TransactionOn` commit / rollback

## Run locally

```bash
# 1) start databases (sqlite needs nothing)
docker compose -f integration/docker-compose.yml up -d

# 2) create the sqlserver test db once (others auto-create)
docker exec -i $(docker ps -qf "ancestor=mcr.microsoft.com/mssql/server:2022-latest") \
  /opt/mssql-tools18/bin/sqlcmd -S localhost -U sa -P 'Strong!Pass123' -C \
  -Q "IF DB_ID('test_xqb_db') IS NULL CREATE DATABASE test_xqb_db"

# 3) run the matrix (from repo root)
cd integration && go mod tidy && go test -v ./...

# 4) run a subset, e.g. sqlite only
XQB_DIALECTS=sqlite go test -v ./...
XQB_DIALECTS=mysql,mariadb go test -v ./...
```

## Env overrides

| Var | Default |
|---|---|
| `XQB_MYSQL_DSN` | `test:test@tcp(127.0.0.1:3306)/test_xqb_db?parseTime=true&multiStatements=true&charset=utf8mb4` |
| `XQB_MARIADB_DSN` | `test:test@tcp(127.0.0.1:3307)/test_xqb_db?parseTime=true&multiStatements=true&charset=utf8mb4` |
| `XQB_POSTGRES_DSN` | `postgres://test:test@127.0.0.1:5432/test_xqb_db?sslmode=disable` |
| `XQB_SQLITE_DSN` | temp file `$(Temp)/xqb_integration_test.db` (auto) |
| `XQB_SQLSERVER_DSN` | `sqlserver://sa:Strong!Pass123@127.0.0.1:1433?database=test_xqb_db` |
| `XQB_DIALECTS` | `all` (or comma list, e.g. `sqlite,mysql`) |

Unreachable databases are **skipped**, not failed — so `go test` passes with
only sqlite available, and CI runs the full matrix.
