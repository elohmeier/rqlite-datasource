package plugin

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/grafana/grafana-plugin-sdk-go/backend/httpclient"
)

// RqliteClient wraps HTTP communication with a rqlite cluster.
type RqliteClient struct {
	httpClient       *http.Client
	baseURL          string
	consistencyLevel string
}

// NewRqliteClient creates a new RqliteClient using Grafana's HTTP client provider.
func NewRqliteClient(baseURL, consistencyLevel string, opts ...httpclient.Options) (*RqliteClient, error) {
	var httpOpts httpclient.Options
	if len(opts) > 0 {
		httpOpts = opts[0]
	}

	client, err := httpclient.New(httpOpts)
	if err != nil {
		return nil, fmt.Errorf("creating HTTP client: %w", err)
	}

	baseURL = strings.TrimRight(baseURL, "/")

	if consistencyLevel == "" {
		consistencyLevel = "weak"
	}

	return &RqliteClient{
		httpClient:       client,
		baseURL:          baseURL,
		consistencyLevel: consistencyLevel,
	}, nil
}

// Query executes a SQL query against rqlite and returns the response.
func (c *RqliteClient) Query(ctx context.Context, sql string) (*RqliteQueryResponse, error) {
	body, err := json.Marshal([]string{sql})
	if err != nil {
		return nil, fmt.Errorf("marshaling query: %w", err)
	}

	url := fmt.Sprintf("%s/db/query?level=%s", c.baseURL, c.consistencyLevel)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing query: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("rqlite returned status %d: %s", resp.StatusCode, string(respBody))
	}

	var result RqliteQueryResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("unmarshaling response: %w", err)
	}

	return &result, nil
}

// CheckReady checks if the rqlite node is ready.
func (c *RqliteClient) CheckReady(ctx context.Context) error {
	url := fmt.Sprintf("%s/readyz", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("creating readiness request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("checking readiness: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("rqlite not ready (status %d): %s", resp.StatusCode, string(body))
	}

	return nil
}
