package utils

import (
	"bufio"
	"io"
	"os"
)

func ForEachLineInFile(path string, f func(string, int)) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	forEachLineInStream(file, f)

	return nil
}

func forEachLineInStream(reader io.Reader, f func(string, int)) {
	bufferedReader := bufio.NewReader(reader)
	for i := 0; true; i++ {
		line, err := bufferedReader.ReadString('\n')
		if err != nil {
			break
		}
		f(line, i)
	}
}
