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

	validateAndAppend := func(path string) error {
		isSkill := false
		isRule := false

		if _, err := os.Stat(filepath.Join(path, "SKILL.md")); err == nil {
			isSkill = true
		}
		if _, err := os.Stat(filepath.Join(path, "AGENTS.md")); err == nil {
			isRule = true
		}

		if isSkill || isRule {
			parent := filepath.Base(filepath.Dir(path))
			if isSkill && parent != "skills" {
				fmt.Fprintf(os.Stderr, "Warning: skipping SKILL.md in %s (parent directory is not 'skills')\n", path)
				return nil
			}
			if isRule && parent != "rules" {
				fmt.Fprintf(os.Stderr, "Warning: skipping AGENTS.md in %s (parent directory is not 'rules')\n", path)
				return nil
			}
			sourceDirs = append(sourceDirs, path)
		}
		return nil
	}

	if skill.Path != "" && !strings.ContainsAny(skill.Path, "*?[") {
		dir := filepath.Join(basePath, filepath.FromSlash(skill.Path))
		if info, err := os.Stat(dir); err != nil || !info.IsDir() {
			return FetchResult{}, fmt.Errorf("local path %q is not a valid directory", dir)
		}
		sourceDirs = []string{dir}
	} else {
		if skill.Path != "" {
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
				if err := validateAndAppend(match); err != nil {
					return FetchResult{}, err
				}
			}
		} else {
			if info, err := os.Stat(basePath); err != nil || !info.IsDir() {
				return FetchResult{}, fmt.Errorf("local source %q is not a valid directory", basePath)
			}
			err := filepath.WalkDir(basePath, func(path string, d os.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if d.IsDir() {
					if d.Name() == ".git" {
						return filepath.SkipDir
					}
					if err := validateAndAppend(path); err != nil {
						return err
					}
				}
				return nil
			})
			if err != nil {
				return FetchResult{}, fmt.Errorf("auto-discovery failed: %w", err)
			}
		}

		if len(sourceDirs) == 0 {
			if skill.Path == "" {
				sourceDirs = []string{basePath}
			} else {
				return FetchResult{}, fmt.Errorf("no valid skills found matching path %q in local source %q", skill.Path, basePath)
			}
		}
	}

	// For local providers, we don't copy to cache. We just use the path directly.
	// We also don't have a CommitSHA.
	return FetchResult{
		SourceDirs: sourceDirs,
		CommitSHA:  "",
	}, nil
}
