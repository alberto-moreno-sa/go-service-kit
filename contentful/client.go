package contentful

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const CMABaseURL = "https://api.contentful.com"

// Client is a base Contentful CMA client. Embed it in your own client
// to get build log and publish functionality for free.
type Client struct {
	SpaceID    string
	Token      string
	HTTPClient *http.Client
}

// NewClient creates a new Contentful CMA client.
func NewClient(spaceID, token string) *Client {
	return &Client{
		SpaceID:    spaceID,
		Token:      token,
		HTTPClient: &http.Client{},
	}
}

// PublishEntry publishes a Contentful entry.
func (c *Client) PublishEntry(ctx context.Context, entryID string, version int) error {
	endpoint := fmt.Sprintf("%s/spaces/%s/environments/master/entries/%s/published",
		CMABaseURL, c.SpaceID, entryID)

	req, err := http.NewRequestWithContext(ctx, "PUT", endpoint, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("X-Contentful-Version", fmt.Sprintf("%d", version))

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("CMA publish failed (%d): could not read body: %w", resp.StatusCode, err)
		}
		return fmt.Errorf("CMA publish failed (%d): %s", resp.StatusCode, string(body))
	}

	return nil
}

// GetEntry fetches a single Contentful entry by ID.
func (c *Client) GetEntry(ctx context.Context, entryID string) (*EntryItem, error) {
	endpoint := fmt.Sprintf("%s/spaces/%s/environments/master/entries/%s",
		CMABaseURL, c.SpaceID, entryID)

	req, err := http.NewRequestWithContext(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.Token)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("CMA get entry failed (%d): could not read body: %w", resp.StatusCode, err)
		}
		return nil, fmt.Errorf("CMA get entry failed (%d): %s", resp.StatusCode, string(body))
	}

	var entry EntryItem
	if err := json.NewDecoder(resp.Body).Decode(&entry); err != nil {
		return nil, fmt.Errorf("decode entry: %w", err)
	}

	return &entry, nil
}

// UpdateEntry updates a Contentful entry using optimistic locking.
// Returns the new version after the update.
func (c *Client) UpdateEntry(ctx context.Context, entryID string, version int, fields map[string]interface{}) (int, error) {
	endpoint := fmt.Sprintf("%s/spaces/%s/environments/master/entries/%s",
		CMABaseURL, c.SpaceID, entryID)

	body := map[string]interface{}{
		"fields": fields,
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return 0, fmt.Errorf("marshal entry body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "PUT", endpoint, bytes.NewReader(bodyBytes))
	if err != nil {
		return 0, err
	}
	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", "application/vnd.contentful.management.v1+json")
	req.Header.Set("X-Contentful-Version", fmt.Sprintf("%d", version))

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return 0, fmt.Errorf("CMA update entry failed (%d): could not read body: %w", resp.StatusCode, err)
		}
		return 0, fmt.Errorf("CMA update entry failed (%d): %s", resp.StatusCode, string(respBody))
	}

	var updated EntryItem
	if err := json.NewDecoder(resp.Body).Decode(&updated); err != nil {
		return 0, fmt.Errorf("decode update response: %w", err)
	}

	return updated.Sys.Version, nil
}
