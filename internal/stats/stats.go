package stats

import (
	"fmt"
	"io"
	"sort"

	"github.com/user/logslice/internal/parser"
)

// Summary holds aggregated statistics for a set of log entries.
type Summary struct {
	Total      int
	FieldCounts map[string]map[string]int
}

// Compute calculates statistics over the given entries for the specified fields.
func Compute(entries []parser.Entry, fields []string) *Summary {
	s := &Summary{
		Total:      len(entries),
		FieldCounts: make(map[string]map[string]int),
	}
	for _, f := range fields {
		s.FieldCounts[f] = make(map[string]int)
	}
	for _, e := range entries {
		for _, f := range fields {
			if val, ok := e.Fields[f]; ok {
				s.FieldCounts[f][fmt.Sprintf("%v", val)]++
			}
		}
	}
	return s
}

// Print writes a human-readable summary to w.
func (s *Summary) Print(w io.Writer) {
	fmt.Fprintf(w, "Total entries: %d\n", s.Total)
	for field, counts := range s.FieldCounts {
		fmt.Fprintf(w, "\nField: %s\n", field)
		keys := make([]string, 0, len(counts))
		for k := range counts {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			fmt.Fprintf(w, "  %-30s %d\n", k, counts[k])
		}
	}
}
