package gotabase

import (
	"database/sql"
	"errors"
)

var connection *connectionHandler

type connectionHandler struct {
	database *sql.DB
}

var _ Connector = (*connectionHandler)(nil)

var (
	connectionNotInitialisedErr  = errors.New("database connection has not been initialised")
	connectionAlreadyInitialised = errors.New("database connection has already been set")
)

func (c *connectionHandler) QueryRow(sql string, args ...interface{}) (Row, error) {
	if connection == nil {
		return nil, connectionNotInitialisedErr
	}

	result := c.database.QueryRow(sql, args...)
	return result, result.Err()
}

func (c *connectionHandler) QueryRows(sql string, args ...interface{}) (Rows, error) {
	if connection == nil {
		return nil, connectionNotInitialisedErr
	}

	return c.database.Query(sql, args...)
}

func (c *connectionHandler) Exec(sql string, args ...interface{}) (Result, error) {
	if connection == nil {
		return nil, connectionNotInitialisedErr
	}

	return c.database.Exec(sql, args...)
}

func InitialiseConnection(connectionString string, driver string) error {
	if connection != nil {
		return connectionAlreadyInitialised
	}

	database, err := sql.Open(driver, connectionString)
	if err != nil {
		return err
	}

	connection = &connectionHandler{
		database: database,
	}
	return nil
}

func GetConnection() Connector {
	if connection == nil {
		panic("database connection not initialised")
	}

	return connection
}

func CloseConnection() error {
	if connection == nil {
		return nil
	}

	if err := connection.database.Close(); err != nil {
		return err
	}
	connection = nil
	return nil
}
