package xqb

import (
	"github.com/iMohamedSheta/xqb/grammar"
	"github.com/iMohamedSheta/xqb/types"
)

var QueryBuilderInstance *QueryBuilder

// QueryBuilder structure with all possible SELECT components
type QueryBuilder struct {
	grammar           grammar.GrammarInterface
	queryType         types.QueryType
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
	aggregateFuncs    []types.AggregateExpr
	subqueries        map[string]*QueryBuilder
	withCTEs          []types.CTE
	jsonExpressions   []types.JSONExpression
	mathExpressions   []types.MathExpression
	conditionalExprs  []types.ConditionalExpr
	stringFuncs       []types.StringFunction
	dateFuncs         []types.DateFunction
	indexHints        []string
	lockType          string
	forceIndex        string
	useIndex          string
	ignoreIndex       string
	procedure         string
	procedureParams   []any
	isUsingDistinct   bool
	isForUpdate       bool
	isLockInShareMode bool
	isHighPriority    bool
	isStraightJoin    bool
	isCalcFoundRows   bool
	comment           string
	ctes              []types.CTE
}

// New creates a new QueryBuilder instance
func New() *QueryBuilder {
	// Get the driver name from the database connection
	driverName := grammar.DriverMySQL // Default to MySQL

	return &QueryBuilder{
		queryType:         types.SELECT,
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
		aggregateFuncs:    []types.AggregateExpr{},
		subqueries:        make(map[string]*QueryBuilder),
		withCTEs:          []types.CTE{},
		jsonExpressions:   []types.JSONExpression{},
		mathExpressions:   []types.MathExpression{},
		conditionalExprs:  []types.ConditionalExpr{},
		stringFuncs:       []types.StringFunction{},
		dateFuncs:         []types.DateFunction{},
		indexHints:        []string{},
		isUsingDistinct:   false,
		isForUpdate:       false,
		isLockInShareMode: false,
		isHighPriority:    false,
		isStraightJoin:    false,
		isCalcFoundRows:   false,
	}
}

// Table creates a new QueryBuilder instance for a specific table
func Table(table string) *QueryBuilder {
	qb := New()
	qb.table = table
	return qb
}

// Reset resets the QueryBuilder instance
func (qb *QueryBuilder) Reset() {
	qb.queryType = types.SELECT
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
	qb.aggregateFuncs = []types.AggregateExpr{}
	qb.subqueries = make(map[string]*QueryBuilder)
	qb.withCTEs = []types.CTE{}
	qb.jsonExpressions = []types.JSONExpression{}
	qb.mathExpressions = []types.MathExpression{}
	qb.conditionalExprs = []types.ConditionalExpr{}
	qb.stringFuncs = []types.StringFunction{}
	qb.dateFuncs = []types.DateFunction{}
	qb.indexHints = []string{}
	qb.lockType = ""
	qb.forceIndex = ""
	qb.useIndex = ""
	qb.ignoreIndex = ""
	qb.procedure = ""
	qb.procedureParams = nil
	qb.isUsingDistinct = false
	qb.isForUpdate = false
	qb.isLockInShareMode = false
	qb.isHighPriority = false
	qb.isStraightJoin = false
	qb.isCalcFoundRows = false
	qb.comment = ""
	qb.ctes = []types.CTE{}
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
		AggregateFuncs:    qb.aggregateFuncs,
		Subqueries:        make(map[string]any), // Convert subqueries to any
		WithCTEs:          qb.withCTEs,
		JSONExpressions:   qb.jsonExpressions,
		MathExpressions:   qb.mathExpressions,
		ConditionalExprs:  qb.conditionalExprs,
		StringFuncs:       qb.stringFuncs,
		DateFuncs:         qb.dateFuncs,
		IndexHints:        qb.indexHints,
		ForceIndex:        qb.forceIndex,
		UseIndex:          qb.useIndex,
		IgnoreIndex:       qb.ignoreIndex,
		IsUsingDistinct:   qb.isUsingDistinct,
		IsForUpdate:       qb.isForUpdate,
		IsLockInShareMode: qb.isLockInShareMode,
		IsHighPriority:    qb.isHighPriority,
		IsStraightJoin:    qb.isStraightJoin,
		IsCalcFoundRows:   qb.isCalcFoundRows,
	}
}

// ToSQL compiles the query to SQL
func (qb *QueryBuilder) ToSQL() (string, []any, error) {
	return qb.grammar.CompileSelect(qb.GetData())
}

// Raw creates a new raw SQL expression
func Raw(sql string, bindings ...any) *types.Expression {
	return &types.Expression{
		SQL:      sql,
		Bindings: bindings,
	}
}
