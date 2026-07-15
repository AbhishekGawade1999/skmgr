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

// Package engine orchestrates the lifecycle of fetching, linking, and locking skills.
package engine

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/AbhishekGawade1999/skmgr/internal/provider"
	"github.com/AbhishekGawade1999/skmgr/internal/types"
)

// Resolver handles fetching and resolving concrete SHAs for skills.
type Resolver struct {
	cacheDir string
}

// NewResolver creates a new Resolver.
func NewResolver(cacheDir string) *Resolver {
	return &Resolver{
		cacheDir: cacheDir,
	}
}

// ResolvedSkill contains the fetched source directory and commit SHA.
type ResolvedSkill struct {
	SkillDependency types.SkillDependency
	SourceDir       string
	CommitSHA       string
	ContentHash     string // Calculated later by installer
}

// Resolve processes a list of skill dependencies, fetches them, and returns resolved information.
// If a lockfile is provided and a skill matches exactly (same source and ref), it can skip fetch
// if --frozen is used, but for now we always fetch to ensure we have the source.
// Detects name conflicts before fetching.
func (r *Resolver) Resolve(skills []types.SkillDependency) ([]ResolvedSkill, error) {
	var resolved []ResolvedSkill
	seen := make(map[string]bool)

	for _, skill := range skills {
		// Get appropriate provider
		prov := provider.GetProvider(skill.Source)

		// Let the provider handle the path validation
		res, err := prov.Fetch(skill, r.cacheDir)
		if err != nil {
			// Check if it's a path not found error from git provider
			if strings.Contains(err.Error(), "not found in repository") {
				return nil, fmt.Errorf("invalid path for skill %q: %w", skill.Name, err)
			}
			return nil, fmt.Errorf("failed to fetch skill %q: %w", skill.Name, err)
		}

		for _, dir := range res.SourceDirs {
			finalName := skill.Name

			// If name is empty, or if this fetch yielded multiple directories (meaning a wildcard was used),
			// we must generate a unique name for each directory.
			if finalName == "" || len(res.SourceDirs) > 1 {
				base := filepath.Base(dir)
				if skill.Name != "" {
					finalName = fmt.Sprintf("%s-%s", skill.Name, base)
				} else {
					finalName = base
				}
			}

			// If it's a single directory and Name was empty, but path didn't have wildcard... Wait, parser requires Name if no wildcard.
			// So this handles both single and multiple correctly.

			if seen[finalName] {
				return nil, fmt.Errorf("duplicate skill name resolved: %q", finalName)
			}
			seen[finalName] = true

			resolvedDep := skill
			resolvedDep.Name = finalName

			// Auto-detect type if a wildcard was used (meaning we might be mixing rules and skills).
			// If both or neither are present, we fall back to whatever was specified in the config.
			if len(res.SourceDirs) > 1 || strings.ContainsAny(skill.Path, "*?[") {
				hasSkill := false
				if _, err := os.Stat(filepath.Join(dir, "SKILL.md")); err == nil {
					hasSkill = true
				}
				hasRule := false
				if _, err := os.Stat(filepath.Join(dir, "AGENTS.md")); err == nil {
					hasRule = true
				}

				if hasRule && !hasSkill {
					resolvedDep.Type = types.TypeRule
				} else if hasSkill && !hasRule {
					resolvedDep.Type = types.TypeSkill
				}
			}

			resolved = append(resolved, ResolvedSkill{
				SkillDependency: resolvedDep,
				SourceDir:       dir,
				CommitSHA:       res.CommitSHA,
			})
		}
	}

	return resolved, nil
}
