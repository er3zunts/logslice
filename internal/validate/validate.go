package validate

import (
	"fmt"
	"regexp"

	"github.com/user/logslice/internal/parser"
)

// Rule defines a validation rule for a field.
type Rule struct {
	Field    string
	Required bool
	Pattern  *regexp.Regexp
	TypeName string // "string", "number", "bool"
}

// Result holds the outcome of validating a single entry.
type Result struct {
	Entry  parser.Entry
	Errors []string
}

// Valid returns true if the entry passed all rules.
func (r Result) Valid() bool {
	return len(r.Errors) == 0
}

// Apply validates each entry against the given rules and returns results.
func Apply(entries []parser.Entry, rules []Rule) []Result {
	results := make([]Result, 0, len(entries))
	for _, entry := range entries {
		result := Result{Entry: entry}
		for _, rule := range rules {
			val, exists := entry.Fields[rule.Field]
			if rule.Required && !exists {
				result.Errors = append(result.Errors, fmt.Sprintf("missing required field: %s", rule.Field))
				continue
			}
			if !exists {
				continue
			}
			strVal := fmt.Sprintf("%v", val)
			if rule.Pattern != nil && !rule.Pattern.MatchString(strVal) {
				result.Errors = append(result.Errors, fmt.Sprintf("field %s value %q does not match pattern %s", rule.Field, strVal, rule.Pattern))
			}
			if rule.TypeName != "" {
				if err := checkType(val, rule.TypeName); err != nil {
					result.Errors = append(result.Errors, fmt.Sprintf("field %s: %s", rule.Field, err))
				}
			}
		}
		results = append(results, result)
	}
	return results
}

func checkType(val interface{}, typeName string) error {
	switch typeName {
	case "string":
		if _, ok := val.(string); !ok {
			return fmt.Errorf("expected string, got %T", val)
		}
	case "number":
		switch val.(type) {
		case float64, int, int64:
		default:
			return fmt.Errorf("expected number, got %T", val)
		}
	case "bool":
		if _, ok := val.(bool); !ok {
			return fmt.Errorf("expected bool, got %T", val)
		}
	}
	return nil
}
