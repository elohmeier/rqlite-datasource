package plugin

import (
	"encoding/json"
	"net/http"

	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
)

// ColumnInfo represents column metadata returned by the /columns endpoint.
type ColumnInfo struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

func (d *Datasource) registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/tables", d.handleTables)
	mux.HandleFunc("/columns", d.handleColumns)
}

func (d *Datasource) handleTables(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	resp, err := d.client.Query(ctx, "SELECT name FROM sqlite_master WHERE type='table' ORDER BY name")
	if err != nil {
		log.DefaultLogger.Error("Failed to query tables", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(resp.Results) == 0 || resp.Results[0].Error != "" {
		errMsg := "no results"
		if len(resp.Results) > 0 {
			errMsg = resp.Results[0].Error
		}
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	tables := make([]string, 0, len(resp.Results[0].Values))
	for _, row := range resp.Results[0].Values {
		if len(row) > 0 {
			if name, ok := row[0].(string); ok {
				tables = append(tables, name)
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tables)
}

func (d *Datasource) handleColumns(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	table := r.URL.Query().Get("table")
	if table == "" {
		http.Error(w, "table parameter is required", http.StatusBadRequest)
		return
	}

	// Use PRAGMA table_info to get column information.
	// PRAGMA returns: cid, name, type, notnull, dflt_value, pk
	resp, err := d.client.Query(ctx, "PRAGMA table_info("+table+")")
	if err != nil {
		log.DefaultLogger.Error("Failed to query columns", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(resp.Results) == 0 || resp.Results[0].Error != "" {
		errMsg := "no results"
		if len(resp.Results) > 0 {
			errMsg = resp.Results[0].Error
		}
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

	// Find column indexes for name and type
	nameIdx := -1
	typeIdx := -1
	for i, col := range resp.Results[0].Columns {
		switch col {
		case "name":
			nameIdx = i
		case "type":
			typeIdx = i
		}
	}

	if nameIdx < 0 {
		http.Error(w, "unexpected PRAGMA result format", http.StatusInternalServerError)
		return
	}

	columns := make([]ColumnInfo, 0, len(resp.Results[0].Values))
	for _, row := range resp.Results[0].Values {
		col := ColumnInfo{}
		if nameIdx < len(row) {
			if name, ok := row[nameIdx].(string); ok {
				col.Name = name
			}
		}
		if typeIdx >= 0 && typeIdx < len(row) {
			if typ, ok := row[typeIdx].(string); ok {
				col.Type = typ
			}
		}
		columns = append(columns, col)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(columns)
}
