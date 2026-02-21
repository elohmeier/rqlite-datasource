package plugin

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRqliteClient_Query(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify method and path
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/db/query" {
			t.Errorf("expected /db/query, got %s", r.URL.Path)
		}

		// Verify consistency level
		level := r.URL.Query().Get("level")
		if level != "strong" {
			t.Errorf("expected level=strong, got %s", level)
		}

		// Verify content type
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type application/json, got %s", r.Header.Get("Content-Type"))
		}

		// Verify body is JSON array of strings
		body, _ := io.ReadAll(r.Body)
		var queries []string
		if err := json.Unmarshal(body, &queries); err != nil {
			t.Errorf("failed to unmarshal body: %v", err)
		}
		if len(queries) != 1 || queries[0] != "SELECT 1" {
			t.Errorf("unexpected query: %v", queries)
		}

		// Return mock response
		resp := RqliteQueryResponse{
			Results: []RqliteResult{
				{
					Columns: []string{"1"},
					Types:   []string{"integer"},
					Values:  [][]interface{}{{float64(1)}},
				},
			},
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := &RqliteClient{
		httpClient:       server.Client(),
		baseURL:          server.URL,
		consistencyLevel: "strong",
	}

	result, err := client.Query(context.Background(), "SELECT 1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(result.Results))
	}
	if len(result.Results[0].Columns) != 1 {
		t.Errorf("expected 1 column, got %d", len(result.Results[0].Columns))
	}
}

func TestRqliteClient_CheckReady(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/readyz" {
			t.Errorf("expected /readyz, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("[+]ok"))
	}))
	defer server.Close()

	client := &RqliteClient{
		httpClient:       server.Client(),
		baseURL:          server.URL,
		consistencyLevel: "weak",
	}

	err := client.CheckReady(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRqliteClient_CheckReady_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte("not ready"))
	}))
	defer server.Close()

	client := &RqliteClient{
		httpClient:       server.Client(),
		baseURL:          server.URL,
		consistencyLevel: "weak",
	}

	err := client.CheckReady(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestRqliteClient_Query_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("internal error"))
	}))
	defer server.Close()

	client := &RqliteClient{
		httpClient:       server.Client(),
		baseURL:          server.URL,
		consistencyLevel: "weak",
	}

	_, err := client.Query(context.Background(), "SELECT 1")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
