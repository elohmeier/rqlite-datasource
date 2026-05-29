# rqlite Datasource for Grafana

[![CI](https://github.com/elohmeier/rqlite-datasource/actions/workflows/ci.yml/badge.svg)](https://github.com/elohmeier/rqlite-datasource/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/elohmeier/rqlite-datasource)](https://goreportcard.com/report/github.com/elohmeier/rqlite-datasource)
[![License](https://img.shields.io/github/license/elohmeier/rqlite-datasource)](https://github.com/elohmeier/rqlite-datasource/blob/main/LICENSE)
[![Grafana Plugin](https://img.shields.io/badge/dynamic/json?logo=grafana&query=$.version&url=https://grafana.com/api/plugins/g42-rqlite-datasource&label=Grafana%20plugin&prefix=v&color=F47A20)](https://grafana.com/grafana/plugins/g42-rqlite-datasource/)

A Grafana backend datasource plugin for [rqlite](https://rqlite.io), the lightweight, distributed relational database built on SQLite.

Use this plugin to query rqlite clusters from Grafana dashboards, explore SQL results as tables or time series, use dashboard variables, and run Grafana-managed alerts against rqlite data.

## Features

- SQL query editor with syntax highlighting
- Visual query builder (table, column, WHERE, GROUP BY, ORDER BY, LIMIT)
- Time series and table format support
- Grafana macros: `$__timeFilter`, `$__timeFrom`, `$__timeTo`, `$__timeGroup`, `$__unixEpochFilter`
- Dashboard variable query support
- Configurable [consistency level](https://rqlite.io/docs/api/read-consistency/) (none, weak, strong, linearizable)
- HTTP Basic Auth support
- Grafana alerting support

## Requirements

- Grafana 12.3.0 or newer
- A reachable [rqlite](https://rqlite.io) node or cluster with the HTTP API enabled

## Install

Install the plugin from the [Grafana plugin catalog](https://grafana.com/grafana/plugins/g42-rqlite-datasource/) or with the Grafana CLI:

```bash
grafana cli plugins install g42-rqlite-datasource
```

Restart Grafana after installing or upgrading the plugin.

## Configure

1. In Grafana, go to **Connections > Data sources** and add the **rqlite** data source.
2. Set the URL to your rqlite HTTP endpoint, for example `http://localhost:4001`.
3. Configure authentication if your rqlite node requires it.
4. Optionally set the rqlite read consistency level and query timeout under **Additional settings**.
5. Click **Save & test**.

## Query

Use code mode for raw SQLite-compatible SQL:

```sql
SELECT
  time,
  value
FROM metrics
WHERE $__timeFilter(time)
ORDER BY time
```

Use builder mode to select a table, columns, filters, grouping, ordering, limits, and offsets without writing SQL manually.

For time series panels, set the query format to **Time series** and list any time columns in the query editor. Time columns can contain Unix timestamps or common string formats such as RFC3339 and `YYYY-MM-DD HH:MM:SS`.

## Macros

| Macro | Output |
| --- | --- |
| `$__timeFilter(column)` | `column >= <from> AND column <= <to>` |
| `$__unixEpochFilter(column)` | Alias for `$__timeFilter(column)` |
| `$__timeFrom` | Dashboard range start as Unix epoch seconds |
| `$__timeTo` | Dashboard range end as Unix epoch seconds |
| `$__timeGroup(column, 5m)` | SQLite-compatible epoch bucket expression |

## Links

- [Grafana plugin catalog](https://grafana.com/grafana/plugins/g42-rqlite-datasource/)
- [GitHub repository](https://github.com/elohmeier/rqlite-datasource)
- [rqlite documentation](https://rqlite.io/docs/)

## Development

### Prerequisites

- [Go](https://go.dev/) 1.23+
- [Mage](https://magefile.org/)
- [Node.js](https://nodejs.org/) 20+
- [pnpm](https://pnpm.io/)
- [Docker](https://www.docker.com/)

### Build

```bash
# Frontend
pnpm install
pnpm run build

# Backend (all platforms)
mage -v

# Backend (specific platform)
mage build:linuxARM64
```

### Run

```bash
docker compose up
```

This starts Grafana on http://localhost:3000 with a pre-configured rqlite datasource, and an rqlite node on http://localhost:4001.

### Test

```bash
# Backend
go test ./pkg/...

# Frontend
pnpm run test:ci

# Lint
pnpm run lint
```
