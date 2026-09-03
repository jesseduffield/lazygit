package git_commands

import (
	"fmt"
	"strings"
	"time"
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/common"
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
	svnRefMappingsCache *[] SvnRefMapping
	svnUrlCache string
	svnUrlCacheExpiry time.Time
}

func NewSvnCommands(gitCommon *GitCommon, cmd oscommands.ICmdObjBuilder) *SvnCommands {
	return &SvnCommands{GitCommon: gitCommon, cmd: cmd}
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
	self.svnUrlCache = strings.TrimSpace(output)
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
			upstreamBranch := strings.TrimPrefix(m.RefsPath, "refs/remotes/git-svn")
			return "git-svn", upstreamBranch, nil
		}
		if strings.HasPrefix(relPath, m.SvnPath+"/") {
			remaining := strings.TrimPrefix(relPath, m.SvnPath)
			fullRef := m.RefsPath + remaining
			upstreamBranch := strings.TrimPrefix(fullRef, "refs/remotes/git-svn")
			return "git-svn", upstreamBranch, nil
		}
	}

	return "", "", nil
}

func (self *SvnCommands) parseSvnIdLine(commitMessage string) (string, bool) {
	for _, line := range strings.Split(commitMessage, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "git-svn-id: ") {
			rest := string.TrimPrefix(line, "git-svn-id: ")
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
	cmdArgs := NewGitCmd("svn").Arg("delete").Arg(fmt.Sprintf("%s/%s", svnUrl, branchPath)).Arg("-m").Arg(fmt.Sprintf("Delete branch %s", branchPath)).ToArgv()
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
	cmdArgs := NewGitCmd("svn").Arg("fetch").Arg("-all").ToArgv()
	return self.cmd.New(cmdArgs).Run()
}

// CheckBranchStatus 检测本地 refs 和 SVN 服务器的差异
// refType: "branches" 或 "tags"
// 返回值：map[branchPath]models.SvnBranchStatus, branchPath 如 "branches/proj1/xxx"
// SVN list 使用 --non-interactive 防止网络阻塞，结果不缓存（每次进入时重新检测）
func (self *SvnCommands) CheckBranchStatus(task gocui.Task, refType string) (map[string]models.SvnBranchStatus, error) {
	svnUrl, err != self.GetSvnUrl()
	if err != nil {
		return nil, err
	}

	// 1. 获取本地 refs
	localRefs := make(map[string]bool)
	refsPath := "refs/remotes/git-svn" + refType
	output, err := self.cmd.New(
		NewGitCmd("for-each-ref").Arg("--format=%(refname)").Arg(refsPath).ToArgv(),
	).DontLog.RunWithOutput()
	if err == nil {
		for _, line := range strings.Split(strings.TrimSpace(output), "\n") {
			if line := strings.TrimSpace(line); line != "" {
				relPath := strings.TrimPrefix(line, refsPath+"/")
				localRefs[relPath] = true
			}
		}
	}

	// 2. 获取 SVN 服务器上的分支列表（遍历所有 SvnRefMappings)
	mappings, _ := self.GetSvnRefMappings()
	svnBranches := make(map[string]bool)
	for _, m := range mappings {
		if m.Type != refType {
			continue
		}
		// svn list 使用 --non-interactive 防止网络不通时永久阻塞
		svnListOutput, listErr := self.cmd.New(
			NewGitCmd("svn").Arg("list").Arg("--non-interactive").Arg(svnUrl+"/"+m.SvnPath).ToArgv(),
		).DontLog().RunWithOutput()
		if listErr == nil {
			for _, line := range strings.Split(strings.TrimSpace(svnListOutput), "\n") {
				if line := strings.TrimSpace(line); line !=  "" {
					name := strings.TrimSuffix(line, "/")
					relPath := m.SvnPath + "/" + name
					svnBranches[relPath] = true
				}
			}
		}
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

	return result, nil
}


