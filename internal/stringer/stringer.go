package stringer

import (
	"fmt"
	"reflect"
	"strings"
)

func String(o any) string {
	t := reflect.TypeOf(o)

	switch t.Kind() {
	case reflect.Struct:
		tVal := reflect.ValueOf(o)

		builder := strings.Builder{}
		builder.WriteString(t.String())
		builder.WriteString("{")

		numField := t.NumField()
		if numField > 0 {
			f, fVal := t.Field(0), tVal.Field(0)
			builder.WriteString(f.Name)
			builder.WriteString(": ")
			builder.WriteString(fmt.Sprint(fVal))

			for i := 1; i < numField; i++ {
				f, fVal := t.Field(i), tVal.Field(i)
				builder.WriteString(", ")
				builder.WriteString(f.Name)
				builder.WriteString(": ")
				builder.WriteString(fmt.Sprint(fVal))
			}
		}

		builder.WriteString("}")
		return builder.String()

	case reflect.Pointer:
		elem := reflect.ValueOf(o).Elem()
		if elem.CanInterface() {
			return String(elem.Interface())
		}
		return fmt.Sprint(o)
	default:
		elem := reflect.ValueOf(o)
		if elem.CanInterface() {
			return fmt.Sprint(elem.Interface())
		}
		return fmt.Sprint(o)
	}
}
