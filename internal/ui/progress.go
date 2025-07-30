package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
)

// ProgressBar represents a progress bar in the terminal
type ProgressBar struct {
	total     int
	current   int
	width     int
	startTime time.Time
}

// NewProgressBar creates a new progress bar
func NewProgressBar(total int) *ProgressBar {
	return &ProgressBar{
		total:     total,
		current:   0,
		width:     50, // Width of the progress bar in characters
		startTime: time.Now(),
	}
}

// Update updates the progress bar with the current progress
func (p *ProgressBar) Update(current int) {
	p.current = current
	p.Render()
}

// Increment increments the progress by 1
func (p *ProgressBar) Increment() {
	if p.current < p.total {
		p.current++
	}
	p.Render()
}

// Render renders the progress bar to the terminal
func (p *ProgressBar) Render() {
	// Calculate percentage
	percent := float64(p.current) / float64(p.total) * 100

	// Calculate the number of filled and empty segments
	filled := int(float64(p.width) * (percent / 100))
	if filled > p.width {
		filled = p.width
	}
	empty := p.width - filled

	// Calculate elapsed time and estimated time remaining
	elapsed := time.Since(p.startTime)
	var remaining time.Duration
	if p.current > 0 {
		totalEstimated := time.Duration(float64(elapsed) * (float64(p.total) / float64(p.current)))
		remaining = totalEstimated - elapsed
	}

	// Create the progress bar string
	progressBar := fmt.Sprintf("\r[%s%s] %3.0f%%  %d/%d  Elapsed: %s  Remaining: %s",
		strings.Repeat("█", filled),
		strings.Repeat(" ", empty),
		percent,
		p.current,
		p.total,
		formatDuration(elapsed),
		formatDuration(remaining),
	)

	// Print the progress bar
	fmt.Print(progressBar)

	// If we're done, print a newline
	if p.current >= p.total {
		fmt.Println()
	}
}

// formatDuration formats a duration in a human-readable format
func formatDuration(d time.Duration) string {
	d = d.Round(time.Second)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second

	if h > 0 {
		return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
	}
	return fmt.Sprintf("%02d:%02d", m, s)
}

// ShowProgressScreen displays a progress screen during file conversion
func (f *FileInputScreen) ShowProgressScreen(filePath string, progressChan <-chan int, doneChan <-chan bool, errorChan <-chan error) error {
	f.screen.Clear()
	f.screen.DisplayWelcome()

	// Display processing message
	fileName := filepath.Base(filePath)
	fmt.Printf("\nProcessing file: %s\n\n", fileName)
	fmt.Println("╔══════════════════════════════════════════════════════════════════════╗")
	fmt.Println("║                         Processing File                            ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════════════╝")
	fmt.Println()

	// Create a progress bar
	progressBar := NewProgressBar(100)

	// Start a goroutine to update the progress
	go func() {
		for {
			select {
			case progress, ok := <-progressChan:
				if !ok {
					return
				}
				progressBar.Update(progress)

			case <-doneChan:
				progressBar.Update(100)
				return

			case err := <-errorChan:
				if err != nil {
					f.screen.ShowError(fmt.Sprintf("Error during conversion: %v", err))
					return
				}
			}
		}
	}()

	// Wait for the conversion to complete
	<-doneChan
	return nil
}
