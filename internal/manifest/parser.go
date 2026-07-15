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

// Package manifest handles reading, writing, and validating skmgr.yml files.
package manifest

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/AbhishekGawade1999/skmgr/internal/types"
	"gopkg.in/yaml.v3"
)

// ManifestFilenames are the accepted manifest filenames, in priority order.
var ManifestFilenames = []string{"skmgr.yml", "skmgr.yaml"}

// FindManifestPath searches for a manifest file in the given directory.
// Returns the full path to the manifest, or an error if none is found.
func FindManifestPath(dir string) (string, error) {
	for _, name := range ManifestFilenames {
		path := filepath.Join(dir, name)
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}
	return "", fmt.Errorf("no manifest found in %s (expected %s)", dir, strings.Join(ManifestFilenames, " or "))
}

// Parse reads and validates a manifest from the given file path.
func Parse(path string) (*types.Manifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading manifest: %w", err)
	}
	return ParseBytes(data)
}

// ParseBytes parses and validates a manifest from raw YAML bytes.
func ParseBytes(data []byte) (*types.Manifest, error) {
	if len(strings.TrimSpace(string(data))) == 0 {
		return nil, fmt.Errorf("manifest is empty")
	}

	var m types.Manifest
	if err := yaml.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("parsing manifest YAML: %w", err)
	}

	if err := validate(&m); err != nil {
		return nil, err
	}

	return &m, nil
}

// validate checks the manifest for required fields and constraint violations.
func validate(m *types.Manifest) error {
	// Name is required.
	if strings.TrimSpace(m.Name) == "" {
		return fmt.Errorf("manifest validation: 'name' is required")
	}

	// Validate each skill entry.
	seen := make(map[string]bool)
	for i, skill := range m.Skills {
		// Name is required.
		if strings.TrimSpace(skill.Name) == "" {
			return fmt.Errorf("manifest validation: skills[%d] is missing 'name'", i)
		}

		// Name must be unique.
		if seen[skill.Name] {
			return fmt.Errorf("manifest validation: duplicate skill name %q", skill.Name)
		}
		seen[skill.Name] = true

		// Source is required.
		if strings.TrimSpace(skill.Source) == "" {
			return fmt.Errorf("manifest validation: skill %q is missing 'source'", skill.Name)
		}

		// Type must be valid if set.
		if skill.Type != "" && !slices.Contains(types.ValidTypes(), skill.Type) {
			return fmt.Errorf("manifest validation: skill %q has invalid type %q (valid: %s)",
				skill.Name, skill.Type, strings.Join(types.ValidTypes(), ", "))
		}

		// Scope must be valid if set.
		if skill.Scope != "" && !slices.Contains(types.ValidScopes(), skill.Scope) {
			return fmt.Errorf("manifest validation: skill %q has invalid scope %q (valid: %s)",
				skill.Name, skill.Scope, strings.Join(types.ValidScopes(), ", "))
		}
	}

	return nil
}
