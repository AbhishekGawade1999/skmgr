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

package lockfile

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/AbhishekGawade1999/skmgr/internal/types"
)

func TestRead_Valid(t *testing.T) {
	path := filepath.Join("..", "..", "testdata", "lockfiles", "valid.lock")
	lock, err := Read(path)
	if err != nil {
		t.Fatalf("Read() error: %v", err)
	}
	if lock == nil {
		t.Fatal("Read() returned nil lockfile")
	}

	if lock.Version != "1" {
		t.Errorf("Version = %q, want 1", lock.Version)
	}
	if len(lock.Entries) != 1 {
		t.Fatalf("Entries length = %d, want 1", len(lock.Entries))
	}

	entry := lock.Entries[0]
	if entry.Name != "test-skill" {
		t.Errorf("Entry Name = %q, want test-skill", entry.Name)
	}
	if entry.CommitSHA != "abc123def456" {
		t.Errorf("Entry CommitSHA = %q", entry.CommitSHA)
	}
}

func TestRead_NotExist(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "does_not_exist.lock")

	lock, err := Read(path)
	if err != nil {
		t.Fatalf("Read() for non-existent file should return nil error, got: %v", err)
	}
	if lock != nil {
		t.Fatalf("Read() for non-existent file should return nil lockfile, got: %v", lock)
	}
}

func TestWriteAndRead(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "skmgr.lock")

	lock := &types.Lockfile{
		Version: types.LockfileVersion,
		Entries: []types.LockEntry{
			{Name: "b-skill", Source: "src-b"},
			{Name: "a-skill", Source: "src-a"},
		},
	}

	if err := Write(path, lock); err != nil {
		t.Fatalf("Write() error: %v", err)
	}

	// Read it back
	readLock, err := Read(path)
	if err != nil {
		t.Fatalf("Read() error: %v", err)
	}

	if readLock.GeneratedAt == "" {
		t.Error("GeneratedAt was not populated by Write()")
	}

	// Verify entries were sorted by name
	if len(readLock.Entries) != 2 {
		t.Fatalf("Entries length = %d", len(readLock.Entries))
	}
	if readLock.Entries[0].Name != "a-skill" || readLock.Entries[1].Name != "b-skill" {
		t.Errorf("Entries were not sorted correctly: %v", readLock.Entries)
	}
}

func TestHashDirectory(t *testing.T) {
	dir := t.TempDir()

	// Create some files
	if err := os.WriteFile(filepath.Join(dir, "file1.txt"), []byte("content1"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(dir, "sub"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "sub", "file2.txt"), []byte("content2"), 0644); err != nil {
		t.Fatal(err)
	}

	hash1, err := HashDirectory(dir)
	if err != nil {
		t.Fatalf("HashDirectory() error: %v", err)
	}
	if hash1 == "" {
		t.Fatal("HashDirectory() returned empty hash")
	}

	// Verify deterministic: run again, should match
	hash2, err := HashDirectory(dir)
	if err != nil {
		t.Fatalf("HashDirectory() error: %v", err)
	}
	if hash1 != hash2 {
		t.Errorf("Hash is not deterministic: %q != %q", hash1, hash2)
	}

	// Modify a file, hash should change
	if err := os.WriteFile(filepath.Join(dir, "file1.txt"), []byte("changed"), 0644); err != nil {
		t.Fatal(err)
	}
	hash3, err := HashDirectory(dir)
	if err != nil {
		t.Fatalf("HashDirectory() error: %v", err)
	}
	if hash1 == hash3 {
		t.Error("Hash did not change after modifying file contents")
	}
}

func TestHashDirectory_NotExist(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "does_not_exist")

	_, err := HashDirectory(path)
	if err == nil {
		t.Fatal("HashDirectory() expected error for non-existent directory, got nil")
	}
}
