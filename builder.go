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
	table             string
	columns           []any
	columnAliases     map[string]string
	where             []types.WhereCondition
	orderBy           []types.OrderBy
	groupBy           []string
	having            []types.Having
	limit             int
	offset            int
	joins             []types.Join
	unions            []types.Union
	bindings          []types.Binding
	distinct          bool
	subqueries        map[string]*QueryBuilder
	withCTEs          []types.CTE
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
		columnAliases:     make(map[string]string),
		where:             []types.WhereCondition{},
		orderBy:           []types.OrderBy{},
		groupBy:           []string{},
		having:            []types.Having{},
		limit:             0,
		offset:            0,
		joins:             []types.Join{},
		unions:            []types.Union{},
		bindings:          []types.Binding{},
		grammar:           grammar.GetGrammar(driverName),
		distinct:          false,
		subqueries:        make(map[string]*QueryBuilder),
		withCTEs:          []types.CTE{},
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
	qb.table = table
	return qb
}

func (qb *QueryBuilder) Table(table string) *QueryBuilder {
	qb.table = table
	return qb
}

// Reset resets the QueryBuilder instance
func (qb *QueryBuilder) Reset() {
	qb.queryType = enums.SELECT
	qb.connection = DBManager().GetDefaultConnectionName()
	qb.table = ""
	qb.columns = nil
	qb.columnAliases = nil
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
	qb.subqueries = nil
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
		ColumnAliases:     qb.columnAliases,
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
		Subqueries:        make(map[string]any),
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

func (qb *QueryBuilder) WithTx(tx *sql.Tx) *QueryBuilder {
	qb.tx = tx
	return qb
}

func (qb *QueryBuilder) Connection(connection string) *QueryBuilder {
	qb.connection = connection
	return qb
}
