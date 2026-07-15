// Copyright 2026 AbhishekGawade1999
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package engine

import (
	"fmt"
	"time"

	"github.com/AbhishekGawade1999/skmgr/internal/linker"
	"github.com/AbhishekGawade1999/skmgr/internal/types"
)

// Engine is the central orchestrator for skmgr commands.
type Engine struct {
	ProjectRoot string
	CacheDir    string
	Resolver    *Resolver
	Installer   *Installer
	Linker      *linker.Linker
}

// NewEngine creates a new Engine.
func NewEngine(projectRoot, cacheDir string) *Engine {
	return &Engine{
		ProjectRoot: projectRoot,
		CacheDir:    cacheDir,
		Resolver:    NewResolver(cacheDir),
		Installer:   NewInstaller(projectRoot),
		Linker:      linker.NewLinker(),
	}
}

// Sync represents the core `skmgr install` flow.
func (e *Engine) Sync(manifest *types.Manifest, existingLock *types.Lockfile, frozen bool) (*types.Lockfile, error) {
	if frozen && existingLock == nil {
		return nil, fmt.Errorf("frozen mode requested but no lockfile exists")
	}

	// 1. Resolve SHAs and fetch sources
	// In a complete implementation, frozen mode would skip fetching if the hash matches.
	resolvedSkills, err := e.Resolver.Resolve(manifest.Skills)
	if err != nil {
		return nil, fmt.Errorf("resolution failed: %w", err)
	}

	newLock := &types.Lockfile{
		Version:     "1",
		GeneratedAt: time.Now().UTC().Format(time.RFC3339),
	}

	// 2. Install skills and link them
	for _, rs := range resolvedSkills {
		// Clean up existing symlinks for all targets just in case we are updating targets
		e.cleanupLinks(rs.SkillDependency.Name, rs.SkillDependency.EffectiveScope(), manifest.Targets)

		// Install to canonical dir
		hash, err := e.Installer.Install(rs)
		if err != nil {
			return nil, fmt.Errorf("install failed for %q: %w", rs.SkillDependency.Name, err)
		}
		rs.ContentHash = hash

		// Create symlinks/merges
		targets := rs.SkillDependency.EffectiveTargets(manifest.Targets)
		for _, target := range targets {
			if rs.SkillDependency.EffectiveType() == types.TypeSkill {
				if err := e.Linker.LinkSkill(target, rs.SkillDependency.Name, rs.SkillDependency.EffectiveScope(), e.ProjectRoot); err != nil {
					return nil, fmt.Errorf("linking skill %q for %q: %w", rs.SkillDependency.Name, target, err)
				}
			} else {
				if err := e.Linker.LinkRule(target, rs.SkillDependency.Name, rs.SkillDependency.EffectiveScope(), e.ProjectRoot); err != nil {
					return nil, fmt.Errorf("linking rule %q for %q: %w", rs.SkillDependency.Name, target, err)
				}
			}
		}

		// Add to lockfile
		newLock.SetEntry(types.LockEntry{
			Name:        rs.SkillDependency.Name,
			CommitSHA:   rs.CommitSHA,
			ContentHash: rs.ContentHash,
		})
	}

	// 3. Clean up orphans (skills in .agents/ that are no longer in the manifest)
	if err := e.Installer.CleanOrphans(manifest, types.ScopeProject); err != nil {
		return nil, fmt.Errorf("cleaning project orphans: %w", err)
	}
	if err := e.Installer.CleanOrphans(manifest, types.ScopeGlobal); err != nil {
		return nil, fmt.Errorf("cleaning global orphans: %w", err)
	}

	return newLock, nil
}

// Remove completely uninstalls a skill.
func (e *Engine) Remove(skillName string, scope string, targets []string) error {
	// 1. Unlink
	e.cleanupLinks(skillName, scope, targets)
	return nil
}

func (e *Engine) cleanupLinks(name string, scope string, targets []string) {
	for _, target := range targets {
		// Try unlinking as both skill and rule to be safe, since we might be converting
		_ = e.Linker.UnlinkSkill(target, name, scope, e.ProjectRoot)
		_ = e.Linker.UnlinkRule(target, name, scope, e.ProjectRoot)
	}
}
