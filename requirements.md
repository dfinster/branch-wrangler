## Branch Wrangler — Requirements, Terminology, and Documentation Blueprint

Below is a complete description of a new software project. It contains:

1. **Branch States Taxonomy** for every branch situation we care about.
2. **Functional & non‑functional requirements** written in concise, testable language.
3. **UI/UX specification** for a full‑screen terminal app.
4. **Additional “don’t‑forget” requirements** that often slip through the cracks.

### Branch States Taxonomy

| Internal ID                      | Display Name (UI)         | Short Definition                                                                                          | Detection Logic                                                                       |
| -------------------------------- | ------------------------- | --------------------------------------------------------------------------------------------------------- | ------------------------------------------------------------------------------------- |
| **DETACHED\_HEAD**               | —                         | `HEAD` is not on any branch.                                                                              | `git symbolic-ref -q HEAD` fails.                                                     |
| **NO\_UPSTREAM**                 | *No Upstream*             | Branch exists only locally and has never been pushed (no tracking set).                                   | `git rev-parse --abbrev-ref --symbolic-full-name @{u}` fails.                         |
| **ORPHAN\_REMOTE\_DELETED**      | *Orphan (remote deleted)* | Branch was pushed at one point but its remote was removed (no PR or PR unmerged).                         | `git rev-parse --verify origin/<branch>` fails **and** upstream is configured.        |
| **IN\_SYNC**                     | *In Sync with Remote*     | Local and `origin/<branch>` point to the same commit.                                                     | `git status --branch --porcelain` shows `ahead=0` **and** `behind=0`.                 |
| **UNPUSHED\_AHEAD**              | *Ahead of Remote*         | Local has commits that aren’t on `origin/<branch>`.                                                       | `git status --branch --porcelain`; `ahead` ≠ 0 **and** `behind` = 0.                  |
| **BEHIND\_REMOTE**               | *Behind Remote*           | Remote has commits that local branch doesn’t.                                                             | `git status --branch --porcelain`; `behind` ≠ 0 **and** `ahead` = 0.                  |
| **DIVERGED**                     | *Diverged from Remote*    | Local and remote each have unique commits.                                                                | `git status --branch --porcelain`; `ahead` > 0 **and** `behind` > 0.                  |
| **DRAFT\_PR**                    | *Draft PR*                | A GitHub PR exists for the branch, but it’s still in draft.                                               | GitHub REST API `GET /repos/:owner/:repo/pulls` → `"draft": true`.                    |
| **OPEN\_PR**                     | *Open PR*                 | A GitHub PR exists and is open (not a draft), awaiting review.                                            | GitHub API → `"state": "open"` **and** `"draft": false`.                              |
| **CLOSED\_PR**                   | *Closed PR*               | A PR was created but closed without being merged.                                                         | GitHub API → `"state": "closed"` **and** `"merged": false`.                           |
| **MERGED\_REMOTE\_EXISTS**       | *Merged (remote kept)*    | PR was merged on GitHub, but the remote branch still exists (auto-delete off).                            | GitHub API → `"merged": true` **and** remote branch lookup succeeds.                  |
| **STALE\_LOCAL**                 | *Merged (remote deleted)* | PR merged and GitHub auto-deleted the branch—safe to delete locally.                                      | GitHub API → `"merged": true` **and** `git rev-parse --verify origin/<branch>` fails. |
| **FULLY\_MERGED\_BASE**          | *Fully Merged Into Base*  | All commits from this branch are already in the base branch (`main`/`develop`), regardless of PR history. | `git merge-base --is-ancestor <branch> <base>` exits 0.                               |
| **NO\_COMMITS** *(New)*          | *Empty Branch*            | Local branch exists but has no commits (newly created, empty state).                                      | `git rev-list --count <branch>` returns 0.                                            |
| **UPSTREAM\_CHANGED** *(New)*    | *Upstream Moved*          | Remote tracking branch was force-pushed or rebased, history diverged significantly.                       | `git rev-list --left-right <branch>...origin/<branch>` shows unrelated history.       |
| **REMOTE\_RENAMED** *(New)*      | *Remote Renamed*          | Remote branch was renamed; local tracking reference outdated.                                             | GitHub API → original remote missing; PR reference shows new name exists.             |
| **UPSTREAM\_GONE** *(Clarified)* | *Upstream Gone*           | Upstream branch explicitly deleted (distinct from orphaned).                                              | Upstream configured, but GitHub explicitly reports 404 on remote branch reference.    |

*Why it matters* — these labels appear in the UI’s sidebar filters and in log / JSON outputs, so keep them stable.

---

### 2. Functional Requirements (FR)


| ID     | Requirement                                                                                                                                                                                                                                                                                                                                                      |
|--------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| FR-1   | *Branch discovery*: Enumerate every local branch, gather metadata (last commit date, author, ahead/behind counts, tracking remote).                                                                                                                                                                                                                              |
| FR-2   | *GitHub reconciliation*: For branches with a configured upstream on GitHub, query the GitHub API to find an associated PR and its state (draft, open, merged, closed-unmerged). Handle API-rate limits with exponential back-off and local caching (TTL configurable, default 15 min).                                                                           |
| FR-3   | *State classification*: Map each branch into exactly **one** `BranchState` (see Branch States Taxonomy).                                                                                                                                                                                                                                                         |
| FR-4   | *Interactive TUI*: Launches full-screen; left panel = branch list with color-coded state badges; right panel = details (commit history preview, PR link, last activity).                                                                                                                                                                                         |
| FR-5   | *Filtering & search*: Toggle by state, text-search by branch name (fuzzy), sort by last activity date.                                                                                                                                                                                                                                                           |
| FR-6   | *Bulk actions*: Multi-select branches by state or individually; If all selected branches are in the same state, enable *Delete* (“d” key) for all selected if `STALE_LOCAL` state, or allow *Force delete* (default "F" key, with confirmation explaining the current branch state) in other states                                                              |
| FR-6.1 | *Individual actions*: When a single branch only is selected, allow these actions: <br/ > - For any branch, even if a local does not exist: *Checkout* (`c`);<br /> - For any branch with a remote: *Open PR in browser* (`o`).                                                                                                                                   |
| FR-7   | *Safety guard*: Deletion is disabled unless branch is in `STALE_LOCAL` *or* user toggles *force* mode (`F`). Always prompt with a dry-run preview (`git branch -d` vs `-D`).                                                                                                                                                                                     |
| FR-8   | *Detached HEAD handling*: If `HEAD` detached, show modal explaining risks and quick keys to checkout `main`, `develop`, or stay.                                                                                                                                                                                                                                 |
| FR-9   | *Config file*: `$XDG_CONFIG_HOME/branch-wrangler/config.yml` (or `%APPDATA%` on Windows) for overrides: default base branches list, color theme, keybindings, GitHub token path.                                                                                                                                                                                 |
| FR-10  | *Headless mode*: See full list at (Headless mode (CLI) options)                                                                                                                                                                                                                                                                                                    |
| FR-11  | *Cross-platform binaries*: Official releases for **macOS**, **Linux**, **Windows** built via `goreleaser`, signed where supported.                                                                                                                                                                                                                               |
| FR-12  | *Token precedence*: On startup, the app shall first check the `GITHUB_TOKEN` environment variable.                                                                                                                                                                                                                                                               |
| FR-13  | *Token precedence*: If unset, it shall check the config file at `$XDG_CONFIG_HOME/branch-wrangler/config.yml` (fallback to `~/.config/branch-wrangler/config.yml`), reading a `token:` field.                                                                                                                                                                    |
| FR-14  | *Token precedence*: If still unauthenticated, it shall enter _OAuth Device Flow_.                                                                                                                                                                                                                                                                                  |
| FR-15  | *Config file format*: The config file shall be valid YAML and support a root-level `token` key whose value is the GitHub access token.                                                                                                                                                                                                                           |
| FR-16  | *OAuth Device Flow*: The app shall implement GitHub’s OAuth Device Flow:<br>1. POST to `/login/device/code` with the client ID;<br>2. Display the returned `user_code` and verification URI to the user;<br>3. Poll `/login/oauth/access_token` until authorization is granted or expired;<br>4. Persist the received token into the config file under `token:`. |
| FR-17  | *OAuth Device Flow*: If the device flow fails (e.g. timeout, network error), the app shall display a clear error and exit non-zero.                                                                                                                                                                                                                              |
| FR-18  | *Token storage & security*: Tokens written to the config file shall overwrite any previous `token:` entry. It shall preserve the previous token value as a comment. If a previous comment exists, it shall make a versioned comment to preserve all previous token history.                                                                                                                                                                                                                                                       |
| FR-19  | *Token storage & security*: The app shall set file permissions on the config directory/files to user-only read/write where the OS supports it (e.g. `chmod 600` on Unix).                                                                                                                                                                                        |
| FR-20  | *Token validation*: After loading a token (from env or config), the app shall make a lightweight “ping” (e.g. `GET /rate_limit`) to verify it’s valid.                                                                                                                                                                                                           |
| FR-21  | *Token validation*: If validation fails (expired/revoked), the app shall issue an error message indicating which type of token (env or config) failed and offer instructions to correct the error, including an option to re-login with OAuth.                                                                                                                                                                                                   |
| FR-22  | *CLI & UX*: Provide a `branch-wrangler --login` command to force an interactive authentication.                                                                                                                                                                                                                                                                    |
| FR-23  | *CLI & UX*: Provide a `branch-wrangler --logout` command that clears any stored token from the config file.                                                                                                                                                                                                                                                        |
| FR-24  | *CLI & UX*: In verbose or help output, document the three authentication methods and advise users on PAT scopes (e.g. `repo`, `workflow`).                                                                                                                                                                                                                       |


### 3. Non‑Functional Requirements (NFR)

| Area          | Requirement                                                                                                             |
|---------------|-------------------------------------------------------------------------------------------------------------------------|
| Performance   | Scan & classify ≤ 200 branches in < 2 s on a 2020 laptop; GitHub API calls concurrent but ≤ 5 in‑flight to avoid abuse. |
| Reliability   | Unit + integration tests; mocked GitHub API; e2e TUI tests via `expect` harness; 90 % line coverage.                    |
| Portability   | No CGO by default; only depend on stdlib + a `tview` TUI lib                                                            |
| Security      | Use the GitHub token specified in `~/.gitconfig` or `~/.github-token`; never log the token; HTTPS enforced.             |
| Accessibility | All commands operable via keyboard; color schemes pass WCAG AA contrast; optional high‑contrast theme.                  |
| Observability | Verbose logs (`--log-level debug`) and structured JSON logs (`--log-format json`).                                      |
| Extensibility | Plugin hooks (Go `plugin` build tag): custom state calculators, extra actions.                                          |
| Documentation | Every exported function has GoDoc; `make docs` runs `godoc -http`; README quick‑start lines < 80 chars.                 |


### 4. UI/UX Specification (TUI)

The user interface is created using the `tview` TUI library and consists of two main panes:

[Please describe the TUI details here]

#### 4.1 Key Map:

| Key            | Action                                                                         |
|----------------|--------------------------------------------------------------------------------|
| ↑/↓, PgUp/PgDn | Move cursor                                                                    |
| Space          | Toggle select                                                                  |
| a              | Show all branches in all states                                                |
| f              | Filter branches to show only one selected state                                |
| /              | Fuzzy search for branch name                                                   |
| d              | Delete selected (if safe)                                                      |
| F              | Toggle *force delete* (with confirmation explaining the current branch state)  |
| u              | Switch to Undo mode and show cached reflog. Allow restore of deleted branches. |
| o              | Open PR in \$BROWSER                                                           |
| c              | Checkout branch                                                                |
| ?              | Help overlay                                                                   |


### 5. Often‑missed but Valuable Requirements

- **Dry‑run _by default_ in TUI mode**: must manually activate destructive mode.
- **Undo**: After a delete operation, keep a cached reflog of refs from last session and allow restore.

### Headless mode (CLI) options

- `branch-wrangler --help` returns help text.
- `branch-wrangler --version` returns version info.
- `branch-wrangler --undo` Read the cached reflog and restore the deleted branches.
- `branch-wrangler --list` returns headless human-readable output.
- `branch-wrangler --log` returns headless verbose debug output.
- `branch-wrangler --json` returns headless verbose debug output in json format.
- `branch-wrangler --github-token-path` to override the default token location.
- `branch-wrangler --base-branches` to override the default base branches list.
- `branch-wrangler --config [path to file]` Override the default configuration file location (`$XDG_CONFIG_HOME/branch‑wrangler/config.yml` or `%APPDATA%` on Windows)
- `branch-wrangler --delete-stale` headless cleanup, deletes branches that are safe to delete.
- `branch-wrangler --delete-stale --dry-run` headless cleanup preview, prints branches that would be deleted.
- `branch‑wrangler --completion bash|zsh|fish` Shell‑completion generators

## Configuration file (`~/.config/branch-wrangler/config.yml`)

```yaml
base_branches: [main, develop]
github_token_path: ~/.github-token
```

### Next Steps

* **Pick a Go TUI toolkit** (`tview`, `bubbletea`, or `gocui`) and generate stub code.
* **Model** your data: `Branch`, `PRInfo`, `BranchState`.
* Scaffold CLI with `spf13/cobra`.
* Feed any section above back into ChatGPT with “write the code for …”.

Good luck wrangling those branches!
