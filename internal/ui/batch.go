package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/spenceriam/HEIC-2-Go/internal/converter"
)

// BatchProcessor handles batch processing of directories
func (f *FileInputScreen) BatchProcessDirectory(inputDir, outputDir string) error {
	// Get all HEIC files in the directory
	files, err := findHEICFiles(inputDir)
	if err != nil {
		return fmt.Errorf("error finding HEIC files: %w", err)
	}

	if len(files) == 0 {
		return fmt.Errorf("no HEIC files found in directory")
	}

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Create a channel to track progress
	progressChan := make(chan int, len(files))
	fileChan := make(chan string, len(files))
	doneChan := make(chan bool)
	errorChan := make(chan error, 1)

	// Start the progress display
	go func() {
		f.showBatchProgress(inputDir, len(files), progressChan, fileChan, doneChan, errorChan)
	}()

	// Process files in parallel
	var wg sync.WaitGroup
	converter := converter.NewHEICConverter(true) // Preserve metadata

	// Start worker goroutines
	numWorkers := 4 // Number of concurrent conversions
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go f.batchWorker(&wg, converter, inputDir, outputDir, progressChan, fileChan, errorChan)
	}

	// Wait for all workers to finish
	wg.Wait()
	close(progressChan)
	doneChan <- true

	// Check for any errors
	select {
	case err := <-errorChan:
		return err
	default:
		return nil
	}
}

// batchWorker processes files from the queue
func (f *FileInputScreen) batchWorker(wg *sync.WaitGroup, conv *converter.HEICConverter, 
	inputDir, outputDir string, progressChan chan<- int, 
	fileChan chan<- string, errorChan chan<- error) {
	defer wg.Done()

	for file := range f.batchQueue {
		// Process the file
		outputFile := filepath.Join(outputDir, filepath.Base(file))
		outputFile = strings.TrimSuffix(outputFile, filepath.Ext(outputFile)) + ".jpg"

		// Convert the file
		if err := conv.Convert(file, outputFile); err != nil {
			errorChan <- fmt.Errorf("error converting %s: %w", file, err)
			continue
		}

		// Update progress
		progressChan <- 1
		fileChan <- filepath.Base(file)
	}
}

// showBatchProgress displays the batch processing progress
func (f *FileInputScreen) showBatchProgress(dir string, totalFiles int, 
	progressChan <-chan int, fileChan <-chan string, 
	doneChan <-chan bool, errorChan chan<- error) {
	
	// Create a progress bar
	progressBar := NewProgressBar(totalFiles)
	progressBar.startTime = time.Now()

	// Track processed files
	processed := 0
	var currentFile string

	for {
		select {
		case n := <-progressChan:
			processed += n
			progressBar.Update(processed)

		case file := <-fileChan:
			currentFile = file

		case <-doneChan:
			progressBar.Update(totalFiles)
			fmt.Printf("\n✅ Successfully processed %d files\n", totalFiles)
			return

		case err := <-errorChan:
			fmt.Printf("\n\n❌ Error: %v\n", err)
			errorChan <- err
			return

		case <-time.After(100 * time.Millisecond):
			// Update the display periodically
			if currentFile != "" {
				fmt.Printf("\rProcessing: %-50s", currentFile)
			}
		}
	}
}

// findHEICFiles finds all HEIC files in a directory
func findHEICFiles(dir string) ([]string, error) {
	var files []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Check if file is HEIC/HEIF
		ext := strings.ToLower(filepath.Ext(path))
		if ext == ".heic" || ext == ".heif" {
			files = append(files, path)
		}

		return nil
	})

	return files, err
}
