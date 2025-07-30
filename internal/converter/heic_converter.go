package converter

import (
	"fmt"
	"image"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/rwcarlsen/goexif/exif"
	heif "github.com/strukturag/libheif/go/heif"
)

// HEICConverter handles the conversion from HEIC to JPG
type HEICConverter struct {
	preserveMetadata bool
}

// NewHEICConverter creates a new HEICConverter instance
func NewHEICConverter(preserveMetadata bool) *HEICConverter {
	return &HEICConverter{
		preserveMetadata: preserveMetadata,
	}
}

// Convert converts a HEIC file to JPG format
func (c *HEICConverter) Convert(inputPath, outputPath string) error {
	// Validate input file
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		return fmt.Errorf("input file does not exist: %s", inputPath)
	}

	// Ensure output directory exists
	outputDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Read HEIC file
	img, metadata, err := c.decodeHEIC(inputPath)
	if err != nil {
		return fmt.Errorf("failed to decode HEIC file: %w", err)
	}

	// Save as JPG
	if err := imaging.Save(img, outputPath, imaging.JPEGQuality(90)); err != nil {
		return fmt.Errorf("failed to save JPG file: %w", err)
	}

	// Preserve metadata if requested
	if c.preserveMetadata && metadata != nil {
		if err := c.writeMetadata(outputPath, metadata); err != nil {
			// Don't fail the entire conversion if metadata can't be written
			// Just log the error and continue
			fmt.Printf("Warning: Failed to write metadata: %v\n", err)
		}
	}

	return nil
}

// decodeHEIC decodes a HEIC file and returns the image and metadata
func (c *HEICConverter) decodeHEIC(path string) (image.Image, *exif.Exif, error) {
	// Open the HEIC file
	file, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	// Get file info for size
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, nil, err
	}

	// Read the entire file into memory
	data := make([]byte, fileInfo.Size())
	if _, err := file.Read(data); err != nil {
		return nil, nil, err
	}

	// Create HEIF context
	ctx := heif.NewContext()
	if err := ctx.ReadFromMemory(data); err != nil {
		return nil, nil, err
	}

	// Get the primary image
	handle, err := ctx.GetPrimaryImageHandle()
	if err != nil {
		return nil, nil, err
	}

	// Decode the image
	img, err := handle.DecodeImage(heif.ColorspaceUndefined, heif.ChromaUndefined, nil)
	if err != nil {
		return nil, nil, err
	}

	// Extract metadata if needed
	var exifData *exif.Exif
	if c.preserveMetadata {
		exifData, _ = c.extractExifMetadata(handle)
	}

	return img, exifData, nil
}

// extractExifMetadata extracts EXIF metadata from HEIC image
func (c *HEICConverter) extractExifMetadata(handle *heif.ImageHandle) (*exif.Exif, error) {
	// Get EXIF data
	exifData, err := handle.GetExif()
	if err != nil {
		return nil, err
	}

	// Parse EXIF data
	return exif.Decode(exifData.Reader())
}

// writeMetadata writes metadata to the output file
func (c *HEICConverter) writeMetadata(path string, exifData *exif.Exif) error {
	// Open the output file for reading and writing
	file, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Get the current file content
	imgData, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	// Create a new buffer with the EXIF data
	var buf []byte
	if exifData != nil {
		buf, err = exifData.MarshalJSON()
		if err != nil {
			return err
		}
	}

	// For now, just print the metadata since we need additional libraries
	// to properly write EXIF to JPEG files
	fmt.Printf("Metadata to be written (not implemented yet): %s\n", string(buf))

	return nil
}

// GetOutputPath generates an output path for the converted file
func (c *HEICConverter) GetOutputPath(inputPath string) string {
	ext := filepath.Ext(inputPath)
	base := strings.TrimSuffix(inputPath, ext)
	return base + ".jpg"
}
