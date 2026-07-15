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

package manifest

import (
	"path/filepath"
	"testing"

	"github.com/AbhishekGawade1999/skmgr/internal/types"
)

func TestWrite(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "skmgr.yml")

	m := &types.Manifest{
		Name:    "test-project",
		Targets: []string{"cursor"},
		Skills: []types.SkillDependency{
			{
				Name:   "s1",
				Source: "src",
			},
		},
	}

	if err := Write(path, m); err != nil {
		t.Fatalf("Write() error: %v", err)
	}

	// Read it back and verify
	parsed, err := Parse(path)
	if err != nil {
		t.Fatalf("Parse() returned error for written manifest: %v", err)
	}

	if parsed.Name != "test-project" {
		t.Errorf("Written Name = %q, want %q", parsed.Name, "test-project")
	}
	if len(parsed.Targets) != 1 || parsed.Targets[0] != "cursor" {
		t.Errorf("Written Targets = %v, want [cursor]", parsed.Targets)
	}
	if len(parsed.Skills) != 1 || parsed.Skills[0].Name != "s1" {
		t.Errorf("Written Skills = %v", parsed.Skills)
	}
}

func TestWrite_ReadOnlyDir(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sub", "skmgr.yml") // sub doesn't exist

	m := &types.Manifest{
		Name: "test",
	}

	err := Write(path, m)
	if err == nil {
		t.Fatal("Write() expected error when directory does not exist, got nil")
	}
}
