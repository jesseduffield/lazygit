package utils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindLazygitRootDirectory(t *testing.T) {
	// This test runs in the pkg/utils directory, so we expect the function to
	// search two levels up for the project root.
	expectedRootDir, err := filepath.Abs(filepath.Join("..", ".."))
	assert.NoError(t, err)

	rootDir, err := FindLazygitRootDirectory()

	assert.NoError(t, err)
	assert.Equal(t, expectedRootDir, rootDir)
}

func TestFindLazygitRootDirectoryOutsideProject(t *testing.T) {
	t.Chdir(t.TempDir())

	_, err := FindLazygitRootDirectory()

	assert.ErrorContains(t, err, "there is no go.mod for a lazygit module")
}

func TestDeclaresLazygitModule(t *testing.T) {
	scenarios := []struct {
		testName string
		contents string
		expected bool
	}{
		{
			testName: "lazygit's go.mod",
			contents: "module github.com/jesseduffield/lazygit\n\ngo 1.25.0\n",
			expected: true,
		},
		{
			testName: "a fork's go.mod",
			contents: "module gitlab.com/somebody-else/lazygit\n\ngo 1.25.0\n",
			expected: true,
		},
		{
			testName: "a fork's go.mod with a major version suffix",
			contents: "module github.com/somebody-else/lazygit/v2\n\ngo 1.25.0\n",
			expected: true,
		},
		{
			testName: "module declaration preceded by a comment",
			contents: "// a comment\n\nmodule github.com/jesseduffield/lazygit\n",
			expected: true,
		},
		{
			testName: "module declaration followed by a comment",
			contents: "module github.com/jesseduffield/lazygit // comment\n\ngo 1.25.0\n",
			expected: true,
		},
		{
			testName: "another project's go.mod",
			contents: "module github.com/jesseduffield/lazydocker\n\ngo 1.25.0\n",
			expected: false,
		},
		{
			testName: "no module declaration",
			contents: "go 1.25.0\n",
			expected: false,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.testName, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "go.mod")
			assert.NoError(t, os.WriteFile(path, []byte(scenario.contents), 0o644))

			assert.Equal(t, scenario.expected, declaresLazygitModule(path))
		})
	}
}

func TestDeclaresLazygitModuleWithoutAGoModFile(t *testing.T) {
	assert.False(t, declaresLazygitModule(filepath.Join(t.TempDir(), "go.mod")))
}
