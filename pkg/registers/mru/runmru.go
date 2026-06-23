package mru

import (
	"fmt"
	"strings"

	"golang.org/x/sys/windows/registry"
)

// ParseRunMRU fetches and parses the RunMRU entries. (Win + R)
// It reads the MRUList value to determine the chronological order
// and then reads the corresponding string values.
func ParseRunMRU() ([]Entry, error) {
	const basePath = `Software\Microsoft\Windows\CurrentVersion\Explorer\RunMRU`

	key, err := registry.OpenKey(registry.CURRENT_USER, basePath, registry.QUERY_VALUE)
	if err != nil {
		if err == registry.ErrNotExist {
			return nil, nil // No RunMRU history
		}
		return nil, fmt.Errorf("failed to open RunMRU key: %w", err)
	}
	defer key.Close()

	mruListStr, _, err := key.GetStringValue("MRUList")
	if err != nil {
		// MRUList doesn't exist or is empty
		return nil, nil
	}

	var entries []Entry

	// The MRUList is a string of characters, e.g., "bac" means 'b' is most recent, then 'a', then 'c'.
	for i, char := range mruListStr {
		valName := string(char)
		data, _, err := key.GetStringValue(valName)
		if err != nil {
			continue
		}

		// RunMRU entries conventionally end with "\1" (e.g., "cmd.exe\1"). We clean it up.
		data = strings.TrimSuffix(data, "\\1")

		entries = append(entries, Entry{
			Order:   i + 1,
			Name:    valName,
			Data:    data,
			MRUType: "RunMRU",
		})
	}

	return entries, nil
}
