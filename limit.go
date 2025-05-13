package xqb

// Limit adds a LIMIT clause
func (qb *QueryBuilder) Limit(limit int) *QueryBuilder {
	qb.limit = limit
	return qb
}

// Offset adds an OFFSET clause
func (qb *QueryBuilder) Offset(offset int) *QueryBuilder {
	qb.offset = offset
	return qb
}

// Skip is an alias for Offset
func (qb *QueryBuilder) Skip(offset int) *QueryBuilder {
	return qb.Offset(offset)
}

// Take is an alias for Limit
func (qb *QueryBuilder) Take(limit int) *QueryBuilder {
	return qb.Limit(limit)
}

// ForPage adds LIMIT and OFFSET clauses for pagination
func (qb *QueryBuilder) ForPage(page int, perPage int) *QueryBuilder {
	return qb.Skip((page - 1) * perPage).Take(perPage)
}
