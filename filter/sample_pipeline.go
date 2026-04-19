package filter

import (
	"bufio"
	"io"
)

// SampleConfig holds options for RunSample.
type SampleConfig struct {
	// Rate keeps 1 out of every Rate records (1 = keep all).
	Rate int
	// Filter is an optional multi-filter expression; empty string means no filtering.
	FilterExpr string
}

// RunSample reads newline-delimited JSON from r, applies optional filtering,
// samples at the configured rate, and writes surviving records to w.
func RunSample(r io.Reader, w io.Writer, cfg SampleConfig) error {
	sampler, err := NewSampler(cfg.Rate)
	if err != nil {
		return err
	}

	var mf *MultiFilter
	if cfg.FilterExpr != "" {
		mf, err = NewMultiFilter(cfg.FilterExpr)
		if err != nil {
			return err
		}
	}

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		rec, err := ParseRecord(line)
		if err != nil {
			continue
		}
		if mf != nil && !mf.Match(rec) {
			continue
		}
		if !sampler.Keep() {
			continue
		}
		if _, err := io.WriteString(w, line+"\n"); err != nil {
			return err
		}
	}
	return scanner.Err()
}
