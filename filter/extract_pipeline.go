package filter

import (
	"fmt"
	"io"
)

// RunExtractPipeline reads JSON log lines from r, applies field extraction using
// the comma-separated fields string, and writes results to w.
// It returns the number of records written and any error encountered.
func RunExtractPipeline(r io.Reader, w io.Writer, fields string, format OutputFormat) (int, error) {
	writer, err := NewWriter(w, format)
	if err != nil {
		return 0, fmt.Errorf("extract pipeline: %w", err)
	}

	ext, err := NewExtractor(fields)
	if err != nil {
		return 0, fmt.Errorf("extract pipeline: %w", err)
	}

	scanner := NewScanner(r, nil)
	count := 0
	for scanner.Scan() {
		record := scanner.Record()
		out := ext.Apply(record)
		if len(out) == 0 {
			continue
		}
		if err := writer.Write(out); err != nil {
			return count, fmt.Errorf("extract pipeline: write: %w", err)
		}
		count++
	}
	if err := scanner.Err(); err != nil {
		return count, fmt.Errorf("extract pipeline: scan: %w", err)
	}
	return count, nil
}
