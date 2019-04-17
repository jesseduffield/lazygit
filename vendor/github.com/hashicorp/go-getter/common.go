package getter

import (
	"io/ioutil"
)

func tmpFile(dir, pattern string) (string, error) {
	f, err := ioutil.TempFile(dir, pattern)
	if err != nil {
		return "", err
	}
	f.Close()
	return f.Name(), nil
}
