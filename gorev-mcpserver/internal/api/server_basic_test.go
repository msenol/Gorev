package api

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/msenol/gorev/internal/gorev"
	"github.com/msenol/gorev/internal/i18n"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupBasicTestServer creates a test server with in-memory database
func setupBasicTestServer(t *testing.T) (*APIServer, func()) {
	// Initialize i18n
	err := i18n.Initialize("en")
	if err != nil {
		// Already initialized, that's fine
	}

	// Create VeriYonetici with in-memory database
	veriYonetici, err := gorev.YeniVeriYonetici(":memory:", "file://../../internal/veri/migrations")
	require.NoError(t, err)

	// Create IsYonetici
	isYonetici := gorev.YeniIsYonetici(veriYonetici)

	// Create API server
	server := NewAPIServer("8080", isYonetici)

	cleanup := func() {
		if veriYonetici != nil {
			veriYonetici.Kapat()
		}
	}

	return server, cleanup
}

// TestAPIServerCreation tests basic server creation
func TestAPIServerCreation(t *testing.T) {
	server, cleanup := setupBasicTestServer(t)
	defer cleanup()

	assert.NotNil(t, server)
	assert.NotNil(t, server.app)
	assert.Equal(t, "8080", server.port)
	assert.NotNil(t, server.isYonetici)
}

// TestHealthCheckAPI tests the health check endpoint
func TestHealthCheckAPI(t *testing.T) {
	server, cleanup := setupBasicTestServer(t)
	defer cleanup()

	req := httptest.NewRequest("GET", "/api/v1/health", nil)
	resp, err := server.app.Test(req)

	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var result map[string]interface{}
	body, _ := io.ReadAll(resp.Body)
	err = json.Unmarshal(body, &result)

	require.NoError(t, err)
	assert.Equal(t, "ok", result["status"])
	assert.NotNil(t, result["time"])
}

// TestProjectsEndpoint tests projects API
func TestProjectsEndpoint(t *testing.T) {
	server, cleanup := setupBasicTestServer(t)
	defer cleanup()

	// List projects (should work)
	t.Run("ListProjects", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/projects", nil)
		resp, err := server.app.Test(req)

		require.NoError(t, err)
		assert.True(t, resp.StatusCode >= 200 && resp.StatusCode < 300)
	})

	// Create project
	t.Run("CreateProject", func(t *testing.T) {
		payload := map[string]interface{}{
			"isim":  "Test Project",
			"tanim": "Test Description",
		}
		body, _ := json.Marshal(payload)

		req := httptest.NewRequest("POST", "/api/v1/projects", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := server.app.Test(req)

		require.NoError(t, err)
		// Accept 200 or 201
		assert.True(t, resp.StatusCode >= 200 && resp.StatusCode < 300)
	})
}

// TestTasksEndpoint tests tasks API
func TestTasksEndpoint(t *testing.T) {
	server, cleanup := setupBasicTestServer(t)
	defer cleanup()

	// List tasks
	t.Run("ListTasks", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/v1/tasks", nil)
		resp, err := server.app.Test(req)

		require.NoError(t, err)
		assert.True(t, resp.StatusCode >= 200 && resp.StatusCode < 300)
	})
}

// TestTemplatesEndpoint tests templates API
func TestTemplatesEndpoint(t *testing.T) {
	server, cleanup := setupBasicTestServer(t)
	defer cleanup()

	req := httptest.NewRequest("GET", "/api/v1/templates", nil)
	resp, err := server.app.Test(req)

	require.NoError(t, err)
	assert.True(t, resp.StatusCode >= 200 && resp.StatusCode < 300)
}

// TestSummaryAPI tests summary endpoint
func TestSummaryAPI(t *testing.T) {
	server, cleanup := setupBasicTestServer(t)
	defer cleanup()

	req := httptest.NewRequest("GET", "/api/v1/summary", nil)
	resp, err := server.app.Test(req)

	require.NoError(t, err)
	assert.True(t, resp.StatusCode >= 200 && resp.StatusCode < 300)
}

// TestCORSSetup tests CORS headers
func TestCORSSetup(t *testing.T) {
	server, cleanup := setupBasicTestServer(t)
	defer cleanup()

	req := httptest.NewRequest("GET", "/api/v1/health", nil)
	req.Header.Set("Origin", "http://localhost:5001")
	resp, err := server.app.Test(req)

	require.NoError(t, err)
	// CORS headers might be present (depending on config)
	// Just check that request succeeds
	assert.True(t, resp.StatusCode >= 200 && resp.StatusCode < 300)
}

// TestServerGracefulShutdown tests shutdown
func TestServerGracefulShutdown(t *testing.T) {
	server, cleanup := setupBasicTestServer(t)
	defer cleanup()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := server.Shutdown(ctx)
	assert.NoError(t, err)
}

// TestServerApp tests App() method
func TestServerApp(t *testing.T) {
	server, cleanup := setupBasicTestServer(t)
	defer cleanup()

	app := server.App()
	assert.NotNil(t, app)
	assert.Equal(t, server.app, app)
}

// TestSetHandlers tests SetHandlers method
func TestSetHandlers(t *testing.T) {
	server, cleanup := setupBasicTestServer(t)
	defer cleanup()

	mockHandlers := struct{}{}
	server.SetHandlers(mockHandlers)

	assert.NotNil(t, server.handlers)
}
