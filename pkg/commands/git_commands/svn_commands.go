package git_commands

import (
	"fmt"
	"strings"
	"sync"
	"time"
	"github.com/jesseduffield/lazygit/pkg/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
)

// SvnRefMapping represents a path mapping in the git-svn configuration.
// e.g. branches = branches/proj1/*:refs/remotes/git-svn/branches/*
// SvnPath is "branches/proj1/*", RefsPath is "refs/remotes/git-svn/branches/*", Type is "branches"
type SvnRefMapping struct {
	SvnPath  string // SVN path with wildcard, e.g. "branches/proj1/*"
	RefsPath string // Git refs path prefix, e.g. "refs/remotes/git-svn/branches"
	Type     string // Type: "trunk" | "branches" | "tags"
}

type SvnCommands struct {
	*GitCommon
	cmd oscommands.ICmdObjBuilder
	// Cache SVN ref mappings to avoid repeatedly parsing git config
	svnRefMappingsCache *[]SvnRefMapping
	svnUrlCache string
	svnUrlCacheExpiry time.Time
	// CheckBranchStatus result cache (keyed by refType), expires after 60s,
	// to avoid issuing svn list network requests on every UI refresh.
	// branches and tag called from different worker goroutines;
	// Document Go map read/write triggers an unrecoverable fatal error,
	// so a mutex is required.
	statusCache         map[string]map[string]models.SvnBranchStatus
	statusCacheExpiry   map[string]time.Time
	statusCacheMutex    sync.Mutex
}

func NewSvnCommands(gitCommon *GitCommon, cmd oscommands.ICmdObjBuilder) *SvnCommands {
	return &SvnCommands{
		GitCommon: gitCommon,
		cmd: cmd,
		statusCache: make(map[string]map[string]models.SvnBranchStatus),
		statusCacheExpiry: make(map[string]time.Time),
	}
}

// GetSvnRemoteName extracts the actual SVN remote name from ref mappings.
// e.g. returns "svn" when RefsPath is "refs/remotes/svn/trunk".
// Code must not hardcode "git-svn" because the svn-remote name can be arbitrary.
func (self *SvnCommands) GetSvnRemoteName() string {
	mappings, err := self.GetSvnRefMappings()
	if err != nil || len(mappings) == 0 {
		return "svn"
	}

	for _, m := range mappings {
		parts := strings.SplitN(m.RefsPath, "/", 4)
		if len(parts) >= 3 && parts[0] == "refs" && parts[1] == "remotes" {
			return parts[2]
		}
	}
	return "svn"
}

// GetSvnUrl retrieves the SVN repository URL from git config.
// Returns e.g. https://svn.example.com/repo
// Result is cached for 60 seconds.
func (self *SvnCommands) GetSvnUrl() (string, error) {
	if self.svnUrlCache != "" && time.Now().Before(self.svnUrlCacheExpiry) {
		return self.svnUrlCache, nil
	}

	output, err := self.cmd.New(
		NewGitCmd("config").Arg("--get", "svn-remote.svn.url").ToArgv(),
	).DontLog().RunWithOutput()
	if err != nil {
		return "", err
	}
	self.svnUrlCache = strings.TrimSuffix(strings.TrimSpace(output), "/")
	self.svnUrlCacheExpiry = time.Now().Add(60 * time.Second)
	return self.svnUrlCache, nil
}

// GetSvnRefMappings parses the svn-remote.svn config and returns refs path
// prefixes for trunk/branches/tags. This is the core method for identifying
// SVN tags and handles non-standard tag directory names (e.g. "release/proj1").
// Example config:
//    [svn-remote "svn"]
//          url = https://svn.example.com/repo
//          fetch = trunk/proj1:refs/remotes/git-svn/trunk
//          branches = branches/proj1/*:refs/remotes/git-svn/branches/*
//          tags = release/proj1/*:refs/remotes/git-svn/tags/*
// Returns:
//      {SvnPath: "trunk/proj1", RefsPath: "refs/remotes/git-svn/trunk", Type: "trunk"}
//      {SvnPath: "branches/proj1/*", RefsPath: "refs/remotes/git-svn/branches", Type: "branches"}
//      {SvnPath: "release/proj1/*", RefsPath: "refs/remotes/git-svn/tags", Type: "tags"}
func (self *SvnCommands) GetSvnRefMappings() ([]SvnRefMapping, error) {
	if self.svnRefMappingsCache != nil {
		return *self.svnRefMappingsCache, nil
	}

	mappings := []SvnRefMapping{}

	// 1. fetch (trunk) - format: trunk/proj1:refs/remotes/git-svn/trunk
	fetchVal, err := self.cmd.New(
		NewGitCmd("config").Arg("--get", "svn-remote.svn.fetch").ToArgv(),
	).DontLog().RunWithOutput()
	if err == nil && strings.TrimSpace(fetchVal) != "" {
		mappings = append(mappings, self.parseSvnRefMapping(strings.TrimSpace(fetchVal), "trunk")...)
	}

	// 2. branches (possibly multiple groups) - git config --get-all
	branchesOutput, _ := self.cmd.New(
		NewGitCmd("config").Arg("--get-all", "svn-remote.svn.branches").ToArgv(),
	).DontLog().RunWithOutput()
	for _, line := range strings.Split(strings.TrimSpace(branchesOutput), "\n") {
		if line := strings.TrimSpace(line); line != "" {
			mappings = append(mappings, self.parseSvnRefMapping(line, "branches")...)
		}
	}

	// 3. tags (possibly multiple groups)
	tagsOutput, _ := self.cmd.New(
		NewGitCmd("config").Arg("--get-all", "svn-remote.svn.tags").ToArgv(),
	).DontLog().RunWithOutput()
	for _, line := range strings.Split(strings.TrimSpace(tagsOutput), "\n") {
		if line := strings.TrimSpace(line); line != "" {
			mappings = append(mappings, self.parseSvnRefMapping(line, "tags")...)
		}
	}

	self.svnRefMappingsCache = &mappings
	return mappings, nil
}

// parseSvnRefMapping parses a single mapping entry.
// Input format: "branches/proj1/*:refs/remotes/git-svn/branches/*"
// Output: SvnPath="branches/proj1*", RefsPath="refs/remotes/git-svn/branches", Type="branches"
func (self *SvnCommands) parseSvnRefMapping(line, defaultType string) []SvnRefMapping {
	parts := strings.Split(line, ":")
	if len(parts) != 2 {
		return nil
	}
	svnPath := strings.TrimSpace(parts[0])
	refsPath := strings.TrimSpace(parts[1])

	// Strip trailing /* wildcard
	svnPathBase := strings.TrimSuffix(svnPath, "/*")
	refsPathBase := strings.TrimSuffix(refsPath, "/*")

	return []SvnRefMapping{{
		SvnPath: svnPathBase,
		RefsPath: refsPathBase,
		Type: defaultType,
	}}
}

// GetTagsRefsPaths returns the list of refs path prefixes for all tags-type mappings.
// Used by TagLoader to scan SVN tags.
func (self *SvnCommands) GetTagsRefsPaths() ([]string, error) {
	mappings, err := self.GetSvnRefMappings()
	if err != nil {
		return nil, err
	}
	var paths []string
	for _, m := range mappings {
		if m.Type == "tags" {
			paths = append(paths, m.RefsPath)
		}
	}
	return paths, nil
}

// GetSvnUpstream derives the SVN remote branch corresponding to local branch
// via the git-svn-id in the commit message. Used in BranchLoader.Load() to
// backfill upstream information.
func (self *SvnCommands) GetSvnUpstream(branchName string) (string, string, error) {
	output, err := self.cmd.New(
		NewGitCmd("log").
		Arg("--grep=git-svn-id").
		Arg("--format=%B").
		Arg("-1").
		Arg(branchName).
		ToArgv(),
	).DontLog().RunWithOutput()
	if err != nil || strings.TrimSpace(output) == "" {
		return "", "", nil
	}

	svnCommitUrl, ok := self.parseSvnIdLine(output)
	if !ok {
		return "", "", nil
	}

	svnRootUrl, err := self.GetSvnUrl()
	if err != nil {
		return "", "", err
	}

	svnRootUrl = strings.TrimSuffix(svnRootUrl, "/")
	relPath := strings.TrimPrefix(svnCommitUrl, svnRootUrl)
	relPath = strings.TrimPrefix(relPath, "/")
	if relPath == ""  {
		return "", "", nil
	}

	mappings, err := self.GetSvnRefMappings()
	if err != nil {
		return "", "", err
	}

	remoteName := self.GetSvnRemoteName()
	refsPrefix := "refs/remotes/" + remoteName + "/"

	for _, m := range mappings {
		if relPath == m.SvnPath {
			upstreamBranch := strings.TrimPrefix(m.RefsPath, refsPrefix)
			return remoteName, upstreamBranch, nil
		}
		if strings.HasPrefix(relPath, m.SvnPath+"/") {
			remaining := strings.TrimPrefix(relPath, m.SvnPath)
			fullRef := m.RefsPath + remaining
			upstreamBranch := strings.TrimPrefix(fullRef, refsPrefix)
			return remoteName, upstreamBranch, nil
		}
	}

	return "", "", nil
}

func (self *SvnCommands) parseSvnIdLine(commitMessage string) (string, bool) {
	for _, line := range strings.Split(commitMessage, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "git-svn-id: ") {
			rest := strings.TrimPrefix(line, "git-svn-id: ")
			parts := strings.SplitN(rest, " ", 2)
			urlWithRev := parts[0]
			atIdx := strings.LastIndex(urlWithRev, "@")
			if atIdx > 0 {
				return urlWithRev[:atIdx], true
			}
			return urlWithRev, true
		}
	}
	return "", false
}

// CreateBranch creates a branch in the SVN repository using git svn branch.
// branchName depends on the --branches config at clone time.
// e.g. if clone specified --branches=branches/proj1, input "xxx" creates branches/proj1/xxx
func (self *SvnCommands) CreateBranch(branchName string) error {
	cmdArgs := NewGitCmd("svn").Arg("branch").Arg("-m").Arg(fmt.Sprintf("Create branch %s", branchName)).Arg(branchName).ToArgv()
	return self.cmd.New(cmdArgs).Run()
}

// DeleteServerBranch deletes a branch from the SVN server.
// branchPath is relative to the SVN root, e.g. "branches/proj1/xxx"
// Implementation: runs svn delete -m "..." <svn-url>/<branchPath>
func (self *SvnCommands) DeleteServerBranch(task gocui.Task, branchPath string) error {
	svnUrl, err := self.GetSvnUrl()
	if err != nil {
		return err
	}
	cmdArgs := []string{"svn", "delete", fmt.Sprintf("%s/%s", svnUrl, branchPath), "-m", fmt.Sprintf("Delete branch %s", branchPath)}
	return self.cmd.New(cmdArgs).PromptOnCredentialRequest(task).Run()
}

// DeleteLocalRef deletes only the local remote-tracking ref (refs/remotes/<svn-remote>/xxx)
// Does not affect the SVN server; safe operation. Uses git update-ref -d rather
// than git branch -D -r because the latter only works for branch-style ref names,
// not for tags paths.
func (self *SvnCommands) DeleteLocalRef(refName string) error {
	cmdArgs := NewGitCmd("update-ref").Arg("-d").Arg(refName).ToArgv()
	return self.cmd.New(cmdArgs).Run()
}

// Fetch runs git svn fetch --all to retrieve SVN updates
func (self *SvnCommands) Fetch() error {
	self.IndexLock().Lock()
	defer self.IndexLock().Unlock()

	cmdArgs := NewGitCmd("svn").Arg("fetch").Arg("--all").ToArgv()
	return self.cmd.New(cmdArgs).Run()
}

// InvalidateStatusCache clears the CheckBranchStatus cache.
// Should be called after git svn fetch to ensure the next check uses fresh data.
func (self *SvnCommands) InvalidateStatusCache() {
	self.statusCacheMutex.Lock()
	defer self.statusCacheMutex.Unlock()
	self.statusCache = make(map[string]map[string]models.SvnBranchStatus)
	self.statusCacheExpiry = make(map[string]time.Time)
}

// PruneStaleRefs deletes local remote-tracking refs that no longer exist on the SVN server.
// refType: "branches" or "tags"
// Returns the list of deleted ref relative paths (matching RemoteBranch.Name format).
func (self *SvnCommands) PruneStaleRefs(refType string) ([]string, error) {
	statuses, err := self.CheckBranchStatus(nil, refType)
	if err != nil {
		return nil, err
	}
	var pruned []string
	refsPrefix := "refs/remotes/" + self.GetSvnRemoteName() + "/"
	for path, status := range statuses {
		if status == models.SvnBranchStatusStale {
			fullRef := refsPrefix + path
			if err := self.DeleteLocalRef(fullRef); err != nil {
				self.GitCommon.Log.Warnf("failed to prune stale ref %s: %v", fullRef, err)
				continue
			}
			pruned = append(pruned, path)
		}
	}
	return pruned, nil
}
// GetMissingRefs returns the list of ref paths that exist on the SVN server
// but have not been fetched locally.
// refType: "branches" or "tags"
func (self *SvnCommands) GetMissingRefs(refType string) ([]string, error) {
	statuses, err := self.CheckBranchStatus(nil, refType)
	if err != nil {
		return nil, err
	}
	var missing []string
	for path, status := range statuses {
		if status == models.SvnBranchStatusMissing {
			missing = append(missing, path)
		}
	}
	return missing, nil
}

// CheckBranchStatus detects differences between local refs and the SVN server.
// refType: "branches" or "tags"
// Returns map[refRelPath]models.SvnBranchStatus, where refRelPath is the relative
// path (full ref with "refs/remotes/git-svn/" prefix stripped, e.g.
// "branches/proj1/xxx","tags/R1.0.0”), matching RemoteBranch.Name /
// Prefix(tag.FullRefName(), "refs/remotes/git-svn/") format.
// Result is cached for 60s; returns error if all svn list calls fail (error
// is not cached, next call retries automatically).
func (self *SvnCommands) CheckBranchStatus(task gocui.Task, refType string) (map[string]models.SvnBranchStatus, error) {
	// Lock the entire method: branches/tags callers run in different worker
	// goroutines, and concurrent Go map read/write triggers and unrecoverable
	// fatal error. The lock also protects concurrent access to the caches
	// used by GetSvnUrl/GetSvnRefMappings with this method.
	self.statusCacheMutex.Lock()
	defer self.statusCacheMutex.Unlock()

	// 0. Return from cache hit (within 60s) to avoid redundant svn list requests
	if cached, ok := self.statusCache[refType]; ok && time.Now().Before(self.statusCacheExpiry[refType]) {
		return cached, nil
	}

	svnUrl, err := self.GetSvnUrl()
	if err != nil {
		return nil, err
	}

	mappings, _ := self.GetSvnRefMappings()
	refsPrefix := "refs/remotes/" + self.GetSvnRemoteName() + "/"

	// 1. Get local refs (iterate over all mappings of this type's RefsPath,
	// to handle configs where refs are mapped to non-standard paths,
	// e.g. tags = tags/*:refs/remotes/svn/releases/*)
	localRefs := make(map[string]bool)
    for _, m := range mappings {
		if m.Type != refType {
			continue;
		}
		output, refErr := self.cmd.New(
			NewGitCmd("for-each-ref").Arg("--format=%(refname)").Arg(m.RefsPath).ToArgv(),
		).DontLog().RunWithOutput()
		if refErr != nil {
			continue;
		}
		for _, line := range strings.Split(strings.TrimSpace(output), "\n") {
			if line := strings.TrimSpace(line); line != "" {
				// Key is normalized to ref relative path
				localRefs[strings.TrimPrefix(line, refsPrefix)] = true
			}
		}
	}
	// Exlude refs belonging to other mapping types to avoid false Stale detection.
	// e.g. if branches RefsPath is refs/remotes/svn, it matches trunk and tags/*,
	// but those belong to trunk and tags types respectively and should not appear
	// in branches localRefs. Only exclude when another mapping's RefsPath falls
	// within the range of the current refType's mapping RefsPath - because only
	// then could for-each-ref return refs of another type. Conversely, if the
	// current refType's RefsPath is more specific (e.g. tags at
	// ref/remotes/svn/tags), for-each-ref won't return other types' refs, so
	// no exlusion is needed.
	for _, other := range mappings {
		if other.Type == refType {
			continue
		}
		isSubPath := false
		for _, m := range mappings {
			if m.Type == refType && (other.RefsPath == m.RefsPath || strings.HasPrefix(other.RefsPath, m.RefsPath+"/")) {
				isSubPath = true
				break
			}
		}
		if !isSubPath {
			continue
		}
		for ref := range localRefs {
			fullRef := refsPrefix + ref
			if fullRef == other.RefsPath || strings.HasPrefix(fullRef, other.RefsPath+"/") {
				delete(localRefs, ref)
			}
		}
	}

	// 2. Get branch list from SVN server (iterate over all SvnRefMappings)
	svnBranches := make(map[string]bool)
	svnListAttempted := false
	svnListOk := false
	for _, m := range mappings {
		if m.Type != refType {
			continue
		}
		svnListAttempted = true
		// Use --non-interactive to prevent permanent blocking on network issues
		svnListOutput, listErr := self.cmd.New(
			[]string{"svn", "list", "--non-interactive", svnUrl+"/"+m.SvnPath},
		).DontLog().RunWithOutput()
		if listErr == nil {
			svnListOk = true
			for _, line := range strings.Split(strings.TrimSpace(svnListOutput), "\n") {
				if line := strings.TrimSpace(line); line !=  "" {
					name := strings.TrimSuffix(line, "/")
					// Map SVN path to ref relative path, matching localRefs key format.
					// m.RefsPath has no trailing slash after wildcard stripping
					// (e.g. refs/remotes/svn), so prepend "/" bofore TrimPrefix
					// refsPrefix (e.g. refs/remotes/svn/), otherwise TrimPrefix
					// won't match and the key becomes the full ref path instead
					// of the relative path.
					svnBranches[strings.TrimPrefix(m.RefsPath+"/", refsPrefix)+name] = true
				}
			}
		} else {
			self.GitCommon.Log.Warnf("svn list failed for %s: %v", svnUrl+"/"+m.SvnPath, listErr)
		}
	}

	// Svn list failed for all paths: svn may be unavailable, network or auth issue.
	// Return empty result with nil error to avoid error dialog.
	// Stale status keeps default (Unknown), does not affect core functionality.
	// Error is not cached, next call will retry automatically.
	if svnListAttempted && !svnListOk {
		self.GitCommon.Log.Warnf("svn list failed for all %s paths, skipping stale detection", refType)
		return map[string]models.SvnBranchStatus{}, nil
	}

	// 3. Compare differences
	result := make(map[string]models.SvnBranchStatus)
	allPaths := make(map[string]bool)
	for k := range localRefs {
		allPaths[k] = true
	}
	for k := range svnBranches {
		allPaths[k] = true
	}

	for path := range allPaths {
		hasLocal := localRefs[path]
		hasSvn := svnBranches[path]
		var status models.SvnBranchStatus
		if hasLocal && hasSvn {
			status = models.SvnBranchStatusOk
		} else if hasLocal && !hasSvn {
			status = models.SvnBranchStatusStale
		} else if !hasLocal && hasSvn {
			status = models.SvnBranchStatusMissing
		} else {
			status = models.SvnBranchStatusUnknown
		}
		result[path] = status
	}

	// Write to cache (60s). At this point SVN-side data is guaranteed reliable:
	// if any svn list was attempted and all failed, we returned early above.
	self.statusCache[refType] = result
	self.statusCacheExpiry[refType] = time.Now().Add(60 * time.Second)

	return result, nil
}


