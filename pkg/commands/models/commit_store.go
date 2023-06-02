package models

import (
	"sync"

	"github.com/jesseduffield/generics/slices"
	"github.com/samber/lo"
)

// This stores immutable commits as a graph for fast lookup
type CommitStore struct {
	commitsMap map[string]ImmutableCommit
	mutex      *sync.RWMutex
}

func NewCommitStore() *CommitStore {
	return &CommitStore{commitsMap: make(map[string]ImmutableCommit), mutex: &sync.RWMutex{}}
}

func (self *CommitStore) Add(commit ImmutableCommit) {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	self.commitsMap[commit.Hash()] = commit
}

func (self *CommitStore) AddSlice(commit []ImmutableCommit) {
	self.mutex.Lock()
	defer self.mutex.Unlock()

	for _, commit := range commit {
		self.commitsMap[commit.Hash()] = commit
	}
}

func (self *CommitStore) GetCommit(hash string) (ImmutableCommit, bool) {
	self.mutex.RLock()
	defer self.mutex.RUnlock()

	commit, ok := self.commitsMap[hash]
	return commit, ok
}

func (self *CommitStore) GetParents(hash string) ([]ImmutableCommit, bool) {
	commit, ok := self.GetCommit(hash)
	if !ok {
		return nil, false
	}

	parents := make([]ImmutableCommit, len(commit.ParentHashes()))
	for i, parentHash := range commit.ParentHashes() {
		parents[i], ok = self.GetCommit(parentHash)
		if !ok {
			return nil, false
		}
	}

	return parents, true
}

type IsAncestorResponse int

const (
	IsAncestorResponseYes IsAncestorResponse = iota
	IsAncestorResponseNo
	// returned when we can't find the commit, or if one of the ancestors along the chain can't be found
	IsAncestorResponseUnknown
)

var IsAncestorResponseStrings = []string{
	"yes",
	"no",
	"unknown",
}

func (self *CommitStore) IsAncestor(hash string, ancestorHash string) IsAncestorResponse {
	if hash == ancestorHash {
		return IsAncestorResponseYes
	}

	commit, ok := self.GetCommit(hash)
	if !ok {
		return IsAncestorResponseUnknown
	}

	if commit.IsRoot() {
		return IsAncestorResponseNo
	}

	parentHashes := commit.ParentHashes()

	// first check the parent hashes themselves: spares us attempting a lookup of the actual parent structs
	for _, parentHash := range parentHashes {
		if parentHash == ancestorHash {
			return IsAncestorResponseYes
		}
	}

	unknown := false
	for _, parentHash := range parentHashes {
		response := self.IsAncestor(parentHash, ancestorHash)
		if response == IsAncestorResponseYes {
			return IsAncestorResponseYes
		}
		if response == IsAncestorResponseUnknown {
			unknown = true
		}
	}

	if unknown {
		return IsAncestorResponseUnknown
	}

	return IsAncestorResponseNo
}

// used for testing
func (self *CommitStore) Slice() []ImmutableCommit {
	self.mutex.RLock()
	defer self.mutex.RUnlock()

	commits := lo.Values(self.commitsMap)

	// sort by hash for deterministic result
	slices.SortFunc(commits, func(a, b ImmutableCommit) bool {
		return a.Hash() < b.Hash()
	})

	return commits
}

func (self *CommitStore) Size() int {
	return len(self.commitsMap)
}

func (self *CommitStore) dfs(hash string, visited map[string]bool) {
	visited[hash] = true
	commit, ok := self.GetCommit(hash)
	if !ok {
		return
	}

	// Traverse the parent hashes
	for _, parentHash := range commit.ParentHashes() {
		if !visited[parentHash] {
			self.dfs(parentHash, visited)
		}
	}
}

// returns subset of candidates which are ancestors of the given hash
func (self *CommitStore) FindAncestors(hash string, candidates []string) map[string]bool {
	visited := make(map[string]bool)
	self.dfs(hash, visited)

	ancestors := map[string]bool{}
	for _, commitHash := range candidates {
		if visited[commitHash] {
			ancestors[commitHash] = true
		}
	}

	return ancestors
}
