// Package sortby provides stable sorting of structured log entries
// by a named field value, supporting both ascending and descending order.
// Numeric fields are compared numerically; string fields lexicographically.
// Entries that are missing the target field are placed at the end of the output.
package sortby
