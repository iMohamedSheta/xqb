package types

type Dialect string

const (
	DialectMySql    Dialect = "mysql"
	DialectPostgres Dialect = "postgres"
)

func (d Dialect) String() string {
	return string(d)
}
