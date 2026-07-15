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

// Package lockfile handles reading and writing the skmgr.lock file.
package lockfile

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/AbhishekGawade1999/skmgr/internal/types"
	yaml "gopkg.in/yaml.v3"
)

// Read loads the lockfile from the given path.
// If the file does not exist, it returns nil, nil.
func Read(path string) (*types.Lockfile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil // Not an error if lockfile doesn't exist yet
		}
		return nil, fmt.Errorf("reading lockfile: %w", err)
	}

	var lock types.Lockfile
	if err := yaml.Unmarshal(data, &lock); err != nil {
		return nil, fmt.Errorf("parsing lockfile YAML: %w", err)
	}

	return &lock, nil
}

// Write saves the lockfile to the given path.
// It automatically updates the GeneratedAt timestamp and sorts entries by name.
func Write(path string, lock *types.Lockfile) error {
	lock.GeneratedAt = time.Now().UTC().Format(time.RFC3339)

	// Sort entries by name for determinism.
	sort.Slice(lock.Entries, func(i, j int) bool {
		return lock.Entries[i].Name < lock.Entries[j].Name
	})

	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)

	if err := enc.Encode(lock); err != nil {
		return fmt.Errorf("encoding lockfile to YAML: %w", err)
	}

	if err := enc.Close(); err != nil {
		return fmt.Errorf("closing YAML encoder: %w", err)
	}

	if err := os.WriteFile(path, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("writing lockfile: %w", err)
	}

	return nil
}

// HashDirectory computes a deterministic SHA-256 hash of a directory's contents.
// It hashes file paths (relative to root) and their contents, ignoring directories.
// The result is prefixed with "sha256:".
func HashDirectory(root string) (string, error) {
	h := sha256.New()

	var files []string
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// Skip directories, we only hash file contents and their relative paths
		if !d.IsDir() {
			rel, err := filepath.Rel(root, path)
			if err != nil {
				return err
			}
			// Use forward slashes for determinism across OSes
			files = append(files, filepath.ToSlash(rel))
		}
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("walking directory for hash: %w", err)
	}

	// Sort files to ensure deterministic hashing
	sort.Strings(files)

	for _, file := range files {
		// Include file path in hash
		h.Write([]byte(file))
		h.Write([]byte{0}) // separator

		// Include file contents
		fullPath := filepath.Join(root, filepath.FromSlash(file))
		f, err := os.Open(fullPath)
		if err != nil {
			return "", fmt.Errorf("opening file for hash %s: %w", file, err)
		}

		if _, err := io.Copy(h, f); err != nil {
			_ = f.Close()
			return "", fmt.Errorf("hashing file contents %s: %w", file, err)
		}
		_ = f.Close()
	}

	return "sha256:" + hex.EncodeToString(h.Sum(nil)), nil
}
