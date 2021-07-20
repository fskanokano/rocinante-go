package bind

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type Binder interface {
	Bind(s interface{}) error
}

func unmarshalValues(s interface{}, data url.Values, tagName string) error {
	t := reflect.TypeOf(s).Elem()
	v := reflect.ValueOf(s).Elem()
	for i := 0; i < v.NumField(); i++ {
		fieldT := t.Field(i)
		fieldV := v.Field(i)

		if !fieldV.CanSet() {
			return errors.New(fmt.Sprintf("invalid field '%s.%s'", t.Name(), fieldT.Name))
		}

		fieldKey := getFieldKey(fieldT, tagName)
		if values, ok := data[fieldKey]; ok {
			if len(values) >= 1 {
				value := values[0]
				valueT := reflect.TypeOf(value)
				if valueT.Kind() == fieldT.Type.Kind() {
					fieldV.Set(reflect.ValueOf(value))
				} else {
					cv, err := typeConversion(t, fieldT, value, tagName, fieldKey)
					if err != nil {
						return err
					}
					fieldV.Set(cv)
				}
			}
		}

	}
	return nil
}

func getFieldKey(t reflect.StructField, tagName string) string {
	if fieldKey, ok := t.Tag.Lookup(tagName); ok {
		return fieldKey
	}
	return strings.ToLower(t.Name)
}

func typeConversion(t reflect.Type, fieldT reflect.StructField, value string, tagName string, fieldKey string) (reflect.Value, error) {
	errS := fmt.Sprintf("the value of '%s.%s' '%s' cannot be resolved into the field '%s.%s' because it is not"+
		" a valid '%s' type.", tagName, fieldKey, value, t.Name(), fieldT.Name, fieldT.Type.Name())
	switch fieldT.Type.Kind() {
	case reflect.Int:
		if cv, err := strconv.Atoi(value); err != nil {
			return reflect.Value{}, errors.New(errS)
		} else {
			return reflect.ValueOf(cv), nil
		}
	case reflect.Int8:
		if cv, err := strconv.ParseInt(value, 10, 8); err != nil {
			return reflect.Value{}, errors.New(errS)
		} else {
			return reflect.ValueOf(int8(cv)), nil
		}
	case reflect.Int16:
		if cv, err := strconv.ParseInt(value, 10, 16); err != nil {
			return reflect.Value{}, errors.New(errS)
		} else {
			return reflect.ValueOf(int16(cv)), nil
		}
	case reflect.Int32:
		if cv, err := strconv.ParseInt(value, 10, 32); err != nil {
			return reflect.Value{}, errors.New(errS)
		} else {
			return reflect.ValueOf(int32(cv)), nil
		}
	case reflect.Int64:
		if cv, err := strconv.ParseInt(value, 10, 64); err != nil {
			return reflect.Value{}, errors.New(errS)
		} else {
			return reflect.ValueOf(cv), nil
		}
	case reflect.Uint:
		if cv, err := strconv.ParseUint(value, 10, 64); err != nil {
			return reflect.Value{}, errors.New(errS)
		} else {
			return reflect.ValueOf(uint(cv)), nil
		}
	case reflect.Uint8:
		if cv, err := strconv.ParseUint(value, 10, 8); err != nil {
			return reflect.Value{}, errors.New(errS)
		} else {
			return reflect.ValueOf(uint8(cv)), nil
		}
	case reflect.Uint16:
		if cv, err := strconv.ParseUint(value, 10, 16); err != nil {
			return reflect.Value{}, errors.New(errS)
		} else {
			return reflect.ValueOf(uint16(cv)), nil
		}
	case reflect.Uint32:
		if cv, err := strconv.ParseUint(value, 10, 32); err != nil {
			return reflect.Value{}, errors.New(errS)
		} else {
			return reflect.ValueOf(uint32(cv)), nil
		}
	case reflect.Uint64:
		if cv, err := strconv.ParseUint(value, 10, 32); err != nil {
			return reflect.Value{}, errors.New(errS)
		} else {
			return reflect.ValueOf(cv), nil
		}
	case reflect.Float64:
		if cv, err := strconv.ParseFloat(value, 64); err != nil {
			return reflect.Value{}, errors.New(errS)
		} else {
			return reflect.ValueOf(cv), nil
		}
	case reflect.Float32:
		if cv, err := strconv.ParseFloat(value, 32); err != nil {
			return reflect.Value{}, errors.New(errS)
		} else {
			return reflect.ValueOf(float32(cv)), nil
		}
	case reflect.Bool:
		if cv, err := strconv.ParseBool(value); err != nil {
			return reflect.Value{}, errors.New(errS)
		} else {
			return reflect.ValueOf(cv), nil
		}
	default:
		switch fieldT.Type.Name() {
		case "Time":
			if cv, err := time.ParseInLocation("2006-01-02 15:04:05", value, time.Local); err != nil {
				return reflect.Value{}, errors.New(errS)
			} else {
				return reflect.ValueOf(cv), nil
			}
		default:
			return reflect.Value{}, errors.New(errS)
		}
	}
}
