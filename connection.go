package xqb

import (
	"database/sql"
	"errors"
	"fmt"
	"sync"
)

var (
	errNoConnection = errors.New("no database connection established")
)

var (
	managerInstance *DBM
	managerOnce     sync.Once
)

// DBManager handles multiple named database connections
type DBM struct {
	defaultConnection string
	connections       map[string]*sql.DB
	mu                sync.RWMutex
}

// db returns the singleton DBManage instance
func DBManager() *DBM {
	managerOnce.Do(func() {
		managerInstance = &DBM{
			connections: make(map[string]*sql.DB),
		}
	})
	return managerInstance
}

// setConnection sets a named DB connection
func (m *DBM) SetConnection(name string, db *sql.DB) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.connections[name] = db
}

// connection returns a *sql.DB by name
func (m *DBM) Connection(name string) (*sql.DB, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	db, ok := m.connections[name]
	if !ok || db == nil {
		return nil, fmt.Errorf("%w: connection %s was not found in the xqb DBManager", errNoConnection, name)
	}
	return db, nil
}

// hasConnection checks if a connection exists
func (m *DBM) HasConnection(name string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	db, ok := m.connections[name]
	if !ok || db == nil {
		return false
	}
	return true
}

// GetDefaultConnection returns the default connection
func (m *DBM) GetDefaultConnection() (*sql.DB, error) {
	return m.Connection(m.defaultConnection)
}

// GetDefaultConnectionName returns the default connection name
func (m *DBM) GetDefaultConnectionName() string {
	return m.defaultConnection
}

// GetDefaultConnection returns the default connection
func GetDefaultConnection() (*sql.DB, error) {
	return DBManager().GetDefaultConnection()
}

// Connection returns a *sql.DB connection by name
func GetConnection(name string) (*sql.DB, error) {
	return DBManager().Connection(name)
}

func (m *DBM) SetDefaultConnection(name string) error {
	if !m.HasConnection(name) {
		return fmt.Errorf("%w: connection  %s was not found in xqb DBManager", errNoConnection, name)
	}

	m.defaultConnection = name
	return nil
}

func AddConnection(name string, db *sql.DB) {
	DBManager().SetConnection(name, db)
}

func HasConnection(name string) bool {
	return DBManager().HasConnection(name)
}

func SetDefault(name string) error {
	return DBManager().SetDefaultConnection(name)
}

func Close(name string) error {
	return DBManager().CloseConnection(name)
}

func CloseAll() error {
	return DBManager().CloseAllConnections()
}

func (m *DBM) CloseAllConnections() error {
	var errs []error
	for name := range m.connections {
		if err := m.closeConnection(name); err != nil {
			errs = append(errs, fmt.Errorf("connection %s - %s ", name, err))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("%w: %s", ErrClosingConnection, errors.Join(errs...))
	}

	return nil
}

// CloseConnection closes a named database connection and returns xqb.ErrClosingConnection if the connection could not be closed
func (m *DBM) CloseConnection(name string) error {
	err := m.closeConnection(name)
	if err != nil {
		return fmt.Errorf("%w: failed to close connection %s - %s", ErrClosingConnection, name, err)
	}
	return nil
}

// closeConnection closes a named database connection and returns any error encountered closing the connection
func (m *DBM) closeConnection(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.connections[name]; ok {
		if m.connections[name] != nil {
			err := m.connections[name].Close()
			if err != nil {
				return err
			}
		}
		delete(m.connections, name)
	}

	return nil
}
