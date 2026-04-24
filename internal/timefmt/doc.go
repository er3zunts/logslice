// Package timefmt reformats timestamp fields within log entries.
//
// Rules are expressed as "field:inputLayout:outputLayout" strings,
// where layouts follow Go's time.Parse reference time convention
// (Mon Jan 2 15:04:05 MST 2006).
package timefmt
