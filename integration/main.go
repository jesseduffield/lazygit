package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/jesseduffield/lazygit/pkg/commands/oscommands"
	"github.com/jesseduffield/lazygit/pkg/integration"
	"github.com/stretchr/testify/assert"
)

// This file can be invoked directly, but you might find it easier to go through
// test/lazyintegration/main.go, which provides a convenient gui wrapper to integration
// tests.
//
// If invoked directly, you can specify a test by passing it as the first argument.
// You can also specify that you want to record a test by passing RECORD_EVENTS=true
// as an env var.

func main() {
	err := test()
	if err != nil {
		panic(err)
	}
}

func test() error {
	rootDir := integration.GetRootDirectory()
	err := os.Chdir(rootDir)
	if err != nil {
		return err
	}

	testDir := filepath.Join(rootDir, "test", "integration")

	osCommand := oscommands.NewDummyOSCommand()
	err = osCommand.RunCommand("go build -o %s", integration.TempLazygitPath())
	if err != nil {
		return err
	}

	tests, err := integration.LoadTests(testDir)
	if err != nil {
		panic(err)
	}

	record := os.Getenv("RECORD_EVENTS") != ""

	updateSnapshots := record || os.Getenv("UPDATE_SNAPSHOTS") != ""

	selectedTestName := os.Args[1]

	for _, test := range tests {
		if selectedTestName != "" && test.Name != selectedTestName {
			continue
		}

		speeds := integration.GetTestSpeeds(test.Speed, updateSnapshots)
		testPath := filepath.Join(testDir, test.Name)
		actualDir := filepath.Join(testPath, "actual")
		expectedDir := filepath.Join(testPath, "expected")
		configDir := filepath.Join(testPath, "used_config")
		log.Printf("testPath: %s, actualDir: %s, expectedDir: %s", testPath, actualDir, expectedDir)

		for i, speed := range speeds {
			log.Printf("%s: attempting test at speed %f\n", test.Name, speed)

			integration.FindOrCreateDir(testPath)
			integration.PrepareIntegrationTestDir(actualDir)
			err := integration.CreateFixture(testPath, actualDir)
			if err != nil {
				return err
			}

			cmd, err := integration.GetLazygitCommand(testPath, rootDir, record, speed)
			if err != nil {
				return err
			}

			cmd.Stdout = os.Stdout
			cmd.Stdin = os.Stdin
			cmd.Stderr = os.Stderr

			if err := cmd.Run(); err != nil {
				return err
			}

			if updateSnapshots {
				err = oscommands.CopyDir(actualDir, expectedDir)
				if err != nil {
					return err
				}
			}

			actual, expected, err := integration.GenerateSnapshots(actualDir, expectedDir)
			if err != nil {
				return err
			}

			if expected == actual {
				fmt.Printf("%s: success at speed %f\n", test.Name, speed)
				break
			}

			// if the snapshots and we haven't tried all playback speeds different we'll retry at a slower speed
			if i == len(speeds)-1 {
				bytes, err := ioutil.ReadFile(filepath.Join(configDir, "development.log"))
				if err != nil {
					return err
				}
				fmt.Println(string(bytes))
				assert.Equal(MockTestingT{}, expected, actual, fmt.Sprintf("expected:\n%s\nactual:\n%s\n", expected, actual))
				os.Exit(1)
			}
		}
	}

	return nil
}

type MockTestingT struct{}

func (t MockTestingT) Errorf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
}
