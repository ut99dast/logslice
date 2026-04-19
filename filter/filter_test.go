package filter

import (
	"testing"
	"time"
)

func TestParseRecord_ValidJSON(t *testing.T) {
	line := `{"level":"info","msg":"started","time":"2024-03-01T10:00:00Z"}`
	r, err := ParseRecord(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Fields["level"] != "info" {
		t.Errorf("expected level=info, got %v", r.Fields["level"])
	}
}

func TestParseRecord_InvalidJSON(t *testing.T) {
	_, err := ParseRecord("not json")
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestTimeField(t *testing.T) {
	line := `{"time":"2024-03-01T10:00:00Z"}`
	r, _ := ParseRecord(line)
	got, ok := r.TimeField("time")
	if !ok {
		t.Fatal("expected time field to parse")
	}
	want := time.Date(2024, 3, 1, 10, 0, 0, 0, time.UTC)
	if !got.Equal(want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestFieldEquals(t *testing.T) {
	line := `{"level":"ERROR"}`
	r, _ := ParseRecord(line)
	if !r.FieldEquals("level", "error") {
		t.Error("expected case-insensitive match")
	}
	if r.FieldEquals("level", "info") {
		t.Error("expected no match")
	}
}

func TestTimeRange_Match(t *testing.T) {
	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)
	tr := TimeRange{Start: &start, End: &end}

	mid := time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)
	if !tr.Match(mid) {
		t.Error("expected mid-year to match")
	}
	before := time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC)
	if tr.Match(before) {
		t.Error("expected before-start to not match")
	}
	after := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	if tr.Match(after) {
		t.Error("expected after-end to not match")
	}
}
