package ui

import (
    "github.com/yourusername/branch-wrangler/pkg/config"
    "github.com/rivo/tview"
)

// Run starts the TUI
func Run(cfg *config.Config) error {
    app := tview.NewApplication()
    // TODO: build your layout, list of branches, keybindings...
    return app.Run()
}