package arguments

import "os"

type Evaler func(string) (string, error)
type Stater func(string) (os.FileInfo, error)
