package xqb

import (
	"fmt"
	"reflect"
)

// ModelInterface represents a model with a table name
type ModelInterface interface {
	Table() string
}

// ModelQ creates a query builder for a model
func ModelQ(model ModelInterface) *QueryBuilder {
	return Table(model.Table())
}

// Bind maps database query results to Go structs
// Supports three patterns:
//   - Single map → Single struct:     Bind(map[string]any{...}, &user)
//   - Slice of maps → Slice of structs: Bind([]map[string]any{...}, &users)
//   - Slice of maps → Single struct:    Bind([]map[string]any{...}, &user) // for JOINs
func Bind(sourceData any, destination any) error {
	binder := &dataBinder{}
	return binder.bind(sourceData, destination)
}

// dataBinder handles the binding logic
type dataBinder struct{}

// bind is the main entry point for binding data to structs
func (b *dataBinder) bind(sourceData any, destination any) error {
	destValue := reflect.ValueOf(destination)

	if err := b.validateDestination(destValue); err != nil {
		return err
	}

	destElement := destValue.Elem()

	switch destElement.Kind() {
	case reflect.Struct:
		return b.bindToStruct(sourceData, destElement)
	case reflect.Slice:
		return b.bindToSlice(sourceData, destElement)
	default:
		return fmt.Errorf("destination must be pointer to struct or slice, got %s", destElement.Kind())
	}
}

// validateDestination ensures the destination is a non-nil pointer
func (b *dataBinder) validateDestination(destValue reflect.Value) error {
	if destValue.Kind() != reflect.Ptr {
		return fmt.Errorf("destination must be a pointer, got %s", destValue.Kind())
	}
	if destValue.IsNil() {
		return fmt.Errorf("destination pointer cannot be nil")
	}
	return nil
}

// bindToStruct handles binding to a single struct
func (b *dataBinder) bindToStruct(sourceData any, structValue reflect.Value) error {
	// Case 1: Single map → Single struct
	if dataMap, ok := sourceData.(map[string]any); ok {
		mapper := &structMapper{}
		return mapper.mapToStruct(dataMap, structValue)
	}

	// Case 2: Slice of maps → Single struct (for JOIN results)
	if dataSlice, ok := sourceData.([]map[string]any); ok {
		aggregator := &relationAggregator{}
		return aggregator.aggregateToStruct(dataSlice, structValue)
	}

	return fmt.Errorf("source data must be map[string]any or []map[string]any for struct binding")
}

// bindToSlice handles binding to a slice of structs
func (b *dataBinder) bindToSlice(sourceData any, sliceValue reflect.Value) error {
	dataSlice, ok := sourceData.([]map[string]any)
	if !ok {
		return fmt.Errorf("source data must be []map[string]any for slice binding")
	}

	mapper := &sliceMapper{}
	return mapper.mapToSlice(dataSlice, sliceValue)
}
