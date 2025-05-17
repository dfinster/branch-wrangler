### Executive Summary

*Branch Wrangler – Local Git‑Branch Management Tool*

---

#### Business Context & Problem

Development teams that work with **short‑lived feature branches** routinely accumulate dozens—or hundreds—of stale local branches. While GitHub can automatically delete a branch once its pull‑request (PR) is merged, a common pattern is for **each developer’s workstation to retain the branch until manually deleted**. This leads to:

* Cluttered `git branch` outputs that slow day‑to‑day workflows.
* Accidental check‑outs of obsolete code, causing wasted debugging time.
* Higher cognitive load when evaluating which branches are still relevant.

---

#### Solution Overview

**Branch Wrangler** is a cross‑platform, full‑screen **terminal UI (TUI)** application written in **Go** that scans every local branch, reconciles it with GitHub, classifies its status, and enables one‑click (or one‑key) cleanup. It targets developers who work over SSH, inside lightweight containers, or on desktop terminals—providing the familiar, keyboard‑centric experience of tools such as *Midnight Commander*, *htop*, or *Atuin*.

---

#### Core Objectives

| Objective                     | Outcome                                                                                                                                |
| ----------------------------- | -------------------------------------------------------------------------------------------------------------------------------------- |
| **Accurate State Detection**  | Automatically identify more than a dozen distinct branch states (e.g., Draft PR, Ahead of Remote, Fully Merged, Stale Local, etc.).                       |
| **Actionable UI**             | Let users filter, search, and bulk‑delete branches that are provably safe to remove, while offering jump‑to‑PR and checkout shortcuts. |
| **Safety & Auditability**     | Run in *dry‑run* mode by default, provide undo for recent deletions, and expose a headless interface for CI or scripting.         |
| **Portability & Performance** | Ship as a single static Go binary with <2 s scan time for 200 branches on commodity laptops, no external services required beyond GitHub’s API. |

---

#### Functional Requirements (Highlights)

1. **Branch Discovery** – Enumerate all local branches; capture ahead/behind counts, last commit metadata, and upstream tracking info.
2. **GitHub Reconciliation** – For branches with an upstream, query GitHub’s REST API to attach PR numbers, URLs, draft/open/merged states, and merge dates—cached locally for 15 minutes with concurrency limits.
3. **State Classification** – Map each branch to exactly one of the canonical *BranchStates* (e.g., `STALE_LOCAL`, `OPEN_PR`, `LOCAL_ONLY`).
4. **Interactive TUI** – Split‑pane interface: list‑view on the left, details pane on the right; keyboard‑only navigation; configurable color themes.
5. **Bulk Actions & Shortcuts** – Delete safe branches (`d`), force‑delete selected branches (`F`), open PR in web browser (`o`), checkout branch (`c`), undo (`u`).
6. **Detached‑HEAD Guard** – Detect detached `HEAD`; present a modal prompting checkout of the default branch at GitHub (usually `main` or `develop`), or continue at user’s own risk.
7. **Config & Extensibility** – YAML config under XDG (`~/.config/branch‑wrangler`); CLI flags for headless operation.

---

#### Non‑Functional Requirements

* **Performance:**  <2 s classification for 200 branches; ≤5 concurrent API calls.
* **Reliability:** 90 % test coverage, mocked GitHub calls, end‑to‑end TUI tests.
* **Security:** OAuth token via env or keychain; never logged; HTTPS enforced.
* **Accessibility:** Full keyboard control, WCAG‑AA contrast defaults, high‑contrast theme.
* **Portability:** No CGO; official binaries for macOS, Linux, Windows.
* **Observability:** Structured JSON logging and adjustable log levels.

---

By delivering an intuitive, safety‑focused tool that integrates seamlessly with existing GitHub workflows and terminal habits, **Branch Wrangler** will reclaim developer focus, reduce mistakes, and keep local repositories lean—ultimately streamlining day‑to‑day Git operations across any team that embraces rapid, PR‑driven development.
