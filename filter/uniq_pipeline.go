package filter

import (
	"fmt"
	"io"
)

// RunUniqPipeline reads from r, deduplicates consecutive records sharing
// the same value for field, and writes results to w in the given format.
// It also prints a summary of processed vs emitted lines to summary.
func RunUniqPipeline(r io.Reader, w io.Writer, summary io.Writer, field string, format OutputFormat) error {
	uniqer, err := NewUniqer(field)
	if err != nil {
		return err
	}

	writer := NewWriter(w, format)
	scanner := NewScanner(r)

	var total, emitted, invalid int

	for scanner.Scan() {
		line := scanner.Text()
		rec, err := ParseRecord(line)
		if err != nil {
			invalid++
			continue
		}
		total++
		out, keep := uniqer.Apply(rec)
		if !keep {
			continue
		}
		if err := writer.Write(out); err != nil {
			return fmt.Errorf("uniq pipeline: write error: %w", err)
		}
		emitted++
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("uniq pipeline: scan error: %w", err)
	}

	if summary != nil {
		fmt.Fprintf(summary, "uniq: total=%d emitted=%d suppressed=%d invalid=%d\n",
			total, emitted, total-emitted, invalid)
	}

	return nil
}
