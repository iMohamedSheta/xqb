package xqb_test

import (
	"database/sql"
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

	t.Run("SQL Nulls extended", func(t *testing.T) {
		require.True(t, m.NullInt16.Valid)
		require.Equal(t, int16(21), m.NullInt16.Int16)

		require.True(t, m.NullInt32.Valid)
		require.Equal(t, int32(22), m.NullInt32.Int32)

		require.True(t, m.NullByte.Valid)
		require.Equal(t, byte(23), m.NullByte.Byte)
	})

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
}
