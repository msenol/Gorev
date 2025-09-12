package gorev

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/msenol/gorev/internal/i18n"
)

// IDEType represents the type of IDE
type IDEType string

const (
	IDETypeVSCode   IDEType = "vscode"
	IDETypeCursor   IDEType = "cursor"
	IDETypeWindsurf IDEType = "windsurf"
)

// IDEInfo holds information about a detected IDE
type IDEInfo struct {
	Type           IDEType `json:"type"`
	Name           string  `json:"name"`
	ExecutablePath string  `json:"executable_path"`
	ConfigPath     string  `json:"config_path"`
	ExtensionsPath string  `json:"extensions_path"`
	Version        string  `json:"version,omitempty"`
	IsInstalled    bool    `json:"is_installed"`
}

// IDEDetector handles IDE detection and path discovery
type IDEDetector struct {
	mu           sync.RWMutex
	detectedIDEs map[IDEType]*IDEInfo
}

// NewIDEDetector creates a new IDE detector
func NewIDEDetector() *IDEDetector {
	return &IDEDetector{
		detectedIDEs: make(map[IDEType]*IDEInfo),
	}
}

// DetectAllIDEs detects all supported IDEs on the system
func (d *IDEDetector) DetectAllIDEs() (map[IDEType]*IDEInfo, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	
	// Clear previous detections
	d.detectedIDEs = make(map[IDEType]*IDEInfo)

	// Detect each IDE type
	ideTypes := []IDEType{IDETypeVSCode, IDETypeCursor, IDETypeWindsurf}

	for _, ideType := range ideTypes {
		ideInfo, err := d.detectIDE(ideType)
		if err == nil && ideInfo != nil {
			d.detectedIDEs[ideType] = ideInfo
		}
	}

	// Return a copy to avoid sharing the map reference
	result := make(map[IDEType]*IDEInfo)
	for k, v := range d.detectedIDEs {
		result[k] = v
	}

	return result, nil
}

// GetDetectedIDE returns information about a specific IDE if detected
func (d *IDEDetector) GetDetectedIDE(ideType IDEType) (*IDEInfo, bool) {
	d.mu.RLock()
	defer d.mu.RUnlock()
	
	ide, exists := d.detectedIDEs[ideType]
	return ide, exists
}

// GetAllDetectedIDEs returns all detected IDEs
func (d *IDEDetector) GetAllDetectedIDEs() map[IDEType]*IDEInfo {
	d.mu.RLock()
	defer d.mu.RUnlock()
	
	// Return a copy to avoid sharing the map reference
	result := make(map[IDEType]*IDEInfo)
	for k, v := range d.detectedIDEs {
		result[k] = v
	}
	return result
}

// detectIDE detects a specific IDE type
func (d *IDEDetector) detectIDE(ideType IDEType) (*IDEInfo, error) {
	switch ideType {
	case IDETypeVSCode:
		return d.detectVSCode()
	case IDETypeCursor:
		return d.detectCursor()
	case IDETypeWindsurf:
		return d.detectWindsurf()
	default:
		return nil, fmt.Errorf(i18n.T("error.ide.unsupportedType", map[string]interface{}{"Type": string(ideType)}))
	}
}

// detectVSCode detects VS Code installation
func (d *IDEDetector) detectVSCode() (*IDEInfo, error) {
	ide := &IDEInfo{
		Type: IDETypeVSCode,
		Name: "Visual Studio Code",
	}

	// Platform-specific paths
	var executablePaths []string
	var configPaths []string
	var extensionsPaths []string

	switch runtime.GOOS {
	case "windows":
		executablePaths = []string{
			filepath.Join(os.Getenv("LOCALAPPDATA"), "Programs", "Microsoft VS Code", "Code.exe"),
			filepath.Join(os.Getenv("PROGRAMFILES"), "Microsoft VS Code", "Code.exe"),
			filepath.Join(os.Getenv("PROGRAMFILES(X86)"), "Microsoft VS Code", "Code.exe"),
		}
		configPaths = []string{
			filepath.Join(os.Getenv("APPDATA"), "Code", "User"),
		}
		extensionsPaths = []string{
			filepath.Join(os.Getenv("USERPROFILE"), ".vscode", "extensions"),
		}
	case "darwin":
		executablePaths = []string{
			"/Applications/Visual Studio Code.app/Contents/Resources/app/bin/code",
			"/usr/local/bin/code",
		}
		configPaths = []string{
			filepath.Join(os.Getenv("HOME"), "Library", "Application Support", "Code", "User"),
		}
		extensionsPaths = []string{
			filepath.Join(os.Getenv("HOME"), ".vscode", "extensions"),
		}
	case "linux":
		executablePaths = []string{
			"/usr/bin/code",
			"/snap/code/current/usr/share/code/bin/code",
			filepath.Join(os.Getenv("HOME"), ".local", "share", "applications", "code"),
		}
		configPaths = []string{
			filepath.Join(os.Getenv("HOME"), ".config", "Code", "User"),
		}
		extensionsPaths = []string{
			filepath.Join(os.Getenv("HOME"), ".vscode", "extensions"),
		}
	}

	// Find executable
	for _, path := range executablePaths {
		if fileExists(path) {
			ide.ExecutablePath = path
			ide.IsInstalled = true
			break
		}
	}

	// Find config path
	for _, path := range configPaths {
		if dirExists(path) {
			ide.ConfigPath = path
			break
		}
	}

	// Find extensions path
	for _, path := range extensionsPaths {
		if dirExists(path) {
			ide.ExtensionsPath = path
			break
		}
	}

	if !ide.IsInstalled {
		return nil, fmt.Errorf(i18n.T("error.ide.notFound", map[string]interface{}{"IDE": "VS Code"}))
	}

	// Try to get version
	ide.Version = d.getIDEVersion(ide.ExecutablePath, []string{"--version"})

	return ide, nil
}

// detectCursor detects Cursor IDE installation
func (d *IDEDetector) detectCursor() (*IDEInfo, error) {
	ide := &IDEInfo{
		Type: IDETypeCursor,
		Name: "Cursor",
	}

	// Platform-specific paths
	var executablePaths []string
	var configPaths []string
	var extensionsPaths []string

	switch runtime.GOOS {
	case "windows":
		executablePaths = []string{
			filepath.Join(os.Getenv("LOCALAPPDATA"), "Programs", "cursor", "Cursor.exe"),
			filepath.Join(os.Getenv("APPDATA"), "Cursor", "Cursor.exe"),
		}
		configPaths = []string{
			filepath.Join(os.Getenv("APPDATA"), "Cursor", "User"),
		}
		extensionsPaths = []string{
			filepath.Join(os.Getenv("USERPROFILE"), ".cursor", "extensions"),
		}
	case "darwin":
		executablePaths = []string{
			"/Applications/Cursor.app/Contents/Resources/app/bin/cursor",
			"/usr/local/bin/cursor",
		}
		configPaths = []string{
			filepath.Join(os.Getenv("HOME"), "Library", "Application Support", "Cursor", "User"),
		}
		extensionsPaths = []string{
			filepath.Join(os.Getenv("HOME"), ".cursor", "extensions"),
		}
	case "linux":
		executablePaths = []string{
			"/usr/bin/cursor",
			filepath.Join(os.Getenv("HOME"), ".local", "bin", "cursor"),
			"/opt/cursor/cursor",
		}
		configPaths = []string{
			filepath.Join(os.Getenv("HOME"), ".config", "Cursor", "User"),
		}
		extensionsPaths = []string{
			filepath.Join(os.Getenv("HOME"), ".cursor", "extensions"),
		}
	}

	// Find executable
	for _, path := range executablePaths {
		if fileExists(path) {
			ide.ExecutablePath = path
			ide.IsInstalled = true
			break
		}
	}

	// Find config path
	for _, path := range configPaths {
		if dirExists(path) {
			ide.ConfigPath = path
			break
		}
	}

	// Find extensions path
	for _, path := range extensionsPaths {
		if dirExists(path) {
			ide.ExtensionsPath = path
			break
		}
	}

	if !ide.IsInstalled {
		return nil, fmt.Errorf(i18n.T("error.ide.notFound", map[string]interface{}{"IDE": "Cursor"}))
	}

	// Try to get version
	ide.Version = d.getIDEVersion(ide.ExecutablePath, []string{"--version"})

	return ide, nil
}

// detectWindsurf detects Windsurf IDE installation
func (d *IDEDetector) detectWindsurf() (*IDEInfo, error) {
	ide := &IDEInfo{
		Type: IDETypeWindsurf,
		Name: "Windsurf",
	}

	// Platform-specific paths
	var executablePaths []string
	var configPaths []string
	var extensionsPaths []string

	switch runtime.GOOS {
	case "windows":
		executablePaths = []string{
			filepath.Join(os.Getenv("LOCALAPPDATA"), "Programs", "Windsurf", "Windsurf.exe"),
			filepath.Join(os.Getenv("APPDATA"), "Windsurf", "Windsurf.exe"),
		}
		configPaths = []string{
			filepath.Join(os.Getenv("APPDATA"), "Windsurf", "User"),
		}
		extensionsPaths = []string{
			filepath.Join(os.Getenv("USERPROFILE"), ".windsurf", "extensions"),
		}
	case "darwin":
		executablePaths = []string{
			"/Applications/Windsurf.app/Contents/Resources/app/bin/windsurf",
			"/usr/local/bin/windsurf",
		}
		configPaths = []string{
			filepath.Join(os.Getenv("HOME"), "Library", "Application Support", "Windsurf", "User"),
		}
		extensionsPaths = []string{
			filepath.Join(os.Getenv("HOME"), ".windsurf", "extensions"),
		}
	case "linux":
		executablePaths = []string{
			"/usr/bin/windsurf",
			filepath.Join(os.Getenv("HOME"), ".local", "bin", "windsurf"),
			"/opt/windsurf/windsurf",
		}
		configPaths = []string{
			filepath.Join(os.Getenv("HOME"), ".config", "Windsurf", "User"),
		}
		extensionsPaths = []string{
			filepath.Join(os.Getenv("HOME"), ".windsurf", "extensions"),
		}
	}

	// Find executable
	for _, path := range executablePaths {
		if fileExists(path) {
			ide.ExecutablePath = path
			ide.IsInstalled = true
			break
		}
	}

	// Find config path
	for _, path := range configPaths {
		if dirExists(path) {
			ide.ConfigPath = path
			break
		}
	}

	// Find extensions path
	for _, path := range extensionsPaths {
		if dirExists(path) {
			ide.ExtensionsPath = path
			break
		}
	}

	if !ide.IsInstalled {
		return nil, fmt.Errorf(i18n.T("error.ide.notFound", map[string]interface{}{"IDE": "Windsurf"}))
	}

	// Try to get version
	ide.Version = d.getIDEVersion(ide.ExecutablePath, []string{"--version"})

	return ide, nil
}

// getIDEVersion attempts to get IDE version using command line
func (d *IDEDetector) getIDEVersion(executablePath string, args []string) string {
	// Bu fonksiyonu şimdilik basit tutuyoruz
	// Daha sonra exec.Command kullanarak version alabilir
	return "unknown"
}

// fileExists checks if a file exists
func fileExists(path string) bool {
	if path == "" {
		return false
	}
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// dirExists checks if a directory exists
func dirExists(path string) bool {
	if path == "" {
		return false
	}
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

// IsExtensionInstalled checks if the Gorev extension is installed in the specified IDE
func (d *IDEDetector) IsExtensionInstalled(ideType IDEType, extensionID string) (bool, error) {
	ide, exists := d.detectedIDEs[ideType]
	if !exists || !ide.IsInstalled {
		return false, fmt.Errorf(i18n.T("error.ide.notDetected", map[string]interface{}{"IDE": string(ideType)}))
	}

	if ide.ExtensionsPath == "" {
		return false, fmt.Errorf(i18n.T("error.ide.extensionsPathNotFound", map[string]interface{}{"IDE": ide.Name}))
	}

	// Check if extension directory exists
	// Extension'lar genellikle şu formatta saklanır: publisherName.extensionName-version
	// Bizim durumumuzda: mehmetsenol.gorev-vscode-x.y.z
	extensionPattern := strings.Split(extensionID, ".")[1] // gorev-vscode kısmını al

	entries, err := os.ReadDir(ide.ExtensionsPath)
	if err != nil {
		return false, fmt.Errorf(i18n.T("error.ide.cannotReadExtensions", map[string]interface{}{
			"Path":  ide.ExtensionsPath,
			"Error": err,
		}))
	}

	for _, entry := range entries {
		if entry.IsDir() && strings.Contains(entry.Name(), extensionPattern) {
			return true, nil
		}
	}

	return false, nil
}

// GetExtensionVersion returns the installed version of the Gorev extension
func (d *IDEDetector) GetExtensionVersion(ideType IDEType, extensionID string) (string, error) {
	ide, exists := d.detectedIDEs[ideType]
	if !exists || !ide.IsInstalled {
		return "", fmt.Errorf(i18n.T("error.ide.notDetected", map[string]interface{}{"IDE": string(ideType)}))
	}

	if ide.ExtensionsPath == "" {
		return "", fmt.Errorf(i18n.T("error.ide.extensionsPathNotFound", map[string]interface{}{"IDE": ide.Name}))
	}

	// Extension directory'sini bul
	extensionPattern := strings.Split(extensionID, ".")[1]

	entries, err := os.ReadDir(ide.ExtensionsPath)
	if err != nil {
		return "", err
	}

	for _, entry := range entries {
		if entry.IsDir() && strings.Contains(entry.Name(), extensionPattern) {
			// Extension directory name'inden version'ı çıkar
			// Format: mehmetsenol.gorev-vscode-x.y.z
			parts := strings.Split(entry.Name(), "-")
			if len(parts) >= 2 {
				return parts[len(parts)-1], nil
			}
		}
	}

	return "", fmt.Errorf(i18n.T("error.ide.extensionNotInstalled", map[string]interface{}{
		"Extension": extensionID,
		"IDE":       ide.Name,
	}))
}
