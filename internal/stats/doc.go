// Package stats provides aggregation and summarisation utilities for
// structured log entries parsed by the parser package.
//
// Use Compute to build a Summary over a slice of entries, optionally
// counting distinct values for one or more named fields. The resulting
// Summary can be printed in a human-readable table via Print.
package stats
