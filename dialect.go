package zinc

import (
	"database/sql"
	"errors"
	"reflect"
	"strconv"
	"strings"
	"unicode"
)

type Dialect interface {
	DriverName() string
	Quote(s string, opts *Options) string
	CompileNamedQuery(q string, opts *Options) (string, []string, error)
	NewDest(ci *sql.ColumnType, opts *Options) any
	CoerceDest(ci *sql.ColumnType, scannedVal reflect.Value, toType reflect.Type, opts *Options) (reflect.Value, error)
}

func dialectOf(driverName string) Dialect {
	switch driverName {
	case "mysql":
		return mysqlDialect{}
	default:
		return nil
	}
}

const (
	bindUnknown = iota
	bindQuestion
	bindDollar
	bindNamed
	bindAt
)

const (
	quoteUnknown = iota
	quoteSingle
	quoteDouble
	quoteBack
)

var allowedBindRunes = []*unicode.RangeTable{unicode.Letter, unicode.Digit}

func quote(s string, quoteType int) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return ""
	}
	var b byte
	switch quoteType {
	case quoteSingle, quoteUnknown:
		b = '\''
	case quoteDouble:
		b = '"'
	case quoteBack:
		b = '`'
	default:
		panic("unhandled quoteType")
	}

	if len(s) >= 2 && s[0] == b && s[len(s)-1] == b {
		return s
	} else {
		return string(b) + s + string(b)
	}
}

func compileNamedQuery(qs []byte, bindType int) (query string, names []string, err error) {
	// 下面的代码来自sqlx
	names = make([]string, 0, 10)
	rebound := make([]byte, 0, len(qs))

	inName := false
	last := len(qs) - 1
	currentVar := 1
	name := make([]byte, 0, 10)

	for i, b := range qs {
		// a ':' while we're in a name is an error
		if b == ':' {
			// if this is the second ':' in a '::' escape sequence, append a ':'
			if inName && i > 0 && qs[i-1] == ':' {
				rebound = append(rebound, ':')
				inName = false
				continue
			} else if inName {
				err = errors.New("unexpected `:` while reading named param at " + strconv.Itoa(i))
				return query, names, err
			}
			inName = true
			name = []byte{}
		} else if inName && i > 0 && b == '=' && len(name) == 0 {
			rebound = append(rebound, ':', '=')
			inName = false
			continue
			// if we're in a name, and this is an allowed character, continue
		} else if inName && (unicode.IsOneOf(allowedBindRunes, rune(b)) || b == '_' || b == '.') && i != last {
			// append the byte to the name if we are in a name and not on the last byte
			name = append(name, b)
			// if we're in a name and it's not an allowed character, the name is done
		} else if inName {
			inName = false
			// if this is the final byte of the string and it is part of the name, then
			// make sure to add it to the name
			if i == last && unicode.IsOneOf(allowedBindRunes, rune(b)) {
				name = append(name, b)
			}
			// add the string representation to the names list
			names = append(names, string(name))
			// add a proper bindvar for the bindType
			switch bindType {
			// oracle only supports named type bind vars even for positional
			case bindNamed:
				rebound = append(rebound, ':')
				rebound = append(rebound, name...)
			case bindQuestion, bindUnknown:
				rebound = append(rebound, '?')
			case bindDollar:
				rebound = append(rebound, '$')
				for _, b := range strconv.Itoa(currentVar) {
					rebound = append(rebound, byte(b))
				}
				currentVar++
			case bindAt:
				rebound = append(rebound, '@', 'p')
				for _, b := range strconv.Itoa(currentVar) {
					rebound = append(rebound, byte(b))
				}
				currentVar++
			default:
				panic("unhandled bindType")
			}
			// add this byte to string unless it was not part of the name
			if i != last {
				rebound = append(rebound, b)
			} else if !unicode.IsOneOf(allowedBindRunes, rune(b)) {
				rebound = append(rebound, b)
			}
		} else {
			// this is a normal byte and should just go onto the rebound query
			rebound = append(rebound, b)
		}
	}

	return string(rebound), names, err
}

// coerce dest
func coerceDest(_ *sql.ColumnType, scannedVal reflect.Value, toType reflect.Type, opts *Options) (reflect.Value, error) {
	if scannedVal.Type() == toType {
		return scannedVal, nil
	}

	// TODO: more coerce
	switch toType {
	case typAny:
		return scannedVal, nil
	case typString:
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
		default:
			panic("xxx")
		}
	default:
		panic("yyyy")
	}
}
