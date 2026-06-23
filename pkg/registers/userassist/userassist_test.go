package userassist

import (
	"encoding/binary"
	"testing"

	"github.com/imvalerio/winskan/pkg/utils"
)

func TestParseBinaryData(t *testing.T) {
	// Construct a dummy 72-byte block
	data := make([]byte, 72)
	
	// RunCount at offset 4
	binary.LittleEndian.PutUint32(data[4:8], 5)
	
	// FocusCount at offset 8
	binary.LittleEndian.PutUint32(data[8:12], 10)
	
	// FocusTimeMs at offset 12
	binary.LittleEndian.PutUint32(data[12:16], 1500)
	
	// LastExecution at offset 60 (Using roughly Jan 1, 2020)
	var filetime uint64 = 132222816000000000
	binary.LittleEndian.PutUint64(data[60:68], filetime)

	name := "C:\\Test\\App.exe"

	entry, ok := ParseBinaryData(name, data)
	if !ok {
		t.Fatalf("ParseBinaryData failed, expected true")
	}

	if entry.Name != name {
		t.Errorf("Expected Name %v, got %v", name, entry.Name)
	}
	if entry.RunCount != 5 {
		t.Errorf("Expected RunCount 5, got %v", entry.RunCount)
	}
	if entry.FocusCount != 10 {
		t.Errorf("Expected FocusCount 10, got %v", entry.FocusCount)
	}
	if entry.FocusTimeMs != 1500 {
		t.Errorf("Expected FocusTimeMs 1500, got %v", entry.FocusTimeMs)
	}

	expectedTime := utils.FiletimeToTime(filetime)
	if !entry.LastRun.Equal(expectedTime) {
		t.Errorf("Expected LastRun %v, got %v", expectedTime, entry.LastRun)
	}
}

func TestParseBinaryDataInvalidLength(t *testing.T) {
	data := make([]byte, 71) // Incorrect length
	_, ok := ParseBinaryData("test", data)
	if ok {
		t.Errorf("ParseBinaryData should fail for incorrect data length")
	}
}
