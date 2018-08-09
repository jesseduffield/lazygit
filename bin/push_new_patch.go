// call from project root with
// go run bin/push_new_patch.go

// goreleaser expects a $GITHUB_TOKEN env variable to be defined
// in order to push the release got github

package main

import (
	"fmt"
	"io/ioutil"
	"log"
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
	runCommand("git", "commit", "-m", "\"bump version to "+newVersion+"\"", "--", "VERSION")
	runCommand("git", "push")
	runCommand("git", "tag", newVersion)
	runCommand("git", "push", "origin", newVersion)
	runCommand("goreleaser", "--rm-dist")
}

func runCommand(args ...string) {
	fmt.Println(strings.Join(args, " "))
	output, err := exec.Command(args[0], args[1:]...).CombinedOutput()
	if err != nil {
		panic(err.Error())
	}
	log.Print(string(output))
}
