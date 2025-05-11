## Branch Wrangler — Requirements, Terminology, and Documentation Blueprint

Below is a complete starter‑pack you can drop into a repo (or paste back into ChatGPT) to drive the actual coding effort. It contains:

1. **Canonical terminology** for every branch situation we care about.
2. **Functional & non‑functional requirements** written in concise, testable language.
3. **UI/UX specification** for a full‑screen terminal app.
4. **Additional “don’t‑forget” requirements** that often slip through the cracks.
5. **Ready‑to‑commit files** ‑ a `README.md` outline and a **User Manual** skeleton.

You can copy each section directly or ask ChatGPT to drill deeper into any part.

---

### 1. Terminology Cheat‑Sheet (“Branch States” taxonomy)

| Internal ID                | Display Name (UI)         | Short Definition                                                                                  | Detection Logic (Git + GitHub)                                                  |
|----------------------------|---------------------------|---------------------------------------------------------------------------------------------------|---------------------------------------------------------------------------------|
| **LOCAL\_ONLY**            | *Local‑Only*              | Branch exists only on the developer’s machine.                                                    | `git rev-parse --verify origin/<branch>` fails.                                 |
| **UNPUSHED\_AHEAD**        | *Ahead of Remote*         | Branch tracks `origin/<branch>` but local has commits remote lacks.                               | `git status --branch --porcelain`; count of `ahead` ≠ 0.                        |
| **DIVERGED**               | *Diverged*                | Local & remote both have exclusive commits.                                                       | `ahead` > 0 **and** `behind` > 0.                                               |
| **DRAFT\_PR**              | *Draft PR*                | Linked GitHub PR is in *draft* state.                                                             | GitHub REST v3 `pulls` API → `"draft": true`.                                   |
| **OPEN\_PR**               | *Open PR*                 | Linked PR is open and ready for review.                                                           | GitHub API → `"state": "open"` and `draft` false.                               |
| **MERGED\_REMOTE\_EXISTS** | *Merged (remote kept)*    | PR merged **but** remote branch not yet auto‑deleted.                                             | PR `"merged": true` **and** branch still present on GitHub.                     |
| **STALE\_LOCAL**           | *Merged (remote deleted)* | PR merged, remote branch is gone, local branch redundant → safe to delete.                        | PR `"merged": true` **and** remote branch lookup fails.                         |
| **FULLY\_MERGED\_BASE**    | *Fully Merged Into Base*  | All commits reachable from the branch are already in `main`/`develop` regardless of PR existence. | `git merge-base --is-ancestor <branch> <base>` exit 0.                          |
| **DETACHED\_HEAD**         | —                         | User’s `HEAD` is detached (not a branch).                                                         | `git symbolic-ref -q HEAD` fails. (Shown as a status banner, not a branch row.) |

*Why it matters* — these labels appear in the UI’s sidebar filters and in log / JSON outputs, so keep them stable.

---

### 2. Functional Requirements (FR)

| ID                                              | Requirement                                                                                                                                                                                                                                                                            |
|-------------------------------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| **FR‑1**                                        | *Branch discovery*: Enumerate every local branch, gather metadata (last commit date, author, ahead/behind counts, tracking remote).                                                                                                                                                    |
| **FR‑2**                                        | *GitHub reconciliation*: For branches with a configured upstream on GitHub, query the GitHub API to find an associated PR and its state (draft, open, merged, closed‑unmerged). Handle API‑rate limits with exponential back‑off and local caching (TTL configurable, default 15 min). |
| **FR‑3**                                        | *State classification*: Map each branch into exactly **one** `BranchState` (see §1).                                                                                                                                                                                                   |
| **FR‑4**                                        | *Interactive TUI*: Launches full‑screen; left panel = branch list with color‑coded state badges; right panel = details (commit history preview, PR link, last activity).                                                                                                               |
| **FR‑5**                                        | *Filtering & search*: Toggle by state, text‑search by branch name (fuzzy), sort by last activity date.                                                                                                                                                                                 |
| **FR‑6**                                        | *Bulk actions*: Multi‑select branches; supported actions per state:\\                                                                                                                                                                                                                  |
| - `STALE_LOCAL` → *Delete* (default “d” key).   |                                                                                                                                                                                                                                                                                        |
| - Any → *Checkout* (`Enter`).                   |                                                                                                                                                                                                                                                                                        |
| - Any with remote → *Open PR in browser* (`o`). |                                                                                                                                                                                                                                                                                        |
| - Any ahead/behind → *Rebase onto base* (`r`).  |                                                                                                                                                                                                                                                                                        |
| **FR‑7**                                        | *Safety guard*: Deletion is disabled unless branch is in `STALE_LOCAL` *or* user toggles *force* mode (`F`). Always prompt with a dry‑run preview (`git branch -d` vs `-D`).                                                                                                           |
| **FR‑8**                                        | *Detached HEAD handling*: If `HEAD` detached, show modal explaining risks and quick keys to checkout `main`, `develop`, or stay.                                                                                                                                                       |
| **FR‑9**                                        | *Config file*: `$XDG_CONFIG_HOME/branch‑wrangler/config.yml` (or `%APPDATA%` on Windows) for overrides: default base branches list, color theme, keybindings, GitHub token path.                                                                                                       |
| **FR‑10**                                       | *Headless mode*: `branch‑wrangler --list --json` prints machine‑readable state info for scripting/CI; `--delete-stale --dry-run` CLI flags replicate TUI actions without UI.                                                                                                           |
| **FR‑11**                                       | *Cross‑platform binaries*: Official releases for **macOS**, **Linux**, **Windows** built via `goreleaser`, signed where supported.                                                                                                                                                     |

---

### 3. Non‑Functional Requirements (NFR)

|  Area         | Requirement                                                                                                             |
| ------------- | ----------------------------------------------------------------------------------------------------------------------- |
| Performance   | Scan & classify ≤ 200 branches in < 2 s on a 2020 laptop; GitHub API calls concurrent but ≤ 5 in‑flight to avoid abuse. |
| Reliability   | Unit + integration tests; mocked GitHub API; e2e TUI tests via `expect` harness; 90 % line coverage.                    |
| Portability   | No CGO by default; only depend on stdlib + a single TUI lib (`tview`, `bubbletea`, or similar).                         |
| Security      | Require GitHub token via env `GITHUB_TOKEN` or keychain; never log it; HTTPS only.                                      |
| Accessibility | All commands operable via keyboard; color schemes pass WCAG AA contrast; optional high‑contrast theme.                  |
| Observability | Verbose logs (`--log-level debug`) and structured JSON logs (`--log-format json`).                                      |
| Extensibility | Plugin hooks (Go `plugin` build tag): custom state calculators, extra actions.                                          |
| Documentation | Every exported function has GoDoc; `make docs` runs `godoc -http`; README quick‑start lines < 80 chars.                 |

---

### 4. UI/UX Specification (TUI)

```
┌ Branch Wrangler ───────────────────────────── [q] Quit ─┐
│ State  │ Branch               │ Last Commit  │  Behind  │
│────────┼──────────────────────┼──────────────┼──────────│
│ 🗑 Safe │ feature/login-clean  │ 2025‑05‑09   │ 0        │
│ 💬 PR   │ feat/new‑api‑auth    │ 2025‑05‑10   │ ↑2 ↓0    │
│ ⌛ Draft│ wip/upload‑refactor   │ 2025‑05‑07   │ 0        │
│ 🏠 Local│ tmp/idea‑scratch     │ 2025‑05‑11   │ 0        │
└──────────────────────────────────────────────────────────┘
Details Pane → shows `git log --oneline -5`, PR URL, comment count.
```

**Key Map** (configurable):

| Key            | Action                    |
| -------------- | ------------------------- |
| ↑/↓, PgUp/PgDn | Move cursor               |
| Space          | Toggle select             |
|  f             | Cycle state filter        |
|  /             | Fuzzy search              |
| d              | Delete selected (if safe) |
| F              | Toggle *force delete*     |
| o              | Open PR in \$BROWSER      |
| Enter          | Checkout branch           |
| r              | Rebase onto base          |
| ?              | Help overlay              |

---

### 5. Often‑missed but Valuable Requirements

* **Dry‑run *by default*** (must be opt‑in to “live” destructive mode).
* **Batch undo**: After a delete operation, keep a temp reflog of refs for N minutes and offer “u” to restore.
* **Auto‑update**: `--self‑update` flag (download from GitHub Releases).
* **CI safe‑guard**: Detect if running in CI (env `CI=true`); auto headless mode.
* **Telemetry opt‑in** for anonymous feature usage (helps prioritise).
* **Shell‑completion** generators (`branch‑wrangler completion bash|zsh|fish`).

---

## 6. `README.md` Starter Template

````markdown
# Branch Wrangler 🏇🗑

**Tame your local Git branches**.
Scan, browse, and clean stale branches with a full‑screen terminal UI.

## Features

* Classifies every branch (local‑only, draft PR, merged & stale, …)
* Links directly to GitHub pull requests
* Bulk delete merged branches safely (dry‑run by default)
* Runs everywhere – single static Go binary

## Installation

```bash
brew install branch-wrangler   # macOS (Tap coming soon)
scoop install branch-wrangler  # Windows
curl -sSL https://github.com/…/branch‑wrangler/releases/download/vX.Y.Z/… -o bw
chmod +x bw && sudo mv bw /usr/local/bin
```

## Quick Start

```bash
branch-wrangler            # launch TUI
branch-wrangler --list     # headless JSON summary
branch-wrangler --delete-stale --dry-run  # preview cleanup
```

## Configuration (`~/.config/branch-wrangler/config.yml`)

```yaml
base_branches: [main, develop]
theme: Dracula
github_token_path: ~/.github-token
```

## Keybindings

| Key     | Description            |
| ------- | ---------------------- |
| `↑` `↓` | Move                   |
| `Space` | Select                 |
| `d`     | Delete selected (safe) |
| `F`     | Force delete           |
| `o`     | Open PR                |
| `?`     | Help                   |

*(See full list in the User Manual.)*

## Screenshots

<img src="docs/screenshot-branches.png" width="650">

## Roadmap

- [ ] Undo stack
- [ ] GitLab / Azure DevOps support
- [ ] Plugin API

## Contributing

1. Fork + `git clone`
2. `make setup && make test`
3. Send a PR ✔️

## License

Proprietary

---

## 7. User Manual Skeleton (docs/user‑manual.md)

1. **Introduction**
   *What is Branch Wrangler, typical workflow, supported platforms.*

2. **Installation**
   Homebrew, Scoop, manual, building from source.

3. **Launching the Application**
   *TUI mode vs headless CLI flags.*

4. **Understanding Branch States**
   Table mapping icons → definitions (reuse §1).

5. **Navigating the Interface**
   Full keybinding reference, search, filters.

6. **Performing Actions**
   Deleting branches (safe vs force), opening PRs, rebasing, undo.

7. **Configuration File Reference**
   Every YAML key, defaults, env vars.

8. **Authentication**
   Setting a GitHub token, scopes required (`repo:status read:org`).

9. **Scripting & Automation**
   Examples using `--json` to prune branches in CI.

10. **Troubleshooting & FAQ**
    Rate‑limit errors, detached HEAD, “branch already deleted” warnings.

11. **Contributing & Support**
    Issue templates, code style, contact info.

---

### Next Steps

* **Pick a Go TUI toolkit** (`tview`, `bubbletea`, or `gocui`) and generate stub code.
* **Model** your data: `Branch`, `PRInfo`, `BranchState`.
* Scaffold CLI with `spf13/cobra`.
* Feed any section above back into ChatGPT with “write the code for …”.

Good luck wrangling those branches!
