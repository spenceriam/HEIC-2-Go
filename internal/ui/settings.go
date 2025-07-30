package ui

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

// Settings holds the application settings
type Settings struct {
	// Image quality (1-100)
	Quality int `json:"quality"`
	// Default output directory
	OutputDir string `json:"output_dir"`
	// Color theme (light/dark)
	Theme string `json:"theme"`
	// Whether to preserve EXIF metadata
	PreserveMetadata bool `json:"preserve_metadata"`
}

// DefaultSettings returns the default application settings
func DefaultSettings() *Settings {
	// Get user's home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}

	return &Settings{
		Quality:          90,
		OutputDir:        filepath.Join(homeDir, "Pictures", "HEIC-2-JPG"),
		Theme:            "dark",
		PreserveMetadata: true,
	}
}

// ShowSettingsMenu displays the settings menu
func (f *FileInputScreen) ShowSettingsMenu() error {
	for {
		f.screen.Clear()
		f.screen.DisplayWelcome()

		// Display current settings
		fmt.Println("\n╔══════════════════════════════════════════════════════════════════════╗")
		fmt.Println("║                         Settings                              ║")
		fmt.Println("╚══════════════════════════════════════════════════════════════════════╝")
		fmt.Println()

		// Show current settings
		fmt.Printf("1. Image Quality: %d%%\n", f.settings.Quality)
		fmt.Printf("2. Output Directory: %s\n", f.settings.OutputDir)
		fmt.Printf("3. Theme: %s\n", strings.Title(f.settings.Theme))
		fmt.Printf("4. Preserve Metadata: %v\n", f.settings.PreserveMetadata)
		fmt.Println("5. Reset to Defaults")
		fmt.Println("6. Back to Main Menu")
		fmt.Print("\nSelect an option (1-6): ")

		// Get user input
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "1":
			f.updateQualitySetting()
		case "2":
			f.updateOutputDirectory()
		case "3":
			f.toggleTheme()
		case "4":
			f.toggleMetadataPreservation()
		case "5":
			f.resetToDefaults()
		case "6":
			return nil
		default:
			fmt.Println("\nInvalid option. Please try again.")
			fmt.Print("Press Enter to continue...")
			reader.ReadString('\n')
		}
	}
}

// updateQualitySetting allows the user to update the image quality setting
func (f *FileInputScreen) updateQualitySetting() {
	f.screen.Clear()
	f.screen.DisplayWelcome()

	fmt.Println("\n╔══════════════════════════════════════════════════════════════════════╗")
	fmt.Println("║                     Image Quality                             ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════════════╝")
	fmt.Println()

	for {
		fmt.Printf("Current quality: %d%%\n", f.settings.Quality)
		fmt.Print("Enter new quality (1-100, 90 recommended): ")

		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "" {
			return
		}

		quality, err := strconv.Atoi(input)
		if err != nil || quality < 1 || quality > 100 {
			fmt.Println("Invalid input. Please enter a number between 1 and 100.")
			continue
		}

		f.settings.Quality = quality
		fmt.Println("\nQuality setting updated successfully!")
		fmt.Print("Press Enter to continue...")
		reader.ReadString('\n')
		return
	}
}

// updateOutputDirectory allows the user to update the output directory
func (f *FileInputScreen) updateOutputDirectory() {
	f.screen.Clear()
	f.screen.DisplayWelcome()

	fmt.Println("\n╔══════════════════════════════════════════════════════════════════════╗")
	fmt.Println("║                   Output Directory                           ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════════════╝")
	fmt.Println()

	fmt.Printf("Current output directory: %s\n\n", f.settings.OutputDir)
	fmt.Println("Enter new output directory path (or press Enter to cancel):")
	fmt.Print("> ")

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "" {
		return
	}

	// Validate the directory
	if info, err := os.Stat(input); err == nil && info.IsDir() {
		f.settings.OutputDir = input
		fmt.Println("\nOutput directory updated successfully!")
	} else {
		fmt.Println("\nError: The specified path is not a valid directory.")
	}

	fmt.Print("Press Enter to continue...")
	reader.ReadString('\n')
}

// toggleTheme toggles between light and dark themes
func (f *FileInputScreen) toggleTheme() {
	if f.settings.Theme == "dark" {
		f.settings.Theme = "light"
		// Update color scheme for light theme
		color.NoColor = false
		color.Set(color.FgBlack)
	} else {
		f.settings.Theme = "dark"
		// Update color scheme for dark theme
		color.NoColor = false
		color.Set(color.FgWhite)
	}
}

// toggleMetadataPreservation toggles the metadata preservation setting
func (f *FileInputScreen) toggleMetadataPreservation() {
	f.settings.PreserveMetadata = !f.settings.PreserveMetadata
	status := "enabled"
	if !f.settings.PreserveMetadata {
		status = "disabled"
	}
	fmt.Printf("\nMetadata preservation has been %s.\n", status)
	fmt.Print("Press Enter to continue...")
	bufio.NewReader(os.Stdin).ReadString('\n')
}

// resetToDefaults resets all settings to their default values
func (f *FileInputScreen) resetToDefaults() {
	defaultSettings := DefaultSettings()
	f.settings = defaultSettings
	fmt.Println("\nAll settings have been reset to their default values.")
	fmt.Print("Press Enter to continue...")
	bufio.NewReader(os.Stdin).ReadString('\n')
}
