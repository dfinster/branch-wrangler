package ui

import (
	"strings"

	"github.com/dfinster/branch-wrangler/internal/git"
)

type FilterMode int

const (
	FilterAll FilterMode = iota
	FilterByState
	FilterBySearch
	FilterByCustom
)

type Filter struct {
	Mode       FilterMode
	States     []git.BranchState
	SearchTerm string
	CustomName string
	IsActive   bool
}

func NewFilter() *Filter {
	return &Filter{
		Mode:     FilterAll,
		States:   []git.BranchState{},
		IsActive: false,
	}
}

func (f *Filter) Apply(branches []git.Branch) []git.Branch {
	if !f.IsActive {
		return branches
	}

	var filtered []git.Branch

	for _, branch := range branches {
		if f.matches(branch) {
			filtered = append(filtered, branch)
		}
	}

	return filtered
}

func (f *Filter) matches(branch git.Branch) bool {
	switch f.Mode {
	case FilterAll:
		return true
	case FilterByState:
		return f.matchesState(branch.State)
	case FilterBySearch:
		return f.matchesSearch(branch.Name)
	case FilterByCustom:
		return f.matchesState(branch.State) || f.matchesSearch(branch.Name)
	}
	return true
}

func (f *Filter) matchesState(state git.BranchState) bool {
	if len(f.States) == 0 {
		return true
	}

	for _, filterState := range f.States {
		if state == filterState {
			return true
		}
	}
	return false
}

func (f *Filter) matchesSearch(branchName string) bool {
	if f.SearchTerm == "" {
		return true
	}

	return strings.Contains(strings.ToLower(branchName), strings.ToLower(f.SearchTerm))
}

func (f *Filter) SetStateFilter(states []git.BranchState) {
	f.Mode = FilterByState
	f.States = states
	f.IsActive = true
	f.SearchTerm = ""
}

func (f *Filter) SetSearchFilter(term string) {
	f.Mode = FilterBySearch
	f.SearchTerm = term
	f.IsActive = term != ""
	f.States = []git.BranchState{}
}

func (f *Filter) SetCustomFilter(name string, states []git.BranchState, searchTerm string) {
	f.Mode = FilterByCustom
	f.CustomName = name
	f.States = states
	f.SearchTerm = searchTerm
	f.IsActive = true
}

func (f *Filter) Clear() {
	f.Mode = FilterAll
	f.States = []git.BranchState{}
	f.SearchTerm = ""
	f.CustomName = ""
	f.IsActive = false
}

func (f *Filter) DisplayName() string {
	if !f.IsActive {
		return "All Branches"
	}

	switch f.Mode {
	case FilterByState:
		if len(f.States) == 1 {
			return f.States[0].DisplayName()
		}
		return "Multiple States"
	case FilterBySearch:
		return "Search: " + f.SearchTerm
	case FilterByCustom:
		return f.CustomName
	}

	return "All Branches"
}

var PredefinedFilters = map[string]Filter{
	"Stale": {
		Mode:       FilterByState,
		States:     []git.BranchState{git.StaleLocal},
		IsActive:   true,
		CustomName: "Stale Branches",
	},
	"PR": {
		Mode:       FilterByState,
		States:     []git.BranchState{git.OpenPR, git.DraftPR, git.ClosedPR},
		IsActive:   true,
		CustomName: "Has PR",
	},
	"Merged": {
		Mode:       FilterByState,
		States:     []git.BranchState{git.MergedRemoteExists, git.StaleLocal, git.FullyMergedBase},
		IsActive:   true,
		CustomName: "Merged Branches",
	},
	"Ahead": {
		Mode:       FilterByState,
		States:     []git.BranchState{git.UnpushedAhead, git.Diverged},
		IsActive:   true,
		CustomName: "Ahead of Remote",
	},
}
