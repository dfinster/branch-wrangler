Here’s a recommended “starter” layout for your Branch Wrangler repository. It follows Go conventions, keeps your TUI code cleanly separated from your core logic, and gives you room for docs, tests, and CI:

branch-wrangler/
├── .github/
│   └── workflows/
│       └── ci.yml              # GitHub Actions for linting, tests, builds
├── cmd/
│   └── branch-wrangler/
│       └── main.go             # the entrypoint: parses flags, kickstarts the TUI
├── internal/
│   ├── git/
│   │   └── scanner.go          # scanning local branches, querying GitHub
│   ├── reconcile/
│   │   └── reconcile.go        # GitHub PR lookup, caching, rate-limit handling
│   ├── classify/
│   │   └── state.go            # BranchState taxonomy and detection logic
│   └── ui/
│       └── tui.go              # tview (or other) UI setup and screen management
├── pkg/                        # any reusable packages you might expose to others
│   └── config/
│       └── config.go           # global config parsing (flags, env, config file)
├── docs/
│   ├── architecture.md         # high-level design and component diagrams
│   ├── requirements.md         # your FRs, SRs, and taxonomies
│   └── user-guide.md           # how to install & use the tool
├── scripts/
│   ├── build.sh                # convenience scripts (e.g. cross-compile)
│   ├── release.sh              # tagging / GitHub release helper
│   └── lint.sh                 # run golangci-lint, go fmt, etc.
├── examples/                   # small example repos or fixtures for testing
│   └── demo/
│       └── README.md
├── testdata/                   # fixtures & sample git repos for unit tests
│   └── repo-with-merged-branches/
├── Dockerfile                  # for containerized builds / demos
├── go.mod
├── go.sum
├── .gitignore
├── README.md
└── LICENSE

What goes where
	•	cmd/branch-wrangler/main.go
Your main package. Wire up flags (e.g. --dry-run, --cache-ttl), initialize your GitHub client, then hand off to internal/ui to start the TUI.
	•	internal/…
Non-exported application code.
	•	git/: interacting with the local Git CLI (branch lists, status).
	•	reconcile/: GitHub API calls, PR lookups, caching + backoff.
	•	classify/: mapping each branch into exactly one BranchState.
	•	ui/: tview or other TUI framework code.
	•	pkg/…
(Optional) code you might publish for others—e.g. shared config loaders or generic utilities.
	•	docs/
Markdown docs:
	•	Architecture decisions, diagrams, sequence flows.
	•	Your requirements.md with the FRs you already defined.
	•	A user guide showing typical workflows and screenshots of the TUI.
	•	scripts/
Helpers for builds, releases, linting, etc. Keep CI YAML minimal by calling these.
	•	examples/ & testdata/
	•	examples/: a small dummy repo showing how to invoke the tool.
	•	testdata/: git repo fixtures for unit tests (you can check these into source).
	•	Top-level files
	•	README.md: project overview, install, usage, contribution guide.
	•	LICENSE (e.g. MIT, Apache 2.0).
	•	.gitignore: ignore bin/, *.exe, *.cache, testdata/*/.git, etc.
	•	go.mod & go.sum: your module declaration.

Getting started

# at the root of branch-wrangler/
git init
go mod init github.com/yourusername/branch-wrangler
git add .
git commit -m "chore: initial project scaffold"

From there:
	1.	Fill in requirements.md in docs/.
	2.	Implement core scanning in internal/git/scanner.go.
	3.	Add a basic main.go that prints “Hello, Branch Wrangler!” as a smoke test.
	4.	Wire up your TUI in internal/ui/tui.go.
	5.	Hook up CI in .github/workflows/ci.yml.

This layout will keep your code organized, testable, and easy to extend as you add more branch-cleanup features. Let me know if you’d like starter templates for any of these files!
