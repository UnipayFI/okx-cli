package common

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const dateTimeLayout = "2006-01-02 15:04:05"

// FormatTime renders a time.Time as a UTC string in the project layout. The OKX
// SDK already decodes API timestamps into time.Time, so this is the single
// place CLI output is made timezone-consistent.
func FormatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.UTC().Format(dateTimeLayout)
}

// ParseTimeFlag accepts either a unix-milliseconds timestamp or a
// "YYYY-MM-DD HH:MM:SS" datetime (parsed in local time) and returns the
// corresponding time.Time. An empty or "0" value yields ok=false so callers can
// treat the flag as unset.
func ParseTimeFlag(flagName string, value string) (time.Time, bool, error) {
	v := strings.TrimSpace(value)
	if v == "" || v == "0" {
		return time.Time{}, false, nil
	}

	if isAllDigits(v) {
		ms, err := strconv.ParseInt(v, 10, 64)
		if err != nil || ms < 0 {
			return time.Time{}, false, invalidTimeFlagError(flagName)
		}
		return time.UnixMilli(ms), true, nil
	}

	t, err := time.ParseInLocation(dateTimeLayout, v, time.Local)
	if err != nil {
		return time.Time{}, false, invalidTimeFlagError(flagName)
	}
	return t, true, nil
}

func invalidTimeFlagError(flagName string) error {
	name := strings.TrimSpace(flagName)
	if name == "" {
		name = "time"
	}
	return fmt.Errorf("invalid %s: expected unix milliseconds timestamp (e.g. 1734495381000) or datetime \"YYYY-MM-DD HH:MM:SS\" (e.g. \"2025-12-18 04:16:21\")", name)
}

func isAllDigits(s string) bool {
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return len(s) > 0
}
