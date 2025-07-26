package types

type Driver string

const (
	DriverMySql    Driver = "mysql"
	DriverPostgres Driver = "postgres"
)

func (d Driver) String() string {
	return string(d)
}
