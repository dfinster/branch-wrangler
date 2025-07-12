package ui

import (
	"fmt"
	"os/exec"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/dfinster/branch-wrangler/internal/git"
)

type ActionMsg struct {
	Action string
	Branch string
	Error  error
}

type ConfirmationMsg struct {
	Action      string
	Branch      string
	Description string
	Dangerous   bool
}

func (m Model) handleActionKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if len(m.filteredBranches) == 0 || m.selected >= len(m.filteredBranches) {
		return m, nil
	}

	selectedBranch := m.filteredBranches[m.selected]

	switch msg.String() {
	case "c":
		return m, m.checkoutBranch(selectedBranch.Name)
	case "d":
		if selectedBranch.State == git.StaleLocal {
			return m, m.deleteBranch(selectedBranch.Name, false)
		} else {
			return m, m.createConfirmation("delete", selectedBranch.Name,
				fmt.Sprintf("Branch '%s' is in state '%s'. Are you sure you want to delete it?",
					selectedBranch.Name, selectedBranch.State.DisplayName()), true)
		}
	case "D":
		return m, m.createConfirmation("force-delete", selectedBranch.Name,
			fmt.Sprintf("Force delete branch '%s'? This cannot be undone.", selectedBranch.Name), true)
	case "o":
		if selectedBranch.PRURL != "" {
			return m, m.openPR(selectedBranch.PRURL)
		}
		return m, nil
	case "u":
		return m, m.showUndoView()
	}

	return m, nil
}

func (m Model) checkoutBranch(branchName string) tea.Cmd {
	return func() tea.Msg {
		cmd := exec.Command("git", "checkout", branchName)
		err := cmd.Run()
		return ActionMsg{
			Action: "checkout",
			Branch: branchName,
			Error:  err,
		}
	}
}

func (m Model) deleteBranch(branchName string, force bool) tea.Cmd {
	return func() tea.Msg {
		var cmd *exec.Cmd
		if force {
			cmd = exec.Command("git", "branch", "-D", branchName)
		} else {
			cmd = exec.Command("git", "branch", "-d", branchName)
		}

		err := cmd.Run()
		return ActionMsg{
			Action: "delete",
			Branch: branchName,
			Error:  err,
		}
	}
}

func (m Model) openPR(url string) tea.Cmd {
	return func() tea.Msg {
		var cmd *exec.Cmd

		// Try different browsers/commands based on OS
		browsers := []string{"xdg-open", "open", "start"}

		for _, browser := range browsers {
			cmd = exec.Command(browser, url)
			if cmd.Run() == nil {
				return ActionMsg{
					Action: "open-pr",
					Branch: url,
					Error:  nil,
				}
			}
		}

		return ActionMsg{
			Action: "open-pr",
			Branch: url,
			Error:  fmt.Errorf("failed to open browser"),
		}
	}
}

func (m Model) createConfirmation(action, branch, description string, dangerous bool) tea.Cmd {
	return func() tea.Msg {
		return ConfirmationMsg{
			Action:      action,
			Branch:      branch,
			Description: description,
			Dangerous:   dangerous,
		}
	}
}

func (m Model) showUndoView() tea.Cmd {
	return func() tea.Msg {
		return ActionMsg{
			Action: "undo",
			Branch: "",
			Error:  fmt.Errorf("undo functionality not implemented yet"),
		}
	}
}

func (m Model) handleBulkActions(selectedBranches []git.Branch, action string) tea.Cmd {
	return func() tea.Msg {
		var errors []error

		for _, branch := range selectedBranches {
			switch action {
			case "delete":
				if branch.State == git.StaleLocal {
					cmd := exec.Command("git", "branch", "-d", branch.Name)
					if err := cmd.Run(); err != nil {
						errors = append(errors, fmt.Errorf("failed to delete %s: %w", branch.Name, err))
					}
				}
			case "force-delete":
				cmd := exec.Command("git", "branch", "-D", branch.Name)
				if err := cmd.Run(); err != nil {
					errors = append(errors, fmt.Errorf("failed to force delete %s: %w", branch.Name, err))
				}
			}
		}

		if len(errors) > 0 {
			return ActionMsg{
				Action: "bulk-" + action,
				Branch: fmt.Sprintf("%d branches", len(selectedBranches)),
				Error:  fmt.Errorf("some operations failed: %v", errors),
			}
		}

		return ActionMsg{
			Action: "bulk-" + action,
			Branch: fmt.Sprintf("%d branches", len(selectedBranches)),
			Error:  nil,
		}
	}
}