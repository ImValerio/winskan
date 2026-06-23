package mru

// Entry represents a single MRU (Most Recently Used) record.
type Entry struct {
	Order   int    // The chronological position in the MRU list (1 is most recent)
	Name    string // The registry value name (e.g., "a", "b", "0", "1")
	Data    string // The decoded string, command, or file path
	MRUType string // The type of MRU list (e.g., "RunMRU", "RecentDocs")
}
