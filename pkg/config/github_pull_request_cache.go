package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

const githubPullRequestsCacheFileName = "github_pull_requests.json"

// CachedPullRequest stores the essential fields of a GitHub pull request.
type CachedPullRequest struct {
	HeadRefName         string `json:"headRefName"`
	Number              int    `json:"number"`
	Title               string `json:"title"`
	State               string `json:"state"`
	ChecksState         string `json:"checksState,omitempty"`
	Url                 string `json:"url"`
	HeadRepositoryOwner string `json:"headRepositoryOwner"`
}

type githubPullRequestCache struct {
	mutex                  sync.Mutex
	path                   string
	pullRequestsByRepoPath map[string][]CachedPullRequest
	loadErr                error
}

func loadGithubPullRequestCache() *githubPullRequestCache {
	path, err := githubPullRequestCachePath()
	if err != nil {
		cache := newGithubPullRequestCache("")
		cache.loadErr = err
		return cache
	}

	cache := newGithubPullRequestCache(path)
	cache.load()
	return cache
}

func githubPullRequestCachePath() (string, error) {
	path, err := stateFilePath(stateFileName)
	if err != nil {
		return "", err
	}

	return filepath.Join(filepath.Dir(path), githubPullRequestsCacheFileName), nil
}

func newGithubPullRequestCache(path string) *githubPullRequestCache {
	return &githubPullRequestCache{
		path:                   path,
		pullRequestsByRepoPath: make(map[string][]CachedPullRequest),
	}
}

func (c *githubPullRequestCache) load() {
	if c.path == "" {
		return
	}

	content, err := os.ReadFile(c.path)
	if err != nil {
		if !os.IsNotExist(err) {
			c.loadErr = fmt.Errorf("reading GitHub pull request cache: %w", err)
		}
		return
	}
	if len(content) == 0 {
		return
	}

	if err := json.Unmarshal(content, &c.pullRequestsByRepoPath); err != nil {
		c.pullRequestsByRepoPath = make(map[string][]CachedPullRequest)
		c.loadErr = fmt.Errorf("parsing GitHub pull request cache: %w", err)
	} else if c.pullRequestsByRepoPath == nil {
		c.pullRequestsByRepoPath = make(map[string][]CachedPullRequest)
	}
}

func (c *githubPullRequestCache) get(repoPath string) []CachedPullRequest {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return append([]CachedPullRequest(nil), c.pullRequestsByRepoPath[repoPath]...)
}

// takeLoadError returns the error, if any, that occurred while loading the
// cache from disk, clearing it so that it is reported only once.
func (c *githubPullRequestCache) takeLoadError() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	loadErr := c.loadErr
	c.loadErr = nil
	return loadErr
}

func (c *githubPullRequestCache) save(repoPath string, pullRequests []CachedPullRequest) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.pullRequestsByRepoPath[repoPath] = append([]CachedPullRequest(nil), pullRequests...)
	if c.path == "" {
		return nil
	}

	content, err := json.MarshalIndent(c.pullRequestsByRepoPath, "", "  ")
	if err != nil {
		return err
	}
	content = append(content, '\n')

	// Apparently when people have read-only permissions they prefer us to fail
	// silently, so don't propagate permission errors.
	if err := os.WriteFile(c.path, content, 0o644); err != nil && !os.IsPermission(err) {
		return err
	}

	return nil
}
