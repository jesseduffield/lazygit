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

	output, err := self.spiceCommands.GetStackBranches()
	if err != nil {
		return nil, err
	}

	branches, err := self.parseBranches(output)
	if err != nil {
		return nil, err
	}

	return self.buildTree(branches), nil
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

	// Build children map: parent -> []children
	children := make(map[string][]string)
	for _, b := range branches {
		if b.Down != nil {
			children[b.Down.Name] = append(children[b.Down.Name], b.Name)
		}
	}

	// Find roots: branches with no Down, or Down pointing to untracked branch (trunk)
	var roots []string
	for _, b := range branches {
		if b.Down == nil {
			roots = append(roots, b.Name)
		} else if _, exists := branchByName[b.Down.Name]; !exists {
			roots = append(roots, b.Name)
		}
	}

	// DFS to build flat list with proper depths
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

		childNames := children[name]
		for i, childName := range childNames {
			dfs(childName, depth+1, i == len(childNames)-1)
		}
	}

	for i, root := range roots {
		dfs(root, 0, i == len(roots)-1)
	}

	return result
}
