package truncate

import (
	"testing"

	"github.com/user/logslice/internal/parser"
)

func makeEntry(fields map[string]interface{}) parser.Entry {
	return parser.Entry{Fields: fields}
}

func TestApply_NoOp_ZeroMaxLen(t *testing.T) {
	entries := []parser.Entry{makeEntry(map[string]interface{}{"msg": "hello world"})}
	out := Apply(entries, Options{MaxLen: 0})
	if out[0].Fields["msg"] != "hello world" {
		t.Errorf("expected no truncation, got %v", out[0].Fields["msg"])
	}
}

func TestApply_AllFields(t *testing.T) {
	entries := []parser.Entry{makeEntry(map[string]interface{}{"msg": "hello world", "level": "info"})}
	out := Apply(entries, Options{MaxLen: 5, Suffix: "..."})
	if out[0].Fields["msg"] != "he..." {
		t.Errorf("unexpected msg: %v", out[0].Fields["msg"])
	}
	if out[0].Fields["level"] != "info" {
		t.Errorf("unexpected level: %v", out[0].Fields["level"])
	}
}

func TestApply_SpecificField(t *testing.T) {
	entries := []parser.Entry{makeEntry(map[string]interface{}{"msg": "hello world", "detail": "some long detail text"})}
	out := Apply(entries, Options{MaxLen: 8, Fields: []string{"detail"}, Suffix: "..."})
	if out[0].Fields["msg"] != "hello world" {
		t.Errorf("msg should be unchanged, got %v", out[0].Fields["msg"])
	}
	if out[0].Fields["detail"] != "some ..." {
		t.Errorf("unexpected detail: %v", out[0].Fields["detail"])
	}
}

func TestApply_NonStringUnchanged(t *testing.T) {
	entries := []parser.Entry{makeEntry(map[string]interface{}{"count": float64(42)})}
	out := Apply(entries, Options{MaxLen: 2, Suffix: "..."})
	if out[0].Fields["count"] != float64(42) {
		t.Errorf("non-string field should be unchanged, got %v", out[0].Fields["count"])
	}
}

func TestApply_ShortStringUnchanged(t *testing.T) {
	entries := []parser.Entry{makeEntry(map[string]interface{}{"msg": "hi"})}
	out := Apply(entries, Options{MaxLen: 10, Suffix: "..."})
	if out[0].Fields["msg"] != "hi" {
		t.Errorf("short string should be unchanged, got %v", out[0].Fields["msg"])
	}
}
