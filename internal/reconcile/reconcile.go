package reconcile

// Reconciler looks up GitHub PR status for branches
type Reconciler struct {
    // TODO: GitHub client, cache, etc.
}

// New returns a new Reconciler
func New() *Reconciler {
    return &Reconciler{}
}

// LookupPRStatus returns PR info for a branch
func (r *Reconciler) LookupPRStatus(branch string) (string, error) {
    // TODO: GitHub API call
    return "", nil
}