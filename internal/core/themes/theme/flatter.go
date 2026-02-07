package theme

import (
	"reflect"
	"strings"
	"theme-engine/internal/core/themes/palette"
)

func flattenStruct(prefix string, v reflect.Value, out map[string]string) {
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}

	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		// skip unexported
		if !field.CanInterface() {
			continue
		}

		name := strings.ToLower(fieldType.Name)
		key := name
		if prefix != "" {
			key = prefix + "." + name
		}

		switch field.Kind() {
		case reflect.String:
			out[key] = field.String()

		case reflect.Struct:
			flattenStruct(key, field, out)
		}
	}
}

func BuildFlattenPalette(t *palette.ResolvedPalette) {
	t.Flat = make(map[string]string);

	flattenStruct("extra", reflect.ValueOf(t), t.Flat);
}

