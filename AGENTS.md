# FreshRSS MCP Server - Agent Guidelines

## Overview

This is an MCP (Model Context Protocol) server for FreshRSS written in Go. It provides tools to interact with FreshRSS through the Fever API, enabling LLMs to manage RSS feeds and articles.

## Project Details

- **Language**: Go (Golang)
- **Transport**: stdio only (no HTTP server)
- **API**: FreshRSS Fever API (`/api/fever.php`)
- **Authentication**: API key (MD5 hash of `username:api_password`)
- **YouTube Integration**: Uses the `yt-dlp` system command for video transcripts (tool is only exposed when `yt-dlp` is available on the host)

## Environment Variables

Required configuration:
- `FRESHRSS_URL`: Base URL of the FreshRSS server (e.g., `https://freshrss.example.net`)
- `FRESHRSS_API_KEY`: Fever API key (pre-computed MD5 hash of `username:api_password`)

## Architecture

### Project Structure
```
mcp-freshrss/
├── main.go              # Entry point, MCP server initialization
├── go.mod               # Go module definition
├── go.sum               # Go dependencies checksum
├── client/
│   └── fever.go         # Fever API client implementation
├── tools/
│   ├── categories.go    # List categories tool
│   ├── feeds.go         # List feeds tool
│   ├── list_items.go    # List items tool
│   ├── get_item.go      # Get item tool
│   ├── items.go         # Mark read/unread tools
│   └── youtube_transcript.go  # YouTube transcript tool
├── models/
│   └── types.go         # Data models for Fever API responses
└── AGENTS.md            # This file
```

## Implemented Tools

### 1. freshrss_list_categories
List all categories (groups) from FreshRSS.
- **Read-only**: Yes
- **API endpoint**: `?api&groups`

### 2. freshrss_list_feeds
List feeds for one or all categories with read/unread filtering.
- **Parameters**:
  - `category_id` (optional): Filter by category ID. If not provided, lists all feeds.
  - `read_status` (optional): Filter by read status (`all`, `read`, `unread`). Default: `all`
- **Read-only**: Yes
- **API endpoints**: `?api&feeds`, `?api&groups`, `?api&unread_item_ids`

### 3. freshrss_list_items
List news items with optional read status filtering.
- **Parameters**:
  - `read_status` (optional): Filter by read status (`all`, `read`, `unread`). Default: `all`
- **Returns**: Item ID, title, and feed_id for each item
- **Read-only**: Yes
- **API endpoints**: `?api&items`, `?api&unread_item_ids`
- **Note**: The Fever API returns a maximum of 50 items per request

### 4. freshrss_get_item
Get a news item by its ID.
- **Parameters**:
  - `item_id` (required): The ID of the item to retrieve
- **Read-only**: Yes
- **API endpoint**: `?api&items&with_ids=<item_id>`

### 5. freshrss_mark_item_read
Mark a news item as read.
- **Parameters**:
  - `item_id` (required): The ID of the item to mark as read
- **Read-only**: No
- **API endpoint**: POST with `mark=item&as=read&id=<item_id>`

### 6. freshrss_mark_item_unread
Mark a news item as unread.
- **Parameters**:
  - `item_id` (required): The ID of the item to mark as unread
- **Read-only**: No
- **API endpoint**: POST with `mark=item&as=unread&id=<item_id>`

### 7. youtube_get_transcript
Get the transcript of a YouTube video as plain text, optimized for LLM analysis.
- **Parameters**:
  - `url` (required): The YouTube video URL or video ID
  - `language` (optional): Language code for the transcript. Supported values: `en` (English) or `fr` (French). Default: `en`
- **Returns**: Plain text transcript without timestamps or metadata
- **Read-only**: Yes
- **Availability**: Only registered if `yt-dlp` is found on `PATH` at startup
- **Implementation**: Runs `yt-dlp --write-auto-sub --skip-download` into a temp directory, then parses the resulting WebVTT file into plain text (strips timestamps, inline timing tags, and deduplicates repeated lines produced by YouTube's auto-caption format)
- **Note**: Not all videos have transcripts. Only newer videos typically have auto-generated captions.

## Fever API Reference

### Authentication
All requests must include an `api_key` parameter in the POST body. The API key is the MD5 hash of `username:api_password` (pre-computed and provided via the `FRESHRSS_API_KEY` environment variable).

The server receives the already-hashed API key and sends it directly to the Fever API.

Example of computing the API key (for user reference):
```bash
api_key=$(echo -n "username:api_password" | md5sum | cut -d' ' -f1)
```

### Request Format
- **Base URL**: `{FRESHRSS_URL}/api/fever.php?api`
- **Method**: POST
- **Body**: `form-data` with `api_key` parameter

### Response Format
All responses are JSON objects containing at minimum:
- `api_version`: API version number (integer)
- `auth`: Authentication status (1 = success, 0 = failure)

### Key Endpoints

| Endpoint | Description |
|----------|-------------|
| `?api&groups` | List all categories/groups |
| `?api&feeds` | List all feeds |
| `?api&items` | List items (with pagination) |
| `?api&unread_item_ids` | Get list of unread item IDs |
| `?api&saved_item_ids` | Get list of saved/starred item IDs |
| POST `mark=item&as=read&id=<id>` | Mark item as read |
| POST `mark=item&as=unread&id=<id>` | Mark item as unread |

## MCP Protocol

This server implements the MCP protocol using stdio transport:
- Communication via standard input/output
- JSON-RPC 2.0 message format
- No HTTP server required

## Coding Guidelines

### Go Best Practices
- Use Go modules for dependency management
- Follow standard Go project layout
- Implement proper error handling
- Use context for cancellation and timeouts
- Write unit tests for critical functions

### Tool Implementation
Each tool must:
1. Have a clear, descriptive name prefixed with `freshrss_`
2. Validate all input parameters
3. Return structured JSON responses
4. Provide meaningful error messages
5. Include proper MCP tool annotations

### Error Handling
- Wrap errors with context using `fmt.Errorf` or `errors.Wrap`
- Return user-friendly error messages to MCP clients
- Log errors to stderr (never stdout in stdio mode)

## Testing

Run tests with:
```bash
go test ./...
```

Build the server:
```bash
go build -o freshrss-mcp
```

## Debugging

To test the Fever API manually:
```bash
# Test authentication
curl -s -F "api_key=YOUR_API_KEY" 'https://your-freshrss.example.net/api/fever.php?api'

# List categories
curl -s -F "api_key=YOUR_API_KEY" 'https://your-freshrss.example.net/api/fever.php?api&groups'

# List feeds
curl -s -F "api_key=YOUR_API_KEY" 'https://your-freshrss.example.net/api/fever.php?api&feeds'
```

## Known Limitations

- The Fever API does not support adding/editing/deleting feeds (only reading)
- Hot Links feature is not supported by FreshRSS
- Maximum 50 items returned per request when fetching items
- YouTube transcripts are not available for all videos (only newer videos typically have transcripts)
- YouTube transcript tool only supports English ('en') and French ('fr') languages

## Resources

- [FreshRSS Fever API Documentation](https://freshrss.github.io/FreshRSS/en/developers/06_Fever_API.html)
- [Original Fever API Documentation](https://web.archive.org/web/20230616124016/https://feedafever.com/api)
- [MCP Specification](https://modelcontextprotocol.io/)
