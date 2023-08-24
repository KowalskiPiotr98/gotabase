package gotabase

// Connector provides an interface for database operations that can be performed during application operation.
// Returned types such as Row, Rows, and Result provide an abstraction layer over database/sql types, but are compatible with them.
type Connector interface {
	QueryRow(sql string, args ...interface{}) (Row, error)
	QueryRows(sql string, args ...interface{}) (Rows, error)
	Exec(sql string, args ...interface{}) (Result, error)
}

type Row interface {
	Scan(dest ...any) error
}

type Rows interface {
	Close() error
	Scan(dest ...any) error
	Next() bool
}

type Result interface {
	LastInsertId() (int64, error)
	RowsAffected() (int64, error)
}
