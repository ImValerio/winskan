package mru

import (
	"encoding/binary"
	"fmt"
	"strconv"
	"unicode/utf16"

	"golang.org/x/sys/windows/registry"
)

// ParseRecentDocs fetches and parses the RecentDocs entries.
// It reads the MRUListEx binary value to determine chronological order.
func ParseRecentDocs() ([]Entry, error) {
	const basePath = `Software\Microsoft\Windows\CurrentVersion\Explorer\RecentDocs`

	key, err := registry.OpenKey(registry.CURRENT_USER, basePath, registry.QUERY_VALUE)
	if err != nil {
		if err == registry.ErrNotExist {
			return nil, nil // No RecentDocs history
		}
		return nil, fmt.Errorf("failed to open RecentDocs key: %w", err)
	}
	defer key.Close()

	mruListExData, _, err := key.GetBinaryValue("MRUListEx")
	if err != nil || len(mruListExData) < 4 {
		return nil, nil
	}

	var mruOrder []uint32
	for i := 0; i < len(mruListExData); i += 4 {
		if i+4 > len(mruListExData) {
			break
		}
		val := binary.LittleEndian.Uint32(mruListExData[i : i+4])
		if val == 0xFFFFFFFF {
			break
		}
		mruOrder = append(mruOrder, val)
	}

	var entries []Entry
	for i, valNum := range mruOrder {
		valName := strconv.FormatUint(uint64(valNum), 10)
		data, _, err := key.GetBinaryValue(valName)
		if err != nil || len(data) == 0 {
			continue
		}

		// Extract the file path from the binary blob
		filePath := extractUTF16String(data)

		entries = append(entries, Entry{
			Order:   i + 1,
			Name:    valName,
			Data:    filePath,
			MRUType: "RecentDocs",
		})
	}

	return entries, nil
}

// extractUTF16String attempts to extract a meaningful UTF-16LE string from the binary blob.
// RecentDocs binary data usually contains the file path as a null-terminated UTF-16 string.
func extractUTF16String(data []byte) string {
	longestStr := ""
	currentStr := make([]uint16, 0)

	for i := 0; i < len(data)-1; i += 2 {
		val := binary.LittleEndian.Uint16(data[i : i+2])
		// Check if it's a printable ASCII character or common European extended character
		if (val >= 0x20 && val <= 0x7E) || (val >= 0xA0 && val <= 0xFF) {
			currentStr = append(currentStr, val)
		} else if val == 0 { // null terminator
			if len(currentStr) > 0 {
				str := string(utf16.Decode(currentStr))
				if len(str) > len(longestStr) {
					longestStr = str
				}
				currentStr = make([]uint16, 0)
			}
		} else {
			// Some non-printable character that breaks the string
			currentStr = make([]uint16, 0)
		}
	}

	if len(currentStr) > 0 {
		str := string(utf16.Decode(currentStr))
		if len(str) > len(longestStr) {
			longestStr = str
		}
	}

	return longestStr
}
