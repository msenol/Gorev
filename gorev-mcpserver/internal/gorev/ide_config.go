package gorev

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// IDEConfig holds IDE management configuration
type IDEConfig struct {
	AutoInstall     bool          `json:"auto_install"`
	AutoUpdate      bool          `json:"auto_update"`
	CheckInterval   time.Duration `json:"check_interval"`
	SupportedIDEs   []string      `json:"supported_ides"`
	ExtensionID     string        `json:"extension_id"`
	DisablePrompts  bool          `json:"disable_prompts"`
	LastUpdateCheck time.Time     `json:"last_update_check"`
}

// DefaultIDEConfig returns the default configuration
func DefaultIDEConfig() *IDEConfig {
	return &IDEConfig{
		AutoInstall:     false,          // Default: don't auto-install
		AutoUpdate:      false,          // Default: don't auto-update
		CheckInterval:   24 * time.Hour, // Check daily
		SupportedIDEs:   []string{"vscode", "cursor", "windsurf"},
		ExtensionID:     "mehmetsenol.gorev-vscode",
		DisablePrompts:  false, // Show prompts by default
		LastUpdateCheck: time.Time{},
	}
}

// IDEConfigManager manages IDE configuration
type IDEConfigManager struct {
	configPath string
	config     *IDEConfig
}

// NewIDEConfigManager creates a new configuration manager
func NewIDEConfigManager() *IDEConfigManager {
	homeDir, _ := os.UserHomeDir()
	configPath := filepath.Join(homeDir, ".gorev", "ide-config.json")

	return &IDEConfigManager{
		configPath: configPath,
		config:     DefaultIDEConfig(),
	}
}

// LoadConfig loads configuration from file
func (cm *IDEConfigManager) LoadConfig() error {
	// If config file doesn't exist, use default config
	if _, err := os.Stat(cm.configPath); os.IsNotExist(err) {
		return cm.SaveConfig() // Create default config file
	}

	data, err := os.ReadFile(cm.configPath)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, cm.config)
}

// SaveConfig saves configuration to file
func (cm *IDEConfigManager) SaveConfig() error {
	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(cm.configPath), 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cm.config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(cm.configPath, data, 0644)
}

// GetConfig returns the current configuration
func (cm *IDEConfigManager) GetConfig() *IDEConfig {
	return cm.config
}

// SetAutoInstall sets auto-install option
func (cm *IDEConfigManager) SetAutoInstall(enable bool) error {
	cm.config.AutoInstall = enable
	return cm.SaveConfig()
}

// SetAutoUpdate sets auto-update option
func (cm *IDEConfigManager) SetAutoUpdate(enable bool) error {
	cm.config.AutoUpdate = enable
	return cm.SaveConfig()
}

// SetDisablePrompts sets prompt disable option
func (cm *IDEConfigManager) SetDisablePrompts(disable bool) error {
	cm.config.DisablePrompts = disable
	return cm.SaveConfig()
}

// SetCheckInterval sets the update check interval
func (cm *IDEConfigManager) SetCheckInterval(interval time.Duration) error {
	cm.config.CheckInterval = interval
	return cm.SaveConfig()
}

// UpdateLastCheckTime updates the last check timestamp
func (cm *IDEConfigManager) UpdateLastCheckTime() error {
	cm.config.LastUpdateCheck = time.Now()
	return cm.SaveConfig()
}

// ShouldCheckForUpdates returns true if it's time to check for updates
func (cm *IDEConfigManager) ShouldCheckForUpdates() bool {
	if cm.config.LastUpdateCheck.IsZero() {
		return true // Never checked before
	}
	return time.Since(cm.config.LastUpdateCheck) >= cm.config.CheckInterval
}

// GetConfigPath returns the configuration file path
func (cm *IDEConfigManager) GetConfigPath() string {
	return cm.configPath
}
