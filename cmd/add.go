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
	"strings"

	"github.com/AbhishekGawade1999/skmgr/internal/engine"
	"github.com/AbhishekGawade1999/skmgr/internal/lockfile"
	"github.com/AbhishekGawade1999/skmgr/internal/manifest"
	"github.com/AbhishekGawade1999/skmgr/internal/types"
	"github.com/spf13/cobra"
)

var (
	addName   string
	addPath   string
	addRef    string
	addType   string
	addScope  string
	addTargets []string
)

var addCmd = &cobra.Command{
	Use:   "add <source>",
	Short: "Add a skill/rule to the manifest and install it",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		source := args[0]
		
		cwd, _ := os.Getwd()
		manifestPath := filepath.Join(cwd, "skmgr.yml")

		m, err := manifest.Parse(manifestPath)
		if err != nil {
			return fmt.Errorf("failed to read skmgr.yml (run 'skmgr init' first?): %w", err)
		}

		// Infer name if not provided
		name := addName
		if name == "" {
			if addPath != "" {
				name = filepath.Base(addPath)
			} else {
				name = filepath.Base(source)
				name = strings.TrimSuffix(name, ".git")
			}
		}

		// Check duplicates
		for _, s := range m.Skills {
			if s.Name == name {
				return fmt.Errorf("skill %q already exists in manifest", name)
			}
		}

		skillType := types.TypeSkill
		if addType == "rule" {
			skillType = types.TypeRule
		}

		scope := types.ScopeProject
		if addScope == "global" || globalScope {
			scope = types.ScopeGlobal
		}

		newSkill := types.SkillDependency{
			Name:    name,
			Source:  source,
			Path:    addPath,
			Ref:     addRef,
			Type:    skillType,
			Scope:   scope,
			Targets: addTargets,
		}

		m.Skills = append(m.Skills, newSkill)

		// Install it immediately
		cacheDir := types.CacheDir()
		e := engine.NewEngine(cwd, cacheDir)

		var existingLock *types.Lockfile
		if l, err := lockfile.Read(filepath.Join(cwd, "skmgr.lock")); err == nil {
			existingLock = l
		}

		newLock, err := e.Sync(m, existingLock, false)
		if err != nil {
			return fmt.Errorf("failed to install newly added skill: %w", err)
		}

		// Save manifest and lockfile only if installation succeeded
		if err := manifest.Write(manifestPath, m); err != nil {
			return err
		}
		if err := lockfile.Write(filepath.Join(cwd, "skmgr.lock"), newLock); err != nil {
			return err
		}

		fmt.Printf("Added and installed %s %q\n", skillType, name)
		return nil
	},
}

func init() {
	addCmd.Flags().StringVarP(&addName, "name", "n", "", "Name of the skill (inferred from path/source if omitted)")
	addCmd.Flags().StringVarP(&addPath, "path", "p", "", "Path within the repository")
	addCmd.Flags().StringVar(&addRef, "ref", "main", "Git branch, tag, or commit SHA")
	addCmd.Flags().StringVarP(&addType, "type", "t", "skill", "Type of dependency ('skill' or 'rule')")
	addCmd.Flags().StringVarP(&addScope, "scope", "s", "project", "Installation scope ('project' or 'global')")
	addCmd.Flags().StringSliceVar(&addTargets, "targets", nil, "Specific agent targets (overrides manifest global targets)")
	
	rootCmd.AddCommand(addCmd)
}
