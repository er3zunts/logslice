// Package extract provides regex-based field extraction from string log fields.
// It applies named capture groups from a regular expression to create new fields
// on each log entry.
package extract

import (
	"fmt"
	"regexp"

	"github.com/yourorg/logslice/internal/parser"
)

// Rule defines a single extraction: the source field, the regex with named
// capture groups, and an optional prefix applied to each captured field name.
type Rule struct {
	Field  string
	Regexp *regexp.Regexp
	Prefix string
}

// ParseRules parses raw rule strings of the form "field:regex" or
// "field:prefix:regex". The regex must contain at least one named capture group.
func ParseRules(raw []string) ([]Rule, error) {
	rules := make([]Rule, 0, len(raw))
	for _, r := range raw {
		parts := splitN(r, ":", 3)
		if len(parts) < 2 {
			return nil, fmt.Errorf("extract: invalid rule %q: expected field:regex or field:prefix:regex", r)
		}
		var field, prefix, pattern string
		if len(parts) == 2 {
			field, pattern = parts[0], parts[1]
		} else {
			field, prefix, pattern = parts[0], parts[1], parts[2]
		}
		re, err := regexp.Compile(pattern)
		if err != nil {
			return nil, fmt.Errorf("extract: invalid regex in rule %q: %w", r, err)
		}
		if len(re.SubexpNames()) < 2 {
			return nil, fmt.Errorf("extract: regex in rule %q has no named capture groups", r)
		}
		rules = append(rules, Rule{Field: field, Regexp: re, Prefix: prefix})
	}
	return rules, nil
}

// Apply runs all extraction rules over the entries, adding captured fields to
// each entry. Entries where the source field is absent or does not match are
// passed through unchanged.
func Apply(entries []parser.Entry, rules []Rule) []parser.Entry {
	if len(rules) == 0 {
		return entries
	}
	out := make([]parser.Entry, 0, len(entries))
	for _, e := range entries {
		for _, rule := range rules {
			v, ok := e.Fields[rule.Field]
			if !ok {
				continue
			}
			s, ok := v.(string)
			if !ok {
				continue
			}
			match := rule.Regexp.FindStringSubmatch(s)
			if match == nil {
				continue
			}
			for i, name := range rule.Regexp.SubexpNames() {
				if name == "" {
					continue
				}
				key := name
				if rule.Prefix != "" {
					key = rule.Prefix + name
				}
				e.Fields[key] = match[i]
			}
		}
		out = append(out, e)
	}
	return out
}

// splitN splits s on sep at most n times, returning up to n parts.
func splitN(s, sep string, n int) []string {
	var parts []string
	for len(parts) < n-1 {
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
