package xqb

import (
	"reflect"
)

// sliceMapper handles mapping a slice of maps to a slice of structs
type sliceMapper struct{}

// mapToSlice binds each map in the slice to a struct
func (m *sliceMapper) mapToSlice(dataSlice []map[string]any, sliceValue reflect.Value) error {
	elementType := sliceValue.Type().Elem()
	mapper := &structMapper{}

	for _, dataMap := range dataSlice {
		// Create new struct instance
		newElement := reflect.New(elementType).Elem()

		// Map data to the struct
		if err := mapper.mapToStruct(dataMap, newElement); err != nil {
			return err
		}

		// Append to result slice
		sliceValue.Set(reflect.Append(sliceValue, newElement))
	}

	return nil
}
