// Package convert transforms log entries between formats (JSON, CSV, logfmt).
package convert

import (
	"encoding/csv"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/user/logslice/internal/parser"
)

// Format represents an output format.
type Format string

const (
	FormatJSON   Format = "json"
	FormatCSV    Format = "csv"
	FormatLogfmt Format = "logfmt"
)

// ToLogfmt writes entries as logfmt lines to w.
func ToLogfmt(entries []parser.Entry, w io.Writer) error {
	for _, e := range entries {
		parts := make([]string, 0, len(e.Fields))
		keys := sortedKeys(e.Fields)
		for _, k := range keys {
			v := fmt.Sprintf("%v", e.Fields[k])
			if strings.ContainsAny(v, " \t\"=") {
				v = fmt.Sprintf("%q", v)
			}
			parts = append(parts, k+"="+v)
		}
		if _, err := fmt.Fprintln(w, strings.Join(parts, " ")); err != nil {
			return err
		}
	}
	return nil
}

// ToCSV writes entries as CSV rows to w. Fields are sorted alphabetically.
func ToCSV(entries []parser.Entry, w io.Writer) error {
	if len(entries) == 0 {
		return nil
	}
	cw := csv.NewWriter(w)
	headers := sortedKeys(entries[0].Fields)
	if err := cw.Write(headers); err != nil {
		return err
	}
	for _, e := range entries {
		row := make([]string, len(headers))
		for i, h := range headers {
			row[i] = fmt.Sprintf("%v", e.Fields[h])
		}
		if err := cw.Write(row); err != nil {
			return err
		}
	}
	cw.Flush()
	return cw.Error()
}

func sortedKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
