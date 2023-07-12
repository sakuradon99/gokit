package web

import (
	"reflect"
	"strconv"
)

func setVal(fieldValue reflect.Value, val string) error {
	switch fieldValue.Interface().(type) {
	case string:
		fieldValue.SetString(val)
	case int64, int32, int16, int8, int:
		int64Val, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return ErrInvalidParams
		}
		fieldValue.SetInt(int64Val)
	case uint64, uint32, uint16, uint8, uint:
		uint64Val, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			return ErrInvalidParams
		}
		fieldValue.SetUint(uint64Val)
	case *string:
		if val != "" {
			p := reflect.New(fieldValue.Type().Elem())
			p.Elem().SetString(val)
			fieldValue.Set(p)
		}
	case *int64, *int32, *int16, *int8, *int:
		int64Val, err := strconv.ParseInt(val, 10, 64)
		if err == nil {
			p := reflect.New(fieldValue.Type().Elem())
			p.Elem().SetInt(int64Val)
			fieldValue.Set(p)
		}
	case *uint64, *uint32, *uint16, *uint8, *uint:
		uint64Val, err := strconv.ParseUint(val, 10, 64)
		if err == nil {
			p := reflect.New(fieldValue.Type().Elem())
			p.Elem().SetUint(uint64Val)
			fieldValue.Set(p)
		}
	}
	return nil
}
