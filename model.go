package xqb

import (
	"database/sql"
	"fmt"
	"reflect"
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

// Bind
func Bind(data any, dest any) error {
	destVal := reflect.ValueOf(dest)
	if destVal.Kind() != reflect.Ptr || destVal.IsNil() {
		return fmt.Errorf("dest must be a non-nil pointer")
	}

	destElem := destVal.Elem()

	switch destElem.Kind() {
	case reflect.Struct:
		// Expecting map[string]any for single struct
		dataMap, ok := data.(map[string]any)
		if !ok {
			return fmt.Errorf("data must be map[string]any for struct binding")
		}
		return bindStruct(dataMap, destElem)

	case reflect.Slice:
		// Expecting []map[string]any for slice binding
		dataSlice, ok := data.([]map[string]any)
		if !ok {
			return fmt.Errorf("data must be []map[string]any for slice binding")
		}

		elemType := destElem.Type().Elem()
		for _, item := range dataSlice {
			newElem := reflect.New(elemType).Elem()
			if err := bindStruct(item, newElem); err != nil {
				return err
			}
			destElem.Set(reflect.Append(destElem, newElem))
		}
		return nil

	default:
		return fmt.Errorf("unsupported destination type: %s", destElem.Kind())
	}
}

// bindStruct binds a struct to a value
func bindStruct(data map[string]any, v reflect.Value) error {
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		col := field.Tag.Get("xqb")
		if col == "" {
			col = field.Name
		}
		if col == "-" {
			continue
		}

		val, ok := data[col]
		if !ok || val == nil {
			continue
		}

		fv := v.Field(i)
		if !fv.CanSet() {
			continue
		}

		switch fv.Type() {
		case reflect.TypeOf(sql.NullString{}):
			if str, ok := val.(string); ok {
				fv.Set(reflect.ValueOf(sql.NullString{String: str, Valid: true}))
			}
		case reflect.TypeOf(sql.NullInt64{}):
			if num, ok := val.(int64); ok {
				fv.Set(reflect.ValueOf(sql.NullInt64{Int64: num, Valid: true}))
			}
		case reflect.TypeOf(sql.NullBool{}):
			if b, ok := val.(bool); ok {
				fv.Set(reflect.ValueOf(sql.NullBool{Bool: b, Valid: true}))
			}
		case reflect.TypeOf(sql.NullFloat64{}):
			if f, ok := val.(float64); ok {
				fv.Set(reflect.ValueOf(sql.NullFloat64{Float64: f, Valid: true}))
			}
		case reflect.TypeOf(sql.NullTime{}):
			if tm, ok := val.(time.Time); ok {
				fv.Set(reflect.ValueOf(sql.NullTime{Time: tm, Valid: true}))
			}
		case reflect.TypeOf(&time.Time{}):
			if tm, ok := val.(time.Time); ok {
				fv.Set(reflect.ValueOf(&tm))
			}
		default:
			if reflect.TypeOf(val).ConvertibleTo(fv.Type()) {
				fv.Set(reflect.ValueOf(val).Convert(fv.Type()))
			}
		}
	}

	return nil
}
