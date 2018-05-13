// Go has various value types including strings,
// integers, floats, booleans, etc. Here are a few
// basic examples.

package main

import (
  "fmt"
  // "log"
  "os/exec"
  "os"
)

func main() {

  var (
    cmdOut []byte
    err    error
  )
  cmdName := "git"
  cmdArgs := []string{"rev-parse", "--verify", "HEAD"}
  if cmdOut, err = exec.Command(cmdName, cmdArgs...).Output(); err != nil {
    fmt.Fprintln(os.Stderr, "There was an error running git rev-parse command: ", err)
    fmt.Println(string(cmdOut))
    os.Exit(1)
  }
  sha := string(cmdOut)
  firstSix := sha[:6]
  fmt.Println("The first six chars of the SHA at HEAD in this repo are", firstSix)

}
