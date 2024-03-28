package zinc

import (
	"database/sql"
	"errors"
	"reflect"
)

type mysqlDialect struct{}

func (d mysqlDialect) DriverName() string {
	return "mysql"
}

func (d mysqlDialect) Quote(s string, _ *Options) string {
	return quote(s, quoteBack)
}

func (d mysqlDialect) CompileNamedQuery(q string, _ *Options) (string, []string, error) {
	return compileNamedQuery([]byte(q), bindQuestion)
}

func (d mysqlDialect) NewDest(ci *sql.ColumnType, _ *Options) any {
	return reflect.New(ci.ScanType()).Interface()
}

func (d mysqlDialect) CoerceDest(ci *sql.ColumnType, scannedVal reflect.Value, toType reflect.Type, opts *Options) (reflect.Value, error) {
	if toType == typAny {
		if sliceContains(mysqlTextTypes, ci.DatabaseTypeName()) {
			switch a := scannedVal.Interface().(type) {
			case []byte:
				s, ok := b2s(a, opts.TextCharset)
				if !ok {
					return reflect.Value{}, errors.New("failed to convert []byte to string")
				}
				return reflect.ValueOf(s), nil
			case sql.RawBytes:
				s, ok := b2s(a, opts.TextCharset)
				if !ok {
					return reflect.Value{}, errors.New("failed to convert sql.RawBytes to string")
				}
				return reflect.ValueOf(s), nil
			}
		}
	}
	return coerceDest(ci, scannedVal, toType, opts)
}

var mysqlTextTypes = []string{
	"CHAR",
	"VARCHAR",
	"TEXT",
}
