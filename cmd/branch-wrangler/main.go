package main

import (
	"context"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/dfinster/branch-wrangler/internal/config"
	"github.com/dfinster/branch-wrangler/internal/git"
	"github.com/dfinster/branch-wrangler/internal/github"
	"github.com/dfinster/branch-wrangler/internal/ui"
)

var rootCmd = &cobra.Command{
	Use:   "branch-wrangler",
	Short: "A cross-platform TUI for managing local Git branches",
	Long:  `Branch Wrangler helps manage local Git branches by reconciling them with GitHub.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := runTUI(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.Flags().Bool("version", false, "Show version information")
	rootCmd.Flags().Bool("list", false, "List branches in headless mode")
	rootCmd.Flags().Bool("json", false, "Output in JSON format")
	rootCmd.Flags().Bool("delete-stale", false, "Delete stale branches")
	rootCmd.Flags().Bool("dry-run", false, "Show what would be deleted without doing it")
	rootCmd.Flags().Bool("login", false, "Force interactive authentication")
	rootCmd.Flags().Bool("logout", false, "Clear stored authentication token")
	rootCmd.Flags().String("config", "", "Override default config file location")
	rootCmd.Flags().String("github-token-path", "", "Override default token location")
	rootCmd.Flags().StringSlice("base-branches", []string{"main", "master", "develop"}, "Override default base branches")
	rootCmd.Flags().String("completion", "", "Generate shell completion (bash|zsh|fish)")
}

func runTUI() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}
	
	gitClient := git.NewClient(cwd)
	if !gitClient.IsGitRepo() {
		return fmt.Errorf("not a git repository")
	}
	
	remoteURL, err := gitClient.GetRemoteURL()
	if err != nil {
		return fmt.Errorf("failed to get remote URL: %w", err)
	}
	
	owner, repo, err := gitClient.ParseGitHubRepo(remoteURL)
	if err != nil {
		return fmt.Errorf("not a GitHub repository: %w", err)
	}
	
	githubClient, err := github.NewCachedClient(owner, repo)
	if err != nil {
		return fmt.Errorf("failed to create GitHub client: %w", err)
	}
	
	classifier := git.NewClassifier(gitClient, githubClient, cfg.BaseBranches)
	
	ctx := context.Background()
	model := ui.NewModel(ctx, classifier)
	
	p := tea.NewProgram(model, tea.WithAltScreen())
	_, err = p.Run()
	return err
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}