package file

import (
	"bufio"
	"fmt"
	"os"
)

// GetLinesFromFilename returns an array of strings (one per every line in the given file)
func GetLinesFromFilename(filename string) ([]string, error) {
	file, err := os.Open(filename)

	if err != nil {
		return nil, fmt.Errorf("error while opening file '%s'", filename)
	}

	defer file.Close()

	var data []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		data = append(data, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error while scanning file '%s': %v", filename, err)
	}

	return data, nil
}
