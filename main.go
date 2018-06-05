package main

import (
  "flag"
  "fmt"
  "time"
)

var (
  startTime time.Time
  debugging bool
)

func main() {
  debuggingPointer := flag.Bool("debug", false, "a boolean")
  flag.Parse()
  debugging = *debuggingPointer
  fmt.Println(homeDirectory() + "/go/src/github.com/jesseduffield/lazygit/development.log")
  devLog("\n\n\n\n\n\n\n\n\n\n")
  startTime = time.Now()
  verifyInGitRepo()
  run()
}
