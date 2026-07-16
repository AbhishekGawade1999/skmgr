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

package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInit_CreatesManifest(t *testing.T) {
	dir := t.TempDir()
	originalWD, _ := os.Getwd()
	_ = os.Chdir(dir)
	defer func() { _ = os.Chdir(originalWD) }()

	cmd := initCmd
	err := cmd.RunE(cmd, []string{})
	if err != nil {
		t.Fatalf("init failed: %v", err)
	}

	if _, err := os.Stat("skmgr.yml"); os.IsNotExist(err) {
		t.Error("skmgr.yml was not created")
	}
	if _, err := os.Stat(filepath.Join(".agents", "skills")); os.IsNotExist(err) {
		t.Error(".agents/skills was not created")
	}

	content, err := os.ReadFile("skmgr.yml")
	if err != nil {
		t.Fatalf("reading skmgr.yml failed: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "\ntargets:\n") {
		t.Errorf("skmgr.yml should contain uncommented targets block, got:\n%s", contentStr)
	}
	if !strings.Contains(contentStr, "#   - cursor") {
		t.Error("skmgr.yml should contain commented cursor target")
	}
	if strings.Contains(contentStr, "\n  - cursor") {
		t.Error("skmgr.yml should not contain uncommented cursor target")
	}
	if !strings.Contains(contentStr, "skills:") {
		t.Error("skmgr.yml should contain skills block")
	}
	if !strings.Contains(contentStr, "# Importing Specific Skill from Repo") {
		t.Error("skmgr.yml should contain commented skill examples")
	}
}

func TestInit_AlreadyInitialized(t *testing.T) {
	dir := t.TempDir()
	originalWD, _ := os.Getwd()
	_ = os.Chdir(dir)
	defer func() { _ = os.Chdir(originalWD) }()

	_ = os.WriteFile("skmgr.yml", []byte("version: 1"), 0644)

	cmd := initCmd
	err := cmd.RunE(cmd, []string{})
	if err == nil {
		t.Error("Expected error when skmgr.yml exists")
	}
}
