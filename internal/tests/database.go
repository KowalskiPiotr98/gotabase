package tests

import (
	"fmt"
	"github.com/KowalskiPiotr98/gotabase"
	_ "github.com/lib/pq"
	"math/rand"
	"os"
	"strconv"
	"testing"
)

func GetDatabaseWithCleanup(t *testing.T) gotabase.Connector {
	db, name := GetDatabase()
	t.Cleanup(func() { DropDatabase(name) })
	return db
}

func GetDatabase() (gotabase.Connector, string) {
	dbName := strconv.Itoa(rand.Int())
	baseConnectionString := getBaseConnectionString()
	err := gotabase.InitialiseConnection(baseConnectionString+dbName, "postgres")
	if err != nil {
		PanicOnErr(gotabase.InitialiseConnection(baseConnectionString+"postgres", "postgres"))
		_, err = gotabase.GetConnection().Exec(fmt.Sprintf("create database \"%s\"", dbName))
		PanicOnErr(err)
		PanicOnErr(gotabase.CloseConnection())
		PanicOnErr(gotabase.InitialiseConnection(baseConnectionString+dbName, "postgres"))
	}
	return gotabase.GetConnection(), dbName
}

func DropDatabase(dbName string) {
	PanicOnErr(gotabase.CloseConnection())
	PanicOnErr(gotabase.InitialiseConnection(getBaseConnectionString()+"postgres", "postgres"))
	_, err := gotabase.GetConnection().Exec(fmt.Sprintf("drop database \"%s\"", dbName))
	PanicOnErr(err)
	PanicOnErr(gotabase.CloseConnection())
}

func getBaseConnectionString() string {
	baseConnectionString := os.Getenv("TEST_POSTGRES")
	if baseConnectionString == "" {
		baseConnectionString = "user=postgres password=postgres sslmode=disable dbname="
	}
	return baseConnectionString
}
