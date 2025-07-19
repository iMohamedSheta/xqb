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
	qb.connection = Manager().defaultConnection
	qb.table = ""
	qb.columns = []any{}
	qb.columnAliases = make(map[string]string)
	qb.where = []types.WhereCondition{}
	qb.orderBy = []types.OrderBy{}
	qb.groupBy = []string{}
	qb.having = []types.Having{}
	qb.limit = 0
	qb.offset = 0
	qb.joins = []types.Join{}
	qb.unions = []types.Union{}
	qb.bindings = []types.Binding{}
	qb.distinct = false
	qb.subqueries = make(map[string]*QueryBuilder)
	qb.withCTEs = []types.CTE{}
	qb.isUsingDistinct = false
	qb.isLockedForUpdate = false
	qb.isInSharedLock = false
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
