package main

import (
	"embed"
	"io/fs"
)

// Embed migrations directory into the binary
// This allows the NPX package to work without external migration files
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
