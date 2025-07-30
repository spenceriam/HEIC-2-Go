# HEIC-2-Go

A cross-platform command-line tool for converting HEIC images to JPG format with metadata preservation.

## Features

- 🖼️ Converts HEIC files to high-quality JPG
- 📊 Preserves all EXIF metadata (GPS, timestamps, camera settings)
- 🖥️ Beautiful ASCII art interface
- ⚡ Fast batch processing
- 🛡️ Handles file conflicts gracefully
- 🔒 Cross-platform admin permission handling

## Installation

1. Install [Go](https://golang.org/dl/) (version 1.16 or later)
2. Clone this repository:
   ```bash
   git clone https://github.com/spenceriam/HEIC-2-Go.git
   cd HEIC-2-Go
   ```
3. Build the application:
   ```bash
   go build -o heic2go ./cmd/heic2go
   ```

## Usage

```bash
# Convert a single file
./heic2go convert image.heic

# Convert all HEIC files in a directory
./heic2go batch /path/to/directory

# Show help
./heic2go --help
```

## Project Structure

```
.
├── cmd/                # Main application entry points
│   └── heic2go/       
│       └── main.go    # Main application
├── internal/          # Private application code
│   ├── app/           # Application logic
│   ├── config/        # Configuration management
│   ├── converter/     # HEIC to JPG conversion
│   └── ui/            # Terminal user interface
├── pkg/               # Public libraries
│   └── version/       # Version information
├── scripts/           # Build and utility scripts
└── test/              # Test files
```

## Development

### Prerequisites

- Go 1.16+
- Git

### Building

```bash
# Build for current platform
go build -o bin/heic2go ./cmd/heic2go

# Cross-compile for all platforms
./scripts/build.sh
```

### Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -coverprofile=coverage.out ./...
```

## License

MIT
