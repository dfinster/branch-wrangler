package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/dfinster/branch-wrangler/internal/git"
)

func (m Model) headerView() string {
	filterDisplay := m.filter.DisplayName()
	if m.filter.IsActive {
		filterDisplay = "Filter: " + filterDisplay
	}

	count := fmt.Sprintf("(%d/%d branches)", len(m.filteredBranches), len(m.branches))

	left := "Branch Wrangler"
	center := filterDisplay
	right := count

	leftStyle := lipgloss.NewStyle().Bold(true)
	centerStyle := lipgloss.NewStyle().Italic(true)
	rightStyle := lipgloss.NewStyle().Faint(true)

	header := lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftStyle.Render(left),
		strings.Repeat(" ", max(0, m.width-len(left)-len(center)-len(right))),
		centerStyle.Render(center),
		strings.Repeat(" ", max(0, m.width-len(left)-len(center)-len(right))),
		rightStyle.Render(right),
	)

	return lipgloss.NewStyle().
		Width(m.width).
		Height(2).
		Border(lipgloss.NormalBorder(), false, false, true, false).
		Padding(0, 1).
		Render(header)
}

func (m Model) filterView() string {
	content := "Filter Options:\n\n"
	content += "a - All branches\n"
	content += "1 - Stale branches\n"
	content += "2 - PR branches\n"
	content += "3 - Merged branches\n"
	content += "4 - Ahead branches\n"
	content += "/ - Search by name\n\n"

	if m.filter.Mode == FilterBySearch {
		content += "Search: " + m.searchInput + "\n"
		content += "Type to search, Enter to apply, Esc to cancel\n"
	}

	content += "\nPress f to close filter menu"
	
	return lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Border(lipgloss.NormalBorder()).
		Padding(1).
		Render(content)
}

func (m Model) getStateColor(state git.BranchState) lipgloss.Color {
	switch state {
	case git.StaleLocal:
		return lipgloss.Color("9")  // Red
	case git.OpenPR:
		return lipgloss.Color("10") // Green
	case git.DraftPR:
		return lipgloss.Color("11") // Yellow
	case git.MergedRemoteExists:
		return lipgloss.Color("12") // Blue
	case git.UnpushedAhead:
		return lipgloss.Color("13") // Magenta
	case git.BehindRemote:
		return lipgloss.Color("14") // Cyan
	case git.Diverged:
		return lipgloss.Color("9")  // Red
	case git.InSync:
		return lipgloss.Color("10") // Green
	default:
		return lipgloss.Color("7")  // White
	}
}

func (m Model) handleFilterKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "escape", "f":
		m.showFilter = false
		return m, nil
	case "a":
		m.filter.Clear()
		m.updateFilteredBranches()
		m.showFilter = false
		return m, nil
	case "1":
		if filter, ok := PredefinedFilters["Stale"]; ok {
			m.filter = &filter
			m.updateFilteredBranches()
			m.showFilter = false
		}
		return m, nil
	case "2":
		if filter, ok := PredefinedFilters["PR"]; ok {
			m.filter = &filter
			m.updateFilteredBranches()
			m.showFilter = false
		}
		return m, nil
	case "3":
		if filter, ok := PredefinedFilters["Merged"]; ok {
			m.filter = &filter
			m.updateFilteredBranches()
			m.showFilter = false
		}
		return m, nil
	case "4":
		if filter, ok := PredefinedFilters["Ahead"]; ok {
			m.filter = &filter
			m.updateFilteredBranches()
			m.showFilter = false
		}
		return m, nil
	case "/":
		m.filter.SetSearchFilter("")
		m.searchInput = ""
		return m, nil
	case "enter":
		if m.filter.Mode == FilterBySearch {
			m.filter.SetSearchFilter(m.searchInput)
			m.updateFilteredBranches()
			m.showFilter = false
		}
		return m, nil
	case "backspace":
		if m.filter.Mode == FilterBySearch && len(m.searchInput) > 0 {
			m.searchInput = m.searchInput[:len(m.searchInput)-1]
			m.filter.SetSearchFilter(m.searchInput)
			m.updateFilteredBranches()
		}
		return m, nil
	default:
		if m.filter.Mode == FilterBySearch && len(msg.String()) == 1 {
			m.searchInput += msg.String()
			m.filter.SetSearchFilter(m.searchInput)
			m.updateFilteredBranches()
		}
		return m, nil
	}
}

func (m *Model) updateFilteredBranches() {
	m.filteredBranches = m.filter.Apply(m.branches)
	if m.selected >= len(m.filteredBranches) {
		m.selected = max(0, len(m.filteredBranches)-1)
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}