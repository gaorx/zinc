package zinc

import (
	"database/sql"
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
	switch ci.DatabaseTypeName() {
	case "CHAR", "VARCHAR":
		return new(string)
	// TODO: more mysql types
	default:
		return reflect.New(ci.ScanType()).Interface()
	}
}

func (d mysqlDialect) CoerceDest(scannedVal reflect.Value, toType reflect.Type, opts *Options) (reflect.Value, error) {
	return coerceDest(scannedVal, toType, opts)
}
