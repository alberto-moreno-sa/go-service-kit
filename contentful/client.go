package contentful

import (
	"context"
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
