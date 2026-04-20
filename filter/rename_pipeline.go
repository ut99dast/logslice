package filter

import (
	"fmt"
	"io"
	"strings"
)

// RenameConfig holds the configuration for a rename pipeline run.
type RenameConfig struct {
	// Mappings is a list of "old=new" field rename expressions.
	Mappings []string
	// Format is the output format (json, pretty, csv).
	Format string
	// CSVFields lists fields to include when Format is "csv".
	CSVFields []string
}

// RunRenamePipeline wires together parsing, renaming, and writing.
func RunRenamePipeline(in io.Reader, out io.Writer, cfg RenameConfig) error {
	if len(cfg.Mappings) == 0 {
		return fmt.Errorf("rename pipeline: no field mappings provided")
	}

	format, err := ParseOutputFormat(cfg.Format)
	if err != nil {
		return fmt.Errorf("rename pipeline: %w", err)
	}

	renamer, err := NewRenamer(cfg.Mappings)
	if err != nil {
		return fmt.Errorf("rename pipeline: %w", err)
	}

	writer, err := NewWriter(out, format, cfg.CSVFields)
	if err != nil {
		return fmt.Errorf("rename pipeline: %w", err)
	}
	defer writer.Flush()

	_ = strings.Join // imported for potential future use

	return RunRename(in, writer, renamer)
}
