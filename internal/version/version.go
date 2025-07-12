package version

import (
	"encoding/json"
	"fmt"
	"runtime"
	"time"
)

// These variables are set via -ldflags at build time
var (
	// Version is the semantic version of the application
	Version = "dev"
	// BuildDate is the date/time when the binary was built
	BuildDate = "unknown"
	// CommitHash is the git commit hash of the build
	CommitHash = "unknown"
	// GoVersion is the Go version used to build the binary
	GoVersion = runtime.Version()
)

// VersionInfo contains all version-related information
type VersionInfo struct {
	Version    string `json:"version"`
	BuildDate  string `json:"build_date"`
	CommitHash string `json:"commit_hash"`
	GoVersion  string `json:"go_version"`
	OS         string `json:"os"`
	Arch       string `json:"arch"`
}

// GetVersion returns the current version string
func GetVersion() string {
	return Version
}

// GetBuildDate returns the build date string
func GetBuildDate() string {
	return BuildDate
}

// GetCommitHash returns the commit hash string
func GetCommitHash() string {
	return CommitHash
}

// GetGoVersion returns the Go version used to build
func GetGoVersion() string {
	return GoVersion
}

// GetFullVersion returns a complete VersionInfo struct
func GetFullVersion() VersionInfo {
	return VersionInfo{
		Version:    Version,
		BuildDate:  BuildDate,
		CommitHash: CommitHash,
		GoVersion:  GoVersion,
		OS:         runtime.GOOS,
		Arch:       runtime.GOARCH,
	}
}

// String returns a human-readable version string
func (v VersionInfo) String() string {
	if BuildDate == "unknown" {
		return fmt.Sprintf("branch-wrangler %s (development build)", v.Version)
	}

	buildTime, err := time.Parse(time.RFC3339, BuildDate)
	if err != nil {
		return fmt.Sprintf("branch-wrangler %s\nBuild: %s\nCommit: %s\nGo: %s",
			v.Version, BuildDate, v.CommitHash, v.GoVersion)
	}

	return fmt.Sprintf("branch-wrangler %s\nBuild: %s\nCommit: %s\nGo: %s",
		v.Version, buildTime.Format("2006-01-02 15:04:05 MST"), v.CommitHash, v.GoVersion)
}

// JSON returns the version information as JSON
func (v VersionInfo) JSON() (string, error) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// IsDevelopmentBuild returns true if this is a development build
func IsDevelopmentBuild() bool {
	return Version == "dev" || BuildDate == "unknown"
}

// GetShortCommit returns the first 7 characters of the commit hash
func GetShortCommit() string {
	if len(CommitHash) >= 7 {
		return CommitHash[:7]
	}
	return CommitHash
}