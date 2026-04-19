// Package flatten provides utilities for flattening nested JSON log entries.
package flatten

import (
	"fmt"
	"strings"

	"github.com/logslice/logslice/internal/parser"
)

// Apply flattens nested map fields in each log entry using dot notation.
// For example, {"a": {"b": 1}} becomes {"a.b": 1}.
// If prefix is non-empty, all keys are prefixed with it.
func Apply(entries []parser.Entry, separator string) []parser.Entry {
	if separator == "" {
		separator = "."
	}
	result := make([]parser.Entry, 0, len(entries))
	for _, e := range entries {
		flat := make(map[string]interface{})
		flattenMap("", e.Fields, flat, separator)
		result = append(result, parser.Entry{
			Raw:    e.Raw,
			Fields: flat,
		})
	}
	return result
}

func flattenMap(prefix string, src map[string]interface{}, dst map[string]interface{}, sep string) {
	for k, v := range src {
		key := k
		if prefix != "" {
			key = strings.Join([]string{prefix, k}, sep)
		}
		switch val := v.(type) {
		case map[string]interface{}:
			flattenMap(key, val, dst, sep)
		default:
			dst[key] = fmt.Sprintf("%v", val)
		}
	}
}
