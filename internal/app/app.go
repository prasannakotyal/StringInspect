package app

import (
	"fmt"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"stringinspect/internal/analysis"
	"stringinspect/internal/export"
	"stringinspect/internal/history"
)

// ViewMode represents the current display mode.
type ViewMode int

const (
	ViewModeTable ViewMode = iota
	ViewModeDetail
	ViewModeCompact
)

func (v ViewMode) String() string {
	switch v {
	case ViewModeTable:
		return "Table"
	case ViewModeDetail:
		return "Detail"
	case ViewModeCompact:
		return "Compact"
	default:
		return "Unknown"
	}
}

// App is the main Bubble Tea model for StringInspect.
type App struct {
	// Input
	input       textinput.Model
	searchInput textinput.Model
	analyzer    *analysis.Analyzer
	history     *history.History

	// State
	characters    []analysis.Character
	cursor        int
	viewMode      ViewMode
	showHelp      bool
	showExport    bool  // Export menu visible
	exportCursor  int   // Selected export format
	showSearch    bool  // Search mode active
	searchMatches []int // Indices of matching characters
	searchCursor  int   // Current match index
	statusMsg     string

	// Export
	exporter *export.Exporter

	// UI
	width  int
	height int
	styles Styles
	keys   KeyMap
	help   help.Model

	// Flags
	ready bool
	err   error
}

// New creates a new App instance.
func New() *App {
	return NewWithContent("")
}

// NewWithContent creates a new App instance with initial content.
func NewWithContent(content string) *App {
	ti := textinput.New()
	ti.Placeholder = "Type or paste text to analyze..."
	ti.Prompt = "> "
	ti.Focus()
	ti.CharLimit = 10000 // Increased for file content
	ti.Width = 60

	// Set initial content if provided
	if content != "" {
		ti.SetValue(content)
	}

	// Search input
	si := textinput.New()
	si.Placeholder = "hex, dec, or char..."
	si.Prompt = "/ "
	si.CharLimit = 50
	si.Width = 30

	h := help.New()
	h.ShowAll = false

	app := &App{
		input:       ti,
		searchInput: si,
		analyzer:    analysis.NewAnalyzer(),
		exporter:    export.NewExporter(),
		history:     history.New(100),
		styles:      DefaultStyles(),
		keys:        DefaultKeyMap(),
		help:        h,
		viewMode:    ViewModeTable,
	}

	// Analyze initial content if provided
	if content != "" {
		app.analyzeInput()
	}

	return app
}

// Init implements tea.Model.
func (a *App) Init() tea.Cmd {
	return textinput.Blink
}

// Update implements tea.Model.
func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		return a.handleKeyPress(msg)

	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		a.help.Width = msg.Width
		// Update input width to fit terminal (with padding)
		inputWidth := msg.Width - 8 // Account for prompt and padding
		if inputWidth > 200 {
			inputWidth = 200 // Cap at reasonable max
		}
		if inputWidth < 20 {
			inputWidth = 20 // Minimum width
		}
		a.input.Width = inputWidth
		a.ready = true
	}

	// Update text input
	var cmd tea.Cmd
	a.input, cmd = a.input.Update(msg)
	cmds = append(cmds, cmd)

	// Analyze input on change
	a.analyzeInput()

	return a, tea.Batch(cmds...)
}

// handleKeyPress processes keyboard input.
func (a *App) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Always allow quit (but not in search mode)
	if key.Matches(msg, a.keys.Quit) && !a.showSearch {
		return a, tea.Quit
	}

	// Handle search mode
	if a.showSearch {
		return a.handleSearchMode(msg)
	}

	// Handle export menu if visible
	if a.showExport {
		return a.handleExportMenu(msg)
	}

	// Toggle help
	if key.Matches(msg, a.keys.Help) {
		a.showHelp = !a.showHelp
		a.help.ShowAll = a.showHelp
		return a, nil
	}

	// If input is focused, let it handle most keys
	if a.input.Focused() {
		// Tab switches to navigation mode
		if key.Matches(msg, a.keys.Tab) {
			if len(a.characters) > 0 {
				a.input.Blur()
			}
			return a, nil
		}

		// History navigation with Up/Down
		if key.Matches(msg, a.keys.Up) {
			prev := a.history.Up(a.input.Value())
			a.input.SetValue(prev)
			a.input.CursorEnd()
			a.analyzeInput()
			return a, nil
		}
		if key.Matches(msg, a.keys.Down) {
			next := a.history.Down()
			a.input.SetValue(next)
			a.input.CursorEnd()
			a.analyzeInput()
			return a, nil
		}

		// Enter commits current input to history
		if key.Matches(msg, a.keys.Enter) {
			a.history.Add(a.input.Value())
			a.history.Reset()
			return a, nil
		}

		// Let input handle the key
		var cmd tea.Cmd
		a.input, cmd = a.input.Update(msg)
		a.analyzeInput()
		return a, cmd
	}

	// Navigation mode
	// Clear status message on navigation (but not on copy/paste)
	clearStatus := true

	switch {
	case key.Matches(msg, a.keys.Tab):
		// Cycle view mode or return to input
		if a.viewMode == ViewModeCompact {
			a.viewMode = ViewModeTable
			a.input.Focus()
		} else {
			a.viewMode++
		}

	case key.Matches(msg, a.keys.Left):
		if a.cursor > 0 {
			a.cursor--
		}

	case key.Matches(msg, a.keys.Right):
		if a.cursor < len(a.characters)-1 {
			a.cursor++
		}

	case key.Matches(msg, a.keys.Home):
		a.cursor = 0

	case key.Matches(msg, a.keys.End):
		if len(a.characters) > 0 {
			a.cursor = len(a.characters) - 1
		}

	case key.Matches(msg, a.keys.PageUp):
		// Move cursor up by page size (based on visible chars)
		pageSize := (a.width - 20) / 10
		if pageSize < 1 {
			pageSize = 1
		}
		a.cursor -= pageSize
		if a.cursor < 0 {
			a.cursor = 0
		}

	case key.Matches(msg, a.keys.PageDown):
		// Move cursor down by page size (based on visible chars)
		pageSize := (a.width - 20) / 10
		if pageSize < 1 {
			pageSize = 1
		}
		a.cursor += pageSize
		if a.cursor >= len(a.characters) {
			a.cursor = len(a.characters) - 1
		}
		if a.cursor < 0 {
			a.cursor = 0
		}

	case key.Matches(msg, a.keys.Copy):
		clearStatus = false
		// Copy selected character info to clipboard
		if a.cursor < len(a.characters) {
			char := a.characters[a.cursor]
			copyText := fmt.Sprintf("%s (U+%04X, 0x%s, %d)", char.Char, char.Dec, char.Hex, char.Dec)
			if err := clipboard.WriteAll(copyText); err != nil {
				a.statusMsg = "Copy failed"
			} else {
				a.statusMsg = fmt.Sprintf("Copied: %s", copyText)
			}
		}

	case key.Matches(msg, a.keys.Paste):
		clearStatus = false
		// Paste from clipboard
		if text, err := clipboard.ReadAll(); err == nil && text != "" {
			a.input.SetValue(text)
			a.analyzeInput()
			a.statusMsg = fmt.Sprintf("Pasted %d chars", len([]rune(text)))
		} else {
			a.statusMsg = "Paste failed"
		}

	case key.Matches(msg, a.keys.Export):
		// Open export menu if we have characters
		if len(a.characters) > 0 {
			a.showExport = true
			a.exportCursor = 0
		} else {
			a.statusMsg = "Nothing to export"
		}
		clearStatus = false

	case key.Matches(msg, a.keys.Search):
		// Enter search mode
		if len(a.characters) > 0 {
			a.showSearch = true
			a.searchInput.SetValue("")
			a.searchInput.Focus()
			a.searchMatches = nil
			a.searchCursor = 0
		} else {
			a.statusMsg = "Nothing to search"
		}
		clearStatus = false

	case key.Matches(msg, a.keys.Enter), key.Matches(msg, a.keys.Escape):
		a.input.Focus()
	}

	if clearStatus {
		a.statusMsg = ""
	}

	return a, nil
}

// analyzeInput processes the current input text.
func (a *App) analyzeInput() {
	input := a.input.Value()
	a.characters = a.analyzer.AnalyzeString(input)

	// Clear status message on input change
	a.statusMsg = ""

	// Keep cursor in bounds
	if a.cursor >= len(a.characters) {
		a.cursor = len(a.characters) - 1
	}
	if a.cursor < 0 {
		a.cursor = 0
	}
}

// handleExportMenu handles keyboard input for the export menu.
func (a *App) handleExportMenu(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if a.exportCursor > 0 {
			a.exportCursor--
		}
	case "down", "j":
		if a.exportCursor < 2 { // 3 formats: 0, 1, 2
			a.exportCursor++
		}
	case "enter":
		// Perform export
		format := export.Format(a.exportCursor)
		filename, err := a.exporter.Export(a.characters, format)
		if err != nil {
			a.statusMsg = fmt.Sprintf("Export failed: %v", err)
		} else {
			a.statusMsg = fmt.Sprintf("Exported to %s", filename)
		}
		a.showExport = false
	case "esc", "q":
		a.showExport = false
	case "1":
		a.exportCursor = 0
		return a.handleExportMenu(tea.KeyMsg{Type: tea.KeyEnter})
	case "2":
		a.exportCursor = 1
		return a.handleExportMenu(tea.KeyMsg{Type: tea.KeyEnter})
	case "3":
		a.exportCursor = 2
		return a.handleExportMenu(tea.KeyMsg{Type: tea.KeyEnter})
	}
	return a, nil
}

// handleSearchMode handles keyboard input for search mode.
func (a *App) handleSearchMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEsc:
		// Cancel search
		a.showSearch = false
		a.searchMatches = nil
		a.input.Focus()
		return a, nil

	case tea.KeyEnter:
		// Confirm search and jump to first match
		if len(a.searchMatches) > 0 {
			a.cursor = a.searchMatches[a.searchCursor]
			a.statusMsg = fmt.Sprintf("Match %d/%d", a.searchCursor+1, len(a.searchMatches))
		}
		a.showSearch = false
		a.searchInput.Blur()
		return a, nil

	case tea.KeyTab:
		// Cycle through matches
		if len(a.searchMatches) > 0 {
			a.searchCursor = (a.searchCursor + 1) % len(a.searchMatches)
			a.cursor = a.searchMatches[a.searchCursor]
		}
		return a, nil
	}

	// Let search input handle the key
	var cmd tea.Cmd
	a.searchInput, cmd = a.searchInput.Update(msg)

	// Perform search on input change
	a.performSearch()

	return a, cmd
}

// performSearch searches for characters matching the search query.
func (a *App) performSearch() {
	query := strings.ToLower(strings.TrimSpace(a.searchInput.Value()))
	if query == "" {
		a.searchMatches = nil
		a.searchCursor = 0
		return
	}

	var matches []int
	for i, char := range a.characters {
		// Match by character
		if strings.ToLower(char.Char) == query {
			matches = append(matches, i)
			continue
		}

		// Match by hex (with or without 0x prefix)
		hexQuery := strings.TrimPrefix(query, "0x")
		if strings.ToLower(char.Hex) == hexQuery {
			matches = append(matches, i)
			continue
		}

		// Match by decimal
		if fmt.Sprintf("%d", char.Dec) == query {
			matches = append(matches, i)
			continue
		}

		// Match by unicode (with or without U+ prefix)
		unicodeQuery := strings.TrimPrefix(strings.ToUpper(query), "U+")
		if strings.TrimPrefix(char.Unicode, "U+") == unicodeQuery {
			matches = append(matches, i)
			continue
		}
	}

	a.searchMatches = matches
	a.searchCursor = 0

	// Jump to first match
	if len(matches) > 0 {
		a.cursor = matches[0]
	}
}

// View implements tea.Model.
func (a *App) View() string {
	if !a.ready {
		return "Initializing..."
	}

	var b strings.Builder

	// Header
	b.WriteString(a.renderHeader())
	b.WriteString("\n\n")

	// Input
	b.WriteString(a.renderInput())
	b.WriteString("\n\n")

	// Content based on view mode
	if len(a.characters) > 0 {
		switch a.viewMode {
		case ViewModeTable:
			b.WriteString(a.renderTableView())
		case ViewModeDetail:
			b.WriteString(a.renderDetailView())
		case ViewModeCompact:
			b.WriteString(a.renderCompactView())
		}
	}

	// Status bar
	b.WriteString("\n\n")
	b.WriteString(a.renderStatusBar())

	// Search overlay
	if a.showSearch {
		b.WriteString("\n\n")
		b.WriteString(a.renderSearchBar())
	}

	// Export menu overlay
	if a.showExport {
		b.WriteString("\n\n")
		b.WriteString(a.renderExportMenu())
	}

	// Help
	if a.showHelp {
		b.WriteString("\n\n")
		b.WriteString(a.help.View(a.keys))
	}

	return a.styles.App.Render(b.String())
}

// renderHeader renders the application header.
func (a *App) renderHeader() string {
	title := a.styles.Header.Render("StringInspect")
	subtitle := a.styles.Muted.Render(" - Character Encoding Analyzer")
	return title + subtitle
}

// renderInput renders the text input field.
func (a *App) renderInput() string {
	return a.input.View()
}

// renderTableView renders the table view of character encodings.
func (a *App) renderTableView() string {
	var b strings.Builder

	// Calculate visible characters based on width
	maxChars := (a.width - 20) / 10
	if maxChars < 1 {
		maxChars = 1
	}
	if maxChars > len(a.characters) {
		maxChars = len(a.characters)
	}

	// Determine scroll offset to keep cursor visible
	start := 0
	if a.cursor >= maxChars {
		start = a.cursor - maxChars + 1
	}
	end := start + maxChars
	if end > len(a.characters) {
		end = len(a.characters)
		start = end - maxChars
		if start < 0 {
			start = 0
		}
	}

	visibleChars := a.characters[start:end]

	// Render rows
	rows := []struct {
		label string
		fn    func(c analysis.Character) string
	}{
		{"Char", func(c analysis.Character) string { return c.Char }},
		{"Hex", func(c analysis.Character) string { return c.Hex }},
		{"Dec", func(c analysis.Character) string { return fmt.Sprintf("%d", c.Dec) }},
		{"Bin", func(c analysis.Character) string { return c.Bin }},
		{"Oct", func(c analysis.Character) string { return c.Oct }},
		{"Unicode", func(c analysis.Character) string { return c.Unicode }},
	}

	for _, row := range rows {
		label := a.styles.TableLabel.Render(row.label)
		b.WriteString(label)

		for i, char := range visibleChars {
			globalIdx := start + i
			value := row.fn(char)

			var style lipgloss.Style
			if globalIdx == a.cursor && !a.input.Focused() {
				style = a.styles.TableSelected
			} else {
				style = a.styles.CharStyle(int(char.Type))
			}

			cell := style.Width(10).Align(lipgloss.Center).Render(value)
			b.WriteString(cell)
		}
		b.WriteString("\n")
	}

	return b.String()
}

// renderDetailView renders a detailed view of the selected character.
func (a *App) renderDetailView() string {
	if a.cursor >= len(a.characters) {
		return ""
	}

	char := a.characters[a.cursor]
	var b strings.Builder

	title := a.styles.Title.Render("Character Details")
	b.WriteString(title)
	b.WriteString("\n\n")

	// Character display
	charStyle := a.styles.CharStyle(int(char.Type))
	charDisplay := charStyle.Bold(true).Padding(1, 3).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorPrimary).
		Render(char.Char)
	b.WriteString(charDisplay)
	b.WriteString("\n\n")

	// Details table
	details := []struct {
		label string
		value string
	}{
		{"Unicode", char.Unicode},
		{"Hexadecimal", "0x" + char.Hex},
		{"Decimal", fmt.Sprintf("%d", char.Dec)},
		{"Octal", "0o" + char.Oct},
		{"Binary", char.Bin},
		{"UTF-8 Bytes", char.UTF8Hex},
		{"Position", fmt.Sprintf("%d (byte: %d)", char.RuneOffset, char.ByteOffset)},
	}

	for _, d := range details {
		label := a.styles.Muted.Width(14).Render(d.label + ":")
		value := a.styles.Printable.Render(d.value)
		b.WriteString(label + " " + value + "\n")
	}

	// Navigation hint
	b.WriteString("\n")
	hint := a.styles.Muted.Render(fmt.Sprintf("← → to navigate (%d/%d)", a.cursor+1, len(a.characters)))
	b.WriteString(hint)

	return b.String()
}

// renderCompactView renders a hex dump style view.
func (a *App) renderCompactView() string {
	var b strings.Builder

	title := a.styles.Title.Render("Compact View (Hex Dump)")
	b.WriteString(title)
	b.WriteString("\n\n")

	// Show offset | hex values | ascii
	charsPerLine := 16
	for i := 0; i < len(a.characters); i += charsPerLine {
		// Offset
		offset := a.styles.Muted.Render(fmt.Sprintf("%04X  ", i))
		b.WriteString(offset)

		// Hex values
		for j := 0; j < charsPerLine; j++ {
			idx := i + j
			if idx < len(a.characters) {
				char := a.characters[idx]
				style := a.styles.CharStyle(int(char.Type))
				if idx == a.cursor && !a.input.Focused() {
					style = a.styles.TableSelected
				}
				hex := style.Render(char.Hex)
				b.WriteString(hex + " ")
			} else {
				b.WriteString("   ")
			}

			// Extra space in middle
			if j == 7 {
				b.WriteString(" ")
			}
		}

		b.WriteString(" │ ")

		// ASCII representation
		for j := 0; j < charsPerLine; j++ {
			idx := i + j
			if idx < len(a.characters) {
				char := a.characters[idx]
				style := a.styles.CharStyle(int(char.Type))
				if idx == a.cursor && !a.input.Focused() {
					style = a.styles.TableSelected
				}

				display := char.Char
				if len(display) > 1 {
					display = "."
				}
				b.WriteString(style.Render(display))
			}
		}

		b.WriteString("\n")
	}

	return b.String()
}

// renderStatusBar renders the status bar.
func (a *App) renderStatusBar() string {
	// Mode indicator
	mode := a.viewMode.String()
	if a.input.Focused() {
		mode = "Input"
	}

	// Character count
	charCount := fmt.Sprintf("%d chars", len(a.characters))

	// Build status
	left := a.styles.Muted.Render(fmt.Sprintf("[%s]", mode))

	// Show status message if present, otherwise show default help hints
	var right string
	if a.statusMsg != "" {
		right = a.styles.Success.Render(a.statusMsg)
	} else {
		right = a.styles.Muted.Render(charCount + " │ F1 help │ q quit")
	}

	gap := a.width - lipgloss.Width(left) - lipgloss.Width(right) - 4
	if gap < 1 {
		gap = 1
	}

	return left + strings.Repeat(" ", gap) + right
}

// renderExportMenu renders the export format selection menu.
func (a *App) renderExportMenu() string {
	var b strings.Builder

	// Menu box
	title := a.styles.Title.Render("Export Format")
	b.WriteString(title)
	b.WriteString("\n\n")

	formats := []struct {
		key    string
		name   string
		desc   string
		format export.Format
	}{
		{"1", "Text", "Plain text table", export.FormatText},
		{"2", "JSON", "Structured JSON", export.FormatJSON},
		{"3", "CSV", "Comma-separated values", export.FormatCSV},
	}

	for i, f := range formats {
		prefix := "  "
		style := a.styles.Muted
		if i == a.exportCursor {
			prefix = "> "
			style = a.styles.Highlighted
		}

		line := fmt.Sprintf("%s[%s] %s - %s", prefix, f.key, f.name, f.desc)
		b.WriteString(style.Render(line))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	hint := a.styles.Muted.Render("↑/↓ select • enter confirm • esc cancel")
	b.WriteString(hint)

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorPrimary).
		Padding(1, 2).
		Render(b.String())
}

// renderSearchBar renders the search input bar.
func (a *App) renderSearchBar() string {
	var b strings.Builder

	title := a.styles.Title.Render("Search")
	b.WriteString(title)
	b.WriteString("\n\n")

	// Search input
	b.WriteString(a.searchInput.View())
	b.WriteString("\n\n")

	// Match count
	if len(a.searchMatches) > 0 {
		matchInfo := fmt.Sprintf("Found %d match(es) - Tab to cycle, Enter to confirm", len(a.searchMatches))
		b.WriteString(a.styles.Success.Render(matchInfo))
	} else if a.searchInput.Value() != "" {
		b.WriteString(a.styles.Error.Render("No matches"))
	} else {
		b.WriteString(a.styles.Muted.Render("Type hex (0x41), decimal (65), or character (A)"))
	}

	b.WriteString("\n\n")
	hint := a.styles.Muted.Render("enter confirm • esc cancel • tab next match")
	b.WriteString(hint)

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(ColorPrimary).
		Padding(1, 2).
		Render(b.String())
}
