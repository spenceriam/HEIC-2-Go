package ui

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// MenuOption represents a single menu option
type MenuOption struct {
	Key         string
	Description string
	Handler     func() error
}

// Menu represents a menu with multiple options
type Menu struct {
	Title   string
	Options []MenuOption
}

// NewMainMenu creates the main menu with all options
func NewMainMenu(screen *Screen) *Menu {
	return &Menu{
		Title: "Main Menu",
		Options: []MenuOption{
			{
				Key:         "1",
				Description: "Convert single file",
				Handler:     screen.handleSingleFile,
			},
			{
				Key:         "2",
				Description: "Convert directory",
				Handler:     screen.handleDirectory,
			},
			{
				Key:         "3",
				Description: "Settings",
				Handler:     screen.handleSettings,
			},
			{
				Key:         "4",
				Description: "Exit",
				Handler:     screen.handleExit,
			},
		},
	}
}

// Display shows the menu and handles user input
func (m *Menu) Display() error {
	for {
		// Clear the screen and show the welcome message
		// (This will be replaced with a proper screen refresh)
		fmt.Println()
		fmt.Println(m.Title)
		fmt.Println(strings.Repeat("-", len(m.Title)))

		// Display all options
		for _, option := range m.Options {
			fmt.Printf("[%s] %s\n", option.Key, option.Description)
		}

		// Get user input
		fmt.Print("\nEnter your choice: ")
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		// Find and execute the selected option
		for _, option := range m.Options {
			if strings.EqualFold(input, option.Key) {
				return option.Handler()
			}
		}

		// If we get here, the input was invalid
		fmt.Printf("\nInvalid option: %s\n\n", input)
	}
}

// handleSingleFile handles the single file conversion option
func (s *Screen) handleSingleFile() error {
	fileInput := NewFileInputScreen(s)
	
	// Show file input screen
	filePath, err := fileInput.Show()
	if err != nil {
		return fmt.Errorf("file selection failed: %w", err)
	}

	// Show processing screen
	outputPath, err := fileInput.ShowProcessingScreen(filePath)
	if err != nil {
		return fmt.Errorf("error showing processing screen: %w", err)
	}

	// TODO: Implement actual file conversion
	// For now, simulate a successful conversion
	// time.Sleep(2 * time.Second) // Simulate processing time

	// Show success screen
	if err := fileInput.ShowSuccessScreen(filePath, outputPath); err != nil {
		return fmt.Errorf("error showing success screen: %w", err)
	}

	return nil
}

// handleDirectory handles the directory conversion option
func (s *Screen) handleDirectory() error {
	s.ShowMessage("Directory conversion selected")
	// TODO: Implement directory conversion
	return nil
}

// handleSettings handles the settings menu
func (s *Screen) handleSettings() error {
	s.ShowMessage("Settings menu selected")
	// TODO: Implement settings
	return nil
}

// handleExit handles the exit option
func (s *Screen) handleExit() error {
	s.ShowMessage("Thank you for using HEIC-2-Go!")
	os.Exit(0)
	return nil
}

// ShowMessage displays a message to the user
func (s *Screen) ShowMessage(message string) {
	fmt.Println("\n" + message)
	fmt.Println("Press Enter to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

// GetInput prompts the user for input and returns the result
func (s *Screen) GetInput(prompt string) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(input), nil
}

// GetIntInput prompts the user for an integer input and returns the result
func (s *Screen) GetIntInput(prompt string, min, max int) (int, error) {
	for {
		input, err := s.GetInput(prompt)
		if err != nil {
			return 0, err
		}

		value, err := strconv.Atoi(input)
		if err != nil {
			fmt.Printf("Please enter a valid number between %d and %d\n", min, max)
			continue
		}

		if value < min || value > max {
			fmt.Printf("Please enter a number between %d and %d\n", min, max)
			continue
		}

		return value, nil
	}
}
