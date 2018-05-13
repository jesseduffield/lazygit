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
)

func plus(a int, b int) int {
  return a + b
}

// from https://gobyexample.com/collection-functions
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
    fmt.Println(r)

    matchedFileStrings := strings.Split(r.FindStringSubmatch(statusString)[1], "\n")
    // fmt.Println(matchedFileStrings)

    matchedFiles := Map(matchedFileStrings, sanitisedFileString)
    // fmt.Println(matchedFiles)
    files = append(files, matchedFiles...)

  }
  fmt.Println(files)
  return files
}

func filesToStage(statusString string) []string {
  targets := []string{"Changes not staged for commit:", "Untracked files:"}
  return filesByMatches(statusString, targets)
}

func runCommand(cmd string) (string, error) {
  splitCmd := strings.Split(cmd, " ")
  fmt.Println(splitCmd)

  cmdOut, err := exec.Command(splitCmd[0], splitCmd[1:]...).Output()

  // if err != nil {
  //   os.Exit(1)
  // }

  return string(cmdOut), err
}

func main() {

  statusString, _ := runCommand("git status")
  fmt.Println(filesToStage(statusString))


}


