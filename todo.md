# HEIC-2-Go Development Tasks

## Project Setup
- [x] Initialize git repository
- [ ] Create project structure
  - [ ] cmd/heic2go/main.go
  - [ ] internal/app/
  - [ ] internal/converter/
  - [ ] internal/ui/
  - [ ] internal/config/
  - [ ] pkg/version/
  - [ ] scripts/
  - [ ] test/
- [ ] Set up Go module
- [ ] Add required dependencies to go.mod

## Core Functionality
- [ ] Implement HEIC to JPG conversion
  - [ ] Basic file conversion
  - [ ] Metadata preservation (EXIF, GPS, etc.)
  - [ ] Quality preservation
  - [ ] Error handling

## User Interface
- [ ] ASCII Art Interface
  - [x] Welcome screen
  - [ ] Main menu
  - [ ] File selection screen
  - [ ] Processing screen
  - [ ] Success/Error screens
  - [ ] Settings menu

## File Operations
- [ ] Single file processing
  - [ ] File path input
  - [ ] File validation
  - [ ] Output file naming
  - [ ] Conflict resolution
- [ ] Batch directory processing
  - [ ] Directory scanning
  - [ ] Progress tracking
  - [ ] Batch conflict resolution
  - [ ] Summary reporting

## System Integration
- [ ] Admin permissions
  - [ ] Windows UAC handling
  - [ ] macOS sudo handling
  - [ ] Linux sudo handling
- [ ] File operations
  - [ ] Cross-platform file handling
  - [ ] File permissions
  - [ ] Temporary file management

## Settings & Configuration
- [ ] User preferences
  - [ ] Default output directory
  - [ ] JPG quality settings
  - [ ] Color theme selection
  - [ ] Configuration persistence

## Testing
- [ ] Unit tests
- [ ] Integration tests
- [ ] Cross-platform testing
- [ ] Test data preparation

## Build & Distribution
- [ ] Cross-compilation scripts
- [ ] Version management
- [ ] Release packaging
- [ ] Installation scripts

## Documentation
- [ ] README.md
- [ ] User guide
- [ ] Developer documentation
- [ ] Command-line help

## Current Status
- Initial project setup in progress
- Git repository initialized
- PRD and user flow documents reviewed
- Todo list created

## Next Steps
1. Set up Go module and project structure
2. Implement basic file conversion functionality
3. Create initial UI components
