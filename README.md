# devharness

Harness and IDE discovery for LLMs, autonomous agents, and editors in the WebTyp ecosystem.

`devharness` is a zero-dependency Go library that provides a single, unified source of truth for detecting installed developer tools (IDEs, agent runtimes, CLI LLMs) and resolving their configuration paths for **Agent Skills** (`SKILL.md`) and **Model Context Protocol (MCP)** servers.

## Why it exists

In the WebTyp ecosystem, two separate tools require knowledge of installed agent environments:
1. **`webtyp.com/devskills`**: CLI that synchronizes and links domain skills (`SKILL.md`) into each agent's skills directory.
2. **`webtyp.com/app`**: Development runtime that auto-configures MCP endpoints (`http://localhost:<port>/mcp`) into each editor and agent config file.

Without `devharness`, both projects duplicated path resolution, VS Code profile scanning, and extension detection. `devharness` centralizes this logic so adding or updating support for an editor (such as Cursor or Antigravity) happens in **one place**.

---

## Supported Harnesses

| ID | Name | Capabilities | Config Path | Skills Path | Detection Mechanism |
| :--- | :--- | :---: | :--- | :--- | :--- |
| `antigravity-vsc` | Antigravity (VS Code) | Skills, MCP | `~/.gemini/config` | `~/.gemini/config/skills` | Extension `google.google-antigravity` in `.vscode/extensions` or VS Code profiles |
| `antigravity` | Antigravity (Desktop) | Skills, MCP | `~/.gemini/antigravity` | `~/.gemini/antigravity/skills` | Directory `~/.gemini/antigravity` exists |
| `claude` | Claude Code | Skills, MCP | `~` (`.claude.json`) | `~/.claude/skills` | `~/.claude` or `~/.claude.json` exists |
| `gemini` | Gemini CLI | Skills | `~/.gemini` | `~/.gemini/skills` | Directory `~/.gemini` exists |
| `opencode` | OpenCode | Skills, MCP | `~/.config/opencode` | `~/.config/opencode/skills` | Directory `~/.config/opencode` exists |
| `codex` | Codex | Skills | `~/.codex` | `~/.codex/skills` | Directory `~/.codex` exists |
| `qwen` | Qwen | Skills | `~/.qwen` | `~/.qwen/skills` | Directory `~/.qwen` exists |
| `agents` | Agents (Vendor-neutral) | Skills | `~/.agents` | `~/.agents/skills` | Directory `~/.agents` exists |
| `vscode` | Visual Studio Code | MCP | Platform User dir | — | VS Code User config directory exists |
| `cursor` | Cursor | MCP | `~/.cursor` / User dir | — | `~/.cursor` or Cursor User dir exists |

### What is the `agents` harness?

`agents` is **not** a single proprietary commercial application. It represents the **vendor-neutral agent convention** adopted by the open-source community:

- **The Problem**: Proprietary tools store configurations in vendor-locked folders (`~/.claude`, `~/.gemini`, `~/.codex`).
- **The Solution**: Multi-agent runners and open-source agents (such as **OpenCode**) look for global skills in `~/.agents/skills/` (and project-level rules in `.agents/`).
- **In `devharness`**: If `~/.agents` exists on the host machine, `devskills` links all WebTyp skills into `~/.agents/skills/`, allowing any compliant open-source agent to discover them automatically without vendor-specific code.

---

## "I want X → Use Y"

| I want to... | Use... |
| :--- | :--- |
| Get all harnesses that support Agent Skills | `devharness.Skills()` |
| Get all harnesses that support MCP servers | `devharness.MCP()` |
| Detect which skills harnesses are installed on the user's system | `devharness.DetectSkills(homeDir)` |
| Detect which MCP harnesses are installed on the user's system | `devharness.DetectMCP(homeDir)` |
| Find a harness by canonical ID or alias (`"antigravity"` -> `"antigravity-vsc"`) | `devharness.Find(name)` |
| Get the skills directory for a harness | `harness.SkillsDir(homeDir)` |
| Get the configuration directory for a harness | `harness.ConfigDir(homeDir)` |
| Check if the Antigravity VS Code extension is installed | `devharness.IsAntigravityVSCodeInstalled(homeDir)` |
| Resolve all config paths for an editor across VS Code profiles | `devharness.FindMCPConfigPaths(basePath, configFileName)` |

---

## Installation

```bash
go get webtyp.com/devharness
```

---

## Usage Examples

### 1. Discovering Installed Harnesses for Skills Sync

```go
package main

import (
	"fmt"
	"os"

	"webtyp.com/devharness"
)

func main() {
	home, _ := os.UserHomeDir()

	for _, h := range devharness.DetectSkills(home) {
		skillsDir, _ := h.SkillsDir(home)
		fmt.Printf("Detected %s (%s) -> Linking skills into: %s\n", h.Name, h.ID, skillsDir)
	}
}
```

### 2. Configuring MCP for Installed IDEs

```go
package main

import (
	"fmt"
	"os"
	"path/filepath"

	"webtyp.com/devharness"
)

func main() {
	home, _ := os.UserHomeDir()

	for _, h := range devharness.DetectMCP(home) {
		configDir, err := h.ConfigDir(home)
		if err != nil {
			continue
		}

		if h.MCP.SkipProfiles {
			configPath := filepath.Join(configDir, h.MCP.FileName)
			fmt.Printf("Write MCP to: %s\n", configPath)
		} else {
			paths, _ := devharness.FindMCPConfigPaths(configDir, h.MCP.FileName)
			for _, p := range paths {
				fmt.Printf("Write MCP to profile: %s\n", p)
			}
		}
	}
}
```

---

## Zero Dependencies

`devharness` depends exclusively on the Go standard library (`os`, `path/filepath`, `runtime`, `strings`, `errors`). It compiles in milliseconds and introduces zero transitive dependencies to consumers.

## License

MIT
