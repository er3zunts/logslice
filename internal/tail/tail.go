// Package tail provides functionality to read the last N lines of a log file.
package tail

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// Lines reads the last n lines from the given file path.
// If n <= 0, all lines are returned.
func Lines(path string, n int) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("tail: open file: %w", err)
	}
	defer f.Close()
	return readLastN(f, n)
}

// readLastN reads from r and returns the last n lines.
func readLastN(r io.Reader, n int) ([]string, error) {
	scanner := bufio.NewScanner(r)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("tail: scan: %w", err)
	}
	if n <= 0 || n >= len(lines) {
		return lines, nil
	}
	return lines[len(lines)-n:], nil
}
