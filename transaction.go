package gotabase

import "database/sql"

type Transaction struct {
	tx *sql.Tx
}

func newTransaction(tx *sql.Tx) *Transaction {
	return &Transaction{tx: tx}
}

var _ Connector = (*Transaction)(nil)

func (t *Transaction) QueryRow(sql string, args ...interface{}) (Row, error) {
	row := t.tx.QueryRow(sql, args...)
	return row, nil
}

func (t *Transaction) QueryRows(sql string, args ...interface{}) (Rows, error) {
	return t.tx.Query(sql, args...)
}

func (t *Transaction) Exec(sql string, args ...interface{}) (Result, error) {
	return t.tx.Exec(sql, args...)
}

func (t *Transaction) Commit() error {
	return t.tx.Commit()
}

func (t *Transaction) Rollback() error {
	return t.tx.Rollback()
}
