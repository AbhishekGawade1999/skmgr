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
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/AbhishekGawade1999/skmgr/internal/types"
)

// LocalProvider fetches skills from a local file path.
type LocalProvider struct{}

// Fetch implements Provider for local sources.
func (p *LocalProvider) Fetch(skill types.SkillDependency, cacheDir string) (FetchResult, error) {
	// Strip file:// prefix if present
	source := strings.TrimPrefix(skill.Source, "file://")

	// Must be an absolute path or relative to current working dir
	absPath, err := filepath.Abs(source)
	if err != nil {
		return FetchResult{}, fmt.Errorf("resolving local path: %w", err)
	}

	// Verify the directory exists
	info, err := os.Stat(absPath)
	if err != nil {
		return FetchResult{}, fmt.Errorf("accessing local source: %w", err)
	}
	if !info.IsDir() {
		return FetchResult{}, fmt.Errorf("local source %q is not a directory", absPath)
	}

	// For local providers, we don't copy to cache. We just use the path directly.
	// We also don't have a CommitSHA.
	return FetchResult{
		SourceDir: absPath,
		CommitSHA: "",
	}, nil
}
