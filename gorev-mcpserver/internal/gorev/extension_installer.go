package gorev

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/msenol/gorev/internal/i18n"
)

// ExtensionInfo holds information about an extension
type ExtensionInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Version     string `json:"version"`
	Publisher   string `json:"publisher"`
	DownloadURL string `json:"download_url"`
	Checksum    string `json:"checksum,omitempty"`
}

// InstallResult represents the result of an extension installation
type InstallResult struct {
	Success   bool   `json:"success"`
	Message   string `json:"message"`
	Extension string `json:"extension"`
	IDE       string `json:"ide"`
	Version   string `json:"version,omitempty"`
}

// ExtensionInstaller handles extension installation and management
type ExtensionInstaller struct {
	detector     *IDEDetector
	downloadPath string
	client       *http.Client
}

// NewExtensionInstaller creates a new extension installer
func NewExtensionInstaller(detector *IDEDetector) *ExtensionInstaller {
	// Create temp directory for downloads
	tempDir := filepath.Join(os.TempDir(), "gorev-extensions")
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		fmt.Printf("Warning: failed to create temp dir %s: %v\n", tempDir, err)
	}

	return &ExtensionInstaller{
		detector:     detector,
		downloadPath: tempDir,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// InstallExtension installs the Gorev extension to the specified IDE
func (ei *ExtensionInstaller) InstallExtension(ctx context.Context, ideType IDEType, extensionInfo *ExtensionInfo) (*InstallResult, error) {
	// Handle nil extension info
	if extensionInfo == nil {
		return nil, fmt.Errorf("extension info cannot be nil")
	}

	result := &InstallResult{
		Extension: extensionInfo.ID,
		IDE:       string(ideType),
	}

	// Check if IDE is detected
	ide, exists := ei.detector.GetDetectedIDE(ideType)
	if !exists {
		result.Message = i18n.T("error.ide.notDetected", map[string]interface{}{"IDE": string(ideType)})
		return result, fmt.Errorf(result.Message)
	}

	// Check if extension is already installed
	isInstalled, err := ei.detector.IsExtensionInstalled(ideType, extensionInfo.ID)
	if err == nil && isInstalled {
		// Check version
		installedVersion, err := ei.detector.GetExtensionVersion(ideType, extensionInfo.ID)
		if err == nil && installedVersion == extensionInfo.Version {
			result.Success = true
			result.Version = installedVersion
			result.Message = i18n.T("success.ide.extensionAlreadyInstalled", map[string]interface{}{
				"Extension": extensionInfo.Name,
				"IDE":       ide.Name,
				"Version":   installedVersion,
			})
			return result, nil
		}
	}

	// Download VSIX file
	vsixPath, err := ei.downloadVSIX(ctx, extensionInfo)
	if err != nil {
		result.Message = i18n.T("error.ide.downloadFailed", map[string]interface{}{
			"Extension": extensionInfo.Name,
			"Error":     err,
		})
		return result, err
	}
	defer func() { _ = os.Remove(vsixPath) }() // Cleanup

	// Install using IDE CLI
	err = ei.installVSIX(ide, vsixPath)
	if err != nil {
		result.Message = i18n.T("error.ide.installFailed", map[string]interface{}{
			"Extension": extensionInfo.Name,
			"IDE":       ide.Name,
			"Error":     err,
		})
		return result, err
	}

	result.Success = true
	result.Version = extensionInfo.Version
	result.Message = i18n.T("success.ide.extensionInstalled", map[string]interface{}{
		"Extension": extensionInfo.Name,
		"IDE":       ide.Name,
		"Version":   extensionInfo.Version,
	})

	return result, nil
}

// InstallToAllIDEs installs the extension to all detected IDEs
func (ei *ExtensionInstaller) InstallToAllIDEs(ctx context.Context, extensionInfo *ExtensionInfo) ([]InstallResult, error) {
	var results []InstallResult

	allIDEs := ei.detector.GetAllDetectedIDEs()
	if len(allIDEs) == 0 {
		return results, fmt.Errorf(i18n.T("error.ide.noIDEsDetected"))
	}

	for ideType := range allIDEs {
		result, err := ei.InstallExtension(ctx, ideType, extensionInfo)
		if result != nil {
			results = append(results, *result)
		}
		if err != nil {
			// Continue with other IDEs even if one fails
			continue
		}
	}

	return results, nil
}

// UninstallExtension removes the extension from the specified IDE
func (ei *ExtensionInstaller) UninstallExtension(ideType IDEType, extensionID string) (*InstallResult, error) {
	result := &InstallResult{
		Extension: extensionID,
		IDE:       string(ideType),
	}

	ide, exists := ei.detector.GetDetectedIDE(ideType)
	if !exists {
		result.Message = i18n.T("error.ide.notDetected", map[string]interface{}{"IDE": string(ideType)})
		return result, fmt.Errorf(result.Message)
	}

	// Check if extension is installed
	isInstalled, err := ei.detector.IsExtensionInstalled(ideType, extensionID)
	if err != nil || !isInstalled {
		result.Message = i18n.T("error.ide.extensionNotInstalled", map[string]interface{}{
			"Extension": extensionID,
			"IDE":       ide.Name,
		})
		return result, fmt.Errorf(result.Message)
	}

	// Uninstall using IDE CLI
	err = ei.uninstallExtension(ide, extensionID)
	if err != nil {
		result.Message = i18n.T("error.ide.uninstallFailed", map[string]interface{}{
			"Extension": extensionID,
			"IDE":       ide.Name,
			"Error":     err,
		})
		return result, err
	}

	result.Success = true
	result.Message = i18n.T("success.ide.extensionUninstalled", map[string]interface{}{
		"Extension": extensionID,
		"IDE":       ide.Name,
	})

	return result, nil
}

// downloadVSIX downloads the VSIX file for the extension
func (ei *ExtensionInstaller) downloadVSIX(ctx context.Context, extensionInfo *ExtensionInfo) (string, error) {
	// Create filename
	filename := fmt.Sprintf("%s-%s.vsix", extensionInfo.ID, extensionInfo.Version)
	filePath := filepath.Join(ei.downloadPath, filename)

	// Check if file already exists and is valid
	if fileExists(filePath) {
		if extensionInfo.Checksum != "" {
			valid, err := ei.verifyChecksum(filePath, extensionInfo.Checksum)
			if err == nil && valid {
				return filePath, nil
			}
		}
		// Remove invalid file
		_ = os.Remove(filePath)
	}

	// Create request with context
	req, err := http.NewRequestWithContext(ctx, "GET", extensionInfo.DownloadURL, nil)
	if err != nil {
		return "", err
	}

	// Set user agent
	req.Header.Set("User-Agent", "Gorev-Extension-Installer/1.0")

	// Download file
	resp, err := ei.client.Do(req)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download failed with status: %d", resp.StatusCode)
	}

	// Create file
	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Printf("Warning: failed to close file %s: %v\n", filePath, err)
		}
	}()

	// Copy data
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		_ = os.Remove(filePath)
		return "", err
	}

	// Verify checksum if provided
	if extensionInfo.Checksum != "" {
		valid, err := ei.verifyChecksum(filePath, extensionInfo.Checksum)
		if err != nil || !valid {
			_ = os.Remove(filePath)
			return "", fmt.Errorf("checksum verification failed")
		}
	}

	return filePath, nil
}

// verifyChecksum verifies the SHA256 checksum of a file
func (ei *ExtensionInstaller) verifyChecksum(filePath, expectedChecksum string) (bool, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return false, err
	}
	defer func() { _ = file.Close() }()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return false, err
	}

	actualChecksum := hex.EncodeToString(hash.Sum(nil))
	return strings.EqualFold(actualChecksum, expectedChecksum), nil
}

// installVSIX installs a VSIX file using the IDE's CLI
func (ei *ExtensionInstaller) installVSIX(ide *IDEInfo, vsixPath string) error {
	if ide.ExecutablePath == "" {
		return fmt.Errorf("IDE executable path not found")
	}

	var args []string
	switch ide.Type {
	case IDETypeVSCode, IDETypeCursor, IDETypeWindsurf:
		args = []string{"--install-extension", vsixPath}
	default:
		return fmt.Errorf("unsupported IDE type: %s", ide.Type)
	}

	cmd := exec.Command(ide.ExecutablePath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// uninstallExtension uninstalls an extension using the IDE's CLI
func (ei *ExtensionInstaller) uninstallExtension(ide *IDEInfo, extensionID string) error {
	if ide.ExecutablePath == "" {
		return fmt.Errorf("IDE executable path not found")
	}

	var args []string
	switch ide.Type {
	case IDETypeVSCode, IDETypeCursor, IDETypeWindsurf:
		args = []string{"--uninstall-extension", extensionID}
	default:
		return fmt.Errorf("unsupported IDE type: %s", ide.Type)
	}

	cmd := exec.Command(ide.ExecutablePath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// ListInstalledExtensions lists all installed extensions in the IDE
func (ei *ExtensionInstaller) ListInstalledExtensions(ideType IDEType) ([]string, error) {
	ide, exists := ei.detector.GetDetectedIDE(ideType)
	if !exists {
		return nil, fmt.Errorf("IDE not detected: %s", ideType)
	}

	if ide.ExecutablePath == "" {
		return nil, fmt.Errorf("IDE executable path not found")
	}

	var args []string
	switch ide.Type {
	case IDETypeVSCode, IDETypeCursor, IDETypeWindsurf:
		args = []string{"--list-extensions"}
	default:
		return nil, fmt.Errorf("unsupported IDE type: %s", ide.Type)
	}

	cmd := exec.Command(ide.ExecutablePath, args...)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	extensions := strings.Split(strings.TrimSpace(string(output)), "\n")
	var result []string
	for _, ext := range extensions {
		if strings.TrimSpace(ext) != "" {
			result = append(result, strings.TrimSpace(ext))
		}
	}

	return result, nil
}

// GetLatestExtensionInfo fetches the latest extension information from GitHub releases
func (ei *ExtensionInstaller) GetLatestExtensionInfo(ctx context.Context, repoOwner, repoName string) (*ExtensionInfo, error) {
	// GitHub releases API URL
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", repoOwner, repoName)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("User-Agent", "Gorev-Extension-Installer/1.0")

	resp, err := ei.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API request failed with status: %d", resp.StatusCode)
	}

	var release struct {
		TagName string `json:"tag_name"`
		Assets  []struct {
			Name               string `json:"name"`
			BrowserDownloadURL string `json:"browser_download_url"`
		} `json:"assets"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, err
	}

	// Find VSIX asset
	var vsixAsset *struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	}

	for _, asset := range release.Assets {
		if strings.HasSuffix(asset.Name, ".vsix") {
			vsixAsset = &asset
			break
		}
	}

	if vsixAsset == nil {
		return nil, fmt.Errorf("no VSIX asset found in release")
	}

	// Parse extension info from asset name
	// Expected format: gorev-vscode-x.y.z.vsix
	baseName := strings.TrimSuffix(vsixAsset.Name, ".vsix")
	parts := strings.Split(baseName, "-")
	if len(parts) < 3 {
		return nil, fmt.Errorf("invalid asset name format: %s", vsixAsset.Name)
	}

	version := strings.TrimPrefix(release.TagName, "v")

	return &ExtensionInfo{
		ID:          "mehmetsenol.gorev-vscode", // This should match package.json
		Name:        "Gorev VS Code Extension",
		Version:     version,
		Publisher:   "mehmetsenol",
		DownloadURL: vsixAsset.BrowserDownloadURL,
	}, nil
}

// GetDownloadPath returns the download directory path
func (ei *ExtensionInstaller) GetDownloadPath() string {
	return ei.downloadPath
}

// Cleanup removes downloaded files
func (ei *ExtensionInstaller) Cleanup() error {
	return os.RemoveAll(ei.downloadPath)
}
