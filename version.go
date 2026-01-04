package main

import "fmt"

// Default build-time variable.
// These values are overridden via ldflags
var (
	Version   = "unknown-version"
	GitCommit = "unknown-commit"
	BuildDate = "unknown-buildtime"
)

func FormatVersion() string {
	return fmt.Sprintf("nmkill version %s, build %s (%s)", Version, GitCommit, BuildDate)
}
