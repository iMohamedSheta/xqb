package xqb_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/iMohamedSheta/xqb"
	"github.com/iMohamedSheta/xqb/shared/types"
	"github.com/stretchr/testify/assert"
)

type User struct {
	ID        int            `xqb:"id"`
	Name      string         `xqb:"name"`
	Email     sql.NullString `xqb:"email"`
	Active    sql.NullBool   `xqb:"active"`
	CreatedAt sql.NullTime   `xqb:"created_at"`
	Password  string         `xqb:"-"` // ignored
}

func (User) Table() string {
	return "users"
}

func Test_Query_WithModel(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		sql, bindings, err := xqb.Model(User{}).SetDialect(dialect).
			Select("id", "name", "email", "active", "created_at").
			Where("username", "=", "ali").
			OrWhere("username", "=", "mohamed").
			Latest("created_at").
			Limit(1).
			AddSelect("password").
			ToSql()

		expectedSql := map[types.Dialect]string{
			types.DialectMySql:    "SELECT `id`, `name`, `email`, `active`, `created_at`, `password` FROM `users` WHERE `username` = ? OR `username` = ? ORDER BY `created_at` DESC LIMIT 1",
			types.DialectPostgres: `SELECT "id", "name", "email", "active", "created_at", "password" FROM "users" WHERE "username" = $1 OR "username" = $2 ORDER BY "created_at" DESC LIMIT 1`,
		}

		assert.Equal(t, expectedSql[dialect], sql)
		assert.Equal(t, []any{"ali", "mohamed"}, bindings)
		assert.NoError(t, err)

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
		err = xqb.Bind(data, &user)

		assert.NoError(t, err)
		assert.Equal(t, 1, user.ID)
		assert.Equal(t, "Ali", user.Name)
		assert.Equal(t, "ali@example.com", user.Email.String)
		assert.Equal(t, true, user.Active.Bool)
		assert.Equal(t, now, user.CreatedAt.Time)
		assert.Empty(t, user.Password) // should be ignored
	})
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

	err = xqb.Bind(map[string]any{}, nil)
	assert.Error(t, err)
}

func TestBind_NestedStructs(t *testing.T) {
	type Address struct {
		Street string `xqb:"street"`
		City   string `xqb:"city"`
		State  string `xqb:"state"`
	}

	type User struct {
		ID       int     `xqb:"id"`
		Name     string  `xqb:"name"`
		Email    string  `xqb:"email"`
		Address  Address `xqb:"address"`
		Password string  `xqb:"-"` // ignored
	}

	data := map[string]any{
		"id":             1,
		"name":           "Ali",
		"email":          "ali@example.com",
		"address.street": "123 Main St",
		"address.city":   "Cairo",
		"address.state":  "Cairo Governorate",
		"password":       "super-secret", // should be ignored
	}

	var user User
	err := xqb.Bind(data, &user)
	assert.NoError(t, err)
	assert.Equal(t, 1, user.ID)
	assert.Equal(t, "Ali", user.Name)
	assert.Equal(t, "ali@example.com", user.Email)
	assert.Equal(t, "123 Main St", user.Address.Street)
	assert.Equal(t, "Cairo", user.Address.City)
	assert.Equal(t, "Cairo Governorate", user.Address.State)
	assert.Empty(t, user.Password)
}

func TestBind_NestedPointerStructs(t *testing.T) {
	type Address struct {
		Street string `xqb:"street"`
		City   string `xqb:"city"`
		State  string `xqb:"state"`
	}

	type User struct {
		ID      int      `xqb:"id"`
		Name    string   `xqb:"name"`
		Email   string   `xqb:"email"`
		Address *Address `xqb:"address"` // pointer to nested struct
	}

	data := map[string]any{
		"id":             2,
		"name":           "Sara",
		"email":          "sara@example.com",
		"address.street": "456 Elm St",
		"address.city":   "Alexandria",
		"address.state":  "Alexandria Governorate",
	}

	var user User
	err := xqb.Bind(data, &user)
	assert.NoError(t, err)
	assert.Equal(t, 2, user.ID)
	assert.Equal(t, "Sara", user.Name)
	assert.Equal(t, "sara@example.com", user.Email)
	assert.NotNil(t, user.Address)
	assert.Equal(t, "456 Elm St", user.Address.Street)
	assert.Equal(t, "Alexandria", user.Address.City)
	assert.Equal(t, "Alexandria Governorate", user.Address.State)
}

func TestBind_CastsAndTimestamps(t *testing.T) {
	type User struct {
		ID        int          `xqb:"id"`
		Name      string       `xqb:"name"`
		CreatedAt sql.NullTime `xqb:"created_at"`
		UpdatedAt sql.NullTime `xqb:"updated_at"`
	}

	now := time.Now()
	data := map[string]any{
		"id":         1,
		"name":       "Ali",
		"created_at": now,
		"updated_at": now,
	}

	var user User
	err := xqb.Bind(data, &user)
	assert.NoError(t, err)
	assert.Equal(t, now, user.CreatedAt.Time)
	assert.Equal(t, now, user.UpdatedAt.Time)
	assert.True(t, user.CreatedAt.Valid)
	assert.True(t, user.UpdatedAt.Valid)
}

func TestBind_CastsAndTime(t *testing.T) {
	type User struct {
		ID        int       `xqb:"id"`
		Name      string    `xqb:"name"`
		CreatedAt time.Time `xqb:"created_at"`
		UpdatedAt time.Time `xqb:"updated_at"`
	}

	now := time.Now()
	data := map[string]any{
		"id":         1,
		"name":       "Ali",
		"created_at": now,
		"updated_at": now,
	}

	var user User
	err := xqb.Bind(data, &user)
	assert.NoError(t, err)
	assert.Equal(t, now, user.CreatedAt)
	assert.Equal(t, now, user.UpdatedAt)
}

func TestBind_SoftDeletes(t *testing.T) {
	type User struct {
		ID        int          `xqb:"id"`
		Name      string       `xqb:"name"`
		DeletedAt sql.NullTime `xqb:"deleted_at"`
	}

	now := time.Now()
	data := map[string]any{
		"id":         1,
		"name":       "Ali",
		"deleted_at": now,
	}

	var user User
	err := xqb.Bind(data, &user)
	assert.NoError(t, err)
	assert.Equal(t, now, user.DeletedAt.Time)
	assert.True(t, user.DeletedAt.Valid)
}

func TestBind_MassAssignmentProtection(t *testing.T) {
	type User struct {
		ID       int    `xqb:"id"`
		Name     string `xqb:"name"`
		Password string `xqb:"-"` // ignored
	}

	data := map[string]any{
		"id":       1,
		"name":     "Ali",
		"password": "secret", // should be ignored
	}

	var user User
	err := xqb.Bind(data, &user)
	assert.NoError(t, err)
	assert.Equal(t, 1, user.ID)
	assert.Equal(t, "Ali", user.Name)
	assert.Empty(t, user.Password) // protected
}

func TestBind_NestedSliceRelations(t *testing.T) {
	type Post struct {
		ID    int    `xqb:"id"`
		Title string `xqb:"title"`
	}

	type User struct {
		ID    int    `xqb:"id"`
		Name  string `xqb:"name"`
		Posts []Post `xqb:"posts"`
	}

	data := map[string]any{
		"id":   1,
		"name": "Ali",
		"posts": []map[string]any{
			{
				"title": "First Post",
			},
			{
				"title": "Second Post",
			},
		},
	}

	var user User
	err := xqb.Bind(data, &user)
	assert.NoError(t, err)
	assert.Len(t, user.Posts, 2)
	assert.Equal(t, "First Post", user.Posts[0].Title)
	assert.Equal(t, "Second Post", user.Posts[1].Title)
}

func TestBind_NestedSliceRelationss(t *testing.T) {
	type Post struct {
		ID    int    `xqb:"id"`
		Name  string `xqb:"name"`
		Title string `xqb:"title"`
	}

	type User struct {
		ID    int    `xqb:"id"`
		Name  string `xqb:"name"`
		Posts []Post `xqb:"posts"`
	}

	data := []map[string]any{
		{
			"id":          1,
			"name":        "Ali",
			"posts_title": "First Post",
		},
		{
			"id":          2,
			"name":        "Ahmed",
			"posts_title": "Second Post",
			"posts_name":  "Third Post",
		},
	}

	var user User
	err := xqb.Bind(data, &user)
	assert.NoError(t, err)
	assert.Len(t, user.Posts, 2)
	assert.Equal(t, "First Post", user.Posts[0].Title)
	assert.Equal(t, "Second Post", user.Posts[1].Title)
}

func TestBind_SpecialCase(t *testing.T) {
	type Post struct {
		ID     int    `xqb:"id"`
		Name   string `xqb:"name"`
		Title  string `xqb:"title"`
		Serial string `xqb:"serial"`
	}

	type User struct {
		ID         int    `xqb:"id"`
		Name       string `xqb:"name"`
		PostSerial string `xqb:"post_serial"`
		Posts      []Post `xqb:"posts" table:"posts"`
	}

	// Fixed: Same user with multiple posts
	data := []map[string]any{
		{
			"id":           1,
			"name":         "Ali",
			"post_serial":  "main_serial",
			"posts_id":     10,
			"posts_title":  "First Post",
			"posts_name":   "Post One",
			"posts_serial": "55555",
		},
		{
			"id":           1, // Same user
			"name":         "Ali",
			"post_serial":  "main_serial",
			"posts_id":     20,
			"posts_title":  "Second Post",
			"posts_name":   "Post Two",
			"posts_serial": "66666",
		},
	}

	var user User
	err := xqb.Bind(data, &user)
	assert.NoError(t, err)
	assert.Equal(t, 1, user.ID)
	assert.Equal(t, "Ali", user.Name)
	assert.Equal(t, "main_serial", user.PostSerial)
	assert.Len(t, user.Posts, 2)
	assert.Equal(t, 10, user.Posts[0].ID)
	assert.Equal(t, "First Post", user.Posts[0].Title)
	assert.Equal(t, "Post One", user.Posts[0].Name)
	assert.Equal(t, "55555", user.Posts[0].Serial)
	assert.Equal(t, 20, user.Posts[1].ID)
	assert.Equal(t, "Second Post", user.Posts[1].Title)
	assert.Equal(t, "Post Two", user.Posts[1].Name)
	assert.Equal(t, "66666", user.Posts[1].Serial)
}

func TestBind_JsonColumns(t *testing.T) {
	type UserSettings struct {
		Logo  string `json:"logo"`
		Color string `json:"color"`
		Size  int    `json:"size"`
	}

	type User struct {
		ID       int          `xqb:"id"`
		Name     string       `xqb:"name"`
		Settings UserSettings `xqb:"settings"`
	}

	// Test Case 1: JSON as string
	t.Run("JSON as string", func(t *testing.T) {
		data := map[string]any{
			"id":       1,
			"name":     "John",
			"settings": `{"logo": "logo.png", "color": "blue", "size": 42}`,
		}

		var user User
		err := xqb.Bind(data, &user)
		assert.NoError(t, err)
		assert.Equal(t, 1, user.ID)
		assert.Equal(t, "John", user.Name)
		assert.Equal(t, "logo.png", user.Settings.Logo)
		assert.Equal(t, "blue", user.Settings.Color)
		assert.Equal(t, 42, user.Settings.Size)
	})

	// Test Case 2: JSON as map (already parsed)
	t.Run("JSON as map", func(t *testing.T) {
		data := map[string]any{
			"id":   2,
			"name": "Jane",
			"settings": map[string]any{
				"logo":  "jane.png",
				"color": "red",
				"size":  24,
			},
		}

		var user User
		err := xqb.Bind(data, &user)
		assert.NoError(t, err)
		assert.Equal(t, 2, user.ID)
		assert.Equal(t, "Jane", user.Name)
		assert.Equal(t, "jane.png", user.Settings.Logo)
		assert.Equal(t, "red", user.Settings.Color)
		assert.Equal(t, 24, user.Settings.Size)
	})

	// Test Case 3: Prefixed columns (like "users.settings")
	t.Run("Prefixed JSON columns", func(t *testing.T) {
		data := map[string]any{
			"users.id":       3,
			"users.name":     "Bob",
			"users.settings": `{"logo": "bob.png", "color": "green", "size": 18}`,
		}

		var user User
		err := xqb.Bind(data, &user)
		assert.NoError(t, err)
		assert.Equal(t, 3, user.ID)
		assert.Equal(t, "Bob", user.Name)
		assert.Equal(t, "bob.png", user.Settings.Logo)
		assert.Equal(t, "green", user.Settings.Color)
		assert.Equal(t, 18, user.Settings.Size)
	})

	// Test Case 4: Empty/null settings
	t.Run("Empty settings", func(t *testing.T) {
		data := map[string]any{
			"id":       4,
			"name":     "Alice",
			"settings": "",
		}

		var user User
		err := xqb.Bind(data, &user)
		assert.NoError(t, err)
		assert.Equal(t, 4, user.ID)
		assert.Equal(t, "Alice", user.Name)
		// Settings should be zero values
		assert.Equal(t, "", user.Settings.Logo)
		assert.Equal(t, "", user.Settings.Color)
		assert.Equal(t, 0, user.Settings.Size)
	})
}
