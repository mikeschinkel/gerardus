package cli

import (
	"reflect"
	"strconv"
)

type Value struct {
	Type   reflect.Kind
	string string
	int    int
}

func NewValue(t reflect.Kind, value any) (v *Value) {
	v = &Value{Type: t}
	switch t {
	case reflect.String:
		v.string = v.asString(value)
	case reflect.Int:
		v.int = v.asInt(value)
	default:
		v.unsupportedType()
	}
	return v
}

func (v Value) IsZero() (zero bool) {
	switch v.Type {
	case reflect.String:
		zero = v.string == ""
	case reflect.Int:
		zero = v.int == 0
	case reflect.Invalid:
		zero = true
	default:
		v.unsupportedType()
	}
	return zero
}

func (v Value) asString(value any) (s string) {
	switch t := value.(type) {
	case string:
		s = t
	case int:
		s = strconv.Itoa(t)
	case nil:
		s = ""
	default:
		panicf("AsString does not (yet?) support type '%T'", value)
	}
	return s
}
func (v Value) asInt(value any) (n int) {
	switch t := value.(type) {
	case string:
		n, _ = strconv.Atoi(t)
	case int:
		n = t
	case nil:
		n = 0
	default:
		panicf("AsString does not (yet?) support type '%T'", value)
	}
	return n
}

func (v Value) Int() int {
	switch v.Type {
	case reflect.Int:
		return v.int
	default:
		v.unsupportedType()
	}
	return 0
}

func (v Value) String() string {
	switch v.Type {
	case reflect.String:
		return v.string
	case reflect.Int:
		return strconv.Itoa(v.int)
	default:
		v.unsupportedType()
	}
	return ""
}

func (v Value) unsupportedType() {
	panicf("Value does not (yet?) support type '%s'", v.Type)
}
