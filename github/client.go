package github

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const apiBaseURL = "https://api.github.com"

// Client is a GitHub REST API client.
type Client struct {
	Token      string
	HTTPClient *http.Client
}

// NewClient creates a new GitHub API client. Token is optional but recommended for higher rate limits.
func NewClient(token string) *Client {
	return &Client{
		Token:      token,
		HTTPClient: &http.Client{},
	}
}

func (c *Client) newRequest(ctx context.Context, method, url string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}
	return req, nil
}

// ListRepos returns all public, non-fork, non-archived repositories for a user.
func (c *Client) ListRepos(ctx context.Context, username string) ([]Repo, error) {
	url := fmt.Sprintf("%s/users/%s/repos?type=public&sort=updated&per_page=100", apiBaseURL, username)

	req, err := c.newRequest(ctx, "GET", url)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("GitHub list repos failed (%d): could not read body: %w", resp.StatusCode, err)
		}
		return nil, fmt.Errorf("GitHub list repos failed (%d): %s", resp.StatusCode, string(body))
	}

	var repos []Repo
	if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		return nil, fmt.Errorf("decode repos: %w", err)
	}

	return repos, nil
}

// GetRepoLanguages returns the language breakdown for a repository.
func (c *Client) GetRepoLanguages(ctx context.Context, owner, repo string) (map[string]int, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/languages", apiBaseURL, owner, repo)

	req, err := c.newRequest(ctx, "GET", url)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("GitHub languages failed (%d): could not read body: %w", resp.StatusCode, err)
		}
		return nil, fmt.Errorf("GitHub languages failed (%d): %s", resp.StatusCode, string(body))
	}

	var languages map[string]int
	if err := json.NewDecoder(resp.Body).Decode(&languages); err != nil {
		return nil, fmt.Errorf("decode languages: %w", err)
	}

	return languages, nil
}

// GetRepoREADME returns the raw README content for a repository.
// Returns an empty string (not an error) if the repository has no README.
func (c *Client) GetRepoREADME(ctx context.Context, owner, repo string) (string, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/readme", apiBaseURL, owner, repo)

	req, err := c.newRequest(ctx, "GET", url)
	if err != nil {
		return "", err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return "", nil
	}

	if resp.StatusCode != 200 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("GitHub readme failed (%d): could not read body: %w", resp.StatusCode, err)
		}
		return "", fmt.Errorf("GitHub readme failed (%d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		Content  string `json:"content"`
		Encoding string `json:"encoding"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decode readme: %w", err)
	}

	if result.Encoding == "base64" {
		decoded, err := base64.StdEncoding.DecodeString(result.Content)
		if err != nil {
			return "", fmt.Errorf("decode base64 readme: %w", err)
		}
		return string(decoded), nil
	}

	return result.Content, nil
}
