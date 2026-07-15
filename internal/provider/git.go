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

package provider

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/AbhishekGawade1999/skmgr/internal/types"
)

// GitProvider fetches skills from remote git repositories.
type GitProvider struct{}

// Fetch implements Provider for git sources.
func (p *GitProvider) Fetch(skill types.SkillDependency, cacheDir string) (FetchResult, error) {
	// 1. Determine unique cache directory for this repo.
	// Hash the source URL to avoid path traversal or invalid characters.
	h := sha256.Sum256([]byte(skill.Source))
	hashStr := hex.EncodeToString(h[:])
	repoDir := filepath.Join(cacheDir, "git", hashStr)

	// Ensure parent dir exists
	if err := os.MkdirAll(repoDir, 0755); err != nil {
		return FetchResult{}, fmt.Errorf("creating git cache dir: %w", err)
	}

	// 2. Clone or fetch
	if _, err := os.Stat(filepath.Join(repoDir, ".git")); os.IsNotExist(err) {
		// Clone fresh
		cmd := exec.Command("git", "clone", "--quiet", skill.Source, repoDir)
		if out, err := cmd.CombinedOutput(); err != nil {
			return FetchResult{}, fmt.Errorf("git clone failed: %v\nOutput: %s", err, string(out))
		}
	} else {
		// Fetch latest
		cmd := exec.Command("git", "fetch", "--quiet", "--all", "--tags")
		cmd.Dir = repoDir
		if out, err := cmd.CombinedOutput(); err != nil {
			return FetchResult{}, fmt.Errorf("git fetch failed: %v\nOutput: %s", err, string(out))
		}
	}

	// 3. Checkout the requested ref (or default branch if empty)
	ref := skill.Ref
	if ref == "" {
		// Determine default branch (usually HEAD)
		cmd := exec.Command("git", "symbolic-ref", "refs/remotes/origin/HEAD")
		cmd.Dir = repoDir
		out, err := cmd.Output()
		if err != nil {
			// Fallback to origin/main or origin/master if symbolic-ref fails
			ref = "origin/main"
		} else {
			ref = strings.TrimSpace(string(out))
			// Extract just the branch name, e.g., refs/remotes/origin/main -> origin/main
			ref = strings.TrimPrefix(ref, "refs/remotes/")
		}
	}

	// In case it's a branch, make sure we use the origin/ version to get the latest fetched
	// But it might be a tag or a sha. Try checking out exactly what was requested first.
	// If it's a branch name like 'main', 'origin/main' is safer, but git checkout is smart.
	// We use git checkout --force to discard any local changes.

	checkoutCmd := exec.Command("git", "checkout", "--force", "--quiet", ref)
	checkoutCmd.Dir = repoDir
	if out, err := checkoutCmd.CombinedOutput(); err != nil {
		// If explicit checkout fails, and it wasn't an explicit origin/ ref, try origin/ref
		if skill.Ref != "" && !strings.HasPrefix(skill.Ref, "origin/") {
			fallbackCmd := exec.Command("git", "checkout", "--force", "--quiet", "origin/"+skill.Ref)
			fallbackCmd.Dir = repoDir
			if fbOut, fbErr := fallbackCmd.CombinedOutput(); fbErr != nil {
				return FetchResult{}, fmt.Errorf("git checkout failed for %q: %v\nOutput: %s\nFallback output: %s",
					ref, err, string(out), string(fbOut))
			}
		} else {
			return FetchResult{}, fmt.Errorf("git checkout failed for %q: %v\nOutput: %s", ref, err, string(out))
		}
	}

	// 4. Get the resolved commit SHA
	revCmd := exec.Command("git", "rev-parse", "HEAD")
	revCmd.Dir = repoDir
	revOut, err := revCmd.Output()
	if err != nil {
		return FetchResult{}, fmt.Errorf("git rev-parse failed: %w", err)
	}
	commitSHA := strings.TrimSpace(string(revOut))

	// 5. Construct the final source directory
	sourceDir := repoDir
	if skill.Path != "" {
		sourceDir = filepath.Join(repoDir, filepath.FromSlash(skill.Path))
		// Verify the path exists inside the repo
		if _, err := os.Stat(sourceDir); os.IsNotExist(err) {
			return FetchResult{}, fmt.Errorf("path %q not found in repository at ref %q", skill.Path, ref)
		}
	}

	return FetchResult{
		SourceDir: sourceDir,
		CommitSHA: commitSHA,
	}, nil
}

// GitInstalled checks if the git executable is available in PATH.
func GitInstalled() bool {
	_, err := exec.LookPath("git")
	return err == nil
}
