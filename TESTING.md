# HEIC-2-Go Test Plan

## Table of Contents
1. [Getting Started](#getting-started)
2. [Test Environment Setup](#test-environment-setup)
3. [Automated Testing](#automated-testing)
4. [Manual Testing](#manual-testing)
5. [Test Data Preparation](#test-data-preparation)
6. [Test Execution](#test-execution)
7. [Reporting & Analysis](#reporting--analysis)
8. [Continuous Integration](#continuous-integration)

## Getting Started

### Prerequisites
- Go 1.19 or later
- Git
- HEIC sample files for testing
- (Optional) Test automation tools (Go test, CI/CD tools)

### Quick Start for Testers
1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/HEIC-2-Go.git
   cd HEIC-2-Go
   ```
2. Install dependencies:
   ```bash
   go mod download
   ```
3. Run automated tests:
   ```bash
   go test ./... -v
   ```

## Test Environment Setup

### Supported Platforms
| OS Version | Architecture | Notes                      |
|------------|--------------|----------------------------|
| Windows 10/11 | x86_64      | Test admin privileges      |
| macOS 12+  | arm64/x86_64 | Test both architectures    |
| Ubuntu 20.04+ | x86_64     | Test with default packages |

### Environment Variables
```bash
# For detailed logging
export HEIC2GO_DEBUG=true

# To test admin functionality
export HEIC2GO_TEST_ADMIN=true
```

## Automated Testing

### 1. Unit Tests

#### 1.1 Converter Package
```bash
# Run all converter tests
go test -v ./internal/converter/...

# Run specific test with coverage
go test -v -run TestHEICToJPGConversion -coverprofile=coverage.out
```

#### 1.2 UI Package
```bash
# Test UI components
go test -v ./internal/ui/...
```

#### 1.3 Error Handling
```bash
# Test error scenarios
go test -v ./internal/errors/...
```

### 2. Integration Tests

#### 2.1 File Operations
```bash
# Test single file conversion with sample data
go test -v -run TestSingleFileConversion ./test/integration

# Test batch processing
go test -v -run TestBatchProcessing ./test/integration
```

## Manual Testing

### 1. Installation Verification

#### 1.1 Fresh Installation
1. **Windows**
   ```powershell
   # In PowerShell as Administrator
   .\scripts\build.ps1
   # Verify installation
   .\bin\heic2go --version
   ```

2. **macOS/Linux**
   ```bash
   chmod +x scripts/build.sh
   ./scripts/build.sh
   # Verify installation
   ./bin/heic2go --version
   ```

### 2. Basic Functionality

#### 2.1 Single File Conversion
1. Prepare a test HEIC file in `testdata/sample.heic`
2. Run the conversion:
   ```bash
   heic2go convert -i testdata/sample.heic -o output.jpg
   ```
3. Verify:
   - Output file exists
   - Image quality is preserved
   - EXIF metadata is intact

#### 2.2 Batch Processing
1. Create a test directory with HEIC files:
   ```bash
   mkdir -p testdata/batch
   # Add sample HEIC files to testdata/batch/
   ```
2. Run batch conversion:
   ```bash
   heic2go batch -i testdata/batch -o output_dir
   ```
3. Verify all files are converted with proper naming

### 3. Error Scenarios

#### 3.1 Invalid Input
```bash
# Non-existent file
heic2go convert -i missing.heic -o out.jpg
# Expected: Clear error message

# Unsupported format
heic2go convert -i test.txt -o out.jpg
# Expected: Format validation error
```

## Test Data Preparation

### Sample Files
1. **Basic Test Set**
   - Small HEIC (<1MB)
   - Medium HEIC (1-5MB)
   - Large HEIC (>5MB)
   - HEIC with GPS metadata
   - HEIC with orientation data

2. **Edge Cases**
   - Corrupted HEIC files
   - Files with special characters in names
   - Files in nested directories

### Test Directories
```
testdata/
  ├── basic/
  │   ├── image1.heic
  │   └── image2.heic
  ├── nested/
  │   └── folder1/
  │       └── image3.heic
  └── special_chars/
      └── "image with spaces.heic"
```

## Test Execution

### Automated Test Suite
```bash
# Run all tests with coverage
./scripts/run_tests.sh

# Generate HTML coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

### Manual Test Checklist

#### Installation
- [ ] Verify installation on all supported platforms
- [ ] Test upgrade from previous versions
- [ ] Verify uninstallation leaves no artifacts

#### Functionality
- [ ] Single file conversion
- [ ] Batch processing
- [ ] Progress reporting
- [ ] Error handling
- [ ] Settings persistence

## Reporting & Analysis

### Test Results
1. **Automated Tests**
   - View console output for test results
   - Check `test-results/` for detailed reports
   - Review coverage reports in `coverage.html`

2. **Manual Testing**
   - Document test cases in `TEST_RESULTS.md`
   - Include screenshots for UI issues
   - Note any performance observations

### Bug Reporting
When reporting issues, include:
1. Test case description
2. Steps to reproduce
3. Expected vs actual results
4. Environment details
5. Logs or screenshots

## Continuous Integration

### GitHub Actions
Workflow files are in `.github/workflows/`:
- `test.yml`: Runs on push/pull requests
- `release.yml`: Creates releases on version tags

### Local CI Verification
```bash
# Run the same checks as CI
./scripts/ci-check.sh
```

## Performance Testing

### Benchmark Tests
```bash
# Run performance benchmarks
go test -bench=. -benchmem ./internal/benchmarks
```

### Memory Profiling
```bash
# Generate memory profile
go test -memprofile=mem.out ./...
# Analyze with pprof
go tool pprof -http=:8080 mem.out
```

## Security Testing

### Input Validation
- [ ] Test with malformed HEIC files
- [ ] Test path traversal attempts
- [ ] Test with extremely large files

### Permission Testing
- [ ] Run with non-admin privileges
- [ ] Test with read-only directories
- [ ] Verify proper cleanup of temp files

## Accessibility

### Keyboard Navigation
- [ ] Tab through all interactive elements
- [ ] Test with screen readers
- [ ] Verify focus indicators are visible

### Color Contrast
- [ ] Verify text is readable in both light/dark themes
- [ ] Check colorblind-friendly color schemes
