package arguments

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
)

func New(args []string, workingDir string, evaler Evaler, stater Stater) (*ParsedArguments, error) {
	if len(args) == 0 {
		return nil, errors.New("argument parsing requires at least one argument")
	}

	fs := flag.NewFlagSet("counterfeiter", flag.ContinueOnError)
	fakeNameFlag := fs.String(
		"fake-name",
		"",
		"The name of the fake struct",
	)

	outputPathFlag := fs.String(
		"o",
		"",
		"The file or directory to which the generated fake will be written",
	)

	packageFlag := fs.Bool(
		"p",
		false,
		"Whether or not to generate a package shim",
	)
	generateFlag := fs.Bool(
		"generate",
		false,
		"Identify all //counterfeiter:generate directives in the current working directory and generate fakes for them",
	)
	headerFlag := fs.String(
		"header",
		"",
		"A path to a file that should be used as a header for the generated fake",
	)
	quietFlag := fs.Bool(
		"q",
		false,
		"Suppress status statements",
	)
	helpFlag := fs.Bool(
		"help",
		false,
		"Display this help",
	)

	err := fs.Parse(args[1:])
	if err != nil {
		return nil, err
	}
	if *helpFlag {
		return nil, errors.New(usage)
	}
	if len(fs.Args()) == 0 && !*generateFlag {
		return nil, errors.New(usage)
	}

	packageMode := *packageFlag
	result := &ParsedArguments{
		PrintToStdOut: any(args, "-"),
		GenerateInterfaceAndShimFromPackageDirectory: packageMode,
		GenerateMode: *generateFlag,
		HeaderFile:   *headerFlag,
		Quiet:        *quietFlag,
	}
	if *generateFlag {
		return result, nil
	}
	err = result.parseSourcePackageDir(packageMode, workingDir, evaler, stater, fs.Args())
	if err != nil {
		return nil, err
	}
	result.parseInterfaceName(packageMode, fs.Args())
	result.parseFakeName(packageMode, *fakeNameFlag, fs.Args())
	result.parseOutputPath(packageMode, workingDir, *outputPathFlag, fs.Args())
	result.parseDestinationPackageName(packageMode, fs.Args())
	result.parsePackagePath(packageMode, fs.Args())
	return result, nil
}

func (a *ParsedArguments) PrettyPrint() {
	b, _ := json.Marshal(a)
	fmt.Println(string(b))
}

func (a *ParsedArguments) parseInterfaceName(packageMode bool, args []string) {
	if packageMode {
		a.InterfaceName = ""
		return
	}
	if len(args) == 1 {
		fullyQualifiedInterface := strings.Split(args[0], ".")
		a.InterfaceName = fullyQualifiedInterface[len(fullyQualifiedInterface)-1]
	} else {
		a.InterfaceName = args[1]
	}
}

func (a *ParsedArguments) parseSourcePackageDir(packageMode bool, workingDir string, evaler Evaler, stater Stater, args []string) error {
	if packageMode {
		a.SourcePackageDir = args[0]
		return nil
	}
	if len(args) <= 1 {
		return nil
	}
	s, err := getSourceDir(args[0], workingDir, evaler, stater)
	if err != nil {
		return err
	}
	a.SourcePackageDir = s
	return nil
}

func (a *ParsedArguments) parseFakeName(packageMode bool, fakeName string, args []string) {
	if packageMode {
		a.parsePackagePath(packageMode, args)
		a.FakeImplName = strings.ToUpper(path.Base(a.PackagePath))[:1] + path.Base(a.PackagePath)[1:]
		return
	}
	if fakeName == "" {
		fakeName = "Fake" + fixupUnexportedNames(a.InterfaceName)
	}
	a.FakeImplName = fakeName
}

func (a *ParsedArguments) parseOutputPath(packageMode bool, workingDir string, outputPath string, args []string) {
	outputPathIsFilename := false
	if strings.HasSuffix(outputPath, ".go") {
		outputPathIsFilename = true
	}
	snakeCaseName := strings.ToLower(camelRegexp.ReplaceAllString(a.FakeImplName, "${1}_${2}"))

	if outputPath != "" {
		if !filepath.IsAbs(outputPath) {
			outputPath = filepath.Join(workingDir, outputPath)
		}
		a.OutputPath = outputPath
		if !outputPathIsFilename {
			a.OutputPath = filepath.Join(a.OutputPath, snakeCaseName+".go")
		}
		return
	}

	if packageMode {
		a.parseDestinationPackageName(packageMode, args)
		a.OutputPath = path.Join(workingDir, a.DestinationPackageName, snakeCaseName+".go")
		return
	}

	d := workingDir
	if len(args) > 1 {
		d = a.SourcePackageDir
	}
	a.OutputPath = filepath.Join(d, packageNameForPath(d), snakeCaseName+".go")
}

func (a *ParsedArguments) parseDestinationPackageName(packageMode bool, args []string) {
	if packageMode {
		a.parsePackagePath(packageMode, args)
		a.DestinationPackageName = path.Base(a.PackagePath) + "shim"
		return
	}

	a.DestinationPackageName = restrictToValidPackageName(filepath.Base(filepath.Dir(a.OutputPath)))
}

func (a *ParsedArguments) parsePackagePath(packageMode bool, args []string) {
	if packageMode {
		a.PackagePath = args[0]
		return
	}
	if len(args) == 1 {
		fullyQualifiedInterface := strings.Split(args[0], ".")
		a.PackagePath = strings.Join(fullyQualifiedInterface[:len(fullyQualifiedInterface)-1], ".")
	} else {
		a.InterfaceName = args[1]
	}

	if a.PackagePath == "" {
		a.PackagePath = a.SourcePackageDir
	}
}

type ParsedArguments struct {
	GenerateInterfaceAndShimFromPackageDirectory bool

	SourcePackageDir string // abs path to the dir containing the interface to fake
	PackagePath      string // package path to the package containing the interface to fake
	OutputPath       string // path to write the fake file to

	DestinationPackageName string // often the base-dir for OutputPath but must be a valid package name

	InterfaceName string // the interface to counterfeit
	FakeImplName  string // the name of the struct implementing the given interface

	PrintToStdOut bool
	GenerateMode  bool
	Quiet         bool

	HeaderFile string
}

func fixupUnexportedNames(interfaceName string) string {
	asRunes := []rune(interfaceName)
	if len(asRunes) == 0 || !unicode.IsLower(asRunes[0]) {
		return interfaceName
	}
	asRunes[0] = unicode.ToUpper(asRunes[0])
	return string(asRunes)
}

var camelRegexp = regexp.MustCompile("([a-z])([A-Z])")

func packageNameForPath(pathToPackage string) string {
	_, packageName := filepath.Split(pathToPackage)
	return packageName + "fakes"
}

func getSourceDir(path string, workingDir string, evaler Evaler, stater Stater) (string, error) {
	if !filepath.IsAbs(path) {
		path = filepath.Join(workingDir, path)
	}

	evaluatedPath, err := evaler(path)
	if err != nil {
		return "", fmt.Errorf("No such file/directory/package [%s]: %v", path, err)
	}

	stat, err := stater(evaluatedPath)
	if err != nil {
		return "", fmt.Errorf("No such file/directory/package [%s]: %v", path, err)
	}

	if !stat.IsDir() {
		return filepath.Dir(path), nil
	}
	return path, nil
}

func any(slice []string, needle string) bool {
	for _, str := range slice {
		if str == needle {
			return true
		}
	}

	return false
}

func restrictToValidPackageName(input string) string {
	return strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' {
			return r
		} else {
			return -1
		}
	}, input)
}
