package plugin

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func setupTestDatasource(t *testing.T, rqliteHandler http.HandlerFunc) (*Datasource, *httptest.Server) {
	t.Helper()

	rqliteServer := httptest.NewServer(rqliteHandler)

	client := &RqliteClient{
		httpClient:       rqliteServer.Client(),
		baseURL:          rqliteServer.URL,
		consistencyLevel: "weak",
	}

	ds := &Datasource{
		client:   client,
		settings: PluginSettings{ConsistencyLevel: "weak"},
	}

	return ds, rqliteServer
}

func TestHandleTables(t *testing.T) {
	ds, rqliteServer := setupTestDatasource(t, func(w http.ResponseWriter, r *http.Request) {
		resp := RqliteQueryResponse{
			Results: []RqliteResult{
				{
					Columns: []string{"name"},
					Types:   []string{"text"},
					Values: [][]interface{}{
						{"users"},
						{"orders"},
						{"products"},
					},
				},
			},
		}
		_ = json.NewEncoder(w).Encode(resp)
	})
	defer rqliteServer.Close()

	req := httptest.NewRequest(http.MethodGet, "/tables", nil)
	rec := httptest.NewRecorder()
	ds.handleTables(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var tables []string
	if err := json.NewDecoder(rec.Body).Decode(&tables); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(tables) != 3 {
		t.Fatalf("expected 3 tables, got %d", len(tables))
	}
	if tables[0] != "users" {
		t.Errorf("expected first table 'users', got %q", tables[0])
	}
}

func TestHandleColumns(t *testing.T) {
	ds, rqliteServer := setupTestDatasource(t, func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("q"); got != "" {
			t.Fatalf("unexpected query string %q", got)
		}

		resp := RqliteQueryResponse{
			Results: []RqliteResult{
				{
					Columns: []string{"cid", "name", "type", "notnull", "dflt_value", "pk"},
					Types:   []string{"integer", "text", "text", "integer", "text", "integer"},
					Values: [][]interface{}{
						{float64(0), "id", "INTEGER", float64(1), nil, float64(1)},
						{float64(1), "name", "TEXT", float64(0), nil, float64(0)},
						{float64(2), "email", "TEXT", float64(0), nil, float64(0)},
					},
				},
			},
		}
		_ = json.NewEncoder(w).Encode(resp)
	})
	defer rqliteServer.Close()

	req := httptest.NewRequest(http.MethodGet, "/columns?table=users", nil)
	rec := httptest.NewRecorder()
	ds.handleColumns(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	var columns []ColumnInfo
	if err := json.NewDecoder(rec.Body).Decode(&columns); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(columns) != 3 {
		t.Fatalf("expected 3 columns, got %d", len(columns))
	}
	if columns[0].Name != "id" || columns[0].Type != "INTEGER" {
		t.Errorf("expected id/INTEGER, got %s/%s", columns[0].Name, columns[0].Type)
	}
}

func TestHandleColumns_MissingTable(t *testing.T) {
	ds := &Datasource{}

	req := httptest.NewRequest(http.MethodGet, "/columns", nil)
	rec := httptest.NewRecorder()
	ds.handleColumns(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

func TestHandleColumns_RejectsUnsafeTableName(t *testing.T) {
	ds := &Datasource{}

	req := httptest.NewRequest(http.MethodGet, "/columns?table=users%3BDROP%20TABLE%20users", nil)
	rec := httptest.NewRecorder()
	ds.handleColumns(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}

	if !strings.Contains(rec.Body.String(), "invalid table parameter") {
		t.Fatalf("unexpected response body: %q", rec.Body.String())
	}
}

func TestHandleColumns_QuotesValidatedTableName(t *testing.T) {
	ds, rqliteServer := setupTestDatasource(t, func(w http.ResponseWriter, r *http.Request) {
		var queries []string
		if err := json.NewDecoder(r.Body).Decode(&queries); err != nil {
			t.Fatalf("failed to decode request body: %v", err)
		}
		if len(queries) != 1 {
			t.Fatalf("expected 1 query, got %d", len(queries))
		}
		if queries[0] != `PRAGMA table_info("users")` {
			t.Fatalf("unexpected query: %q", queries[0])
		}

		resp := RqliteQueryResponse{
			Results: []RqliteResult{
				{
					Columns: []string{"cid", "name", "type"},
					Types:   []string{"integer", "text", "text"},
					Values:  [][]interface{}{{float64(0), "id", "INTEGER"}},
				},
			},
		}
		_ = json.NewEncoder(w).Encode(resp)
	})
	defer rqliteServer.Close()

	req := httptest.NewRequest(http.MethodGet, "/columns?table=users", nil)
	rec := httptest.NewRecorder()
	ds.handleColumns(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
}

func TestHandleTables_HidesBackendErrorDetails(t *testing.T) {
	ds, rqliteServer := setupTestDatasource(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("dial tcp 10.0.0.12:4001: connect: connection refused"))
	})
	defer rqliteServer.Close()

	req := httptest.NewRequest(http.MethodGet, "/tables", nil)
	rec := httptest.NewRecorder()
	ds.handleTables(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rec.Code)
	}
	if strings.Contains(rec.Body.String(), "10.0.0.12") {
		t.Fatalf("backend details leaked to client: %q", rec.Body.String())
	}
	if !strings.Contains(rec.Body.String(), genericQueryErrorMessage) {
		t.Fatalf("expected generic error message, got %q", rec.Body.String())
	}
}
