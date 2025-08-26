package xqb

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
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
func Model(model ModelInterface) *QueryBuilder {
	return Table(model.Table())
}

// Bind maps data to destination struct or slice
// Examples:
//   - Bind(map[string]any{"name": "John"}, &user)
//   - Bind([]map[string]any{{...}}, &users)
//   - Bind([]map[string]any{{...}}, &user) // aggregates relational data into single struct
func Bind(data any, dest any) error {
	destVal := reflect.ValueOf(dest)
	if destVal.Kind() != reflect.Ptr || destVal.IsNil() {
		return fmt.Errorf("dest must be a non-nil pointer")
	}

	destElem := destVal.Elem()

	switch destElem.Kind() {
	case reflect.Struct:
		if dataMap, ok := data.(map[string]any); ok {
			// Single map to single struct
			return bindStruct(dataMap, destElem)
		} else if dataSlice, ok := data.([]map[string]any); ok {
			// Slice of maps to single struct (aggregate relational data)
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
func bindSlice(dataSlice []map[string]any, destSlice reflect.Value) error {
	elemType := destSlice.Type().Elem()

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

		if !exists || dataValue == nil {
			continue
		}

		// Handle nested struct binding
		if valueToSet.Kind() == reflect.Struct && !isSQLNull(valueToSet.Type()) {
			// Case 1: Data value is already a map
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

			// Case 3: Fall back to dot notation binding
			if err := bindNestedStruct(data, valueToSet, columnName); err != nil {
				return err
			}
			continue
		}

		// Set field value based on type
		if err := setFieldValue(valueToSet, dataValue); err != nil {
			return err
		}
	}

	return nil
}

// getFieldColumnName returns column name from xqb tag or converts field name to snake_case
func getFieldColumnName(field reflect.StructField) string {
	if tag := field.Tag.Get("xqb"); tag != "" {
		return tag
	}
	return toSnakeCase(field.Name)
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

// setFieldValue sets field value based on its type
func setFieldValue(fieldValue reflect.Value, dataValue any) error {
	// Handle slice fields
	if fieldValue.Kind() == reflect.Slice {
		return setSliceFieldValue(fieldValue, dataValue)
	}

	// Handle SQL null types
	if isSQLNull(fieldValue.Type()) {
		setSQLNull(fieldValue, dataValue)
		return nil
	}

	// Handle regular types
	if reflect.TypeOf(dataValue).ConvertibleTo(fieldValue.Type()) {
		fieldValue.Set(reflect.ValueOf(dataValue).Convert(fieldValue.Type()))
	}

	return nil
}

// setSliceFieldValue handles slice field binding
func setSliceFieldValue(sliceValue reflect.Value, dataValue any) error {
	items, ok := dataValue.([]any)
	if !ok {
		// Handle []map[string]any case
		if mapSlice, ok := dataValue.([]map[string]any); ok {
			items = make([]any, len(mapSlice))
			for i, m := range mapSlice {
				items[i] = m
			}
		} else {
			return nil
		}
	}

	sliceType := sliceValue.Type()
	elemType := sliceType.Elem()
	newSlice := reflect.MakeSlice(sliceType, 0, len(items))

	for _, item := range items {
		elem := reflect.New(elemType).Elem()

		// Handle struct elements
		if itemMap, ok := item.(map[string]any); ok && elem.Kind() == reflect.Struct {
			if err := bindStruct(itemMap, elem); err != nil {
				return err
			}
		} else if reflect.TypeOf(item).ConvertibleTo(elemType) {
			// Handle primitive elements
			elem.Set(reflect.ValueOf(item).Convert(elemType))
		}

		newSlice = reflect.Append(newSlice, elem)
	}

	sliceValue.Set(newSlice)
	return nil
}

// isSQLNull checks if the type is sql.Null*
func isSQLNull(t reflect.Type) bool {
	switch t {
	case reflect.TypeOf(sql.NullString{}),
		reflect.TypeOf(sql.NullBool{}),
		reflect.TypeOf(sql.NullInt64{}),
		reflect.TypeOf(sql.NullFloat64{}),
		reflect.TypeOf(sql.NullTime{}):
		return true
	}
	return false
}

// setSQLNull sets the appropriate value for sql.Null* types
func setSQLNull(fv reflect.Value, val any) {
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
		if n, ok := val.(int64); ok {
			fv.Set(reflect.ValueOf(sql.NullInt64{Int64: n, Valid: true}))
		} else if n, ok := val.(int); ok {
			fv.Set(reflect.ValueOf(sql.NullInt64{Int64: int64(n), Valid: true}))
		}
	case reflect.TypeOf(sql.NullFloat64{}):
		if f, ok := val.(float64); ok {
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
	}
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
