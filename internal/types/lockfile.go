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

package types

// Lockfile represents the skmgr.lock file.
// It records the exact resolved state of every skill dependency for
// reproducible installs. Commit this file to version control.
type Lockfile struct {
	// Version is the lockfile format version (currently "1").
	Version string `yaml:"version"`

	// GeneratedAt is the ISO 8601 timestamp of when the lockfile was last written.
	GeneratedAt string `yaml:"generated_at"`

	// Entries contains one entry per installed skill/rule.
	Entries []LockEntry `yaml:"entries,omitempty"`
}

// LockEntry records the resolved state of a single skill or rule.
type LockEntry struct {
	// Name matches the SkillDependency.Name from the manifest.
	Name string `yaml:"name"`

	// Source is the git URL or local path the skill was fetched from.
	Source string `yaml:"source"`

	// Path is the subdirectory within the source repo (if monorepo).
	Path string `yaml:"path,omitempty"`

	// CommitSHA is the exact git commit SHA that was checked out.
	// For local sources this may be empty.
	CommitSHA string `yaml:"commit_sha,omitempty"`

	// ContentHash is the SHA-256 hash of the installed skill directory contents.
	// Used to detect local modifications and verify integrity.
	ContentHash string `yaml:"content_hash"`

	// ResolvedAt is the ISO 8601 timestamp of when this entry was resolved.
	ResolvedAt string `yaml:"resolved_at"`
}

// LockfileVersion is the current lockfile format version.
const LockfileVersion = "1"

// NewLockfile creates a new empty Lockfile with the current version.
func NewLockfile() *Lockfile {
	return &Lockfile{
		Version: LockfileVersion,
	}
}

// FindEntry returns the LockEntry for the given skill name, or nil if not found.
func (l *Lockfile) FindEntry(name string) *LockEntry {
	for i := range l.Entries {
		if l.Entries[i].Name == name {
			return &l.Entries[i]
		}
	}
	return nil
}

// SetEntry adds or replaces a LockEntry by name.
func (l *Lockfile) SetEntry(entry LockEntry) {
	for i := range l.Entries {
		if l.Entries[i].Name == entry.Name {
			l.Entries[i] = entry
			return
		}
	}
	l.Entries = append(l.Entries, entry)
}

// RemoveEntry removes the LockEntry with the given name.
// Returns true if an entry was removed, false if not found.
func (l *Lockfile) RemoveEntry(name string) bool {
	for i := range l.Entries {
		if l.Entries[i].Name == name {
			l.Entries = append(l.Entries[:i], l.Entries[i+1:]...)
			return true
		}
	}
	return false
}
