package git_commands

import (
	"os"
	"path/filepath"
	"strings"
)

type BisectCommands struct {
	*GitCommon
}

func NewBisectCommands(gitCommon *GitCommon) *BisectCommands {
	return &BisectCommands{
		GitCommon: gitCommon,
	}
}

// This command is pretty cheap to run so we're not storing the result anywhere.
// But if it becomes problematic we can chang that.
func (self *BisectCommands) GetInfo() *BisectInfo {
	return self.GetInfoForGitDir(self.repoPaths.WorktreeGitDirPath())
}

func (self *BisectCommands) GetInfoForGitDir(gitDir string) *BisectInfo {
	var err error
	info := &BisectInfo{started: false, log: self.Log, newTerm: "bad", oldTerm: "good"}
	// we return nil if we're not in a git bisect session.
	// we know we're in a session by the presence of a .git/BISECT_START file

	bisectStartPath := filepath.Join(gitDir, "BISECT_START")
	exists, err := self.os.FileExists(bisectStartPath)
	if err != nil {
		self.Log.Infof("error getting git bisect info: %s", err.Error())
		return info
	}

	if !exists {
		return info
	}

	startContent, err := os.ReadFile(bisectStartPath)
	if err != nil {
		self.Log.Infof("error getting git bisect info: %s", err.Error())
		return info
	}

	info.started = true
	info.start = strings.TrimSpace(string(startContent))

	termsContent, err := os.ReadFile(filepath.Join(gitDir, "BISECT_TERMS"))
	if err != nil {
		// old git versions won't have this file so we default to bad/good
	} else {
		splitContent := strings.Split(string(termsContent), "\n")
		info.newTerm = splitContent[0]
		info.oldTerm = splitContent[1]
	}

	bisectRefsDir := filepath.Join(gitDir, "refs", "bisect")
	files, err := os.ReadDir(bisectRefsDir)
	if err != nil {
		self.Log.Infof("error getting git bisect info: %s", err.Error())
		return info
	}

	info.statusMap = make(map[string]BisectStatus)
	for _, file := range files {
		status := BisectStatusSkipped
		name := file.Name()
		path := filepath.Join(bisectRefsDir, name)

		fileContent, err := os.ReadFile(path)
		if err != nil {
			self.Log.Infof("error getting git bisect info: %s", err.Error())
			return info
		}

		sha := strings.TrimSpace(string(fileContent))

		if name == info.newTerm {
			status = BisectStatusNew
		} else if strings.HasPrefix(name, info.oldTerm+"-") {
			status = BisectStatusOld
		} else if strings.HasPrefix(name, "skipped-") {
			status = BisectStatusSkipped
		}

		info.statusMap[sha] = status
	}

	currentContent, err := os.ReadFile(filepath.Join(gitDir, "BISECT_EXPECTED_REV"))
	if err != nil {
		self.Log.Infof("error getting git bisect info: %s", err.Error())
		return info
	}
	currentSha := strings.TrimSpace(string(currentContent))
	info.current = currentSha

	return info
}

func (self *BisectCommands) Reset() error {
	cmdArgs := NewGitCmd("bisect").Arg("reset").ToArgv()

	return self.cmd.New(cmdArgs).StreamOutput().Run()
}

func (self *BisectCommands) Mark(ref string, term string) error {
	cmdArgs := NewGitCmd("bisect").Arg(term, ref).ToArgv()

	return self.cmd.New(cmdArgs).
		IgnoreEmptyError().
		StreamOutput().
		Run()
}

func (self *BisectCommands) Skip(ref string) error {
	return self.Mark(ref, "skip")
}

func (self *BisectCommands) Start() error {
	cmdArgs := NewGitCmd("bisect").Arg("start").ToArgv()

	return self.cmd.New(cmdArgs).StreamOutput().Run()
}

func (self *BisectCommands) StartWithTerms(oldTerm string, newTerm string) error {
	cmdArgs := NewGitCmd("bisect").Arg("start").
		Arg("--term-old=" + oldTerm).
		Arg("--term-new=" + newTerm).
		ToArgv()

	return self.cmd.New(cmdArgs).StreamOutput().Run()
}

// tells us whether we've found our problem commit(s). We return a string slice of
// commit sha's if we're done, and that slice may have more that one item if
// skipped commits are involved.
func (self *BisectCommands) IsDone() (bool, []string, error) {
	info := self.GetInfo()
	if !info.Bisecting() {
		return false, nil, nil
	}

	newSha := info.GetNewSha()
	if newSha == "" {
		return false, nil, nil
	}

	// if we start from the new commit and reach the a good commit without
	// coming across any unprocessed commits, then we're done
	done := false
	candidates := []string{}

	cmdArgs := NewGitCmd("rev-list").Arg(newSha).ToArgv()
	err := self.cmd.New(cmdArgs).RunAndProcessLines(func(line string) (bool, error) {
		sha := strings.TrimSpace(line)

		if status, ok := info.statusMap[sha]; ok {
			switch status {
			case BisectStatusSkipped, BisectStatusNew:
				candidates = append(candidates, sha)
				return false, nil
			case BisectStatusOld:
				done = true
				return true, nil
			}
		} else {
			return true, nil
		}

		// should never land here
		return true, nil
	})
	if err != nil {
		return false, nil, err
	}

	return done, candidates, nil
}

// tells us whether the 'start' ref that we'll be sent back to after we're done
// bisecting is actually a descendant of our current bisect commit. If it's not, we need to
// render the commits from the bad commit.
func (self *BisectCommands) ReachableFromStart(bisectInfo *BisectInfo) bool {
	cmdArgs := NewGitCmd("merge-base").
		Arg("--is-ancestor", bisectInfo.GetNewSha(), bisectInfo.GetStartSha()).
		ToArgv()

	err := self.cmd.New(cmdArgs).DontLog().Run()

	return err == nil
}
