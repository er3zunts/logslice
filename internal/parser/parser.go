package parser

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// Entry represents a single parsed log line as a generic map.
type Entry map[string]interface{}

// String returns the entry serialized back to a JSON string.
func (e Entry) String() string {
	b, err := json.Marshal(e)
	if err != nil {
		return fmt.Sprintf("%v", map[string]interface{}(e))
	}
	return string(b)
}

// ParseJSONLines reads newline-delimited JSON from r and returns a slice of Entry.
func ParseJSONLines(r io.Reader) ([]Entry, error) {
	var entries []Entry
	scanner := bufio.NewScanner(r)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		var entry Entry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			return nil, fmt.Errorf("line %d: invalid JSON: %w", lineNum, err)
		}
		entries = append(entries, entry)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("reading input: %w", err)
	}
	return entries, nil
}
