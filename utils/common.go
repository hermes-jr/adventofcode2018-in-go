package utils

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

// ReadFile automatically appends level_XX directory and reads
// file contents into slice of strings
func ReadFile(fileName string) []string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	fname := filepath.Base(ex)
	re := regexp.MustCompile("(level_\\d\\d)")
	match := re.FindStringSubmatch(fname)

	var result []string
	result = make([]string, 0, 100)

	file, err := os.Open(filepath.Join(match[1], fileName))
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			return
		}
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		result = append(result, scanner.Text())
	}

	return result
}
