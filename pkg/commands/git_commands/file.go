package git_commands

import (
	"os"
	"strconv"
	"strings"

	"github.com/go-errors/errors"
	"github.com/jesseduffield/lazygit/pkg/config"
	"github.com/jesseduffield/lazygit/pkg/utils"
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

func (self *FileCommands) GetEditCmdStrLegacy(filename string, lineNumber int) (string, error) {
	editor := self.UserConfig.OS.EditCommand

	if editor == "" {
		editor = self.config.GetCoreEditor()
	}
	if editor == "" {
		editor = self.os.Getenv("GIT_EDITOR")
	}
	if editor == "" {
		editor = self.os.Getenv("VISUAL")
	}
	if editor == "" {
		editor = self.os.Getenv("EDITOR")
	}
	if editor == "" {
		if err := self.cmd.New([]string{"which", "vi"}).DontLog().Run(); err == nil {
			editor = "vi"
		}
	}
	if editor == "" {
		return "", errors.New("No editor defined in config file, $GIT_EDITOR, $VISUAL, $EDITOR, or git config")
	}

	templateValues := map[string]string{
		"editor":   editor,
		"filename": self.cmd.Quote(filename),
		"line":     strconv.Itoa(lineNumber),
	}

	editCmdTemplate := self.UserConfig.OS.EditCommandTemplate
	if len(editCmdTemplate) == 0 {
		switch editor {
		case "emacs", "nano", "vi", "vim", "nvim":
			editCmdTemplate = "{{editor}} +{{line}} -- {{filename}}"
		case "subl":
			editCmdTemplate = "{{editor}} -- {{filename}}:{{line}}"
		case "code":
			editCmdTemplate = "{{editor}} -r --goto -- {{filename}}:{{line}}"
		default:
			editCmdTemplate = "{{editor}} -- {{filename}}"
		}
	}
	return utils.ResolvePlaceholderString(editCmdTemplate, templateValues), nil
}

func (self *FileCommands) GetEditCmdStr(filename string) (string, bool) {
	// Legacy support for old config; to be removed at some point
	if self.UserConfig.OS.Edit == "" && self.UserConfig.OS.EditCommandTemplate != "" {
		if cmdStr, err := self.GetEditCmdStrLegacy(filename, 1); err == nil {
			return cmdStr, true
		}
	}

	template, editInTerminal := config.GetEditTemplate(&self.UserConfig.OS, self.guessDefaultEditor)

	templateValues := map[string]string{
		"filename": self.cmd.Quote(filename),
	}

	cmdStr := utils.ResolvePlaceholderString(template, templateValues)
	return cmdStr, editInTerminal
}

func (self *FileCommands) GetEditAtLineCmdStr(filename string, lineNumber int) (string, bool) {
	// Legacy support for old config; to be removed at some point
	if self.UserConfig.OS.EditAtLine == "" && self.UserConfig.OS.EditCommandTemplate != "" {
		if cmdStr, err := self.GetEditCmdStrLegacy(filename, lineNumber); err == nil {
			return cmdStr, true
		}
	}

	template, editInTerminal := config.GetEditAtLineTemplate(&self.UserConfig.OS, self.guessDefaultEditor)

	templateValues := map[string]string{
		"filename": self.cmd.Quote(filename),
		"line":     strconv.Itoa(lineNumber),
	}

	cmdStr := utils.ResolvePlaceholderString(template, templateValues)
	return cmdStr, editInTerminal
}

func (self *FileCommands) GetEditAtLineAndWaitCmdStr(filename string, lineNumber int) string {
	// Legacy support for old config; to be removed at some point
	if self.UserConfig.OS.EditAtLineAndWait == "" && self.UserConfig.OS.EditCommandTemplate != "" {
		if cmdStr, err := self.GetEditCmdStrLegacy(filename, lineNumber); err == nil {
			return cmdStr
		}
	}

	template := config.GetEditAtLineAndWaitTemplate(&self.UserConfig.OS, self.guessDefaultEditor)

	templateValues := map[string]string{
		"filename": self.cmd.Quote(filename),
		"line":     strconv.Itoa(lineNumber),
	}

	cmdStr := utils.ResolvePlaceholderString(template, templateValues)
	return cmdStr
}

func (self *FileCommands) GetOpenDirInEditorCmdStr(path string) string {
	template := config.GetOpenDirInEditorTemplate(&self.UserConfig.OS, self.guessDefaultEditor)

	templateValues := map[string]string{
		"dir": self.cmd.Quote(path),
	}

	cmdStr := utils.ResolvePlaceholderString(template, templateValues)
	return cmdStr
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
