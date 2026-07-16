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
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/spf13/cobra"
)

const initManifestTemplate = `name: {{ .Name }}
version: "1"

targets:
#   - cursor
#   - gemini
#   - claude-code
#   - copilot

skills:
# # Importing Full Repo Skills
#   - name: anthropics
#     source: https://github.com/anthropics/skills.git
#     ref: main

# # Importing Specific Skill from Repo
#   - name: skill-creator
#     source: https://github.com/anthropics/skills.git
#     path: skills/skill-creator
#     ref: main

# # Importing a Local Skill
#   - name: my-local-skill
#     source: file://./path/to/my/local-skill

# # When selected global scope, it will install that skill globally.
# # By default it's installed project wide only
#   - name: skill-creator
#     source: https://github.com/anthropics/skills.git
#     path: skills/skill-creator
#     ref: main
#     scope: global
`

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new skmgr.yml in the current directory",
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("getting current directory: %w", err)
		}

		manifestPath := filepath.Join(cwd, "skmgr.yml")
		if _, err := os.Stat(manifestPath); err == nil {
			return fmt.Errorf("skmgr.yml already exists")
		}

		// Create canonical dirs
		_ = os.MkdirAll(filepath.Join(cwd, ".agents", "skills"), 0755)
		_ = os.MkdirAll(filepath.Join(cwd, ".agents", "rules"), 0755)

		type templateData struct {
			Name string
		}

		data := templateData{
			Name: filepath.Base(cwd),
		}

		tmpl, err := template.New("manifest").Parse(initManifestTemplate)
		if err != nil {
			return fmt.Errorf("parsing template: %w", err)
		}

		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, data); err != nil {
			return fmt.Errorf("executing template: %w", err)
		}

		if err := os.WriteFile(manifestPath, buf.Bytes(), 0644); err != nil {
			return fmt.Errorf("writing manifest: %w", err)
		}

		fmt.Printf("Initialized skmgr.yml for project %s\n", data.Name)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
