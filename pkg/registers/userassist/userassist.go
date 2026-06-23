package userassist

import (
	"encoding/binary"
	"fmt"
	"strings"
	"time"

	"github.com/imvalerio/winskan/pkg/utils"
	"golang.org/x/sys/windows/registry"
)

const (
	GUIDExecutables = "{CEBFF5CD-ACE2-4F4F-9178-9926F41749EA}"
	GUIDShortcuts   = "{F4E57C4B-2036-45F0-A9AB-443BCFE33D9F}"
)

// Entry represents a single parsed UserAssist record.
type Entry struct {
	Name        string
	RunCount    uint32
	FocusCount  uint32
	FocusTimeMs uint32
	LastRun     time.Time
}

// Parse fetches and parses UserAssist entries for the given GUID.
func Parse(guid string) ([]Entry, error) {
	basePath := fmt.Sprintf(`Software\Microsoft\Windows\CurrentVersion\Explorer\UserAssist\%s\Count`, guid)

	key, err := registry.OpenKey(registry.CURRENT_USER, basePath, registry.QUERY_VALUE)
	if err != nil {
		return nil, fmt.Errorf("failed to open registry key: %w", err)
	}
	defer key.Close()

	valNames, err := key.ReadValueNames(-1)
	if err != nil {
		return nil, fmt.Errorf("failed to read value names: %w", err)
	}

	var entries []Entry

	for _, valName := range valNames {
		decodedName := utils.Rot13(valName)
		if strings.HasPrefix(decodedName, "UEME_") {
			continue // Skip UEME entries as per original behavior
		}

		data, _, err := key.GetBinaryValue(valName)
		if err != nil {
			continue
		}

		entry, ok := ParseBinaryData(decodedName, data)
		if ok {
			entries = append(entries, entry)
		}
	}

	return entries, nil
}

// ParseBinaryData parses the 72-byte UserAssist binary data.
func ParseBinaryData(name string, data []byte) (Entry, bool) {
	if len(data) != 72 {
		return Entry{}, false
	}

	runCount := binary.LittleEndian.Uint32(data[4:8])
	focusCount := binary.LittleEndian.Uint32(data[8:12])
	focusTimeMs := binary.LittleEndian.Uint32(data[12:16])
	lastExecutionRaw := binary.LittleEndian.Uint64(data[60:68])
	lastExecution := utils.FiletimeToTime(lastExecutionRaw)

	return Entry{
		Name:        name,
		RunCount:    runCount,
		FocusCount:  focusCount,
		FocusTimeMs: focusTimeMs,
		LastRun:     lastExecution,
	}, true
}
