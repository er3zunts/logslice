// Package mask provides field value masking for sensitive log data.
package mask

import (
	"regexp"
	"strings"

	"github.com/user/logslice/internal/parser"
)

// Rule defines how a field should be masked.
type Rule struct {
	Field   string
	Pattern *regexp.Regexp // if nil, mask entire value
	Replace string
}

// Apply masks sensitive fields in each log entry according to the given rules.
func Apply(entries []parser.Entry, rules []Rule) []parser.Entry {
	out := make([]parser.Entry, 0, len(entries))
	for _, e := range entries {
		masked := maskEntry(e, rules)
		out = append(out, masked)
	}
	return out
}

func maskEntry(e parser.Entry, rules []Rule) parser.Entry {
	fields := make(map[string]interface{}, len(e.Fields))
	for k, v := range e.Fields {
		fields[k] = v
	}
	for _, r := range rules {
		val, ok := fields[r.Field]
		if !ok {
			continue
		}
		str, ok := val.(string)
		if !ok {
			continue
		}
		if r.Pattern == nil {
			fields[r.Field] = strings.Repeat("*", len(str))
		} else {
			repl := r.Replace
			if repl == "" {
				repl = "***"
			}
			fields[r.Field] = r.Pattern.ReplaceAllString(str, repl)
		}
	}
	return parser.Entry{Fields: fields}
}

// ParseRules parses mask rules from strings of the form "field" or "field=pattern".
func ParseRules(specs []string) ([]Rule, error) {
	rules := make([]Rule, 0, len(specs))
	for _, s := range specs {
		parts := strings.SplitN(s, "=", 2)
		r := Rule{Field: parts[0]}
		if len(parts) == 2 {
			re, err := regexp.Compile(parts[1])
			if err != nil {
				return nil, err
			}
			r.Pattern = re
		}
		rules = append(rules, r)
	}
	return rules, nil
}
