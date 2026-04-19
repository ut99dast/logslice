package filter

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestNewTailer_Invalid(t *testing.T) {
	_, err := NewTailer(0)
	if err == nil {
		t.Fatal("expected error for n=0")
	}
}

func TestTailer_FewerThanN(t *testing.T) {
	tailer, _ := NewTailer(5)
	for i := 0; i < 3; i++ {
		tailer.Add(map[string]interface{}{"i": i})
	}
	records := tailer.Records()
	if len(records) != 3 {
		t.Fatalf("expected 3 records, got %d", len(records))
	}
}

func TestTailer_ExactlyN(t *testing.T) {
	tailer, _ := NewTailer(3)
	for i := 0; i < 3; i++ {
		tailer.Add(map[string]interface{}{"i": float64(i)})
	}
	records := tailer.Records()
	if len(records) != 3 {
		t.Fatalf("expected 3, got %d", len(records))
	}
	for idx, rec := range records {
		if rec["i"] != float64(idx) {
			t.Errorf("record %d: expected i=%d, got %v", idx, idx, rec["i"])
		}
	}
}

func TestTailer_MoreThanN(t *testing.T) {
	tailer, _ := NewTailer(3)
	for i := 0; i < 7; i++ {
		tailer.Add(map[string]interface{}{"i": float64(i)})
	}
	records := tailer.Records()
	if len(records) != 3 {
		t.Fatalf("expected 3, got %d", len(records))
	}
	expected := []float64{4, 5, 6}
	for idx, rec := range records {
		if rec["i"] != expected[idx] {
			t.Errorf("pos %d: expected %v got %v", idx, expected[idx], rec["i"])
		}
	}
}

func TestTailer_Reset(t *testing.T) {
	tailer, _ := NewTailer(3)
	tailer.Add(map[string]interface{}{"x": 1})
	tailer.Reset()
	if len(tailer.Records()) != 0 {
		t.Fatal("expected empty after reset")
	}
}

func TestRunTail(t *testing.T) {
	lines := []string{
		`{"msg":"a"}`,
		`not json`,
		`{"msg":"b"}`,
		`{"msg":"c"}`,
		`{"msg":"d"}`,
	}
	input := strings.NewReader(strings.Join(lines, "\n"))
	var out bytes.Buffer
	if err := RunTail(input, &out, 2); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var records []map[string]interface{}
	dec := json.NewDecoder(&out)
	for dec.More() {
		var rec map[string]interface{}
		if err := dec.Decode(&rec); err != nil {
			t.Fatal(err)
		}
		records = append(records, rec)
	}
	if len(records) != 2 {
		t.Fatalf("expected 2 records, got %d", len(records))
	}
	if records[0]["msg"] != "c" || records[1]["msg"] != "d" {
		t.Errorf("unexpected records: %v", records)
	}
}
