package zinc

import (
	"reflect"
)

func fromPtr[T any](p *T) T {
	if p == nil {
		var empty T
		return empty
	} else {
		return *p
	}
}

func deref(t reflect.Type) reflect.Type {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

func isAny(t reflect.Type) bool {
	return t == typAny
}

func isAnyPtr(t reflect.Type) bool {
	if t.Kind() != reflect.Ptr {
		return false
	}
	return isAny(t.Elem())
}

func isPrimitive(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.String:
		return true
	case reflect.Bool:
		return true
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true
	case reflect.Float32, reflect.Float64:
		return true
	default:
		if t == typBytes || t == typRawBytes {
			return true
		} else {
			return false
		}
	}
}

func isPrimitivePtr(t reflect.Type) (reflect.Type, bool) {
	if t.Kind() != reflect.Ptr {
		return nil, false
	}
	et := t.Elem()
	if !isPrimitive(et) {
		return nil, false
	}
	return et, true
}

func isRowMap(t reflect.Type) bool {
	if t.Kind() != reflect.Map {
		return false
	}
	if t.Key() != typString {
		return false
	}
	return true
}

func isRowMapPtr(t reflect.Type) (reflect.Type, bool) {
	if t.Kind() != reflect.Ptr {
		return nil, false
	}
	et := t.Elem()
	if !isRowMap(et) {
		return nil, false
	}
	return et, true
}

func isSlice(t reflect.Type) bool {
	return t.Kind() == reflect.Slice
}

func isSlicePtr(t reflect.Type) (reflect.Type, bool) {
	if t.Kind() != reflect.Ptr {
		return nil, false
	}
	et := t.Elem()
	if !isSlice(et) {
		return nil, false
	}
	return et, true
}

func isStruct(t reflect.Type) bool {
	return t.Kind() == reflect.Struct
}

func isStructPtr(t reflect.Type) (reflect.Type, bool) {
	if t.Kind() != reflect.Ptr {
		return nil, false
	}
	et := t.Elem()
	if !isStruct(et) {
		return nil, false
	}
	return et, true
}

func isStructPtrPtr(t reflect.Type) (reflect.Type, bool) {
	if t.Kind() != reflect.Ptr {
		return nil, false
	}
	return isStructPtr(t.Elem())
}
