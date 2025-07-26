package xqb

import (
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/iMohamedSheta/xqb/dialects"
	"github.com/iMohamedSheta/xqb/shared/enums"
	xqbErr "github.com/iMohamedSheta/xqb/shared/errors"
	"github.com/iMohamedSheta/xqb/shared/types"
)

type QueryBuilderSettings struct {
	mu                    sync.RWMutex
	onBeforeQueryCallback func(qb *QueryBuilder)
	onAfterQueryCallback  func(query *QueryExecuted)
}

func NewQueryBuilderSettings() *QueryBuilderSettings {
	return &QueryBuilderSettings{}
}

var defaultSettings = NewQueryBuilderSettings()

func DefaultSettings() *QueryBuilderSettings {
	return defaultSettings
}

func (settings *QueryBuilderSettings) OnBeforeQuery(onBeforeQuery func(qb *QueryBuilder)) {
	settings.mu.Lock()
	defer settings.mu.Unlock()
	settings.onBeforeQueryCallback = onBeforeQuery
}

func (settings *QueryBuilderSettings) OnAfterQuery(onAfterQuery func(query *QueryExecuted)) {
	settings.mu.Lock()
	defer settings.mu.Unlock()
	settings.onAfterQueryCallback = onAfterQuery
}

func (settings *QueryBuilderSettings) GetOnBeforeQuery() func(qb *QueryBuilder) {
	settings.mu.RLock()
	defer settings.mu.RUnlock()
	return settings.onBeforeQueryCallback
}

func (settings *QueryBuilderSettings) GetOnAfterQuery() func(query *QueryExecuted) {
	settings.mu.RLock()
	defer settings.mu.RUnlock()
	return settings.onAfterQueryCallback
}

// QueryBuilder structure with all possible SELECT components
type QueryBuilder struct {
	connection      string
	settings        *QueryBuilderSettings
	dialect         dialects.DialectInterface
	queryType       enums.QueryType
	table           *types.Table
	columns         []any
	where           []*types.WhereCondition
	orderBy         []*types.OrderBy
	groupBy         []string
	having          []*types.Having
	limit           int
	offset          int
	joins           []*types.Join
	unions          []*types.Union
	bindings        []*types.Binding
	distinct        bool
	withCTEs        []*types.CTE
	isUsingDistinct bool
	tx              *sql.Tx
	errors          []error
	deleteFrom      []string
	options         map[types.Option]any // field for flexible Sql extensions
	insertedValues  []map[string]any
	updatedBindings []*types.Binding
}

func (qb *QueryBuilder) GetDialect() dialects.DialectInterface {
	return qb.dialect
}

func (qb *QueryBuilder) GetSettings() *QueryBuilderSettings {
	if qb.settings != nil {
		return qb.settings
	}
	return DefaultSettings()
}

func (qb *QueryBuilder) GetConnection() string {
	return qb.connection
}

func (qb *QueryBuilder) GetTable() *types.Table {
	return qb.table
}

func (qb *QueryBuilder) WithSettings(settings *QueryBuilderSettings) *QueryBuilder {
	qb.settings = settings
	return qb
}

type QueryExecuted struct {
	Sql        string
	Bindings   []any
	Time       time.Duration
	Connection string
	Dialect    types.Driver
	Err        error
}

// New creates a new QueryBuilder instance
func New() *QueryBuilder {
	// Get the driver name from the database connection
	driverName := types.DriverMySql // Default to MySql

	return &QueryBuilder{
		queryType:       enums.SELECT,
		columns:         []any{},
		where:           nil,
		orderBy:         nil,
		groupBy:         []string{},
		having:          nil,
		limit:           0,
		offset:          0,
		joins:           nil,
		unions:          nil,
		bindings:        nil,
		dialect:         dialects.GetDialect(driverName),
		distinct:        false,
		withCTEs:        nil,
		isUsingDistinct: false,
		tx:              nil,
		errors:          nil,
		deleteFrom:      nil,
		options:         make(map[types.Option]any),
		settings:        DefaultSettings(),
		connection:      DBManager().GetDefaultConnectionName(),
		table:           nil,
	}
}

func Query() *QueryBuilder {
	return New()
}

type wrappedDialect interface {
	Wrap(string) string
}

func (qb *QueryBuilder) Wrap(value string) string {
	if dialect, ok := qb.GetDialect().(wrappedDialect); ok {
		return dialect.Wrap(value)
	}
	return value
}

// Table creates a new QueryBuilder instance for a specific table
func Table(table string) *QueryBuilder {
	qb := New()
	qb.table = &types.Table{Name: table}
	return qb
}

func (qb *QueryBuilder) Table(table string) *QueryBuilder {
	qb.table = &types.Table{Name: table}
	return qb
}

// Reset resets the QueryBuilder instance
func (qb *QueryBuilder) Reset() {
	qb.queryType = enums.SELECT
	qb.connection = DBManager().GetDefaultConnectionName()
	qb.table = nil
	qb.columns = nil
	qb.where = nil
	qb.orderBy = nil
	qb.groupBy = nil
	qb.having = nil
	qb.limit = 0
	qb.offset = 0
	qb.joins = nil
	qb.unions = nil
	qb.bindings = nil
	qb.distinct = false
	qb.withCTEs = nil
	qb.isUsingDistinct = false
	qb.errors = nil
	qb.deleteFrom = nil
	qb.tx = nil
	qb.options = make(map[types.Option]any)
	qb.settings = DefaultSettings()
}

// GetData returns the QueryBuilderData for use by grammars
func (qb *QueryBuilder) GetData() *types.QueryBuilderData {
	return &types.QueryBuilderData{
		QueryType:       qb.queryType,
		Table:           qb.table,
		Columns:         qb.columns,
		Where:           qb.where,
		OrderBy:         qb.orderBy,
		GroupBy:         qb.groupBy,
		Having:          qb.having,
		Limit:           qb.limit,
		Offset:          qb.offset,
		Joins:           qb.joins,
		Unions:          qb.unions,
		Bindings:        qb.bindings,
		Distinct:        qb.distinct,
		WithCTEs:        qb.withCTEs,
		IsUsingDistinct: qb.isUsingDistinct,
		Errors:          qb.errors,
		DeleteFrom:      qb.deleteFrom,
		Options:         qb.options,
		InsertedValues:  qb.insertedValues,
		UpdatedBindings: qb.updatedBindings,
	}
}

func (qb *QueryBuilder) SetDialect(dialect types.Driver) *QueryBuilder {
	qb.dialect = dialects.GetDialect(dialect)
	return qb
}

// ToSql compiles the query to Sql
func (qb *QueryBuilder) ToSql() (string, []any, error) {
	if before := qb.GetSettings().GetOnBeforeQuery(); before != nil {
		before(qb)
	}
	start := time.UnixMicro(time.Now().UnixMicro())

	sql, bindings, err := qb.GetDialect().Build(qb.GetData())
	queryExecuted := &QueryExecuted{
		Sql:        sql,
		Bindings:   bindings,
		Time:       time.Duration(time.Now().UnixNano() - start.UnixNano()),
		Connection: qb.connection,
		Dialect:    qb.dialect.GetDriver(),
		Err:        err,
	}

	if after := qb.GetSettings().GetOnAfterQuery(); after != nil {
		after(queryExecuted)
	}

	return sql, bindings, err
}

// To Expression compiles the query to Sql and returns the Expression
func (qb *QueryBuilder) ToRawExpr() *types.Expression {
	sql, bindings, err := qb.ToSql()
	if err != nil {
		qb.errors = append(qb.errors, err)
		return nil
	}
	return &types.Expression{
		Sql:      sql,
		Bindings: bindings,
	}
}

func (qb *QueryBuilder) WithTx(tx *sql.Tx) *QueryBuilder {
	qb.tx = tx
	return qb
}

func (qb *QueryBuilder) Connection(connection string) *QueryBuilder {
	qb.connection = connection
	return qb
}

func (qb *QueryBuilder) SetOption(key types.Option, value any) {
	if qb.options == nil {
		qb.options = make(map[types.Option]any)
	}

	qb.options[key] = value
}

func (qb *QueryBuilder) GetOption(key types.Option) (any, bool) {
	val, ok := qb.options[key]
	return val, ok
}

func (qb *QueryBuilder) appendError(err error) {
	qb.errors = append(qb.errors, err)
}

// ToSqlQuery returns the query to Sql and replace bindings with actual values
func (qb *QueryBuilder) ToSqlView() (string, error) {
	sql, bindings, err := qb.ToSql()
	if err != nil {
		return "", err
	}
	finalSql, err := InjectBindings(qb.dialect.GetDriver(), sql, bindings)
	if err != nil {
		return "", err
	}

	return finalSql, nil
}

func InjectBindings(dialect types.Driver, sql string, bindings []any) (string, error) {
	var finalSql string
	switch dialect {
	case types.DriverMySql:
		// Replace `?` one by one with corresponding value
		for _, b := range bindings {
			sql = strings.Replace(sql, "?", formatBinding(b), 1)
		}
		finalSql = sql

	case types.DriverPostgres:
		// Replace `$1`, `$2`, ... with corresponding value
		for i, b := range bindings {
			placeholder := fmt.Sprintf("$%d", i+1)
			sql = strings.Replace(sql, placeholder, formatBinding(b), 1)
		}
		finalSql = sql

	default:
		return "", fmt.Errorf("%w: unsupported dialect please update ToSqlView method to support your dialect", xqbErr.ErrUnsupportedFeature)
	}

	return finalSql, nil
}

func formatBinding(value any) string {
	switch v := value.(type) {
	case string:
		return "'" + strings.ReplaceAll(v, "'", "''") + "'"
	case bool:
		return fmt.Sprintf("%t", v) // "true" / "false"
	case time.Time:
		return "'" + v.Format("2006-01-02 15:04:05") + "'"
	case nil:
		return "NULL"
	default:
		return fmt.Sprintf("%v", v)
	}
}
