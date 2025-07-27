package xqb

import (
	"database/sql"
	"errors"
	"fmt"
	"sync"

	xqbErr "github.com/iMohamedSheta/xqb/shared/errors"
	"github.com/iMohamedSheta/xqb/shared/types"
)

type Connection struct {
	Name    string
	DB      *sql.DB
	Dialect types.Dialect
}

type DBM struct {
	mu                sync.RWMutex
	defaultConnection string
	connections       map[string]*Connection
}

var (
	managerOnce     sync.Once
	managerInstance *DBM
)

func DBManager() *DBM {
	managerOnce.Do(func() {
		managerInstance = &DBM{
			defaultConnection: "default",
			connections:       make(map[string]*Connection),
		}
	})
	return managerInstance
}

func (m *DBM) GetConnections() map[string]*Connection {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.connections
}

func (m *DBM) SetConnections(conns map[string]*Connection) error {
	if len(conns) == 0 {
		return fmt.Errorf("%w: invalid connections", xqbErr.ErrNoConnection)
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.connections = conns
	return nil
}

func (m *DBM) SetConnection(conn *Connection) error {
	if conn == nil || conn.Name == "" {
		return fmt.Errorf("%w: invalid connection", xqbErr.ErrNoConnection)
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.connections[conn.Name] = conn
	return nil
}

func (m *DBM) Connection(name string) (*Connection, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	conn, ok := m.connections[name]
	if !ok || conn == nil {
		return nil, fmt.Errorf("%w: connection %q not found", xqbErr.ErrNoConnection, name)
	}
	return conn, nil
}

func (m *DBM) ConnectionDB(name string) (*sql.DB, error) {
	conn, err := m.Connection(name)
	if err != nil {
		return nil, err
	}
	if conn.DB == nil {
		return nil, fmt.Errorf("%w: connection %q has no sql.DB", xqbErr.ErrNoConnection, name)
	}
	return conn.DB, nil
}

func (m *DBM) SetDialect(name string, dialect types.Dialect) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	conn, ok := m.connections[name]
	if !ok {
		return fmt.Errorf("%w: unknown connection %q", xqbErr.ErrNoConnection, name)
	}
	conn.Dialect = dialect
	return nil
}

func (m *DBM) GetDialect(name string) (types.Dialect, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	conn, ok := m.connections[name]
	if !ok || conn == nil {
		return "", fmt.Errorf("%w: dialect for connection %q not found", xqbErr.ErrNoConnection, name)
	}
	return conn.Dialect, nil
}

func (m *DBM) HasConnection(name string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	conn, ok := m.connections[name]
	return ok && conn != nil && conn.DB != nil
}

func (m *DBM) GetDefaultConnection() (*Connection, error) {
	if m.defaultConnection == "" {
		return nil, fmt.Errorf("%w: no default connection set", xqbErr.ErrNoConnection)
	}
	return m.Connection(m.defaultConnection)
}

func (m *DBM) GetDefaultConnectionName() string {
	return m.defaultConnection
}

func (m *DBM) SetDefaultConnection(name string) error {
	if !m.HasConnection(name) {
		return fmt.Errorf("%w: connection %q not found", xqbErr.ErrNoConnection, name)
	}
	m.defaultConnection = name
	return nil
}

func (m *DBM) CloseConnection(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	conn, ok := m.connections[name]
	if !ok || conn == nil {
		return fmt.Errorf("%w: connection %q not found", xqbErr.ErrNoConnection, name)
	}
	if conn.DB != nil {
		if err := conn.DB.Close(); err != nil {
			return err
		}
	}
	delete(m.connections, name)
	return nil
}

func (m *DBM) CloseAllConnections() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	var errs []error
	for name, conn := range m.connections {
		if conn != nil && conn.DB != nil {
			if err := conn.DB.Close(); err != nil {
				errs = append(errs, fmt.Errorf("%s: %v", name, err))
			}
		}
	}
	m.connections = make(map[string]*Connection)
	if len(errs) > 0 {
		return fmt.Errorf("%w: %v", xqbErr.ErrClosingConnection, errors.Join(errs...))
	}
	return nil
}

// --- Global Helpers (keep your current API) ---

func AddConnection(conn *Connection) error {
	return DBManager().SetConnection(conn)
}

func GetConnection(name string) (*Connection, error) {
	return DBManager().Connection(name)
}

func GetConnectionDB(name string) (*sql.DB, error) {
	return DBManager().ConnectionDB(name)
}

func GetDefaultConnection() (*Connection, error) {
	return DBManager().GetDefaultConnection()
}

func SetDefaultConnection(name string) error {
	return DBManager().SetDefaultConnection(name)
}

func HasConnection(name string) bool {
	return DBManager().HasConnection(name)
}

func Close(name string) error {
	return DBManager().CloseConnection(name)
}

func CloseAll() error {
	return DBManager().CloseAllConnections()
}
