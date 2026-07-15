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
	"bytes"
	"fmt"
	"os"

	"github.com/AbhishekGawade1999/skmgr/internal/types"
	"gopkg.in/yaml.v3"
)

// Write saves the manifest back to the given file path.
// It uses yaml.v3 to encode the struct to YAML.
// Note: yaml.v3 will drop comments when marshalling from a struct.
func Write(path string, m *types.Manifest) error {
	var buf bytes.Buffer
	enc := yaml.NewEncoder(&buf)
	enc.SetIndent(2)

	if err := enc.Encode(m); err != nil {
		return fmt.Errorf("encoding manifest to YAML: %w", err)
	}

	if err := enc.Close(); err != nil {
		return fmt.Errorf("closing YAML encoder: %w", err)
	}

	// Write the file, creating it if it doesn't exist, truncating if it does.
	if err := os.WriteFile(path, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("writing manifest file: %w", err)
	}

	return nil
}
