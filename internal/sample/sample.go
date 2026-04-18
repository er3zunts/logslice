// Package sample provides log entry sampling by rate or count.
package sample

import (
	"math/rand"

	"github.com/yourorg/logslice/internal/parser"
)

// Options controls sampling behaviour.
type Options struct {
	// Rate is a value between 0.0 and 1.0; entries are kept with this probability.
	Rate float64
	// Every keeps every Nth entry (1 = all, 2 = every other, etc.).
	Every int
	// Head keeps only the first N entries (0 = disabled).
	Head int
	// Tail keeps only the last N entries (0 = disabled).
	Tail int
}

// Apply returns a sampled subset of entries according to opts.
func Apply(entries []parser.Entry, opts Options) []parser.Entry {
	if opts.Head > 0 {
		if opts.Head >= len(entries) {
			return entries
		}
		return entries[:opts.Head]
	}
	if opts.Tail > 0 {
		if opts.Tail >= len(entries) {
			return entries
		}
		return entries[len(entries)-opts.Tail:]
	}
	if opts.Every > 1 {
		return everyN(entries, opts.Every)
	}
	if opts.Rate > 0 && opts.Rate < 1.0 {
		return byRate(entries, opts.Rate)
	}
	return entries
}

func everyN(entries []parser.Entry, n int) []parser.Entry {
	out := make([]parser.Entry, 0, len(entries)/n+1)
	for i, e := range entries {
		if i%n == 0 {
			out = append(out, e)
		}
	}
	return out
}

func byRate(entries []parser.Entry, rate float64) []parser.Entry {
	out := make([]parser.Entry, 0, int(float64(len(entries))*rate)+1)
	for _, e := range entries {
		if rand.Float64() < rate {
			out = append(out, e)
		}
	}
	return out
}
