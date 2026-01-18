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

	format := self.UserConfig().Git.Spice.LogFormat
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

	// Find leaves: branches with no Ups (top of stack)
	var leaves []string
	for _, b := range branches {
		if len(b.Ups) == 0 {
			leaves = append(leaves, b.Name)
		}
	}

	// DFS to build flat list with proper depths, starting from leaves (top) going down
	var result []*models.SpiceStackItem
	var dfs func(name string, depth int, isLast bool)
	dfs = func(name string, depth int, isLast bool) {
		branch, exists := branchByName[name]
		if !exists {
			return
		}

		item := &models.SpiceStackItem{
			Name:    name,
			Depth:   depth,
			IsLast:  isLast,
			Current: branch.Current,
		}

		if branch.Down != nil {
			item.NeedsRestack = branch.Down.NeedsRestack
		}
		if branch.Change != nil {
			item.PRNumber = branch.Change.ID
			item.PRURL = branch.Change.URL
			item.PRStatus = branch.Change.Status
		}
		if branch.Push != nil {
			item.Ahead = branch.Push.Ahead
			item.Behind = branch.Push.Behind
			item.NeedsPush = branch.Push.NeedsPush
		}

		result = append(result, item)

		// Add commits if in long format and commits exist
		if self.UserConfig().Git.Spice.LogFormat == "long" && len(branch.Commits) > 0 {
			for i, commit := range branch.Commits {
				commitItem := &models.SpiceStackItem{
					Name:          name, // Keep branch name for context
					Depth:         depth + 1,
					IsLast:        i == len(branch.Commits)-1 && branch.Down == nil,
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

		// Traverse down (to parent/base)
		if branch.Down != nil {
			if downBranch, exists := branchByName[branch.Down.Name]; exists {
				dfs(downBranch.Name, depth+1, true)
			}
		}
	}

	for i, leaf := range leaves {
		dfs(leaf, 0, i == len(leaves)-1)
	}

	return result
}
