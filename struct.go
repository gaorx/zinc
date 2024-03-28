package zinc

import (
	"fmt"
	"reflect"
	"sync"
)

type Struct struct {
	Name   string         `json:"name"`
	Fields []*StructField `json:"fields"`
}

type StructField struct {
	Paths    []*StructFieldPath `json:"paths"`
	Type     reflect.Type       `json:"-"`
	TypeName string             `json:"type"`
	TagCol   string             `json:"zcol"`
	TagCols  []string           `json:"zcols"`
}

type StructFieldPath struct {
	Name  string `json:"name"`
	Index []int  `json:"index"`
}

func ParseStruct(t reflect.Type) (*Struct, bool) {
	if t == nil {
		return nil, false
	} else if isStruct(t) {
		return parseStruct1(t), true
	} else if t1, ok := isStructPtr(t); ok {
		return ParseStruct(t1)
	} else if t1, ok := isStructPtrPtr(t); ok {
		return ParseStruct(t1)
	} else {
		return nil, false
	}
}

var (
	structCache      = map[reflect.Type]*Struct{}
	structCacheMutex = sync.RWMutex{}
)

func parseStruct1(t reflect.Type) *Struct {
	var cached *Struct
	lockR(&structCacheMutex, func() {
		cached = structCache[t]
	})
	if cached != nil {
		return cached
	}
	var r *Struct
	lockW(&structCacheMutex, func() {
		s := parseStruct0(t)
		structCache[t] = s
		r = s
	})
	return r
}

func parseStruct0(t reflect.Type) *Struct {
	fmt.Println("**********")
	s := &Struct{
		Name: t.Name(),
	}
	parseStruct0To(t, nil, s)
	return s
}

func parseStruct0To(t reflect.Type, currentPaths []*StructFieldPath, target *Struct) {
	fieldNum := t.NumField()
	for i := 0; i < fieldNum; i++ {
		f := t.Field(i)
		if !f.IsExported() {
			continue
		}

		if f.Anonymous {
			embedded := derefDeep(f.Type)
			if isStruct(embedded) {
				parseStruct0To(embedded, append(currentPaths, &StructFieldPath{Name: f.Name, Index: f.Index}), target)
			}
		} else {
			col := f.Tag.Get("zcol")
			cols := splitNonEmpty(f.Tag.Get("zcols"), ",")
			target.Fields = append(target.Fields, &StructField{
				Paths:    append(cloneSlice(currentPaths), &StructFieldPath{Name: f.Name, Index: f.Index}),
				Type:     f.Type,
				TypeName: f.Type.String(),
				TagCol:   col,
				TagCols:  cols,
			})
		}
	}
}
