package config

import (
	"os"
	"sync"
)

// ServerMode defines how the server handles workspace databases
type ServerMode string

const (
	// ModeLocal creates a .gorev/gorev.db in each workspace directory
	// Best for: Local development where each project has its own database
	ModeLocal ServerMode = "local"

	// ModeCentralized uses a single database with workspace_id isolation
	// Best for: Docker, remote servers, multi-user environments
	ModeCentralized ServerMode = "centralized"
)

// ServerConfig holds the server configuration
type ServerConfig struct {
	// Mode determines how workspace databases are managed
	Mode ServerMode

	// CentralizedDBPath is the database path when in centralized mode
	// Default: /data/gorev.db (Docker) or ~/.gorev/gorev.db (local)
	CentralizedDBPath string

	// Port is the HTTP API port
	Port string

	// AllowLocalPaths enables path-based workspace creation (security)
	// When false, only workspace_id based operations are allowed
	AllowLocalPaths bool
}

var (
	globalConfig *ServerConfig
	configMu     sync.RWMutex
)

// DefaultConfig returns the default server configuration
func DefaultConfig() *ServerConfig {
	mode := ModeLocal
	if envMode := os.Getenv("GOREV_MODE"); envMode != "" {
		mode = ServerMode(envMode)
	}

	dbPath := os.Getenv("GOREV_DB_PATH")
	if dbPath == "" {
		if mode == ModeCentralized {
			dbPath = "/data/gorev.db"
		}
		// For local mode, dbPath is determined per-workspace
	}

	port := os.Getenv("GOREV_API_PORT")
	if port == "" {
		port = "5082"
	}

	return &ServerConfig{
		Mode:              mode,
		CentralizedDBPath: dbPath,
		Port:              port,
		AllowLocalPaths:   mode == ModeLocal,
	}
}

// SetGlobalConfig sets the global server configuration
func SetGlobalConfig(cfg *ServerConfig) {
	configMu.Lock()
	defer configMu.Unlock()
	globalConfig = cfg
}

// GetGlobalConfig returns the global server configuration
func GetGlobalConfig() *ServerConfig {
	configMu.RLock()
	defer configMu.RUnlock()
	if globalConfig == nil {
		return DefaultConfig()
	}
	return globalConfig
}

// IsCentralizedMode returns true if server is in centralized mode
func IsCentralizedMode() bool {
	return GetGlobalConfig().Mode == ModeCentralized
}

// IsLocalMode returns true if server is in local mode
func IsLocalMode() bool {
	return GetGlobalConfig().Mode == ModeLocal
}
