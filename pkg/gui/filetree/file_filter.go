package filetree

import (
	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/sahilm/fuzzy"
	"github.com/samber/lo"
)

type filePathSource struct {
	files []*models.File
}

func (s *filePathSource) String(i int) string {
	return s.files[i].Path
}

func (s *filePathSource) Len() int {
	return len(s.files)
}

func filterFilesByText(files []*models.File, filter string, useFuzzySearch bool) []*models.File {
	source := &filePathSource{files: files}
	matches := utils.FindFrom(filter, source, useFuzzySearch)
	return lo.Map(matches, func(match fuzzy.Match, _ int) *models.File {
		return files[match.Index]
	})
}

type commitFilePathSource struct {
	files []*models.CommitFile
}

func (s *commitFilePathSource) String(i int) string {
	return s.files[i].Path
}

func (s *commitFilePathSource) Len() int {
	return len(s.files)
}

func filterCommitFilesByText(files []*models.CommitFile, filter string, useFuzzySearch bool) []*models.CommitFile {
	source := &commitFilePathSource{files: files}
	matches := utils.FindFrom(filter, source, useFuzzySearch)
	return lo.Map(matches, func(match fuzzy.Match, _ int) *models.CommitFile {
		return files[match.Index]
	})
}
