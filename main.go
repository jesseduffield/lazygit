package main

import "github.com/fatih/color"

func main() {
  verifyInGitRepo()
  a, b := gitUpstreamDifferenceCount()
  colorLog(color.FgRed, a, b)
  devLog("\n\n\n\n\n\n\n\n\n\n")
  run()
}
