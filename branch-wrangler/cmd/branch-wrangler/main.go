package main

import (
    "flag"
    "fmt"
    "os"

    "github.com/yourusername/branch-wrangler/pkg/config"
    "github.com/yourusername/branch-wrangler/internal/ui"
)

func main() {
    // parse flags
    cfg := config.New()
    flag.BoolVar(&cfg.DryRun, "dry-run", true, "show branches to delete without deleting")
    flag.Parse()

    // kick off TUI
    if err := ui.Run(cfg); err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v
", err)
        os.Exit(1)
    }
}