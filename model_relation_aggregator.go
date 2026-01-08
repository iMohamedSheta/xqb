package xqb

import (
	"reflect"
	"strings"
)

// relationAggregator handles aggregating JOIN results into a single struct
// Example: []map{{"id": 1, "posts_title": "A"}, {"id": 1, "posts_title": "B"}}
//
//	-> User{ID: 1, Posts: []Post{{Title: "A"}, {Title: "B"}}}
type relationAggregator struct{}

// aggregateToStruct combines flat relational data into nested struct with slices
func (a *relationAggregator) aggregateToStruct(dataSlice []map[string]any, structValue reflect.Value) error {
	if len(dataSlice) == 0 {
		return nil
	}

	// Step 1: Bind main struct fields from first row
	mapper := &structMapper{}
	if err := mapper.mapToStruct(dataSlice[0], structValue); err != nil {
		return err
	}

	// Step 2: Aggregate slice fields from all rows
	return a.aggregateSliceFields(dataSlice, structValue)
}

// aggregateSliceFields processes all slice fields in the struct
func (a *relationAggregator) aggregateSliceFields(dataSlice []map[string]any, structValue reflect.Value) error {
	structType := structValue.Type()

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		fieldValue := structValue.Field(i)

		if !a.isAggregableSliceField(fieldValue) {
			continue
		}

		columnName := getColumnNameForField(field)
		if columnName == "-" {
			continue
		}

		// Get table name from tag, fallback to column name
		tableName := a.getTableName(field, columnName)

		// Aggregate related data into this slice field
		if err := a.aggregateRelatedData(dataSlice, fieldValue, tableName); err != nil {
			return err
		}
	}

	return nil
}

// isAggregableSliceField checks if a field can be aggregated
func (a *relationAggregator) isAggregableSliceField(fieldValue reflect.Value) bool {
	return fieldValue.CanSet() && fieldValue.Kind() == reflect.Slice
}

// getTableName extracts table name from field tag or uses column name
func (a *relationAggregator) getTableName(field reflect.StructField, columnName string) string {
	tableName := field.Tag.Get("table")
	if tableName == "" {
		tableName = columnName
	}
	return tableName
}

// aggregateRelatedData collects all related records into a slice
func (a *relationAggregator) aggregateRelatedData(dataSlice []map[string]any, sliceField reflect.Value, tableName string) error {
	sliceType := sliceField.Type()
	elementType := sliceType.Elem()

	// Get valid field names for the element type
	validColumns := a.getValidColumns(elementType)
	columnPrefix := tableName + "_"

	// Build the new slice
	newSlice := reflect.MakeSlice(sliceType, 0, len(dataSlice))
	mapper := &structMapper{}

	for _, rowData := range dataSlice {
		// Extract data for this element
		elementData := a.extractElementData(rowData, columnPrefix, validColumns)

		// Only create element if we found data
		if len(elementData) > 0 {
			element := reflect.New(elementType).Elem()

			if err := mapper.mapToStruct(elementData, element); err != nil {
				return err
			}

			newSlice = reflect.Append(newSlice, element)
		}
	}

	sliceField.Set(newSlice)
	return nil
}

// getValidColumns returns all valid column names for a struct type
func (a *relationAggregator) getValidColumns(structType reflect.Type) map[string]bool {
	validColumns := make(map[string]bool)

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		columnName := getColumnNameForField(field)
		if columnName != "-" {
			validColumns[columnName] = true
		}
	}

	return validColumns
}

// extractElementData filters row data for columns belonging to this element
// Example: {"posts_title": "A", "posts_id": 1} -> {"title": "A", "id": 1}
func (a *relationAggregator) extractElementData(rowData map[string]any, prefix string, validColumns map[string]bool) map[string]any {
	elementData := make(map[string]any)

	for key, value := range rowData {
		// Match columns with the table prefix
		if strings.HasPrefix(key, prefix) && len(key) > len(prefix) {
			columnName := strings.TrimPrefix(key, prefix)

			// Only include if this column exists in the struct
			if validColumns[columnName] {
				elementData[columnName] = value
			}
		}
	}

	return elementData
}
