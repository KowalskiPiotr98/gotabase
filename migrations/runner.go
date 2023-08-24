package migrations

import (
	"fmt"
	database "github.com/KowalskiPiotr98/gotabase"
	"strings"
)

var (
	MigrationCreator = func(migrationBodySql string, currentMigration int) string {
		return fmt.Sprintf("begin transaction;\n"+
			"%s\n"+
			"insert into migrations (id) values (%d);\n"+
			"commit;",
			migrationBodySql,
			currentMigration)
	}
	LatestMigrationSelectorSql = "select id from migrations order by id desc limit 1"
	IsInitialMigrationError    = func(err error) bool {
		return strings.HasPrefix(err.Error(), "pq: relation") && strings.HasSuffix(err.Error(), "does not exist")
	}
)

func Migrate(connector database.Connector, fileProvider MigrationFileProvider) error {
	latestApplied, err := getLatestAppliedMigration(connector)
	if err != nil {
		if !IsInitialMigrationError(err) {
			return err
		}
		latestApplied = -1
	}

	latestAvailable, err := getLatestAvailableMigration(fileProvider)
	if latestApplied == latestAvailable {
		return nil
	}

	currentMigration := latestApplied + 1
	for currentMigration <= latestAvailable {
		migrationSql, err := getMigrationSql(fileProvider, currentMigration)
		if err != nil {
			return err
		}

		_, err = connector.Exec(MigrationCreator(migrationSql, currentMigration))
		if err != nil {
			return err
		}
	}
	return nil
}

func getLatestAppliedMigration(connector database.Connector) (int, error) {
	result, err := connector.QueryRow(LatestMigrationSelectorSql)
	if err != nil {
		return 0, err
	}
	var latest int
	if err = result.Scan(&latest); err != nil {
		return 0, err
	}
	return latest, nil
}
