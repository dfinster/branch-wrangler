package git

// Scanner examines the local Git repo for branches
type Scanner struct {
    // TODO: fields for exec.Command, repo path, etc.
}

// NewScanner returns a new Scanner
func NewScanner() *Scanner {
    return &Scanner{}
}

// ListBranches returns all local branches
func (s *Scanner) ListBranches() ([]string, error) {
    // TODO: run `git branch --format="%(refname:short)"`
    return nil, nil
}