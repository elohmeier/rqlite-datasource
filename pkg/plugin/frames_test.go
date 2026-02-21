package plugin

import (
	"testing"
	"time"
)

func TestResultToFrame_BasicTypes(t *testing.T) {
	result := &RqliteResult{
		Columns: []string{"id", "name", "score"},
		Types:   []string{"integer", "text", "real"},
		Values: [][]interface{}{
			{float64(1), "alice", float64(9.5)},
			{float64(2), "bob", float64(8.3)},
		},
	}

	frame, err := ResultToFrame(result, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(frame.Fields) != 3 {
		t.Fatalf("expected 3 fields, got %d", len(frame.Fields))
	}

	if frame.Fields[0].Name != "id" {
		t.Errorf("expected field name 'id', got %q", frame.Fields[0].Name)
	}
	if frame.Fields[1].Name != "name" {
		t.Errorf("expected field name 'name', got %q", frame.Fields[1].Name)
	}
	if frame.Fields[2].Name != "score" {
		t.Errorf("expected field name 'score', got %q", frame.Fields[2].Name)
	}

	if frame.Fields[0].Len() != 2 {
		t.Errorf("expected 2 rows, got %d", frame.Fields[0].Len())
	}

	// Check values
	v := frame.Fields[0].At(0).(*int64)
	if *v != 1 {
		t.Errorf("expected id=1, got %d", *v)
	}
	name := frame.Fields[1].At(0).(*string)
	if *name != "alice" {
		t.Errorf("expected name='alice', got %q", *name)
	}
	score := frame.Fields[2].At(0).(*float64)
	if *score != 9.5 {
		t.Errorf("expected score=9.5, got %f", *score)
	}
}

func TestResultToFrame_NullValues(t *testing.T) {
	result := &RqliteResult{
		Columns: []string{"id", "name"},
		Types:   []string{"integer", "text"},
		Values: [][]interface{}{
			{float64(1), nil},
			{nil, "bob"},
		},
	}

	frame, err := ResultToFrame(result, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// First row: id=1, name=nil
	v := frame.Fields[0].At(0).(*int64)
	if *v != 1 {
		t.Errorf("expected id=1, got %d", *v)
	}
	name := frame.Fields[1].At(0)
	if name.(*string) != nil {
		t.Errorf("expected nil name, got %v", name)
	}

	// Second row: id=nil, name="bob"
	id2 := frame.Fields[0].At(1)
	if id2.(*int64) != nil {
		t.Errorf("expected nil id, got %v", id2)
	}
}

func TestResultToFrame_TimeColumns_Epoch(t *testing.T) {
	result := &RqliteResult{
		Columns: []string{"time", "value"},
		Types:   []string{"integer", "real"},
		Values: [][]interface{}{
			{float64(1700000000), float64(42.0)},
		},
	}

	frame, err := ResultToFrame(result, []string{"time"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// time field should be parsed as time.Time
	tv := frame.Fields[0].At(0).(*time.Time)
	if tv == nil {
		t.Fatal("expected non-nil time")
	}
	expected := time.Unix(1700000000, 0).UTC()
	if !tv.Equal(expected) {
		t.Errorf("expected time %v, got %v", expected, *tv)
	}
}

func TestResultToFrame_TimeColumns_RFC3339(t *testing.T) {
	result := &RqliteResult{
		Columns: []string{"ts", "value"},
		Types:   []string{"text", "real"},
		Values: [][]interface{}{
			{"2023-11-14T12:00:00Z", float64(42.0)},
		},
	}

	frame, err := ResultToFrame(result, []string{"ts"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tv := frame.Fields[0].At(0).(*time.Time)
	if tv == nil {
		t.Fatal("expected non-nil time")
	}
	expected, _ := time.Parse(time.RFC3339, "2023-11-14T12:00:00Z")
	if !tv.Equal(expected) {
		t.Errorf("expected time %v, got %v", expected, *tv)
	}
}

func TestResultToFrame_TimeColumns_NullTime(t *testing.T) {
	result := &RqliteResult{
		Columns: []string{"time", "value"},
		Types:   []string{"integer", "real"},
		Values: [][]interface{}{
			{nil, float64(42.0)},
		},
	}

	frame, err := ResultToFrame(result, []string{"time"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tv := frame.Fields[0].At(0)
	if tv.(*time.Time) != nil {
		t.Errorf("expected nil time, got %v", tv)
	}
}

func TestResultToFrame_Error(t *testing.T) {
	result := &RqliteResult{
		Error: "some SQL error",
	}

	_, err := ResultToFrame(result, nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestResultToFrame_EmptyResult(t *testing.T) {
	result := &RqliteResult{
		Columns: []string{"id", "name"},
		Types:   []string{"integer", "text"},
		Values:  [][]interface{}{},
	}

	frame, err := ResultToFrame(result, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if frame.Fields[0].Len() != 0 {
		t.Errorf("expected 0 rows, got %d", frame.Fields[0].Len())
	}
}
