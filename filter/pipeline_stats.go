package filter

import "io"

// PipelineWithStats wraps a Pipeline and collects Stats during Run.
type PipelineWithStats struct {
	pipeline *Pipeline
	filter   Filter
	writer   *Writer
	Stats    Stats
}

// Filter is the interface satisfied by MultiFilter, TimeRange, FieldFilter etc.
type Filter interface {
	Match(record map[string]interface{}) bool
}

// NewPipelineWithStats creates a PipelineWithStats.
func NewPipelineWithStats(r io.Reader, f Filter, w *Writer) *PipelineWithStats {
	return &PipelineWithStats{
		pipeline: NewPipeline(r, nil, w),
		filter:   f,
		writer:   w,
	}
}

// Run processes all lines, writes matches, and populates Stats.
func (p *PipelineWithStats) Run(r io.Reader) error {
	scanner := NewScanner(r, p.filter)
	for {
		rec, err := scanner.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			p.Stats.Add(false, true)
			continue
		}
		matched := p.filter.Match(rec)
		p.Stats.Add(matched, false)
		if matched {
			if werr := p.writer.Write(rec); werr != nil {
				return werr
			}
		}
	}
	return p.writer.Flush()
}

// Reset clears the collected Stats, allowing the PipelineWithStats to be reused.
func (p *PipelineWithStats) Reset() {
	p.Stats = Stats{}
}
