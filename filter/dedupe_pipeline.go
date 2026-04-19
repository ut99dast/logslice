package filter

import (
	"bufio"
	"io"
)

// DedupeOptions configures deduplication pipeline behaviour.
type DedupeOptions struct {
	// Fields to use as deduplication key. Empty means whole record.
	Fields []string
}

// RunDedupe reads JSON log lines from r, deduplicates them, and writes
// unique records to w using the provided Writer. Stats are returned.
func RunDedupe(r io.Reader, w *Writer, opts DedupeOptions) (*Stats, error) {
	deduper := NewDeduplicator(opts.Fields...)
	stats := &Stats{}

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		record, err := ParseRecord(line)
		if err != nil {
			stats.Invalid++
			continue
		}
		stats.Valid++
		if deduper.IsDuplicate(record) {
			continue
		}
		stats.Matched++
		if err := w.Write(record); err != nil {
			return stats, err
		}
	}
	if err := scanner.Err(); err != nil {
		return stats, err
	}
	w.Flush()
	return stats, nil
}
