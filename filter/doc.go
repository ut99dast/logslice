// Package filter provides primitives for parsing, filtering, transforming,
// aggregating, and deduplicating structured (JSON) log records.
//
// # Parsing
//
// ParseRecord parses a single JSON log line into a map.
//
// # Filtering
//
// TimeRange and FieldFilter implement the Filter interface and can be
// composed with MultiFilter for AND-semantics across multiple conditions.
//
// # Transformation
//
// Transformer supports field-level operations: rename, drop, add, and
// template-based computed fields. TransformPipeline chains transformers.
//
// # Aggregation
//
// RunAggregate groups records by a key field and computes counts.
//
// # Deduplication
//
// Deduplicator tracks seen records by key fields (or whole-record hash)
// and drops duplicates. RunDedupe integrates deduplication into a
// streaming pipeline with Stats reporting.
//
// # Output
//
// Writer supports JSON, pretty-printed JSON, and CSV output formats.
//
// # Statistics
//
// Stats tracks valid, invalid, and matched record counts across a run.
package filter
