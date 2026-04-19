package filter

import (
	"testing"
)

func TestDeduplicator_NoDuplicates(t *testing.T) {
	d := NewDeduplicator("id")
	r1 := map[string]interface{}{"id": "1", "msg": "hello"}
	r2 := map[string]interface{}{"id": "2", "msg": "world"}
	if d.IsDuplicate(r1) {
		t.Error("expected r1 not to be duplicate")
	}
	if d.IsDuplicate(r2) {
		t.Error("expected r2 not to be duplicate")
	}
	if d.Dropped != 0 {
		t.Errorf("expected 0 dropped, got %d", d.Dropped)
	}
}

func TestDeduplicator_WithDuplicates(t *testing.T) {
	d := NewDeduplicator("id")
	r := map[string]interface{}{"id": "abc", "msg": "test"}
	if d.IsDuplicate(r) {
		t.Error("first occurrence should not be duplicate")
	}
	if !d.IsDuplicate(r) {
		t.Error("second occurrence should be duplicate")
	}
	if d.Dropped != 1 {
		t.Errorf("expected 1 dropped, got %d", d.Dropped)
	}
	if d.Total != 2 {
		t.Errorf("expected total 2, got %d", d.Total)
	}
}

func TestDeduplicator_WholeRecord(t *testing.T) {
	d := NewDeduplicator()
	r1 := map[string]interface{}{"msg": "hello"}
	r2 := map[string]interface{}{"msg": "hello"}
	d.IsDuplicate(r1)
	if !d.IsDuplicate(r2) {
		t.Error("identical records should be deduplicated")
	}
}

func TestDeduplicator_MissingField(t *testing.T) {
	d := NewDeduplicator("id")
	r1 := map[string]interface{}{"msg": "no id"}
	r2 := map[string]interface{}{"msg": "also no id"}
	d.IsDuplicate(r1)
	if !d.IsDuplicate(r2) {
		t.Error("records with missing key field should collide")
	}
}

func TestDeduplicator_Reset(t *testing.T) {
	d := NewDeduplicator("id")
	r := map[string]interface{}{"id": "x"}
	d.IsDuplicate(r)
	d.Reset()
	if d.IsDuplicate(r) {
		t.Error("after reset, record should not be duplicate")
	}
}
