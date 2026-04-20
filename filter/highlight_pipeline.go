package filter

import (
	"bufio"
	"encoding/json"
	"io"
)

// RunHighlight reads JSON lines from r, applies the Highlighter to each valid
// record, and writes the (possibly annotated) JSON lines to w.
// Invalid lines are skipped silently.
func RunHighlight(r io.Reader, w io.Writer, h *Highlighter) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Bytes()
		var record map[string]interface{}
		if err := json.Unmarshal(line, &record); err != nil {
			continue
		}
		annotated := h.Apply(record)
		out, err := json.Marshal(annotated)
		if err != nil {
			continue
		}
		if _, err := w.Write(append(out, '\n')); err != nil {
			return err
		}
	}
	return scanner.Err()
}
