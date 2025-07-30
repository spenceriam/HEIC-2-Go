// Package ui handles the terminal user interface
package ui

import (
	"fmt"
	"os"
	"strings"

	"github.com/common-nighthawk/go-figure"
	"github.com/fatih/color"
)

// Screen represents a terminal screen
type Screen struct {
	width  int
	height int
}

// NewScreen creates a new Screen instance
func NewScreen() *Screen {
	// Default dimensions, can be updated with terminal size detection
	return &Screen{
		width:  80,
		height: 24,
	}
}

// Clear clears the terminal screen
func (s *Screen) Clear() {
	// This will be replaced with proper terminal clearing
	fmt.Print("\033[H\033[2J")
}

// DisplayWelcome shows the welcome screen
func (s *Screen) DisplayWelcome() {
	s.Clear()

	// Create the ASCII art
	fig := figure.NewFigure("HEIC-2-Go", "doom", true)
	asciiArt := fig.String()

	// Create the box border
	border := strings.Repeat("═", s.width-2)
	titleBox := fmt.Sprintf("╔%s╗\n║%s║\n╚%s╝", 
		border, 
		s.centerText("Convert HEIC files to JPG", s.width-2), 
		border)

	// Display the ASCII art and title box
	color.Cyan(asciiArt)
	color.HiCyan(titleBox)
	fmt.Println()
}

// ShowMenu displays the main menu and handles user input
func (s *Screen) ShowMenu() error {
	for {
		s.DisplayWelcome()
		
		// Create and display the menu
		menu := NewMainMenu(s)
		
		// Show menu options
		options := ""
		for _, option := range menu.Options {
			options += fmt.Sprintf("  [%s] %s\n", option.Key, option.Description)
		}
		
		// Display the menu in a box
		menuBox := fmt.Sprintf(`
┌─────────────────────────────────────────────────────────────────────┐
%s└─────────────────────────────────────────────────────────────────────┘

Enter your choice (1-%d): `, 
			options, 
			len(menu.Options))
		
		// Get user selection
		choice, err := s.GetIntInput(menuBox, 1, len(menu.Options))
		if err != nil {
			s.ShowMessage(fmt.Sprintf("Error: %v", err))
			continue
		}
		
		// Execute the selected option
		for _, option := range menu.Options {
			if option.Key == fmt.Sprintf("%d", choice) {
				if err := option.Handler(); err != nil {
					s.ShowMessage(fmt.Sprintf("Error: %v", err))
				}
				break
			}
		}
	}
}

// centerText centers text within a given width
func (s *Screen) centerText(text string, width int) string {
	if len(text) >= width {
		return text[:width]
	}

	padding := (width - len(text)) / 2
	return strings.Repeat(" ", padding) + text + strings.Repeat(" ", width-len(text)-padding)
}

// ShowError displays an error message
func (s *Screen) ShowError(message string) {
	color.Red("\nError: %s\n", message)
	fmt.Println("Press Enter to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
