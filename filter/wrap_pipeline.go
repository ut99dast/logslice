package filter

import (
	"fmt"
	"io"
	"strings"
)

// RunWrapPipeline parses wrap arguments of the form "field:prefix:suffix" and
// applies all wrappers sequentially to each record read from r, writing to w.
//
// Each expr must have exactly three colon-separated parts: field, prefix, suffix.
// Prefix and suffix may be empty strings (e.g. "msg:[:]").
func RunWrapPipeline(r io.Reader, w io.Writer, exprs []string, format OutputFormat) error {
	if len(exprs) == 0 {
		return fmt.Errorf("wrap: at least one expression required")
	}

	wrappers := make([]*Wrapper, 0, len(exprs))
	for _, expr := range exprs {
		parts := strings.SplitN(expr, ":", 3)
		if len(parts) != 3 {
			return fmt.Errorf("wrap: invalid expression %q: expected field:prefix:suffix", expr)
		}
		wrapper, err := NewWrapper(parts[0], parts[1], parts[2])
		if err != nil {
			return err
		}
		wrappers = append(wrappers, wrapper)
	}

	writer, err := NewWriter(w, format)
	if err != nil {
		return err
	}

	scanner := NewScanner(r)
	for scanner.Scan() {
		rec, err := scanner.Record()
		if err != nil {
			continue
		}
		for _, wr := range wrappers {
			rec = wr.Apply(rec)
		}
		if err := writer.Write(rec); err != nil {
			return err
		}
	}
	return writer.Flush()
}
