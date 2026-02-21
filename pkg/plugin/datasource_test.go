package plugin

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
)

func TestDatasource_QueryData(t *testing.T) {
	rqliteServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := RqliteQueryResponse{
			Results: []RqliteResult{
				{
					Columns: []string{"time", "value"},
					Types:   []string{"integer", "real"},
					Values: [][]interface{}{
						{float64(1700000000), float64(42.5)},
						{float64(1700000060), float64(43.1)},
					},
				},
			},
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer rqliteServer.Close()

	ds := &Datasource{
		client: &RqliteClient{
			httpClient:       rqliteServer.Client(),
			baseURL:          rqliteServer.URL,
			consistencyLevel: "weak",
		},
		settings: PluginSettings{ConsistencyLevel: "weak"},
	}

	qm := QueryModel{
		RawSQL:      "SELECT time, value FROM metrics WHERE $__timeFilter(time)",
		Format:      "time_series",
		TimeColumns: []string{"time"},
	}
	qmJSON, _ := json.Marshal(qm)

	req := &backend.QueryDataRequest{
		Queries: []backend.DataQuery{
			{
				RefID: "A",
				JSON:  qmJSON,
				TimeRange: backend.TimeRange{
					From: time.Unix(1700000000, 0),
					To:   time.Unix(1700000120, 0),
				},
			},
		},
	}

	resp, err := ds.QueryData(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	resA, ok := resp.Responses["A"]
	if !ok {
		t.Fatal("expected response for refID 'A'")
	}

	if resA.Error != nil {
		t.Fatalf("unexpected error in response: %v", resA.Error)
	}

	if len(resA.Frames) != 1 {
		t.Fatalf("expected 1 frame, got %d", len(resA.Frames))
	}

	frame := resA.Frames[0]
	if len(frame.Fields) != 2 {
		t.Fatalf("expected 2 fields, got %d", len(frame.Fields))
	}

	if frame.Fields[0].Len() != 2 {
		t.Errorf("expected 2 rows, got %d", frame.Fields[0].Len())
	}
}

func TestDatasource_QueryData_EmptySQL(t *testing.T) {
	ds := &Datasource{
		client: &RqliteClient{
			httpClient:       http.DefaultClient,
			baseURL:          "http://localhost:1234",
			consistencyLevel: "weak",
		},
	}

	qm := QueryModel{RawSQL: ""}
	qmJSON, _ := json.Marshal(qm)

	req := &backend.QueryDataRequest{
		Queries: []backend.DataQuery{
			{
				RefID: "A",
				JSON:  qmJSON,
			},
		},
	}

	resp, err := ds.QueryData(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	resA := resp.Responses["A"]
	if resA.Error == nil {
		t.Fatal("expected error for empty SQL")
	}
}

func TestDatasource_CheckHealth(t *testing.T) {
	rqliteServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("[+]ok"))
	}))
	defer rqliteServer.Close()

	ds := &Datasource{
		client: &RqliteClient{
			httpClient:       rqliteServer.Client(),
			baseURL:          rqliteServer.URL,
			consistencyLevel: "weak",
		},
	}

	result, err := ds.CheckHealth(context.Background(), &backend.CheckHealthRequest{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != backend.HealthStatusOk {
		t.Errorf("expected OK status, got %v: %s", result.Status, result.Message)
	}
	if result.Message != "rqlite is ready" {
		t.Errorf("unexpected message: %s", result.Message)
	}
}

func TestDatasource_CheckHealth_Error(t *testing.T) {
	rqliteServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
		_, _ = w.Write([]byte("not ready"))
	}))
	defer rqliteServer.Close()

	ds := &Datasource{
		client: &RqliteClient{
			httpClient:       rqliteServer.Client(),
			baseURL:          rqliteServer.URL,
			consistencyLevel: "weak",
		},
	}

	result, err := ds.CheckHealth(context.Background(), &backend.CheckHealthRequest{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Status != backend.HealthStatusError {
		t.Errorf("expected Error status, got %v", result.Status)
	}
}
