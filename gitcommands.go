// Go has various value types including strings,
// integers, floats, booleans, etc. Here are a few
// basic examples.

package main

import (
  "fmt"
  // "log"
  "os/exec"
  "os"
  "strings"
  "regexp"
  "runtime"
)

// Map (from https://gobyexample.com/collection-functions)
func Map(vs []string, f func(string) string) []string {
    vsm := make([]string, len(vs))
    for i, v := range vs {
        vsm[i] = f(v)
    }
    return vsm
}

func sanitisedFileString(fileString string) string {
  r := regexp.MustCompile("\\s| \\(new commits\\)|.* ")
  fileString = r.ReplaceAllString(fileString, "")
  return fileString
}

func filesByMatches(statusString string, targets []string) []string {
  files := make([]string, 0)
  for _, target := range targets {
    if strings.Index(statusString, target) == -1 {
      continue
    }
    r := regexp.MustCompile("(?s)" + target + ".*?\n\n(.*?)\n\n")
    // fmt.Println(r)

    matchedFileStrings := strings.Split(r.FindStringSubmatch(statusString)[1], "\n")
    // fmt.Println(matchedFileStrings)

    matchedFiles := Map(matchedFileStrings, sanitisedFileString)
    // fmt.Println(matchedFiles)
    files = append(files, matchedFiles...)

  }

  breakHere()

  // fmt.Println(files)
  return files
}

func breakHere() {
  if len(os.Args) > 1 && os.Args[1] == "debug" {
    runtime.Breakpoint()
  }
}

func getFilesToStage(statusString string) []string {
  targets := []string{"Changes not staged for commit:", "Untracked files:"}
  return filesByMatches(statusString, targets)
}

func getFilesToUnstage(statusString string) []string {
  targets := []string{"Changes to be committed:"}
  return filesByMatches(statusString, targets)
}

func runCommand(cmd string) (string, error) {
  splitCmd := strings.Split(cmd, " ")
  cmdOut, err := exec.Command(splitCmd[0], splitCmd[1:]...).Output()
  return string(cmdOut), err
}

func getDiff(file string, cached bool) string {
  devLog(file)
  cachedArg := ""
  if cached {
    cachedArg = "--cached "
  }
  s, err := runCommand("git diff " + cachedArg + file)
  if err != nil {
    // for now we assume an error means the file was deleted
    return "deleted"
  }
  return s
}

func stageFile(file string) error {
  devLog("staging " + file)
  _, err := runCommand("git add " + file)
  return err
}

func unStageFile(file string) error {
  _, err := runCommand("git reset HEAD " + file)
  return err
}

func testGettingFiles() {

  statusString, _ := runCommand("git status")
  fmt.Println(getFilesToStage(statusString))
  fmt.Println(getFilesToUnstage(statusString))

  runCommand("git add hello-world.go")
}


