# Datafrost

A native desktop database GUI for browsing tables, running SQL, and exploring schemas.

---

## User guide

### Install

**macOS (Apple Silicon — M1/M2/M3+)**

Download the latest `.tar.gz` from [GitHub Releases](https://github.com/3-lines-studio/datafrost/releases).

```bash
tar -xzf datafrost-macos-arm64-v*.tar.gz
sudo cp datafrost-macos-arm64-v*/datafrost /usr/local/bin/
chmod +x /usr/local/bin/datafrost
```

Or via Homebrew (when available in the tap):

```bash
brew tap 3-lines-studio/tap
brew install datafrost
```

On first launch, macOS may block the app because it is unsigned. Go to **System Settings → Privacy & Security** and click **Open Anyway**, or right-click the binary and choose **Open**.

Intel Macs are not covered by pre-built releases — see [Building from source](#building-from-source).

**Linux**

Install WebKit2GTK and GTK 3 first:

```bash
# Ubuntu/Debian
sudo apt-get install libwebkit2gtk-4.0-37 libgtk-3-0

# Fedora/RHEL
sudo dnf install webkit2gtk3 gtk3

# Arch
sudo pacman -S webkit2gtk gtk3
```

Download the Linux package from [GitHub Releases](https://github.com/3-lines-studio/datafrost/releases), extract it, and follow the included `README.txt` to add a desktop shortcut.

### Run

```bash
datafrost
```

Check the installed version:

```bash
datafrost --version
```

### Connect to a database

1. Click **+** in the sidebar to add a connection.
2. Choose a database type and fill in the credentials.
3. Click **Test** to verify, then **Save**.

| Database | What you need |
|----------|---------------|
| **SQLite** | Path to a local `.db` file |
| **Turso** | Database URL (`libsql://...`) and auth token |
| **PostgreSQL** | Connection URL *or* host, port, database, username, password, and SSL mode |
| **BigQuery** | Project ID, dataset, and service account JSON (paste or upload a `.json` file) |

Credentials are stored locally on your machine — nothing is sent to external servers.

> **Read-only.** Datafrost only runs `SELECT`, `WITH`, and (for SQLite/Turso) `PRAGMA` queries. You cannot insert, update, or delete data through the app.

### Browse tables

1. Click a saved connection in the sidebar to connect.
2. Expand the **Tables** list and click a table name to open it in a new tab.
3. Use the filter bar to narrow rows by column (`=`, `!=`, `>`, `<`, `LIKE`, `IS NULL`, etc.).
4. Paginate through results at the bottom of the table view.
5. Right-click a table (or use the menu) to open its **Schema** tab — columns, indexes, and constraints.

### Run SQL queries

1. With a connection active, click **New Query** (or press `Cmd/Ctrl + T`).
2. Write SQL in the editor and press **Run** (or `Cmd/Ctrl + Enter`).
3. Results appear below the editor. Use **Copy as CSV** or **Copy as JSON** to export them.
4. Press **Save** to store the query under that connection for later reuse.

Saved queries appear in the sidebar. Click one to open it; use the menu to rename or delete.

### Keyboard shortcuts

| Shortcut | Action |
|----------|--------|
| `Cmd/Ctrl + T` | New query tab |
| `Cmd/Ctrl + W` | Close active tab |
| `Cmd/Ctrl + Enter` | Run query |
| `Shift + Alt + F` | Format SQL |

On macOS, standard Edit menu shortcuts (Undo, Cut, Copy, Paste) work in the query editor.

### Theme and layout

Toggle light/dark mode from the sidebar. The sidebar width is resizable — drag the handle and your preference is saved automatically. Open tabs are restored the next time you connect to the same database.

### Reset or troubleshoot

**Start fresh** — removes all saved connections, queries, and preferences:

```bash
datafrost reset
```

**Config file location** (if you need to back up or delete it manually):

```
~/.config/datafrost/config.db                          # Linux
~/Library/Application Support/datafrost/config.db      # macOS
```

**Linux won't start** — confirm WebKit2GTK is installed (see [Install](#install) above) and run `datafrost` from a terminal to read error messages.

**macOS "developer cannot be verified"** — right-click the binary → **Open**, then confirm.

---

## Features

- **Multiple databases** — SQLite, Turso, PostgreSQL, and BigQuery in one app
- **Connection management** — Create, edit, delete, and test connections
- **Table browser** — Paginated views with column filters
- **SQL editor** — Syntax highlighting, formatting, and saved queries
- **Schema viewer** — Columns, indexes, and constraints per table
- **Tabbed workspace** — Query, table, and schema tabs restored per connection
- **Export** — Copy results as CSV or JSON
- **Light / dark theme** with a resizable sidebar

---

## For developers

Built with Go and React, packaged as a lightweight native app via [Bifrost](https://github.com/3-lines-studio/bifrost) and [webview](https://github.com/webview/webview).

### Prerequisites

- [Go](https://go.dev/) 1.25+
- [Bun](https://bun.sh/)
- **macOS:** Xcode command-line tools (CGO / native menu)
- **Linux:** `libgtk-3-dev`, `libwebkit2gtk-4.0-dev`

### Setup

```bash
git clone https://github.com/3-lines-studio/datafrost.git
cd datafrost
bun install
make doctor   # verify environment
```

### Commands

| Command | Description |
|---------|-------------|
| `make dev` | Hot reload (`BIFROST_DEV=1` + [Air](https://github.com/air-verse/air)) |
| `make build` | Bifrost frontend build + Go binary → `./tmp/app` |
| `make start` | Build and run |
| `make reset` | Reset local config database |

### Building from source

```bash
bun install
go run github.com/3-lines-studio/bifrost/cmd/build@latest .
go build -o datafrost .
./datafrost
```

### Architecture

```
internal/
├── core/entity/          # Domain models
├── usecase/              # Business logic
└── adapter/
    ├── database/         # SQLite, Turso, Postgres, BigQuery
    ├── http/             # REST API (Chi)
    └── repository/       # Local config.db persistence
```

```
┌─────────────────────────────────────────┐
│  webview (native window)                │
│  ┌───────────────────────────────────┐  │
│  │  React UI (Bifrost + Tailwind)    │  │
│  └──────────────┬────────────────────┘  │
│                 │ REST /api/*           │
│  ┌──────────────▼────────────────────┐  │
│  │  Go (Chi + usecases + adapters)   │  │
│  └──────────────┬────────────────────┘  │
│  ┌──────────────▼──────┐  ┌───────────┐  │
│  │  config.db (SQLite) │  │  Remote   │  │
│  └─────────────────────┘  └───────────┘  │
└─────────────────────────────────────────┘
```

To add a database adapter: implement `entity.DatabaseAdapter` in `internal/adapter/database/`, register it in `factory.go` with an `AdapterInfo` UI config, and the connection dialog picks up the fields automatically.

### API overview

All endpoints are local under `/api`:

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/adapters` | List database adapters |
| `GET/POST` | `/api/connections` | List / create connections |
| `PUT/DELETE` | `/api/connections/{id}` | Update / delete |
| `POST` | `/api/connections/test` | Test credentials |
| `POST` | `/api/connections/{id}/test` | Test existing connection |
| `GET` | `/api/connections/{id}/tables` | List tables |
| `GET` | `/api/connections/{id}/tables/{name}` | Paginated table data |
| `GET` | `/api/connections/{id}/tables/{name}/schema` | Table schema |
| `POST` | `/api/connections/{id}/query` | Execute SQL |
| `GET/POST` | `/api/connections/{id}/queries` | Saved queries |
| `GET/POST` | `/api/connections/{id}/tabs` | Open tabs |
| `GET/POST` | `/api/theme` | Theme preference |
| `GET/POST` | `/api/layouts/{key}` | Panel layout |

### Tech stack

| Layer | Technology |
|-------|------------|
| Desktop shell | [webview_go](https://github.com/webview/webview_go) |
| Backend | Go, [Chi](https://github.com/go-chi/chi) |
| Frontend | React 19, TypeScript, [Bifrost](https://github.com/3-lines-studio/bifrost) |
| UI | shadcn/ui, Tailwind CSS 4 |
| State | Zustand, TanStack Query |
| SQL editor | react-simple-code-editor, Prism.js, sql-formatter |

### Releasing

Tagged pushes (`v*`) trigger the [Release workflow](.github/workflows/release.yml), building Linux and macOS ARM64 packages for GitHub Releases and updating the [Homebrew tap](https://github.com/3-lines-studio/homebrew-tap).

## License

See the repository for license information.
