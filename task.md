# skmgr — Task Tracker

## Phase 1: Project Scaffolding & Core Types ✅
- [x] Initialize `go.mod` with module path and dependencies (cobra)
- [x] Create `main.go` entry point
- [x] Create `cmd/root.go` — root cobra command with global flags
- [x] Replace LICENSE (MIT → Apache 2.0)
- [x] Verify with `go build ./...` and `go vet ./...`
- [x] Smoke test: `skmgr --version` → `skmgr version dev`
- [x] Smoke test: `skmgr -h` → full help with Usage and Flags

## Phase 2: Internal Domain Types ✅
- [x] Create `internal/types/manifest.go` — Manifest + SkillDependency structs with helpers
- [x] Create `internal/types/lockfile.go` — Lockfile + LockEntry structs with CRUD helpers
- [x] Create `internal/types/config.go` — AgentDef definitions for Cursor, Gemini, Claude Code, Copilot
- [x] Create `internal/types/manifest_test.go` — 17 test cases, all passing
- [x] Verify with `go test -v -race ./internal/types/...` — PASS
- [x] Verify with `go build ./...` and `go vet ./...` — clean
