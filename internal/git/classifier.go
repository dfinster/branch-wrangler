package git

import (
	"context"
	"fmt"

	"github.com/dfinster/branch-wrangler/internal/github"
)

type Classifier struct {
	gitClient    *Client
	githubClient *github.CachedClient
	baseBranches []string
}

func NewClassifier(gitClient *Client, githubClient *github.CachedClient, baseBranches []string) *Classifier {
	return &Classifier{
		gitClient:    gitClient,
		githubClient: githubClient,
		baseBranches: baseBranches,
	}
}

func (c *Classifier) ClassifyBranch(ctx context.Context, branch *Branch) error {
	_, isDetached, err := c.gitClient.GetCurrentBranch()
	if err != nil {
		return err
	}

	if isDetached {
		branch.State = DetachedHead
		return nil
	}
	
	if branch.TrackingRef == "" {
		branch.State = NoUpstream
		return nil
	}

	if branch.CommitCount == 0 {
		branch.State = NoCommits
		return nil
	}

	remoteExists := c.gitClient.RemoteExists(branch.Name)
	if !remoteExists && branch.TrackingRef != "" {
		branch.State = OrphanRemoteDeleted
		return nil
	}

	for _, base := range c.baseBranches {
		if c.gitClient.IsMergedIntoBase(branch.Name, base) {
			branch.State = FullyMergedBase
			return nil
		}
	}

	prs, err := c.githubClient.GetPullRequestsForBranch(ctx, branch.Name)
	if err != nil {
		return c.classifyByGitStatus(branch)
	}

	if len(prs) > 0 {
		return c.classifyByPR(ctx, branch, prs[0])
	}

	return c.classifyByGitStatus(branch)
}

func (c *Classifier) classifyByGitStatus(branch *Branch) error {
	if branch.Ahead == 0 && branch.Behind == 0 {
		branch.State = InSync
		return nil
	}

	if branch.Ahead > 0 && branch.Behind == 0 {
		branch.State = UnpushedAhead
		return nil
	}

	if branch.Ahead == 0 && branch.Behind > 0 {
		branch.State = BehindRemote
		return nil
	}

	if branch.Ahead > 0 && branch.Behind > 0 {
		branch.State = Diverged
		return nil
	}

	branch.State = InSync
	return nil
}

func (c *Classifier) classifyByPR(ctx context.Context, branch *Branch, pr github.PullRequest) error {
	branch.PRNumber = pr.Number
	branch.PRTitle = pr.Title
	branch.PRURL = pr.URL

	if pr.State == "open" {
		if pr.Draft {
			branch.State = DraftPR
		} else {
			branch.State = OpenPR
		}
		return nil
	}

	if pr.State == "closed" {
		if pr.Merged {
			remoteExists, err := c.githubClient.BranchExists(ctx, branch.Name)
			if err != nil {
				return err
			}

			if remoteExists {
				branch.State = MergedRemoteExists
			} else {
				branch.State = StaleLocal
			}
		} else {
			branch.State = ClosedPR
		}
		return nil
	}

	return c.classifyByGitStatus(branch)
}

func (c *Classifier) ClassifyAllBranches(ctx context.Context) ([]Branch, error) {
	branches, err := c.gitClient.ListBranches()
	if err != nil {
		return nil, fmt.Errorf("failed to list branches: %w", err)
	}

	for i := range branches {
		if err := c.ClassifyBranch(ctx, &branches[i]); err != nil {
			return nil, fmt.Errorf("failed to classify branch %s: %w", branches[i].Name, err)
		}
	}

	return branches, nil
}