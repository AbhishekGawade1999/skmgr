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

// Package provider abstracts fetching skills from various sources (git, local).
package provider

import "github.com/AbhishekGawade1999/skmgr/internal/types"

// FetchResult contains the outcome of fetching a skill.
type FetchResult struct {
	// SourceDir is the absolute path to the fetched skill contents.
	// For git, this is inside the skmgr cache. For local, it's the target directory itself.
	SourceDir string

	// CommitSHA is the resolved git commit, if applicable.
	// Empty for local sources.
	CommitSHA string
}

// Provider defines the interface for fetching skills from a source.
type Provider interface {
	// Fetch retrieves the skill and returns the local path to its contents.
	// cacheDir is the global ~/.skmgr/cache directory.
	Fetch(skill types.SkillDependency, cacheDir string) (FetchResult, error)
}
