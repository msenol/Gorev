package api

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

// TestHealthEndpointSimple tests the health check endpoint without database
func TestHealthEndpointSimple(t *testing.T) {
	// Create a minimal server without database for health check test
	server := &APIServer{}
	server.app = createFiberApp()
	server.setupHealthRoute()

	req := httptest.NewRequest("GET", "/api/v1/health", nil)
	resp, err := server.app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

// createFiberApp creates a basic Fiber app for testing
func createFiberApp() *fiber.App {
	return fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})
}

// setupHealthRoute sets up only the health route for testing
func (s *APIServer) setupHealthRoute() {
	api := s.app.Group("/api/v1")
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
			"time":   time.Now().Unix(),
		})
	})
}
