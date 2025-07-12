package git

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	workingDir string
}

func NewClient(workingDir string) *Client {
	return &Client{workingDir: workingDir}
}

func (c *Client) IsGitRepo() bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	cmd.Dir = c.workingDir
	return cmd.Run() == nil
}

func (c *Client) GetCurrentBranch() (string, bool, error) {
	cmd := exec.Command("git", "symbolic-ref", "-q", "HEAD")
	cmd.Dir = c.workingDir
	output, err := cmd.Output()
	if err != nil {
		return "", true, nil
	}

	branchRef := strings.TrimSpace(string(output))
	if strings.HasPrefix(branchRef, "refs/heads/") {
		return strings.TrimPrefix(branchRef, "refs/heads/"), false, nil
	}

	return "", true, nil
}

func (c *Client) ListBranches() ([]Branch, error) {
	cmd := exec.Command("git", "for-each-ref", "--format=%(refname:short)|%(committerdate:iso)|%(authorname)|%(upstream:short)|%(HEAD)", "refs/heads/")
	cmd.Dir = c.workingDir
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list branches: %w", err)
	}

	var branches []Branch
	scanner := bufio.NewScanner(bytes.NewReader(output))

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "|")
		if len(parts) < 5 {
			continue
		}

		name := parts[0]
		commitDate, _ := time.Parse("2006-01-02 15:04:05 -0700", parts[1])
		author := parts[2]
		upstream := parts[3]
		isCurrent := parts[4] == "*"

		branch := Branch{
			Name:        name,
			LastCommit:  commitDate,
			Author:      author,
			TrackingRef: upstream,
			IsCurrent:   isCurrent,
		}

		if upstream != "" {
			ahead, behind, err := c.getAheadBehind(name, upstream)
			if err == nil {
				branch.Ahead = ahead
				branch.Behind = behind
			}
		}

		commitCount, err := c.getCommitCount(name)
		if err == nil {
			branch.CommitCount = commitCount
		}

		branches = append(branches, branch)
	}

	return branches, scanner.Err()
}

func (c *Client) getAheadBehind(local, remote string) (int, int, error) {
	cmd := exec.Command("git", "rev-list", "--left-right", "--count", fmt.Sprintf("%s...%s", local, remote))
	cmd.Dir = c.workingDir
	output, err := cmd.Output()
	if err != nil {
		return 0, 0, err
	}

	parts := strings.Fields(strings.TrimSpace(string(output)))
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("unexpected output format")
	}

	ahead, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, err
	}

	behind, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, err
	}

	return ahead, behind, nil
}

func (c *Client) getCommitCount(branch string) (int, error) {
	cmd := exec.Command("git", "rev-list", "--count", branch)
	cmd.Dir = c.workingDir
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	count, err := strconv.Atoi(strings.TrimSpace(string(output)))
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (c *Client) RemoteExists(branch string) bool {
	cmd := exec.Command("git", "rev-parse", "--verify", fmt.Sprintf("origin/%s", branch))
	cmd.Dir = c.workingDir
	return cmd.Run() == nil
}

func (c *Client) IsMergedIntoBase(branch, base string) bool {
	cmd := exec.Command("git", "merge-base", "--is-ancestor", branch, base)
	cmd.Dir = c.workingDir
	return cmd.Run() == nil
}

func (c *Client) GetRemoteURL() (string, error) {
	cmd := exec.Command("git", "config", "--get", "remote.origin.url")
	cmd.Dir = c.workingDir
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}

func (c *Client) ParseGitHubRepo(remoteURL string) (owner, repo string, err error) {
	url := strings.TrimSpace(remoteURL)

	if strings.HasPrefix(url, "git@github.com:") {
		url = strings.TrimPrefix(url, "git@github.com:")
	} else if strings.HasPrefix(url, "https://github.com/") {
		url = strings.TrimPrefix(url, "https://github.com/")
	} else {
		return "", "", fmt.Errorf("not a GitHub repository URL")
	}

	url = strings.TrimSuffix(url, ".git")

	parts := strings.Split(url, "/")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid GitHub repository format")
	}

	return parts[0], parts[1], nil
}