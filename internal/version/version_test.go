package version

import (
	"encoding/json"
	"runtime"
	"strings"
	"testing"
)

func TestGetVersion(t *testing.T) {
	version := GetVersion()
	if version == "" {
		t.Error("GetVersion() returned empty string")
	}
}

func TestGetBuildDate(t *testing.T) {
	buildDate := GetBuildDate()
	if buildDate == "" {
		t.Error("GetBuildDate() returned empty string")
	}
}

func TestGetCommitHash(t *testing.T) {
	commitHash := GetCommitHash()
	if commitHash == "" {
		t.Error("GetCommitHash() returned empty string")
	}
}

func TestGetGoVersion(t *testing.T) {
	goVersion := GetGoVersion()
	if goVersion == "" {
		t.Error("GetGoVersion() returned empty string")
	}
	// Should match runtime.Version()
	if goVersion != runtime.Version() {
		t.Errorf("GetGoVersion() = %v, want %v", goVersion, runtime.Version())
	}
}

func TestGetFullVersion(t *testing.T) {
	versionInfo := GetFullVersion()

	if versionInfo.Version == "" {
		t.Error("VersionInfo.Version is empty")
	}
	if versionInfo.BuildDate == "" {
		t.Error("VersionInfo.BuildDate is empty")
	}
	if versionInfo.CommitHash == "" {
		t.Error("VersionInfo.CommitHash is empty")
	}
	if versionInfo.GoVersion == "" {
		t.Error("VersionInfo.GoVersion is empty")
	}
	if versionInfo.OS == "" {
		t.Error("VersionInfo.OS is empty")
	}
	if versionInfo.Arch == "" {
		t.Error("VersionInfo.Arch is empty")
	}

	// Check that OS and Arch match runtime values
	if versionInfo.OS != runtime.GOOS {
		t.Errorf("VersionInfo.OS = %v, want %v", versionInfo.OS, runtime.GOOS)
	}
	if versionInfo.Arch != runtime.GOARCH {
		t.Errorf("VersionInfo.Arch = %v, want %v", versionInfo.Arch, runtime.GOARCH)
	}
}

func TestVersionInfoString(t *testing.T) {
	versionInfo := GetFullVersion()
	str := versionInfo.String()

	if str == "" {
		t.Error("VersionInfo.String() returned empty string")
	}

	// Should contain the application name
	if !strings.Contains(str, "branch-wrangler") {
		t.Error("VersionInfo.String() should contain 'branch-wrangler'")
	}

	// Should contain version
	if !strings.Contains(str, versionInfo.Version) {
		t.Error("VersionInfo.String() should contain version")
	}
}

func TestVersionInfoStringDevelopmentBuild(t *testing.T) {
	// Save original values
	originalVersion := Version
	originalBuildDate := BuildDate

	// Set to development values
	Version = "dev"
	BuildDate = "unknown"

	defer func() {
		// Restore original values
		Version = originalVersion
		BuildDate = originalBuildDate
	}()

	versionInfo := GetFullVersion()
	str := versionInfo.String()

	if !strings.Contains(str, "development build") {
		t.Error("Development build should be indicated in version string")
	}
}

func TestVersionInfoJSON(t *testing.T) {
	versionInfo := GetFullVersion()
	jsonStr, err := versionInfo.JSON()

	if err != nil {
		t.Errorf("VersionInfo.JSON() returned error: %v", err)
	}

	if jsonStr == "" {
		t.Error("VersionInfo.JSON() returned empty string")
	}

	// Should be valid JSON
	var parsed VersionInfo
	if err := json.Unmarshal([]byte(jsonStr), &parsed); err != nil {
		t.Errorf("VersionInfo.JSON() produced invalid JSON: %v", err)
	}

	// Parsed should match original
	if parsed.Version != versionInfo.Version {
		t.Error("Parsed JSON version doesn't match original")
	}
	if parsed.OS != versionInfo.OS {
		t.Error("Parsed JSON OS doesn't match original")
	}
}

func TestIsDevelopmentBuild(t *testing.T) {
	// Save original values
	originalVersion := Version
	originalBuildDate := BuildDate

	defer func() {
		// Restore original values
		Version = originalVersion
		BuildDate = originalBuildDate
	}()

	// Test development build detection
	Version = "dev"
	BuildDate = "unknown"
	if !IsDevelopmentBuild() {
		t.Error("Should detect development build when Version='dev' and BuildDate='unknown'")
	}

	// Test with dev version but known build date
	Version = "dev"
	BuildDate = "2023-01-01T00:00:00Z"
	if !IsDevelopmentBuild() {
		t.Error("Should detect development build when Version='dev'")
	}

	// Test with release version but unknown build date
	Version = "v1.0.0"
	BuildDate = "unknown"
	if !IsDevelopmentBuild() {
		t.Error("Should detect development build when BuildDate='unknown'")
	}

	// Test release build
	Version = "v1.0.0"
	BuildDate = "2023-01-01T00:00:00Z"
	if IsDevelopmentBuild() {
		t.Error("Should not detect development build for release version with known build date")
	}
}

func TestGetShortCommit(t *testing.T) {
	// Save original value
	originalCommitHash := CommitHash

	defer func() {
		// Restore original value
		CommitHash = originalCommitHash
	}()

	// Test with long commit hash
	CommitHash = "abcdef1234567890abcdef1234567890abcdef12"
	short := GetShortCommit()
	if short != "abcdef1" {
		t.Errorf("GetShortCommit() = %v, want %v", short, "abcdef1")
	}

	// Test with short commit hash
	CommitHash = "abc123"
	short = GetShortCommit()
	if short != "abc123" {
		t.Errorf("GetShortCommit() = %v, want %v", short, "abc123")
	}

	// Test with empty commit hash
	CommitHash = ""
	short = GetShortCommit()
	if short != "" {
		t.Errorf("GetShortCommit() = %v, want empty string", short)
	}
}
