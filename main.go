package main

import (
  "time"
)

// StartTime : The starting time of the app
var StartTime time.Time

func main() {
  devLog("\n\n\n\n\n\n\n\n\n\n")
  StartTime = time.Now()
  verifyInGitRepo()
  run()
}
