package operations

import (
	"github.com/KowalskiPiotr98/gotabase"
	"github.com/KowalskiPiotr98/gotabase/logger"
)

type HasIdSetter[T any] interface {
	SetId(id T)
}

// QueryRows is a helper function to run a multiple rows query based on a query string.
func QueryRows[T any](connector gotabase.Connector, scanner func(row gotabase.Row) (*T, error), query string, args ...any) ([]*T, error) {
	rows, err := connector.QueryRows(query, args...)
	if err != nil {
		return nil, Errors.HandleError(err)
	}
	return scanRows(rows, scanner)
}

// QueryRow is a helper function to run a single row query.
func QueryRow[T any](connector gotabase.Connector, scanner func(row gotabase.Row) (*T, error), query string, args ...any) (*T, error) {
	row, err := connector.QueryRow(query, args...)
	if err != nil {
		return nil, Errors.HandleError(err)
	}

	result, err := scanner(row)
	if err != nil {
		return nil, Errors.HandleError(err)
	}
	return result, nil
}

// CreateRowWithId creates a new data in the database.
// The query is expected to return a new id, that will be set in the object using the HasIdSetter interface method.
func CreateRowWithId[TId any, T HasIdSetter[TId]](connector gotabase.Connector, object T, query string, args ...any) error {
	row, err := connector.QueryRow(query, args...)
	if err != nil {
		return Errors.HandleError(err)
	}

	var id TId
	if err := row.Scan(&id); err != nil {
		logger.LogWarn("Failed to scan new object id: %v", err)
		return Errors.HandleError(err)
	}
	object.SetId(id)

	return nil
}

// CreateRowWithScan creates a new data in the database.
// The query is expected to return a row, that will match the provided scanner function.
func CreateRowWithScan[T any](connector gotabase.Connector, object *T, scanner func(row gotabase.Row, object *T) error, query string, args ...any) error {
	row, err := connector.QueryRow(query, args...)
	if err != nil {
		return Errors.HandleError(err)
	}

	if err = scanner(row, object); err != nil {
		return Errors.HandleError(err)
	}

	return nil
}

// CreateRow creates a new data in the database.
func CreateRow(connector gotabase.Connector, query string, args ...any) error {
	result, err := connector.Exec(query, args...)
	if err != nil {
		return Errors.HandleError(err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return Errors.HandleError(err)
	}

	if affected != 1 {
		return Errors.RowNumberUnexpectedErr
	}

	return nil
}

// UpdateRow runs the query to update a single row in the database.
// If none or more than one row are updated, then error will be returned.
func UpdateRow(connector gotabase.Connector, query string, args ...any) error {
	return runSingleRowAffectedQuery(connector, query, args...)
}

// DeleteRow runs the query to remove a single row from the database.
// If none or more than one row are affected, then error will be returned.
func DeleteRow(connector gotabase.Connector, query string, args ...any) error {
	return runSingleRowAffectedQuery(connector, query, args...)
}

// DeleteRows runs the query to remove at last one row from the database.
// If none or more than one row are affected, then error will be returned.
func DeleteRows(connector gotabase.Connector, query string, args ...any) error {
	result, err := connector.Exec(query, args...)
	if err != nil {
		return Errors.HandleError(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return Errors.HandleError(err)
	}

	if rowsAffected == 0 {
		return Errors.DataNotFoundErr
	}
	return nil
}

// TryDelete runs the query attempting to delete something.
// If the query runs successfully, no error is returned, regardless of the number of affected rows.
func TryDelete(connector gotabase.Connector, query string, args ...any) error {
	_, err := connector.Exec(query, args...)
	if err != nil {
		return Errors.HandleError(err)
	}
	return nil
}

func runSingleRowAffectedQuery(connector gotabase.Connector, query string, args ...any) error {
	result, err := connector.Exec(query, args...)
	if err != nil {
		return Errors.HandleError(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return Errors.HandleError(err)
	}

	if rowsAffected > 1 {
		return Errors.RowNumberUnexpectedErr
	}
	if rowsAffected == 0 {
		return Errors.DataNotFoundErr
	}
	return nil
}

func scanRows[T any](rows gotabase.Rows, scanner func(row gotabase.Row) (*T, error)) ([]*T, error) {
	defer rows.Close()
	result := make([]*T, 0)

	for rows.Next() {
		item, err := scanner(rows)
		if err != nil {
			return nil, Errors.HandleError(err)
		}

		result = append(result, item)
	}

	return result, nil
}
