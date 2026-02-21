package plugin

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/data"
)

// ResultToFrame converts a rqlite result to a Grafana data frame.
func ResultToFrame(result *RqliteResult, timeColumns []string) (*data.Frame, error) {
	if result.Error != "" {
		return nil, fmt.Errorf("rqlite query error: %s", result.Error)
	}

	frame := data.NewFrame("response")

	timeColSet := make(map[string]bool, len(timeColumns))
	for _, tc := range timeColumns {
		timeColSet[strings.ToLower(tc)] = true
	}

	// Build fields based on column types
	fields := make([]*data.Field, len(result.Columns))
	for i, col := range result.Columns {
		colType := ""
		if i < len(result.Types) {
			colType = strings.ToLower(result.Types[i])
		}

		if timeColSet[strings.ToLower(col)] {
			fields[i] = data.NewField(col, nil, make([]*time.Time, 0, len(result.Values)))
		} else {
			fields[i] = newFieldForType(col, colType, len(result.Values))
		}
	}

	// Fill in values
	for _, row := range result.Values {
		for colIdx, field := range fields {
			var val interface{}
			if colIdx < len(row) {
				val = row[colIdx]
			}
			appendValue(field, val, timeColSet[strings.ToLower(result.Columns[colIdx])])
		}
	}

	frame.Fields = fields
	return frame, nil
}

func newFieldForType(name, colType string, capacity int) *data.Field {
	switch {
	case strings.Contains(colType, "int"):
		return data.NewField(name, nil, make([]*int64, 0, capacity))
	case strings.Contains(colType, "real") || strings.Contains(colType, "float") || strings.Contains(colType, "double") || strings.Contains(colType, "numeric"):
		return data.NewField(name, nil, make([]*float64, 0, capacity))
	case strings.Contains(colType, "blob"):
		return data.NewField(name, nil, make([]*string, 0, capacity))
	default:
		// text, varchar, and anything else → string
		return data.NewField(name, nil, make([]*string, 0, capacity))
	}
}

func appendValue(field *data.Field, val interface{}, isTimeCol bool) {
	if isTimeCol {
		appendTimeValue(field, val)
		return
	}

	if val == nil {
		appendNilValue(field)
		return
	}

	switch field.Type() {
	case data.FieldTypeNullableInt64:
		v := toInt64(val)
		field.Append(&v)
	case data.FieldTypeNullableFloat64:
		v := toFloat64(val)
		field.Append(&v)
	case data.FieldTypeNullableString:
		v := fmt.Sprintf("%v", val)
		field.Append(&v)
	default:
		v := fmt.Sprintf("%v", val)
		field.Append(&v)
	}
}

func appendNilValue(field *data.Field) {
	switch field.Type() {
	case data.FieldTypeNullableInt64:
		field.Append((*int64)(nil))
	case data.FieldTypeNullableFloat64:
		field.Append((*float64)(nil))
	case data.FieldTypeNullableString:
		field.Append((*string)(nil))
	case data.FieldTypeNullableTime:
		field.Append((*time.Time)(nil))
	default:
		field.Append((*string)(nil))
	}
}

func appendTimeValue(field *data.Field, val interface{}) {
	if val == nil {
		field.Append((*time.Time)(nil))
		return
	}

	t := parseTime(val)
	field.Append(&t)
}

func parseTime(val interface{}) time.Time {
	switch v := val.(type) {
	case float64:
		return unixToTime(v)
	case string:
		// Try RFC3339 first
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			return t
		}
		// Try RFC3339Nano
		if t, err := time.Parse(time.RFC3339Nano, v); err == nil {
			return t
		}
		// Try common SQLite datetime format
		if t, err := time.Parse("2006-01-02 15:04:05", v); err == nil {
			return t
		}
		if t, err := time.Parse("2006-01-02", v); err == nil {
			return t
		}
		return time.Time{}
	default:
		return time.Time{}
	}
}

func unixToTime(v float64) time.Time {
	// Detect whether the value is seconds, milliseconds, microseconds, or nanoseconds
	// by magnitude. Seconds: < 1e12, Milliseconds: < 1e15, Microseconds: < 1e18
	abs := math.Abs(v)
	switch {
	case abs < 1e12:
		// Unix seconds (possibly fractional)
		sec := int64(v)
		nsec := int64((v - float64(sec)) * 1e9)
		return time.Unix(sec, nsec).UTC()
	case abs < 1e15:
		// Milliseconds
		return time.UnixMilli(int64(v)).UTC()
	case abs < 1e18:
		// Microseconds
		return time.UnixMicro(int64(v)).UTC()
	default:
		// Nanoseconds
		return time.Unix(0, int64(v)).UTC()
	}
}

func toInt64(val interface{}) int64 {
	switch v := val.(type) {
	case float64:
		return int64(v)
	case int64:
		return v
	case string:
		return 0
	default:
		return 0
	}
}

func toFloat64(val interface{}) float64 {
	switch v := val.(type) {
	case float64:
		return v
	case int64:
		return float64(v)
	case string:
		return 0
	default:
		return 0
	}
}
