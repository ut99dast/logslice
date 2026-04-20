package filter

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
)

// RunMask reads JSON log lines from r, applies the masker to each valid record,
// and writes the resulting JSON lines to w.
// Invalid JSON lines are skipped and counted in the returned error summary.
func RunMask(r io.Reader, w io.Writer, masker *Masker) error {
	scanner := bufio.NewScanner(r)
	var skipped int
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		var record map[string]interface{}
		if err := json.Unmarshal(line, &record); err != nil {
			skipped++
			continue
		}
		masked := masker.Apply(record)
		out, err := json.Marshal(masked)
		if err != nil {
			skipped++
			continue
		}
		if _, err := fmt.Fprintf(w, "%s\n", out); err != nil {
			return err
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	if skipped > 0 {
		return fmt.Errorf("mask: skipped %d invalid line(s)", skipped)
	}
	return nil
}
