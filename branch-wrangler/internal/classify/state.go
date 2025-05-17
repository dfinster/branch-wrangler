package classify

// State is one of the branch states
type State string

const (
    LocalOnly State = "LocalOnly"
    InSync    State = "InSync"
    // TODO: add all your states...
)

// Classifier determines the state of a branch
type Classifier struct{}

// New returns a new Classifier
func New() *Classifier { return &Classifier{} }

// Detect inspects branch vs. remote & PR status
func (c *Classifier) Detect(branch string) (State, error) {
    // TODO: use git.Scanner + reconcile.Reconciler
    return LocalOnly, nil
}