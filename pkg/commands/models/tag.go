package models

import (
	"strings"
)

// Tag : A git tag
type Tag struct {
	Name string
	// this is either the first line of the message of an annotated tag, or the
	// first line of a commit message for a lightweight tag
	Message string
	// FullRefNameOverride 当非空时，FullRefName() 返回此值而非默认的 refs/tags/<Name>
	// 用于 SVN 类型的 tag，其引用路径在 refs/remotes/git-svn/tags/* 而非 refs/tags/*
	FullRefNameOverride string
	// StaleStatus 表示 SVN tag 的 stale 状态
	StaleStatus SvnBranchStatus
}

func (t *Tag) FullRefName() string {
	return "refs/tags/" + t.RefName()
}

func (t *Tag) RefName() string {
	return t.Name
}

func (t *Tag) ShortRefName() string {
	return t.RefName()
}

func (t *Tag) ParentRefName() string {
	return t.RefName() + "^"
}

func (t *Tag) ID() string {
	return t.RefName()
}

func (t *Tag) URN() string {
	return "tag-" + t.ID()
}

func (t *Tag) Description() string {
	return t.Message
}

// IsSvnTag 判断是否为 SVN tag
func (t *Tag) IsSvnTag() bool {
	return t.FullRefNameOverride != "" && strings.HasPrefix(t.FullRefNameOverride, "refs/remotes/")
}
