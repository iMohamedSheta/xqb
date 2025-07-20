package xqb

import (
	"database/sql"

	"github.com/iMohamedSheta/xqb/grammar"
	"github.com/iMohamedSheta/xqb/shared/enums"
	"github.com/iMohamedSheta/xqb/shared/types"
)

// QueryBuilder structure with all possible SELECT components
type QueryBuilder struct {
	connection        string
	grammar           grammar.GrammarInterface
	queryType         enums.QueryType
	table             *types.Table
	columns           []any
	where             []*types.WhereCondition
	orderBy           []*types.OrderBy
	groupBy           []string
	having            []*types.Having
	limit             int
	offset            int
	joins             []*types.Join
	unions            []*types.Union
	bindings          []*types.Binding
	distinct          bool
	withCTEs          []*types.CTE
	isUsingDistinct   bool
	isLockedForUpdate bool
	isInSharedLock    bool
	tx                *sql.Tx
	errors            []error
}

// New creates a new QueryBuilder instance
func New() *QueryBuilder {
	// Get the driver name from the database connection
	driverName := grammar.DriverMySQL // Default to MySQL

	return &QueryBuilder{
		queryType:         enums.SELECT,
		columns:           []any{},
		where:             nil,
		orderBy:           nil,
		groupBy:           []string{},
		having:            nil,
		limit:             0,
		offset:            0,
		joins:             nil,
		unions:            nil,
		bindings:          nil,
		grammar:           grammar.GetGrammar(driverName),
		distinct:          false,
		withCTEs:          nil,
		isUsingDistinct:   false,
		isLockedForUpdate: false,
		isInSharedLock:    false,
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
	qb.isLockedForUpdate = false
	qb.isInSharedLock = false
	qb.errors = nil
}

// GetData returns the QueryBuilderData for use by grammars
func (qb *QueryBuilder) GetData() *types.QueryBuilderData {
	return &types.QueryBuilderData{
		QueryType:         qb.queryType,
		Table:             qb.table,
		Columns:           qb.columns,
		Where:             qb.where,
		OrderBy:           qb.orderBy,
		GroupBy:           qb.groupBy,
		Having:            qb.having,
		Limit:             qb.limit,
		Offset:            qb.offset,
		Joins:             qb.joins,
		Unions:            qb.unions,
		Bindings:          qb.bindings,
		Distinct:          qb.distinct,
		WithCTEs:          qb.withCTEs,
		IsUsingDistinct:   qb.isUsingDistinct,
		IsLockedForUpdate: qb.isLockedForUpdate,
		IsInSharedLock:    qb.isInSharedLock,
		Errors:            qb.errors,
	}
}

func (qb *QueryBuilder) SetDialect(dialect grammar.Driver) {
	qb.grammar = grammar.GetGrammar(dialect)
}

// ToSQL compiles the query to SQL
func (qb *QueryBuilder) ToSQL() (string, []any, error) {
	return qb.grammar.CompileSelect(qb.GetData())
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
