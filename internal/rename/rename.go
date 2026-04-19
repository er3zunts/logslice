// Package rename provides field renaming for log entries.
package rename

import "github.com/logslice/logslice/internal/parser"

// Rule maps an old field name to a new field name.
type Rule struct {
	From string
	To   string
}

// Apply renames fields in each entry according to the provided rules.
// If a source field does not exist the entry is left unchanged.
// If the destination field already exists it will be overwritten.
func Apply(entries []parser.Entry, rules []Rule) []parser.Entry {
	if len(rules) == 0 {
		return entries
	}
	out := make([]parser.Entry, 0, len(entries))
	for _, e := range entries {
		fields := make(map[string]interface{}, len(e.Fields))
		for k, v := range e.Fields {
			fields[k] = v
		}
		for _, r := range rules {
			if val, ok := fields[r.From]; ok {
				delete(fields, r.From)
				fields[r.To] = val
			}
		}
		out = append(out, parser.Entry{Fields: fields})
	}
	return out
}
