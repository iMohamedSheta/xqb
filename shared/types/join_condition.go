package types

type JoinConditionKind int

const (
	JoinConditionOn    JoinConditionKind = iota // ON col1 op col2   (column = column, no binding)
	JoinConditionWhere                          // ON col op ?        (column = value, binding)
	JoinConditionRaw                            // raw SQL fragment
	JoinConditionGroup                          // nested (...)
)

type JoinCondition struct {
	Kind      JoinConditionKind
	First     string // left-hand column (On, Where)
	Operator  string
	Second    string // right-hand column (On only)
	Value     any    // bound value (Where only)
	Connector WhereConditionEnum
	Raw       ExpressionInterface
	Group     []*JoinCondition // nested group (Group only)
}
