
## Authentication

branch-wrangler needs a GitHub access token to call the GitHub API. You can authenticate in one of three ways. branch-wrangler will try them in order:

1. **Environment Variable**: If you already have a GitHub Personal Access Token (PAT), set it in your shell before running branch-wrangler:

    `export GITHUB_TOKEN=ghp_yourLongRandomToken`

    or

    `export GH_TOKEN=ghp_yourLongRandomToken`

    branch-wrangler will detect this and use it automatically.

2. **Config-File Token**: If you prefer a persistent setting, add your token to branch-wrangler’s config file (default location $XDG_CONFIG_HOME/branch-wrangler/config.yml, or ~/.config/branch-wrangler/config.yml):

    ```
    # ~/.config/branch-wrangler/config.yml
    token: ghp_yourLongRandomToken
    ```

    branch-wrangler will read this file and authenticate with the given token if no environment variable is set.

3. **Interactive Login (OAuth Device Flow)**: If you haven’t created a PAT, run:

    `branch-wrangler login`

    branch-wrangler will:
	    - Contact GitHub to request a device code.
	    - Prompt you to visit GitHub’s verification URL and enter a short code.
	    - Poll GitHub until you authorize the app.
	    - Receive an OAuth token and save it into your config file for future runs.

    Once authenticated, branch-wrangler will reuse the saved token until it expires or you revoke it.
