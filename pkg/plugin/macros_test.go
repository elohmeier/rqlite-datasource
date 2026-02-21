package plugin

import (
	"testing"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

func TestApplyMacros_TimeFilter(t *testing.T) {
	tr := backend.TimeRange{
		From: time.Unix(1000, 0),
		To:   time.Unix(2000, 0),
	}

	sql := "SELECT * FROM t WHERE $__timeFilter(ts)"
	result := ApplyMacros(sql, tr, 60000)
	expected := "SELECT * FROM t WHERE ts >= 1000 AND ts <= 2000"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestApplyMacros_UnixEpochFilter(t *testing.T) {
	tr := backend.TimeRange{
		From: time.Unix(500, 0),
		To:   time.Unix(600, 0),
	}

	sql := "SELECT * FROM t WHERE $__unixEpochFilter(created_at)"
	result := ApplyMacros(sql, tr, 60000)
	expected := "SELECT * FROM t WHERE created_at >= 500 AND created_at <= 600"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestApplyMacros_TimeFromTo(t *testing.T) {
	tr := backend.TimeRange{
		From: time.Unix(1000, 0),
		To:   time.Unix(2000, 0),
	}

	sql := "SELECT * FROM t WHERE ts BETWEEN $__timeFrom AND $__timeTo"
	result := ApplyMacros(sql, tr, 60000)
	expected := "SELECT * FROM t WHERE ts BETWEEN 1000 AND 2000"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestApplyMacros_TimeGroup(t *testing.T) {
	tr := backend.TimeRange{
		From: time.Unix(1000, 0),
		To:   time.Unix(2000, 0),
	}

	tests := []struct {
		name     string
		sql      string
		expected string
	}{
		{
			name:     "5 minutes",
			sql:      "SELECT $__timeGroup(ts, 5m), COUNT(*) FROM t GROUP BY 1",
			expected: "SELECT (CAST(ts / 300 AS INTEGER) * 300), COUNT(*) FROM t GROUP BY 1",
		},
		{
			name:     "1 hour",
			sql:      "SELECT $__timeGroup(ts, 1h), COUNT(*) FROM t GROUP BY 1",
			expected: "SELECT (CAST(ts / 3600 AS INTEGER) * 3600), COUNT(*) FROM t GROUP BY 1",
		},
		{
			name:     "1 day",
			sql:      "SELECT $__timeGroup(ts, 1d), COUNT(*) FROM t GROUP BY 1",
			expected: "SELECT (CAST(ts / 86400 AS INTEGER) * 86400), COUNT(*) FROM t GROUP BY 1",
		},
		{
			name:     "plain seconds",
			sql:      "SELECT $__timeGroup(ts, 60), COUNT(*) FROM t GROUP BY 1",
			expected: "SELECT (CAST(ts / 60 AS INTEGER) * 60), COUNT(*) FROM t GROUP BY 1",
		},
		{
			name:     "$__interval",
			sql:      "SELECT $__timeGroup(ts, $__interval), COUNT(*) FROM t GROUP BY 1",
			expected: "SELECT (CAST(ts / 60 AS INTEGER) * 60), COUNT(*) FROM t GROUP BY 1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ApplyMacros(tt.sql, tr, 60000)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestApplyMacros_MultipleMacros(t *testing.T) {
	tr := backend.TimeRange{
		From: time.Unix(1000, 0),
		To:   time.Unix(2000, 0),
	}

	sql := "SELECT $__timeGroup(ts, 5m) as time, value FROM t WHERE $__timeFilter(ts) GROUP BY 1"
	result := ApplyMacros(sql, tr, 60000)
	expected := "SELECT (CAST(ts / 300 AS INTEGER) * 300) as time, value FROM t WHERE ts >= 1000 AND ts <= 2000 GROUP BY 1"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestApplyMacros_NoMacros(t *testing.T) {
	tr := backend.TimeRange{
		From: time.Unix(1000, 0),
		To:   time.Unix(2000, 0),
	}

	sql := "SELECT * FROM users"
	result := ApplyMacros(sql, tr, 60000)
	if result != sql {
		t.Errorf("expected %q, got %q", sql, result)
	}
}

func TestParseInterval(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"1s", 1},
		{"5m", 300},
		{"1h", 3600},
		{"1d", 86400},
		{"60", 60},
		{"$__interval", 60},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := parseInterval(tt.input, 60000)
			if result != tt.expected {
				t.Errorf("parseInterval(%q): expected %d, got %d", tt.input, tt.expected, result)
			}
		})
	}
}
