// Package filter provides primitives for parsing, filtering, transforming,
// aggregating, and annotating structured (JSON) log records.
//
// # Core types
//
// ParseRecord parses a single JSON log line into a map[string]interface{}.
//
// # Filtering
//
// TimeRange and FieldFilter let callers select records by timestamp range or
// by arbitrary field equality / regex. MultiFilter composes multiple filters
// with AND semantics.
//
// # Transformation
//
// Transformer supports RenameField, DropField, AddField and more.
// Truncator, Flattener, Masker, Caster, and Highlighter each address a
// specific field-level transformation need.
//
// # Aggregation & deduplication
//
// AggregateResult and RunAggregate count records grouped by a field value.
// Deduplicator removes repeated records based on a key field or whole-record
// identity.
//
// # Sampling, limiting, sorting
//
// Sampler keeps every N-th record. Limiter stops after N records. Sorter
// orders records by a numeric or string field. Tailer / Header return the
// last / first N records respectively.
//
// # Annotation
//
// Highlighter annotates matching substrings inside a chosen field using one
// of three modes: bracket ([[match]]), upper (MATCH), or mark (>>>value).
//
// # Output
//
// Writer supports JSON, pretty-printed JSON, and CSV output formats.
package filter
