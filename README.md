# rqlite Datasource for Grafana

A Grafana backend datasource plugin for [rqlite](https://rqlite.io), the lightweight, distributed relational database built on SQLite.

## Features

- SQL query editor with syntax highlighting
- Visual query builder (table, column, WHERE, GROUP BY, ORDER BY, LIMIT)
- Time series and table format support
- Grafana macros: `$__timeFilter`, `$__timeFrom`, `$__timeTo`, `$__timeGroup`, `$__unixEpochFilter`
- Configurable [consistency level](https://rqlite.io/docs/api/read-consistency/) (none, weak, strong, linearizable)
- HTTP Basic Auth support
- Grafana alerting support

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
