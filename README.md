# Branch Wrangler

A terminal UI tool for cleaning up stale local Git branches by reconciling with GitHub.

## Features
- Scan all local branches and detect state (ahead/behind, merged, orphaned).
- Lookup GitHub PR status (open, draft, merged, closed).
- Interactive TUI for one-key branch cleanup.
- Dry-run mode & safe-delete options.

## Quickstart

```bash
git clone git@github.com:yourusername/branch-wrangler.git
cd branch-wrangler
go build ./cmd/branch-wrangler
./branch-wrangler --help
```

See **docs/user-guide.md** for full usage.