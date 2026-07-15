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
	"fmt"
	"os"
	"path/filepath"

	"github.com/AbhishekGawade1999/skmgr/internal/manifest"
	"github.com/AbhishekGawade1999/skmgr/internal/types"
	"github.com/spf13/cobra"
)

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

		// Detect targets
		var targets []string
		if _, err := os.Stat(filepath.Join(cwd, ".cursor")); err == nil {
			targets = append(targets, "cursor")
		}
		if _, err := os.Stat(filepath.Join(cwd, ".gemini")); err == nil {
			targets = append(targets, "gemini")
		}
		if _, err := os.Stat(filepath.Join(cwd, ".claude")); err == nil {
			targets = append(targets, "claude-code")
		}
		if _, err := os.Stat(filepath.Join(cwd, ".github")); err == nil {
			targets = append(targets, "copilot")
		}

		if len(targets) == 0 {
			fmt.Println("Warning: No agent directories detected. Creating empty targets list.")
		}

		// Create canonical dirs
		_ = os.MkdirAll(filepath.Join(cwd, ".agents", "skills"), 0750)
		_ = os.MkdirAll(filepath.Join(cwd, ".agents", "rules"), 0750)

		m := &types.Manifest{
			Version: "1",
			Name:    filepath.Base(cwd),
			Targets: targets,
			Skills:  []types.SkillDependency{},
		}

		if err := manifest.Write(manifestPath, m); err != nil {
			return fmt.Errorf("writing manifest: %w", err)
		}

		fmt.Printf("Initialized skmgr.yml for project %s\n", m.Name)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
