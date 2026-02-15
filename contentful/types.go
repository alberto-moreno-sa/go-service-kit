package contentful

// BuildLogEntry represents a single execution record.
type BuildLogEntry struct {
	Service         string `json:"service"`
	Timestamp       string `json:"timestamp"`
	TriggeredBy     string `json:"triggeredBy"`
	ForceUpdate     bool   `json:"forceUpdate"`
	TranslationUsed bool   `json:"translationUsed"`
	NewAdded        int    `json:"newAdded"`
	TotalAfterSync  int    `json:"totalAfterSync"`
	Status          string `json:"status"`
}

// BuildLogResult holds the fetched build log along with entry metadata
// needed for the fetch-mutate-put update pattern.
type BuildLogResult struct {
	Entries   []BuildLogEntry
	EntryID   string
	Version   int
	RawFields map[string]interface{}
}

// EntriesResponse is the CMA response for entry queries.
type EntriesResponse struct {
	Items []EntryItem `json:"items"`
	Total int         `json:"total"`
}

// EntryItem represents a single Contentful entry.
type EntryItem struct {
	Sys    EntrySys               `json:"sys"`
	Fields map[string]interface{} `json:"fields"`
}

// EntrySys holds system metadata for an entry.
type EntrySys struct {
	ID      string `json:"id"`
	Version int    `json:"version"`
}
