package filter

// Sampler selects every Nth matching record.
type Sampler struct {
	rate  int
	count int
}

// NewSampler creates a Sampler that keeps 1 out of every rate records.
// rate must be >= 1.
func NewSampler(rate int) (*Sampler, error) {
	if rate < 1 {
		return nil, fmt.Errorf("sample rate must be >= 1, got %d", rate)
	}
	return &Sampler{rate: rate}, nil
}

// Keep returns true if this record should be kept.
func (s *Sampler) Keep() bool {
	s.count++
	return s.count%s.rate == 1
}

// Reset resets the internal counter.
func (s *Sampler) Reset() {
	s.count = 0
}
