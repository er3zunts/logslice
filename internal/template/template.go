// Package template provides log entry formatting using Go text/template syntax.
package template

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/nicholasgasior/logslice/internal/parser"
)

// Renderer holds a compiled template for reuse across entries.
type Renderer struct {
	tmpl *template.Template
}

// New compiles the given template string and returns a Renderer.
// Returns an error if the template is invalid.
func New(tmplStr string) (*Renderer, error) {
	funcMap := template.FuncMap{
		"default": func(def, val interface{}) interface{} {
			if val == nil || val == "" {
				return def
			}
			return val
		},
		"upper": func(s string) string {
			return fmt.Sprintf("%s", bytes.ToUpper([]byte(s)))
		},
	}

	t, err := template.New("entry").Funcs(funcMap).Parse(tmplStr)
	if err != nil {
		return nil, fmt.Errorf("template parse error: %w", err)
	}
	return &Renderer{tmpl: t}, nil
}

// Render applies the compiled template to a single log entry and returns
// the resulting string. The entry's Fields map is passed as template data.
func (r *Renderer) Render(entry parser.Entry) (string, error) {
	var buf bytes.Buffer
	if err := r.tmpl.Execute(&buf, entry.Fields); err != nil {
		return "", fmt.Errorf("template render error: %w", err)
	}
	return buf.String(), nil
}

// Apply renders each entry using the Renderer and returns lines of output.
// Entries that fail to render are skipped; errors are collected and returned.
func (r *Renderer) Apply(entries []parser.Entry) ([]string, []error) {
	results := make([]string, 0, len(entries))
	var errs []error
	for _, e := range entries {
		line, err := r.Render(e)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		results = append(results, line)
	}
	return results, errs
}
