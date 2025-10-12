package querystring

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"google.golang.org/protobuf/proto"
)

var defaultUnmarshalOptions = &UnmarshalOptions{
	StructTag:    "json",
	ProtoOptions: proto.UnmarshalOptions{DiscardUnknown: true},
}

var strictUnmarshalOptions = &UnmarshalOptions{
	StructTag:    "json",
	ProtoOptions: proto.UnmarshalOptions{DiscardUnknown: false},
}

func Unmarshal(q url.Values, data any) error {
	return defaultUnmarshalOptions.Unmarshal(q, data)
}

func UnmarshalStrict(q url.Values, data any) error {
	return strictUnmarshalOptions.Unmarshal(q, data)
}

type UnmarshalOptions struct {
	StructTag    string
	ProtoOptions proto.UnmarshalOptions
}

func (u *UnmarshalOptions) Unmarshal(query url.Values, value any) error {
	refv := reflect.ValueOf(value)
	if refv.Kind() != reflect.Ptr {
		return fmt.Errorf("data must be a pointer to a struct")
	}

	refv = refv.Elem()
	if refv.Kind() != reflect.Struct {
		return fmt.Errorf("data must be a pointer to a struct")
	}

	reft := refv.Type()

	for i := 0; i < reft.NumField(); i++ {
		value := refv.Field(i)

		if !value.CanSet() {
			continue
		}

		param := u.queryParamName(reft.Field(i))
		if param == "-" {
			continue
		}

		if err := u.setFieldValue(value, query[param]); err != nil {
			return fmt.Errorf("error setting field %s: %w", reft.Field(i).Name, err)
		}
	}

	// for backwards compatibility entire query might be provided in q parameter
	if q, ok := query["q"]; ok && len(q) > 0 {
		_ = u.unmarshalValue([]byte(q[0]), value)
	}

	return nil
}

func (u *UnmarshalOptions) queryParamName(field reflect.StructField) string {
	if val := field.Tag.Get(u.StructTag); val != "" {
		return strings.Split(val, ",")[0]
	}

	return field.Name
}

func (u *UnmarshalOptions) setFieldValue(fv reflect.Value, values []string) error {
	ft := fv.Type()

	if ft.Kind() == reflect.Slice || ft.Kind() == reflect.Array {
		return u.setSliceValue(fv, values)
	}

	if len(values) == 0 {
		return nil
	}

	value := values[0]

	switch ft.Kind() {
	case reflect.String:
		fv.SetString(value)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		val, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fmt.Errorf("cannot parse %q as int: %w", value, err)
		}

		fv.SetInt(val)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		val, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return fmt.Errorf("cannot parse %q as uint: %w", value, err)
		}

		fv.SetUint(val)
	case reflect.Float32, reflect.Float64:
		val, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("cannot parse %q as float: %w", value, err)
		}

		fv.SetFloat(val)
	case reflect.Bool:
		val, err := strconv.ParseBool(value)
		if err != nil {
			return fmt.Errorf("cannot parse %q as bool: %w", value, err)
		}

		fv.SetBool(val)
	default:
		return u.setComplexValue(fv, value)
	}

	return nil
}

func (u *UnmarshalOptions) setSliceValue(fv reflect.Value, values []string) error {
	slice := reflect.MakeSlice(fv.Type(), len(values), len(values))

	for index, value := range values {
		element := slice.Index(index)

		switch fv.Type().Elem().Kind() {
		case reflect.String:
			element.SetString(value)

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			val, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				return fmt.Errorf("cannot parse %q as int: %w", value, err)
			}

			element.SetInt(val)

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			val, err := strconv.ParseUint(value, 10, 64)
			if err != nil {
				return fmt.Errorf("cannot parse %q as uint: %w", value, err)
			}

			element.SetUint(val)

		case reflect.Float32, reflect.Float64:
			val, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return fmt.Errorf("cannot parse %q as float: %w", value, err)
			}

			element.SetFloat(val)

		case reflect.Bool:
			val, err := strconv.ParseBool(value)
			if err != nil {
				return fmt.Errorf("cannot parse %q as bool: %w", value, err)
			}

			element.SetBool(val)

		default:
			if err := u.setComplexValue(element, value); err != nil {
				return err
			}
		}
	}

	fv.Set(slice)

	return nil
}

func (u *UnmarshalOptions) setComplexValue(fv reflect.Value, value string) error {
	if fv.Kind() == reflect.Ptr {
		if fv.IsNil() {
			fv.Set(reflect.New(fv.Type().Elem()))
		}

		fv = fv.Elem()
	}

	return u.unmarshalValue([]byte(value), fv.Addr().Interface())
}

func (u *UnmarshalOptions) unmarshalValue(data []byte, value any) error {
	if m, ok := value.(proto.Message); ok {
		return u.ProtoOptions.Unmarshal(data, m)
	}

	return json.Unmarshal(data, value)
}
