package env

import (
	"os"
)

// This package encapsulates accessing/mutating the ENV of the program.

func GetGitDirEnv() string {
	return os.Getenv("GIT_DIR")
}

func GetGitWorkTreeEnv() string {
	return os.Getenv("GIT_WORK_TREE")
}

func SetGitDirEnv(value string) {
	os.Setenv("GIT_DIR", value)
}

func SetGitWorkTreeEnv(value string) {
	os.Setenv("GIT_WORK_TREE", value)
}

func UnsetGitDirEnvs() {
	_ = os.Unsetenv("GIT_DIR")
	_ = os.Unsetenv("GIT_WORK_TREE")
}
