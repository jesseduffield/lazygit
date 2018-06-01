package main

import (

  // "log"
  "fmt"
  "os"
  "os/exec"
  "strings"
  "time"

  "github.com/fatih/color"
)

// GitFile : A staged/unstaged file
type GitFile struct {
  Name               string
  HasStagedChanges   bool
  HasUnstagedChanges bool
  Tracked            bool
  Deleted            bool
  DisplayString      string
}

// Branch : A git branch
type Branch struct {
  Name          string
  Type          string
  BaseBranch    string
  DisplayString string
  DisplayColor  color.Attribute
}

// Commit : A git commit
type Commit struct {
  Sha           string
  Name          string
  Pushed        bool
  DisplayString string
}

func devLog(objects ...interface{}) {
  localLog(color.FgWhite, "/Users/jesseduffieldduffield/go/src/github.com/jesseduffield/gitgot/development.log", objects...)
}

func colorLog(colour color.Attribute, objects ...interface{}) {
  localLog(colour, "/Users/jesseduffieldduffield/go/src/github.com/jesseduffield/gitgot/development.log", objects...)
}

func commandLog(objects ...interface{}) {
  localLog(color.FgWhite, "/Users/jesseduffieldduffield/go/src/github.com/jesseduffield/gitgot/commands.log", objects...)
  // localLog(color.FgWhite, "/Users/jesseduffieldduffield/go/src/github.com/jesseduffield/gitgot/development.log", objects...)
}

func localLog(colour color.Attribute, path string, objects ...interface{}) {
  f, _ := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
  defer f.Close()
  for _, object := range objects {
    colorFunction := color.New(colour).SprintFunc()
    f.WriteString(colorFunction(fmt.Sprint(object)) + "\n")
  }
}

// Map (from https://gobyexample.com/collection-functions)
func Map(vs []string, f func(string) string) []string {
  vsm := make([]string, len(vs))
  for i, v := range vs {
    vsm[i] = f(v)
  }
  return vsm
}

func includes(list []string, a string) bool {
  for _, b := range list {
    if b == a {
      return true
    }
  }
  return false
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

func runDirectCommand(command string) (string, error) {
  timeStart := time.Now()

  commandLog(command)
  cmdOut, err := exec.Command("bash", "-c", command).CombinedOutput()
  devLog("run direct command time for command: ", command, time.Now().Sub(timeStart))

  return string(cmdOut), err
}

func branchStringParts(branchString string) (string, string) {
  splitBranchName := strings.Split(branchString, "\t")
  return splitBranchName[0], splitBranchName[1]
}

// branchPropertiesFromName : returns branch type, base, and color
func branchPropertiesFromName(name string) (string, string, color.Attribute) {
  if strings.Contains(name, "feature/") {
    return "feature", "develop", color.FgGreen
  } else if strings.Contains(name, "bugfix/") {
    return "bugfix", "develop", color.FgYellow
  } else if strings.Contains(name, "hotfix/") {
    return "hotfix", "master", color.FgRed
  }
  return "other", name, color.FgWhite
}

func coloredString(str string, colour color.Attribute) string {
  return color.New(colour).SprintFunc()(fmt.Sprint(str))
}

func withPadding(str string, padding int) string {
  return str + strings.Repeat(" ", padding-len(str))
}

func branchFromLine(line string, index int) Branch {
  recency, name := branchStringParts(line)
  branchType, branchBase, colour := branchPropertiesFromName(name)
  if index == 0 {
    recency = "  *"
  }
  displayString := withPadding(recency, 4) + coloredString(name, colour)
  return Branch{
    Name:          name,
    Type:          branchType,
    BaseBranch:    branchBase,
    DisplayString: displayString,
    DisplayColor:  colour,
  }
}

func getGitBranches() []Branch {
  branches := make([]Branch, 0)
  // check if there are any branches
  branchCheck, _ := runDirectCommand("git branch")
  if branchCheck == "" {
    return branches
  }
  rawString, _ := runDirectCommand(getBranchesCommand)
  branchLines := splitLines(rawString)
  for i, line := range branchLines {
    branches = append(branches, branchFromLine(line, i))
  }
  return branches
}

func getGitStatusFiles() []GitFile {
  statusOutput, _ := getGitStatus()
  statusStrings := splitLines(statusOutput)
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

func gitCheckout(branch string, force bool) (string, error) {
  forceArg := ""
  if force {
    forceArg = "--force "
  }
  return runCommand("git checkout " + forceArg + branch)
}

func runCommand(command string) (string, error) {
  startTime := time.Now()
  commandLog(command)
  splitCmd := strings.Split(command, " ")
  cmdOut, err := exec.Command(splitCmd[0], splitCmd[1:]...).CombinedOutput()
  devLog("run command time: ", time.Now().Sub(startTime))
  return string(cmdOut), err
}

func openFile(filename string) (string, error) {
  return runCommand("open " + filename)
}

func sublimeOpenFile(filename string) (string, error) {
  return runCommand("subl " + filename)
}

func getBranchDiff(branch string, baseBranch string) (string, error) {
  return runCommand("git diff --color " + baseBranch + "..." + branch)
}

func verifyInGitRepo() {
  if output, err := runCommand("git status"); err != nil {
    fmt.Println(output)
    os.Exit(1)
  }
}

func getCommits() []Commit {
  pushables := gitCommitsToPush()
  log := getLog()
  commits := make([]Commit, 0)
  // now we can split it up and turn it into commits
  lines := splitLines(log)
  for _, line := range lines {
    splitLine := strings.Split(line, " ")
    sha := splitLine[0]
    pushed := includes(pushables, sha)
    commits = append(commits, Commit{
      Sha:           sha,
      Name:          strings.Join(splitLine[1:], " "),
      Pushed:        pushed,
      DisplayString: strings.Join(splitLine, " "),
    })
  }
  return commits
}

func getLog() string {
  // currently limiting to 30 for performance reasons
  // TODO: add lazyloading when you scroll down
  result, err := runDirectCommand("git log --oneline -30")
  if err != nil {
    // assume if there is an error there are no commits yet for this branch
    return ""
  }
  return result
}

func gitIgnore(filename string) {
  if _, err := runDirectCommand("echo '" + filename + "' >> .gitignore"); err != nil {
    panic(err)
  }
}

func gitShow(sha string) string {
  result, err := runDirectCommand("git show --color " + sha)
  if err != nil {
    panic(err)
  }
  return result
}

func getDiff(file GitFile) string {
  cachedArg := ""
  if file.HasStagedChanges {
    cachedArg = "--cached "
  }
  deletedArg := ""
  if file.Deleted {
    deletedArg = "-- "
  }
  trackedArg := ""
  if !file.Tracked {
    trackedArg = "--no-index /dev/null "
  }
  command := "git diff --color " + cachedArg + deletedArg + trackedArg + file.Name
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

func removeFile(file GitFile) error {
  // if the file isn't tracked, we assume you want to delete it
  if !file.Tracked {
    _, err := runCommand("rm -rf ./" + file.Name)
    return err
  }
  // if the file is tracked, we assume you want to just check it out
  _, err := runCommand("git checkout " + file.Name)
  return err
}

func gitCommit(message string) error {
  _, err := runDirectCommand("git commit -m \"" + message + "\"")
  return err
}

func gitPull() (string, error) {
  return runDirectCommand("git pull --no-edit")
}

func gitPush() (string, error) {
  return runDirectCommand("git push -u")
}

func gitSquashPreviousTwoCommits(message string) (string, error) {
  return runDirectCommand("git reset --soft head^ && git commit --amend -m \"" + message + "\"")
}

func gitRenameCommit(message string) (string, error) {
  return runDirectCommand("git commit --allow-empty --amend -m \"" + message + "\"")
}

func gitUpstreamDifferenceCount() (string, string) {
  // TODO: deal with these errors which appear when we haven't yet pushed a feature branch
  pushableCount, err := runDirectCommand("git rev-list @{u}..head --count")
  if err != nil {
    return "?", "?"
  }
  pullableCount, err := runDirectCommand("git rev-list head..@{u} --count")
  if err != nil {
    return "?", "?"
  }
  return strings.Trim(pullableCount, " \n"), strings.Trim(pushableCount, " \n")
}

func gitCommitsToPush() []string {
  pushables, err := runDirectCommand("git rev-list @{u}..head --abbrev-commit")
  if err != nil {
    return make([]string, 0)
  }
  return splitLines(pushables)
}

func gitCurrentBranchName() string {
  branchName, err := runDirectCommand("git rev-parse --abbrev-ref HEAD")
  // if there is an error, assume there are no branches yet
  if err != nil {
    return ""
  }
  return branchName
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
  done | sed 's/ days /d /g' | sed 's/ weeks /w /g' | sed 's/ hours /h /g' | sed 's/ minutes /m /g' | sed 's/ seconds /m /g' | sed 's/ago//g' | tr -d ' '
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
