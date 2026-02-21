# md2azw3

A lightweight HTTP service that converts Markdown files to AZW3 (Kindle) format. Pure Go implementation using [gomarkdown](https://github.com/gomarkdown/markdown) for Markdown parsing and [leotaku/mobi](https://github.com/leotaku/mobi) for KF8/AZW3 generation. No external binaries required.

## Quick Start

```bash
docker compose up --build
```

The server starts on port `8081`.

## API

### `POST /convert`

Converts a Markdown file (with optional cover image) to AZW3.

**Request:** `multipart/form-data`

| Field      | Type   | Required | Description           |
|------------|--------|----------|-----------------------|
| `markdown` | file   | Yes      | The `.md` file        |
| `cover`    | file   | No       | Cover image (jpg/png) |
| `title`    | string | No       | Book title            |
| `author`   | string | No       | Author name           |

**Response:** The converted `.azw3` file as a download.

**Example:**

```bash
curl -X POST \
  -F "markdown=@book.md" \
  -F "cover=@cover.jpg" \
  -F "title=My Book" \
  -F "author=John Doe" \
  http://localhost:8081/convert \
  -o book.azw3
```

### `GET /health`

Returns `{"status": "ok"}` when the service is running.

## Configuration

All configuration is done via environment variables:

| Variable                         | Default | Description             |
|----------------------------------|---------|-------------------------|
| `HTTP_PORT`                      | `8081`  | HTTP server port        |
| `IS_PRODUCTION_MODE`             | `false` | Production mode flag    |
| `LOGGER_LEVEL`                   | `debug` | Log level               |
| `LOGGER_IS_PRETTY_PRINT`         | `false` | JSON formatted logs     |
| `LOGGER_IS_REPORT_CALLER_MODE`   | `false` | Include caller info     |

## Development

```bash
# Hot reload
make dev

# Format code
make fmt

# Run tests
make test

# Lint
make lint
```

## Build

```bash
# Local Docker build
make build

# Multi-platform build and push
make push
```
