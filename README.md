# FreshRSS MCP Server

A Model Context Protocol (MCP) server for FreshRSS written in Go. This server enables LLMs to interact with FreshRSS through the Fever API, providing tools to manage RSS feeds and articles.

## Features

- **List Categories**: Retrieve all categories/groups from your FreshRSS instance
- **List Feeds**: List feeds with optional filtering by category and read status
- **Mark as Read**: Mark individual news items as read
- **Mark as Unread**: Mark individual news items as unread

## Requirements

- Go 1.21 or higher
- FreshRSS server with Fever API enabled
- API credentials (username and API password)

## Installation

### Build from Source

```bash
git clone https://github.com/yourusername/mcp-freshrss.git
cd mcp-freshrss
go build -o freshrss-mcp
```

### Install with Go

```bash
go install github.com/yourusername/mcp-freshrss@latest
```

## Configuration

### Environment Variables

Set the following environment variables before running the server:

| Variable | Description | Example |
|----------|-------------|---------|
| `FRESHRSS_URL` | Base URL of your FreshRSS server | `https://freshrss.example.net` |
| `FRESHRSS_API_KEY` | Fever API key (pre-computed MD5 hash) | `abc123def456...` |

### Generating the API Key

The `FRESHRSS_API_KEY` environment variable must contain the **pre-computed** MD5 hash of your FreshRSS username and API password concatenated with a colon.

**Compute the hash:**

```bash
# Linux
echo -n "username:api_password" | md5sum | cut -d' ' -f1

# macOS
echo -n "username:api_password" | md5

# Using Python
python3 -c "import hashlib; print(hashlib.md5('username:api_password'.encode()).hexdigest())"
```

**Important**:
- Use your FreshRSS **API password**, not your login password
- You can set/generate the API password in FreshRSS under **Settings** → **User profile** → **API password**
- The server expects the already-hashed value, not the raw `username:password` string

### Enabling Fever API in FreshRSS

1. Log into your FreshRSS instance
2. Go to **Settings** → **User profile**
3. Check **Enable API access**
4. Set or generate an **API password**
5. Save changes

## Usage

### Running the Server

```bash
export FRESHRSS_URL="https://your-freshrss.example.net"
export FRESHRSS_API_KEY="your_md5_hash_here"

./freshrss-mcp
```

The server communicates via stdio, making it compatible with any MCP client.

### Integration with MCP Clients

#### Claude Desktop

Add to your Claude Desktop configuration (`~/Library/Application Support/Claude/claude_desktop_config.json` on macOS):

```json
{
  "mcpServers": {
    "freshrss": {
      "command": "/path/to/freshrss-mcp",
      "env": {
        "FRESHRSS_URL": "https://your-freshrss.example.net",
        "FRESHRSS_API_KEY": "your_md5_hash_here"
      }
    }
  }
}
```

#### Other MCP Clients

The server uses stdio transport and is compatible with any MCP-compliant client. Configure your client to run the `freshrss-mcp` binary with the required environment variables.

## Available Tools

### freshrss_list_categories

List all categories (groups) from your FreshRSS instance.

**Parameters**: None

**Example Response**:
```json
{
  "categories": [
    {"id": 1, "title": "Technology"},
    {"id": 2, "title": "News"},
    {"id": 3, "title": "Entertainment"}
  ]
}
```

### freshrss_list_feeds

List feeds with optional filtering.

**Parameters**:
- `category_id` (optional, integer): Filter by category ID
- `read_status` (optional, string): Filter by read status
  - `all` (default): Show all feeds
  - `read`: Show feeds with read items
  - `unread`: Show feeds with unread items

**Example**:
```
List all feeds with unread items in category 1
```

**Example Response**:
```json
{
  "feeds": [
    {
      "id": 1,
      "title": "Example Blog",
      "url": "https://example.com/feed",
      "site_url": "https://example.com",
      "unread_count": 5
    }
  ]
}
```

### freshrss_get_item

Get a news item by its ID.

**Parameters**:
- `item_id` (required, integer): The ID of the item to retrieve

**Example**:
```
Get item 12345
```

**Example Response**:
```json
{
  "id": "12345",
  "feed_id": 1,
  "title": "Article Title",
  "author": "Author Name",
  "url": "https://example.com/article",
  "html": "<p>Article content...</p>",
  "is_read": false,
  "is_saved": false,
  "created_on_time": 1234567890
}
```

### freshrss_mark_item_read

Mark a news item as read.

**Parameters**:
- `item_id` (required, integer): The ID of the item to mark as read

**Example**:
```
Mark item 12345 as read
```

### freshrss_mark_item_unread

Mark a news item as unread.

**Parameters**:
- `item_id` (required, integer): The ID of the item to mark as unread

**Example**:
```
Mark item 12345 as unread
```

## API Reference

This server uses the FreshRSS Fever API implementation. Key endpoints:

| Operation | Endpoint | Method |
|-----------|----------|--------|
| List categories | `?api&groups` | POST |
| List feeds | `?api&feeds` | POST |
| List unread IDs | `?api&unread_item_ids` | POST |
| Mark item read | POST body: `mark=item&as=read&id=<id>` | POST |
| Mark item unread | POST body: `mark=item&as=unread&id=<id>` | POST |

For full API documentation, see:
- [FreshRSS Fever API Documentation](https://freshrss.github.io/FreshRSS/en/developers/06_Fever_API.html)
- [Original Fever API Documentation](https://web.archive.org/web/20230616124016/https://feedafever.com/api)

## Development

### Project Structure

```
mcp-freshrss/
├── main.go              # Entry point
├── client/
│   └── fever.go         # Fever API client
├── tools/
│   ├── categories.go    # Categories tool
│   ├── feeds.go         # Feeds tool
│   └── items.go         # Item tools
└── models/
    └── types.go         # Data models
```

### Building

```bash
go build -o freshrss-mcp
```

### Testing

```bash
go test ./...
```

## Limitations

- The Fever API does not support adding, editing, or deleting feeds
- Hot Links feature is not supported by FreshRSS
- Maximum 50 items returned per request

## License

MIT License

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Support

For issues related to:
- **This MCP server**: Open an issue on GitHub
- **FreshRSS**: See [FreshRSS documentation](https://freshrss.github.io/)
- **Fever API**: See the [API documentation](https://freshrss.github.io/FreshRSS/en/developers/06_Fever_API.html)
