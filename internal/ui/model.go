package ui

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dfinster/branch-wrangler/internal/git"
)

type Model struct {
	branches     []git.Branch
	selected     int
	width        int
	height       int
	showHelp     bool
	filter       string
	ctx          context.Context
	classifier   *git.Classifier
	loading      bool
	err          error
}

type LoadBranchesMsg struct {
	branches []git.Branch
	err      error
}

func NewModel(ctx context.Context, classifier *git.Classifier) Model {
	return Model{
		branches:   []git.Branch{},
		selected:   0,
		ctx:        ctx,
		classifier: classifier,
		loading:    true,
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
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.selected > 0 {
				m.selected--
			}
		case "down", "j":
			if m.selected < len(m.branches)-1 {
				m.selected++
			}
		case "?":
			m.showHelp = !m.showHelp
		case "r":
			m.loading = true
			return m, m.loadBranches()
		}
	
	case LoadBranchesMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err
		} else {
			m.branches = msg.branches
			m.err = nil
		}
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
	
	leftPane := m.branchListView()
	rightPane := m.branchDetailsView()
	
	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		leftPane,
		rightPane,
	)
}

func (m Model) branchListView() string {
	var content string
	
	if len(m.branches) == 0 {
		content = "No branches found"
	} else {
		for i, branch := range m.branches {
			cursor := " "
			if i == m.selected {
				cursor = ">"
			}
			
			state := branch.State.DisplayName()
			content += lipgloss.NewStyle().
				Width(m.width/2-4).
				Render(cursor+" "+branch.Name+" ["+state+"]") + "\n"
		}
	}
	
	return lipgloss.NewStyle().
		Width(m.width/2).
		Height(m.height).
		Border(lipgloss.NormalBorder()).
		Padding(1).
		Render(content)
}

func (m Model) branchDetailsView() string {
	var content string
	
	if len(m.branches) == 0 || m.selected >= len(m.branches) {
		content = "No branch selected"
	} else {
		branch := m.branches[m.selected]
		content = "Branch: " + branch.Name + "\n"
		content += "State: " + branch.State.DisplayName() + "\n"
		content += "Last Commit: " + branch.LastCommit.Format("2006-01-02 15:04:05") + "\n"
		content += "Author: " + branch.Author + "\n"
		
		if branch.Ahead > 0 {
			content += "Ahead: " + string(rune(branch.Ahead)) + "\n"
		}
		if branch.Behind > 0 {
			content += "Behind: " + string(rune(branch.Behind)) + "\n"
		}
		if branch.PRNumber > 0 {
			content += "PR: #" + string(rune(branch.PRNumber)) + " - " + branch.PRTitle + "\n"
		}
		if branch.PRURL != "" {
			content += "URL: " + branch.PRURL + "\n"
		}
	}
	
	return lipgloss.NewStyle().
		Width(m.width/2).
		Height(m.height).
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

Actions (coming soon):
  space   Select branch
  d       Delete selected
  c       Checkout branch
  o       Open PR in browser
  f       Filter branches
  /       Search branches

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