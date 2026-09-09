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

	for _, m := range mappings {
		if relPath == m.SvnPath {
			upstreamBranch := strings.TrimPrefix(m.RefsPath, "refs/remotes/git-svn/")
			return "git-svn", upstreamBranch, nil
		}
		if strings.HasPrefix(relPath, m.SvnPath+"/") {
			remaining := strings.TrimPrefix(relPath, m.SvnPath)
			fullRef := m.RefsPath + remaining
			upstreamBranch := strings.TrimPrefix(fullRef, "refs/remotes/git-svn/")
			return "git-svn", upstreamBranch, nil
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

// DeleteLocalRef 仅删除本地远程跟踪引用（refs/remotes/git-svn/xxx）
// 不影响 SVN 服务器，安全操作
func (self *SvnCommands) DeleteLocalRef(refName string) error {
	cmdArgs := NewGitCmd("branch").Arg("-D").Arg("-r").Arg(refName).ToArgv()
	return self.cmd.New(cmdArgs).Run()
}

// Fetch 执行 git svn fetch -all 获取 SVN 更新
func (self *SvnCommands) Fetch() error {
	cmdArgs := NewGitCmd("svn").Arg("fetch").Arg("--all").ToArgv()
	return self.cmd.New(cmdArgs).Run()
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

	// 1. 获取本地 refs（遍历该类型所有 mapping 的 RefsPath，
	// 兼容refs 端影射到非标准路径的配置，如 tags = tags/*:refs/remotes/git-svn/releases/*）
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
				localRefs[strings.TrimPrefix(line, "refs/remotes/git-svn/")] = true
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
					svnBranches[strings.TrimPrefix(m.RefsPath, "refs/remotes/git-svn/")+"/"+name] = true
				}
			}
		}
	}

	// svn list 全部失败（网络不通、认证失败等）时返回错误并中止，
	// 避免把”全部 Stale”的误导性结果当作真实状态展示（错误不缓存，下次自动重时）
	if svnListAttempted && !svnListOk {
		return nil, fmt.Errorf("svn list failed for all %s paths (network or auth error?)", refType)
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


