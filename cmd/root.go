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

	"github.com/spf13/cobra"
)

// Version is set at build time via -ldflags.
var Version = "dev"

// Global flags accessible to all subcommands.
var (
	globalScope bool
	verbose     bool
	quiet       bool
)

// rootCmd is the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "skmgr",
	Short: "The framework-agnostic skill manager for AI agents",
	Long: `skmgr manages AI agent skills and rules as declarative dependencies
pulled from any git repository. Declare them in skmgr.yml, run skmgr install,
and every team member gets the same agent setup.

Skills are stored canonically in .agents/skills/ and symlinked to each
agent's native directory (.cursor/skills/, .claude/skills/, etc.) to
avoid duplication.`,
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

// Execute runs the root command. Called from main.go.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&globalScope, "global", "g", false, "Operate on global scope (~/.agents/)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "Suppress non-error output")

	rootCmd.Version = Version
	rootCmd.SetVersionTemplate("skmgr version {{.Version}}\n")
}
