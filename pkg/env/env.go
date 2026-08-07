package env

import (
	"os"
)

// This package encapsulates accessing/mutating the ENV of the program.

// The variables with which git can be told where a repo is, rather than having
// it find out from the working directory.
const (
	GitDirEnvVar      = "GIT_DIR"
	GitWorkTreeEnvVar = "GIT_WORK_TREE"
)

func GetGitDirEnv() string {
	return os.Getenv(GitDirEnvVar)
}

func SetGitDirEnv(value string) {
	os.Setenv(GitDirEnvVar, value)
}

func GetWorkTreeEnv() string {
	return os.Getenv(GitWorkTreeEnvVar)
}

func SetWorkTreeEnv(value string) {
	os.Setenv(GitWorkTreeEnvVar, value)
}

func UnsetGitLocationEnvVars() {
	_ = os.Unsetenv(GitDirEnvVar)
	_ = os.Unsetenv(GitWorkTreeEnvVar)
}
