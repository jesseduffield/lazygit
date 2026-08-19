package env

import (
	"os"
	"strings"
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

// GetGitLocationEnvVars returns the location variables that are set, as
// "NAME=value" entries.
func GetGitLocationEnvVars() []string {
	envVars := []string{}
	for _, name := range []string{GitDirEnvVar, GitWorkTreeEnvVar} {
		if value := os.Getenv(name); value != "" {
			envVars = append(envVars, name+"="+value)
		}
	}
	return envVars
}

// SetGitLocationEnvVars sets the location variables from "NAME=value" entries,
// clearing both first so that only what is given remains. Passing nothing is
// how you say the repo is to be found from the working directory.
func SetGitLocationEnvVars(envVars []string) {
	UnsetGitLocationEnvVars()
	for _, envVar := range envVars {
		if name, value, ok := strings.Cut(envVar, "="); ok {
			os.Setenv(name, value)
		}
	}
}
