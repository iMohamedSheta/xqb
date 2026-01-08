package xqb

import (
	"database/sql"
	"reflect"
	"time"
)

// fieldValueSetter handles setting different types of field values
type fieldValueSetter struct{}

// setFieldValue assigns a value to a field based on its type
func (s *fieldValueSetter) setFieldValue(fieldValue reflect.Value, dataValue any) error {
	// Priority order for type handling
	handlers := []func(reflect.Value, any) (bool, error){
		s.trySetJSONField,
		s.trySetSliceField,
		s.trySetScannerField,
		s.trySetSQLNullField,
		s.trySetTimeField,
		s.trySetConvertibleField,
		s.trySetNumericField,
		s.trySetBoolField,
		s.trySetStringField,
	}

	for _, handler := range handlers {
		handled, err := handler(fieldValue, dataValue)
		if err != nil {
			return err
		}
		if handled {
			return nil
		}
	}

	// Value not convertible - skip silently
	return nil
}

// trySetJSONField handles JSON/JSONB fields
func (s *fieldValueSetter) trySetJSONField(fieldValue reflect.Value, dataValue any) (bool, error) {
	if !isJSONFieldType(fieldValue.Type()) {
		return false, nil
	}

	jsonSetter := &jsonFieldSetter{}
	return true, jsonSetter.setJSONValue(fieldValue, dataValue)
}

// trySetSliceField handles slice fields
func (s *fieldValueSetter) trySetSliceField(fieldValue reflect.Value, dataValue any) (bool, error) {
	if fieldValue.Kind() != reflect.Slice {
		return false, nil
	}

	sliceSetter := &sliceFieldSetter{}
	return true, sliceSetter.setSliceValue(fieldValue, dataValue)
}

// trySetScannerField handles sql.Scanner implementations
func (s *fieldValueSetter) trySetScannerField(fieldValue reflect.Value, dataValue any) (bool, error) {
	scannerType := reflect.TypeOf((*sql.Scanner)(nil)).Elem()
	if !fieldValue.Type().Implements(scannerType) {
		return false, nil
	}

	scanner := reflect.New(fieldValue.Type()).Interface().(sql.Scanner)
	if err := scanner.Scan(dataValue); err != nil {
		return true, err
	}

	fieldValue.Set(reflect.ValueOf(scanner).Elem())
	return true, nil
}

// trySetSQLNullField handles sql.Null* types
func (s *fieldValueSetter) trySetSQLNullField(fieldValue reflect.Value, dataValue any) (bool, error) {
	if !isSQLNullType(fieldValue.Type()) {
		return false, nil
	}

	nullSetter := &sqlNullSetter{}
	nullSetter.setSQLNullValue(fieldValue, dataValue)
	return true, nil
}

// trySetTimeField handles time.Time fields
func (s *fieldValueSetter) trySetTimeField(fieldValue reflect.Value, dataValue any) (bool, error) {
	if fieldValue.Type() != reflect.TypeOf(time.Time{}) {
		return false, nil
	}

	timeSetter := &timeFieldSetter{}
	return timeSetter.setTimeValue(fieldValue, dataValue)
}

// trySetConvertibleField handles directly convertible types
func (s *fieldValueSetter) trySetConvertibleField(fieldValue reflect.Value, dataValue any) (bool, error) {
	dataValueReflect := reflect.ValueOf(dataValue)
	if !dataValueReflect.IsValid() {
		return true, nil // nil value - skip
	}

	if dataValueReflect.Type().ConvertibleTo(fieldValue.Type()) {
		fieldValue.Set(dataValueReflect.Convert(fieldValue.Type()))
		return true, nil
	}

	return false, nil
}

// trySetNumericField handles numeric conversions
func (s *fieldValueSetter) trySetNumericField(fieldValue reflect.Value, dataValue any) (bool, error) {
	converter := &numericConverter{}

	switch fieldValue.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if value, ok := converter.toInt64(dataValue); ok {
			fieldValue.SetInt(value)
			return true, nil
		}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		if value, ok := converter.toUint64(dataValue); ok {
			fieldValue.SetUint(value)
			return true, nil
		}

	case reflect.Float32, reflect.Float64:
		if value, ok := converter.toFloat64(dataValue); ok {
			fieldValue.SetFloat(value)
			return true, nil
		}
	}

	return false, nil
}

// trySetBoolField handles boolean conversions
func (s *fieldValueSetter) trySetBoolField(fieldValue reflect.Value, dataValue any) (bool, error) {
	if fieldValue.Kind() != reflect.Bool {
		return false, nil
	}

	boolValue := s.convertToBool(dataValue)
	if boolValue != nil {
		fieldValue.SetBool(*boolValue)
		return true, nil
	}

	return false, nil
}

// convertToBool converts various types to bool
func (s *fieldValueSetter) convertToBool(value any) *bool {
	switch v := value.(type) {
	case bool:
		return &v
	case int, int8, int16, int32, int64:
		result := reflect.ValueOf(v).Int() != 0
		return &result
	case uint, uint8, uint16, uint32, uint64:
		result := reflect.ValueOf(v).Uint() != 0
		return &result
	case float32, float64:
		result := reflect.ValueOf(v).Float() != 0
		return &result
	}
	return nil
}

// trySetStringField handles string assignment
func (s *fieldValueSetter) trySetStringField(fieldValue reflect.Value, dataValue any) (bool, error) {
	if fieldValue.Kind() != reflect.String {
		return false, nil
	}

	if str, ok := dataValue.(string); ok {
		fieldValue.SetString(str)
		return true, nil
	}

	return false, nil
}
