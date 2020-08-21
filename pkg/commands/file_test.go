package commands

import (
	"testing"

	"github.com/sirupsen/logrus"
)

// TestMergeGITStatus tests if 2 git status merge successfull
func TestMergeGITStatus(t *testing.T) {
	type scenario struct {
		testName string
		a        string
		b        string
		expected string
	}

	scenarios := []scenario{
		{
			"Nothing",
			"  ",
			"  ",
			"  ",
		}, {
			"Untracked",
			"?",
			" ?",
			"??",
		}, {
			"Untracked and delete",
			"D?",
			"?D",
			"D?",
		}, {
			"Update",
			"RU",
			"UD",
			"UU",
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			res := MergeGITStatus(s.a, s.b)
			if res != s.expected {
				t.Errorf("\"%s\" is not equal to \"%s\"", res, s.expected)
			}
		})
	}
}

// TestHeight tests if dir input returns the correct height
func TestHeight(t *testing.T) {
	type scenario struct {
		testName       string
		dir            *Dir
		expectedHeight int
	}

	scenarios := []scenario{
		{
			"Empty dir",
			&Dir{},
			0, // The dir doesn't have any content
		}, {
			"Dir with subdir",
			&Dir{SubDirs: []*Dir{{}}},
			1, // The dir has 1 sub dir
		}, {
			"Dir with file",
			&Dir{Files: []*File{{}}},
			1, // The dir has 1 file
		}, {
			"Dir with file and subdir",
			&Dir{Files: []*File{{}}, SubDirs: []*Dir{{}}},
			2, // The dir has 1 file and sub dir
		}, {
			"Dir with subdir with subdir",
			&Dir{SubDirs: []*Dir{{SubDirs: []*Dir{{}}}}},
			2, // The dir has 1 sub dir with 1 sub dir
		}, {
			"Dir with subdir with subdir with files",
			&Dir{SubDirs: []*Dir{{SubDirs: []*Dir{{Files: []*File{{}}}}}}},
			3, // The dir has 1 sub dir with 1 sub dir with 1 file
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			res := s.dir.Height()
			if res != s.expectedHeight {
				t.Errorf("%d is not equal to %d", res, s.expectedHeight)
			}
		})
	}
}

// TestMatchPath checks if a path can be matched
func TestMatchPath(t *testing.T) {
	file1 := &File{
		Name: "file 1",
	}
	testDir := NewDir()
	folder1 := testDir.NewDir("folder 1")
	folder2 := testDir.NewDir("folder 2")
	folder3 := testDir.NewDir("folder 3")
	subFolder1 := folder2.NewDir("sub folder 1")
	subFolder1.AddFile(file1)

	type scenario struct {
		testName     string
		path         []int
		expectedDir  *Dir
		expectedFile *File
	}

	scenarios := []scenario{
		{
			"match root",
			[]int{},
			nil,
			nil,
		}, {
			"folder 1",
			[]int{0},
			folder1,
			nil,
		}, {
			"folder 2",
			[]int{1},
			folder2,
			nil,
		}, {
			"folder 3",
			[]int{2},
			folder3,
			nil,
		}, {
			"sub folder 1",
			[]int{1, 0},
			subFolder1,
			nil,
		}, {
			"sub folder file 1",
			[]int{1, 0, 0},
			nil,
			file1,
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			file, dir := testDir.MatchPath(s.path)
			if s.expectedDir != nil {
				if dir == nil {
					t.Errorf("expected dir %s but got nil", s.expectedDir.Name)
				} else if dir != s.expectedDir {
					t.Errorf("expected dir %s but got %s", s.expectedDir.Name, dir.Name)
				}
			} else if dir != nil {
				t.Errorf("expected nil dir but got dir %s", dir.Name)
			}
			if s.expectedFile != nil {
				if file == nil {
					t.Errorf("expected file %s but got nil", s.expectedFile.Name)
				} else if file != s.expectedFile {
					t.Errorf("expected file %s but got %s", s.expectedFile.Name, file.Name)
				}
			} else if file != nil {
				t.Errorf("expected nil file but got file %s", file.Name)
			}
		})
	}
}

// TestCombine tests if dirs can be combined if 1 dir only has 1 subfolder
func TestCombine(t *testing.T) {
	root := NewDir()
	subDir1 := root.NewDir("1")
	subDir2 := subDir1.NewDir("2")
	subDir3 := subDir2.NewDir("3")
	subDir3.AddFile(&File{
		Name: "file 1",
	})
	subDir4 := subDir3.NewDir("4")
	subDir4.AddFile(&File{
		Name: "file 1",
	})

	type scenario struct {
		testName     string
		dir          *Dir
		expectedName string
	}

	scenarios := []scenario{
		{
			"root",
			root,
			"",
		}, {
			"combined dir",
			subDir1,
			"1/2/3",
		}, {
			"Not combined dir",
			subDir4,
			"4",
		},
	}

	root.Combine(logrus.NewEntry(logrus.New()))

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			if s.dir.Name != s.expectedName {
				t.Errorf("expected name \"%s\" but got \"%s\"", s.expectedName, s.dir.Name)
			}
		})
	}
}

// TestFilesToTree tests if a list of files can be changed to a dir
func TestFilesToTree(t *testing.T) {
	files := []*File{
		{Name: "a/file1"},
		{Name: "a/b/file2"},
		{Name: "a/b/file3"},
		{Name: "a/b/c/file4"},
		{Name: "z/file1"},
		{Name: "z/a/b/file2"},
		{Name: "z/a/b/c/file3"},
	}
	root := FilesToTree(logrus.NewEntry(logrus.New()), files)

	type scenario struct {
		testName     string
		dir          *Dir
		expectedName string
	}

	scenarios := []scenario{
		{
			"root",
			root,
			"",
		}, {
			"subdir a",
			root.SubDirs[0],
			"a",
		}, {
			"subdir z",
			root.SubDirs[1],
			"z",
		}, {
			"subdir z/a/b",
			root.SubDirs[1].SubDirs[0],
			"a/b",
		}, {
			"subdir a/b",
			root.SubDirs[0].SubDirs[0],
			"b",
		},
	}

	for _, s := range scenarios {
		t.Run(s.testName, func(t *testing.T) {
			if s.dir.Name != s.expectedName {
				t.Errorf("expected name \"%s\" but got \"%s\"", s.expectedName, s.dir.Name)
			}
		})
	}

}

// TestFileGetY tests if the correct y position can be returned from a file
func TestFileGetY(t *testing.T) {
	f1 := &File{Name: "a/file1"}
	f2 := &File{Name: "a/b/file2"}
	f3 := &File{Name: "a/b/file3"}
	f4 := &File{Name: "a/b/c/file4"}
	f5 := &File{Name: "z/file1"}
	f6 := &File{Name: "z/a/b/file2"}
	f7 := &File{Name: "z/a/b/c/file3"}
	files := []*File{f1, f2, f3, f4, f5, f6, f7}
	FilesToTree(logrus.NewEntry(logrus.New()), files)

	type scenario struct {
		file      *File
		expectedY int
	}

	scenarios := []scenario{
		{f1, 1},
		{f2, 2},
		{f3, 3},
		{f4, 3},
		{f5, 8},
		{f6, 9},
		{f7, 10},
	}

	for _, s := range scenarios {
		t.Run("y of "+s.file.Name, func(t *testing.T) {
			y := s.file.GetY()
			if y != s.expectedY {
				t.Errorf("expected y %d but got %d for dir \"%s\"", s.expectedY, y, s.file.Name)
			}
		})
	}
}
