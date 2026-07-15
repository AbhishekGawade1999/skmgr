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

package linker

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

const (
	gitignoreStartDelim = "# skmgr:managed"
	gitignoreEndDelim   = "# skmgr:end"
)

// ensureGitignoreEntry ensures the given path is present in the managed
// section of the project's .gitignore file.
func ensureGitignoreEntry(projectRoot string, entryPath string) error {
	// Normalize path to use forward slashes for .gitignore
	entryPath = filepath.ToSlash(entryPath)
	gitignorePath := filepath.Join(projectRoot, ".gitignore")

	var lines []string
	if data, err := os.ReadFile(gitignorePath); err == nil {
		lines = strings.Split(string(data), "\n")
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("reading .gitignore: %w", err)
	}

	var managedEntries []string
	inManaged := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if trimmed == gitignoreStartDelim {
			inManaged = true
			continue
		}

		if trimmed == gitignoreEndDelim {
			inManaged = false
			continue
		}

		if inManaged {
			if trimmed != "" {
				managedEntries = append(managedEntries, trimmed)
			}
		}
	}

	// Add the new entry if it's not already there
	if !slices.Contains(managedEntries, entryPath) {
		managedEntries = append(managedEntries, entryPath)
	}

	// Sort entries for determinism
	slices.Sort(managedEntries)

	// Reconstruct the file
	// Let's rewrite the logic to be more like merger.go to preserve position
	outputLines := rebuildGitignore(lines, managedEntries)

	output := strings.Join(outputLines, "\n")
	if !strings.HasSuffix(output, "\n") {
		output += "\n"
	}

	if err := os.WriteFile(gitignorePath, []byte(output), 0644); err != nil {
		return fmt.Errorf("writing .gitignore: %w", err)
	}

	return nil
}

// removeGitignoreEntry removes a specific path from the managed section.
// If it was the last path, it removes the entire managed section.
func removeGitignoreEntry(projectRoot string, entryPath string) error {
	entryPath = filepath.ToSlash(entryPath)
	gitignorePath := filepath.Join(projectRoot, ".gitignore")

	data, err := os.ReadFile(gitignorePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("reading .gitignore: %w", err)
	}

	lines := strings.Split(string(data), "\n")
	var managedEntries []string
	inManaged := false

	// First pass: extract managed entries
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == gitignoreStartDelim {
			inManaged = true
			continue
		}
		if trimmed == gitignoreEndDelim {
			inManaged = false
			continue
		}
		if inManaged && trimmed != "" {
			managedEntries = append(managedEntries, trimmed)
		}
	}

	// Remove the target entry
	var keptEntries []string
	removed := false
	for _, e := range managedEntries {
		if e != entryPath {
			keptEntries = append(keptEntries, e)
		} else {
			removed = true
		}
	}

	if !removed {
		return nil // Nothing to do
	}

	outputLines := rebuildGitignore(lines, keptEntries)

	// Clean trailing empty lines
	for len(outputLines) > 0 && outputLines[len(outputLines)-1] == "" {
		outputLines = outputLines[:len(outputLines)-1]
	}

	output := strings.Join(outputLines, "\n")
	if output != "" && !strings.HasSuffix(output, "\n") {
		output += "\n"
	}

	if err := os.WriteFile(gitignorePath, []byte(output), 0644); err != nil {
		return fmt.Errorf("writing .gitignore: %w", err)
	}

	return nil
}

// rebuildGitignore reconstructs the file lines, updating the managed section in place.
// If keptEntries is empty, the managed section is omitted entirely.
func rebuildGitignore(originalLines []string, keptEntries []string) []string {
	var newLines []string
	inManaged := false
	managedFound := false

	for _, line := range originalLines {
		trimmed := strings.TrimSpace(line)

		if trimmed == gitignoreStartDelim {
			inManaged = true
			managedFound = true
			if len(keptEntries) > 0 {
				newLines = append(newLines, gitignoreStartDelim)
				newLines = append(newLines, keptEntries...)
				newLines = append(newLines, gitignoreEndDelim)
			}
			continue
		}

		if trimmed == gitignoreEndDelim {
			inManaged = false
			continue
		}

		if !inManaged {
			newLines = append(newLines, line)
		}
	}

	// If managed section wasn't found and we have entries to add, append them
	if !managedFound && len(keptEntries) > 0 {
		if len(newLines) > 0 && newLines[len(newLines)-1] != "" {
			newLines = append(newLines, "")
		}
		newLines = append(newLines, gitignoreStartDelim)
		newLines = append(newLines, keptEntries...)
		newLines = append(newLines, gitignoreEndDelim)
	}

	return newLines
}
