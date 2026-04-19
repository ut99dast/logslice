// Package filter provides primitives for parsing, filtering, transforming,
// and outputting structured (JSON) log records.
//
// # Parsing
//
// ParseRecord decodes a single JSON log line into a map.
//
// # Filtering
//
// Filters implement a common interface and can be combined with MultiFilter:
//   - TimeRange – matches records whose timestamp falls within a time window.
//   - FieldFilter – matches records where a named field equals a given value.
//   - MultiFilter – composes multiple filters with AND semantics.
//
// # Transforming
//
// Transformer applies an ordered chain of TransformFuncs to each record:
//   - RenameField – renames a field key.
//   - DropField   – removes a field from the record.
//   - AddField    – sets a field to a static value.
//   - RequireField – returns an error when a field is absent.
//
// # Scanning
//
// Scanner reads an io.Reader line by line, applies a filter, and emits
// matching records to a channel consumed by a Pipeline.
//
// # Output
//
// Writer serialises records to JSON, pretty-printed JSON, or CSV.
// Use ParseOutputFormat to resolve a format name from a CLI flag.
//
// # Statistics
//
// Stats tracks total lines read, valid JSON lines, and matched lines,
// and can print a summary to any io.Writer.
package filter
