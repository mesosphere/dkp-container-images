package cli

import (
	"bufio"
	"io"
	"os"
	"strings"
)

// OpenFileOrStdin returns a reader from given file path or Stdin in path
// argument equals to `-`.
func OpenFileOrStdin(path string) (io.Reader, error) {
	if path == "-" {
		return os.Stdin, nil
	}

	input, err := os.Open(path)
	return input, err
}

// ReadImages parses lines to array of images. Empty lines and lines with `#` prefix (commented out)
// are ignored.
func ReadImages(input io.Reader) ([]string, error) {
	scanner := bufio.NewScanner(input)
	scanner.Split(bufio.ScanLines)

	images := []string{}
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "#") {
			continue
		}
		images = append(images, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return images, nil
}
