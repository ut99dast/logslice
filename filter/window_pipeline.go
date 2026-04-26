package filter

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
)

// RunWindow reads all valid JSON records from r, groups them into windows
// using the given field, size, and overlap, then writes each window as a
// JSON object with "window" (index) and "records" keys to w.
func RunWindow(r io.Reader, w io.Writer, field string, size, overlap int) error {
	windower, err := NewWindower(field, size, overlap)
	if err != nil {
		return err
	}

	var records []map[string]interface{}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		rec, err := ParseRecord(line)
		if err != nil {
			continue
		}
		records = append(records, rec)
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("scan error: %w", err)
	}

	if field != "" {
		SortRecordsByField(records, field)
	}

	groups := windower.Apply(records)
	bw := bufio.NewWriter(w)
	for _, g := range groups {
		out := map[string]interface{}{
			"window":  strconv.Itoa(g.Index),
			"records": g.Records,
		}
		b, err := json.Marshal(out)
		if err != nil {
			return fmt.Errorf("marshal error: %w", err)
		}
		_, _ = bw.Write(b)
		_ = bw.WriteByte('\n')
	}
	return bw.Flush()
}
