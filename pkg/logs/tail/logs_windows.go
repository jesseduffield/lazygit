//go:build windows
// +build windows

package tail

import (
	"bufio"
	"log"
	"os"
	"strings"
	"time"

	"github.com/aybabtme/humanlog"
)

func tailLogsForPlatform(logFilePath string, opts *humanlog.HandlerOptions) {
	var lastModified int64 = 0
	var lastOffset int64 = 0
	for {
		stat, err := os.Stat(logFilePath)
		if err != nil {
			log.Fatal(err)
		}
		if stat.ModTime().Unix() > lastModified {
			err = tailFrom(lastOffset, logFilePath, opts)
			if err != nil {
				log.Fatal(err)
			}
		}
		lastOffset = stat.Size()
		time.Sleep(1 * time.Second)
	}
}

func openAndSeek(filepath string, offset int64) (*os.File, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}

	_, err = file.Seek(offset, 0)
	if err != nil {
		_ = file.Close()
		return nil, err
	}
	return file, nil
}

func tailFrom(lastOffset int64, logFilePath string, opts *humanlog.HandlerOptions) error {
	file, err := openAndSeek(logFilePath, lastOffset)
	if err != nil {
		return err
	}

	fileScanner := bufio.NewScanner(file)
	var lines []string
	for fileScanner.Scan() {
		lines = append(lines, fileScanner.Text())
	}
	file.Close()
	lineCount := len(lines)
	lastTen := lines
	if lineCount > 10 {
		lastTen = lines[lineCount-10:]
	}
	for _, line := range lastTen {
		reader := strings.NewReader(line)
		if err := humanlog.Scanner(reader, os.Stdout, opts); err != nil {
			log.Fatal(err)
		}
	}
	return nil
}
