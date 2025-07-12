package ui

import (
	"context"
	"fmt"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/dfinster/branch-wrangler/internal/git"
)

type Model struct {
	branches          []git.Branch
	filteredBranches  []git.Branch
	selected          int
	selectedBranches  map[int]bool
	width             int
	height            int
	showHelp          bool
	showFilter        bool
	showConfirmDialog bool
	filter            *Filter
	searchInput       string
	ctx               context.Context
	classifier        *git.Classifier
	loading           bool
	err               error
	lastAction        string
	confirmation      ConfirmationMsg
}

type LoadBranchesMsg struct {
	branches []git.Branch
	err      error
}

func NewModel(ctx context.Context, classifier *git.Classifier) Model {
	return Model{
		branches:         []git.Branch{},
		filteredBranches: []git.Branch{},
		selected:         0,
		selectedBranches: make(map[int]bool),
		ctx:              ctx,
		classifier:       classifier,
		loading:          true,
		filter:           NewFilter(),
	}
}

func (m Model) Init() tea.Cmd {
	return m.loadBranches()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		if m.showFilter {
			return m.handleFilterKeys(msg)
		}

		if m.showConfirmDialog {
			return m.handleConfirmKeys(msg)
		}

		// Handle action keys first
		if newModel, cmd := m.handleActionKeys(msg); cmd != nil {
			return newModel, cmd
		}

		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.selected > 0 {
				m.selected--
			}
		case "down", "j":
			if m.selected < len(m.filteredBranches)-1 {
				m.selected++
			}
		case "?":
			m.showHelp = !m.showHelp
		case "r":
			m.loading = true
			return m, m.loadBranches()
		case "f":
			m.showFilter = !m.showFilter
		case "a":
			m.filter.Clear()
			m.updateFilteredBranches()
		case "/":
			m.filter.SetSearchFilter("")
			m.showFilter = true
		case "1":
			if filter, ok := PredefinedFilters["Stale"]; ok {
				m.filter = &filter
				m.updateFilteredBranches()
			}
		case "2":
			if filter, ok := PredefinedFilters["PR"]; ok {
				m.filter = &filter
				m.updateFilteredBranches()
			}
		case "3":
			if filter, ok := PredefinedFilters["Merged"]; ok {
				m.filter = &filter
				m.updateFilteredBranches()
			}
		case "4":
			if filter, ok := PredefinedFilters["Ahead"]; ok {
				m.filter = &filter
				m.updateFilteredBranches()
			}
		case "space":
			if _, exists := m.selectedBranches[m.selected]; exists {
				delete(m.selectedBranches, m.selected)
			} else {
				m.selectedBranches[m.selected] = true
			}
		}

	case LoadBranchesMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err
		} else {
			m.branches = msg.branches
			m.updateFilteredBranches()
			m.err = nil
		}
		return m, nil

	case ActionMsg:
		m.lastAction = fmt.Sprintf("%s: %s", msg.Action, msg.Branch)
		if msg.Error != nil {
			m.err = msg.Error
		} else {
			m.loading = true
			return m, m.loadBranches()
		}
		return m, nil

	case ConfirmationMsg:
		m.confirmation = msg
		m.showConfirmDialog = true
		return m, nil
	}

	return m, nil
}

func (m Model) View() string {
	if m.loading {
		return "Loading branches..."
	}

	if m.err != nil {
		return "Error: " + m.err.Error()
	}

	if m.showHelp {
		return m.helpView()
	}

	if m.showFilter {
		return m.filterView()
	}

	if m.showConfirmDialog {
		return m.confirmationView()
	}

	header := m.headerView()
	leftPane := m.branchListView()
	rightPane := m.branchDetailsView()

	content := lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftPane,
		rightPane,
	)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		content,
	)
}

func (m Model) branchListView() string {
	var content string

	if len(m.filteredBranches) == 0 {
		content = "No branches match filter"
	} else {
		for i, branch := range m.filteredBranches {
			cursor := " "
			if i == m.selected {
				cursor = ">"
			}

			checkbox := " "
			if _, isSelected := m.selectedBranches[i]; isSelected {
				checkbox = "✓"
			}

			state := branch.State.DisplayName()
			stateColor := m.getStateColor(branch.State)

			line := fmt.Sprintf("%s%s %s", cursor, checkbox, branch.Name)
			if state != "" {
				line += fmt.Sprintf(" [%s]", state)
			}

			content += lipgloss.NewStyle().
				Width(m.width/2-4).
				Foreground(stateColor).
				Render(line) + "\n"
		}
	}

	return lipgloss.NewStyle().
		Width(m.width / 2).
		Height(m.height - 3).
		Border(lipgloss.NormalBorder()).
		Padding(1).
		Render(content)
}

func (m Model) branchDetailsView() string {
	var content string

	if len(m.filteredBranches) == 0 || m.selected >= len(m.filteredBranches) {
		content = "No branch selected"
	} else {
		branch := m.filteredBranches[m.selected]
		content = "Branch: " + branch.Name + "\n"
		content += "State: " + branch.State.DisplayName() + "\n"
		content += "Last Commit: " + branch.LastCommit.Format("2006-01-02 15:04:05") + "\n"
		content += "Author: " + branch.Author + "\n"

		if branch.Ahead > 0 {
			content += "Ahead: " + strconv.Itoa(branch.Ahead) + "\n"
		}
		if branch.Behind > 0 {
			content += "Behind: " + strconv.Itoa(branch.Behind) + "\n"
		}
		if branch.PRNumber > 0 {
			content += "PR: #" + strconv.Itoa(branch.PRNumber) + " - " + branch.PRTitle + "\n"
		}
		if branch.PRURL != "" {
			content += "URL: " + branch.PRURL + "\n"
		}
	}

	return lipgloss.NewStyle().
		Width(m.width / 2).
		Height(m.height - 3).
		Border(lipgloss.NormalBorder()).
		Padding(1).
		Render(content)
}

func (m Model) helpView() string {
	help := `Branch Wrangler Help

Navigation:
  ↑/k     Move up
  ↓/j     Move down
  r       Refresh branches
  ?       Toggle help
  q       Quit

Filtering:
  f       Toggle filter menu
  a       Show all branches
  /       Search branches
  1       Stale branches
  2       PR branches
  3       Merged branches
  4       Ahead branches

Actions:
  space   Select/unselect branch
  c       Checkout branch
  d       Delete branch (safe)
  D       Force delete branch
  o       Open PR in browser
  u       Undo (coming soon)

Press ? to close help`

	return lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Border(lipgloss.NormalBorder()).
		Padding(1).
		Render(help)
}

func (m Model) loadBranches() tea.Cmd {
	return func() tea.Msg {
		branches, err := m.classifier.ClassifyAllBranches(m.ctx)
		return LoadBranchesMsg{branches: branches, err: err}
	}
}
