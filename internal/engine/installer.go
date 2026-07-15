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
	"io"
	"os"
	"path/filepath"

	"github.com/AbhishekGawade1999/skmgr/internal/lockfile"
	"github.com/AbhishekGawade1999/skmgr/internal/types"
)

// Installer handles moving fetched skills into the .agents canonical directories.
type Installer struct {
	projectRoot string
}

// NewInstaller creates a new Installer.
func NewInstaller(projectRoot string) *Installer {
	return &Installer{
		projectRoot: projectRoot,
	}
}

// Install copies the skill from the cache source to its canonical directory and computes the content hash.
func (i *Installer) Install(resolved ResolvedSkill) (string, error) {
	scope := resolved.SkillDependency.EffectiveScope()
	
	var destBase string
	if resolved.SkillDependency.EffectiveType() == types.TypeSkill {
		destBase = types.CanonicalSkillsDir(scope, i.projectRoot)
	} else {
		destBase = types.CanonicalRulesDir(scope, i.projectRoot)
	}

	destDir := filepath.Join(destBase, resolved.SkillDependency.Name)

	// Update-in-place: remove existing if present
	if err := os.RemoveAll(destDir); err != nil {
		return "", fmt.Errorf("removing existing directory %q: %w", destDir, err)
	}

	// Ensure parent dir exists
	if err := os.MkdirAll(destBase, 0755); err != nil {
		return "", fmt.Errorf("creating canonical directory %q: %w", destBase, err)
	}

	// Copy from cache to dest
	if err := copyDir(resolved.SourceDir, destDir); err != nil {
		return "", fmt.Errorf("copying skill from %q to %q: %w", resolved.SourceDir, destDir, err)
	}

	// Compute deterministic content hash
	hash, err := lockfile.HashDirectory(destDir)
	if err != nil {
		return "", fmt.Errorf("computing content hash for %q: %w", destDir, err)
	}

	return hash, nil
}

// CleanOrphans removes any directories in the canonical locations that are not present in the manifest.
func (i *Installer) CleanOrphans(manifest *types.Manifest, scope string) error {
	validNames := make(map[string]bool)
	for _, s := range manifest.Skills {
		if s.EffectiveScope() == scope {
			validNames[s.Name] = true
		}
	}

	cleanDir := func(baseDir string) error {
		entries, err := os.ReadDir(baseDir)
		if err != nil {
			if os.IsNotExist(err) {
				return nil
			}
			return err
		}

		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}
			if !validNames[entry.Name()] {
				// Orphan found, remove it
				orphanPath := filepath.Join(baseDir, entry.Name())
				if err := os.RemoveAll(orphanPath); err != nil {
					return fmt.Errorf("removing orphan %q: %w", orphanPath, err)
				}
			}
		}
		return nil
	}

	if err := cleanDir(types.CanonicalSkillsDir(scope, i.projectRoot)); err != nil {
		return err
	}
	if err := cleanDir(types.CanonicalRulesDir(scope, i.projectRoot)); err != nil {
		return err
	}

	return nil
}

// copyDir recursively copies a directory tree.
func copyDir(src string, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Don't copy .git directories from the cache
		if info.IsDir() && info.Name() == ".git" {
			return filepath.SkipDir
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		
		targetPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(targetPath, info.Mode())
		}

		if info.Mode()&os.ModeSymlink != 0 {
			// Skip symlinks to avoid complexity or infinite loops
			return nil
		}

		// It's a file, copy it
		return copyFile(path, targetPath, info.Mode())
	})
}

func copyFile(src, dst string, mode os.FileMode) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err = io.Copy(out, in); err != nil {
		return err
	}
	return out.Sync()
}
