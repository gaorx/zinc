package zinc

import (
	"database/sql"
	"reflect"
)

type Mapper func(src *Src, dest *Dest, opts *Options) error

type Src struct {
	Rows        *sql.Rows
	Columns     []string
	ColumnTypes []*sql.ColumnType
}

type Dest struct {
	Target any
	Kind   int
	Type   reflect.Type
}

const (
	StructPtr    = 1
	Map          = 2
	PrimitivePtr = 3
	SlicePtr     = 4
)

func (src *Src) fetchColumns(force bool) error {
	if force || len(src.Columns) <= 0 {
		if cols, err := src.Rows.Columns(); err != nil {
			return err
		} else {
			src.Columns = cols
		}
	}
	if force || len(src.ColumnTypes) <= 0 {
		if colTypes, err := src.Rows.ColumnTypes(); err != nil {
			return err
		} else {
			src.ColumnTypes = colTypes
		}
	}
	return nil
}

func (src *Src) newDestSlice(opts *Options) []any {
	dialect := opts.Dialect
	slice := make([]any, len(src.Columns))
	for i := range slice {
		ci := src.ColumnTypes[i]
		slice[i] = dialect.NewDest(ci, opts)
	}
	return slice
}

func (src *Src) NextResultSet() (bool, error) {
	if hasNext := src.Rows.NextResultSet(); hasNext {
		return true, src.fetchColumns(true)
	} else {
		return false, nil
	}
}

func defaultMapper(src *Src, dest *Dest, opts *Options) error {
	destSlice := src.newDestSlice(opts)
	if err := src.Rows.Scan(destSlice...); err != nil {
		return err
	}
	coerceDest := func(ci *sql.ColumnType, scannedVal reflect.Value, toType reflect.Type) (reflect.Value, error) {
		return opts.Dialect.CoerceDest(ci, scannedVal, toType, opts)
	}
	switch dest.Kind {
	case StructPtr:
		panic("not implemented")
	case Map:
		dv := reflect.ValueOf(dest.Target)
		for i := range src.Columns {
			ci, cn, dv1 := src.ColumnTypes[i], src.Columns[i], reflect.ValueOf(destSlice[i]).Elem()
			targetVal, err := coerceDest(ci, dv1, dest.Type.Elem())
			if err != nil {
				return ErrCoerceDest
			}
			dv.SetMapIndex(reflect.ValueOf(cn), targetVal)
		}
		return nil
	case PrimitivePtr:
		dv := reflect.ValueOf(dest.Target)
		dv0 := reflect.ValueOf(destSlice[0]).Elem()
		ci0 := src.ColumnTypes[0]
		targetVal, err := coerceDest(ci0, dv0, dest.Type.Elem())
		if err != nil {
			return ErrCoerceDest
		}
		dv.Elem().Set(targetVal)
		return nil
	case SlicePtr:
		dv := reflect.ValueOf(dest.Target)
		rv := reflect.MakeSlice(dest.Type, 0, len(destSlice))
		for i := range src.Columns {
			ci1 := src.ColumnTypes[i]
			dv1 := reflect.ValueOf(destSlice[i]).Elem()
			targetVal, err := coerceDest(ci1, dv1, dest.Type.Elem())
			if err != nil {
				return ErrCoerceDest
			}
			rv = reflect.Append(rv, targetVal)
		}
		dv.Elem().Set(rv)
		return nil
	default:
		panic("unreachable")
	}
}
