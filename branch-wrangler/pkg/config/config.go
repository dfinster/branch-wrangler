package config

// Config holds CLI flags and environment settings
type Config struct {
    DryRun bool
}

// New returns default config
func New() *Config {
    return &Config{DryRun: true}
}