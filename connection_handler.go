package gotabase

import (
	"database/sql"
	"errors"
	"github.com/KowalskiPiotr98/gotabase/logger"
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

	logger.LogInfo("Initialising database connection...")
	database, err := sql.Open(driver, connectionString)
	if err != nil {
		logger.LogWarn("Failed to open database connection: %v", err)
		return err
	}
	err = database.Ping()
	if err != nil {
		logger.LogWarn("Failed to ping database: %v", err)
		database.Close()
		return err
	}

	logger.LogInfo("Database connection established")
	connection = &connectionHandler{
		database: database,
	}
	return nil
}

func GetConnection() Connector {
	if connection == nil {
		logger.LogPanic(connectionNotInitialisedErr.Error())
	}

	return connection
}

func BeginTransaction() (*Transaction, error) {
	if connection == nil {
		logger.LogPanic(connectionNotInitialisedErr.Error())
	}

	tx, err := connection.database.Begin()
	if err != nil {
		logger.LogWarn("Failed to begin transaction: %v", err)
		return nil, err
	}
	return newTransaction(tx), nil
}

func CloseConnection() error {
	if connection == nil {
		return nil
	}

	if err := connection.database.Close(); err != nil {
		logger.LogWarn("Failed to close database connection: %v", err)
		return err
	}
	connection = nil
	return nil
}
