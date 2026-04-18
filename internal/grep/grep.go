// Package grep provides pattern matching across log entry fields.
package grep

import (
	"regexp"
	"strings"

	"github.com/user/logslice/internal/parser"
)

// Options controls how matching is performed.
type Options struct {
	Pattern     string
	Fields      []string // if empty, search all fields
	IgnoreCase  bool
	InvertMatch bool
}

// Match reports whether the entry matches the given options.
func Match(entry parser.Entry, opts Options) (bool, error) {
	pattern := opts.Pattern
	if opts.IgnoreCase {
		pattern = "(?i)" + pattern
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return false, err
	}

	matched := matchEntry(entry, re, opts.Fields)
	if opts.InvertMatch {
		return !matched, nil
	}
	return matched, nil
}

func matchEntry(entry parser.Entry, re *regexp.Regexp, fields []string) bool {
	if len(fields) == 0 {
		for _, v := range entry.Fields {
			if re.MatchString(strings.TrimSpace(v)) {
				return true
			}
		}
		return false
	}
	for _, f := range fields {
		val, ok := entry.Fields[f]
		if ok && re.MatchString(strings.TrimSpace(val)) {
			return true
		}
	}
	return false
}

// Filter returns only entries that match opts.
func Filter(entries []parser.Entry, opts Options) ([]parser.Entry, error) {
	var result []parser.Entry
	for _, e := range entries {
		ok, err := Match(e, opts)
		if err != nil {
			return nil, err
		}
		if ok {
			result = append(result, e)
		}
	}
	return result, nil
}
