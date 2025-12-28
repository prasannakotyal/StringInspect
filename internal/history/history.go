// Package history manages input history for the application.
package history

// History stores previous input strings for navigation.
type History struct {
	entries []string
	cursor  int    // Current position in history (-1 means not browsing)
	limit   int    // Maximum entries to store
	current string // Temporarily stores current input while browsing
}

// New creates a new History with the specified limit.
func New(limit int) *History {
	if limit < 1 {
		limit = 100
	}
	return &History{
		entries: make([]string, 0, limit),
		cursor:  -1,
		limit:   limit,
	}
}

// Add adds a new entry to the history.
// Empty strings and duplicates of the last entry are ignored.
func (h *History) Add(entry string) {
	if entry == "" {
		return
	}

	// Don't add if it's the same as the last entry
	if len(h.entries) > 0 && h.entries[len(h.entries)-1] == entry {
		return
	}

	// Add to history
	h.entries = append(h.entries, entry)

	// Trim if over limit
	if len(h.entries) > h.limit {
		h.entries = h.entries[1:]
	}

	// Reset cursor
	h.cursor = -1
	h.current = ""
}

// Up moves up in history (to older entries).
// Returns the entry at the new position, or empty string if at the beginning.
// currentInput is saved on first Up press so it can be restored.
func (h *History) Up(currentInput string) string {
	if len(h.entries) == 0 {
		return currentInput
	}

	// First time pressing up - save current input
	if h.cursor == -1 {
		h.current = currentInput
		h.cursor = len(h.entries) - 1
	} else if h.cursor > 0 {
		h.cursor--
	}

	return h.entries[h.cursor]
}

// Down moves down in history (to newer entries).
// Returns the entry at the new position, or the saved current input if at the end.
func (h *History) Down() string {
	if h.cursor == -1 {
		return h.current
	}

	h.cursor++

	// If we've moved past the last entry, return to current input
	if h.cursor >= len(h.entries) {
		h.cursor = -1
		return h.current
	}

	return h.entries[h.cursor]
}

// Reset resets the history browsing state.
func (h *History) Reset() {
	h.cursor = -1
	h.current = ""
}

// Len returns the number of entries in history.
func (h *History) Len() int {
	return len(h.entries)
}

// IsBrowsing returns true if currently browsing history.
func (h *History) IsBrowsing() bool {
	return h.cursor != -1
}
