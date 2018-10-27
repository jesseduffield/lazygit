package procs

import (
	"fmt"
	"os"
	"strings"
)

// ParseEnv takes an environment []string and converts it to a map[string]string.
func ParseEnv(environ []string) map[string]string {
	env := make(map[string]string)
	for _, e := range environ {
		pair := strings.SplitN(e, "=", 2)

		// There is a chance we can get an env with empty values
		if len(pair) == 2 {
			env[pair[0]] = pair[1]
		}
	}
	return env
}

// Env takes a map[string]string and converts it to a []string that
// can be used with exec.Cmd. The useEnv boolean flag will include the
// current process environment, overlaying the provided env
// map[string]string.
func Env(env map[string]string, useEnv bool) []string {
	envlist := []string{}

	// update our env by loading our env and overriding any values in
	// the provided env.
	if useEnv {
		environ := ParseEnv(os.Environ())
		for k, v := range env {
			environ[k] = v
		}
		env = environ
	}

	for key, val := range env {
		if key == "" {
			continue
		}
		envlist = append(envlist, fmt.Sprintf("%s=%s", key, val))
	}

	return envlist
}
