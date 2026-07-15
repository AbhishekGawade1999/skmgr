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

	"github.com/AbhishekGawade1999/skmgr/internal/linker"
	"github.com/AbhishekGawade1999/skmgr/internal/types"
	"github.com/spf13/cobra"
)

var cleanCache bool
var cleanAll bool

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean the project's installed agents and/or the global cache",
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current working directory: %w", err)
		}

		if cleanAll {
			cleanCache = true
		}

		if cleanCache {
			cacheDir := types.CacheDir()
			if err := os.RemoveAll(cacheDir); err != nil {
				return fmt.Errorf("failed to clear cache: %w", err)
			}
			fmt.Printf("Cache cleaned (%s)\n", cacheDir)
		}

		// Only clean the project if --cache wasn't the sole flag specified.
		cleanProject := !cleanCache || cleanAll

		// If --cache is false, cleanProject is true.
		// If --cache is true and --all is false, cleanProject is false.
		// If --all is true, both are true.
		// However, cobra binds to variables, but wait, the way I wrote it:
		// if cleanAll is true, cleanCache is set to true. Then cleanProject = false || true => true.
		// So this logic works.

		if cleanProject {
			agentsDir := filepath.Join(cwd, ".agents")
			if err := os.RemoveAll(agentsDir); err != nil {
				return fmt.Errorf("failed to remove .agents directory: %w", err)
			}

			// Clean broken links left behind by removing .agents
			l := linker.NewLinker()
			if err := l.CleanBrokenLinks(types.ScopeProject, cwd); err != nil {
				return fmt.Errorf("failed to clean broken links: %w", err)
			}
			fmt.Printf("Project cleaned (.agents/ removed and broken symlinks pruned)\n")
		}

		return nil
	},
}

func init() {
	cleanCmd.Flags().BoolVar(&cleanCache, "cache", false, "Clean the global download cache instead of the project")
	cleanCmd.Flags().BoolVar(&cleanAll, "all", false, "Clean both the project and the global cache")
	rootCmd.AddCommand(cleanCmd)
}
