package git_config

import (
	"os"
	"testing"
)

func TestGitConfigCount(t *testing.T) {
	type scenario struct {
		TestName             string
		Envs                 map[string]string
		ExpectedConfigValues map[string]string
	}
	scenarios := []scenario{
		{
			TestName: "1 config",
			Envs: map[string]string{
				"GIT_CONFIG_COUNT":   "1",
				"GIT_CONFIG_KEY_0":   "user.signingkey",
				"GIT_CONFIG_VALUE_0": "abc0123",
			},
			ExpectedConfigValues: map[string]string{
				"user.signingkey": "abc0123",
			},
		},
		{
			TestName: "2 configs",
			Envs: map[string]string{
				"GIT_CONFIG_COUNT":   "2",
				"GIT_CONFIG_KEY_0":   "custom.foo",
				"GIT_CONFIG_VALUE_0": "bob",
				"GIT_CONFIG_KEY_1":   "custom.bar",
				"GIT_CONFIG_VALUE_1": "alice",
			},
			ExpectedConfigValues: map[string]string{
				"custom.foo": "bob",
				"custom.bar": "alice",
			},
		},
	}
	for _, s := range scenarios {
		t.Run(s.TestName, func(t *testing.T) {
			for k, v := range s.Envs {
				os.Setenv(k, v)
			}
			for k, v := range s.ExpectedConfigValues {
				cmd := getGitConfigCmd(k)
				res, err := runGitConfigCmd(cmd)
				if err != nil {
					t.Error(err)
				}
				if v != res {
					t.Errorf("expected: %s, got: %s", v, res)
				}
			}
		})
	}
}
