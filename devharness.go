package devharness

import (
	"errors"
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
