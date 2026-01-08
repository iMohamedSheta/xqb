package xqb_test

import (
	"database/sql"
	"encoding/json"
	"testing"
	"time"

	"github.com/iMohamedSheta/xqb"
	"github.com/stretchr/testify/require"
)

type AllTypes struct {
	// Primitives
	IntVal     int     `xqb:"int_val"`
	Int8Val    int8    `xqb:"int8_val"`
	Int16Val   int16   `xqb:"int16_val"`
	Int32Val   int32   `xqb:"int32_val"`
	Int64Val   int64   `xqb:"int64_val"`
	UintVal    uint    `xqb:"uint_val"`
	Uint8Val   uint8   `xqb:"uint8_val"`
	Uint16Val  uint16  `xqb:"uint16_val"`
	Uint32Val  uint32  `xqb:"uint32_val"`
	Uint64Val  uint64  `xqb:"uint64_val"`
	Float32Val float32 `xqb:"float32_val"`
	Float64Val float64 `xqb:"float64_val"`
	BoolVal    bool    `xqb:"bool_val"`
	StringVal  string  `xqb:"string_val"`

	// Pointers [Not Supported] if you want nullable type use default sql package nullable types sql.NullType
	// PtrInt    *int       `xqb:"ptr_int"`
	// PtrBool   *bool      `xqb:"ptr_bool"`
	// PtrString *string    `xqb:"ptr_string"`
	// PtrTime   *time.Time `xqb:"ptr_time"`

	// SQL Nulls
	NullStr   sql.NullString  `xqb:"null_str"`
	NullInt16 sql.NullInt16   `xqb:"null_int16"`
	NullInt32 sql.NullInt32   `xqb:"null_int32"`
	NullByte  sql.NullByte    `xqb:"null_byte"`
	NullInt   sql.NullInt64   `xqb:"null_int"`
	NullFloat sql.NullFloat64 `xqb:"null_float"`
	NullBool  sql.NullBool    `xqb:"null_bool"`
	NullTime  sql.NullTime    `xqb:"null_time"`

	// Time
	TimeVal time.Time `xqb:"time_val"`

	// JSON/JSONB fields
	RawJson        []byte          `xqb:"raw_json"`         // For JSONB/JSON as bytes
	RawJsonMessage json.RawMessage `xqb:"raw_json_message"` // Alternative JSON storage
	JsonMap        map[string]any  `xqb:"json_map"`         // Will unmarshal JSON into map
	JsonbData      []byte          `xqb:"jsonb_data"`       // PostgreSQL JSONB
	MetaData       map[string]any  `xqb:"metadata"`         // Common pattern for metadata

	// Collections
	SliceStr []string       `xqb:"slice_str"`
	SliceInt []int          `xqb:"slice_int"`
	ArrayInt [3]int         `xqb:"array_int"`
	MapAny   map[string]any `xqb:"map_any"`

	// Custom
	CustomString MyString `xqb:"custom_string"`
	CustomInt    MyInt    `xqb:"custom_int"`

	// Struct embedding + nested
	Embedded
	Nested NestedStruct `xqb:"nested"`

	// Ignored / unsupported
	Iface any       `xqb:"-"`
	Func  func()    `xqb:"-"`
	Chan  chan bool `xqb:"-"`
}

type Embedded struct {
	EmbeddedVal string `xqb:"embedded_val"`
}

type NestedStruct struct {
	NestedVal string `xqb:"nested_val"`
}

type MyString string
type MyInt int

func TestBind_AllTypes(t *testing.T) {
	now := time.Now()

	data := map[string]any{
		// Primitives
		"int_val":     1,
		"int8_val":    int8(2),
		"int16_val":   int16(3),
		"int32_val":   int32(4),
		"int64_val":   int64(5),
		"uint_val":    uint(6),
		"uint8_val":   uint8(7),
		"uint16_val":  uint16(8),
		"uint32_val":  uint32(9),
		"uint64_val":  uint64(10),
		"float32_val": float32(11.11),
		"float64_val": 12.34,
		"bool_val":    true,
		"string_val":  "hello-world",

		// SQL Nulls
		"null_str":   "nullable",
		"null_int16": int16(21),
		"null_int32": int32(22),
		"null_byte":  byte(23),
		"null_int":   123,
		"null_float": 45.67,
		"null_bool":  true,
		"null_time":  now,

		// Time
		"time_val": now,

		// Collections
		"slice_str": []any{"a", "b"},
		"slice_int": []any{1, 2},
		"array_int": [3]int{1, 2, 3},
		"map_any":   map[string]any{"x": 1, "y": "yes"},

		// Custom
		"custom_string": "custom-val",
		"custom_int":    77,

		// Embedded + nested
		"embedded_val":      "embed-me",
		"nested.nested_val": "nested-me",
	}

	var m AllTypes
	err := xqb.Bind(data, &m)
	require.NoError(t, err)

	// Check Primitive types
	t.Run("Primitive types", func(t *testing.T) {
		require.Equal(t, 1, m.IntVal)
		require.Equal(t, int8(2), m.Int8Val)
		require.Equal(t, int16(3), m.Int16Val)
		require.Equal(t, int32(4), m.Int32Val)
		require.Equal(t, int64(5), m.Int64Val)
		require.Equal(t, uint(6), m.UintVal)
		require.Equal(t, uint8(7), m.Uint8Val)
		require.Equal(t, uint16(8), m.Uint16Val)
		require.Equal(t, uint32(9), m.Uint32Val)
		require.Equal(t, uint64(10), m.Uint64Val)
		require.InEpsilon(t, float32(11.11), m.Float32Val, 0.0001)
		require.InEpsilon(t, 12.34, m.Float64Val, 0.0001)
		require.Equal(t, true, m.BoolVal)
		require.Equal(t, "hello-world", m.StringVal)
	})

	// Check SQL Nulls
	t.Run("SQL Nulls extended", func(t *testing.T) {
		require.True(t, m.NullInt16.Valid)
		require.Equal(t, int16(21), m.NullInt16.Int16)

		require.True(t, m.NullInt32.Valid)
		require.Equal(t, int32(22), m.NullInt32.Int32)

		require.True(t, m.NullByte.Valid)
		require.Equal(t, byte(23), m.NullByte.Byte)

		require.True(t, m.NullStr.Valid)
		require.Equal(t, "nullable", m.NullStr.String)

		require.True(t, m.NullInt.Valid)
		require.Equal(t, int64(123), m.NullInt.Int64)

		require.True(t, m.NullFloat.Valid)
		require.InEpsilon(t, 45.67, m.NullFloat.Float64, 0.0001)

		require.True(t, m.NullBool.Valid)
		require.Equal(t, true, m.NullBool.Bool)

		require.True(t, m.NullTime.Valid)
		require.WithinDuration(t, now, m.NullTime.Time, time.Second)
	})

	// Check Zero/null values
	t.Run("Zero/null values", func(t *testing.T) {
		data2 := map[string]any{
			"null_str":   nil,
			"null_int":   nil,
			"null_float": nil,
			"null_bool":  nil,
			"null_time":  nil,
		}
		var m2 AllTypes
		err := xqb.Bind(data2, &m2)
		require.NoError(t, err)
		require.False(t, m2.NullStr.Valid)
		require.False(t, m2.NullInt.Valid)
		require.False(t, m2.NullFloat.Valid)
		require.False(t, m2.NullBool.Valid)
		require.False(t, m2.NullTime.Valid)
	})

	// Check Time
	t.Run("Time", func(t *testing.T) {
		require.WithinDuration(t, now, m.TimeVal, time.Second)
	})

	// Check Collections
	t.Run("Collections", func(t *testing.T) {
		require.Equal(t, []string{"a", "b"}, m.SliceStr)
		require.Equal(t, []int{1, 2}, m.SliceInt)
		require.Equal(t, [3]int{1, 2, 3}, m.ArrayInt)
		require.Equal(t, map[string]interface{}{"x": 1.0, "y": "yes"}, m.MapAny)
	})

	// Check Custom types
	t.Run("Custom types", func(t *testing.T) {
		require.Equal(t, MyString("custom-val"), m.CustomString)
		require.Equal(t, MyInt(77), m.CustomInt)
	})

	// Check Embedded and Nested
	t.Run("Embedded and Nested", func(t *testing.T) {
		require.Equal(t, "embed-me", m.EmbeddedVal)
		require.Equal(t, "nested-me", m.Nested.NestedVal)
	})

	// Check String to Float conversion
	t.Run("String to Float conversion", func(t *testing.T) {
		data3 := map[string]any{
			"float32_val": "13.37",
			"float64_val": []byte("42.42"),
		}

		var m3 AllTypes
		err := xqb.Bind(data3, &m3)
		require.NoError(t, err)

		require.InEpsilon(t, float32(13.37), m3.Float32Val, 0.0001)
		require.InEpsilon(t, float64(42.42), m3.Float64Val, 0.0001)
	})

}
func TestBind_JSON_JSONB(t *testing.T) {
	t.Run("JSON from bytes", func(t *testing.T) {
		data := map[string]any{
			"raw_json": []byte(`{"name":"John","age":30}`),
		}

		var m AllTypes
		err := xqb.Bind(data, &m)
		require.NoError(t, err)
		require.JSONEq(t, `{"name":"John","age":30}`, string(m.RawJson))
	})

	t.Run("JSON from string", func(t *testing.T) {
		data := map[string]any{
			"raw_json": `{"city":"NYC","country":"USA"}`,
		}

		var m AllTypes
		err := xqb.Bind(data, &m)
		require.NoError(t, err)
		require.JSONEq(t, `{"city":"NYC","country":"USA"}`, string(m.RawJson))
	})

	t.Run("JSON to map", func(t *testing.T) {
		data := map[string]any{
			"json_map": []byte(`{"key1":"value1","key2":"value2","count":42}`),
		}

		var m AllTypes
		err := xqb.Bind(data, &m)
		require.NoError(t, err)
		require.Equal(t, "value1", m.JsonMap["key1"])
		require.Equal(t, "value2", m.JsonMap["key2"])
		require.Equal(t, float64(42), m.JsonMap["count"])
	})

	t.Run("Map to JSON bytes", func(t *testing.T) {
		data := map[string]any{
			"metadata": map[string]any{
				"version": "1.0",
				"active":  true,
				"tags":    []string{"go", "database"},
			},
		}

		var m AllTypes
		err := xqb.Bind(data, &m)
		require.NoError(t, err)
		require.Equal(t, "1.0", m.MetaData["version"])
		require.Equal(t, true, m.MetaData["active"])
	})

	t.Run("JSONB from PostgreSQL", func(t *testing.T) {
		// Simulate JSONB data from PostgreSQL (typically comes as []byte)
		data := map[string]any{
			"jsonb_data": []byte(`{"user_id":123,"preferences":{"theme":"dark","lang":"en"}}`),
		}

		var m AllTypes
		err := xqb.Bind(data, &m)
		require.NoError(t, err)

		var parsed map[string]any
		err = json.Unmarshal(m.JsonbData, &parsed)
		require.NoError(t, err)
		require.Equal(t, float64(123), parsed["user_id"])
	})

	t.Run("Invalid JSON handling", func(t *testing.T) {
		data := map[string]any{
			"raw_json": []byte(`{invalid json}`),
		}

		var m AllTypes
		err := xqb.Bind(data, &m)
		require.Error(t, err)
		require.Contains(t, err.Error(), "invalid JSON")
	})

	t.Run("Empty JSON", func(t *testing.T) {
		data := map[string]any{
			"raw_json": []byte{},
			"json_map": "",
		}

		var m AllTypes
		err := xqb.Bind(data, &m)
		require.NoError(t, err)
		require.Empty(t, m.RawJson)
		require.Nil(t, m.JsonMap)
	})

	t.Run("JSON array", func(t *testing.T) {
		data := map[string]any{
			"raw_json_message": []byte(`["item1","item2","item3"]`),
		}

		var m AllTypes
		err := xqb.Bind(data, &m)
		require.NoError(t, err)

		var arr []string
		err = json.Unmarshal(m.RawJsonMessage, &arr)
		require.NoError(t, err)
		require.Equal(t, []string{"item1", "item2", "item3"}, arr)
	})

	t.Run("Nested JSON objects", func(t *testing.T) {
		data := map[string]any{
			"json_map": `{
				"user": {
					"name": "Alice",
					"details": {
						"age": 25,
						"city": "Boston"
					}
				}
			}`,
		}

		var m AllTypes
		err := xqb.Bind(data, &m)
		require.NoError(t, err)

		user := m.JsonMap["user"].(map[string]any)
		require.Equal(t, "Alice", user["name"])

		details := user["details"].(map[string]any)
		require.Equal(t, float64(25), details["age"])
		require.Equal(t, "Boston", details["city"])
	})

	t.Run("JSON RawMessage from string", func(t *testing.T) {
		data := map[string]any{
			"raw_json_message": `{"status":"active","id":999}`,
		}

		var m AllTypes
		err := xqb.Bind(data, &m)
		require.NoError(t, err)

		var obj map[string]any
		err = json.Unmarshal(m.RawJsonMessage, &obj)
		require.NoError(t, err)
		require.Equal(t, "active", obj["status"])
		require.Equal(t, float64(999), obj["id"])
	})

	t.Run("Nil JSON values", func(t *testing.T) {
		data := map[string]any{
			"raw_json":         nil,
			"raw_json_message": nil,
			"json_map":         nil,
		}

		var m AllTypes
		err := xqb.Bind(data, &m)
		require.NoError(t, err)
		require.Nil(t, m.RawJson)
		require.Nil(t, m.RawJsonMessage)
		require.Nil(t, m.JsonMap)
	})

	t.Run("Complex nested structure", func(t *testing.T) {
		data := map[string]any{
			"metadata": map[string]any{
				"config": map[string]any{
					"database": map[string]any{
						"host":     "localhost",
						"port":     5432,
						"ssl":      true,
						"replicas": []any{"db1", "db2", "db3"},
					},
					"cache": map[string]any{
						"ttl":     3600,
						"enabled": true,
					},
				},
			},
		}

		var m AllTypes
		err := xqb.Bind(data, &m)
		require.NoError(t, err)

		config := m.MetaData["config"].(map[string]any)
		database := config["database"].(map[string]any)

		require.Equal(t, "localhost", database["host"])
		require.Equal(t, float64(5432), database["port"])
		require.Equal(t, true, database["ssl"])

		replicas := database["replicas"].([]any)
		require.Len(t, replicas, 3)
		require.Equal(t, "db1", replicas[0])
	})
}
