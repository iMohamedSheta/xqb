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

func Test_Query_WithModelQ(t *testing.T) {
	forEachDialect(t, func(t *testing.T, dialect types.Dialect) {
		sql, bindings, err := xqb.ModelQ(User{}).SetDialect(dialect).
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

func TestBind_SingleModelQ(t *testing.T) {
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

func TestBind_SliceModelQ(t *testing.T) {
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
		Logo  string `xqb:"logo" json:"logo"`
		Color string `xqb:"color" json:"color"`
		Size  int    `xqb:"size" json:"size"`
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

func TestBind_JsonArrayColumns(t *testing.T) {
	type Plan struct {
		ID       int      `xqb:"id"`
		Name     string   `xqb:"name"`
		Features []string `xqb:"features"`
	}

	// Case 1: JSON array as string (from DB jsonb)
	t.Run("JSON array as string", func(t *testing.T) {
		data := map[string]any{
			"id":       1,
			"name":     "Basic",
			"features": `["دعم اساسي", "تحليلات معيارية"]`,
		}

		var plan Plan
		err := xqb.Bind(data, &plan)
		assert.NoError(t, err)
		assert.Equal(t, 1, plan.ID)
		assert.Equal(t, "Basic", plan.Name)
		assert.NotNil(t, plan.Features)
		assert.Len(t, plan.Features, 2)
		assert.Equal(t, "دعم اساسي", plan.Features[0])
		assert.Equal(t, "تحليلات معيارية", plan.Features[1])
	})

	// Case 2: JSON array already decoded (map -> slice)
	t.Run("JSON array as slice", func(t *testing.T) {
		data := map[string]any{
			"id":       2,
			"name":     "Pro",
			"features": []any{"Support", "Analytics"},
		}

		var plan Plan
		err := xqb.Bind(data, &plan)
		assert.NoError(t, err)
		assert.Equal(t, 2, plan.ID)
		assert.Equal(t, "Pro", plan.Name)
		assert.Len(t, plan.Features, 2)
		assert.Equal(t, "Support", plan.Features[0])
		assert.Equal(t, "Analytics", plan.Features[1])
	})

	// Case 3: Empty/null features
	t.Run("Empty features", func(t *testing.T) {
		data := map[string]any{
			"id":       3,
			"name":     "Empty",
			"features": "",
		}

		var plan Plan
		err := xqb.Bind(data, &plan)
		assert.NoError(t, err)
		assert.Equal(t, 3, plan.ID)
		assert.Equal(t, "Empty", plan.Name)
		assert.Empty(t, plan.Features)
	})
}

func TestBind_JsonAdvanced(t *testing.T) {
	type Address struct {
		Street string `json:"street"`
		City   string `json:"city"`
		State  string `json:"state"`
	}

	type Profile struct {
		Username string         `json:"username"`
		Age      int            `json:"age"`
		Active   bool           `json:"active"`
		Address  Address        `json:"address"`
		Hobbies  []string       `json:"hobbies"`
		Tags     []string       `json:"tags"`
		Metadata map[string]any `json:"metadata"`
	}

	type User struct {
		ID      int     `xqb:"id"`
		Name    string  `xqb:"name"`
		Profile Profile `xqb:"profile"`
	}

	// Case 1: JSON as string with nested objects
	t.Run("Nested JSON string", func(t *testing.T) {
		data := map[string]any{
			"id":   1,
			"name": "Ali",
			"profile": `{
				"username": "ali123",
				"age": 30,
				"active": true,
				"address": {"street": "123 Main St", "city": "Cairo", "state": "Cairo Governorate"},
				"hobbies": ["reading","gaming"],
				"tags": ["vip","premium"],
				"metadata": {"key1": "value1","key2": 42}
			}`,
		}

		var user User
		err := xqb.Bind(data, &user)
		assert.NoError(t, err)
		assert.Equal(t, "Ali", user.Name)
		assert.Equal(t, "ali123", user.Profile.Username)
		assert.Equal(t, 30, user.Profile.Age)
		assert.True(t, user.Profile.Active)
		assert.Equal(t, "123 Main St", user.Profile.Address.Street)
		assert.Equal(t, "Cairo", user.Profile.Address.City)
		assert.Equal(t, "Cairo Governorate", user.Profile.Address.State)
		assert.Equal(t, []string{"reading", "gaming"}, user.Profile.Hobbies)
		assert.Equal(t, []string{"vip", "premium"}, user.Profile.Tags)
		assert.Equal(t, map[string]any{"key1": "value1", "key2": float64(42)}, user.Profile.Metadata)
	})

	// Case 2: JSON as map already parsed
	t.Run("Nested JSON map", func(t *testing.T) {
		data := map[string]any{
			"id":   2,
			"name": "Sara",
			"profile": map[string]any{
				"username": "sara321",
				"age":      25,
				"active":   false,
				"address": map[string]any{
					"street": "456 Elm St",
					"city":   "Alexandria",
					"state":  "Alexandria Governorate",
				},
				"hobbies": []any{"swimming", "painting"},
				"tags":    []any{"standard"},
				"metadata": map[string]any{
					"keyA": "valueA",
					"keyB": true,
				},
			},
		}

		var user User
		err := xqb.Bind(data, &user)
		assert.NoError(t, err)
		assert.Equal(t, "Sara", user.Name)
		assert.Equal(t, "sara321", user.Profile.Username)
		assert.Equal(t, 25, user.Profile.Age)
		assert.False(t, user.Profile.Active)
		assert.Equal(t, "456 Elm St", user.Profile.Address.Street)
		assert.Equal(t, "Alexandria", user.Profile.Address.City)
		assert.Equal(t, "Alexandria Governorate", user.Profile.Address.State)
		assert.Equal(t, []string{"swimming", "painting"}, user.Profile.Hobbies)
		assert.Equal(t, []string{"standard"}, user.Profile.Tags)
		assert.Equal(t, map[string]any{"keyA": "valueA", "keyB": true}, user.Profile.Metadata)
	})

	// Case 3: Empty JSON / missing fields
	t.Run("Empty JSON", func(t *testing.T) {
		data := map[string]any{
			"id":      3,
			"name":    "Bob",
			"profile": "",
		}

		var user User
		err := xqb.Bind(data, &user)
		assert.NoError(t, err)
		assert.Equal(t, 3, user.ID)
		assert.Equal(t, "Bob", user.Name)
		assert.Equal(t, "", user.Profile.Username)
		assert.Equal(t, 0, user.Profile.Age)
		assert.False(t, user.Profile.Active)
		assert.Equal(t, Address{}, user.Profile.Address)
		assert.Empty(t, user.Profile.Hobbies)
		assert.Empty(t, user.Profile.Tags)
		assert.Empty(t, user.Profile.Metadata)
	})

	// Case 4: JSON array of objects
	t.Run("JSON array of objects", func(t *testing.T) {
		type Post struct {
			ID    int    `xqb:"id"`
			Title string `xqb:"title"`
		}

		type Blog struct {
			ID    int    `xqb:"id"`
			Name  string `xqb:"name"`
			Posts []Post `xqb:"posts"`
		}

		data := map[string]any{
			"id":   1,
			"name": "Tech Blog",
			"posts": []map[string]any{
				{"id": 101, "title": "First Post"},
				{"id": 102, "title": "Second Post"},
			},
		}

		var blog Blog
		err := xqb.Bind(data, &blog)
		assert.NoError(t, err)
		assert.Len(t, blog.Posts, 2)
		xqb.Dump(blog.Posts)
		assert.Equal(t, 101, blog.Posts[0].ID)
		assert.Equal(t, "First Post", blog.Posts[0].Title)
		assert.Equal(t, 102, blog.Posts[1].ID)
		assert.Equal(t, "Second Post", blog.Posts[1].Title)
	})
}
