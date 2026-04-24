// Package coalesce provides functionality to merge multiple fields into a
// single output field, using the first non-empty value found.
package coalesce

import (
	"fmt"
	"strings"

	"github.com/logslice/logslice/internal/parser"
)

// Rule defines a coalesce operation: pick the first non-empty value from
// Sources and write it to Dest. If none of the source fields are present or
// non-empty the Dest field is left untouched (or omitted).
type Rule struct {
	Dest    string
	Sources []string
}

// ParseRules converts raw string specs of the form "dest:src1,src2,..." into
// a slice of Rule values. It returns an error for any malformed spec.
func ParseRules(specs []string) ([]Rule, error) {
	rules := make([]Rule, 0, len(specs))
	for _, spec := range specs {
		dest, rest, ok := strings.Cut(spec, ":")
		if !ok || dest == "" || rest == "" {
			return nil, fmt.Errorf("coalesce: invalid rule %q (want dest:src1,src2,...)", spec)
		}
		srcs := strings.Split(rest, ",")
		for i, s := range srcs {
			srcs[i] = strings.TrimSpace(s)
		}
		rules = append(rules, Rule{Dest: strings.TrimSpace(dest), Sources: srcs})
	}
	return rules, nil
}

// Apply runs all rules against every entry and returns a new slice with the
// coalesced fields populated.
func Apply(entries []parser.Entry, rules []Rule) []parser.Entry {
	if len(rules) == 0 {
		return entries
	}
	out := make([]parser.Entry, len(entries))
	for i, e := range entries {
		newFields := make(map[string]interface{}, len(e.Fields))
		for k, v := range e.Fields {
			newFields[k] = v
		}
		for _, r := range rules {
			for _, src := range r.Sources {
				val, exists := newFields[src]
				if !exists {
					continue
				}
				if s, ok := val.(string); ok && s == "" {
					continue
				}
				newFields[r.Dest] = val
				break
			}
		}
		out[i] = parser.Entry{Fields: newFields}
	}
	return out
}
