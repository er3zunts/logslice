// Package enrich adds derived fields to log entries based on existing field values.
package enrich

import (
	"fmt"
	"strings"

	"github.com/user/logslice/internal/parser"
)

// Rule defines a single enrichment: a new field derived from a template.
type Rule struct {
	// TargetField is the name of the field to add or overwrite.
	TargetField string
	// Template is a string like "${level}-${service}" where ${field} tokens
	// are replaced with the corresponding entry values.
	Template string
}

// ParseRules parses a slice of "field=template" strings into Rules.
func ParseRules(specs []string) ([]Rule, error) {
	rules := make([]Rule, 0, len(specs))
	for _, s := range specs {
		parts := strings.SplitN(s, "=", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return nil, fmt.Errorf("invalid enrich rule %q: expected field=template", s)
		}
		rules = append(rules, Rule{TargetField: parts[0], Template: parts[1]})
	}
	return rules, nil
}

// Apply enriches each entry in entries according to rules, returning a new slice.
func Apply(entries []parser.Entry, rules []Rule) []parser.Entry {
	if len(rules) == 0 {
		return entries
	}
	out := make([]parser.Entry, len(entries))
	for i, e := range entries {
		newFields := make(map[string]interface{}, len(e.Fields)+len(rules))
		for k, v := range e.Fields {
			newFields[k] = v
		}
		for _, r := range rules {
			newFields[r.TargetField] = expand(r.Template, e.Fields)
		}
		out[i] = parser.Entry{Fields: newFields}
	}
	return out
}

// expand replaces ${key} tokens in template with values from fields.
func expand(template string, fields map[string]interface{}) string {
	result := template
	for k, v := range fields {
		token := "${" + k + "}"
		result = strings.ReplaceAll(result, token, fmt.Sprintf("%v", v))
	}
	return result
}
