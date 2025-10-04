package test

import (
	"embed"
	"io/fs"
)

// Embed migrations directory into test binary
// This ensures integration tests can always access migrations
//
//go:embed migrations
var embeddedMigrations embed.FS

// getEmbeddedMigrationsFS returns the embedded migrations filesystem
func getEmbeddedMigrationsFS() (fs.FS, error) {
	migrationsFS, err := fs.Sub(embeddedMigrations, "migrations")
	if err != nil {
		return nil, err
	}
	return migrationsFS, nil
}
