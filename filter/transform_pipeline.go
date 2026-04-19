package filter

import "io"

// TransformPipeline combines a Scanner with a Transformer and Writer,
// processing each matched record through a set of transformations before output.
type TransformPipeline struct {
	scanner     *Scanner
	transformer *Transformer
	writer      *Writer
	stats        *Stats
}

// NewTransformPipeline creates a TransformPipeline that reads from r, applies
// the given filter, transforms each matching record, and writes to w.
func NewTransformPipeline(r io.Reader, f Filter, t *Transformer, w *Writer) *TransformPipeline {
	return &TransformPipeline{
		scanner:     NewScanner(r, f),
		transformer: t,
		writer:      w,
		stats:        &Stats{},
	}
}

// Run executes the pipeline: scan, transform, write.
// It returns the first write error encountered, or any scanner error.
func (tp *TransformPipeline) Run() error {
	for tp.scanner.Scan() {
		record := tp.scanner.Record()
		tp.stats.Add(true)

		transformed := tp.transformer.Apply(record)

		if err := tp.writer.Write(transformed); err != nil {
			return err
		}
	}

	if err := tp.scanner.Err(); err != nil {
		return err
	}

	return tp.writer.Flush()
}

// Stats returns the collected statistics for this pipeline run.
func (tp *TransformPipeline) Stats() *Stats {
	return tp.stats
}
