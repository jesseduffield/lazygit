package arguments

const usage = `
USAGE
	counterfeiter
		[-generate>] [-o <output-path>] [-p] [--fake-name <fake-name>]
		[-header <header-file>]
		[<source-path>] <interface> [-]

ARGUMENTS
	source-path
		Path to the file or directory containing the interface to fake.
		In package mode (-p), source-path should instead specify the path
		of the input package; alternatively you can use the package name
		(e.g. "os") and the path will be inferred from your GOROOT.

	interface
		If source-path is specified: Name of the interface to fake.
		If no source-path is specified: Fully qualified interface path of the interface to fake.
    If -p is specified, this will be the name of the interface to generate.

	example:
		# writes "FakeStdInterface" to ./packagefakes/fake_std_interface.go
		counterfeiter package/subpackage.StdInterface

	'-' argument
		Write code to standard out instead of to a file

OPTIONS
	-generate
		Identify all //counterfeiter:generate directives in .go file in the
		current working directory and generate fakes for them. You can pass
		arguments as usual.

		NOTE: This is not the same as //go:generate directives
		(used with the 'go generate' command), but it can be combined with
		go generate by adding the following to a .go file:

		# runs counterfeiter in generate mode
		//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate

	example:
		Add the following to a .go file:

		//counterfeiter:generate . MyInterface
		//counterfeiter:generate . MyOtherInterface
		//counterfeiter:generate . MyThirdInterface

		# run counterfeiter
		counterfeiter -generate
		# writes "FakeMyInterface" to ./mypackagefakes/fake_my_interface.go
		# writes "FakeMyOtherInterface" to ./mypackagefakes/fake_my_other_interface.go
		# writes "FakeMyThirdInterface" to ./mypackagefakes/fake_my_third_interface.go

	-o
		Path to the file or directory for the generated fakes.
		This also determines the package name that will be used.
		By default, the generated fakes will be generated in
		the package "xyzfakes" which is nested in package "xyz",
		where "xyz" is the name of referenced package.

	example:
		# writes "FakeMyInterface" to ./mySpecialFakesDir/specialFake.go
		counterfeiter -o ./mySpecialFakesDir/specialFake.go ./mypackage MyInterface

		# writes "FakeMyInterface" to ./mySpecialFakesDir/fake_my_interface.go
		counterfeiter -o ./mySpecialFakesDir ./mypackage MyInterface

	-p
		Package mode:  When invoked in package mode, counterfeiter
		will generate an interface and shim implementation from a
		package in your module.  Counterfeiter finds the public methods
		in the package <source-path> and adds those method signatures
		to the generated interface <interface-name>.

	example:
		# generates os.go (interface) and osshim.go (shim) in ${PWD}/osshim
		counterfeiter -p os
		# now generate fake in ${PWD}/osshim/os_fake (fake_os.go)
		go generate osshim/...

	-header
		Path to the file which should be used as a header for all generated fakes.
		By default, no special header is used.
		This is useful to e.g. add a licence header to every fake.

		If the generate mode is used and both the "go:generate" and the
		"counterfeiter:generate" specify a header file, the header file from the
		"counterfeiter:generate" line takes precedence.

	example:
		# having the following code in a package ...
		//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -header ./generic.go.txt -generate
		//counterfeiter:generate -header ./specific.go.txt . MyInterface
		//counterfeiter:generate . MyOtherInterface
		//counterfeiter:generate . MyThirdInterface

		# ... generating the fakes ...
		go generate .

		# writes "FakeMyInterface" with ./specific.go.txt as a header
		# writes "FakeMyOtherInterface" & "FakeMyThirdInterface" with ./generic.go.txt as a header

	--fake-name
		Name of the fake struct to generate. By default, 'Fake' will
		be prepended to the name of the original interface. (ignored in
		-p mode)

	example:
		# writes "CoolThing" to ./mypackagefakes/cool_thing.go
		counterfeiter --fake-name CoolThing ./mypackage MyInterface
`
