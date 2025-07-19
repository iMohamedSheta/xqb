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
	managerInstance *DBManager
	managerOnce     sync.Once
)

// DBManager handles multiple named database connections
type DBManager struct {
	defaultConnection string
	connections       map[string]*sql.DB
	mu                sync.RWMutex
}

// db returns the singleton DBManage instance
func Manager() *DBManager {
	managerOnce.Do(func() {
		managerInstance = &DBManager{
			connections: make(map[string]*sql.DB),
		}
	})
	return managerInstance
}

// setConnection sets a named DB connection
func (m *DBManager) SetConnection(name string, db *sql.DB) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.connections[name] = db
}

// connection returns a *sql.DB by name
func (m *DBManager) connection(name string) (*sql.DB, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	db, ok := m.connections[name]
	if !ok || db == nil {
		return nil, fmt.Errorf("%w: %s", errNoConnection, name)
	}
	return db, nil
}

// closeConnection closes and removes a named DB connection
func (m *DBManager) CloseConnection(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.connections[name]; ok {
		if m.connections[name] != nil {
			_ = m.connections[name].Close()
		}
		delete(m.connections, name)
	}

	return nil
}

// hasConnection checks if a connection exists
func (m *DBManager) HasConnection(name string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, ok := m.connections[name]
	return ok
}

// Connection returns a *sql.DB connection by name
func Connection(name string) (*sql.DB, error) {
	dbManager := Manager()
	if name == "default" || name == "" || !dbManager.HasConnection(name) {
		name = dbManager.defaultConnection
	}
	return dbManager.connection(name)
}

func (m *DBManager) SetDefaultConnection(name string) {
	if name == "" || !m.HasConnection(name) {
		name = "default"
	}
	m.defaultConnection = name
}

func AddConnection(name string, db *sql.DB) {
	Manager().SetConnection(name, db)
}

func HasConnection(name string) bool {
	return Manager().HasConnection(name)
}

func SetDefault(name string) {
	Manager().SetDefaultConnection(name)
}

func Close(name string) error {
	return Manager().CloseConnection(name)
}

func CloseAll() error {
	return Manager().CloseAllConnections()
}

func (m *DBManager) CloseAllConnections() error {
	for name := range m.connections {
		if err := m.CloseConnection(name); err != nil {
			return err
		}
	}
	return nil
}
