package models

import "path/filepath"

type SubmoduleConfig struct {
	Name string
	Path string
	Url  string

	ParentModule *SubmoduleConfig // nil if top-level
}

func (r *SubmoduleConfig) FullName() string {
	if r.ParentModule != nil {
		return r.ParentModule.FullName() + "/" + r.Name
	}

	return r.Name
}

func (r *SubmoduleConfig) FullPath() string {
	if r.ParentModule != nil {
		return r.ParentModule.FullPath() + "/" + r.Path
	}

	return r.Path
}

func (r *SubmoduleConfig) RefName() string {
	return r.FullName()
}

func (r *SubmoduleConfig) ID() string {
	return r.RefName()
}

func (r *SubmoduleConfig) Description() string {
	return r.RefName()
}

func (r *SubmoduleConfig) GitDirPath(repoGitDirPath string) string {
	parentPath := repoGitDirPath
	if r.ParentModule != nil {
		parentPath = r.ParentModule.GitDirPath(repoGitDirPath)
	}

	return filepath.Join(parentPath, "modules", r.Name)
}
