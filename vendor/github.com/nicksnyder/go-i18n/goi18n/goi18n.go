package main

import (
	"flag"
	"fmt"
	"os"
)

type command interface {
	execute() error
	parse(arguments []string)
}

func main() {
	flag.Usage = usage

	if len(os.Args) == 1 {
		usage()
	}

	var cmd command

	switch os.Args[1] {
	case "merge":
		cmd = &mergeCommand{}
		cmd.parse(os.Args[2:])
	case "constants":
		cmd = &constantsCommand{}
		cmd.parse(os.Args[2:])
	default:
		cmd = &mergeCommand{}
		cmd.parse(os.Args[1:])
	}

	if err := cmd.execute(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func usage() {
	fmt.Printf(`goi18n manages translation files.

Usage:

    goi18n merge     Merge translation files
    goi18n constants Generate constant file from translation file

For more details execute:

    goi18n [command] -help

`)
	os.Exit(1)
}
