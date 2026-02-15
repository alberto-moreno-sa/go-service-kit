# go-service-kit

Shared Go SDK for common service utilities. Currently provides a Contentful CMA client with build log tracking.

## Install

```bash
go get github.com/alberto-moreno-sa/go-service-kit
```

## Packages

### `contentful`

Base Contentful CMA client with build log support.

```go
import "github.com/alberto-moreno-sa/go-service-kit/contentful"
```

#### Client

```go
client := contentful.NewClient(spaceID, token)
```

The client exposes `SpaceID`, `Token`, and `HTTPClient` fields so you can embed it in your own client and build domain-specific methods on top.

#### Build Log

Track execution metadata in a Contentful `buildLog` content type:

```go
// Fetch existing log
result, _ := client.GetBuildLog(ctx)

// Append entry
entries := append(result.Entries, contentful.BuildLogEntry{
    Service:     "my-service",
    Timestamp:   time.Now().UTC().Format(time.RFC3339),
    TriggeredBy: "github-actions",
    Status:      "success",
})

// Update and publish
version, _ := client.UpdateBuildLog(ctx, result, entries)
client.PublishEntry(ctx, entryID, version)
```

#### Publish Entry

```go
err := client.PublishEntry(ctx, entryID, version)
```

#### Embedding

Embed the SDK client in your own to get build log and publish for free:

```go
type MyClient struct {
    *contentful.Client
}

func NewMyClient(spaceID, token string) *MyClient {
    return &MyClient{Client: contentful.NewClient(spaceID, token)}
}

// Add your own domain-specific methods
func (c *MyClient) GetItems(ctx context.Context) { ... }
```

## Contentful Setup

The build log requires a `buildLog` content model in Contentful with a `logInfo` JSON field.

### BuildLogEntry fields

| Field | Type | Description |
|---|---|---|
| `service` | string | Service that generated the log |
| `timestamp` | string | ISO 8601 timestamp |
| `triggeredBy` | string | `"local"` or `"github-actions"` |
| `forceUpdate` | bool | Whether force mode was used |
| `translationUsed` | bool | Whether translation was used |
| `newAdded` | int | Items added in this run |
| `totalAfterSync` | int | Total items after sync |
| `status` | string | `"success"` |

## License

[MIT](LICENSE)
