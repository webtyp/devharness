package devharness_test

import (
	"os"
	"path/filepath"
	"testing"

	"webtyp.com/devharness"
)

func TestCatalog_All(t *testing.T) {
	all := devharness.All()
	if len(all) < 10 {
		t.Fatalf("expected at least 10 harnesses, got %d", len(all))
	}

	ids := map[string]bool{}
	for _, h := range all {
		ids[h.ID] = true
	}

	expected := []string{
		"antigravity-vsc",
		"antigravity",
		"claude",
		"gemini",
		"codex",
		"qwen",
		"opencode",
		"agents",
		"vscode",
		"cursor",
	}

	for _, exp := range expected {
		if !ids[exp] {
			t.Errorf("expected harness ID %q in All()", exp)
		}
	}
}

func TestCatalog_Skills(t *testing.T) {
	skills := devharness.Skills()
	for _, h := range skills {
		if !h.SupportsSkills {
			t.Errorf("harness %s returned in Skills() but SupportsSkills is false", h.ID)
		}
	}

	ids := map[string]bool{}
	for _, h := range skills {
		ids[h.ID] = true
	}

	if !ids["antigravity-vsc"] {
		t.Error("expected antigravity-vsc in Skills()")
	}
	if ids["vscode"] {
		t.Error("vscode should not be in Skills()")
	}
	if ids["cursor"] {
		t.Error("cursor should not be in Skills()")
	}
}

func TestCatalog_MCP(t *testing.T) {
	mcp := devharness.MCP()
	for _, h := range mcp {
		if !h.SupportsMCP {
			t.Errorf("harness %s returned in MCP() but SupportsMCP is false", h.ID)
		}
	}

	ids := map[string]bool{}
	for _, h := range mcp {
		ids[h.ID] = true
	}

	if !ids["antigravity-vsc"] {
		t.Error("expected antigravity-vsc in MCP()")
	}
	if !ids["vscode"] {
		t.Error("expected vscode in MCP()")
	}
	if !ids["cursor"] {
		t.Error("expected cursor in MCP()")
	}
	if ids["gemini"] {
		t.Error("gemini should not be in MCP()")
	}
}

func TestFind(t *testing.T) {
	tests := []struct {
		query   string
		wantID  string
		wantErr bool
	}{
		{"antigravity-vsc", "antigravity-vsc", false},
		{"antigravity", "antigravity-vsc", false}, // Alias
		{"Antigravity", "antigravity-vsc", false}, // Case-insensitive
		{"claude", "claude", false},
		{"claude-code", "claude", false}, // Alias
		{"vsc", "vscode", false},         // Alias
		{"vscode", "vscode", false},
		{"cursor", "cursor", false},
		{"nonexistent", "", true},
	}

	for _, tt := range tests {
		h, found := devharness.Find(tt.query)
		if tt.wantErr {
			if found {
				t.Errorf("Find(%q) expected not found, got %s", tt.query, h.ID)
			}
		} else {
			if !found {
				t.Fatalf("Find(%q) expected found, got false", tt.query)
			}
			if h.ID != tt.wantID {
				t.Errorf("Find(%q) = %s, want %s", tt.query, h.ID, tt.wantID)
			}
		}
	}
}

func TestIsAntigravityVSCodeInstalled_ExtensionDir(t *testing.T) {
	tempHome := t.TempDir()

	if devharness.IsAntigravityVSCodeInstalled(tempHome) {
		t.Error("expected false for empty directory")
	}

	extDir := filepath.Join(tempHome, ".vscode", "extensions", "google.google-antigravity-1.6.0")
	if err := os.MkdirAll(extDir, 0755); err != nil {
		t.Fatal(err)
	}

	if !devharness.IsAntigravityVSCodeInstalled(tempHome) {
		t.Error("expected true when extension dir exists")
	}
}

func TestIsAntigravityVSCodeInstalled_Profiles(t *testing.T) {
	tempHome := t.TempDir()

	vscUserDir, err := devharness.GetVSCodeConfigPath(tempHome)
	if err != nil {
		t.Fatalf("failed to get vscode config path: %v", err)
	}

	profileDir := filepath.Join(vscUserDir, "profiles", "custom-profile")
	if err := os.MkdirAll(profileDir, 0755); err != nil {
		t.Fatal(err)
	}
	extJSON := `[{"identifier":{"id":"google.google-antigravity"}}]`
	if err := os.WriteFile(filepath.Join(profileDir, "extensions.json"), []byte(extJSON), 0644); err != nil {
		t.Fatal(err)
	}

	if !devharness.IsAntigravityVSCodeInstalled(tempHome) {
		t.Error("expected true when extension is listed in profile extensions.json")
	}
}

func TestGetAntigravityVSCodeConfigPath(t *testing.T) {
	tempHome := t.TempDir()

	if _, err := devharness.GetAntigravityVSCodeConfigPath(tempHome); err == nil {
		t.Error("expected error when extension is not installed")
	}

	extDir := filepath.Join(tempHome, ".vscode", "extensions", "google.google-antigravity-1.6.0")
	if err := os.MkdirAll(extDir, 0755); err != nil {
		t.Fatal(err)
	}

	cfgDir, err := devharness.GetAntigravityVSCodeConfigPath(tempHome)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := filepath.Join(tempHome, ".gemini", "config")
	if cfgDir != expected {
		t.Errorf("got %s, want %s", cfgDir, expected)
	}

	if _, err := os.Stat(cfgDir); err != nil {
		t.Errorf("expected directory %s to be created", cfgDir)
	}
}

func TestGetCursorConfigPath(t *testing.T) {
	tempHome := t.TempDir()

	// Case 1: ~/.cursor directory exists
	dotCursor := filepath.Join(tempHome, ".cursor")
	if err := os.MkdirAll(dotCursor, 0755); err != nil {
		t.Fatal(err)
	}

	cfgDir, err := devharness.GetCursorConfigPath(tempHome)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfgDir != dotCursor {
		t.Errorf("got %s, want %s", cfgDir, dotCursor)
	}
}

func TestFindMCPConfigPaths(t *testing.T) {
	tempHome := t.TempDir()

	// Non-existent directory
	if _, err := devharness.FindMCPConfigPaths(filepath.Join(tempHome, "missing"), "mcp.json"); err == nil {
		t.Error("expected error for non-existent directory")
	}

	// Single config file (no profiles)
	baseDir := filepath.Join(tempHome, "ide")
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		t.Fatal(err)
	}
	paths, err := devharness.FindMCPConfigPaths(baseDir, "mcp.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(paths) != 1 || paths[0] != filepath.Join(baseDir, "mcp.json") {
		t.Errorf("expected [%s], got %v", filepath.Join(baseDir, "mcp.json"), paths)
	}

	// Profiles directory
	profile1 := filepath.Join(baseDir, "profiles", "p1")
	profile2 := filepath.Join(baseDir, "profiles", "p2")
	if err := os.MkdirAll(profile1, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(profile2, 0755); err != nil {
		t.Fatal(err)
	}

	paths, err = devharness.FindMCPConfigPaths(baseDir, "mcp.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(paths) != 2 {
		t.Fatalf("expected 2 profile paths, got %d", len(paths))
	}
}

func TestDetectSkills_And_DetectMCP(t *testing.T) {
	tempHome := t.TempDir()

	// Empty home -> 0 detected
	if len(devharness.DetectSkills(tempHome)) != 0 {
		t.Errorf("expected 0 skills harnesses in empty home")
	}
	if len(devharness.DetectMCP(tempHome)) != 0 {
		t.Errorf("expected 0 MCP harnesses in empty home")
	}

	// Create ~/.claude
	if err := os.MkdirAll(filepath.Join(tempHome, ".claude"), 0755); err != nil {
		t.Fatal(err)
	}

	skillsH := devharness.DetectSkills(tempHome)
	if len(skillsH) != 1 || skillsH[0].ID != "claude" {
		t.Errorf("expected [claude], got %v", skillsH)
	}

	mcpH := devharness.DetectMCP(tempHome)
	if len(mcpH) != 1 || mcpH[0].ID != "claude" {
		t.Errorf("expected [claude] in MCP, got %v", mcpH)
	}
}
