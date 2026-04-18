package highlight_test

import (
	"strings"
	"testing"

	"github.com/user/logslice/internal/highlight"
)

func TestLevelColor(t *testing.T) {
	cases := []struct {
		level    string
		wantColor string
	}{
		{"error", highlight.Red},
		{"FATAL", highlight.Red},
		{"warn", highlight.Yellow},
		{"WARNING", highlight.Yellow},
		{"info", highlight.Green},
		{"debug", highlight.Cyan},
		{"trace", highlight.Cyan},
		{"unknown", highlight.Reset},
	}
	for _, tc := range cases {
		got := highlight.LevelColor(tc.level)
		if got != tc.wantColor {
			t.Errorf("LevelColor(%q) = %q, want %q", tc.level, got, tc.wantColor)
		}
	}
}

func TestColorize(t *testing.T) {
	result := highlight.Colorize(highlight.Red, "hello")
	if !strings.Contains(result, "hello") {
		t.Error("expected text in colorized output")
	}
	if !strings.HasPrefix(result, highlight.Red) {
		t.Error("expected color prefix")
	}
	if !strings.HasSuffix(result, highlight.Reset) {
		t.Error("expected reset suffix")
	}
}

func TestHighlightField(t *testing.T) {
	out := highlight.HighlightField("level", "error", highlight.Red)
	stripped := highlight.Strip(out)
	if stripped != "level=error" {
		t.Errorf("stripped HighlightField = %q, want %q", stripped, "level=error")
	}
}

func TestStrip(t *testing.T) {
	input := "\033[31mhello\033[0m world"
	got := highlight.Strip(input)
	if got != "hello world" {
		t.Errorf("Strip() = %q, want %q", got, "hello world")
	}
}

func TestStrip_NoEscapes(t *testing.T) {
	input := "plain text"
	if got := highlight.Strip(input); got != input {
		t.Errorf("Strip() modified plain text: %q", got)
	}
}
