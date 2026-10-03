package devharness

import (
	"path/filepath"
)

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
