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

func UnsetGitDirEnv() {
	_ = os.Unsetenv("GIT_DIR")
}
