package typecast

import (
	"fmt"
	"strconv"

	"github.com/logslice/logslice/internal/parser"
)

// Rule defines a field and the target type to cast it to.
type Rule struct {
	Field  string
	Target string // "string", "int", "float", "bool"
}

// Apply casts fields in each entry according to the given rules.
func Apply(entries []parser.Entry, rules []Rule) []parser.Entry {
	result := make([]parser.Entry, 0, len(entries))
	for _, e := range entries {
		fields := make(map[string]interface{}, len(e.Fields))
		for k, v := range e.Fields {
			fields[k] = v
		}
		for _, r := range rules {
			val, ok := fields[r.Field]
			if !ok {
				continue
			}
			if casted, err := cast(val, r.Target); err == nil {
				fields[r.Field] = casted
			}
		}
		result = append(result, parser.Entry{Fields: fields})
	}
	return result
}

func cast(val interface{}, target string) (interface{}, error) {
	s := fmt.Sprintf("%v", val)
	switch target {
	case "string":
		return s, nil
	case "int":
		return strconv.ParseInt(s, 10, 64)
	case "float":
		return strconv.ParseFloat(s, 64)
	case "bool":
		return strconv.ParseBool(s)
	}
	return nil, fmt.Errorf("unknown target type: %s", target)
}

// ParseRules parses rules from strings like "field:type".
func ParseRules(specs []string) ([]Rule, error) {
	rules := make([]Rule, 0, len(specs))
	for _, s := range specs {
		for i := len(s) - 1; i >= 0; i-- {
			if s[i] == ':' {
				rules = append(rules, Rule{Field: s[:i], Target: s[i+1:]})
				goto next
			}
		}
		return nil, fmt.Errorf("invalid rule %q: expected field:type", s)
	next:
	}
	return rules, nil
}
