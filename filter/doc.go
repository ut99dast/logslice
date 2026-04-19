// Package filter provides primitives for parsing and filtering structured
// (JSON) log records.
//
// # Record Parsing
//
// ParseRecord decodes a single JSON log line into a map[string]interface{}.
//
// # Time Filtering
//
// NewTimeRange builds a TimeRange that can test whether a record's timestamp
// field falls within a given [from, to] window. Both boundaries are optional.
//
// # Field Filtering
//
// NewFieldFilter parses expressions of the form "field=value" or
// "field!=value" into a FieldFilter. Multiple filters can be combined with
// NewMultiFilter, which applies AND semantics across all contained filters.
//
// # Typical usage
//
//	tr, _ := filter.NewTimeRange(from, to)
//	mf, _ := filter.NewMultiFilter([]string{"level=error"})
//
//	record, err := filter.ParseRecord(line)
//	if err == nil && tr.Match(record) && mf.Match(record) {
//		// emit record
//	}
package filter
