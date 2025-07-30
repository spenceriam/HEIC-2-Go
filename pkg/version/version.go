// Package version defines the version information for HEIC-2-Go
package version

import "fmt"

// These values are set during build time using -ldflags
var (
	// Version is the semantic version of the application
	Version = "0.1.0"

	// Commit is the git commit hash
	Commit = "dev"

	// BuildTime is the build timestamp
	BuildTime = ""
)

// String returns a formatted version string
func String() string {
	return fmt.Sprintf("v%s (commit: %s, built: %s)", Version, Commit, BuildTime)
}
