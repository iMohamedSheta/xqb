package xqb_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/iMohamedSheta/xqb"
	"github.com/stretchr/testify/assert"
)

type User struct {
	ID        int            `xqb:"id"`
	Name      string         `xqb:"name"`
	Email     sql.NullString `xqb:"email"`
	Active    sql.NullBool   `xqb:"active"`
	CreatedAt sql.NullTime   `xqb:"created_at"`
	Password  string         `xqb:"-"` // should be ignored
}

func (User) Table() string {
	return "users"
}

func Test_Query_WithModel(t *testing.T) {
	sql, bindings, _ := xqb.Model(User{}).
		Select("id", "name", "email", "active", "created_at").
		Where("username", "=", "ali").
		OrWhere("username", "=", "mohamed").
		Latest("created_at").
		Limit(1).
		AddSelect("password").
		ToSQL()

	assert.Equal(t, "SELECT id, name, email, active, created_at, password FROM users WHERE username = ? OR username = ? ORDER BY created_at DESC LIMIT 1", sql)
	assert.Equal(t, []any{"ali", "mohamed"}, bindings)
	now := time.Now()
	data := map[string]any{
		"id":         1,
		"name":       "Ali",
		"email":      "ali@example.com",
		"active":     true,
		"created_at": now,
		"password":   "super-secret", // should not be set
	}

	var user User
	err := xqb.Bind(data, &user)

	assert.NoError(t, err)
	assert.Equal(t, 1, user.ID)
	assert.Equal(t, "Ali", user.Name)
	assert.Equal(t, "ali@example.com", user.Email.String)
	assert.Equal(t, true, user.Active.Bool)
	assert.Equal(t, now, user.CreatedAt.Time)
	assert.Empty(t, user.Password) // should be ignored
}

func TestBind_SingleModel(t *testing.T) {
	now := time.Now()
	data := map[string]any{
		"id":         1,
		"name":       "Ali",
		"email":      "ali@example.com",
		"active":     true,
		"created_at": now,
		"password":   "super-secret", // should not be set
	}

	var user User
	err := xqb.Bind(data, &user)

	assert.NoError(t, err)
	assert.Equal(t, 1, user.ID)
	assert.Equal(t, "Ali", user.Name)
	assert.True(t, user.Email.Valid)
	assert.Equal(t, "ali@example.com", user.Email.String)
	assert.True(t, user.Active.Bool)
	assert.True(t, user.CreatedAt.Valid)
	assert.Equal(t, now, user.CreatedAt.Time)
	assert.Empty(t, user.Password) // should be ignored
}

func TestBind_SliceModel(t *testing.T) {
	data := []map[string]any{
		{"id": 1, "name": "Ali"},
		{"id": 2, "name": "Sara"},
	}

	var users []User
	err := xqb.Bind(data, &users)

	assert.NoError(t, err)
	assert.Len(t, users, 2)
	assert.Equal(t, 1, users[0].ID)
	assert.Equal(t, "Ali", users[0].Name)
	assert.Equal(t, 2, users[1].ID)
	assert.Equal(t, "Sara", users[1].Name)
}

func TestBind_ErrorCases(t *testing.T) {
	var notPointer User
	err := xqb.Bind(map[string]any{}, notPointer)
	assert.Error(t, err)

	var invalidType string
	err = xqb.Bind(map[string]any{}, &invalidType)
	assert.Error(t, err)

	err = xqb.Bind([]map[string]any{}, &User{}) // dest is not a slice
	assert.Error(t, err)

	err = xqb.Bind(map[string]any{}, nil)
	assert.Error(t, err)
}
