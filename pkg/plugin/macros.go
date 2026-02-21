package plugin

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

var (
	timeFilterRegex      = regexp.MustCompile(`\$__timeFilter\((\w+)\)`)
	unixEpochFilterRegex = regexp.MustCompile(`\$__unixEpochFilter\((\w+)\)`)
	timeGroupRegex       = regexp.MustCompile(`\$__timeGroup\((\w+)\s*,\s*([^)]+)\)`)
	timeFromRegex        = regexp.MustCompile(`\$__timeFrom`)
	timeToRegex          = regexp.MustCompile(`\$__timeTo`)
)

// ApplyMacros replaces Grafana macros in a SQL string with SQLite-compatible expressions.
func ApplyMacros(sql string, timeRange backend.TimeRange, intervalMS int64) string {
	fromUnix := timeRange.From.Unix()
	toUnix := timeRange.To.Unix()

	// $__timeFilter(col) → col >= <from> AND col <= <to>
	sql = timeFilterRegex.ReplaceAllStringFunc(sql, func(match string) string {
		sub := timeFilterRegex.FindStringSubmatch(match)
		if len(sub) < 2 {
			return match
		}
		col := sub[1]
		return fmt.Sprintf("%s >= %d AND %s <= %d", col, fromUnix, col, toUnix)
	})

	// $__unixEpochFilter(col) → same as $__timeFilter
	sql = unixEpochFilterRegex.ReplaceAllStringFunc(sql, func(match string) string {
		sub := unixEpochFilterRegex.FindStringSubmatch(match)
		if len(sub) < 2 {
			return match
		}
		col := sub[1]
		return fmt.Sprintf("%s >= %d AND %s <= %d", col, fromUnix, col, toUnix)
	})

	// $__timeGroup(col, interval) → (CAST(col / N AS INTEGER) * N)
	sql = timeGroupRegex.ReplaceAllStringFunc(sql, func(match string) string {
		sub := timeGroupRegex.FindStringSubmatch(match)
		if len(sub) < 3 {
			return match
		}
		col := sub[1]
		intervalStr := strings.TrimSpace(sub[2])
		seconds := parseInterval(intervalStr, intervalMS)
		if seconds <= 0 {
			seconds = 60 // default to 1 minute
		}
		return fmt.Sprintf("(CAST(%s / %d AS INTEGER) * %d)", col, seconds, seconds)
	})

	// $__timeFrom → unix epoch seconds
	sql = timeFromRegex.ReplaceAllString(sql, strconv.FormatInt(fromUnix, 10))

	// $__timeTo → unix epoch seconds
	sql = timeToRegex.ReplaceAllString(sql, strconv.FormatInt(toUnix, 10))

	return sql
}

// parseInterval converts an interval string to seconds.
// Supports: "1s", "5m", "1h", "1d", plain integer seconds, "$__interval".
func parseInterval(s string, intervalMS int64) int64 {
	s = strings.TrimSpace(s)

	if s == "$__interval" {
		return intervalMS / 1000
	}

	// Try plain integer (seconds)
	if n, err := strconv.ParseInt(s, 10, 64); err == nil {
		return n
	}

	// Try Go duration format
	if d, err := time.ParseDuration(s); err == nil {
		return int64(d.Seconds())
	}

	// Try custom format with 'd' suffix for days
	if strings.HasSuffix(s, "d") {
		if n, err := strconv.ParseInt(strings.TrimSuffix(s, "d"), 10, 64); err == nil {
			return n * 86400
		}
	}

	return 0
}
