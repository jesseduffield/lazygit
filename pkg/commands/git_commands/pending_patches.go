package git_commands

import (
	"bytes"
	"fmt"
	"mime"
	"net/mail"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/models"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/stefanhaller/git-todo-parser/todo"
)

// getPendingPatches returns the patches that "git am" hasn't applied yet; the
// last of them is the one that it stopped at. They usually come from another
// repository, so they have no hashes that we could show or look up.
func (self *CommitLoader) getPendingPatches(hashPool *utils.StringPool) []*models.Commit {
	dir := filepath.Join(self.repoPaths.WorktreeGitDirPath(), "rebase-apply")

	readNumber := func(name string) (int, error) {
		content, err := self.readFile(filepath.Join(dir, name))
		if err != nil {
			return 0, err
		}
		return strconv.Atoi(strings.TrimSpace(string(content)))
	}

	next, err := readNumber("next")
	if err != nil {
		self.Log.Errorf("error occurred reading rebase-apply/next: %s", err.Error())
		return nil
	}
	last, err := readNumber("last")
	if err != nil {
		self.Log.Errorf("error occurred reading rebase-apply/last: %s", err.Error())
		return nil
	}

	commits := []*models.Commit{}
	for number := next; number <= last; number++ {
		// git names the patch files with (at least) four digits
		patchFile := filepath.Join(dir, fmt.Sprintf("%04d", number))
		content, err := self.readFile(patchFile)
		if err != nil {
			self.Log.Errorf("error occurred reading %s: %s", patchFile, err.Error())
			continue
		}

		status := models.StatusApplyingPatch
		if number == next {
			status = models.StatusConflicted
		}
		commits = utils.Prepend(commits, patchToCommit(hashPool, content, patchFile, status))
	}

	return commits
}

// patchToCommit makes a commit for the list of pending patches from a patch
// file that "git am" has split off from the mailbox it is applying.
func patchToCommit(hashPool *utils.StringPool, content []byte, patchFile string, status models.CommitStatus) *models.Commit {
	opts := models.NewCommitOpts{
		Status:    status,
		Action:    todo.Pick,
		PatchFile: patchFile,
	}

	// The file starts with a line like "From <hash> Mon Sep 17 00:00:00 2001",
	// which isn't a header
	if bytes.HasPrefix(content, []byte("From ")) {
		_, content, _ = bytes.Cut(content, []byte("\n"))
	}

	message, err := mail.ReadMessage(bytes.NewReader(content))
	if err != nil {
		return models.NewCommit(hashPool, opts)
	}

	subject := message.Header.Get("Subject")
	if decodedSubject, err := new(mime.WordDecoder).DecodeHeader(subject); err == nil {
		subject = decodedSubject
	}
	opts.Name = cleanUpPatchSubject(subject)

	if author, err := mail.ParseAddress(message.Header.Get("From")); err == nil {
		opts.AuthorName = author.Name
		opts.AuthorEmail = author.Address
	}

	if date, err := message.Header.Date(); err == nil {
		opts.UnixTimestamp = date.Unix()
	}

	return models.NewCommit(hashPool, opts)
}

// cleanUpPatchSubject removes the parts of a patch's subject that "git am"
// doesn't use for the commit message, like "[PATCH 1/2]" or "Re:".
func cleanUpPatchSubject(subject string) string {
	for {
		subject = strings.TrimSpace(subject)
		switch {
		case strings.HasPrefix(subject, "["):
			end := strings.Index(subject, "]")
			if end < 0 {
				return subject
			}
			subject = subject[end+1:]
		case len(subject) >= 3 && strings.EqualFold(subject[:3], "re:"):
			subject = subject[3:]
		default:
			return subject
		}
	}
}
