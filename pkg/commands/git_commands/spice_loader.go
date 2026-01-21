package git_commands

import (
	"bufio"
	"bytes"
	"encoding/json"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/common"
)

type SpiceStackLoader struct {
	*common.Common
	spiceCommands *SpiceCommands
}

func NewSpiceStackLoader(
	common *common.Common,
	spiceCommands *SpiceCommands,
) *SpiceStackLoader {
	return &SpiceStackLoader{
		Common:        common,
		spiceCommands: spiceCommands,
	}
}

// Load returns spice stack items, or nil if git-spice is not available
func (self *SpiceStackLoader) Load() ([]*models.SpiceStackItem, error) {
	if !self.spiceCommands.IsAvailable() || !self.spiceCommands.IsInitialized() {
		return nil, nil
	}

	format := self.getEffectiveLogFormat()
	output, err := self.spiceCommands.GetStackBranches(format)
	if err != nil {
		self.Log.Errorf("Failed to get stack branches: %v", err)
		return nil, err
	}

	self.Log.Infof("git-spice output: %s", output)

	branches, err := self.parseBranches(output)
	if err != nil {
		self.Log.Errorf("Failed to parse branches: %v", err)
		return nil, err
	}

	self.Log.Infof("Parsed %d branches", len(branches))

	result := self.buildTree(branches)
	self.Log.Infof("Built tree with %d items", len(result))

	return result, nil
}

// getEffectiveLogFormat returns the log format, checking AppState first, then UserConfig
func (self *SpiceStackLoader) getEffectiveLogFormat() string {
	if format := self.AppState.Spice.LogFormat; format != "" {
		return format
	}
	if format := self.UserConfig().Git.Spice.LogFormat; format != "" {
		return format
	}
	return "short"
}

// parseBranches parses the newline-delimited JSON from gs log
func (self *SpiceStackLoader) parseBranches(output string) ([]*models.SpiceBranchJSON, error) {
	var branches []*models.SpiceBranchJSON
	scanner := bufio.NewScanner(bytes.NewReader([]byte(output)))
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		var branch models.SpiceBranchJSON
		if err := json.Unmarshal([]byte(line), &branch); err != nil {
			self.Log.Warnf("Failed to parse git-spice JSON: %s", err)
			continue
		}
		branches = append(branches, &branch)
	}
	return branches, scanner.Err()
}

// buildTree converts flat list of branches into tree structure with proper depths and last-sibling markers
func (self *SpiceStackLoader) buildTree(branches []*models.SpiceBranchJSON) []*models.SpiceStackItem {
	if len(branches) == 0 {
		return nil
	}

	// Build lookup map
	branchByName := make(map[string]*models.SpiceBranchJSON)
	for _, b := range branches {
		branchByName[b.Name] = b
	}

	// Find roots: branches whose Down is nil or points outside tracked branches
	var roots []string
	for _, b := range branches {
		if b.Down == nil || branchByName[b.Down.Name] == nil {
			roots = append(roots, b.Name)
		}
	}

	// Track visited nodes to prevent duplicates
	visited := make(map[string]bool)

	// DFS to build flat list with proper depths, starting from roots going up
	// Using fliptree algorithm: children are added before parents
	var result []*models.SpiceStackItem
	var dfs func(name string, depth int, siblingIndex int)
	dfs = func(name string, depth int, siblingIndex int) {
		if visited[name] {
			return
		}
		visited[name] = true

		branch, exists := branchByName[name]
		if !exists {
			return
		}

		// FLIPTREE: Traverse children FIRST (before adding parent)
		for i, up := range branch.Ups {
			if upBranch, exists := branchByName[up.Name]; exists {
				dfs(upBranch.Name, depth+1, i)
			}
		}

		// Add the parent node AFTER children have been added
		item := &models.SpiceStackItem{
			Name:         name,
			Depth:        depth,
			SiblingIndex: siblingIndex,
			Current:      branch.Current,
		}

		if branch.Down != nil {
			item.NeedsRestack = branch.Down.NeedsRestack
		}
		if branch.Change != nil {
			item.PRNumber = branch.Change.ID
			item.PRStatus = branch.Change.Status
		}
		if branch.Push != nil {
			item.Ahead = branch.Push.Ahead
			item.Behind = branch.Push.Behind
		}

		result = append(result, item)

		// Add commits if in long format and commits exist (AFTER the branch)
		if self.getEffectiveLogFormat() == "long" && len(branch.Commits) > 0 {
			for _, commit := range branch.Commits {
				commitItem := &models.SpiceStackItem{
					Name:          name, // Keep branch name for context
					Depth:         depth + 1,
					IsCommit:      true,
					CommitSha:     commit.Sha,
					CommitSubject: commit.Subject,
				}
				// Use short SHA (7 chars) if longer
				if len(commitItem.CommitSha) > 7 {
					commitItem.CommitSha = commitItem.CommitSha[:7]
				}
				result = append(result, commitItem)
			}
		}
	}

	for i, root := range roots {
		dfs(root, 0, i)
	}

	// No need to reverse - fliptree naturally produces children-before-parents order
	return result
}
