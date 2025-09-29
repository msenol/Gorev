package api

import (
	"embed"
	"io/fs"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
)

// ServeStaticFiles serves the embedded web UI static files
func ServeStaticFiles(app *fiber.App, webDistFS embed.FS) error {
	// Get the dist subdirectory from embedded FS
	// webDistFS contains "web/dist" structure, we need just "dist"
	distFS, err := fs.Sub(webDistFS, "web/dist")
	if err != nil {
		log.Printf("Error creating sub FS: %v", err)
		return err
	}

	// Serve static files from root path
	// This must be registered AFTER API routes
	app.Use("/", filesystem.New(filesystem.Config{
		Root:   http.FS(distFS),
		Browse: false,
		Index:  "index.html",
	}))

	return nil
}
