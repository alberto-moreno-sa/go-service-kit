package contentful

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// GetBuildLog fetches the build log entry.
func (c *Client) GetBuildLog(ctx context.Context) (*BuildLogResult, error) {
	endpoint := fmt.Sprintf("%s/spaces/%s/environments/master/entries", CMABaseURL, c.SpaceID)

	params := url.Values{}
	params.Set("content_type", "buildLog")
	params.Set("limit", "1")

	req, err := http.NewRequestWithContext(ctx, "GET", endpoint+"?"+params.Encode(), nil)
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
			return nil, fmt.Errorf("CMA build log query failed (%d): could not read body: %w", resp.StatusCode, err)
		}
		return nil, fmt.Errorf("CMA build log query failed (%d): %s", resp.StatusCode, string(body))
	}

	var result EntriesResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode build log response: %w", err)
	}

	if len(result.Items) == 0 {
		return &BuildLogResult{}, nil
	}

	entry := result.Items[0]

	logField, ok := entry.Fields["logInfo"]
	if !ok {
		return &BuildLogResult{
			EntryID:   entry.Sys.ID,
			Version:   entry.Sys.Version,
			RawFields: entry.Fields,
		}, nil
	}

	localeMap, ok := logField.(map[string]interface{})
	if !ok {
		return &BuildLogResult{
			EntryID:   entry.Sys.ID,
			Version:   entry.Sys.Version,
			RawFields: entry.Fields,
		}, nil
	}

	rawContent, ok := localeMap["en-US"]
	if !ok {
		for _, v := range localeMap {
			rawContent = v
			break
		}
	}

	contentBytes, err := json.Marshal(rawContent)
	if err != nil {
		return nil, fmt.Errorf("marshal build log content: %w", err)
	}

	var entries []BuildLogEntry
	if err := json.Unmarshal(contentBytes, &entries); err != nil {
		return nil, fmt.Errorf("unmarshal build log entries: %w", err)
	}

	return &BuildLogResult{
		Entries:   entries,
		EntryID:   entry.Sys.ID,
		Version:   entry.Sys.Version,
		RawFields: entry.Fields,
	}, nil
}

// UpdateBuildLog updates the build log entry using the fetch-mutate-put pattern.
func (c *Client) UpdateBuildLog(ctx context.Context, result *BuildLogResult, entries []BuildLogEntry) (int, error) {
	endpoint := fmt.Sprintf("%s/spaces/%s/environments/master/entries/%s",
		CMABaseURL, c.SpaceID, result.EntryID)

	fields := make(map[string]interface{})
	for k, v := range result.RawFields {
		fields[k] = v
	}
	fields["logInfo"] = map[string]interface{}{
		"en-US": entries,
	}

	body := map[string]interface{}{
		"fields": fields,
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return 0, fmt.Errorf("marshal build log body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "PUT", endpoint, bytes.NewReader(bodyBytes))
	if err != nil {
		return 0, err
	}
	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", "application/vnd.contentful.management.v1+json")
	req.Header.Set("X-Contentful-Version", fmt.Sprintf("%d", result.Version))

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return 0, fmt.Errorf("CMA build log update failed (%d): could not read body: %w", resp.StatusCode, err)
		}
		return 0, fmt.Errorf("CMA build log update failed (%d): %s", resp.StatusCode, string(respBody))
	}

	var updated EntryItem
	if err := json.NewDecoder(resp.Body).Decode(&updated); err != nil {
		return 0, fmt.Errorf("decode build log update response: %w", err)
	}

	return updated.Sys.Version, nil
}

// CreateBuildLog creates a new buildLog entry.
func (c *Client) CreateBuildLog(ctx context.Context, entries []BuildLogEntry) (string, int, error) {
	endpoint := fmt.Sprintf("%s/spaces/%s/environments/master/entries", CMABaseURL, c.SpaceID)

	body := map[string]interface{}{
		"fields": map[string]interface{}{
			"logInfo": map[string]interface{}{"en-US": entries},
		},
	}

	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return "", 0, fmt.Errorf("marshal build log body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewReader(bodyBytes))
	if err != nil {
		return "", 0, err
	}
	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", "application/vnd.contentful.management.v1+json")
	req.Header.Set("X-Contentful-Content-Type", "buildLog")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", 0, fmt.Errorf("CMA build log create failed (%d): could not read body: %w", resp.StatusCode, err)
		}
		return "", 0, fmt.Errorf("CMA build log create failed (%d): %s", resp.StatusCode, string(respBody))
	}

	var created EntryItem
	if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
		return "", 0, fmt.Errorf("decode build log create response: %w", err)
	}

	return created.Sys.ID, created.Sys.Version, nil
}
