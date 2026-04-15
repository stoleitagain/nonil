package nonil

import (
	"reflect"
	"strings"
)

func NoNil(input interface{}) map[string]interface{} {
	out := make(map[string]interface{})
	buildMap(reflect.ValueOf(input), out)
	return out
}

func buildMap(v reflect.Value, out map[string]interface{}) {
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return
		}
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return
	}

	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		fv := v.Field(i)
		ft := t.Field(i)

		if ft.PkgPath != "" {
			continue
		}

		tag := ft.Tag.Get("json")
		name, opts := parseTag(tag)

		if name == "-" {
			continue
		}

		if name == "" {
			name = ft.Name
		}

		switch fv.Kind() {
		case reflect.Ptr, reflect.Interface, reflect.Map, reflect.Slice:
			if fv.IsNil() {
				continue
			}
		}

		if fv.Kind() == reflect.Struct && !isPrimitiveStruct(fv) {
			nested := make(map[string]interface{})
			buildMap(fv, nested)

			if len(nested) == 0 && opts["omitempty"] {
				continue
			}

			out[name] = nested
			continue
		}

		val := fv.Interface()
		if fv.Kind() == reflect.Ptr {
			val = fv.Elem().Interface()
		}

		if opts["omitempty"] && isZeroValue(fv) {
			continue
		}

		out[name] = val
	}
}

func parseTag(tag string) (string, map[string]bool) {
	opts := make(map[string]bool)

	if tag == "" {
		return "", opts
	}

	parts := strings.Split(tag, ",")
	name := parts[0]

	for _, opt := range parts[1:] {
		opts[opt] = true
	}

	return name, opts
}

func isZeroValue(v reflect.Value) bool {
	return v.IsZero()
}

func isPrimitiveStruct(v reflect.Value) bool {
	return v.Type().PkgPath() == "time"
}