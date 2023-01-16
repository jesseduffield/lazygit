package git_commands

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
)

type GitVersion struct {
	Major, Minor, Patch int
	Additional          string
}

func GetGitVersion(osCommand *oscommands.OSCommand) (*GitVersion, error) {
	versionStr, _, err := osCommand.Cmd.New("git --version").RunWithOutputs()
	if err != nil {
		return nil, err
	}

	version, err := ParseGitVersion(versionStr)
	if err != nil {
		return nil, err
	}

	return version, nil
}

func ParseGitVersion(versionStr string) (*GitVersion, error) {
	// versionStr should be something like:
	// git version 2.39.0
	// git version 2.37.1 (Apple Git-137.1)
	re := regexp.MustCompile(`[^\d]+(\d+)(\.\d+)?(\.\d+)?(.*)`)
	matches := re.FindStringSubmatch(versionStr)

	if len(matches) < 5 {
		return nil, errors.New("unexpected git version format: " + versionStr)
	}

	v := &GitVersion{}
	var err error

	if v.Major, err = strconv.Atoi(matches[1]); err != nil {
		return nil, err
	}
	if len(matches[2]) > 1 {
		if v.Minor, err = strconv.Atoi(matches[2][1:]); err != nil {
			return nil, err
		}
	}
	if len(matches[3]) > 1 {
		if v.Patch, err = strconv.Atoi(matches[3][1:]); err != nil {
			return nil, err
		}
	}
	v.Additional = strings.Trim(matches[4], " \r\n")

	return v, nil
}

func (v *GitVersion) IsOlderThan(major, minor, patch int) bool {
	actual := v.Major*1000*1000 + v.Minor*1000 + v.Patch
	required := major*1000*1000 + minor*1000 + patch
	return actual < required
}
