package utils

import (
	"testing"
	"time"
)

func TestRot13(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "lowercase",
			input:    "uryyb",
			expected: "hello",
		},
		{
			name:     "uppercase",
			input:    "URYYB",
			expected: "HELLO",
		},
		{
			name:     "mixed case",
			input:    "Uryyb Jbeyq!",
			expected: "Hello World!",
		},
		{
			name:     "UserAssist specific prefix",
			input:    "HRZR_EHACNGU", // UEME_RUNPATH
			expected: "UEME_RUNPATH",
		},
		{
			name:     "numbers and symbols",
			input:    "123!@#",
			expected: "123!@#",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Rot13(tt.input)
			if got != tt.expected {
				t.Errorf("Rot13() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestFiletimeToTime(t *testing.T) {
	tests := []struct {
		name     string
		ft       uint64
		expected time.Time
	}{
		{
			name:     "zero value",
			ft:       0,
			expected: time.Time{},
		},
		{
			name:     "Windows Epoch",
			ft:       0,
			expected: time.Time{},
		},
		{
			name:     "Unix Epoch",
			ft:       116444736000000000,
			expected: time.Unix(0, 0).UTC(),
		},
		{
			name:     "Known date",
			ft:       132222816000000000,
			expected: time.Date(2019, 12, 31, 16, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FiletimeToTime(tt.ft)
			if !got.Equal(tt.expected) {
				t.Errorf("FiletimeToTime() = %v, want %v", got, tt.expected)
			}
		})
	}
}
