package main

import (
	"fmt"
	"os"

	"github.com/spenceriam/HEIC-2-Go/internal/ui"
	"github.com/spenceriam/HEIC-2-Go/pkg/version"
)

const (
	appName = "HEIC-2-Go"
)

func main() {
	// Initialize the terminal UI
	screen := ui.NewScreen()

	// Start the main menu
	if err := screen.ShowMenu(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
