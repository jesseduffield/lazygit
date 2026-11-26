package git_commands

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/common"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestGetPendingPatches(t *testing.T) {
	dir := filepath.Join("/repo", ".git", "rebase-apply")
	files := map[string]string{
		filepath.Join(dir, "next"): "2\n",
		filepath.Join(dir, "last"): "3\n",
		filepath.Join(dir, "0001"): "From 1111111 Mon Sep 17 00:00:00 2001\nSubject: [PATCH 1/3] one\n\n",
		filepath.Join(dir, "0002"): "From 2222222 Mon Sep 17 00:00:00 2001\nSubject: [PATCH 2/3] two\n\n",
		filepath.Join(dir, "0003"): "From 3333333 Mon Sep 17 00:00:00 2001\nSubject: [PATCH 3/3] three\n\n",
	}
	loader := &CommitLoader{
		Common:    common.NewDummyCommon(),
		GitCommon: &GitCommon{repoPaths: MockRepoPaths("/repo")},
		readFile: func(filename string) ([]byte, error) {
			if content, ok := files[filename]; ok {
				return []byte(content), nil
			}
			return nil, errors.New("no such file")
		},
	}

	commits := loader.getPendingPatches(&utils.StringPool{})

	assert.Equal(t, []string{"three", "two"}, lo.Map(commits, func(c *models.Commit, _ int) string { return c.Name }))
	assert.Equal(t, []models.CommitStatus{models.StatusApplyingPatch, models.StatusConflicted},
		[]models.CommitStatus{commits[0].Status, commits[1].Status})
	assert.Equal(t, filepath.Join(dir, "0003"), commits[0].PatchFile)
	assert.Equal(t, filepath.Join(dir, "0002"), commits[1].PatchFile)
}

func TestPatchToCommit(t *testing.T) {
	scenarios := []struct {
		testName              string
		content               string
		expectedName          string
		expectedAuthorName    string
		expectedAuthorEmail   string
		expectedUnixTimestamp int64
	}{
		{
			testName: "patch made by git format-patch",
			content: "From 1234567 Mon Sep 17 00:00:00 2001\n" +
				"From: Jane Doe <jane@example.com>\n" +
				"Date: Sat, 7 Nov 2026 11:25:00 +0100\n" +
				"Subject: [PATCH 2/2] Add a file\n" +
				"\n" +
				"---\n" +
				" file | 1 +\n",
			expectedName:          "Add a file",
			expectedAuthorName:    "Jane Doe",
			expectedAuthorEmail:   "jane@example.com",
			expectedUnixTimestamp: 1794047100,
		},
		{
			testName: "subject and author name with non-ASCII characters",
			content: "From 1234567 Mon Sep 17 00:00:00 2001\n" +
				"From: =?UTF-8?q?J=C3=BCrgen?= <juergen@example.com>\n" +
				"Subject: [PATCH] =?UTF-8?q?Fix=20the=20=C3=A4?=\n" +
				"\n",
			expectedName:        "Fix the ä",
			expectedAuthorName:  "Jürgen",
			expectedAuthorEmail: "juergen@example.com",
		},
		{
			testName: "long subject that spans several lines",
			content: "From 1234567 Mon Sep 17 00:00:00 2001\n" +
				"Subject: [PATCH v2 1/3] A subject that is long enough to be\n" +
				" folded\n" +
				"\n",
			expectedName: "A subject that is long enough to be folded",
		},
		{
			testName:     "reply",
			content:      "Subject: Re: [PATCH] Fix it\n\n",
			expectedName: "Fix it",
		},
		{
			testName:     "not a mail message",
			content:      "garbage",
			expectedName: "",
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			commit := patchToCommit(&utils.StringPool{}, []byte(s.content), "/repo/.git/rebase-apply/0001", models.StatusApplyingPatch)

			assert.Equal(t, s.expectedName, commit.Name)
			assert.Equal(t, s.expectedAuthorName, commit.AuthorName)
			assert.Equal(t, s.expectedAuthorEmail, commit.AuthorEmail)
			assert.Equal(t, s.expectedUnixTimestamp, commit.UnixTimestamp)
			assert.Equal(t, "", commit.Hash())
			assert.Equal(t, "/repo/.git/rebase-apply/0001", commit.PatchFile)
		})
	}
}
