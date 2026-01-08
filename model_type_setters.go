package xqb

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"time"
)

// jsonFieldSetter handles JSON and JSONB field types
type jsonFieldSetter struct{}

func (s *jsonFieldSetter) setJSONValue(fieldValue reflect.Value, dataValue any) error {
	if dataValue == nil {
		return nil
	}

	jsonBytes, err := s.convertToJSONBytes(dataValue)
	if err != nil {
		return err
	}

	if len(jsonBytes) > 0 && !json.Valid(jsonBytes) {
		return fmt.Errorf("invalid JSON data")
	}

	return s.assignJSONToField(fieldValue, jsonBytes)
}

func (s *jsonFieldSetter) convertToJSONBytes(dataValue any) ([]byte, error) {
	switch v := dataValue.(type) {
	case []byte:
		return v, nil

	case string:
		if v == "" {
			return nil, nil
		}
		return []byte(v), nil

	case json.RawMessage:
		return []byte(v), nil

	case map[string]any, []any:
		return json.Marshal(v)

	default:
		return json.Marshal(v)
	}
}

func (s *jsonFieldSetter) assignJSONToField(fieldValue reflect.Value, jsonBytes []byte) error {
	switch fieldValue.Type() {
	case reflect.TypeOf([]byte{}):
		fieldValue.SetBytes(jsonBytes)

	case reflect.TypeOf(json.RawMessage{}):
		fieldValue.Set(reflect.ValueOf(json.RawMessage(jsonBytes)))

	default:
		if fieldValue.Kind() == reflect.Map {
			return s.unmarshalToMap(fieldValue, jsonBytes)
		}
	}

	return nil
}

func (s *jsonFieldSetter) unmarshalToMap(mapValue reflect.Value, jsonBytes []byte) error {
	if len(jsonBytes) == 0 {
		return nil
	}

	if mapValue.IsNil() {
		mapValue.Set(reflect.MakeMap(mapValue.Type()))
	}

	mapPtr := reflect.New(mapValue.Type())
	if err := json.Unmarshal(jsonBytes, mapPtr.Interface()); err != nil {
		return fmt.Errorf("failed to unmarshal JSON into map: %w", err)
	}

	mapValue.Set(mapPtr.Elem())
	return nil
}

// sliceFieldSetter handles slice field types
type sliceFieldSetter struct{}

func (s *sliceFieldSetter) setSliceValue(sliceValue reflect.Value, dataValue any) error {
	// Try to decode JSON strings/bytes first
	if decodedValue, decoded := s.tryDecodeJSON(dataValue); decoded {
		dataValue = decodedValue
	}

	reflectValue := reflect.ValueOf(dataValue)
	if !reflectValue.IsValid() {
		return nil
	}

	if !s.isSliceOrArray(reflectValue.Kind()) {
		return nil
	}

	return s.populateSlice(sliceValue, reflectValue)
}

func (s *sliceFieldSetter) tryDecodeJSON(dataValue any) (any, bool) {
	switch v := dataValue.(type) {
	case string:
		if v != "" {
			var decoded any
			if err := json.Unmarshal([]byte(v), &decoded); err == nil {
				return decoded, true
			}
		}
	case []byte:
		if len(v) > 0 {
			var decoded any
			if err := json.Unmarshal(v, &decoded); err == nil {
				return decoded, true
			}
		}
	}
	return nil, false
}

func (s *sliceFieldSetter) isSliceOrArray(kind reflect.Kind) bool {
	return kind == reflect.Slice || kind == reflect.Array
}

func (s *sliceFieldSetter) populateSlice(sliceValue reflect.Value, sourceValue reflect.Value) error {
	elementType := sliceValue.Type().Elem()
	newSlice := reflect.MakeSlice(sliceValue.Type(), 0, sourceValue.Len())
	mapper := &structMapper{}

	for i := 0; i < sourceValue.Len(); i++ {
		item := sourceValue.Index(i).Interface()
		element := reflect.New(elementType).Elem()

		if err := s.setElement(element, item, mapper); err != nil {
			return err
		}

		newSlice = reflect.Append(newSlice, element)
	}

	sliceValue.Set(newSlice)
	return nil
}

func (s *sliceFieldSetter) setElement(element reflect.Value, item any, mapper *structMapper) error {
	// Handle struct elements
	if element.Kind() == reflect.Struct {
		return s.setStructElement(element, item, mapper)
	}

	// Handle primitive elements
	return s.setPrimitiveElement(element, item)
}

func (s *sliceFieldSetter) setStructElement(element reflect.Value, item any, mapper *structMapper) error {
	if mapData, ok := item.(map[string]any); ok {
		return mapper.mapToStruct(mapData, element)
	}

	// Try to unmarshal JSON string
	if str, ok := item.(string); ok && str != "" {
		var mapData map[string]any
		if err := json.Unmarshal([]byte(str), &mapData); err == nil {
			return mapper.mapToStruct(mapData, element)
		}
	}

	return nil
}

func (s *sliceFieldSetter) setPrimitiveElement(element reflect.Value, item any) error {
	reflectValue := reflect.ValueOf(item)
	if !reflectValue.IsValid() {
		return nil
	}

	elementType := element.Type()

	// Try direct conversion
	if reflectValue.Type().ConvertibleTo(elementType) {
		element.Set(reflectValue.Convert(elementType))
		return nil
	}

	// Try numeric conversions
	converter := &numericConverter{}
	return converter.convertAndSet(element, item)
}

// timeFieldSetter handles time.Time fields
type timeFieldSetter struct{}

func (s *timeFieldSetter) setTimeValue(fieldValue reflect.Value, dataValue any) (bool, error) {
	// Direct time.Time value
	if timeVal, ok := dataValue.(time.Time); ok {
		fieldValue.Set(reflect.ValueOf(timeVal))
		return true, nil
	}

	// Try parsing string
	if str, ok := dataValue.(string); ok && str != "" {
		if parsedTime, err := time.Parse(time.RFC3339, str); err == nil {
			fieldValue.Set(reflect.ValueOf(parsedTime))
			return true, nil
		}
	}

	return false, nil
}

// sqlNullSetter handles sql.Null* types
type sqlNullSetter struct{}

func (s *sqlNullSetter) setSQLNullValue(fieldValue reflect.Value, dataValue any) {
	converter := &numericConverter{}

	switch fieldValue.Type() {
	case reflect.TypeOf(sql.NullString{}):
		if str, ok := dataValue.(string); ok {
			fieldValue.Set(reflect.ValueOf(sql.NullString{String: str, Valid: true}))
		}

	case reflect.TypeOf(sql.NullBool{}):
		if b, ok := dataValue.(bool); ok {
			fieldValue.Set(reflect.ValueOf(sql.NullBool{Bool: b, Valid: true}))
		}

	case reflect.TypeOf(sql.NullInt64{}):
		if n, ok := converter.toInt64(dataValue); ok {
			fieldValue.Set(reflect.ValueOf(sql.NullInt64{Int64: n, Valid: true}))
		}

	case reflect.TypeOf(sql.NullInt32{}):
		if n, ok := converter.toInt64(dataValue); ok {
			fieldValue.Set(reflect.ValueOf(sql.NullInt32{Int32: int32(n), Valid: true}))
		}

	case reflect.TypeOf(sql.NullInt16{}):
		if n, ok := converter.toInt64(dataValue); ok {
			fieldValue.Set(reflect.ValueOf(sql.NullInt16{Int16: int16(n), Valid: true}))
		}

	case reflect.TypeOf(sql.NullByte{}):
		if n, ok := converter.toInt64(dataValue); ok {
			fieldValue.Set(reflect.ValueOf(sql.NullByte{Byte: byte(n), Valid: true}))
		}

	case reflect.TypeOf(sql.NullFloat64{}):
		if f, ok := converter.toFloat64(dataValue); ok {
			fieldValue.Set(reflect.ValueOf(sql.NullFloat64{Float64: f, Valid: true}))
		}

	case reflect.TypeOf(sql.NullTime{}):
		if tm, ok := dataValue.(time.Time); ok {
			valid := !tm.IsZero()
			fieldValue.Set(reflect.ValueOf(sql.NullTime{Time: tm, Valid: valid}))
		}
	}
}
