// Package pivot provides functionality to group and count log entries
// by a specified field, producing a frequency table as output entries.
package pivot

import (
	"fmt"
	"sort"

	"github.com/user/logslice/internal/parser"
)

// Result holds a single row of the pivot output.
type Result struct {
	Value string
	Count int
}

// ByField groups entries by the given field and returns a slice of Results
// sorted by count descending. Entries missing the field are grouped under
// the empty string key unless skipMissing is true.
func ByField(entries []parser.Entry, field string, skipMissing bool) []Result {
	counts := make(map[string]int)

	for _, e := range entries {
		v, ok := e.Fields[field]
		if !ok {
			if skipMissing {
				continue
			}
			counts[""]++
			continue
		}
		counts[fmt.Sprintf("%v", v)]++
	}

	results := make([]Result, 0, len(counts))
	for val, cnt := range counts {
		results = append(results, Result{Value: val, Count: cnt})
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].Count != results[j].Count {
			return results[i].Count > results[j].Count
		}
		return results[i].Value < results[j].Value
	})

	return results
}

// ToEntries converts pivot Results into parser.Entry values so they can be
// written by the standard output writer.
func ToEntries(results []Result, field string) []parser.Entry {
	out := make([]parser.Entry, len(results))
	for i, r := range results {
		out[i] = parser.Entry{
			Fields: map[string]interface{}{
				field:   r.Value,
				"count": r.Count,
			},
		}
	}
	return out
}
