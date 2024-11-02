package operations

import (
	"github.com/KowalskiPiotr98/gotabase"
	"github.com/KowalskiPiotr98/gotabase/internal/tests"
	"github.com/stretchr/testify/assert"
	"testing"
)

type testData struct {
	Id int
}

func (t *testData) SetId(id int) {
	t.Id = id
}

func scanTest(row gotabase.Row) (*testData, error) {
	var item testData
	if err := row.Scan(&item.Id); err != nil {
		return nil, err
	}
	return &item, nil
}

func TestQueryRows(t *testing.T) {
	t.Run("Rows queried and returned", func(t *testing.T) {
		db := tests.GetDatabaseWithCleanup(t)
		makeTestTable(db)
		_, err := db.Exec("insert into test(id) values (12)")
		tests.PanicOnErr(err)
		rows, err := QueryRows(db, scanTest, "select id from test")
		assert.NoError(t, err)
		assert.Len(t, rows, 1)
	})
}

func TestQueryRow(t *testing.T) {
	t.Run("Single row returned", func(t *testing.T) {
		db := tests.GetDatabaseWithCleanup(t)
		makeTestTable(db)
		_, err := db.Exec("insert into test(id) values (12)")
		tests.PanicOnErr(err)
		data, err := QueryRow(db, scanTest, "select id from test")
		assert.NoError(t, err)
		assert.Equal(t, 12, data.Id)
	})
	t.Run("Row not found, error returned", func(t *testing.T) {
		setPostgresHandlers()
		db := tests.GetDatabaseWithCleanup(t)
		makeTestTable(db)
		_, err := QueryRow(db, scanTest, "select id from test")
		assert.Equal(t, Errors.DataNotFoundErr, err)
	})
}

func TestCreateRowWithId(t *testing.T) {
	t.Run("New object created with id read", func(t *testing.T) {
		db := tests.GetDatabaseWithCleanup(t)
		makeTestTable(db)
		testDataObject := &testData{}
		err := CreateRowWithId(db, testDataObject, "insert into test (id) values (12) returning id")
		assert.NoError(t, err)
		assert.Equal(t, 12, testDataObject.Id)
	})
	t.Run("Duplicate object, error returned", func(t *testing.T) {
		setPostgresHandlers()
		db := tests.GetDatabaseWithCleanup(t)
		makeTestTable(db)
		_, err := db.Exec("insert into test(id) values (12)")
		tests.PanicOnErr(err)
		testDataObject := &testData{}
		err = CreateRowWithId(db, testDataObject, "insert into test (id) values (12) returning id")
		assert.Equal(t, Errors.DataAlreadyExistErr, err)
	})
}

func TestDeleteRow(t *testing.T) {
	t.Run("Row deleted", func(t *testing.T) {
		setPostgresHandlers()
		db := tests.GetDatabaseWithCleanup(t)
		makeTestTable(db)
		_, err := db.Exec("insert into test(id) values (12)")
		tests.PanicOnErr(err)
		err = DeleteRow(db, "delete from test where id = 12")
		assert.NoError(t, err)
	})
	t.Run("Row not found, error returned", func(t *testing.T) {
		setPostgresHandlers()
		db := tests.GetDatabaseWithCleanup(t)
		makeTestTable(db)
		err := DeleteRow(db, "delete from test where id = 12")
		assert.Equal(t, Errors.DataNotFoundErr, err)
	})
}
