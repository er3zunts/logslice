// Package compute provides field computation for log entries,
// allowing new fields to be derived from existing ones.
package compute

import (
	"fmt"
	"strconv"

	"github.com/user/logslice/internal/parser"
)

// Rule defines a computed field: the target field name and an expression.
type Rule struct {
	Field string
	Expr  string // supported: "field1+field2" (concat), "field1-field2" (numeric diff)
}

// ParseRules parses strings like "latency_ms=end-start" into Rule structs.
func ParseRules(specs []string) ([]Rule, error) {
	rules := make([]Rule, 0, len(specs))
	for _, s := range specs {
		field, expr, ok := cut(s, "=")
		if !ok {
			return nil, fmt.Errorf("invalid compute rule %q: expected field=expr", s)
		}
		rules = append(rules, Rule{Field: field, Expr: expr})
	}
	return rules, nil
}

// Apply evaluates each rule against every entry and sets the computed field.
func Apply(entries []parser.Entry, rules []Rule) []parser.Entry {
	if len(rules) == 0 {
		return entries
	}
	out := make([]parser.Entry, len(entries))
	for i, e := range entries {
		fields := copyFields(e.Fields)
		for _, r := range rules {
			if v, err := eval(fields, r.Expr); err == nil {
				fields[r.Field] = v
			}
		}
		out[i] = parser.Entry{Fields: fields}
	}
	return out
}

func eval(fields map[string]interface{}, expr string) (interface{}, error) {
	if l, r, ok := cut(expr, "-"); ok {
		lv, err1 := toFloat(fields, l)
		rv, err2 := toFloat(fields, r)
		if err1 != nil || err2 != nil {
			return nil, fmt.Errorf("non-numeric operands")
		}
		return lv - rv, nil
	}
	if l, r, ok := cut(expr, "+"); ok {
		ls := fieldStr(fields, l)
		rs := fieldStr(fields, r)
		return ls + rs, nil
	}
	return nil, fmt.Errorf("unsupported expression: %s", expr)
}

func toFloat(fields map[string]interface{}, key string) (float64, error) {
	v, ok := fields[key]
	if !ok {
		return 0, fmt.Errorf("missing field %s", key)
	}
	switch n := v.(type) {
	case float64:
		return n, nil
	case string:
		return strconv.ParseFloat(n, 64)
	}
	return 0, fmt.Errorf("not numeric")
}

func fieldStr(fields map[string]interface{}, key string) string {
	if v, ok := fields[key]; ok {
		return fmt.Sprintf("%v", v)
	}
	return ""
}

func copyFields(src map[string]interface{}) map[string]interface{} {
	m := make(map[string]interface{}, len(src))
	for k, v := range src {
		m[k] = v
	}
	return m
}

func cut(s, sep string) (string, string, bool) {
	for i := 0; i < len(s); i++ {
		if string(s[i]) == sep {
			return s[:i], s[i+1:], true
		}
	}
	return "", "", false
}
