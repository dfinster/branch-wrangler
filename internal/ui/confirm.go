package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m Model) confirmationView() string {
	title := "Confirmation Required"
	
	var actionColor lipgloss.Color
	if m.confirmation.Dangerous {
		actionColor = lipgloss.Color("9") // Red
	} else {
		actionColor = lipgloss.Color("10") // Green
	}
	
	content := lipgloss.NewStyle().
		Bold(true).
		Foreground(actionColor).
		Render(title) + "\n\n"
	
	content += m.confirmation.Description + "\n\n"
	content += "Press 'y' to confirm, 'n' to cancel"
	
	return lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(actionColor).
		Padding(2).
		Margin(2).
		Render(content)
}

func (m Model) handleConfirmKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y":
		m.showConfirmDialog = false
		switch m.confirmation.Action {
		case "delete":
			return m, m.deleteBranch(m.confirmation.Branch, false)
		case "force-delete":
			return m, m.deleteBranch(m.confirmation.Branch, true)
		}
		return m, nil
	case "n", "escape":
		m.showConfirmDialog = false
		return m, nil
	}
	return m, nil
}