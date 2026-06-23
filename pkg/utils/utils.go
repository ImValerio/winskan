package utils

import (
	"strings"
	"time"
)

// Rot13 decodes or encodes the ROT13 encoded registry value names
func Rot13(s string) string {
	var result strings.Builder
	for _, r := range s {
		if r >= 'a' && r <= 'z' {
			result.WriteRune('a' + (r-'a'+13)%26)
		} else if r >= 'A' && r <= 'Z' {
			result.WriteRune('A' + (r-'A'+13)%26)
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// FiletimeToTime converts a Windows FILETIME uint64 to a standard Go time.Time.
// A Windows FILETIME represents the number of 100-nanosecond intervals since January 1, 1601 (UTC).
func FiletimeToTime(ft uint64) time.Time {
	if ft == 0 {
		return time.Time{}
	}
	const unixEpochOffset = 116444736000000000
	nsec := int64(ft-unixEpochOffset) * 100
	return time.Unix(0, nsec).UTC()
}
