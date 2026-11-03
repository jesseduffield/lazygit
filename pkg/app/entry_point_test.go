package app

import (
	"runtime/debug"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMergeBuildInfo(t *testing.T) {
	scenarios := []struct {
		name        string
		buildInfo   BuildInfo
		goBuildInfo *debug.BuildInfo
		expected    BuildInfo
	}{
		{
			name:      "version set by ldflags",
			buildInfo: BuildInfo{Version: "0.66.2", BuildSource: "binaryRelease"},
			goBuildInfo: &debug.BuildInfo{
				Main:     debug.Module{Version: "(devel)"},
				Settings: []debug.BuildSetting{{Key: "vcs.revision", Value: "0123456789abcdef"}},
			},
			expected: BuildInfo{Version: "0.66.2", BuildSource: "binaryRelease"},
		},
		{
			name:      "built from a git checkout",
			buildInfo: BuildInfo{BuildSource: "unknown"},
			goBuildInfo: &debug.BuildInfo{
				Main: debug.Module{Version: "v0.66.3-0.20261103120000-0123456789ab+dirty"},
				Settings: []debug.BuildSetting{
					{Key: "vcs.revision", Value: "0123456789abcdef"},
					{Key: "vcs.time", Value: "2026-11-03T12:00:00Z"},
				},
			},
			expected: BuildInfo{
				Commit:      "0123456789abcdef",
				Date:        "2026-11-03T12:00:00Z",
				Version:     "01234567",
				BuildSource: "unknown",
			},
		},
		{
			name:        "built with go install",
			buildInfo:   BuildInfo{BuildSource: "unknown"},
			goBuildInfo: &debug.BuildInfo{Main: debug.Module{Version: "v0.66.2"}},
			expected:    BuildInfo{Version: "0.66.2", BuildSource: "goInstall"},
		},
		{
			name:      "built with go install from a commit",
			buildInfo: BuildInfo{BuildSource: "unknown"},
			goBuildInfo: &debug.BuildInfo{
				Main: debug.Module{Version: "v0.66.3-0.20261103120000-0123456789ab"},
			},
			expected: BuildInfo{Version: "unversioned", BuildSource: "unknown"},
		},
		{
			name:        "no build info",
			buildInfo:   BuildInfo{BuildSource: "unknown"},
			goBuildInfo: nil,
			expected:    BuildInfo{Version: "unversioned", BuildSource: "unknown"},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			buildInfo := s.buildInfo
			mergeBuildInfo(&buildInfo, s.goBuildInfo)
			assert.Equal(t, s.expected, buildInfo)
		})
	}
}
