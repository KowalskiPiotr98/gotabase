package migrations

import (
	"errors"
	"fmt"
	"github.com/KowalskiPiotr98/gotabase/logger"
	"slices"
	"strconv"
	"strings"
)

var (
	migrationNotFound = errors.New("migration with this id was not found")
)

func getLatestAvailableMigration(migrations MigrationFileProvider) (int, error) {
	available, err := getAvailableMigrations(migrations)
	if err != nil {
		return 0, nil
	}
	return slices.Max(*available), nil
}

func getMigrationSql(migrations MigrationFileProvider, i int) (string, error) {
	fileBytes, err := migrations.ReadFile(fmt.Sprintf("sql/%d.sql", i))
	if err != nil {
		logger.LogWarn("Unable to read migration file: %v", err)
		return "", migrationNotFound
	}
	return string(fileBytes), nil
}

func getAvailableMigrations(migrations MigrationFileProvider) (*[]int, error) {
	dirContents, err := migrations.ReadDir("sql")
	if err != nil {
		logger.LogWarn("Unable to read migration directory: %v", err)
		return nil, migrationNotFound
	}

	available := make([]int, 0)
	for _, file := range dirContents {
		nameString := strings.TrimPrefix(strings.TrimSuffix(file.Name(), ".sql"), "sql/")
		nameInt, err := strconv.Atoi(nameString)
		if err != nil {
			logger.LogWarn("Unable to parse migration file name: %v", err)
			return nil, err
		}
		available = append(available, nameInt)
	}
	return &available, nil
}
