package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ConflictResolution represents the user's choice for handling file conflicts
type ConflictResolution int

const (
	// ConflictOverwrite indicates the user wants to overwrite the existing file
	ConflictOverwrite ConflictResolution = iota + 1
	// ConflictRename indicates the user wants to save with a new name
	ConflictRename
	// ConflictSkip indicates the user wants to skip this file
	ConflictSkip
)

// handleFileConflict handles file conflict resolution
func (f *FileInputScreen) handleFileConflict(outputPath string) (string, ConflictResolution, error) {
	// If the output file doesn't exist, no conflict
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		return outputPath, 0, nil
	}

	// File exists, show conflict resolution menu
	for {
		f.screen.Clear()
		f.screen.DisplayWelcome()

		// Display conflict message
		fmt.Printf("\n⚠️  File already exists: %s\n\n", filepath.Base(outputPath))
		fmt.Println("╔══════════════════════════════════════════════════════════════════════╗")
		fmt.Println("║                     File Already Exists                             ║")
		fmt.Println("╚══════════════════════════════════════════════════════════════════════╝")
		fmt.Println()
		fmt.Printf("The file '%s' already exists.\n\n", filepath.Base(outputPath))
		fmt.Println("How would you like to proceed?")
		fmt.Println("1. Overwrite the existing file")
		fmt.Println("2. Save with a different name")
		fmt.Println("3. Skip this file")
		fmt.Println("4. Cancel all")
		fmt.Print("\nEnter your choice (1-4): ")

		// Get user input
		var choice int
		_, err := fmt.Scanln(&choice)
		if err != nil {
			fmt.Println("Invalid input. Please enter a number between 1 and 4.")
			continue
		}

		switch choice {
		case 1: // Overwrite
			return outputPath, ConflictOverwrite, nil
		case 2: // Rename
			newPath, err := f.promptForNewFilename(outputPath)
			if err != nil {
				return "", 0, err
			}
			return newPath, ConflictRename, nil
		case 3: // Skip
			return "", ConflictSkip, nil
		case 4: // Cancel
			return "", 0, fmt.Errorf("operation cancelled by user")
		default:
			fmt.Println("Invalid choice. Please enter a number between 1 and 4.")
		}
	}
}

// promptForNewFilename prompts the user for a new filename
func (f *FileInputScreen) promptForNewFilename(originalPath string) (string, error) {
	for {
		f.screen.Clear()
		f.screen.DisplayWelcome()

		ext := filepath.Ext(originalPath)
		dir := filepath.Dir(originalPath)
		base := strings.TrimSuffix(filepath.Base(originalPath), ext)

		fmt.Printf("\nEnter new filename (without extension): %s_", base)
		fmt.Print("\n\nNew filename: ")

		var newName string
		_, err := fmt.Scanln(&newName)
		if err != nil {
			return "", err
		}

		// If empty, use the original base name
		if newName == "" {
			newName = base
		}

		// Create the new path
		newPath := filepath.Join(dir, newName+ext)

		// Check if the new name is the same as the original
		if newPath == originalPath {
			fmt.Println("New filename is the same as the original. Please choose a different name.")
			continue
		}

		// Check if the new path already exists
		if _, err := os.Stat(newPath); err == nil {
			fmt.Printf("File '%s' already exists. Please choose a different name.\n", filepath.Base(newPath))
			continue
		}

		return newPath, nil
	}
}

// prepareOutputPath handles file conflicts and returns the final output path
func (f *FileInputScreen) prepareOutputPath(inputPath, outputDir string) (string, error) {
	// If output directory is not specified, use the input file's directory
	if outputDir == "" {
		outputDir = filepath.Dir(inputPath)
	}

	// Create the output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create output directory: %w", err)
	}

	// Generate the output path
	ext := filepath.Ext(inputPath)
	base := strings.TrimSuffix(filepath.Base(inputPath), ext)
	outputPath := filepath.Join(outputDir, base+".jpg")

	// Check for conflicts and handle them
	if _, err := os.Stat(outputPath); err == nil {
		// File exists, handle conflict
		newPath, resolution, err := f.handleFileConflict(outputPath)
		if err != nil {
			return "", err
		}

		switch resolution {
		case ConflictOverwrite:
			// Remove the existing file
			if err := os.Remove(outputPath); err != nil {
				return "", fmt.Errorf("failed to remove existing file: %w", err)
			}
			return outputPath, nil
		case ConflictRename:
			return newPath, nil
		case ConflictSkip:
			return "", nil // Empty path indicates skip
		}
	}

	return outputPath, nil
}
