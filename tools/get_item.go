package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mickaelroger/mcp-freshrss/client"
)

func NewGetItemTool(feverClient *client.FeverClient) mcp.Tool {
	return mcp.NewTool(
		"freshrss_get_item",
		mcp.WithDescription("Get a news item by its ID. Returns the full item details including title, author, URL, content, and read status."),
		mcp.WithNumber("item_id",
			mcp.Required(),
			mcp.Description("The ID of the item to retrieve"),
		),
	)
}

func HandleGetItem(feverClient *client.FeverClient) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		itemID, err := request.RequireInt("item_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		item, err := feverClient.GetItem(ctx, itemID)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get item %d: %v", itemID, err)), nil
		}

		result := map[string]interface{}{
			"id":              item.ID,
			"feed_id":         item.FeedID,
			"title":           item.Title,
			"author":          item.Author,
			"url":             item.URL,
			"html":            item.HTML,
			"is_read":         item.IsRead == 1,
			"is_saved":        item.IsSaved == 1,
			"created_on_time": item.CreatedOnTime,
		}

		jsonData, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to format response: %v", err)), nil
		}

		return mcp.NewToolResultText(string(jsonData)), nil
	}
}
