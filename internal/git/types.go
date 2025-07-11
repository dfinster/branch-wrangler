package git

import "time"

type BranchState string

const (
	DetachedHead        BranchState = "DETACHED_HEAD"
	NoUpstream          BranchState = "NO_UPSTREAM"
	OrphanRemoteDeleted BranchState = "ORPHAN_REMOTE_DELETED"
	InSync              BranchState = "IN_SYNC"
	UnpushedAhead       BranchState = "UNPUSHED_AHEAD"
	BehindRemote        BranchState = "BEHIND_REMOTE"
	Diverged            BranchState = "DIVERGED"
	DraftPR             BranchState = "DRAFT_PR"
	OpenPR              BranchState = "OPEN_PR"
	ClosedPR            BranchState = "CLOSED_PR"
	MergedRemoteExists  BranchState = "MERGED_REMOTE_EXISTS"
	StaleLocal          BranchState = "STALE_LOCAL"
	FullyMergedBase     BranchState = "FULLY_MERGED_BASE"
	NoCommits           BranchState = "NO_COMMITS"
	UpstreamChanged     BranchState = "UPSTREAM_CHANGED"
	RemoteRenamed       BranchState = "REMOTE_RENAMED"
	UpstreamGone        BranchState = "UPSTREAM_GONE"
)

func (s BranchState) DisplayName() string {
	switch s {
	case DetachedHead:
		return "—"
	case NoUpstream:
		return "No Upstream"
	case OrphanRemoteDeleted:
		return "Orphan (remote deleted)"
	case InSync:
		return "In Sync with Remote"
	case UnpushedAhead:
		return "Ahead of Remote"
	case BehindRemote:
		return "Behind Remote"
	case Diverged:
		return "Diverged from Remote"
	case DraftPR:
		return "Draft PR"
	case OpenPR:
		return "Open PR"
	case ClosedPR:
		return "Closed PR"
	case MergedRemoteExists:
		return "Merged (remote kept)"
	case StaleLocal:
		return "Merged (remote deleted)"
	case FullyMergedBase:
		return "Fully Merged Into Base"
	case NoCommits:
		return "Empty Branch"
	case UpstreamChanged:
		return "Upstream Moved"
	case RemoteRenamed:
		return "Remote Renamed"
	case UpstreamGone:
		return "Upstream Gone"
	default:
		return string(s)
	}
}

type Branch struct {
	Name         string
	State        BranchState
	LastCommit   time.Time
	Author       string
	Ahead        int
	Behind       int
	TrackingRef  string
	PRNumber     int
	PRTitle      string
	PRURL        string
	IsCurrent    bool
	CommitCount  int
	LastCommitSHA string
}

type GitStatus struct {
	CurrentBranch string
	IsDetached    bool
	Ahead         int
	Behind        int
	HasUpstream   bool
	UpstreamRef   string
}