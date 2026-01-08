package xqb

import (
	"database/sql"
	"encoding/json"
	"reflect"
	"strconv"
	"strings"
)

// numericConverter handles numeric type conversions
type numericConverter struct{}

func (c *numericConverter) toInt64(value any) (int64, bool) {
	switch v := value.(type) {
	case nil:
		return 0, true
	case int:
		return int64(v), true
	case int8:
		return int64(v), true
	case int16:
		return int64(v), true
	case int32:
		return int64(v), true
	case int64:
		return v, true
	case uint:
		return int64(v), true
	case uint8:
		return int64(v), true
	case uint16:
		return int64(v), true
	case uint32:
		return int64(v), true
	case uint64:
		return int64(v), true
	case float32:
		return int64(v), true
	case float64:
		return int64(v), true
	case json.Number:
		if i, err := v.Int64(); err == nil {
			return i, true
		}
		if f, err := v.Float64(); err == nil {
			return int64(f), true
		}
	case string:
		if i, err := strconv.ParseInt(v, 10, 64); err == nil {
			return i, true
		}
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return int64(f), true
		}
	}
	return 0, false
}

func (c *numericConverter) toUint64(value any) (uint64, bool) {
	switch v := value.(type) {
	case nil:
		return 0, true
	case int:
		return uint64(v), true
	case int8:
		return uint64(v), true
	case int16:
		return uint64(v), true
	case int32:
		return uint64(v), true
	case int64:
		return uint64(v), true
	case uint:
		return uint64(v), true
	case uint8:
		return uint64(v), true
	case uint16:
		return uint64(v), true
	case uint32:
		return uint64(v), true
	case uint64:
		return v, true
	case float32:
		return uint64(v), true
	case float64:
		return uint64(v), true
	case json.Number:
		if i, err := v.Int64(); err == nil {
			return uint64(i), true
		}
	}
	return 0, false
}

func (c *numericConverter) toFloat64(value any) (float64, bool) {
	switch v := value.(type) {
	case nil:
		return 0, true
	case float32:
		return float64(v), true
	case float64:
		return v, true
	case int:
		return float64(v), true
	case int8:
		return float64(v), true
	case int16:
		return float64(v), true
	case int32:
		return float64(v), true
	case int64:
		return float64(v), true
	case uint:
		return float64(v), true
	case uint8:
		return float64(v), true
	case uint16:
		return float64(v), true
	case uint32:
		return float64(v), true
	case uint64:
		return float64(v), true
	case json.Number:
		if f, err := v.Float64(); err == nil {
			return f, true
		}
	case string:
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f, true
		}
	case []byte:
		if f, err := strconv.ParseFloat(string(v), 64); err == nil {
			return f, true
		}
	}
	return 0, false
}

func (c *numericConverter) convertAndSet(element reflect.Value, item any) error {
	switch element.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if n, ok := c.toInt64(item); ok {
			element.SetInt(n)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		if n, ok := c.toUint64(item); ok {
			element.SetUint(n)
		}
	case reflect.Float32, reflect.Float64:
		if f, ok := c.toFloat64(item); ok {
			element.SetFloat(f)
		}
	case reflect.String:
		if s, ok := item.(string); ok {
			element.SetString(s)
		}
	default:
		// Try direct conversion as fallback
		rv := reflect.ValueOf(item)
		if rv.IsValid() && rv.Type().ConvertibleTo(element.Type()) {
			element.Set(rv.Convert(element.Type()))
		}
	}
	return nil
}

// Type checking utilities

func isJSONFieldType(fieldType reflect.Type) bool {
	// []byte for JSON/JSONB
	if fieldType == reflect.TypeOf([]byte{}) {
		return true
	}

	// json.RawMessage
	if fieldType == reflect.TypeOf(json.RawMessage{}) {
		return true
	}

	// map[string]any
	if fieldType.Kind() == reflect.Map && fieldType.Key().Kind() == reflect.String {
		return true
	}

	return false
}

func isSQLNullType(fieldType reflect.Type) bool {
	switch fieldType {
	case reflect.TypeOf(sql.NullString{}),
		reflect.TypeOf(sql.NullBool{}),
		reflect.TypeOf(sql.NullInt64{}),
		reflect.TypeOf(sql.NullInt32{}),
		reflect.TypeOf(sql.NullInt16{}),
		reflect.TypeOf(sql.NullByte{}),
		reflect.TypeOf(sql.NullFloat64{}),
		reflect.TypeOf(sql.NullTime{}):
		return true
	}
	return false
}

// Field naming utilities

func getColumnNameForField(field reflect.StructField) string {
	// Check for explicit xqb tag
	if tag := field.Tag.Get("xqb"); tag != "" {
		return tag
	}

	// Convert field name to snake_case
	return convertToSnakeCase(field.Name)
}

func convertToSnakeCase(camelCase string) string {
	var result []rune

	for i, char := range camelCase {
		if i > 0 && char >= 'A' && char <= 'Z' {
			result = append(result, '_')
		}
		result = append(result, char)
	}

	return strings.ToLower(string(result))
}
