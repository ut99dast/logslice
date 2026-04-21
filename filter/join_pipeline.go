package filter

import (
	"bufio"
	"fmt"
	"io"
)

// RunJoinPipeline reads NDJSON from r, enriches each valid record using
// joiner, and writes output to w. Invalid lines are skipped silently.
// Returns the number of records written and any write error.
func RunJoinPipeline(r io.Reader, joiner *Joiner, w *Writer) (int, error) {
	scanner := bufio.NewScanner(r)
	written := 0
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		rec, err := ParseRecord(line)
		if err != nil {
			continue
		}
		enriched := joiner.Apply(rec)
		if err := w.Write(enriched); err != nil {
			return written, fmt.Errorf("join pipeline: write: %w", err)
		}
		written++
	}
	if err := scanner.Err(); err != nil {
		return written, fmt.Errorf("join pipeline: scan: %w", err)
	}
	return written, nil
}
