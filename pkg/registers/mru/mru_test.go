package mru

import (
	"testing"
)

func TestParseRunMRU(t *testing.T) {
	entries, err := ParseRunMRU()
	if err != nil {
		t.Fatalf("ParseRunMRU failed: %v", err)
	}

	t.Logf("Found %d RunMRU entries", len(entries))
	for _, entry := range entries {
		t.Logf("[%d] %s: %s", entry.Order, entry.Name, entry.Data)
	}
}

func TestParseRecentDocs(t *testing.T) {
	entries, err := ParseRecentDocs()
	if err != nil {
		t.Fatalf("ParseRecentDocs failed: %v", err)
	}

	t.Logf("Found %d RecentDocs entries", len(entries))
	for _, entry := range entries {
		t.Logf("[%d] %s: %s", entry.Order, entry.Name, entry.Data)
	}
}
