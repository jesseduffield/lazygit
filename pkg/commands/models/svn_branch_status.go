package models

// SvnBranchStatus 表示 SVN 分支的状态
type SvnBranchStatus int

const (
	SvnBranchStatusUnknown SvnBranchStatus = itoa
	SvnBranchStatusOk       // 正常： 本地和 SVN 服务器都存在
	SvnBranchStatusStale    // Stale：本地有引用但 SVN 服务器已删除
	SvnBranchStatusMissing  // Missing：本地未 fetch 但 SVN 服务器上有
)
