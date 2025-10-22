package operations

import (
	"errors"
	"github.com/KowalskiPiotr98/gotabase/logger"
	"github.com/lib/pq"
	"strings"
)

var Errors errorConfig

type errorConfig struct {
	// DataNotFoundErr indicates that the requested data couldn't be found.
	DataNotFoundErr error
	// DataAlreadyExistsErr indicates that the requested data modification or creation conflicts with already existing data.
	DataAlreadyExistErr error
	// DataUserErr indicates that data is being used in another row (as a foreign key) and cannot be modified or removed.
	DataUsedErr error
	// RowNumberUnexpectedErr indicates that an unexpected number of rows was affected and you should probably consider aborting the transation.
	RowNumberUnexpectedErr error

	handlers []func(err error) error
}

func init() {
	Errors = errorConfig{
		DataNotFoundErr:        errors.New("requested data was not found in the database"),
		DataAlreadyExistErr:    errors.New("this data already exists in the database"),
		DataUsedErr:            errors.New("this data is being used in another object and cannot be removed"),
		RowNumberUnexpectedErr: errors.New("unexpected number of rows affected"),
		handlers:               make([]func(err error) error, 0),
	}
}

// RegisterHandler adds a new handler function to the last place in the handlers array.
// The handler function should either return an error to indicate the handling of the passed error, or nil if the database error is not recognised.
func (config *errorConfig) RegisterHandler(handler func(err error) error) {
	config.handlers = append(config.handlers, handler)
}

// HandleError will look through the list of registered error handlers in order of registration.
// If one handler is capable of converting the error, then that converted error will be returned.
// Otherwise, the original error will be echoed back and warning logged.
func (config *errorConfig) HandleError(err error) error {
	if err == nil {
		return nil
	}

	for _, handler := range config.handlers {
		handled := handler(err)
		if handled != nil {
			return handled
		}
	}

	logger.LogWarn("Unknown database error occurred: %v", err)
	return err
}

func (config *errorConfig) RegisterDefaultPostgresHandlers() {
	config.RegisterHandler(func(err error) error {
		// this handles scans of 0 rows
		if err.Error() == "sql: no rows in result set" {
			return config.DataNotFoundErr
		}
		// this handles foreign keys errors
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && (pgErr.Code == "23503" || pgErr.Code == "23001") {
			// needs check to see what kind of FK error is this
			if pgErr.Constraint != "" && strings.Contains(pgErr.Message, "update or delete") {
				return config.DataUsedErr
			}
			return config.DataNotFoundErr
		}
		return nil
	})
	config.RegisterHandler(func(err error) error {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return config.DataAlreadyExistErr
		}
		return nil
	})
}
