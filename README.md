<div align="center">
  <h1>skmgr</h1>
  <p><b>The framework-agnostic skill manager for AI agents</b></p>

  [![CI](https://github.com/AbhishekGawade1999/skmgr/actions/workflows/ci.yml/badge.svg)](https://github.com/AbhishekGawade1999/skmgr/actions/workflows/ci.yml)
  [![Release](https://img.shields.io/github/v/release/AbhishekGawade1999/skmgr)](https://github.com/AbhishekGawade1999/skmgr/releases)
  [![Go Version](https://img.shields.io/github/go-mod/go-version/AbhishekGawade1999/skmgr)](https://go.dev/)
  [![License](https://img.shields.io/github/license/AbhishekGawade1999/skmgr)](https://github.com/AbhishekGawade1999/skmgr/blob/main/LICENSE)
</div>

---

`skmgr` is a package manager for AI agent skills and rules. It allows you to declare your agent dependencies in a `skmgr.yml` manifest, lock them to specific git commits in a `skmgr.lock` file, and share them across your entire team.

Instead of copy-pasting markdown files or duplicating agent configurations, `skmgr` centralizes them in an `.agents/` directory and intelligently symlinks them into native agent directories like `.cursor/`, `.gemini/`, and `.claude/`.

## ✨ Features

- **Declarative YAML Manifest:** Define your skills and rules in a single `skmgr.yml` file.
- **Deterministic Lockfile:** Reproducible installs across your team using strict git SHAs and deep directory hashing.
- **Any Git Remote:** Pull dependencies from GitHub, GitLab, Bitbucket, or even local file paths.
- **Monorepo Support:** Target specific subdirectories within large repositories using the `path:` directive.
- **Agent-Agnostic:** Works seamlessly with Cursor, Gemini, Claude Code, GitHub Copilot, and more.
- **Centralized Management:** Stores skills in `.agents/` and creates non-destructive symlinks, drastically reducing duplication.
- **Intelligent Rule Merging:** Supports merging single-file markdown rules (like `CLAUDE.md`) using `<!-- skmgr:start -->` delimiters without destroying local edits.
- **Global & Project Scopes:** Install skills per-project (`.agents/`) or globally (`~/.agents/`).

## 🆚 Quick Comparison

| Feature | `skmgr` | APM (Agent Package Manager) | `gh skill` |
|---------|---------|-----------------------------|------------|
| **Primary Use Case** | Cross-agent dependency management | Node.js based agent extensions | GitHub CLI specific skills |
| **Manifest File** | `skmgr.yml` | `package.json` | None |
| **Lockfile Support** | ✅ Yes (`skmgr.lock`) | ✅ Yes | ❌ No |
| **Symlink Strategy** | ✅ Yes (avoids duplication) | ❌ Copies files | ❌ Copies files |
| **Framework Agnostic** | ✅ Yes | ❌ Framework specific | ❌ GitHub specific |
| **Git & Local Sources** | ✅ Yes | ⚠️ Mostly npm | ✅ GitHub only |

---

## 📦 Installation

### macOS (Homebrew)
```bash
brew install AbhishekGawade1999/tap/skmgr
```

### Quick Install (curl | sh)
```bash
curl -fsSL https://raw.githubusercontent.com/AbhishekGawade1999/skmgr/main/scripts/install.sh | sh
```

### Go Install
If you have Go 1.22+ installed:
```bash
go install github.com/AbhishekGawade1999/skmgr@latest
```

### Linux Package Managers (apt / yum)
We host native package repositories via Gemfury.

**Ubuntu / Debian:**
```bash
echo "deb [trusted=yes] https://apt.fury.io/abhishekgawade1999/ /" | sudo tee /etc/apt/sources.list.d/skmgr.list
sudo apt-get update
sudo apt-get install skmgr
```

**CentOS / RHEL / Fedora:**
```bash
echo "[skmgr]
name=skmgr repository
baseurl=https://yum.fury.io/abhishekgawade1999/
enabled=1
gpgcheck=0" | sudo tee /etc/yum.repos.d/skmgr.repo
sudo yum install skmgr
```

### Manual Download
Pre-compiled binaries for Linux, macOS, and Windows are available on the [Releases page](https://github.com/AbhishekGawade1999/skmgr/releases).

---

## 🚀 Quickstart

1. **Initialize `skmgr` in your project:**
   ```bash
   skmgr init
   ```
   *This detects your existing agent directories (like `.cursor/`) and creates a `skmgr.yml` manifest.*

2. **Add a skill:**
   ```bash
   skmgr add https://github.com/anthropics/skills.git --path skills/frontend-design --name frontend-design
   ```

3. **Install skills:**
   ```bash
   skmgr install
   ```

4. **Verify your setup:**
   ```bash
   skmgr list
   ```
   ```text
   NAME               TYPE    SCOPE     REF     STATUS      TARGETS
   frontend-design    skill   project   main    ✅ current   cursor, gemini
   ```

---

## 📄 Manifest Reference (`skmgr.yml`)

The `skmgr.yml` file is the source of truth for your agent dependencies.

```yaml
version: "1"
name: my-awesome-project

# Default agents to symlink skills to
targets:
  - cursor
  - gemini

skills:
  # A standard skill pulled from a git repo
  - name: react-best-practices
    source: https://github.com/my-org/react-skills.git
    ref: v1.2.0

  # A skill pulled from a specific subdirectory (monorepo support)
  - name: backend-api
    source: https://github.com/my-org/monorepo.git
    path: agents/skills/backend
    targets: [cursor] # Override targets just for this skill

  # A global rule pulled from a local directory
  - name: coding-standards
    source: file:///Users/me/local-skills/standards
    type: rule
    scope: global
```

### Fields
- `name` (string): The project identifier.
- `targets` (list): Default agents to symlink to (`cursor`, `gemini`, `claude-code`, `copilot`).
- `skills` (list): The list of dependencies.
  - `name`: Local alias for the skill (must be unique).
  - `source`: Git URL or local `file://` path.
  - `path`: (Optional) Subdirectory within the source repository.
  - `ref`: (Optional) Git branch, tag, or commit SHA. Defaults to default branch.
  - `type`: `skill` (directory of files) or `rule` (single file instructions). Defaults to `skill`.
  - `scope`: `project` (installs to `.agents/`) or `global` (installs to `~/.agents/`). Defaults to `project`.
  - `targets`: (Optional) Overrides the manifest-level targets for this specific skill.

---

## 🔒 Lockfile (`skmgr.lock`)

When you run `skmgr install` or `skmgr update`, `skmgr` generates a `skmgr.lock` file.

- **Purpose:** Records the exact git commit SHAs and computes deep SHA-256 directory hashes to guarantee that every developer on your team gets the exact same files.
- **Commit it:** You should commit `skmgr.lock` to your version control.
- **Frozen Installs:** In CI/CD pipelines or strict environments, run `skmgr install --frozen` to install exactly what is in the lockfile without contacting remotes to resolve refs.

---

## 💻 CLI Reference

| Command | Synopsis |
|---------|----------|
| `skmgr init` | Creates `skmgr.yml` and `.agents/` directories. Auto-detects target agents based on existing folders in your repo. |
| `skmgr add <source>` | Adds a dependency to the manifest and installs it immediately. Use `--name`, `--path`, `--ref`, `--type`, and `--scope` flags to customize. |
| `skmgr remove <name>` | Removes a skill from the manifest and cleanly deletes it from the cache, `.agents/`, and all symlinked target directories. |
| `skmgr install` | Installs all skills defined in `skmgr.yml`. Use `--frozen` to strictly adhere to `skmgr.lock`. |
| `skmgr update [name]` | Updates all skills (or a specific skill) to their latest matching git references, regenerating the lockfile. |
| `skmgr list` | Tabular output of all skills, their current resolution status, and target agents. Use `--json` for programmatic consumption. |
| `skmgr sync` | Quickly synchronizes the `.agents/` and target directories with the state of the manifest without fully re-resolving external refs unless necessary. |

---

## 🔗 How It Works: Symlinks & Architecture

`skmgr` is designed to be non-destructive and avoid file duplication.

1. **Resolution:** Skills are fetched and cached in `~/.skmgr/cache/`.
2. **Installation:** Skills are securely copied to the canonical `.agents/skills/` directory.
3. **Linking:** `skmgr` creates symlinks from the agent-specific directories (e.g., `.cursor/skills/my-skill`) back to the canonical `.agents/` directory.

### `.gitignore` Management
`skmgr` automatically manages your project's `.gitignore` file. When it creates symlinks in `.cursor/` or `.gemini/`, it ensures those symlinks are ignored by git so you don't accidentally commit them. It wraps these rules in safe `### skmgr managed ###` blocks.

### Windows Fallback
On Windows, where symlink creation often requires Developer Mode or Administrator privileges, `skmgr` automatically falls back to creating NTFS Junction Points.

---

## 📜 Rules vs Skills

`skmgr` differentiates between `type: skill` and `type: rule`.

- **Skills (`type: skill`)**: Directories of files (e.g., `SKILL.md`, scripts, examples). These are symlinked as entire directories into the target agent folders.
- **Rules (`type: rule`)**: Single-file instructions (often named `rule.md`).
  
**Intelligent Merging for Rules:**
Some agents, like Claude Code, rely on a single `.claude.md` or `CLAUDE.md` file rather than a directory of skills. For these targets, `skmgr` uses an intelligent merge strategy. It injects the rule content directly into the target file, wrapping it in `<!-- skmgr:start:rule-name -->` and `<!-- skmgr:end:rule-name -->` delimiters. This allows you to have manually written project instructions in `CLAUDE.md` alongside automatically managed rules.

---

## 🤖 Agent Compatibility

| Agent Target | Skill Directory | Rule Strategy |
|--------------|-----------------|---------------|
| `cursor` | `.cursor/rules/` | Symlink file |
| `gemini` | `.gemini/skills/` | Symlink file |
| `claude-code`| `N/A` | Merge into `CLAUDE.md` |
| `copilot` | `.github/copilot-instructions.md` | Merge into file |

*Note: New agents can easily be added to `skmgr` by defining an `AgentDef` in the registry.*

---

## 🔄 CI/CD Usage

`skmgr` is designed for automated environments. To ensure reproducible builds in your pipelines:

```bash
# In your GitHub Actions or CI pipeline
- name: Install Agent Skills
  run: skmgr install --frozen
```

---

## 🤝 Contributing

Contributions are welcome!

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Ensure all tests pass (`go test -v -race ./...`)
4. Ensure the linter is happy (`golangci-lint run`)
5. Commit your changes
6. Open a Pull Request

**Adding a New Agent:** To add a new agent target, simply add a new `AgentDef` struct to `internal/types/config.go` and update the initialization block. No core logic changes required!

---

## ⚖️ License

Copyright 2026 AbhishekGawade1999

Licensed under the Apache License, Version 2.0. See the [LICENSE](LICENSE) file for details.
