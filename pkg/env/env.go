package env

import (
	"os"
)

// This package encapsulates accessing/mutating the ENV of the program.

func GetGitDirEnv() string {
	return os.Getenv("GIT_DIR")
}

func SetGitDirEnv(value string) {
	os.Setenv("GIT_DIR", value)
}

func GetWorkTreeEnv() string {
	return os.Getenv("GIT_WORK_TREE")
}

func SetWorkTreeEnv(value string) {
	os.Setenv("GIT_WORK_TREE", value)
}

func UnsetGitLocationEnvVars() {
	_ = os.Unsetenv("GIT_DIR")
	_ = os.Unsetenv("GIT_WORK_TREE")
}
