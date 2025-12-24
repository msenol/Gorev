package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// ConnectionMode defines how clients connect to the daemon
type ConnectionMode string

const (
	// ModeAuto automatically detects the best connection method
	ModeAuto ConnectionMode = "auto"

	// ModeConnectionLocal uses local gorev binary
	ModeConnectionLocal ConnectionMode = "local"

	// ModeDocker uses Docker container
	ModeDocker ConnectionMode = "docker"

	// ModeRemote connects to an already running remote daemon
	ModeRemote ConnectionMode = "remote"
)

// SharedConfig holds configuration shared between all gorev components
// This is stored in ~/.gorev/config.json and used by VS Code extension,
// MCP proxy, and CLI tools
type SharedConfig struct {
	// ConnectionMode is the preferred connection mode
	ConnectionMode ConnectionMode `json:"connection_mode"`

	// LocalServerPath is the path to local gorev binary (for local mode)
	LocalServerPath string `json:"local_server_path,omitempty"`

	// DockerComposeFile is the path to docker-compose.yml (for docker mode)
	DockerComposeFile string `json:"docker_compose_file,omitempty"`

	// ServerPort is the default server port
	ServerPort string `json:"server_port,omitempty"`

	// LastUpdated is the timestamp of last update
	LastUpdated int64 `json:"last_updated"`
}

var (
	sharedConfig     *SharedConfig
	sharedConfigOnce sync.Once
	sharedConfigPath string
)

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		sharedConfigPath = "/tmp/.gorev/config.json"
	} else {
		sharedConfigPath = filepath.Join(home, ".gorev", "config.json")
	}
}

// GetSharedConfig returns the shared configuration (lazy loaded)
func GetSharedConfig() *SharedConfig {
	sharedConfigOnce.Do(func() {
		sharedConfig = loadSharedConfig()
	})
	return sharedConfig
}

// loadSharedConfig loads shared config from file or returns defaults
func loadSharedConfig() *SharedConfig {
	data, err := os.ReadFile(sharedConfigPath)
	if err != nil {
		if os.IsNotExist(err) {
			return defaultSharedConfig()
		}
		// On read error, return defaults
		return defaultSharedConfig()
	}

	var config SharedConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return defaultSharedConfig()
	}
	return &config
}

// defaultSharedConfig returns the default shared configuration
func defaultSharedConfig() *SharedConfig {
	return &SharedConfig{
		ConnectionMode: ModeAuto,
		ServerPort:     "5082",
	}
}

// SaveSharedConfig saves the shared configuration to file
func SaveSharedConfig(config *SharedConfig) error {
	config.LastUpdated = UnixNow()

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(sharedConfigPath), 0755); err != nil {
		return err
	}

	return os.WriteFile(sharedConfigPath, data, 0644)
}

// GetEffectiveConnectionMode returns the effective connection mode
// Priority: ENV > Shared Config File > Default
func GetEffectiveConnectionMode() ConnectionMode {
	// 1. Check environment variable first (highest priority)
	if envMode := os.Getenv("GOREV_CONNECTION_MODE"); envMode != "" {
		return ConnectionMode(envMode)
	}

	// 2. Use shared config file
	return GetSharedConfig().ConnectionMode
}

// GetEffectiveServerPort returns the effective server port
// Priority: ENV > Shared Config > CLI Flag > Default
func GetEffectiveServerPort() string {
	// 1. Check environment variable
	if port := os.Getenv("GOREV_API_PORT"); port != "" {
		return port
	}

	// 2. Use shared config
	return GetSharedConfig().ServerPort
}

// UnixNow returns current Unix timestamp
func UnixNow() int64 {
	return UnixNowFunc()
}

// UnixNowFunc is a variable for testing
var UnixNowFunc = func() int64 {
	return time.Now().Unix()
}
