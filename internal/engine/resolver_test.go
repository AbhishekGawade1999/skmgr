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
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/AbhishekGawade1999/skmgr/internal/types"
)

func TestResolve_LocalSource(t *testing.T) {
	cache := t.TempDir()
	source := t.TempDir()

	_ = os.WriteFile(filepath.Join(source, "SKILL.md"), []byte("hello"), 0644)

	r := NewResolver(cache)
	skills := []types.SkillDependency{
		{
			Name:   "local-skill",
			Source: "file://" + source,
		},
	}

	res, err := r.Resolve(skills)
	if err != nil {
		t.Fatalf("Resolve failed: %v", err)
	}

	if len(res) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(res))
	}

	if res[0].CommitSHA != "" {
		t.Error("Local sources should not have a commit SHA")
	}
}

func TestResolve_ConflictingNames(t *testing.T) {
	r := NewResolver(t.TempDir())
	skills := []types.SkillDependency{
		{Name: "duplicate"},
		{Name: "duplicate"},
	}

	_, err := r.Resolve(skills)
	if err == nil {
		t.Fatal("Expected error for conflicting names, got nil")
	}
	if !strings.Contains(err.Error(), "duplicate skill name") {
		t.Errorf("Expected duplicate error, got: %v", err)
	}
}

func TestResolve_InvalidPath(t *testing.T) {
	cache := t.TempDir()
	source := t.TempDir()

	r := NewResolver(cache)
	skills := []types.SkillDependency{
		{
			Name:   "local-skill",
			Source: "file://" + source,
			Path:   "nonexistent",
		},
	}
	// For LocalProvider, we don't extract paths the same way in provider. Fetch just verifies `source`.
	// Let's actually test an invalid local path at the source level.
	skills[0].Source = "file://" + filepath.Join(source, "nonexistent")

	_, err := r.Resolve(skills)
	if err == nil {
		t.Fatal("Expected error for invalid path, got nil")
	}
}
