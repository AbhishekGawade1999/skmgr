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
	"text/tabwriter"

	"github.com/AbhishekGawade1999/skmgr/internal/lockfile"
	"github.com/AbhishekGawade1999/skmgr/internal/manifest"
	"github.com/spf13/cobra"
)

var listJson bool

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List installed skills with status",
	RunE: func(cmd *cobra.Command, args []string) error {
		cwd, _ := os.Getwd()
		
		m, err := manifest.Parse(filepath.Join(cwd, "skmgr.yml"))
		if err != nil {
			return fmt.Errorf("failed to read skmgr.yml: %w", err)
		}

		l, _ := lockfile.Read(filepath.Join(cwd, "skmgr.lock"))

		if listJson {
			// Minimal JSON implementation
			fmt.Println("[")
			for i, s := range m.Skills {
				fmt.Printf(`  {"name": "%s", "type": "%s", "scope": "%s"}`, s.Name, s.Type, s.Scope)
				if i < len(m.Skills)-1 {
					fmt.Println(",")
				} else {
					fmt.Println()
				}
			}
			fmt.Println("]")
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "NAME\tTYPE\tSCOPE\tREF\tSTATUS\tTARGETS")

		for _, s := range m.Skills {
			status := "❌ missing"
			if l != nil {
				if entry := l.FindEntry(s.Name); entry != nil {
					status = "✅ current"
				}
			}

			targets := strings.Join(s.EffectiveTargets(m.Targets), ", ")
			
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n", 
				s.Name, 
				s.EffectiveType(), 
				s.EffectiveScope(), 
				s.Ref,
				status,
				targets,
			)
		}
		w.Flush()

		return nil
	},
}

func init() {
	listCmd.Flags().BoolVar(&listJson, "json", false, "Output as JSON")
	rootCmd.AddCommand(listCmd)
}
