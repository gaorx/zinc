package zinc

import (
	"database/sql"
	"errors"
)

var (
	ErrUnsupportedDriver = errors.New("unsupported driver")
	ErrBind              = errors.New("bind error")
	ErrInvalidDest       = errors.New("invalid dest")
	ErrNoRows            = sql.ErrNoRows
	ErrCoerceDest        = errors.New("coerce dest error")
)
