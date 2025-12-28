package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"stringinspect/internal/app"
)

func main() {
	// Parse command line flags
	filePath := flag.String("f", "", "Path to file to analyze")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "StringInspect - Interactive Character Encoding Analyzer\n\n")
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s                    # Start interactive mode\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -f file.txt        # Analyze file contents\n", os.Args[0])
	}
	flag.Parse()

	// Create the application
	var a *app.App
	if *filePath != "" {
		// Read file contents
		content, err := os.ReadFile(*filePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
			os.Exit(1)
		}
		a = app.NewWithContent(string(content))
	} else {
		a = app.New()
	}

	// Create and run the program
	p := tea.NewProgram(a, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running application: %v\n", err)
		os.Exit(1)
	}
}
