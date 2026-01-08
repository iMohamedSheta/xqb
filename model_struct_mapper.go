package xqb

import (
	"encoding/json"
	"reflect"
	"strings"
	"time"
)

// structMapper handles mapping a single map to a struct
type structMapper struct{}

// mapToStruct binds map data to a struct using reflection
func (m *structMapper) mapToStruct(dataMap map[string]any, structValue reflect.Value) error {
	structType := structValue.Type()

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		fieldValue := structValue.Field(i)

		// Handle embedded structs
		if m.isEmbeddedStruct(field, fieldValue) {
			if err := m.mapToStruct(dataMap, fieldValue); err != nil {
				return err
			}
			continue
		}

		if !fieldValue.CanSet() {
			continue
		}

		columnName := getColumnNameForField(field)
		if columnName == "-" {
			continue
		}

		// Initialize pointer fields
		m.initializePointerIfNeeded(fieldValue)

		// Get the actual value to set (dereference if pointer)
		targetValue := m.getTargetValue(fieldValue)

		// Find the data value
		dataValue, exists := m.findDataValue(dataMap, columnName)

		// Handle nested structs
		if m.isNestedStruct(targetValue) {
			if err := m.handleNestedStruct(dataMap, targetValue, dataValue, exists, columnName); err != nil {
				return err
			}
			continue
		}

		// For primitive fields, we need the data value
		if !exists || dataValue == nil {
			continue
		}

		// Set the field value
		valueSetter := &fieldValueSetter{}
		if err := valueSetter.setFieldValue(targetValue, dataValue); err != nil {
			return err
		}
	}

	return nil
}

// isEmbeddedStruct checks if a field is an embedded struct
func (m *structMapper) isEmbeddedStruct(field reflect.StructField, fieldValue reflect.Value) bool {
	return field.Anonymous && fieldValue.Kind() == reflect.Struct && fieldValue.CanSet()
}

// initializePointerIfNeeded initializes a nil pointer field
func (m *structMapper) initializePointerIfNeeded(fieldValue reflect.Value) {
	if fieldValue.Kind() == reflect.Ptr && fieldValue.IsNil() {
		fieldValue.Set(reflect.New(fieldValue.Type().Elem()))
	}
}

// getTargetValue dereferences pointer to get the actual value to work with
func (m *structMapper) getTargetValue(fieldValue reflect.Value) reflect.Value {
	if fieldValue.Kind() == reflect.Ptr {
		return fieldValue.Elem()
	}
	return fieldValue
}

// findDataValue searches for the data value in the map
// Tries direct column name first, then table.column patterns
func (m *structMapper) findDataValue(dataMap map[string]any, columnName string) (any, bool) {
	// Try direct column name first
	if value, exists := dataMap[columnName]; exists {
		return value, true
	}

	// Try table.column pattern
	for key, value := range dataMap {
		if strings.Contains(key, ".") {
			parts := strings.Split(key, ".")
			if len(parts) == 2 && parts[1] == columnName {
				return value, true
			}
		}
	}

	return nil, false
}

// isNestedStruct checks if the target should be treated as a nested struct
func (m *structMapper) isNestedStruct(targetValue reflect.Value) bool {
	if targetValue.Kind() != reflect.Struct {
		return false
	}
	// Exclude SQL null types and time.Time
	if isSQLNullType(targetValue.Type()) {
		return false
	}
	if targetValue.Type() == reflect.TypeOf(time.Time{}) {
		return false
	}
	return true
}

// handleNestedStruct processes nested struct fields
func (m *structMapper) handleNestedStruct(dataMap map[string]any, structValue reflect.Value, dataValue any, exists bool, columnName string) error {
	// Case 1: Data value is already a map
	if exists && dataValue != nil {
		if mapData, ok := dataValue.(map[string]any); ok {
			return m.mapToStruct(mapData, structValue)
		}

		// Case 2: Data value is a JSON string
		if jsonStr, ok := dataValue.(string); ok && jsonStr != "" {
			var mapData map[string]any
			if err := json.Unmarshal([]byte(jsonStr), &mapData); err == nil {
				return m.mapToStruct(mapData, structValue)
			}
		}
	}

	// Case 3: Dot notation binding (e.g., "profile.name" -> Profile.Name)
	dotNotationMapper := &dotNotationMapper{}
	return dotNotationMapper.mapDotNotation(dataMap, structValue, columnName)
}

// dotNotationMapper handles dot notation for nested structs
type dotNotationMapper struct{}

// mapDotNotation extracts fields with dot notation prefix and maps them
// Example: data["profile.name"] -> struct.Profile.Name
func (m *dotNotationMapper) mapDotNotation(dataMap map[string]any, structValue reflect.Value, prefix string) error {
	nestedData := m.extractNestedData(dataMap, prefix)

	if len(nestedData) > 0 {
		mapper := &structMapper{}
		return mapper.mapToStruct(nestedData, structValue)
	}

	return nil
}

// extractNestedData collects all keys with the given prefix
func (m *dotNotationMapper) extractNestedData(dataMap map[string]any, prefix string) map[string]any {
	nestedData := make(map[string]any)
	dotPrefix := prefix + "."

	for key, value := range dataMap {
		if strings.HasPrefix(key, dotPrefix) {
			nestedKey := strings.TrimPrefix(key, dotPrefix)
			nestedData[nestedKey] = value
		}
	}

	return nestedData
}
