package migrations

import "io/fs"

// MigrationFileProvider is an interface used for .sql migration file retrieval.
// See readme for usage.
type MigrationFileProvider interface {
	ReadDir(name string) ([]fs.DirEntry, error)
	ReadFile(name string) ([]byte, error)
}
