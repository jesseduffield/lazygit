<p align="center">

<img src="/logo.png" />
<br />
<a href="https://goreportcard.com/report/github.com/integrii/flaggy"><img src="https://goreportcard.com/badge/github.com/integrii/flaggy"></a>
<a href="https://travis-ci.org/integrii/flaggy"><img src="https://travis-ci.org/integrii/flaggy.svg?branch=master"></a>
<a href="http://godoc.org/github.com/integrii/flaggy"><img src="https://camo.githubusercontent.com/d48cccd1ce67ddf8ba7fc356ec1087f3f7aa6d12/68747470733a2f2f676f646f632e6f72672f6769746875622e636f6d2f6c696c65696f2f6c696c653f7374617475732e737667"></a>
<a href="http://unlicense.org/"><img src="https://img.shields.io/badge/license-Unlicense-blue.svg"></a>
<a href="https://cover.run/go?repo=github.com%2Fintegrii%2Fflaggy&tag=golang-1.10"><img src="https://cover.run/go/github.com/integrii/flaggy.svg?style=flat&tag=golang-1.10"></a>
<a href="https://github.com/avelino/awesome-go"><img src="https://cdn.rawgit.com/sindresorhus/awesome/d7305f38d29fed78fa85652e3a63e154dd8e8829/media/badge.svg"></a>
</p>

Sensible and _fast_ command-line flag parsing with excellent support for **subcommands** and **positional values**. Flags can be at any position. Flaggy has no required project or package layout like [Cobra requires](https://github.com/spf13/cobra/issues/641), and **no external dependencies**!

Check out the [godoc](http://godoc.org/github.com/integrii/flaggy), [examples directory](https://github.com/integrii/flaggy/tree/master/examples), and [examples in this readme](https://github.com/integrii/flaggy#super-simple-example) to get started quickly. You can also read the Flaggy introduction post with helpful examples [on my weblog](https://ericgreer.info/post/a-better-flags-package-for-go/).

# Installation

`go get -u github.com/integrii/flaggy`

# Key Features

- Very easy to use ([see examples below](https://github.com/integrii/flaggy#super-simple-example))
- 35 different flag types supported
- Any flag can be at any position
- Pretty and readable help output by default
- Positional subcommands
- Positional parameters
- Suggested subcommands when a subcommand is typo'd
- Nested subcommands
- Both global and subcommand specific flags
- Both global and subcommand specific positional parameters
- [Customizable help templates for both the global command and subcommands](https://github.com/integrii/flaggy/blob/master/examples/customTemplate/main.go)
- Customizable appended/prepended help messages for both the global command and subcommands
- Simple function that displays help followed by a custom message string
- Flags and subcommands may have both a short and long name
- Unlimited trailing arguments after a `--`
- Flags can use a single dash or double dash (`--flag`, `-flag`, `-f`, `--f`)
- Flags can have `=` assignment operators, or use a space (`--flag=value`, `--flag value`)
- Flags support single quote globs with spaces (`--flag 'this is all one value'`)
- Flags of slice types can be passed multiple times (`-f one -f two -f three`)
- Optional but default version output with `--version`
- Optional but default help output with `-h` or `--help`
- Optional but default help output when any invalid or unknown parameter is passed
- It's _fast_. All flag and subcommand parsing takes less than `1ms` in most programs.

# Example Help Output

```
testCommand - Description goes here.  Get more information at http://flaggy.
This is a prepend for help

  Usage:
    testCommand [subcommandA|subcommandB|subcommandC] [testPositionalA] [testPositionalB]

  Positional Variables:
    testPositionalA   Test positional A does some things with a positional value. (Required)
    testPositionalB   Test positional B does some less than serious things with a positional value.

  Subcommands:
    subcommandA (a)   Subcommand A is a command that does stuff
    subcommandB (b)   Subcommand B is a command that does other stuff
    subcommandC (c)   Subcommand C is a command that does SERIOUS stuff

  Flags:
       --version        Displays the program version string.
    -h --help           Displays help with available flag, subcommand, and positional value parameters.
    -s --stringFlag     This is a test string flag that does some stringy string stuff.
    -i --intFlg         This is a test int flag that does some interesting int stuff. (default: 5)
    -b --boolFlag       This is a test bool flag that does some booly bool stuff. (default: true)
    -d --durationFlag   This is a test duration flag that does some untimely stuff. (default: 1h23s)

This is an append for help
This is a help add-on message
```

# Super Simple Example

`./yourApp -f test`

```go
// Declare variables and their defaults
var stringFlag = "defaultValue"

// Add a flag
flaggy.String(&stringFlag, "f", "flag", "A test string flag")

// Parse the flag
flaggy.Parse()

// Use the flag
print(stringFlag)
```


# Example with Subcommand

`./yourApp subcommandExample -f test`

```go
// Declare variables and their defaults
var stringFlag = "defaultValue"

// Create the subcommand
subcommand := flaggy.NewSubcommand("subcommandExample")

// Add a flag to the subcommand
subcommand.String(&stringFlag, "f", "flag", "A test string flag")

// Add the subcommand to the parser at position 1
flaggy.AttachSubcommand(subcommand, 1)

// Parse the subcommand and all flags
flaggy.Parse()

// Use the flag
print(stringFlag)
```

# Example with Nested Subcommands, Various Flags and Trailing Arguments

`./yourApp subcommandExample --flag=5 nestedSubcommand -t test -y -- trailingArg`

```go
// Declare variables and their defaults
var stringFlagF = "defaultValueF"
var intFlagT = 3
var boolFlagB bool

// Create the subcommands
subcommandExample := flaggy.NewSubcommand("subcommandExample")
nestedSubcommand := flaggy.NewSubcommand("nestedSubcommand")

// Add a flag to both subcommands
subcommandExample.String(&stringFlagF, "t", "testFlag", "A test string flag")
nestedSubcommand.Int(&intFlagT, "f", "flag", "A test int flag")

// add a global bool flag for fun
flaggy.Bool(&boolFlagB, "y", "yes", "A sample boolean flag")

// attach the nested subcommand to the parent subcommand at position 1
subcommandExample.AttachSubcommand(nestedSubcommand, 1)
// attach the base subcommand to the parser at position 1
flaggy.AttachSubcommand(subcommandExample, 1)

// Parse everything, then use the flags and trailing arguments
flaggy.Parse()
print(stringFlagF)
print(intFlagT)
print(boolFlagB)
print(flaggy.TrailingArguments[0])
```

# Supported Flag Types

Flaggy has specific flag types for all basic types included in go as well as a slice of any of those types.  This includes all of the following types:

- string and []string
- bool and []bool
- all int types and all []int types
- all float types and all []float types
- all uint types and all []uint types

Other more specific types can also be used as flag types.  They will be automatically parsed using the standard parsing functions included with those types in those packages.  This includes:

- net.IP
- []net.IP
- net.HardwareAddr
- []net.HardwareAddr
- net.IPMask
- []net.IPMask
- time.Duration
- []time.Duration

# An Example Program

Best practice when using flaggy includes setting your program's name, description, and version (at build time) as shown in this example program.

```go
package main

import "github.com/integrii/flaggy"

// Make a variable for the version which will be set at build time.
var version = "unknown"

// Keep subcommands as globals so you can easily check if they were used later on.
var mySubcommand *flaggy.Subcommand

// Setup the variables you want your incoming flags to set.
var testVar string

// If you would like an environment variable as the default for a value, just populate the flag
// with the value of the environment by default.  If the flag corresponding to this value is not
// used, then it will not be changed.
var myVar = os.Getenv("MY_VAR")


func init() {
  // Set your program's name and description.  These appear in help output.
  flaggy.SetName("Test Program")
  flaggy.SetDescription("A little example program")

  // You can disable various things by changing bools on the default parser
  // (or your own parser if you have created one).
  flaggy.DefaultParser.ShowHelpOnUnexpected = false

  // You can set a help prepend or append on the default parser.
  flaggy.DefaultParser.AdditionalHelpPrepend = "http://github.com/integrii/flaggy"
  
  // Add a flag to the main program (this will be available in all subcommands as well).
  flaggy.String(&testVar, "tv", "testVariable", "A variable just for testing things!")

  // Create any subcommands and set their parameters.
  mySubcommand = flaggy.NewSubcommand("mySubcommand")
  mySubcommand.Description = "My great subcommand!"
  
  // Add a flag to the subcommand.
  mySubcommand.String(&myVar, "mv", "myVariable", "A variable just for me!")

  // Set the version and parse all inputs into variables.
  flaggy.SetVersion(version)
  flaggy.Parse()
}

func main(){
    if mySubcommand.Used {
      ...
    }
}
```

Then, you can use the following build command to set the `version` variable in the above program at build time.

```bash
# build your app and set the version string
$ go build -ldflags='-X main.version=1.0.3-a3db3'
$ ./yourApp version
Version: 1.0.3-a3db3
$ ./yourApp --help
Test Program - A little example program
http://github.com/integrii/flaggy
```

# Contributions

Please feel free to open an issue if you find any bugs or see any features that make sense. Pull requests will be reviewed and accepted if they make sense, but it is always wise to submit a proposal issue before any major changes.
