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

	"github.com/AbhishekGawade1999/skmgr/internal/engine"
	"github.com/AbhishekGawade1999/skmgr/internal/lockfile"
	"github.com/AbhishekGawade1999/skmgr/internal/manifest"
	"github.com/AbhishekGawade1999/skmgr/internal/types"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove <name>",
	Short: "Remove a skill/rule from the manifest and uninstall it",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		
		cwd, _ := os.Getwd()
		manifestPath := filepath.Join(cwd, "skmgr.yml")

		m, err := manifest.Parse(manifestPath)
		if err != nil {
			return fmt.Errorf("failed to read skmgr.yml: %w", err)
		}

		// Find it
		idx := -1
		var removedSkill types.SkillDependency
		for i, s := range m.Skills {
			if s.Name == name {
				idx = i
				removedSkill = s
				break
			}
		}

		if idx == -1 {
			return fmt.Errorf("skill %q not found in manifest", name)
		}

		// Remove from manifest
		m.Skills = append(m.Skills[:idx], m.Skills[idx+1:]...)

		// Uninstall it
		cacheDir := types.CacheDir()
		e := engine.NewEngine(cwd, cacheDir)

		targets := removedSkill.EffectiveTargets(m.Targets)
		if err := e.Remove(name, removedSkill.EffectiveScope(), targets); err != nil {
			return fmt.Errorf("failed to uninstall skill: %w", err)
		}

		// Save manifest
		if err := manifest.Write(manifestPath, m); err != nil {
			return err
		}

		// Update lockfile
		lockPath := filepath.Join(cwd, "skmgr.lock")
		if l, err := lockfile.Read(lockPath); err == nil && l != nil {
			l.RemoveEntry(name)
			_ = lockfile.Write(lockPath, l)
		}

		// Clean orphans to remove from .agents/
		_ = e.Installer.CleanOrphans(m, types.ScopeProject)
		_ = e.Installer.CleanOrphans(m, types.ScopeGlobal)

		fmt.Printf("Removed %q\n", name)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
