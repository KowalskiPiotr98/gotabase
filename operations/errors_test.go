package operations

import (
	"github.com/KowalskiPiotr98/gotabase"
	"github.com/KowalskiPiotr98/gotabase/internal/tests"
	"github.com/stretchr/testify/assert"
	"testing"
)

func makeTestTable(connector gotabase.Connector) {
	_, err := connector.Exec("create table test (id integer primary key); create table test2 (test_id integer references test(id) on delete restrict);")
	tests.PanicOnErr(err)
}

func setPostgresHandlers() {
	Errors.RegisterDefaultPostgresHandlers()
}

func TestErrorConfig_HandleError(t *testing.T) {
	setPostgresHandlers()

	t.Run("Default postgres handlers", func(t *testing.T) {
		t.Run("Data not found error during select", func(t *testing.T) {
			db := tests.GetDatabaseWithCleanup(t)
			makeTestTable(db)
			query := `select * from test where id = 0`
			row, err := db.QueryRow(query)
			tests.PanicOnErr(err)
			var testInt int
			err = row.Scan(&testInt)
			assert.Equal(t, Errors.DataNotFoundErr, Errors.HandleError(err))
		})
		t.Run("Data not found during foreign key insert", func(t *testing.T) {
			db := tests.GetDatabaseWithCleanup(t)
			makeTestTable(db)
			_, err := db.Exec("insert into test2 (test_id) values (12)")
			assert.Equal(t, Errors.DataNotFoundErr, Errors.HandleError(err))
		})
		t.Run("Data already exists error", func(t *testing.T) {
			db := tests.GetDatabaseWithCleanup(t)
			makeTestTable(db)
			_, err := db.Exec("insert into test (id) values (12)")
			tests.PanicOnErr(err)
			_, err = db.Exec("insert into test (id) values (12)")
			assert.Equal(t, Errors.DataAlreadyExistErr, Errors.HandleError(err))
		})
		t.Run("Data used in FK error", func(t *testing.T) {
			db := tests.GetDatabaseWithCleanup(t)
			makeTestTable(db)
			_, err := db.Exec("insert into test (id) values (12)")
			tests.PanicOnErr(err)
			_, err = db.Exec("insert into test2 (test_id) values (12)")
			tests.PanicOnErr(err)
			_, err = db.Exec("delete from test where id = 12")
			assert.Equal(t, Errors.DataUsedErr, Errors.HandleError(err))
		})
	})
}
