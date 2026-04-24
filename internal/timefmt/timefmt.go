// Package timefmt provides utilities for reformatting timestamp fields
// in log entries from one time layout to another.
package timefmt

import (
	"fmt"
	"time"

	"github.com/nicholasgasior/logslice/internal/parser"
)

// Rule defines a single timestamp reformat operation.
type Rule struct {
	Field     string
	InputFmt  string
	OutputFmt string
}

// ParseRules parses a slice of raw rule strings in the form
// "field:inputFmt:outputFmt" and returns a slice of Rule.
func ParseRules(raw []string) ([]Rule, error) {
	var rules []Rule
	for _, s := range raw {
		parts := splitN(s, ":", 3)
		if len(parts) != 3 {
			return nil, fmt.Errorf("timefmt: invalid rule %q, expected field:inputFmt:outputFmt", s)
		}
		rules = append(rules, Rule{
			Field:     parts[0],
			InputFmt:  parts[1],
			OutputFmt: parts[2],
		})
	}
	return rules, nil
}

// Apply reformats timestamp fields in each entry according to the given rules.
// Entries whose target field is missing or unparseable are passed through unchanged.
func Apply(entries []parser.Entry, rules []Rule) []parser.Entry {
	out := make([]parser.Entry, 0, len(entries))
	for _, e := range entries {
		ne := parser.Entry{Fields: make(map[string]interface{}, len(e.Fields))}
		for k, v := range e.Fields {
			ne.Fields[k] = v
		}
		for _, r := range rules {
			val, ok := ne.Fields[r.Field]
			if !ok {
				continue
			}
			str, ok := val.(string)
			if !ok {
				continue
			}
			t, err := time.Parse(r.InputFmt, str)
			if err != nil {
				continue
			}
			ne.Fields[r.Field] = t.Format(r.OutputFmt)
		}
		out = append(out, ne)
	}
	return out
}

// splitN splits s by sep at most n times, returning up to n parts.
func splitN(s, sep string, n int) []string {
	var parts []string
	for i := 0; i < n-1; i++ {
		idx := indexOf(s, sep)
		if idx < 0 {
			break
		}
		parts = append(parts, s[:idx])
		s = s[idx+len(sep):]
	}
	parts = append(parts, s)
	return parts
}

func indexOf(s, sub string) int {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}
