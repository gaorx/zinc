package zinc

import (
	"context"
	"database/sql"
)

type DB struct {
	kind    int
	db      *sql.DB
	tx      *sql.Tx
	st      *sql.Stmt
	options *Options
	ctx     context.Context
}

const (
	kindDB = 1
	kindTx = 2
	kindSt = 3
)

func New(driverName string, db *sql.DB, opts *Options) (*DB, error) {
	opts1 := fromPtr(opts)
	if opts1.Dialect == nil {
		dialect := dialectOf(driverName)
		if dialect == nil {
			return nil, ErrUnsupportedDriver
		}
		opts1.Dialect = dialect
	}
	if opts1.NameResolver == nil {
		opts1.NameResolver = DefaultNameResolver
	}
	return &DB{
		kind:    kindDB,
		db:      db,
		tx:      nil,
		st:      nil,
		options: &opts1,
	}, nil
}

func Open(driverName string, dsn string, opts *Options) (*DB, error) {
	db, err := sql.Open(driverName, dsn)
	if err != nil {
		return nil, err
	}
	return New(driverName, db, opts)
}

func (db *DB) Close() error {
	switch db.kind {
	case kindDB:
		return db.db.Close()
	case kindSt:
		return db.st.Close()
	default:
		panic("invalid kind")
	}
}

func (db *DB) DB() *sql.DB {
	return db.db
}

func (db *DB) Dialect() Dialect {
	return db.options.Dialect
}

func (db *DB) WithContext(ctx context.Context) *DB {
	cloned := *db
	cloned.ctx = ctx
	return &cloned
}

func (db *DB) clone(newKind int, newDB *sql.DB, newTx *sql.Tx, newStmt *sql.Stmt) *DB {
	cloned := *db
	cloned.kind = newKind
	cloned.db = newDB
	cloned.tx = newTx
	cloned.st = newStmt
	return &cloned
}

func (db *DB) getCtx() context.Context {
	if db.ctx == nil {
		return context.Background()
	}
	return db.ctx
}
