package app

import "github.com/charmbracelet/lipgloss"

// Color palette - Charmbracelet-inspired dark theme
var (
	ColorPrimary    = lipgloss.Color("#7D56F4") // Purple - selection, focus
	ColorSuccess    = lipgloss.Color("#73F59F") // Green - confirmations
	ColorError      = lipgloss.Color("#FF4672") // Red - errors
	ColorWarning    = lipgloss.Color("#FDFF90") // Yellow - warnings
	ColorSubtle     = lipgloss.Color("#383838") // Dark gray - borders
	ColorMuted      = lipgloss.Color("#929292") // Gray - muted text
	ColorText       = lipgloss.Color("#EEEEEE") // Off-white - default text
	ColorWhitespace = lipgloss.Color("#00E2C7") // Cyan - whitespace chars
	ColorControl    = lipgloss.Color("#FF7698") // Pink/red - control chars
	ColorExtended   = lipgloss.Color("#FDFF90") // Yellow - extended ASCII
	ColorBackground = lipgloss.Color("#1a1a1a") // Dark background
)

// Styles holds all application styles.
type Styles struct {
	// App-level styles
	App       lipgloss.Style
	Header    lipgloss.Style
	StatusBar lipgloss.Style

	// Content styles
	Title       lipgloss.Style
	Subtitle    lipgloss.Style
	Muted       lipgloss.Style
	Error       lipgloss.Style
	Success     lipgloss.Style
	Highlighted lipgloss.Style

	// Character type styles
	Printable  lipgloss.Style
	Whitespace lipgloss.Style
	Control    lipgloss.Style
	Extended   lipgloss.Style

	// Table styles
	TableHeader   lipgloss.Style
	TableCell     lipgloss.Style
	TableSelected lipgloss.Style
	TableLabel    lipgloss.Style

	// Input styles
	InputPrompt lipgloss.Style
	InputText   lipgloss.Style

	// Help styles
	HelpKey  lipgloss.Style
	HelpDesc lipgloss.Style
	HelpSep  lipgloss.Style
}

// DefaultStyles returns the default application styles.
func DefaultStyles() Styles {
	return Styles{
		// App-level styles
		App: lipgloss.NewStyle().
			Padding(1, 2),

		Header: lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary).
			MarginBottom(1),

		StatusBar: lipgloss.NewStyle().
			Foreground(ColorMuted).
			MarginTop(1),

		// Content styles
		Title: lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorText),

		Subtitle: lipgloss.NewStyle().
			Foreground(ColorMuted),

		Muted: lipgloss.NewStyle().
			Foreground(ColorMuted),

		Error: lipgloss.NewStyle().
			Foreground(ColorError),

		Success: lipgloss.NewStyle().
			Foreground(ColorSuccess),

		Highlighted: lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorText).
			Background(ColorPrimary).
			Padding(0, 1),

		// Character type styles
		Printable: lipgloss.NewStyle().
			Foreground(ColorText),

		Whitespace: lipgloss.NewStyle().
			Foreground(ColorWhitespace),

		Control: lipgloss.NewStyle().
			Foreground(ColorControl),

		Extended: lipgloss.NewStyle().
			Foreground(ColorExtended),

		// Table styles
		TableHeader: lipgloss.NewStyle().
			Bold(true).
			Foreground(ColorPrimary).
			BorderBottom(true).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(ColorSubtle),

		TableCell: lipgloss.NewStyle().
			Padding(0, 1),

		TableSelected: lipgloss.NewStyle().
			Background(ColorPrimary).
			Foreground(ColorText).
			Bold(true).
			Padding(0, 1),

		TableLabel: lipgloss.NewStyle().
			Foreground(ColorMuted).
			Width(8),

		// Input styles
		InputPrompt: lipgloss.NewStyle().
			Foreground(ColorPrimary).
			Bold(true),

		InputText: lipgloss.NewStyle().
			Foreground(ColorText),

		// Help styles
		HelpKey: lipgloss.NewStyle().
			Foreground(ColorPrimary).
			Bold(true),

		HelpDesc: lipgloss.NewStyle().
			Foreground(ColorMuted),

		HelpSep: lipgloss.NewStyle().
			Foreground(ColorSubtle),
	}
}

// CharStyle returns the appropriate style for a character type.
func (s Styles) CharStyle(charType int) lipgloss.Style {
	switch charType {
	case 1: // Whitespace
		return s.Whitespace
	case 2: // Control
		return s.Control
	case 3: // Extended
		return s.Extended
	default: // Printable
		return s.Printable
	}
}
