package github

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/go-github/v68/github"
	"golang.org/x/oauth2"
)

const (
	GitHubClientID = "Ov23liOMSfFBJ3w6rK1U"
	DeviceCodeURL  = "https://github.com/login/device/code"
	TokenURL       = "https://github.com/login/oauth/access_token"
)

type AuthConfig struct {
	Token    string
	TokenEnv string
	Config   *oauth2.Config
}

func NewAuthConfig() *AuthConfig {
	return &AuthConfig{
		Config: &oauth2.Config{
			ClientID: GitHubClientID,
			Scopes:   []string{"repo", "workflow"},
			Endpoint: oauth2.Endpoint{
				AuthURL:       "https://github.com/login/oauth/authorize",
				TokenURL:      TokenURL,
				DeviceAuthURL: DeviceCodeURL,
			},
		},
	}
}

func (a *AuthConfig) GetToken() (string, error) {
	if token := os.Getenv("GITHUB_TOKEN"); token != "" {
		return token, nil
	}

	return "", fmt.Errorf("GITHUB_TOKEN environment variable not set - please set it to your GitHub Personal Access Token")
}

func (a *AuthConfig) DeviceFlow(ctx context.Context) (string, error) {
	return "", fmt.Errorf("device flow not implemented - please set GITHUB_TOKEN environment variable")
}

func (a *AuthConfig) SaveToken(token string) error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	return writeTokenToConfig(configPath, token)
}

func (a *AuthConfig) ValidateToken(token string) error {
	client := github.NewClient(oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)))

	_, _, err := client.RateLimit.Get(context.Background())
	return err
}

func getConfigPath() (string, error) {
	configDir := os.Getenv("XDG_CONFIG_HOME")
	if configDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		configDir = filepath.Join(homeDir, ".config")
	}

	appConfigDir := filepath.Join(configDir, "branch-wrangler")
	if err := os.MkdirAll(appConfigDir, 0700); err != nil {
		return "", err
	}

	return filepath.Join(appConfigDir, "config.yml"), nil
}

func writeTokenToConfig(configPath, token string) error {
	return fmt.Errorf("config file writing not implemented yet")
}
