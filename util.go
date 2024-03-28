package zinc

import (
	"reflect"
	"strings"
	"sync"
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
		return t.Elem()
	}
	return t
}

func derefDeep(t reflect.Type) reflect.Type {
	for t.Kind() == reflect.Ptr {
		return derefDeep(t.Elem())
	}
	return t
}

func cloneSlice[T any](slice []T) []T {
	if slice == nil {
		return nil
	}
	r := make([]T, len(slice))
	copy(r, slice)
	return r
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

func sliceContains[T comparable](slice []T, target T) bool {
	for _, v := range slice {
		if v == target {
			return true
		}
	}
	return false
}

func b2s(b []byte, charset string) (string, bool) {
	// TODO: charset
	return string(b), true
}

func lockR(m *sync.RWMutex, f func()) {
	m.RLock()
	defer m.RUnlock()
	f()
}

func lockW(m *sync.RWMutex, f func()) {
	m.Lock()
	defer m.Unlock()
	f()
}

func splitNonEmpty(s string, sep string) []string {
	if s == "" {
		return nil
	}
	var r []string
	for _, e := range strings.Split(s, sep) {
		e = strings.TrimSpace(e)
		if e != "" {
			r = append(r, e)
		}
	}
	return r
}
