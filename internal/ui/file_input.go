package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

// FileInputScreen handles the file input interface
type FileInputScreen struct {
	screen *Screen
}

// NewFileInputScreen creates a new file input screen
func NewFileInputScreen(screen *Screen) *FileInputScreen {
	return &FileInputScreen{
		screen: screen,
	}
}

// Show displays the file input screen
func (f *FileInputScreen) Show() (string, error) {
	// Create admin manager
	adminMgr := app.NewAdminManager()

	// Check if we have admin privileges
	isAdmin, err := adminMgr.IsAdmin()
	if err != nil {
		return "", fmt.Errorf("error checking admin status: %w", err)
	}

	// If not admin, request elevation
	if !isAdmin {
		f.screen.Clear()
		f.screen.DisplayWelcome()
		
		fmt.Println("\n🔒 Admin privileges required")
		fmt.Println("This application needs admin/root privileges to access image files.")
		fmt.Println("You may be prompted to enter your password.")
		fmt.Print("\nPress Enter to continue...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')

		if err := adminMgr.RequestAdmin(); err != nil {
			return "", fmt.Errorf("failed to get admin privileges: %w", err)
		}
		
		// If we get here, the user canceled the admin prompt
		return "", nil
	}

	// Main file input loop
	for {
		f.screen.Clear()
		f.screen.DisplayWelcome()

		// Display the file input prompt
		prompt := `
╔══════════════════════════════════════════════════════════════════════╗
║                            Convert Single File                       ║
╚══════════════════════════════════════════════════════════════════════╝

┌─────────────────────────────────────────────────────────────────────┐
│  📁 Enter file path or drag & drop:                                 │
│  ________________________________________________________________  │
│                                                                     │
│  💡 Tip: You can also type 'browse' to open file picker            │
└─────────────────────────────────────────────────────────────────────┘

File path: `

		// Get user input
		input, err := f.screen.GetInput(prompt)
		if err != nil {
			return "", fmt.Errorf("error getting input: %w", err)
		}

		// Handle browse command
		if strings.EqualFold(input, "browse") {
			// TODO: Implement file picker dialog
			f.screen.ShowMessage("File picker will be implemented in a future version.")
			continue
		}

		// Clean up the input (handle drag and drop quotes on Windows)
		input = strings.TrimSpace(input)
		input = strings.Trim(input, `"'`)

		// Validate the file
		if input == "" {
			f.screen.ShowError("Please enter a file path")
			continue
		}

		// Check if file exists and is a regular file
		fileInfo, err := os.Stat(input)
		if err != nil {
			if os.IsNotExist(err) {
				f.screen.ShowError("File does not exist")
			} else {
				f.screen.ShowError(fmt.Sprintf("Error accessing file: %v", err))
			}
			continue
		}

		// Check if it's a directory
		if fileInfo.IsDir() {
			f.screen.ShowError("Path is a directory, please select a file")
			continue
		}

		// Check file extension first for quick validation
	ext := strings.ToLower(filepath.Ext(input))
		if ext != ".heic" && ext != ".heif" {
			f.screen.ShowError("File must have a .heic or .heif extension")
			continue
		}

		// Validate that it's actually a HEIC/HEIF file
		isValid, err := converter.IsValidHEIC(input)
		if err != nil {
			f.screen.ShowError(fmt.Sprintf("Error validating HEIC file: %v", err))
			continue
		}

		if !isValid {
			f.screen.ShowError("The selected file is not a valid HEIC/HEIF file")
			continue
		}

		// Return the validated file path
		return input, nil
	}
}

// ShowProcessingScreen displays the file processing screen
func (f *FileInputScreen) ShowProcessingScreen(filePath string) (string, error) {
	// Extract the file name for display
	fileName := filepath.Base(filePath)

	// Generate output path (same directory, change extension to .jpg)
	outputPath := strings.TrimSuffix(filePath, filepath.Ext(filePath)) + ".jpg"

	// Display processing screen
	f.screen.Clear()
	f.screen.DisplayWelcome()

	// Show processing message
	fmt.Printf("\nProcessing file: %s\n\n", fileName)
	fmt.Println("╔══════════════════════════════════════════════════════════════════════╗")
	fmt.Println("║                         Processing File                            ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════════════╝")
	fmt.Println()

	// Show progress bar placeholder
	progressBar := "[                                                  ] 0%"
	fmt.Println(progressBar)

	// Return the output path where the file will be saved
	return outputPath, nil
}

// ShowSuccessScreen displays the success screen after conversion
func (f *FileInputScreen) ShowSuccessScreen(inputPath, outputPath string) error {
	// Get file sizes for display
	inputSize, err := getFileSize(inputPath)
	if err != nil {
		return err
	}

	outputSize, err := getFileSize(outputPath)
	if err != nil {
		return err
	}

	// Display success message
	f.screen.Clear()
	f.screen.DisplayWelcome()

	fmt.Printf("\n✅ Conversion successful!\n\n")
	fmt.Println("╔══════════════════════════════════════════════════════════════════════╗")
	fmt.Println("║                       Conversion Complete                           ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════════════╝")
	fmt.Println()

	// Show file information
	fmt.Printf("📄 Original: %s (%.2f MB)\n", filepath.Base(inputPath), float64(inputSize)/(1024*1024))
	fmt.Printf("💾 Saved as: %s (%.2f MB)\n\n", filepath.Base(outputPath), float64(outputSize)/(1024*1024))

	// Show success message with green color
	successMsg := "File converted successfully!"
	color.Green(successMsg)
	fmt.Println()

	// Wait for user to continue
	fmt.Println("Press Enter to return to the main menu...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

	return nil
}

// getFileSize returns the size of a file in bytes
func getFileSize(filePath string) (int64, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return 0, err
	}
	return fileInfo.Size(), nil
}
