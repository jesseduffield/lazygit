// call from project root with
// go run scripts/push_new_patch/main.go

// goreleaser expects a $GITHUB_TOKEN env variable to be defined
// in order to push the release got github

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func main() {

	version, err := ioutil.ReadFile("VERSION")
	if err != nil {
		log.Panicln(err.Error())
	}
	stringVersion := string(version)
	fmt.Println("VERSION was " + stringVersion)

	runCommand("git", "pull")

	splitVersion := strings.Split(stringVersion, ".")
	patch := splitVersion[len(splitVersion)-1]
	newPatch, err := strconv.Atoi(patch)
	splitVersion[len(splitVersion)-1] = strconv.FormatInt(int64(newPatch)+1, 10)
	newVersion := strings.Join(splitVersion, ".")

	err = ioutil.WriteFile("VERSION", []byte(newVersion), 0644)
	if err != nil {
		log.Panicln(err.Error())
	}

	runCommand("git", "add", "VERSION")
	runCommand("git", "commit", "-m", "bump version to "+newVersion, "--", "VERSION")
	runCommand("git", "push")
	runCommand("git", "tag", newVersion)
	runCommand("git", "push", "origin", newVersion)
	runCommand("goreleaser", "--rm-dist")
	runCommand("rm", "-rf", "dist")
}

func runCommand(args ...string) {
	fmt.Println(strings.Join(args, " "))
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		panic(err.Error())
	}
	err = cmd.Wait()
	if err != nil {
		panic(err.Error())
	}
}
