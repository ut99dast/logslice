package filter

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// GroupByPipelineConfig holds configuration for the groupby pipeline command.
type GroupByPipelineConfig struct {
	Field string
}

// RunGroupByPipeline wires up the groupby pipeline from the provided reader and
// writer, using the given field name. It prints a summary table of value
// frequencies for the specified field.
func RunGroupByPipeline(r io.Reader, w io.Writer, cfg GroupByPipelineConfig) error {
	field := strings.TrimSpace(cfg.Field)
	if field == "" {
		return fmt.Errorf("groupby pipeline: field must not be empty")
	}
	return RunGroupBy(r, w, field)
}

// RunGroupByFromArgs is a convenience wrapper that reads from stdin and writes
// to stdout, suitable for direct CLI dispatch.
func RunGroupByFromArgs(field string) error {
	if strings.TrimSpace(field) == "" {
		return fmt.Errorf("groupby: --field is required")
	}
	return RunGroupByPipeline(os.Stdin, os.Stdout, GroupByPipelineConfig{
		Field: field,
	})
}
