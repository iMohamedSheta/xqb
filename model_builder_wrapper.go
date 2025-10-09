package xqb

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/iMohamedSheta/xqb/shared/types"
)

// ModelBuilder is a generic query builder that provides type-safe operations
type ModelBuilder[T any] struct {
	*QueryBuilder
}

func NewModel[T ModelInterface]() *ModelBuilder[T] {
	var model T
	return &ModelBuilder[T]{
		QueryBuilder: Table(model.Table()),
	}
}

func ModelQuery[T ModelInterface]() *ModelBuilder[T] {
	return NewModel[T]()
}

func Model[T ModelInterface]() *ModelBuilder[T] {
	return NewModel[T]()
}

func (mq *ModelBuilder[T]) Table(table string) *ModelBuilder[T] {
	mq.QueryBuilder.Table(table)
	return mq
}

func (mq *ModelBuilder[T]) Connection(connection string) *ModelBuilder[T] {
	mq.QueryBuilder.Connection(connection)
	return mq
}

func (mq *ModelBuilder[T]) WithContext(ctx context.Context) *ModelBuilder[T] {
	mq.QueryBuilder.WithContext(ctx)
	return mq
}

func (mq *ModelBuilder[T]) WithTx(tx *sql.Tx) *ModelBuilder[T] {
	mq.QueryBuilder.WithTx(tx)
	return mq
}

func (mq *ModelBuilder[T]) WithSettings(settings *QueryBuilderSettings) *ModelBuilder[T] {
	mq.QueryBuilder.WithSettings(settings)
	return mq
}

func (mq *ModelBuilder[T]) AllowDangerous() *ModelBuilder[T] {
	mq.QueryBuilder.AllowDangerous()
	return mq
}

// Clone returns a copy of the ModelBuilder
func (mq *ModelBuilder[T]) Clone() *ModelBuilder[T] {
	return &ModelBuilder[T]{
		QueryBuilder: mq.QueryBuilder.Clone(),
	}
}

func (mq *ModelBuilder[T]) SetDialect(dialect types.Dialect) *ModelBuilder[T] {
	mq.QueryBuilder.SetDialect(dialect)
	return mq
}

// Select Methods
func (mq *ModelBuilder[T]) Select(columns ...any) *ModelBuilder[T] {
	mq.QueryBuilder.Select(columns...)
	return mq
}

func (mq *ModelBuilder[T]) AddSelect(columns ...any) *ModelBuilder[T] {
	mq.QueryBuilder.AddSelect(columns...)
	return mq
}

func (mq *ModelBuilder[T]) SelectSub(subQuery *QueryBuilder, alias string) *ModelBuilder[T] {
	mq.QueryBuilder.SelectSub(subQuery, alias)
	return mq
}

func (mq *ModelBuilder[T]) SelectRaw(sql string, bindings ...any) *ModelBuilder[T] {
	mq.QueryBuilder.SelectRaw(sql, bindings...)
	return mq
}

func (mq *ModelBuilder[T]) AddSelectRaw(sql string, bindings ...any) *ModelBuilder[T] {
	mq.QueryBuilder.AddSelectRaw(sql, bindings...)
	return mq
}

func (qb *ModelBuilder[T]) From(table string) *ModelBuilder[T] {
	qb.QueryBuilder.From(table)
	return qb
}

func (qb *ModelBuilder[T]) FromSubquery(subQuery *QueryBuilder, alias string) *ModelBuilder[T] {
	qb.QueryBuilder.FromSubquery(subQuery, alias)
	return qb
}

func (mq *ModelBuilder[T]) Distinct() *ModelBuilder[T] {
	mq.QueryBuilder.Distinct()
	return mq
}

// =============================================================================
// Join Methods
// =============================================================================

func (mq *ModelBuilder[T]) Join(table string, condition string, values ...any) *ModelBuilder[T] {
	mq.QueryBuilder.Join(table, condition, values...)
	return mq
}

func (mq *ModelBuilder[T]) LeftJoin(table string, condition string, values ...any) *ModelBuilder[T] {
	mq.QueryBuilder.LeftJoin(table, condition, values...)
	return mq
}

func (mq *ModelBuilder[T]) RightJoin(table string, condition string, values ...any) *ModelBuilder[T] {
	mq.QueryBuilder.RightJoin(table, condition, values...)
	return mq
}

func (mq *ModelBuilder[T]) FullJoin(table string, condition string, values ...any) *ModelBuilder[T] {
	mq.QueryBuilder.FullJoin(table, condition, values...)
	return mq
}

func (mq *ModelBuilder[T]) CrossJoin(table string) *ModelBuilder[T] {
	mq.QueryBuilder.CrossJoin(table)
	return mq
}

func (mq *ModelBuilder[T]) JoinSub(sub *QueryBuilder, alias, condition string, values ...any) *ModelBuilder[T] {
	mq.QueryBuilder.JoinSub(sub, alias, condition, values...)
	return mq
}

func (mq *ModelBuilder[T]) LeftJoinSub(sub *QueryBuilder, alias, condition string, values ...any) *ModelBuilder[T] {
	mq.QueryBuilder.LeftJoinSub(sub, alias, condition, values...)
	return mq
}

func (mq *ModelBuilder[T]) RightJoinSub(sub *QueryBuilder, alias, condition string, values ...any) *ModelBuilder[T] {
	mq.QueryBuilder.RightJoinSub(sub, alias, condition, values...)
	return mq
}

func (mq *ModelBuilder[T]) FullJoinSub(sub *QueryBuilder, alias, condition string, values ...any) *ModelBuilder[T] {
	mq.QueryBuilder.FullJoinSub(sub, alias, condition, values...)
	return mq
}

func (mq *ModelBuilder[T]) CrossJoinSub(sub *QueryBuilder, alias string) *ModelBuilder[T] {
	mq.QueryBuilder.CrossJoinSub(sub, alias)
	return mq
}

func (mq *ModelBuilder[T]) JoinExpr(expr *types.Expression, condition any, values ...any) *ModelBuilder[T] {
	mq.QueryBuilder.JoinExpr(expr, condition, values...)
	return mq
}

func (mq *ModelBuilder[T]) LeftJoinExpr(expr *types.Expression, condition any, values ...any) *ModelBuilder[T] {
	mq.QueryBuilder.LeftJoinExpr(expr, condition, values...)
	return mq
}

func (mq *ModelBuilder[T]) RightJoinExpr(expr *types.Expression, condition any, values ...any) *ModelBuilder[T] {
	mq.QueryBuilder.RightJoinExpr(expr, condition, values...)
	return mq
}

func (mq *ModelBuilder[T]) FullJoinExpr(expr *types.Expression, condition any, values ...any) *ModelBuilder[T] {
	mq.QueryBuilder.FullJoinExpr(expr, condition, values...)
	return mq
}

func (mq *ModelBuilder[T]) CrossJoinExpr(expr *types.Expression) *ModelBuilder[T] {
	mq.QueryBuilder.CrossJoinExpr(expr)
	return mq
}

// ============================================================================
// ORDER BY Methods
// ============================================================================

func (mq *ModelBuilder[T]) OrderBy(column string, direction string) *ModelBuilder[T] {
	mq.QueryBuilder.OrderBy(column, direction)
	return mq
}

func (mq *ModelBuilder[T]) OrderByDesc(column string) *ModelBuilder[T] {
	mq.QueryBuilder.OrderByDesc(column)
	return mq
}

func (mq *ModelBuilder[T]) OrderByAsc(column string) *ModelBuilder[T] {
	mq.QueryBuilder.OrderByAsc(column)
	return mq
}

func (mq *ModelBuilder[T]) OrderByRaw(sql string, bindings ...any) *ModelBuilder[T] {
	mq.QueryBuilder.OrderByRaw(sql, bindings...)
	return mq
}

func (mq *ModelBuilder[T]) Latest(column string) *ModelBuilder[T] {
	mq.QueryBuilder.Latest(column)
	return mq
}

func (mq *ModelBuilder[T]) Oldest(column string) *ModelBuilder[T] {
	mq.QueryBuilder.Oldest(column)
	return mq
}

// ============================================================================
// WHERE Methods
// ============================================================================

func (mq *ModelBuilder[T]) Where(column any, operator string, value any) *ModelBuilder[T] {
	mq.QueryBuilder.Where(column, operator, value)
	return mq
}

func (mq *ModelBuilder[T]) OrWhere(column any, operator string, value any) *ModelBuilder[T] {
	mq.QueryBuilder.OrWhere(column, operator, value)
	return mq
}

func (mq *ModelBuilder[T]) WhereValue(column string, operator string, value any) *ModelBuilder[T] {
	mq.QueryBuilder.WhereValue(column, operator, value)
	return mq
}

func (mq *ModelBuilder[T]) OrWhereValue(column string, operator string, value any) *ModelBuilder[T] {
	mq.QueryBuilder.OrWhereValue(column, operator, value)
	return mq
}

func (mq *ModelBuilder[T]) WhereSub(column string, operator string, sub *QueryBuilder) *ModelBuilder[T] {
	mq.QueryBuilder.WhereSub(column, operator, sub)
	return mq
}

func (mq *ModelBuilder[T]) OrWhereSub(column string, operator string, sub *QueryBuilder) *ModelBuilder[T] {
	mq.QueryBuilder.OrWhereSub(column, operator, sub)
	return mq
}

func (mq *ModelBuilder[T]) WhereExpr(column string, operator string, expr *types.Expression) *ModelBuilder[T] {
	mq.QueryBuilder.WhereExpr(column, operator, expr)
	return mq
}

func (mq *ModelBuilder[T]) OrWhereExpr(column string, operator string, expr *types.Expression) *ModelBuilder[T] {
	mq.QueryBuilder.OrWhereExpr(column, operator, expr)
	return mq
}

func (mq *ModelBuilder[T]) WhereRaw(sql string, bindings ...any) *ModelBuilder[T] {
	mq.QueryBuilder.WhereRaw(sql, bindings...)
	return mq
}

func (mq *ModelBuilder[T]) OrWhereRaw(sql string, bindings ...any) *ModelBuilder[T] {
	mq.QueryBuilder.OrWhereRaw(sql, bindings...)
	return mq
}

func (mq *ModelBuilder[T]) WhereNull(column string) *ModelBuilder[T] {
	mq.QueryBuilder.WhereNull(column)
	return mq
}

func (mq *ModelBuilder[T]) OrWhereNull(column string) *ModelBuilder[T] {
	mq.QueryBuilder.OrWhereNull(column)
	return mq
}

func (mq *ModelBuilder[T]) WhereNotNull(column string) *ModelBuilder[T] {
	mq.QueryBuilder.WhereNotNull(column)
	return mq
}

func (mq *ModelBuilder[T]) OrWhereNotNull(column string) *ModelBuilder[T] {
	mq.QueryBuilder.OrWhereNotNull(column)
	return mq
}

func (mq *ModelBuilder[T]) WhereIn(column string, values any) *ModelBuilder[T] {
	mq.QueryBuilder.WhereIn(column, values)
	return mq
}

func (mq *ModelBuilder[T]) OrWhereIn(column string, values any) *ModelBuilder[T] {
	mq.QueryBuilder.OrWhereIn(column, values)
	return mq
}

func (mq *ModelBuilder[T]) WhereNotIn(column string, values any) *ModelBuilder[T] {
	mq.QueryBuilder.WhereNotIn(column, values)
	return mq
}

func (mq *ModelBuilder[T]) OrWhereNotIn(column string, values any) *ModelBuilder[T] {
	mq.QueryBuilder.OrWhereNotIn(column, values)
	return mq
}

func (mq *ModelBuilder[T]) WhereInQuery(column string, sub *QueryBuilder) *ModelBuilder[T] {
	mq.QueryBuilder.WhereInQuery(column, sub)
	return mq
}

func (mq *ModelBuilder[T]) OrWhereInQuery(column string, sub *QueryBuilder) *ModelBuilder[T] {
	mq.QueryBuilder.OrWhereInQuery(column, sub)
	return mq
}

func (mq *ModelBuilder[T]) WhereNotInQuery(column string, sub *QueryBuilder) *ModelBuilder[T] {
	mq.QueryBuilder.WhereNotInQuery(column, sub)
	return mq
}

func (mq *ModelBuilder[T]) OrWhereNotInQuery(column string, sub *QueryBuilder) *ModelBuilder[T] {
	mq.QueryBuilder.OrWhereNotInQuery(column, sub)
	return mq
}

func (mq *ModelBuilder[T]) WhereTrue(column string) *ModelBuilder[T] {
	mq.QueryBuilder.WhereTrue(column)
	return mq
}

func (mq *ModelBuilder[T]) OrWhereTrue(column string) *ModelBuilder[T] {
	mq.QueryBuilder.OrWhereTrue(column)
	return mq
}

func (mq *ModelBuilder[T]) WhereFalse(column string) *ModelBuilder[T] {
	mq.QueryBuilder.WhereFalse(column)
	return mq
}

func (mq *ModelBuilder[T]) OrWhereFalse(column string) *ModelBuilder[T] {
	mq.QueryBuilder.OrWhereFalse(column)
	return mq
}

func (mq *ModelBuilder[T]) WhereBetween(column string, min, max any) *ModelBuilder[T] {
	mq.QueryBuilder.WhereBetween(column, min, max)
	return mq
}

func (mq *ModelBuilder[T]) OrWhereBetween(column string, min, max any) *ModelBuilder[T] {
	mq.QueryBuilder.OrWhereBetween(column, min, max)
	return mq
}

func (mq *ModelBuilder[T]) WhereNotBetween(column string, min, max any) *ModelBuilder[T] {
	mq.QueryBuilder.WhereNotBetween(column, min, max)
	return mq
}

func (mq *ModelBuilder[T]) OrWhereNotBetween(column string, min, max any) *ModelBuilder[T] {
	mq.QueryBuilder.OrWhereNotBetween(column, min, max)
	return mq
}

func (mq *ModelBuilder[T]) WhereExists(subquery any) *ModelBuilder[T] {
	mq.QueryBuilder.WhereExists(subquery)
	return mq
}

func (mq *ModelBuilder[T]) OrWhereExists(subquery any) *ModelBuilder[T] {
	mq.QueryBuilder.OrWhereExists(subquery)
	return mq
}

func (mq *ModelBuilder[T]) WhereNotExists(subquery any) *ModelBuilder[T] {
	mq.QueryBuilder.WhereNotExists(subquery)
	return mq
}

func (mq *ModelBuilder[T]) OrWhereNotExists(subquery any) *ModelBuilder[T] {
	mq.QueryBuilder.OrWhereNotExists(subquery)
	return mq
}

func (mq *ModelBuilder[T]) WhereGroup(fn func(qb *QueryBuilder)) *ModelBuilder[T] {
	mq.QueryBuilder.WhereGroup(fn)
	return mq
}

func (mq *ModelBuilder[T]) OrWhereGroup(fn func(qb *QueryBuilder)) *ModelBuilder[T] {
	mq.QueryBuilder.OrWhereGroup(fn)
	return mq
}

// ============================================================================
// GroupBy && Having && Limit Methods
// ============================================================================

func (mq *ModelBuilder[T]) GroupBy(columns ...any) *ModelBuilder[T] {
	mq.QueryBuilder.GroupBy(columns...)
	return mq
}

func (mq *ModelBuilder[T]) Having(column any, operator string, value any) *ModelBuilder[T] {
	mq.QueryBuilder.Having(column, operator, value)
	return mq
}

func (mq *ModelBuilder[T]) HavingRaw(sql string, bindings ...any) *ModelBuilder[T] {
	mq.QueryBuilder.HavingRaw(sql, bindings...)
	return mq
}

func (mq *ModelBuilder[T]) OrHaving(column any, operator string, value any) *ModelBuilder[T] {
	mq.QueryBuilder.OrHaving(column, operator, value)
	return mq
}

func (mq *ModelBuilder[T]) OrHavingRaw(sql string, bindings ...any) *ModelBuilder[T] {
	mq.QueryBuilder.OrHavingRaw(sql, bindings...)
	return mq
}

// Limit & Offset Methods
func (mq *ModelBuilder[T]) Limit(limit int) *ModelBuilder[T] {
	mq.QueryBuilder.Limit(limit)
	return mq
}

func (mq *ModelBuilder[T]) Offset(offset int) *ModelBuilder[T] {
	mq.QueryBuilder.Offset(offset)
	return mq
}

func (mq *ModelBuilder[T]) Take(limit int) *ModelBuilder[T] {
	mq.QueryBuilder.Take(limit)
	return mq
}

func (mq *ModelBuilder[T]) Skip(offset int) *ModelBuilder[T] {
	mq.QueryBuilder.Skip(offset)
	return mq
}

// ForPage adds LIMIT and OFFSET clauses for pagination
func (qb *ModelBuilder[T]) ForPage(page int, perPage int) *ModelBuilder[T] {
	return qb.Skip((page - 1) * perPage).Take(perPage)
}

// ============================================================================
// LOCK Methods
// ============================================================================

func (mq *ModelBuilder[T]) LockForUpdate() *ModelBuilder[T] {
	mq.QueryBuilder.LockForUpdate()
	return mq
}

func (mq *ModelBuilder[T]) SharedLock() *ModelBuilder[T] {
	mq.QueryBuilder.SharedLock()
	return mq
}

func (mq *ModelBuilder[T]) LockNoKeyUpdate() *ModelBuilder[T] {
	mq.QueryBuilder.LockNoKeyUpdate()
	return mq
}

func (mq *ModelBuilder[T]) LockKeyShare() *ModelBuilder[T] {
	mq.QueryBuilder.LockKeyShare()
	return mq
}

func (mq *ModelBuilder[T]) NoWaitLocked() *ModelBuilder[T] {
	mq.QueryBuilder.NoWaitLocked()
	return mq
}

func (mq *ModelBuilder[T]) SkipLocked() *ModelBuilder[T] {
	mq.QueryBuilder.SkipLocked()
	return mq
}

// ============================================================================
// UNION Methods
// ============================================================================

func (mq *ModelBuilder[T]) Union(secondaryQuery ...*QueryBuilder) *ModelBuilder[T] {
	mq.QueryBuilder.Union(secondaryQuery...)
	return mq
}

func (mq *ModelBuilder[T]) UnionRaw(sql string, bindings ...any) *ModelBuilder[T] {
	mq.QueryBuilder.UnionRaw(sql, bindings...)
	return mq
}

func (mq *ModelBuilder[T]) UnionAll(secondaryQuery ...*QueryBuilder) *ModelBuilder[T] {
	mq.QueryBuilder.UnionAll(secondaryQuery...)
	return mq
}

func (mq *ModelBuilder[T]) UnionAllRaw(sql string, bindings ...any) *ModelBuilder[T] {
	mq.QueryBuilder.UnionAllRaw(sql, bindings...)
	return mq
}

func (mq *ModelBuilder[T]) ExceptUnion(secondaryQuery ...*QueryBuilder) *ModelBuilder[T] {
	mq.QueryBuilder.ExceptUnion(secondaryQuery...)
	return mq
}

func (mq *ModelBuilder[T]) ExceptUnionAll(secondaryQuery ...*QueryBuilder) *ModelBuilder[T] {
	mq.QueryBuilder.ExceptUnionAll(secondaryQuery...)
	return mq
}

func (mq *ModelBuilder[T]) ExceptUnionRaw(sql string, all bool, bindings ...any) *ModelBuilder[T] {
	mq.QueryBuilder.ExceptUnionRaw(sql, all, bindings...)
	return mq
}

func (mq *ModelBuilder[T]) IntersectUnion(secondaryQuery ...*QueryBuilder) *ModelBuilder[T] {
	mq.QueryBuilder.IntersectUnion(secondaryQuery...)
	return mq
}

func (mq *ModelBuilder[T]) IntersectUnionAll(secondaryQuery ...*QueryBuilder) *ModelBuilder[T] {
	mq.QueryBuilder.IntersectUnionAll(secondaryQuery...)
	return mq
}

func (mq *ModelBuilder[T]) IntersectUnionRaw(sql string, all bool, bindings ...any) *ModelBuilder[T] {
	mq.QueryBuilder.IntersectUnionRaw(sql, all, bindings...)
	return mq
}

// ============================================================================
// Result Methods - These need wrapping to return T instead of map[string]any
// ============================================================================

// Get executes the query and returns all results as a slice of T
func (mq *ModelBuilder[T]) Get() ([]T, error) {
	data, err := mq.QueryBuilder.Get()
	if err != nil {
		return nil, err
	}

	var results []T
	if err := Bind(data, &results); err != nil {
		return nil, fmt.Errorf("%w: Get() failed to bind results: %v", ErrInvalidResult, err)
	}

	return results, nil
}

// First executes the query and returns the first result as type T
func (mq *ModelBuilder[T]) First() (*T, error) {
	data, err := mq.QueryBuilder.First()
	if err != nil {
		return nil, err
	}

	var model T
	if err := Bind(data, &model); err != nil {
		return nil, fmt.Errorf("%w: First() failed to bind result: %v", ErrInvalidResult, err)
	}

	return &model, nil
}

// Find finds the first result by ID
func (mq *ModelBuilder[T]) Find(id any) (*T, error) {
	data, err := mq.QueryBuilder.Find(id)
	if err != nil {
		return nil, err
	}

	var model T
	if err := Bind(data, &model); err != nil {
		return nil, fmt.Errorf("%w: Find() failed to bind result: %v", ErrInvalidResult, err)
	}

	return &model, nil
}

// FindOrFail finds the first result by ID or returns a "not found" error
func (mq *ModelBuilder[T]) FindOrFail(id any) (*T, error) {
	data, err := mq.QueryBuilder.FindOrFail(id)
	if err != nil {
		return nil, err
	}

	var model T
	if err := Bind(data, &model); err != nil {
		return nil, fmt.Errorf("%w: FindOrFail() failed to bind result: %v", ErrInvalidResult, err)
	}

	return &model, nil
}

// Paginate returns paginated results with optional count metadata
func (mq *ModelBuilder[T]) Paginate(perPage int, page int, countBy string) ([]T, map[string]any, error) {
	data, meta, err := mq.QueryBuilder.Paginate(perPage, page, countBy)
	if err != nil {
		return nil, nil, err
	}

	var results []T
	if err := Bind(data, &results); err != nil {
		return nil, nil, fmt.Errorf("%w: Paginate() failed to bind results: %v", ErrInvalidResult, err)
	}

	return results, meta, nil
}

// Chunks processes results in batches and calls the closure for each chunk
func (mq *ModelBuilder[T]) Chunks(chunkSize int, closure func(results []T) error) error {
	return mq.QueryBuilder.Chunks(chunkSize, func(results []map[string]any) error {
		var models []T
		if err := Bind(results, &models); err != nil {
			return fmt.Errorf("%w: Chunks() failed to bind results: %v", ErrInvalidResult, err)
		}
		return closure(models)
	})
}

// QB provides access to raw QueryBuilder methods while maintaining type safety
func (mq *ModelBuilder[T]) Q(fn func(*QueryBuilder) *QueryBuilder) *ModelBuilder[T] {
	fn(mq.QueryBuilder)
	return mq
}
