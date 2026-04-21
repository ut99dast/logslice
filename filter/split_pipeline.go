package filter

import (
	"fmt"
	"io"
)

// SplitArgs holds the CLI arguments for the split command.
type SplitArgs struct {
	Field     string
	Delimiter string
	OutField  string
	Format    string
}

// RunSplitPipeline wires together input scanning, splitting, and output writing.
func RunSplitPipeline(in io.Reader, out io.Writer, args SplitArgs) error {
	if args.Field == "" {
		return fmt.Errorf("split: --field is required")
	}
	if args.Delimiter == "" {
		return fmt.Errorf("split: --delimiter is required")
	}

	fmt, err := ParseOutputFormat(args.Format)
	if err != nil {
		return err
	}

	scanner := NewScanner(in)
	writer := NewWriter(out, fmt)

	if err := RunSplit(scanner, writer, args.Field, args.Delimiter, args.OutField); err != nil {
		return err
	}
	return writer.Flush()
}
