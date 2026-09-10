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

// SvnRefMapping 表示git-svn配置中的一组路径映射
// 例如branches = branches/proj1/*:refs/remotes/git-svn/branches/*
// SvnPath 为 "branches/proj1/*", RefsPath 为 "refs/remotes/git-svn/branches/*", Type 为 "branches"
type SvnRefMapping struct {
	SvnPath  string // SVN 路径，含通配符，如"branches/proj1/*"
	RefsPath string // Git refs 路径前缀，如”"refs/remotes/git-svn/branches"
	Type     string // 类型："trunk" | "branches" | "tags"
}

type SvnCommands struct {
	*GitCommon
	cmd oscommands.ICmdObjBuilder
	// 缓存 SVN ref 映射，避免重复解析git config
	svnRefMappingsCache *[]SvnRefMapping
	svnUrlCache string
	svnUrlCacheExpiry time.Time
	// CheckBranchStatus 结果缓存（key 为 refType），60秒过期，
	// 避免每次界面刷新都发起 svn list 网络请求；
	// branches 与 tags 在不同 worker协程并发调用，
	// Go map 并发读写会触发不可恢复的 fatal error，故必须加锁
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

// GetSvnRemoteName 从 ref 影射中提取实际的 SVN 远程名。
// 例如 RefsPath 为 “refs/remotes/svn/trunk” 时返回 “svn”。
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

// GetSvnUrl 从git config获取SVN仓库URL
// 返回值如 https://svn.example.com/repo
// 结果缓存60秒
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

// GetSvnRefMappings 解析 svn-remote.svn 配置，返回trunk/branches/tags的refs路径前缀
// 这是识别SVN tags的核心方法，能处理非标准 tags 目录名（如 "release/proj1"）
// 示例配置：
//    [svn-remote "svn"]
//    		url = https://svn.example.com/repo
//    		fetch = trunk/proj1:refs/remotes/git-svn/trunk
//    		branches = branches/proj1/*:refs/remotes/git-svn/branches/*
//    		tags = release/proj1/*:refs/remotes/git-svn/tags/*
// 返回：
//   	{SvnPath: "trunk/proj1", RefsPath: "refs/remotes/git-svn/trunk", Type: "trunk"}
//   	{SvnPath: "branches/proj1/*", RefsPath: "refs/remotes/git-svn/branches", Type: "branches"}
//   	{SvnPath: "release/proj1/*", RefsPath: "refs/remotes/git-svn/tags", Type: "tags"}
func (self *SvnCommands) GetSvnRefMappings() ([]SvnRefMapping, error) {
	if self.svnRefMappingsCache != nil {
		return *self.svnRefMappingsCache, nil
	}

	mappings := []SvnRefMapping{}

	// 1. fetch (trunk) - 格式：trunk/proj1:refs/remotes/git-svn/trunk
	fetchVal, err := self.cmd.New(
		NewGitCmd("config").Arg("--get", "svn-remote.svn.fetch").ToArgv(),
	).DontLog().RunWithOutput()
	if err == nil && strings.TrimSpace(fetchVal) != "" {
		mappings = append(mappings, self.parseSvnRefMapping(strings.TrimSpace(fetchVal), "trunk")...)
	}

	// 2. branches （可能多组） - git config --get-all
	branchesOutput, _ := self.cmd.New(
		NewGitCmd("config").Arg("--get-all", "svn-remote.svn.branches").ToArgv(),
	).DontLog().RunWithOutput()
	for _, line := range strings.Split(strings.TrimSpace(branchesOutput), "\n") {
		if line := strings.TrimSpace(line); line != "" {
			mappings = append(mappings, self.parseSvnRefMapping(line, "branches")...)
		}
	}

	// 3. tags （可能多组）
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

// parseSvnRefMapping 解析单条映射
// 输入格式： "branches/proj1/*:refs/remotes/git-svn/branches/*"
// 输出：SvnPath="branches/proj1*", RefsPath="refs/remotes/git-svn/branches", Type="branches"
func (self *SvnCommands) parseSvnRefMapping(line, defaultType string) []SvnRefMapping {
	parts := strings.Split(line, ":")
	if len(parts) != 2 {
		return nil
	}
	svnPath := strings.TrimSpace(parts[0])
	refsPath := strings.TrimSpace(parts[1])

	// 去掉末尾的 /* 通配符
	svnPathBase := strings.TrimSuffix(svnPath, "/*")
	refsPathBase := strings.TrimSuffix(refsPath, "/*")

	return []SvnRefMapping{{
		SvnPath: svnPathBase,
		RefsPath: refsPathBase,
		Type: defaultType,
	}}
}

// GetTagsRefsPaths 返回所有tags类型的 refs 路径前缀列表
// 用于 TagLoader 扫描 SVN tags
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

// GetSvnUpstream 通过 commit message 中的 git-svn-id 反推本地分支对应的 SVN 远程分支
// 用于 BranchLoader.Load() 循环中回填 upstream 信息
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

// CreateBranch 使用 git svn branch 在 SVN 仓库创建分支
// branchName 取决于 clone 时的 --branches 配置
// 例如 clone 时指定 --branches=branches/proj1，则输入 "xxx" 创建 branches/proj1/xxx
func (self *SvnCommands) CreateBranch(branchName string) error {
	cmdArgs := NewGitCmd("svn").Arg("branch").Arg("-m").Arg(fmt.Sprintf("Create branch %s", branchName)).Arg(branchName).ToArgv()
	return self.cmd.New(cmdArgs).Run()
}

// DeleteServerBranch 从 SVN 服务器删除分支
// branchPath 是相对于 SVN 根的路径，如 "branches/proj1/xxx"
// 实现：执行 svn delete -m "..." <svn-url>/<branchPath>
func (self *SvnCommands) DeleteServerBranch(task gocui.Task, branchPath string) error {
	svnUrl, err := self.GetSvnUrl()
	if err != nil {
		return err
	}
	cmdArgs := []string{"svn", "delete", fmt.Sprintf("%s/%s", svnUrl, branchPath), "-m", fmt.Sprintf("Delete branch %s", branchPath)}
	return self.cmd.New(cmdArgs).PromptOnCredentialRequest(task).Run()
}

// DeleteLocalRef 仅删除本地远程跟踪引用（refs/remotes/<svn-remote>/xxx）
// 不影响 SVN 服务器，安全操作。使用 git update-ref -d 而非 git branch -D -r，
// 因为后者只适用于 branch 格式的 ref 名，不适用于 tags 路径。
func (self *SvnCommands) DeleteLocalRef(refName string) error {
	cmdArgs := NewGitCmd("update-ref").Arg("-d").Arg(refName).ToArgv()
	return self.cmd.New(cmdArgs).Run()
}

// Fetch 执行 git svn fetch -all 获取 SVN 更新
func (self *SvnCommands) Fetch() error {
	cmdArgs := NewGitCmd("svn").Arg("fetch").Arg("--all").ToArgv()
	return self.cmd.New(cmdArgs).Run()
}

// InvalidateStatusCache 清除 CheckBranchStatus 的缓存。
// 应在 git svn fetch 之后调用，确保下次检测使用最新数据。
func (self *SvnCommands) InvalidateStatusCache() {
	self.statusCacheMutex.Lock()
	defer self.statusCacheMutex.Unlock()
	self.statusCache = make(map[string]map[string]models.SvnBranchStatus)
	self.statusCacheExpiry = make(map[string]time.Time)
}

// PruneStaleRefs 删除 SVN 服务器上已不存在的本地远程跟踪引用。
// refType: “branches” 或 “tags”
// 返回已删除的 ref 相对路径列表（与 RemoteBranch.Name 格式一致）。
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
// GetMissingRefs 返回 SVN 服务器上存在但本地未 fetch 的 ref 路径列表。
// refType: “branches” 或 “tags”
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

// CheckBranchStatus 检测本地 refs 和 SVN 服务器的差异
// refType: "branches" 或 "tags"
// 返回值：map[refRelPath]models.SvnBranchStatus, refRelPath 为 ref 相对路
//（完整 ref 去掉 “refs/remotes/git-svn/” 前缀，如 “branches/proj1/xxx”、”tags/R1.0.0”），
// 与 RemoteBranch.Name / TrimPrefix(tag.FullRefName(), “refs/remotes/git-svn/”)格式一致
// 结果缓存 60 秒； svn list 全部失败时返回错误（错误不缓存，下次调用自动重试）径
func (self *SvnCommands) CheckBranchStatus(task gocui.Task, refType string) (map[string]models.SvnBranchStatus, error) {
	// 整个方法加锁： branches/tags 两个调用方在不同 worker 协程，Go map并发读写
	// 会直接触发不可恢复的 fatal error；加锁同时保护函数内
	// GetSvnUrl/GetSvnRefMappings 既有缓存在此路径上的并发访问
	self.statusCacheMutex.Lock()
	defer self.statusCacheMutex.Unlock()

	// 0. 命中缓存直接返回（60 秒内），避免重复发起 svn list 网络请求
	if cached, ok := self.statusCache[refType]; ok && time.Now().Before(self.statusCacheExpiry[refType]) {
		return cached, nil
	}

	svnUrl, err := self.GetSvnUrl()
	if err != nil {
		return nil, err
	}

	mappings, _ := self.GetSvnRefMappings()
	refsPrefix := "refs/remotes/" + self.GetSvnRemoteName() + "/"

	// 1. 获取本地 refs（遍历该类型所有 mapping 的 RefsPath，
	// 兼容refs 端影射到非标准路径的配置，如 tags = tags/*:refs/remotes/svn/releases/*）
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
				// Key 统一为 ref 相对路径
				localRefs[strings.TrimPrefix(line, refsPrefix)] = true
			}
		}
	}
	// 排除属于其他类型 mapping 的 ref，避免误判为 Stale。
	// 例如 branches 的 RefsPath 为 refs/remotes/svn, 会匹配到 trunk 和 tags/*。
	// 但它们分别属于 trunk 和 tags类型，不应出现在 branches 的 localRefs 中。
	// 仅当其他 mapping  的 RefsPath 落在当前 refType 某个 mapping 的 RefsPath 范围内时
	// 才需要排除——因为只有此时 for-each-ref 才可能返回属于其他类型的 ref。
	// 反之若当前 refType 的 RefsPath 更具体（如 tags 的 refs/remotes/svn/tags），
	// for-each-ref 不会返回其他类型的 ref，无需排除。
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
		if isSubPath {
			continue
		}
		for ref := range localRefs {
			fullRef := refsPrefix + ref
			if fullRef == other.RefsPath || strings.HasPrefix(fullRef, other.RefsPath+"/") {
				delete(localRefs, ref)
			}
		}
	}

	// 2. 获取 SVN 服务器上的分支列表（遍历所有 SvnRefMappings)
	svnBranches := make(map[string]bool)
	svnListAttempted := false
	svnListOk := false
	for _, m := range mappings {
		if m.Type != refType {
			continue
		}
		svnListAttempted = true
		// svn list 使用 --non-interactive 防止网络不通时永久阻塞
		svnListOutput, listErr := self.cmd.New(
			[]string{"svn", "list", "--non-interactive", svnUrl+"/"+m.SvnPath},
		).DontLog().RunWithOutput()
		if listErr == nil {
			svnListOk = true
			for _, line := range strings.Split(strings.TrimSpace(svnListOutput), "\n") {
				if line := strings.TrimSpace(line); line !=  "" {
					name := strings.TrimSuffix(line, "/")
					// SVN 路径正向影射为 ref 相对路径，与 localRefs 的 key 格式统一
					// m.RefsPath 去掉通配符后无尾斜杠（如 refs/remotes/svn），
					// 需先补 “/” 再 TrimPrefix refsPrefix （如 refs/remotes/svn/），
					// 否则 TrimPrefix 不匹配，key 变成完整 ref 路径而非相对路径。
					svnBranches[strings.TrimPrefix(m.RefsPath+"/", refsPrefix)+name] = true
				}
			}
		} else {
			self.GitCommon.Log.Warnf("svn list failed for %s: %v", svnUrl+"/"+m.SvnPath, listErr)
		}
	}

	// Svn list failed for all path: svn may be unavailable, network or auth issue.
	// Return empty result with nil error to avoid error dialog.
	// Stale status keeps default (Unknown), does not affect core functionality.
	// Error is not cached, next call will retry automatically.
	if svnListAttempted && !svnListOk {
		self.GitCommon.Log.Warnf("svn list failed for all %s paths, skipping stale detection", refType)
		return map[string]models.SvnBranchStatus{}, nil
	}

	// 3. 对比差异
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

	// 写入缓存（60 秒）。到达此处时 SVN 侧数据必然可信：
	// 若有过 svn list 且全部失败，上方已提前返回
	self.statusCache[refType] = result
	self.statusCacheExpiry[refType] = time.Now().Add(60 * time.Second)

	return result, nil
}


