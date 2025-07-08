package git_commands

import (
	"os"
	"strconv"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/utils"
	"github.com/samber/lo"
)

type FileCommands struct {
	*GitCommon
}

func NewFileCommands(gitCommon *GitCommon) *FileCommands {
	return &FileCommands{
		GitCommon: gitCommon,
	}
}

// Cat obtains the content of a file
func (self *FileCommands) Cat(fileName string) (string, error) {
	buf, err := os.ReadFile(fileName)
	if err != nil {
		return "", nil
	}
	return string(buf), nil
}

func (self *FileCommands) GetEditCmdStr(filenames []string) (string, bool) {
	template, suspend := config.GetEditTemplate(self.os.Platform.Shell, &self.UserConfig().OS, self.guessDefaultEditor)
	quotedFilenames := lo.Map(filenames, func(filename string, _ int) string { return self.cmd.Quote(filename) })

	templateValues := map[string]string{
		"filename": strings.Join(quotedFilenames, " "),
	}

	cmdStr := utils.ResolvePlaceholderString(template, templateValues)
	return cmdStr, suspend
}

func (self *FileCommands) GetEditAtLineCmdStr(filename string, lineNumber int) (string, bool) {
	template, suspend := config.GetEditAtLineTemplate(self.os.Platform.Shell, &self.UserConfig().OS, self.guessDefaultEditor)

	templateValues := map[string]string{
		"filename": self.cmd.Quote(filename),
		"line":     strconv.Itoa(lineNumber),
	}

	cmdStr := utils.ResolvePlaceholderString(template, templateValues)
	return cmdStr, suspend
}

func (self *FileCommands) GetEditAtLineAndWaitCmdStr(filename string, lineNumber int) string {
	template := config.GetEditAtLineAndWaitTemplate(self.os.Platform.Shell, &self.UserConfig().OS, self.guessDefaultEditor)

	templateValues := map[string]string{
		"filename": self.cmd.Quote(filename),
		"line":     strconv.Itoa(lineNumber),
	}

	cmdStr := utils.ResolvePlaceholderString(template, templateValues)
	return cmdStr
}

func (self *FileCommands) GetOpenDirInEditorCmdStr(path string) (string, bool) {
	template, suspend := config.GetOpenDirInEditorTemplate(self.os.Platform.Shell, &self.UserConfig().OS, self.guessDefaultEditor)

	templateValues := map[string]string{
		"dir": self.cmd.Quote(path),
	}

	cmdStr := utils.ResolvePlaceholderString(template, templateValues)
	return cmdStr, suspend
}

func (self *FileCommands) guessDefaultEditor() string {
	// Try to query a few places where editors get configured
	editor := self.config.GetCoreEditor()
	if editor == "" {
		editor = self.os.Getenv("GIT_EDITOR")
	}
	if editor == "" {
		editor = self.os.Getenv("VISUAL")
	}
	if editor == "" {
		editor = self.os.Getenv("EDITOR")
	}

	if editor != "" {
		// At this point, it might be more than just the name of the editor;
		// e.g. it might be "code -w" or "vim -u myvim.rc". So assume that
		// everything up to the first space is the editor name.
		editor = strings.Split(editor, " ")[0]
	}

	return editor
}
