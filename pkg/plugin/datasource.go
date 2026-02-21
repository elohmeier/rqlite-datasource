package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/resource/httpadapter"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

var (
	_ backend.QueryDataHandler    = (*Datasource)(nil)
	_ backend.CheckHealthHandler  = (*Datasource)(nil)
	_ backend.CallResourceHandler = (*Datasource)(nil)
	_ instancemgmt.InstanceDisposer = (*Datasource)(nil)
)

// Datasource is the rqlite datasource plugin implementation.
type Datasource struct {
	client          *RqliteClient
	resourceHandler backend.CallResourceHandler
	settings        PluginSettings
}

// NewDatasource creates a new datasource instance.
func NewDatasource(ctx context.Context, settings backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	var pluginSettings PluginSettings
	if err := json.Unmarshal(settings.JSONData, &pluginSettings); err != nil {
		return nil, fmt.Errorf("unmarshaling settings: %w", err)
	}

	httpOpts, err := settings.HTTPClientOptions(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting HTTP client options: %w", err)
	}

	client, err := NewRqliteClient(settings.URL, pluginSettings.ConsistencyLevel, httpOpts)
	if err != nil {
		return nil, fmt.Errorf("creating rqlite client: %w", err)
	}

	ds := &Datasource{
		client:   client,
		settings: pluginSettings,
	}

	mux := http.NewServeMux()
	ds.registerRoutes(mux)
	ds.resourceHandler = httpadapter.New(mux)

	return ds, nil
}

// Dispose cleans up resources.
func (d *Datasource) Dispose() {}

// QueryData handles multiple queries and returns multiple responses.
func (d *Datasource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	response := backend.NewQueryDataResponse()

	for _, q := range req.Queries {
		res := d.query(ctx, q)
		response.Responses[q.RefID] = res
	}

	return response, nil
}

func (d *Datasource) query(ctx context.Context, query backend.DataQuery) backend.DataResponse {
	var qm QueryModel
	if err := json.Unmarshal(query.JSON, &qm); err != nil {
		return backend.ErrDataResponse(backend.StatusBadRequest, fmt.Sprintf("json unmarshal: %v", err))
	}

	rawSQL := qm.RawSQL
	if rawSQL == "" {
		return backend.ErrDataResponse(backend.StatusBadRequest, "query is empty")
	}

	// Apply macros
	rawSQL = ApplyMacros(rawSQL, query.TimeRange, query.Interval.Milliseconds())

	// Execute query
	result, err := d.client.Query(ctx, rawSQL)
	if err != nil {
		return backend.ErrDataResponse(backend.StatusInternal, fmt.Sprintf("query execution: %v", err))
	}

	if len(result.Results) == 0 {
		return backend.DataResponse{}
	}

	// Convert to data frame
	frame, err := ResultToFrame(&result.Results[0], qm.TimeColumns)
	if err != nil {
		return backend.ErrDataResponse(backend.StatusInternal, fmt.Sprintf("converting result: %v", err))
	}

	if qm.Format == "time_series" {
		frame.Meta = &data.FrameMeta{
			Type: data.FrameTypeTimeSeriesWide,
		}
	}

	return backend.DataResponse{Frames: data.Frames{frame}}
}

// CheckHealth handles health checks.
func (d *Datasource) CheckHealth(ctx context.Context, _ *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	if err := d.client.CheckReady(ctx); err != nil {
		return &backend.CheckHealthResult{
			Status:  backend.HealthStatusError,
			Message: fmt.Sprintf("rqlite health check failed: %v", err),
		}, nil
	}

	return &backend.CheckHealthResult{
		Status:  backend.HealthStatusOk,
		Message: "rqlite is ready",
	}, nil
}

// CallResource handles resource calls for the visual query builder.
func (d *Datasource) CallResource(ctx context.Context, req *backend.CallResourceRequest, sender backend.CallResourceResponseSender) error {
	return d.resourceHandler.CallResource(ctx, req, sender)
}
