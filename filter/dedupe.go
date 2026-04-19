package filter

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
)

// Deduplicator tracks seen records and filters duplicates.
type Deduplicator struct {
	keys  []string
	seen  map[string]struct{}
	Total int
	Dropped int
}

// NewDeduplicator creates a Deduplicator that deduplicates by the given fields.
// If no fields are provided, the entire record is used as the key.
func NewDeduplicator(fields ...string) *Deduplicator {
	return &Deduplicator{
		keys: fields,
		seen: make(map[string]struct{}),
	}
}

// IsDuplicate returns true if the record has been seen before.
func (d *Deduplicator) IsDuplicate(record map[string]interface{}) bool {
	d.Total++
	key := d.recordKey(record)
	if _, exists := d.seen[key]; exists {
		d.Dropped++
		return true
	}
	d.seen[key] = struct{}{}
	return false
}

func (d *Deduplicator) recordKey(record map[string]interface{}) string {
	if len(d.keys) == 0 {
		b, _ := json.Marshal(record)
		h := sha256.Sum256(b)
		return hex.EncodeToString(h[:])
	}
	parts := make([]string, 0, len(d.keys))
	for _, k := range d.keys {
		v, ok := record[k]
		if !ok {
			parts = append(parts, fmt.Sprintf("%s=<nil>", k))
		} else {
			parts = append(parts, fmt.Sprintf("%s=%v", k, v))
		}
	}
	return strings.Join(parts, "|")
}

// Reset clears the seen set.
func (d *Deduplicator) Reset() {
	d.seen = make(map[string]struct{})
	d.Total = 0
	d.Dropped = 0
}
