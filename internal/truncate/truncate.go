// Package truncate provides field value truncation for log entries.
package truncate

import (
	"fmt"
	"strings"

	"github.com/user/logslice/internal/parser"
)

// Options controls truncation behavior.
type Options struct {
	MaxLen int
	Fields []string // if empty, truncate all string fields
	Suffix string   // appended when truncated, e.g. "..."
}

// Apply truncates string field values in each entry according to opts.
func Apply(entries []parser.Entry, opts Options) []parser.Entry {
	if opts.MaxLen <= 0 {
		return entries
	}
	if opts.Suffix == "" {
		opts.Suffix = "..."
	}

	fieldSet := make(map[string]bool, len(opts.Fields))
	for _, f := range opts.Fields {
		fieldSet[f] = true
	}

	result := make([]parser.Entry, len(entries))
	for i, e := range entries {
		newFields := make(map[string]interface{}, len(e.Fields))
		for k, v := range e.Fields {
			if len(opts.Fields) == 0 || fieldSet[k] {
				if s, ok := v.(string); ok {
					v = truncateString(s, opts.MaxLen, opts.Suffix)
				}
			}
			newFields[k] = v
		}
		result[i] = parser.Entry{Fields: newFields}
	}
	return result
}

func truncateString(s string, maxLen int, suffix string) string {
	if len(s) <= maxLen {
		return s
	}
	cutoff := maxLen - len(suffix)
	if cutoff < 0 {
		cutoff = 0
	}
	return fmt.Sprintf("%s%s", strings.TrimRight(s[:cutoff], " "), suffix)
}
