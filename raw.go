package zinc

import (
	"database/sql"
	"fmt"
	"reflect"
)

func (db *DB) RawExec(dest any, q string, args ...any) error {
	uArgs, optsModifier := United(args...)
	bound, boundArgs, err := db.Bind(q, uArgs)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrBind, err.Error())
	}
	opts := copyOptions(db.options, optsModifier)
	var sqlRes sql.Result
	switch db.kind {
	case kindDB:
		sqlRes, err = logDo(
			db.getCtx(),
			q, uArgs,
			bound, boundArgs,
			opts,
			func() (sql.Result, error) {
				return db.db.ExecContext(db.getCtx(), bound, boundArgs...)
			},
		)
	case kindTx:
		sqlRes, err = logDo(
			db.getCtx(),
			q, uArgs,
			bound, boundArgs,
			opts,
			func() (sql.Result, error) {
				return db.tx.ExecContext(db.getCtx(), bound, boundArgs...)
			},
		)
	case kindSt:
		sqlRes, err = logDo(
			db.getCtx(),
			q, uArgs,
			bound, boundArgs,
			opts,
			func() (sql.Result, error) {
				return db.st.ExecContext(db.getCtx(), boundArgs...)
			},
		)
	default:
		panic("unreachable")
	}
	if err != nil {
		return err
	}
	switch d := dest.(type) {
	case nil:
		// do nothing
	case *sql.Result:
		*d = sqlRes
	case *Result:
		if d != nil {
			if execRes, err := newResult(sqlRes); err != nil {
				return err
			} else {
				*d = *execRes
			}
		}
	case **Result:
		if d != nil {
			if *d == nil {
				*d = &Result{}
			}
			if execRes, err := newResult(sqlRes); err != nil {
				return err
			} else {
				**d = *execRes
			}
		}
	default:
		return ErrInvalidDest
	}
	return nil
}

func (db *DB) RawQueryOne(dest any, q string, args ...any) error {
	uArgs, optsModifier := United(args...)
	bound, boundArgs, err := db.Bind(q, uArgs)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrBind, err.Error())
	}
	opts := copyOptions(db.options, optsModifier)
	var rows *sql.Rows
	switch db.kind {
	case kindDB:
		rows, err = logDo(
			db.getCtx(),
			q, uArgs,
			bound, boundArgs,
			opts,
			func() (*sql.Rows, error) {
				return db.db.QueryContext(db.getCtx(), bound, boundArgs...)
			},
		)
	case kindTx:
		rows, err = logDo(
			db.getCtx(),
			q, uArgs,
			bound, boundArgs,
			opts,
			func() (*sql.Rows, error) {
				return db.tx.QueryContext(db.getCtx(), bound, boundArgs...)
			},
		)
	case kindSt:
		rows, err = logDo(
			db.getCtx(),
			q, uArgs,
			bound, boundArgs,
			opts,
			func() (*sql.Rows, error) {
				return db.st.QueryContext(db.getCtx(), boundArgs...)
			},
		)
	default:
		panic("unreachable")
	}
	if err != nil {
		return err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	if dest == nil {
		// do nothing
		return nil
	} else if f, ok := dest.(func(*sql.Rows) error); ok {
		// custom map
		if f == nil {
			return nil
		}
		if rows.Next() {
			if err := f(rows); err != nil {
				return err
			}
			if err := rows.Err(); err != nil {
				return err
			}
			return nil
		} else {
			return ErrNoRows
		}
	} else {
		// use mapper to map
		mapper := getMapper(dest, opts)
		src := Src{Rows: rows}
		if err := src.fetchColumns(false); err != nil {
			return err
		}
		mapRow := makeMapRowFunc(reflect.TypeOf(dest), mapper, opts)
		if rows.Next() {
			if err := mapRow(&src, dest); err != nil {
				return err
			}
			if err := rows.Err(); err != nil {
				return err
			}
		} else {
			return ErrNoRows
		}
	}
	return nil
}

func (db *DB) RawQueryAll(dest any, q string, args ...any) error {
	uArgs, optsModifier := United(args...)
	bound, boundArgs, err := db.Bind(q, uArgs)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrBind, err.Error())
	}
	opts := copyOptions(db.options, optsModifier)
	var rows *sql.Rows
	switch db.kind {
	case kindDB:
		rows, err = logDo(
			db.getCtx(),
			q, uArgs,
			bound, boundArgs,
			opts,
			func() (*sql.Rows, error) {
				return db.db.QueryContext(db.getCtx(), bound, boundArgs...)
			},
		)
	case kindTx:
		rows, err = logDo(
			db.getCtx(),
			q, uArgs,
			bound, boundArgs,
			opts,
			func() (*sql.Rows, error) {
				return db.tx.QueryContext(db.getCtx(), bound, boundArgs...)
			},
		)
	case kindSt:
		rows, err = logDo(
			db.getCtx(),
			q, uArgs,
			bound, boundArgs,
			opts,
			func() (*sql.Rows, error) {
				return db.st.QueryContext(db.getCtx(), boundArgs...)
			},
		)
	default:
		panic("unreachable")
	}
	if err != nil {
		return err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	if dest == nil {
		// do nothing
		return nil
	} else if f, ok := dest.(func(*sql.Rows) error); ok {
		// custom map
		if f == nil {
			return nil
		}
		for rows.Next() {
			if err := f(rows); err != nil {
				return err
			}
		}
		if err := rows.Err(); err != nil {
			return err
		}
		return nil
	} else {
		dv := reflect.ValueOf(dest)
		if dv.IsNil() {
			return nil
		}
		sliceTyp, ok := isSlicePtr(dv.Type())
		if !ok {
			return ErrInvalidDest
		}
		et := sliceTyp.Elem()

		// use mapper to map
		mapper := getMapper(dest, opts)
		src := Src{Rows: rows}
		if err := src.fetchColumns(false); err != nil {
			return err
		}
		mapRow := makeMapRowFunc(et, mapper, opts)
		if ok := isStruct(et); ok {
			elemDest := reflect.New(et).Interface()
			for rows.Next() {
				if err := mapRow(&src, elemDest); err != nil {
					return err
				}
				dv.Elem().Set(reflect.Append(dv.Elem(), reflect.ValueOf(elemDest).Elem()))
			}
			return rows.Err()
		} else if et1, ok := isStructPtr(et); ok {
			for rows.Next() {
				elemDest := reflect.New(et1).Interface()
				if err := mapRow(&src, elemDest); err != nil {
					return err
				}
				dv.Elem().Set(reflect.Append(dv.Elem(), reflect.ValueOf(elemDest)))
			}
			return rows.Err()
		} else if ok := isRowMap(et); ok {
			for rows.Next() {
				elemDest := reflect.MakeMap(et).Interface()
				if err := mapRow(&src, elemDest); err != nil {
					return err
				}
				dv.Elem().Set(reflect.Append(dv.Elem(), reflect.ValueOf(elemDest)))
			}
			return rows.Err()
		} else if ok := isPrimitive(et); ok {
			for rows.Next() {
				elemDest := reflect.New(et).Interface()
				if err := mapRow(&src, elemDest); err != nil {
					return err
				}
				dv.Elem().Set(reflect.Append(dv.Elem(), reflect.ValueOf(elemDest).Elem()))
			}
			return rows.Err()
		} else if ok := isAny(et); ok {
			for rows.Next() {
				elemDest := map[string]any{}
				if err := mapRow(&src, elemDest); err != nil {
					return err
				}
				dv.Elem().Set(reflect.Append(dv.Elem(), reflect.ValueOf(elemDest)))
			}
			return rows.Err()
		} else {
			return ErrInvalidDest
		}
	}
	return nil
}

func getMapper(dest any, opts *Options) Mapper {
	var mapper Mapper
	if m, ok := dest.(Mapper); ok {
		mapper = m
	} else if m, ok := dest.(func(*Src, *Dest, *Options) error); ok {
		mapper = m
	} else if m, ok := dest.(func(*Src, *Dest) error); ok {
		mapper = func(src *Src, dest *Dest, opts *Options) error {
			return m(src, dest)
		}
	} else {
		mapper = opts.Mapper
	}
	if mapper == nil {
		mapper = defaultMapper
	}
	return mapper
}

func makeMapRowFunc(dt reflect.Type, mapper Mapper, opts *Options) func(src *Src, dest any) error {
	if et, ok := isStructPtr(dt); ok {
		return func(src *Src, dest any) error {
			dv := reflect.ValueOf(dest)
			if dv.IsNil() {
				return nil
			}
			d := Dest{Target: dest, Kind: StructPtr, Type: et}
			return mapper(src, &d, opts)
		}
	} else if et, ok := isStructPtrPtr(dt); ok {
		return func(src *Src, dest any) error {
			dv := reflect.ValueOf(dest)
			if dv.IsNil() {
				return nil
			} else {
				if dv.Elem().IsNil() {
					dv.Elem().Set(reflect.New(et))
				}
			}
			d := Dest{Target: dv.Elem().Interface(), Kind: StructPtr, Type: et}
			return mapper(src, &d, opts)
		}
	} else if ok := isRowMap(dt); ok {
		return func(src *Src, dest any) error {
			dv := reflect.ValueOf(dest)
			if dv.IsNil() {
				return nil
			}
			d := Dest{Target: dest, Kind: Map, Type: dt}
			return mapper(src, &d, opts)
		}
	} else if et, ok := isRowMapPtr(dt); ok {
		return func(src *Src, dest any) error {
			dv := reflect.ValueOf(dest)
			if dv.IsNil() {
				return nil
			}
			if dv.Elem().IsNil() {
				dv.Elem().Set(reflect.MakeMap(et))
			}
			d := Dest{Target: dv.Elem().Interface(), Kind: Map, Type: et}
			return mapper(src, &d, opts)
		}
	} else if et, ok := isPrimitivePtr(dt); ok {
		return func(src *Src, dest any) error {
			dv := reflect.ValueOf(dest)
			if dv.IsNil() {
				return nil
			}
			d := Dest{Target: dest, Kind: PrimitivePtr, Type: et}
			return mapper(src, &d, opts)
		}
	} else if et, ok := isSlicePtr(dt); ok {
		return func(src *Src, dest any) error {
			dv := reflect.ValueOf(dest)
			if dv.IsNil() {
				return nil
			}
			d := Dest{Target: dest, Kind: SlicePtr, Type: et}
			return mapper(src, &d, opts)
		}
	} else if ok := isAnyPtr(dt); ok {
		return func(src *Src, dest any) error {
			dv := reflect.ValueOf(dest)
			if dv.IsNil() {
				return nil
			} else {
				dv.Elem().Set(reflect.ValueOf(map[string]any{}))
			}
			d := Dest{Target: dv.Elem(), Kind: Map, Type: dv.Elem().Type()}
			return mapper(src, &d, opts)
		}
	} else {
		return func(src *Src, dest any) error {
			return ErrInvalidDest
		}
	}
}
