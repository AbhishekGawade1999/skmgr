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
	"strings"
)

// GetProvider returns the appropriate Provider implementation based on the source URL.
func GetProvider(source string) Provider {
	// If it explicitly starts with file://, it's local.
	if strings.HasPrefix(source, "file://") {
		return &LocalProvider{}
	}

	// If it starts with / or ./ or ../, it's local.
	if strings.HasPrefix(source, "/") || strings.HasPrefix(source, "./") || strings.HasPrefix(source, "../") {
		return &LocalProvider{}
	}

	// For Windows, check for drive letters (e.g., C:\)
	if len(source) > 2 && source[1] == ':' && (source[2] == '\\' || source[2] == '/') {
		return &LocalProvider{}
	}

	// Otherwise, assume it's a git URL (https://, git@, etc.)
	return &GitProvider{}
}
