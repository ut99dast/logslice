package filter

import (
	"testing"
)

func baseRecord() map[string]interface{} {
	return map[string]interface{}{
		"level": "info",
		"msg":   "hello",
		"ts":    "2024-01-01T00:00:00Z",
	}
}

func TestRenameField(t *testing.T) {
	fn := RenameField("msg", "message")
	out, err := fn(baseRecord())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["msg"]; ok {
		t.Error("old field 'msg' should be removed")
	}
	if out["message"] != "hello" {
		t.Errorf("expected 'hello', got %v", out["message"])
	}
}

func TestRenameField_Missing(t *testing.T) {
	fn := RenameField("nonexistent", "new")
	out, err := fn(baseRecord())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["new"]; ok {
		t.Error("field 'new' should not be added when source is missing")
	}
}

func TestDropField(t *testing.T) {
	fn := DropField("level")
	out, err := fn(baseRecord())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["level"]; ok {
		t.Error("field 'level' should be dropped")
	}
}

func TestAddField(t *testing.T) {
	fn := AddField("env", "production")
	out, err := fn(baseRecord())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["env"] != "production" {
		t.Errorf("expected 'production', got %v", out["env"])
	}
}

func TestRequireField_Present(t *testing.T) {
	fn := RequireField("level")
	_, err := fn(baseRecord())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRequireField_Missing(t *testing.T) {
	fn := RequireField("missing_field")
	_, err := fn(baseRecord())
	if err == nil {
		t.Error("expected error for missing required field")
	}
}

func TestTransformer_Apply(t *testing.T) {
	tr := NewTransformer(
		RenameField("msg", "message"),
		DropField("ts"),
		AddField("version", 2),
	)
	out, err := tr.Apply(baseRecord())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["message"] != "hello" {
		t.Errorf("rename failed: %v", out)
	}
	if _, ok := out["ts"]; ok {
		t.Error("ts should be dropped")
	}
	if out["version"] != 2 {
		t.Errorf("version not set: %v", out)
	}
}

func TestTransformer_Empty(t *testing.T) {
	tr := NewTransformer()
	out, err := tr.Apply(baseRecord())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["level"] != "info" {
		t.Error("record should be unchanged")
	}
}
