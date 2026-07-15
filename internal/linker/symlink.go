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

// Package linker manages symlinking and merging skill/rule contents.
package linker

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// createSymlink creates a symlink at linkPath pointing to target.
// If a symlink already exists at linkPath, it is checked. If it points to
// the wrong target, it is replaced. If a regular file or directory exists
// at linkPath, an error is returned.
// It automatically creates any missing parent directories for linkPath.
func createSymlink(target string, linkPath string) error {
	// Ensure parent directories exist
	if err := os.MkdirAll(filepath.Dir(linkPath), 0755); err != nil {
		return fmt.Errorf("creating parent directories for symlink: %w", err)
	}

	// Check if something already exists at linkPath
	info, err := os.Lstat(linkPath)
	if err == nil {
		// Something exists
		if info.Mode()&os.ModeSymlink != 0 {
			// It's a symlink. Does it point to the right place?
			existingTarget, err := os.Readlink(linkPath)
			if err != nil {
				return fmt.Errorf("reading existing symlink: %w", err)
			}
			if existingTarget == target {
				// Already correct, no-op
				return nil
			}
			// Wrong target, remove it so we can recreate it
			if err := os.Remove(linkPath); err != nil {
				return fmt.Errorf("removing incorrect symlink: %w", err)
			}
		} else {
			// It's a regular file or directory, not a symlink. Refuse to overwrite.
			return fmt.Errorf("path %q already exists and is not a symlink", linkPath)
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		// Some other error occurred (e.g. permission denied)
		return fmt.Errorf("statting symlink path: %w", err)
	}

	// Create the symlink
	// Note: Windows requires Developer Mode or Administrator privileges to create symlinks.
	// We rely on standard os.Symlink; if it fails on Windows, it will return an error here.
	if err := os.Symlink(target, linkPath); err != nil {
		return fmt.Errorf("creating symlink: %w", err)
	}

	return nil
}

// isSymlink returns true if the path exists and is a symlink.
func isSymlink(path string) bool {
	info, err := os.Lstat(path)
	if err != nil {
		return false
	}
	return info.Mode()&os.ModeSymlink != 0
}

// removeSymlink safely removes a symlink. It returns an error if the path
// exists but is not a symlink (to prevent accidental deletion of real files).
func removeSymlink(path string) error {
	info, err := os.Lstat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil // Already gone
		}
		return fmt.Errorf("statting path to remove: %w", err)
	}

	if info.Mode()&os.ModeSymlink == 0 {
		return fmt.Errorf("refusing to remove %q: it is not a symlink", path)
	}

	if err := os.Remove(path); err != nil {
		return fmt.Errorf("removing symlink: %w", err)
	}

	return nil
}
