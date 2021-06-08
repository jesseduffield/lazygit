package utils

import (
	"bufio"
	"os"
)

func ForEachLineInFile(path string, f func(string, int)) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	for i := 0; true; i++ {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		f(line, i)
	}

	return nil
}
