package converter

import (
	"bytes"
	"encoding/binary"
	"errors"
	"os"
)

// HEIC signature (file signature for HEIF/HEIC files)
var (
	heicSignature = []byte{'f', 't', 'y', 'p'}
	heicBrands    = [][]byte{
		{'h', 'e', 'i', 'c'}, // HEIC
		{'h', 'e', 'i', 'x'}, // HEIF still image
		{'m', 'i', 'f', '1'}, // HEIF image sequence
		{'m', 's', 'f', '1'}, // HEIF image sequence
	}
)

// IsValidHEIC checks if the given file is a valid HEIC/HEIF file
func IsValidHEIC(filePath string) (bool, error) {
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return false, err
	}
	defer file.Close()

	// Read the first 16 bytes (should be enough to identify the file type)
	buffer := make([]byte, 16)
	n, err := file.Read(buffer)
	if err != nil || n < 16 {
		return false, errors.New("file too small or unreadable")
	}

	// Check for the 'ftyp' signature at the start of the file
	if !bytes.Equal(buffer[4:8], heicSignature) {
		return false, errors.New("not a valid HEIC/HEIF file (missing 'ftyp' signature)")
	}

	// Check for known HEIC/HEIF brands
	for i := 0; i < len(heicBrands); i++ {
		// Brands are typically at offset 8-11 or 12-15 in the file
		if bytes.Equal(buffer[8:12], heicBrands[i]) || 
		   (len(buffer) > 15 && bytes.Equal(buffer[12:16], heicBrands[i])) {
			return true, nil
		}
	}

	return false, errors.New("not a supported HEIC/HEIF variant")
}

// GetFileSize returns the size of the file in bytes
func GetFileSize(filePath string) (int64, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return 0, err
	}
	return fileInfo.Size(), nil
}

// ReadFileHeader reads the first n bytes from a file
func ReadFileHeader(filePath string, n int) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	header := make([]byte, n)
	_, err = file.Read(header)
	if err != nil {
		return nil, err
	}

	return header, nil
}

// IsBigEndian checks if the system is big endian
func IsBigEndian() bool {
	var i int32 = 0x01020304
	u := (*[4]byte)(&i)
	return u[0] == 0x01
}

// ReadUint32 reads a 32-bit unsigned integer from a byte slice
func ReadUint32(data []byte, offset int) uint32 {
	if IsBigEndian() {
		return binary.BigEndian.Uint32(data[offset : offset+4])
	}
	return binary.LittleEndian.Uint32(data[offset : offset+4])
}
