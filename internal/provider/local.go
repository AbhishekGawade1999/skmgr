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

package provider

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/AbhishekGawade1999/skmgr/internal/types"
)

// LocalProvider fetches skills from a local file path.
type LocalProvider struct{}

// Fetch implements Provider for local sources.
func (p *LocalProvider) Fetch(skill types.SkillDependency, cacheDir string) (FetchResult, error) {
	// Strip file:// prefix if present
	source := strings.TrimPrefix(skill.Source, "file://")

	basePath, err := filepath.Abs(source)
	if err != nil {
		return FetchResult{}, fmt.Errorf("resolving local path: %w", err)
	}

	var sourceDirs []string

	if skill.Path != "" {
		if strings.ContainsAny(skill.Path, "*?[") {
			globPattern := filepath.Join(basePath, filepath.FromSlash(skill.Path))
			matches, err := filepath.Glob(globPattern)
			if err != nil {
				return FetchResult{}, fmt.Errorf("invalid glob pattern %q: %w", skill.Path, err)
			}

			for _, match := range matches {
				info, err := os.Stat(match)
				if err != nil || !info.IsDir() {
					continue
				}

				hasSkill := false
				if _, err := os.Stat(filepath.Join(match, "SKILL.md")); err == nil {
					hasSkill = true
				} else if _, err := os.Stat(filepath.Join(match, "AGENTS.md")); err == nil {
					hasSkill = true
				}

				if hasSkill {
					sourceDirs = append(sourceDirs, match)
				}
			}
			if len(sourceDirs) == 0 {
				return FetchResult{}, fmt.Errorf("no valid skills found matching path %q in local source %q", skill.Path, basePath)
			}
		} else {
			dir := filepath.Join(basePath, filepath.FromSlash(skill.Path))
			if info, err := os.Stat(dir); err != nil || !info.IsDir() {
				return FetchResult{}, fmt.Errorf("local path %q is not a valid directory", dir)
			}
			sourceDirs = []string{dir}
		}
	} else {
		if info, err := os.Stat(basePath); err != nil || !info.IsDir() {
			return FetchResult{}, fmt.Errorf("local source %q is not a valid directory", basePath)
		}
		sourceDirs = []string{basePath}
	}

	// For local providers, we don't copy to cache. We just use the path directly.
	// We also don't have a CommitSHA.
	return FetchResult{
		SourceDirs: sourceDirs,
		CommitSHA:  "",
	}, nil
}
