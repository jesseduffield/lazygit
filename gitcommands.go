// Go has various value types including strings,
// integers, floats, booleans, etc. Here are a few
// basic examples.

package main

import (

  // "log"
  "fmt"
  "os/exec"
  "strings"
)

// GitFile : A staged/unstaged file
type GitFile struct {
  Name               string
  DisplayString      string
  HasStagedChanges   bool
  HasUnstagedChanges bool
  Tracked            bool
  Deleted            bool
}

// Branch : A git branch
type Branch struct {
  Name          string
  DisplayString string
  Type          string
  BaseBranch    string
}

// Map (from https://gobyexample.com/collection-functions)
func Map(vs []string, f func(string) string) []string {
  vsm := make([]string, len(vs))
  for i, v := range vs {
    vsm[i] = f(v)
  }
  return vsm
}

func mergeGitStatusFiles(oldGitFiles, newGitFiles []GitFile) []GitFile {
  if len(oldGitFiles) == 0 {
    return newGitFiles
  }

  result := make([]GitFile, 0)
  for _, oldGitFile := range oldGitFiles {
    for _, newGitFile := range newGitFiles {
      if oldGitFile.Name == newGitFile.Name {
        result = append(result, newGitFile)
        break
      }
    }
  }
  return result
}

func getGitBranchOutput() (string, error) {
  cmdOut, err := exec.Command("bash", "-c", getBranchesCommand).Output()
  return string(cmdOut), err
}

func branchNameFromString(branchString string) string {
  // because this has the recency at the beginning,
  // we need to split and take the second part
  splitBranchName := strings.Split(branchString, "\t")
  return splitBranchName[len(splitBranchName)-1]
}

func getGitBranches() []Branch {
  branches := make([]Branch, 0)
  rawString, _ := getGitBranchOutput()
  branchLines := splitLines(rawString)
  for _, line := range branchLines {
    name := branchNameFromString(line)
    var branchType string
    var baseBranch string
    if strings.Contains(line, "feature/") {
      branchType = "feature"
      baseBranch = "develop"
    } else if strings.Contains(line, "bugfix/") {
      branchType = "bugfix"
      baseBranch = "develop"
    } else if strings.Contains(line, "hotfix/") {
      branchType = "hotfix"
      baseBranch = "master"
    } else {
      branchType = "other"
      baseBranch = name
    }
    branches = append(branches, Branch{name, line, branchType, baseBranch})
  }
  devLog(fmt.Sprint(branches))
  return branches
}

func getGitStatusFiles() []GitFile {
  statusOutput, _ := getGitStatus()
  statusStrings := splitLines(statusOutput)
  devLog(fmt.Sprint(statusStrings))
  // a file can have both staged and unstaged changes
  // I'll probably end up ignoring the unstaged flag for now but might revisit
  // tracked, staged, unstaged

  gitFiles := make([]GitFile, 0)

  for _, statusString := range statusStrings {
    stagedChange := statusString[0:1]
    unstagedChange := statusString[1:2]
    filename := statusString[3:]
    tracked := statusString[0:2] != "??"
    gitFile := GitFile{
      Name:               filename,
      DisplayString:      statusString,
      HasStagedChanges:   tracked && stagedChange != " ",
      HasUnstagedChanges: !tracked || unstagedChange != " ",
      Tracked:            tracked,
      Deleted:            unstagedChange == "D" || stagedChange == "D",
    }
    gitFiles = append(gitFiles, gitFile)
  }
  return gitFiles
}

func gitCheckout(branch string, force bool) error {
  forceArg := ""
  if force {
    forceArg = "--force "
  }
  _, err := runCommand("git checkout " + forceArg + branch)
  return err
}

func runCommand(cmd string) (string, error) {
  splitCmd := strings.Split(cmd, " ")
  cmdOut, err := exec.Command(splitCmd[0], splitCmd[1:]...).Output()
  devLog(cmd)
  devLog(string(cmdOut))
  return string(cmdOut), err
}

func getBranchDiff(branch string, baseBranch string) (string, error) {
  return runCommand("git diff --color " + baseBranch + "..." + branch)
}

func getDiff(file GitFile) string {
  cachedArg := ""
  if file.HasStagedChanges {
    cachedArg = "--cached "
  }
  deletedArg := ""
  if file.Deleted || !file.Tracked {
    deletedArg = "--no-index /dev/null "
  }
  command := "git diff --color " + cachedArg + deletedArg + file.Name
  s, err := runCommand(command)
  if err != nil {
    // for now we assume an error means the file was deleted
    return s
  }
  return s
}

func stageFile(file string) error {
  _, err := runCommand("git add " + file)
  return err
}

func unStageFile(file string) error {
  _, err := runCommand("git reset HEAD " + file)
  return err
}

func getGitStatus() (string, error) {
  return runCommand("git status --untracked-files=all --short")
}

const getBranchesCommand = `set -e
git reflog -n100 --pretty='%cr|%gs' --grep-reflog='checkout: moving' HEAD | {
  seen=":"
  git_dir="$(git rev-parse --git-dir)"
  while read line; do
    date="${line%%|*}"
    branch="${line##* }"
    if ! [[ $seen == *:"${branch}":* ]]; then
      seen="${seen}${branch}:"
      if [ -f "${git_dir}/refs/heads/${branch}" ]; then
        printf "%s\t%s\n" "$date" "$branch"
      fi
    fi
  done | sed 's/ days /d /g' | sed 's/ weeks /w /g' | sed 's/ hours /h /g' | sed 's/ minutes /m /g' | sed 's/ago//g' | tr -d ' '
}
`

// func main() {
//   getGitStatusFiles()
// }

// func devLog(s string) {
//   f, _ := os.OpenFile("development.log", os.O_APPEND|os.O_WRONLY, 0644)
//   defer f.Close()

//   f.WriteString(s + "\n")
// }
