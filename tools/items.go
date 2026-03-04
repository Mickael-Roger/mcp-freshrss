package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mickaelroger/mcp-freshrss/client"
)

func NewMarkItemReadTool(feverClient *client.FeverClient) mcp.Tool {
	return mcp.NewTool(
		"freshrss_mark_item_read",
		mcp.WithDescription("Mark a news item as read by its ID. This operation cannot be undone unless you mark the item as unread again."),
		mcp.WithNumber("item_id",
			mcp.Required(),
			mcp.Description("The ID of the item to mark as read"),
		),
	)
}

func HandleMarkItemRead(feverClient *client.FeverClient) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		itemID, err := request.RequireInt("item_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		if err := feverClient.MarkItemRead(ctx, itemID); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to mark item %d as read: %v", itemID, err)), nil
		}

		result := map[string]interface{}{
			"success": true,
			"item_id": itemID,
			"message": fmt.Sprintf("Item %d has been marked as read", itemID),
		}

		jsonData, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to format response: %v", err)), nil
		}

		return mcp.NewToolResultText(string(jsonData)), nil
	}
}

func NewMarkItemUnreadTool(feverClient *client.FeverClient) mcp.Tool {
	return mcp.NewTool(
		"freshrss_mark_item_unread",
		mcp.WithDescription("Mark a news item as unread by its ID."),
		mcp.WithNumber("item_id",
			mcp.Required(),
			mcp.Description("The ID of the item to mark as unread"),
		),
	)
}

func HandleMarkItemUnread(feverClient *client.FeverClient) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		itemID, err := request.RequireInt("item_id")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		if err := feverClient.MarkItemUnread(ctx, itemID); err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to mark item %d as unread: %v", itemID, err)), nil
		}

		result := map[string]interface{}{
			"success": true,
			"item_id": itemID,
			"message": fmt.Sprintf("Item %d has been marked as unread", itemID),
		}

		jsonData, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to format response: %v", err)), nil
		}

		return mcp.NewToolResultText(string(jsonData)), nil
	}
}
