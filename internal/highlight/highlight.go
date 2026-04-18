// Package highlight provides ANSI color highlighting for log output.
package highlight

import (
	"fmt"
	"strings"
)

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Cyan   = "\033[36m"
	Bold   = "\033[1m"
)

// LevelColor returns an ANSI color code for a known log level string.
func LevelColor(level string) string {
	switch strings.ToLower(level) {
	case "error", "fatal", "critical":
		return Red
	case "warn", "warning":
		return Yellow
	case "info":
		return Green
	case "debug", "trace":
		return Cyan
	default:
		return Reset
	}
}

// Colorize wraps text with the given ANSI color and resets afterward.
func Colorize(color, text string) string {
	return fmt.Sprintf("%s%s%s", color, text, Reset)
}

// HighlightField returns a formatted key=value string with the value colorized.
func HighlightField(key, value, color string) string {
	return fmt.Sprintf("%s%s%s=%s%s%s", Bold, key, Reset, color, value, Reset)
}

// Strip removes all ANSI escape codes from s.
func Strip(s string) string {
	var b strings.Builder
	inEsc := false
	for _, r := range s {
		if r == '\033' {
			inEsc = true
			continue
		}
		if inEsc {
			if r == 'm' {
				inEsc = false
			}
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}
