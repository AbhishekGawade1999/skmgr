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

// Package types defines the core domain types for skmgr.
package types

// Manifest represents the skmgr.yml manifest file.
// It declares the project's skill and rule dependencies along with
// the target agents they should be symlinked to.
type Manifest struct {
	// Name is the project identifier (required).
	Name string `yaml:"name"`

	// Version is an optional project version string.
	Version string `yaml:"version,omitempty"`

	// Targets lists the agent names that skills should be symlinked to
	// (e.g., "cursor", "gemini", "claude-code", "copilot").
	// Individual skills can override this with their own Targets field.
	Targets []string `yaml:"targets,omitempty"`

	// Skills is the list of skill and rule dependencies.
	Skills []SkillDependency `yaml:"skills,omitempty"`
}

// SkillDependency represents a single skill or rule entry in the manifest.
type SkillDependency struct {
	// Name is the local alias and directory name for this skill (required).
	// Must be unique within a manifest.
	Name string `yaml:"name"`

	// Source is the git URL or local path to fetch the skill from (required).
	// Examples:
	//   - https://github.com/user/repo.git
	//   - git@github.com:user/repo.git
	//   - https://gitlab.internal.com/team/skills.git
	//   - file:///Users/me/skills/custom-skill
	Source string `yaml:"source"`

	// Path is the subdirectory within the source repository containing the skill.
	// Used for monorepo support. If empty, the repo root is used.
	Path string `yaml:"path,omitempty"`

	// Ref is the git reference to check out: a tag (v1.2.0), branch (main),
	// or commit SHA (abc123def). If empty, defaults to the repo's default branch.
	Ref string `yaml:"ref,omitempty"`

	// Type is either "skill" or "rule". Defaults to "skill".
	// Skills are stored in .agents/skills/ and rules in .agents/rules/.
	Type string `yaml:"type,omitempty"`

	// Scope is either "project" or "global". Defaults to "project".
	// Project scope installs to .agents/ in the project root.
	// Global scope installs to ~/.agents/ in the user's home directory.
	Scope string `yaml:"scope,omitempty"`

	// Targets overrides the manifest-level Targets for this specific skill.
	// If empty, the manifest-level Targets are used.
	Targets []string `yaml:"targets,omitempty"`
}

// Valid values for SkillDependency.Type.
const (
	TypeSkill = "skill"
	TypeRule  = "rule"
)

// Valid values for SkillDependency.Scope.
const (
	ScopeProject = "project"
	ScopeGlobal  = "global"
)

// ValidTypes returns the set of valid Type values.
func ValidTypes() []string {
	return []string{TypeSkill, TypeRule}
}

// ValidScopes returns the set of valid Scope values.
func ValidScopes() []string {
	return []string{ScopeProject, ScopeGlobal}
}

// EffectiveType returns the Type, defaulting to "skill" if empty.
func (s *SkillDependency) EffectiveType() string {
	if s.Type == "" {
		return TypeSkill
	}
	return s.Type
}

// EffectiveScope returns the Scope, defaulting to "project" if empty.
func (s *SkillDependency) EffectiveScope() string {
	if s.Scope == "" {
		return ScopeProject
	}
	return s.Scope
}

// EffectiveTargets returns the skill's Targets, falling back to the
// manifest-level targets if the skill doesn't override them.
func (s *SkillDependency) EffectiveTargets(manifestTargets []string) []string {
	if len(s.Targets) > 0 {
		return s.Targets
	}
	return manifestTargets
}
