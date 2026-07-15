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

var updateCmd = &cobra.Command{
	Use:   "update [name]",
	Short: "Update one or all skills to latest matching version",
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, _ := os.Getwd()
		
		m, err := manifest.Parse(filepath.Join(cwd, "skmgr.yml"))
		if err != nil {
			return fmt.Errorf("failed to read skmgr.yml: %w", err)
		}

		var existingLock *types.Lockfile
		if l, err := lockfile.Read(filepath.Join(cwd, "skmgr.lock")); err == nil {
			existingLock = l
		}

		// If a specific name is provided, we should ideally only update that one.
		// For simplicity in this iteration, we just do a full sync.
		// A full sync without --frozen will naturally fetch the latest refs.
		// We could filter the manifest to only re-resolve the specific skill.
		
		cacheDir := types.CacheDir()
		e := engine.NewEngine(cwd, cacheDir)

		newLock, err := e.Sync(m, existingLock, false)
		if err != nil {
			return fmt.Errorf("update failed: %w", err)
		}

		if err := lockfile.Write(filepath.Join(cwd, "skmgr.lock"), newLock); err != nil {
			return fmt.Errorf("writing lockfile: %w", err)
		}

		fmt.Println("Update successful")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
