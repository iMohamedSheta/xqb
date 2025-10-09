package xqb

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

/*
| ----------------------------------------------
| Model support try by xqb
| ----------------------------------------------
*/

// Model Interface represents a model
type ModelInterface interface {
	// Table returns the table name of the model
	Table() string
}

// Model returns the model of the query builder
func ModelQ(model ModelInterface) *QueryBuilder {
	return Table(model.Table())
}

// Bind maps data to destination struct of models or slice of models
// Examples:
//   - Bind(map[string]any{"name": "John"}, &user)
//   - Bind([]map[string]any{{...}}, &users)
//   - Bind([]map[string]any{{...}}, &user) // aggregates relational data into single struct
func Bind(data any, dest any) error {
	destVal := reflect.ValueOf(dest)

	// dest must be a non-nil pointer reference to a model (&model || &[]model)
	if destVal.Kind() != reflect.Ptr || destVal.IsNil() {
		return fmt.Errorf("dest must be a non-nil pointer")
	}

	destElem := destVal.Elem()

	// the destination model should be struct or slice of structs example: &User{} || &[]User{} anything other than that is not valid
	switch destElem.Kind() {
	case reflect.Struct:
		if dataMap, ok := data.(map[string]any); ok {
			// Single map to single struct
			return bindStruct(dataMap, destElem)
		} else if dataSlice, ok := data.([]map[string]any); ok {
			// Slice of maps to single struct (aggregate relational data)
			// some data can be []map[string]any not always map[string]any for one struct model
			// cause it handle joins data into list of maps and repeat a lot of data
			// so we handle it by binding first data map to struct and then search of values that
			// expect many values and bind all the values to it
			return bindRelationalDataToStruct(dataSlice, destElem)
		}
		return fmt.Errorf("data must be map[string]any or []map[string]any for struct binding")

	case reflect.Slice:
		// Slice binding
		dataSlice, ok := data.([]map[string]any)
		if !ok {
			return fmt.Errorf("data must be []map[string]any for slice binding")
		}
		return bindSlice(dataSlice, destElem)

	default:
		return fmt.Errorf("unsupported destination type: %s", destElem.Kind())
	}
}

// bindRelationalDataToStruct aggregates flat relational data into a single struct with nested slices
// Example: []map[string]any{{"id": 1, "posts_title": "Post1"}, {"id": 1, "posts_title": "Post2"}}
// -> User{ID: 1, Posts: []Post{{Title: "Post1"}, {Title: "Post2"}}}
func bindRelationalDataToStruct(dataSlice []map[string]any, structValue reflect.Value) error {
	if len(dataSlice) == 0 {
		return nil
	}

	// First, bind the main struct fields from the first row
	if err := bindStruct(dataSlice[0], structValue); err != nil {
		return err
	}

	// Then, handle slice fields by aggregating data from all rows
	structType := structValue.Type()
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		fieldValue := structValue.Field(i)

		if !fieldValue.CanSet() || fieldValue.Kind() != reflect.Slice {
			continue
		}

		columnName := getFieldColumnName(field)
		if columnName == "-" {
			continue
		}

		// Then, handle slice fields by aggregating data from all rows
		structType := structValue.Type()
		for i := 0; i < structType.NumField(); i++ {
			field := structType.Field(i)
			fieldValue := structValue.Field(i)

			if !fieldValue.CanSet() || fieldValue.Kind() != reflect.Slice {
				continue
			}

			columnName := getFieldColumnName(field)
			if columnName == "-" {
				continue
			}

			// Get table name from field tag, fallback to column name
			tableName := field.Tag.Get("table")
			if tableName == "" {
				tableName = columnName
			}

			// Aggregate slice data from all rows
			if err := aggregateSliceFromRelationalData(dataSlice, fieldValue, tableName); err != nil {
				return err
			}
		}
	}
	return nil
}

// aggregateSliceFromRelationalData collects slice elements from flat relational data
// Uses table name to identify related columns and validates against struct fields
// Example: Posts []Post `xqb:"posts" table:"posts"` -> matches posts_title, posts_id but skips invalid fields
func aggregateSliceFromRelationalData(dataSlice []map[string]any, sliceField reflect.Value, tableName string) error {
	sliceType := sliceField.Type()
	elemType := sliceType.Elem()
	newSlice := reflect.MakeSlice(sliceType, 0, len(dataSlice))

	// Get valid fields from the target struct type
	validFields := getStructFieldColumns(elemType)
	exactPrefix := tableName + "_"

	for _, rowData := range dataSlice {
		elem := reflect.New(elemType).Elem()
		elemData := make(map[string]any)

		for key, value := range rowData {
			// Match columns that start with table prefix
			if strings.HasPrefix(key, exactPrefix) && len(key) > len(exactPrefix) {
				elemKey := strings.TrimPrefix(key, exactPrefix)

				// Only include if this field exists in the target struct
				if _, exists := validFields[elemKey]; exists {
					elemData[elemKey] = value
				}
			}
		}

		// Only add if we found related data
		if len(elemData) > 0 {
			if err := bindStruct(elemData, elem); err != nil {
				return err
			}
			newSlice = reflect.Append(newSlice, elem)
		}
	}

	sliceField.Set(newSlice)
	return nil
}

// getStructFieldColumns returns a map of column names to field names for a struct type
// This helps validate which columns actually belong to the struct
func getStructFieldColumns(structType reflect.Type) map[string]string {
	columns := make(map[string]string)

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		columnName := getFieldColumnName(field)
		if columnName != "-" {
			columns[columnName] = field.Name
		}
	}

	return columns
}

// bindSlice is a way to bind a slice of structs to a slice of structs
// it uses the same logic as bindStruct but loops through the slice of structs
func bindSlice(dataSlice []map[string]any, destSlice reflect.Value) error {
	elemType := destSlice.Type().Elem() // Example User as type for []User list of users models

	for _, itemMap := range dataSlice {
		newElem := reflect.New(elemType).Elem()
		if err := bindStruct(itemMap, newElem); err != nil {
			return err
		}
		destSlice.Set(reflect.Append(destSlice, newElem))
	}
	return nil
}

// bindStruct binds map data to struct using field tags and reflection
func bindStruct(data map[string]any, v reflect.Value) error {
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)

		// If field is anonymous (embedded), recurse into it as if it's part of parent
		if field.Anonymous && fieldValue.Kind() == reflect.Struct && fieldValue.CanSet() {
			// Recurse into embedded struct fields using same data map
			if err := bindStruct(data, fieldValue); err != nil {
				return err
			}
			continue
		}

		if !fieldValue.CanSet() {
			continue
		}

		// Get column name from tag or convert field name
		columnName := getFieldColumnName(field)
		if columnName == "-" {
			continue
		}

		// Handle pointer fields - initialize if nil
		if fieldValue.Kind() == reflect.Ptr && fieldValue.IsNil() {
			fieldValue.Set(reflect.New(fieldValue.Type().Elem()))
		}

		// Get the actual value to work with
		valueToSet := fieldValue
		if fieldValue.Kind() == reflect.Ptr {
			valueToSet = fieldValue.Elem()
		}

		// Look for data value (try direct column name and prefixed versions)
		var dataValue any
		var exists bool

		// Try direct column name first
		if dataValue, exists = data[columnName]; !exists {
			// Try with potential table prefixes
			for key, value := range data {
				// Match "table.column" pattern where column matches our field
				if strings.Contains(key, ".") {
					parts := strings.Split(key, ".")
					if len(parts) == 2 && parts[1] == columnName {
						dataValue = value
						exists = true
						break
					}
				}
			}
		}

		// Handle nested struct binding ONLY if it's a struct and we don't have direct data
		if valueToSet.Kind() == reflect.Struct && !isSQLNull(valueToSet.Type()) && valueToSet.Type() != reflect.TypeOf(time.Time{}) {
			if exists && dataValue != nil {
				if mapData, ok := dataValue.(map[string]any); ok {
					if err := bindStruct(mapData, valueToSet); err != nil {
						return err
					}
					continue
				}

				// Case 2: Data value is a JSON string
				if jsonStr, ok := dataValue.(string); ok && jsonStr != "" {
					var mapData map[string]any
					if err := json.Unmarshal([]byte(jsonStr), &mapData); err == nil {
						if err := bindStruct(mapData, valueToSet); err != nil {
							return err
						}
						continue
					}
				}
			}

			// Case 3: Fall back to dot notation binding (whether we have direct data or not)
			if err := bindNestedStruct(data, valueToSet, columnName); err != nil {
				return err
			}
			continue
		}

		// For non-struct fields, we need the direct data value
		if !exists || dataValue == nil {
			continue
		}

		// Set field value based on type
		if err := setFieldValue(valueToSet, dataValue); err != nil {
			return err
		}
	}

	return nil
}

// bindNestedStruct handles nested struct fields with dot notation
// Example: data["profile.name"] -> struct.Profile.Name
func bindNestedStruct(data map[string]any, structValue reflect.Value, prefix string) error {
	nestedMap := make(map[string]any)
	dotPrefix := prefix + "."

	for key, value := range data {
		if strings.HasPrefix(key, dotPrefix) {
			nestedKey := strings.TrimPrefix(key, dotPrefix)
			nestedMap[nestedKey] = value
		}
	}

	if len(nestedMap) > 0 {
		return bindStruct(nestedMap, structValue)
	}
	return nil
}

// setFieldValue handles field value based on type
func setFieldValue(fieldValue reflect.Value, dataValue any) error {
	// If field is a slice, delegate
	if fieldValue.Kind() == reflect.Slice {
		return setSliceFieldValue(fieldValue, dataValue)
	}

	// If the field implements sql.Scanner it was handled in bindStruct;
	// but as a fallback, check again for non-addressable scanner (rare)
	scannerType := reflect.TypeOf((*sql.Scanner)(nil)).Elem()
	if fieldValue.Type().Implements(scannerType) {
		// non-pointer scanner (uncommon) - create and scan
		scanner := reflect.New(fieldValue.Type()).Interface().(sql.Scanner)
		if err := scanner.Scan(dataValue); err != nil {
			return err
		}
		fieldValue.Set(reflect.ValueOf(scanner).Elem())
		return nil
	}

	// Handle sql.Null... concrete types (legacy path)
	if isSQLNull(fieldValue.Type()) {
		setSQLNull(fieldValue, dataValue)
		return nil
	}

	// time.Time handling: accept time.Time or parse RFC strings if provided
	if fieldValue.Type() == reflect.TypeOf(time.Time{}) {
		if timeVal, ok := dataValue.(time.Time); ok {
			fieldValue.Set(reflect.ValueOf(timeVal))
			return nil
		}
		if s, ok := dataValue.(string); ok && s != "" {
			// Try RFC3339 parse (best-effort). You can expand formats if needed.
			if tm, err := time.Parse(time.RFC3339, s); err == nil {
				fieldValue.Set(reflect.ValueOf(tm))
				return nil
			}
		}
	}

	// Normal convertible types (handle numeric widening)
	dv := reflect.ValueOf(dataValue)
	if !dv.IsValid() {
		return nil
	}

	// If direct convertible, convert and set
	if dv.Type().ConvertibleTo(fieldValue.Type()) {
		fieldValue.Set(dv.Convert(fieldValue.Type()))
		return nil
	}

	// Attempt common numeric conversions (e.g., int -> int64)
	switch fieldValue.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if num, ok := toInt64(dataValue); ok {
			fieldValue.SetInt(num)
			return nil
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		if unum, ok := toUint64(dataValue); ok {
			fieldValue.SetUint(unum)
			return nil
		}
	case reflect.Float32, reflect.Float64:
		if f, ok := toFloat64(dataValue); ok {
			fieldValue.SetFloat(f)
			return nil
		}
	case reflect.Bool:
		if b, ok := dataValue.(bool); ok {
			fieldValue.SetBool(b)
			return nil
		}
	case reflect.String:
		if s, ok := dataValue.(string); ok {
			fieldValue.SetString(s)
			return nil
		}
	}

	// not convertible - ignore silently
	return nil
}

// ---------- setSliceFieldValue (updated) ----------
func setSliceFieldValue(sliceValue reflect.Value, dataValue any) error {
	// Accept:
	// - []any
	// - []string, []int, etc
	// - []map[string]any (for struct elements)
	// - single value (wrap into one-element slice) - optional
	// Special case: JSON string/bytes

	switch v := dataValue.(type) {
	case string:
		if v != "" {
			var decoded any
			if err := json.Unmarshal([]byte(v), &decoded); err == nil {
				return setSliceFieldValue(sliceValue, decoded)
			}
		}
	case []byte:
		if len(v) > 0 {
			var decoded any
			if err := json.Unmarshal(v, &decoded); err == nil {
				return setSliceFieldValue(sliceValue, decoded)
			}
		}
	}

	v := reflect.ValueOf(dataValue)
	if !v.IsValid() {
		return nil
	}

	// If it's not a slice/array, nothing to do
	kind := v.Kind()
	if kind != reflect.Slice && kind != reflect.Array {
		return nil
	}

	elemType := sliceValue.Type().Elem()
	newSlice := reflect.MakeSlice(sliceValue.Type(), 0, v.Len())

	for i := 0; i < v.Len(); i++ {
		item := v.Index(i).Interface()
		elem := reflect.New(elemType).Elem()

		// If element is a struct and item is map -> bindStruct
		if elem.Kind() == reflect.Struct {
			if m, ok := item.(map[string]any); ok {
				if err := bindStruct(m, elem); err != nil {
					return err
				}
			} else {
				// If item is JSON string, try to unmarshal into map and bind
				if s, ok := item.(string); ok && s != "" {
					var m map[string]any
					if err := json.Unmarshal([]byte(s), &m); err == nil {
						if err := bindStruct(m, elem); err != nil {
							return err
						}
					}
				}
			}
		} else {
			// Primitive/convertible element - try direct conversion
			rv := reflect.ValueOf(item)
			if rv.IsValid() && rv.Type().ConvertibleTo(elemType) {
				elem.Set(rv.Convert(elemType))
			} else {
				switch elem.Kind() {
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					if n, ok := toInt64(item); ok {
						elem.SetInt(n)
					}
				case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
					if n, ok := toUint64(item); ok {
						elem.SetUint(n)
					}
				case reflect.Float32, reflect.Float64:
					if f, ok := toFloat64(item); ok {
						elem.SetFloat(f)
					}
				case reflect.String:
					if s, ok := item.(string); ok {
						elem.SetString(s)
					}
				default:
					rv := reflect.ValueOf(item)
					if rv.IsValid() && rv.Type().ConvertibleTo(elemType) {
						elem.Set(rv.Convert(elemType))
					}
				}

			}
		}

		newSlice = reflect.Append(newSlice, elem)
	}

	sliceValue.Set(newSlice)
	return nil
}

// ---------- isSQLNull (updated to include other Null types) ----------
func isSQLNull(t reflect.Type) bool {
	switch t {
	case reflect.TypeOf(sql.NullString{}),
		reflect.TypeOf(sql.NullBool{}),
		reflect.TypeOf(sql.NullInt64{}),
		reflect.TypeOf(sql.NullFloat64{}),
		reflect.TypeOf(sql.NullTime{}),
		reflect.TypeOf(sql.NullInt16{}),
		reflect.TypeOf(sql.NullInt32{}),
		reflect.TypeOf(sql.NullByte{}):
		return true
	}
	return false
}

// ---------- setSQLNull (improved) ----------
func setSQLNull(fv reflect.Value, val any) {
	// Normalize numbers and strings to appropriate underlying type
	switch fv.Type() {
	case reflect.TypeOf(sql.NullString{}):
		if s, ok := val.(string); ok {
			fv.Set(reflect.ValueOf(sql.NullString{String: s, Valid: true}))
		}
	case reflect.TypeOf(sql.NullBool{}):
		if b, ok := val.(bool); ok {
			fv.Set(reflect.ValueOf(sql.NullBool{Bool: b, Valid: true}))
		}
	case reflect.TypeOf(sql.NullInt64{}):
		if n, ok := toInt64(val); ok {
			fv.Set(reflect.ValueOf(sql.NullInt64{Int64: n, Valid: true}))
		}
	case reflect.TypeOf(sql.NullFloat64{}):
		if f, ok := toFloat64(val); ok {
			fv.Set(reflect.ValueOf(sql.NullFloat64{Float64: f, Valid: true}))
		}
	case reflect.TypeOf(sql.NullTime{}):
		if tm, ok := val.(time.Time); ok {
			var valid bool
			if !tm.IsZero() {
				valid = true
			}
			fv.Set(reflect.ValueOf(sql.NullTime{Time: tm, Valid: valid}))
		}
	// Additional ones:
	case reflect.TypeOf(sql.NullInt16{}):
		if n, ok := toInt64(val); ok {
			fv.Set(reflect.ValueOf(sql.NullInt16{Int16: int16(n), Valid: true}))
		}
	case reflect.TypeOf(sql.NullInt32{}):
		if n, ok := toInt64(val); ok {
			fv.Set(reflect.ValueOf(sql.NullInt32{Int32: int32(n), Valid: true}))
		}
	case reflect.TypeOf(sql.NullByte{}):
		if n, ok := toInt64(val); ok {
			fv.Set(reflect.ValueOf(sql.NullByte{Byte: byte(n), Valid: true}))
		}
	}
}

// getFieldColumnName returns column name from xqb tag or converts field name to snake_case
func getFieldColumnName(field reflect.StructField) string {
	if tag := field.Tag.Get("xqb"); tag != "" {
		return tag
	}
	return toSnakeCase(field.Name)
}

// toSnakeCase converts CamelCase to snake_case
func toSnakeCase(str string) string {
	var result []rune
	for i, r := range str {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result = append(result, '_')
		}
		result = append(result, r)
	}
	return strings.ToLower(string(result))
}

func toInt64(v any) (int64, bool) {
	switch n := v.(type) {
	case int:
		return int64(n), true
	case int8:
		return int64(n), true
	case int16:
		return int64(n), true
	case int32:
		return int64(n), true
	case int64:
		return n, true
	case uint:
		return int64(n), true
	case uint8:
		return int64(n), true
	case uint16:
		return int64(n), true
	case uint32:
		return int64(n), true
	case uint64:
		return int64(n), true
	case float32:
		return int64(n), true
	case float64:
		return int64(n), true
	case json.Number:
		if i, err := n.Int64(); err == nil {
			return i, true
		}
		if f, err := n.Float64(); err == nil {
			return int64(f), true
		}
	case string:
		if i, err := strconv.ParseInt(n, 10, 64); err == nil {
			return i, true
		}
		if f, err := strconv.ParseFloat(n, 64); err == nil {
			return int64(f), true
		}
	}
	return 0, false
}

func toUint64(v any) (uint64, bool) {
	switch n := v.(type) {
	case int:
		return uint64(n), true
	case int8:
		return uint64(n), true
	case int16:
		return uint64(n), true
	case int32:
		return uint64(n), true
	case int64:
		return uint64(n), true
	case uint:
		return uint64(n), true
	case uint8:
		return uint64(n), true
	case uint16:
		return uint64(n), true
	case uint32:
		return uint64(n), true
	case uint64:
		return n, true
	case float32:
		return uint64(n), true
	case float64:
		return uint64(n), true
	case json.Number:
		if i, err := n.Int64(); err == nil {
			return uint64(i), true
		}
	}
	return 0, false
}

func toFloat64(v any) (float64, bool) {
	switch n := v.(type) {
	case float32:
		return float64(n), true
	case float64:
		return n, true
	case int:
		return float64(n), true
	case int8:
		return float64(n), true
	case int16:
		return float64(n), true
	case int32:
		return float64(n), true
	case int64:
		return float64(n), true
	case uint:
		return float64(n), true
	case uint8:
		return float64(n), true
	case uint16:
		return float64(n), true
	case uint32:
		return float64(n), true
	case uint64:
		return float64(n), true
	case json.Number:
		if f, err := n.Float64(); err == nil {
			return f, true
		}
	case string:
		if f, err := strconv.ParseFloat(n, 64); err == nil {
			return f, true
		}
	case []byte:
		if f, err := strconv.ParseFloat(string(n), 64); err == nil {
			return f, true
		}
	}
	return 0, false
}
