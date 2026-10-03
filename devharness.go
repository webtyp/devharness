package devharness

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// MCPConfig holds MCP-specific configuration parameters for a harness.
type MCPConfig struct {
	FileName     string
	ServersKey   string
	URLKey       string
	ExtraFields  map[string]any
	HasInputs    bool
	SkipProfiles bool
}

// Harness represents an LLM, agent runtime, or IDE execution environment.
type Harness struct {
	ID      string
	Name    string
	Aliases []string

	// Skills capabilities
	SupportsSkills bool
	GetSkillsDir   func(homeDir string) (string, error)

	// MCP capabilities
	SupportsMCP  bool
	GetConfigDir func(homeDir string) (string, error)
	MCP          MCPConfig

	// Detection
	IsInstalled func(homeDir string) bool
}

// SkillsDir returns the directory where Agent Skills should be installed or linked.
func (h Harness) SkillsDir(homeDir string) (string, error) {
	if !h.SupportsSkills || h.GetSkillsDir == nil {
		return "", errors.New("harness " + h.ID + " does not support skills")
	}
	return h.GetSkillsDir(homeDir)
}

// ConfigDir returns the configuration directory for this harness.
func (h Harness) ConfigDir(homeDir string) (string, error) {
	if h.GetConfigDir == nil {
		return "", errors.New("harness " + h.ID + " has no config directory")
	}
	return h.GetConfigDir(homeDir)
}

// Installed returns whether this harness is currently detected on the system.
func (h Harness) Installed(homeDir string) bool {
	if h.IsInstalled == nil {
		return false
	}
	return h.IsInstalled(homeDir)
}

var registry = []Harness{
	{
		ID:             "antigravity-vsc",
		Name:           "Antigravity (VS Code)",
		Aliases:        []string{"antigravity"},
		SupportsSkills: true,
		GetSkillsDir: func(homeDir string) (string, error) {
			return filepath.Join(homeDir, ".gemini", "config", "skills"), nil
		},
		SupportsMCP:  true,
		GetConfigDir: GetAntigravityVSCodeConfigPath,
		MCP: MCPConfig{
			FileName:     "mcp_config.json",
			ServersKey:   "mcpServers",
			URLKey:       "serverUrl",
			SkipProfiles: true,
		},
		IsInstalled: IsAntigravityVSCodeInstalled,
	},
	{
		ID:             "antigravity",
		Name:           "Antigravity",
		SupportsSkills: true,
		GetSkillsDir: func(homeDir string) (string, error) {
			return filepath.Join(homeDir, ".gemini", "antigravity", "skills"), nil
		},
		SupportsMCP: true,
		GetConfigDir: func(homeDir string) (string, error) {
			return filepath.Join(homeDir, ".gemini", "antigravity"), nil
		},
		MCP: MCPConfig{
			FileName:   "mcp_config.json",
			ServersKey: "mcpServers",
			URLKey:     "serverUrl",
		},
		IsInstalled: func(homeDir string) bool {
			return dirExists(filepath.Join(homeDir, ".gemini", "antigravity"))
		},
	},
	{
		ID:             "claude",
		Name:           "Claude Code",
		Aliases:        []string{"claude-code"},
		SupportsSkills: true,
		GetSkillsDir: func(homeDir string) (string, error) {
			return filepath.Join(homeDir, ".claude", "skills"), nil
		},
		SupportsMCP: true,
		GetConfigDir: func(homeDir string) (string, error) {
			return homeDir, nil
		},
		MCP: MCPConfig{
			FileName:     ".claude.json",
			ServersKey:   "mcpServers",
			URLKey:       "url",
			ExtraFields:  map[string]any{"type": "http"},
			SkipProfiles: true,
		},
		IsInstalled: func(homeDir string) bool {
			return dirExists(filepath.Join(homeDir, ".claude")) || fileExists(filepath.Join(homeDir, ".claude.json"))
		},
	},
	{
		ID:             "gemini",
		Name:           "Gemini CLI",
		SupportsSkills: true,
		GetSkillsDir: func(homeDir string) (string, error) {
			return filepath.Join(homeDir, ".gemini", "skills"), nil
		},
		IsInstalled: func(homeDir string) bool {
			return dirExists(filepath.Join(homeDir, ".gemini"))
		},
	},
	{
		ID:             "codex",
		Name:           "Codex",
		SupportsSkills: true,
		GetSkillsDir: func(homeDir string) (string, error) {
			return filepath.Join(homeDir, ".codex", "skills"), nil
		},
		IsInstalled: func(homeDir string) bool {
			return dirExists(filepath.Join(homeDir, ".codex"))
		},
	},
	{
		ID:             "qwen",
		Name:           "Qwen",
		SupportsSkills: true,
		GetSkillsDir: func(homeDir string) (string, error) {
			return filepath.Join(homeDir, ".qwen", "skills"), nil
		},
		IsInstalled: func(homeDir string) bool {
			return dirExists(filepath.Join(homeDir, ".qwen"))
		},
	},
	{
		ID:             "opencode",
		Name:           "OpenCode",
		SupportsSkills: true,
		GetSkillsDir: func(homeDir string) (string, error) {
			return filepath.Join(homeDir, ".config", "opencode", "skills"), nil
		},
		SupportsMCP: true,
		GetConfigDir: func(homeDir string) (string, error) {
			return filepath.Join(homeDir, ".config", "opencode"), nil
		},
		MCP: MCPConfig{
			FileName:     "opencode.json",
			ServersKey:   "mcp",
			URLKey:       "url",
			ExtraFields:  map[string]any{"type": "remote", "enabled": true},
			SkipProfiles: true,
		},
		IsInstalled: func(homeDir string) bool {
			return dirExists(filepath.Join(homeDir, ".config", "opencode"))
		},
	},
	{
		// agents is the vendor-neutral convention (~/.agents/skills) used by
		// open-source agents (such as OpenCode) and multi-agent tools to auto-load
		// skills without being tied to a single vendor-specific directory.
		ID:             "agents",
		Name:           "Agents (Vendor-neutral)",
		SupportsSkills: true,
		GetSkillsDir: func(homeDir string) (string, error) {
			return filepath.Join(homeDir, ".agents", "skills"), nil
		},
		IsInstalled: func(homeDir string) bool {
			return dirExists(filepath.Join(homeDir, ".agents"))
		},
	},
	{
		ID:          "vscode",
		Name:        "Visual Studio Code",
		Aliases:     []string{"vsc"},
		SupportsMCP: true,
		GetConfigDir: func(homeDir string) (string, error) {
			return GetVSCodeConfigPath(homeDir)
		},
		MCP: MCPConfig{
			FileName:     "mcp.json",
			ServersKey:   "servers",
			URLKey:       "url",
			ExtraFields:  map[string]any{"type": "http", "autoStart": true},
			HasInputs:    true,
			SkipProfiles: false,
		},
		IsInstalled: func(homeDir string) bool {
			dir, err := GetVSCodeConfigPath(homeDir)
			return err == nil && dirExists(dir)
		},
	},
	{
		ID:          "cursor",
		Name:        "Cursor",
		SupportsMCP: true,
		GetConfigDir: func(homeDir string) (string, error) {
			return GetCursorConfigPath(homeDir)
		},
		MCP: MCPConfig{
			FileName:     "mcp.json",
			ServersKey:   "mcpServers",
			URLKey:       "url",
			SkipProfiles: true,
		},
		IsInstalled: func(homeDir string) bool {
			dir, err := GetCursorConfigPath(homeDir)
			return (err == nil && dirExists(dir)) || dirExists(filepath.Join(homeDir, ".cursor"))
		},
	},
	{
		ID:             "pi",
		Name:           "Pi Agent",
		SupportsSkills: true,
		GetSkillsDir: func(homeDir string) (string, error) {
			return filepath.Join(homeDir, ".pi", "agent", "skills"), nil
		},
		SupportsMCP: true,
		GetConfigDir: func(homeDir string) (string, error) {
			return filepath.Join(homeDir, ".pi", "agent"), nil
		},
		MCP: MCPConfig{
			FileName:     "mcp.json",
			ServersKey:   "mcpServers",
			URLKey:       "url",
			SkipProfiles: true,
		},
		IsInstalled: func(homeDir string) bool {
			return dirExists(filepath.Join(homeDir, ".pi")) || dirExists(filepath.Join(homeDir, ".pi", "agent"))
		},
	},
	{
		ID:          "windsurf",
		Name:        "Windsurf",
		SupportsMCP: true,
		GetConfigDir: func(homeDir string) (string, error) {
			return filepath.Join(homeDir, ".codeium", "windsurf"), nil
		},
		MCP: MCPConfig{
			FileName:     "mcp_config.json",
			ServersKey:   "mcpServers",
			URLKey:       "serverUrl",
			SkipProfiles: true,
		},
		IsInstalled: func(homeDir string) bool {
			return dirExists(filepath.Join(homeDir, ".codeium", "windsurf"))
		},
	},
	{
		ID:          "zed",
		Name:        "Zed",
		SupportsMCP: true,
		GetConfigDir: func(homeDir string) (string, error) {
			return filepath.Join(homeDir, ".config", "zed"), nil
		},
		MCP: MCPConfig{
			FileName:     "settings.json",
			ServersKey:   "context_servers",
			URLKey:       "url",
			SkipProfiles: true,
		},
		IsInstalled: func(homeDir string) bool {
			return dirExists(filepath.Join(homeDir, ".config", "zed"))
		},
	},
	{
		ID:             "hermes",
		Name:           "Hermes Agent",
		Aliases:        []string{"hermes-agent"},
		SupportsSkills: true,
		GetSkillsDir: func(homeDir string) (string, error) {
			return filepath.Join(homeDir, ".hermes", "skills"), nil
		},
		SupportsMCP: true,
		GetConfigDir: func(homeDir string) (string, error) {
			return filepath.Join(homeDir, ".hermes"), nil
		},
		MCP: MCPConfig{
			FileName:     "mcp.json",
			ServersKey:   "mcpServers",
			URLKey:       "url",
			SkipProfiles: true,
		},
		IsInstalled: func(homeDir string) bool {
			return dirExists(filepath.Join(homeDir, ".hermes"))
		},
	},
	{
		ID:             "deepseek",
		Name:           "DeepSeek",
		Aliases:        []string{"deepseek-tui"},
		SupportsSkills: true,
		GetSkillsDir: func(homeDir string) (string, error) {
			return filepath.Join(homeDir, ".deepseek", "skills"), nil
		},
		SupportsMCP: true,
		GetConfigDir: func(homeDir string) (string, error) {
			return filepath.Join(homeDir, ".deepseek"), nil
		},
		MCP: MCPConfig{
			FileName:     "mcp.json",
			ServersKey:   "mcpServers",
			URLKey:       "url",
			SkipProfiles: true,
		},
		IsInstalled: func(homeDir string) bool {
			return dirExists(filepath.Join(homeDir, ".deepseek"))
		},
	},
}

// All returns all registered harnesses.
func All() []Harness {
	out := make([]Harness, len(registry))
	copy(out, registry)
	return out
}

// Skills returns all harnesses that support Agent Skills.
func Skills() []Harness {
	var out []Harness
	for _, h := range registry {
		if h.SupportsSkills {
			out = append(out, h)
		}
	}
	return out
}

// MCP returns all harnesses that support MCP configuration.
func MCP() []Harness {
	var out []Harness
	for _, h := range registry {
		if h.SupportsMCP {
			out = append(out, h)
		}
	}
	return out
}

// DetectInstalled returns all harnesses detected on the host system.
func DetectInstalled(homeDir string) []Harness {
	var out []Harness
	for _, h := range registry {
		if h.Installed(homeDir) {
			out = append(out, h)
		}
	}
	return out
}

// DetectSkills returns all installed harnesses that support Agent Skills.
func DetectSkills(homeDir string) []Harness {
	var out []Harness
	for _, h := range Skills() {
		if h.Installed(homeDir) {
			out = append(out, h)
		}
	}
	return out
}

// DetectMCP returns all installed harnesses that support MCP configuration.
func DetectMCP(homeDir string) []Harness {
	var out []Harness
	for _, h := range MCP() {
		if h.Installed(homeDir) {
			out = append(out, h)
		}
	}
	return out
}

// Find searches for a harness by canonical ID or alias.
func Find(name string) (Harness, bool) {
	for _, h := range registry {
		if strings.EqualFold(h.ID, name) {
			return h, true
		}
		for _, alias := range h.Aliases {
			if strings.EqualFold(alias, name) {
				return h, true
			}
		}
	}
	return Harness{}, false
}

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
