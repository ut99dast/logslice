package filter

import (
	"bufio"
	"encoding/json"
	"io"
)

// FlattenPipelineOptions configures the flatten pipeline.
type FlattenPipelineOptions struct {
	MaxDepth int
	Format   OutputFormat
}

// RunFlattenPipeline reads JSON log lines from r, flattens each record,
// and writes results to w using the specified output format.
func RunFlattenPipeline(r io.Reader, w io.Writer, opts FlattenPipelineOptions) error {
	fl, err := NewFlattener(opts.MaxDepth)
	if err != nil {
		return err
	}

	writer, err := NewWriter(w, opts.Format, nil)
	if err != nil {
		return err
	}
	defer writer.Flush()

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		var rec map[string]any
		if err := json.Unmarshal([]byte(line), &rec); err != nil {
			continue
		}
		flat := fl.Flatten(rec)
		if err := writer.Write(flat); err != nil {
			return err
		}
	}
	return scanner.Err()
}
