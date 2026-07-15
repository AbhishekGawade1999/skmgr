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
	"testing"
)

func TestInit_CreatesManifest(t *testing.T) {
	dir := t.TempDir()
	originalWD, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(originalWD)

	os.MkdirAll(".cursor", 0755)

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
}

func TestInit_AlreadyInitialized(t *testing.T) {
	dir := t.TempDir()
	originalWD, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(originalWD)

	os.WriteFile("skmgr.yml", []byte("version: 1"), 0644)

	cmd := initCmd
	err := cmd.RunE(cmd, []string{})
	if err == nil {
		t.Error("Expected error when skmgr.yml exists")
	}
}
