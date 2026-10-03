package devharness

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// GetVSCodeConfigPath returns the platform-specific VS Code User directory path.
func GetVSCodeConfigPath(homeDir string) (string, error) {
	switch runtime.GOOS {
	case "linux":
		return filepath.Join(homeDir, ".config", "Code", "User"), nil
	case "darwin":
		return filepath.Join(homeDir, "Library", "Application Support", "Code", "User"), nil
	case "windows":
		appData := os.Getenv("APPDATA")
		if appData != "" {
			return filepath.Join(appData, "Code", "User"), nil
		}
		return filepath.Join(homeDir, "AppData", "Roaming", "Code", "User"), nil
	default:
		return "", errors.New("unsupported platform: " + runtime.GOOS)
	}
}

// GetCursorConfigPath returns the platform-specific Cursor configuration directory path.
func GetCursorConfigPath(homeDir string) (string, error) {
	dotCursor := filepath.Join(homeDir, ".cursor")
	if dirExists(dotCursor) {
		return dotCursor, nil
	}

	switch runtime.GOOS {
	case "linux":
		return filepath.Join(homeDir, ".config", "Cursor", "User"), nil
	case "darwin":
		return filepath.Join(homeDir, "Library", "Application Support", "Cursor", "User"), nil
	case "windows":
		appData := os.Getenv("APPDATA")
		if appData != "" {
			return filepath.Join(appData, "Cursor", "User"), nil
		}
		return filepath.Join(homeDir, "AppData", "Roaming", "Cursor", "User"), nil
	default:
		return dotCursor, nil
	}
}

// GetAntigravityVSCodeConfigPath returns the Antigravity config directory path (~/.gemini/config)
// for the Antigravity VS Code extension when it is installed.
func GetAntigravityVSCodeConfigPath(homeDir string) (string, error) {
	if !IsAntigravityVSCodeInstalled(homeDir) {
		return "", errors.New("antigravity VS Code extension is not installed")
	}
	configDir := filepath.Join(homeDir, ".gemini", "config")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "", err
	}
	return configDir, nil
}

// IsAntigravityVSCodeInstalled checks if the Google Antigravity extension is installed in VS Code.
func IsAntigravityVSCodeInstalled(homeDir string) bool {
	// 1. Check standard ~/.vscode/extensions directory
	extDir := filepath.Join(homeDir, ".vscode", "extensions")
	if entries, err := os.ReadDir(extDir); err == nil {
		for _, entry := range entries {
			if strings.HasPrefix(entry.Name(), "google.google-antigravity") {
				return true
			}
		}
	}

	// 2. Check VS Code User directory profiles and sync data
	vscUserDir, err := GetVSCodeConfigPath(homeDir)
	if err == nil {
		// Check CachedExtensionVSIXs (in parent of User dir, e.g. ~/.config/Code/CachedExtensionVSIXs)
		codeDir := filepath.Dir(vscUserDir)
		cachedVsixDir := filepath.Join(codeDir, "CachedExtensionVSIXs")
		if entries, err := os.ReadDir(cachedVsixDir); err == nil {
			for _, entry := range entries {
				if strings.HasPrefix(entry.Name(), "google.google-antigravity") {
					return true
				}
			}
		}

		// Check profiles/*/extensions.json
		profilesDir := filepath.Join(vscUserDir, "profiles")
		if entries, err := os.ReadDir(profilesDir); err == nil {
			for _, entry := range entries {
				if entry.IsDir() {
					extFile := filepath.Join(profilesDir, entry.Name(), "extensions.json")
					if containsExtensionID(extFile, "google.google-antigravity") {
						return true
					}
				}
			}
		}

		// Check User/extensions.json and sync/extensions/lastSyncextensions.json
		for _, rel := range []string{
			"extensions.json",
			filepath.Join("sync", "extensions", "lastSyncextensions.json"),
		} {
			if containsExtensionID(filepath.Join(vscUserDir, rel), "google.google-antigravity") {
				return true
			}
		}
	}

	return false
}

// FindMCPConfigPaths resolves all config file paths based on IDE profile structure.
func FindMCPConfigPaths(basePath string, configFileName string) ([]string, error) {
	if !dirExists(basePath) {
		return nil, errors.New("directory not found")
	}

	profilesPath := filepath.Join(basePath, "profiles")
	if !dirExists(profilesPath) {
		return []string{filepath.Join(basePath, configFileName)}, nil
	}

	entries, err := os.ReadDir(profilesPath)
	if err != nil {
		return nil, err
	}

	var configPaths []string
	for _, entry := range entries {
		if entry.IsDir() {
			configPaths = append(configPaths, filepath.Join(profilesPath, entry.Name(), configFileName))
		}
	}

	if len(configPaths) == 0 {
		return []string{filepath.Join(basePath, configFileName)}, nil
	}

	return configPaths, nil
}

func containsExtensionID(filePath string, extID string) bool {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return false
	}
	return strings.Contains(string(data), extID)
}

func dirExists(path string) bool {
	fi, err := os.Stat(path)
	return err == nil && fi.IsDir()
}

func fileExists(path string) bool {
	fi, err := os.Stat(path)
	return err == nil && !fi.IsDir()
}
