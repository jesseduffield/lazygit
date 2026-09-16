package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGithubPullRequestCachePath(t *testing.T) {
	stateDir := t.TempDir()
	t.Setenv("CONFIG_DIR", stateDir)

	path, err := githubPullRequestCachePath()

	assert.NoError(t, err)
	assert.Equal(t, filepath.Join(stateDir, githubPullRequestsCacheFileName), path)
}

func TestGithubPullRequestCache(t *testing.T) {
	path := filepath.Join(t.TempDir(), githubPullRequestsCacheFileName)
	cache := newGithubPullRequestCache(path)
	repoOnePullRequests := []CachedPullRequest{{
		HeadRefName:         "first-branch",
		Number:              1,
		Title:               "First pull request",
		State:               "OPEN",
		Url:                 "https://github.com/owner/repo/pull/1",
		HeadRepositoryOwner: "owner",
	}}
	repoTwoPullRequests := []CachedPullRequest{{
		HeadRefName:         "second-branch",
		Number:              2,
		Title:               "Second pull request",
		State:               "MERGED",
		Url:                 "https://github.com/other/repo/pull/2",
		HeadRepositoryOwner: "other",
	}}

	assert.NoError(t, cache.save("/repo/one", repoOnePullRequests))
	assert.NoError(t, cache.save("/repo/two", repoTwoPullRequests))

	content, err := os.ReadFile(path)
	assert.NoError(t, err)
	assert.Equal(t, `{
  "/repo/one": [
    {
      "headRefName": "first-branch",
      "number": 1,
      "title": "First pull request",
      "state": "OPEN",
      "url": "https://github.com/owner/repo/pull/1",
      "headRepositoryOwner": "owner"
    }
  ],
  "/repo/two": [
    {
      "headRefName": "second-branch",
      "number": 2,
      "title": "Second pull request",
      "state": "MERGED",
      "url": "https://github.com/other/repo/pull/2",
      "headRepositoryOwner": "other"
    }
  ]
}
`, string(content))

	reloadedCache := newGithubPullRequestCache(path)
	reloadedCache.load()
	assert.Equal(t, repoOnePullRequests, reloadedCache.get("/repo/one"))
	assert.Equal(t, repoTwoPullRequests, reloadedCache.get("/repo/two"))
	assert.NoError(t, reloadedCache.takeLoadError())
}

func TestGithubPullRequestCacheIgnoresMalformedContent(t *testing.T) {
	path := filepath.Join(t.TempDir(), githubPullRequestsCacheFileName)
	assert.NoError(t, os.WriteFile(path, []byte("{"), 0o644))

	cache := newGithubPullRequestCache(path)
	cache.load()

	assert.ErrorContains(t, cache.takeLoadError(), "parsing GitHub pull request cache")
	assert.Empty(t, cache.get("/repo"))
	assert.NoError(t, cache.save("/repo", []CachedPullRequest{{Number: 1}}))

	reloadedCache := newGithubPullRequestCache(path)
	reloadedCache.load()
	assert.Equal(t, []CachedPullRequest{{Number: 1}}, reloadedCache.get("/repo"))
	assert.NoError(t, reloadedCache.takeLoadError())
}

func TestGithubPullRequestCacheDoesNotModifyAppState(t *testing.T) {
	stateDir := t.TempDir()
	t.Setenv("CONFIG_DIR", stateDir)
	statePath := filepath.Join(stateDir, stateFileName)
	stateContent := []byte("recentrepos:\n  - /repo\n")
	assert.NoError(t, os.WriteFile(statePath, stateContent, 0o644))

	cache := loadGithubPullRequestCache()
	assert.NoError(t, cache.save("/repo", []CachedPullRequest{{Number: 1}}))

	actualStateContent, err := os.ReadFile(statePath)
	assert.NoError(t, err)
	assert.Equal(t, stateContent, actualStateContent)
	_, err = os.Stat(filepath.Join(stateDir, githubPullRequestsCacheFileName))
	assert.NoError(t, err)
}

func TestGithubPullRequestCacheSerializesConcurrentSaves(t *testing.T) {
	path := filepath.Join(t.TempDir(), githubPullRequestsCacheFileName)
	cache := newGithubPullRequestCache(path)
	const repoCount = 20
	var waitGroup sync.WaitGroup
	errs := make(chan error, repoCount)

	for index := range repoCount {
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			repoPath := fmt.Sprintf("/repo/%d", index)
			errs <- cache.save(repoPath, []CachedPullRequest{{Number: index}})
		}()
	}
	waitGroup.Wait()
	close(errs)
	for err := range errs {
		assert.NoError(t, err)
	}

	reloadedCache := newGithubPullRequestCache(path)
	reloadedCache.load()
	for index := range repoCount {
		repoPath := fmt.Sprintf("/repo/%d", index)
		assert.Equal(t, []CachedPullRequest{{Number: index}}, reloadedCache.get(repoPath))
	}
	assert.NoError(t, reloadedCache.takeLoadError())
}
