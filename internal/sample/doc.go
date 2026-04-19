// Package sample provides utilities for sampling structured log entries
// by head count, tail count, stride (every Nth), or random rate.
//
// # Sampling Strategies
//
// Head: retain the first N entries.
//
// Tail: retain the last N entries.
//
// Stride: retain every Nth entry (e.g., stride=3 keeps entries 0, 3, 6, …).
//
// Random: retain each entry independently with probability p ∈ [0.0, 1.0].
package sample
