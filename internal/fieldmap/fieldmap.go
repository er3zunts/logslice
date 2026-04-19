// Package fieldmap provides functionality to remap or reorder log entry fields.
package fieldmap

import "github.com/logslice/logslice/internal/parser"

// Rule defines a field inclusion with an optional alias.
type Rule struct {
	Field string
	Alias string
}

// Apply returns new entries containing only the fields specified by rules.
// If a rule has a non-empty Alias, the field is stored under that name.
func Apply(entries []parser.Entry, rules []Rule) []parser.Entry {
	if len(rules) == 0 {
		return entries
	}
	out := make([]parser.Entry, 0, len(entries))
	for _, e := range entries {
		fields := make(map[string]interface{}, len(rules))
		for _, r := range rules {
			val, ok := e.Fields[r.Field]
			if !ok {
				continue
			}
			key := r.Field
			if r.Alias != "" {
				key = r.Alias
			}
			fields[key] = val
		}
		out = append(out, parser.Entry{Fields: fields})
	}
	return out
}

// ParseRules parses a slice of "field:alias" or "field" strings into Rules.
func ParseRules(specs []string) []Rule {
	rules := make([]Rule, 0, len(specs))
	for _, s := range specs {
		for i := 0; i < len(s); i++ {
			if s[i] == ':' {
				rules = append(rules, Rule{Field: s[:i], Alias: s[i+1:]})
				goto next
			}
		}
		rules = append(rules, Rule{Field: s})
	next:
	}
	return rules
}
