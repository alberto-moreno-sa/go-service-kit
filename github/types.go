package github

import "time"

// Repo represents a GitHub repository from the REST API.
type Repo struct {
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	HTMLURL     string    `json:"html_url"`
	Homepage    *string   `json:"homepage"`
	Language    *string   `json:"language"`
	Topics      []string  `json:"topics"`
	Fork        bool      `json:"fork"`
	Archived    bool      `json:"archived"`
	Size        int       `json:"size"`
	UpdatedAt   time.Time `json:"updated_at"`
	PushedAt    time.Time `json:"pushed_at"`
}
