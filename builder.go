package xqb

import (
	"database/sql"

	"github.com/iMohamedSheta/xqb/dialects"
	"github.com/iMohamedSheta/xqb/shared/enums"
	"github.com/iMohamedSheta/xqb/shared/types"
)

// QueryBuilder structure with all possible SELECT components
type QueryBuilder struct {
	connection      string
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
	options         map[types.Option]any // field for flexible SQL extensions
}

// New creates a new QueryBuilder instance
func New() *QueryBuilder {
	// Get the driver name from the database connection
	driverName := types.DriverMySQL // Default to MySQL

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
	}
}

func Query() *QueryBuilder {
	return New()
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
	}
}

func (qb *QueryBuilder) SetDialect(dialect types.Driver) *QueryBuilder {
	qb.dialect = dialects.GetDialect(dialect)
	return qb
}

// ToSQL compiles the query to SQL
func (qb *QueryBuilder) ToSQL() (string, []any, error) {
	return qb.dialect.Build(qb.GetData())
}

// To Expression compiles the query to SQL and returns the Expression
func (qb *QueryBuilder) ToRawExpr() *types.Expression {
	sql, bindings, err := qb.ToSQL()
	if err != nil {
		qb.errors = append(qb.errors, err)
		return nil
	}
	return &types.Expression{
		SQL:      sql,
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
