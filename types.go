package zinc

import (
	"database/sql"
	"reflect"
	"time"
)

var (
	typBool     = reflect.TypeOf(false)
	typInt      = reflect.TypeOf(0)
	typInt8     = reflect.TypeOf(int8(0))
	typInt16    = reflect.TypeOf(int16(0))
	typInt32    = reflect.TypeOf(int32(0))
	typInt64    = reflect.TypeOf(int64(0))
	typUint     = reflect.TypeOf(uint(0))
	typUint8    = reflect.TypeOf(uint8(0))
	typUint16   = reflect.TypeOf(uint16(0))
	typUint32   = reflect.TypeOf(uint32(0))
	typUint64   = reflect.TypeOf(uint64(0))
	typFloat32  = reflect.TypeOf(float32(0.0))
	typFloat64  = reflect.TypeOf(0.0)
	typString   = reflect.TypeOf("")
	typBytes    = reflect.TypeOf([]byte{})
	typRawBytes = reflect.TypeOf(sql.RawBytes{})
	typAny      = reflect.TypeOf((*any)(nil)).Elem()
	typTime     = reflect.TypeOf(time.Time{})
)
