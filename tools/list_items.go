package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mickaelroger/mcp-freshrss/client"
	"github.com/mickaelroger/mcp-freshrss/models"
)

func NewListItemsTool(feverClient *client.FeverClient) mcp.Tool {
	return mcp.NewTool(
		"freshrss_list_items",
		mcp.WithDescription("List news items with optional read status filtering. Returns item ID, title, and feed_id for each item. Note: The Fever API returns max 50 items per request."),
		mcp.WithString("read_status",
			mcp.Description("Filter by read status: 'all' (default), 'read', or 'unread'"),
			mcp.Enum("all", "read", "unread"),
		),
	)
}

func HandleListItems(feverClient *client.FeverClient) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		readStatus := request.GetString("read_status", "all")

		var items []models.Item
		var err error

		if readStatus == "unread" {
			unreadIDs, err := feverClient.GetUnreadItemStringIDs(ctx)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to get unread items: %v", err)), nil
			}

			items, err = feverClient.GetItemsByIDs(ctx, unreadIDs)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to get items: %v", err)), nil
			}
		} else {
			items, err = feverClient.GetItems(ctx)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to list items: %v", err)), nil
			}
		}

		var unreadIDs map[string]bool
		if readStatus == "read" {
			unreadIDStrings, err := feverClient.GetUnreadItemStringIDs(ctx)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to get unread items: %v", err)), nil
			}
			unreadIDs = make(map[string]bool)
			for _, id := range unreadIDStrings {
				unreadIDs[id] = true
			}
		}

		var filteredItems []map[string]interface{}
		for _, item := range items {
			isUnread := unreadIDs == nil || unreadIDs[item.ID]

			if readStatus == "read" && isUnread {
				continue
			}

			itemMap := map[string]interface{}{
				"id":      item.ID,
				"title":   item.Title,
				"feed_id": item.FeedID,
			}

			if readStatus != "all" {
				itemMap["is_read"] = !isUnread
			}

			filteredItems = append(filteredItems, itemMap)
		}

		result := map[string]interface{}{
			"items": filteredItems,
			"count": len(filteredItems),
		}

		jsonData, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to format response: %v", err)), nil
		}

		return mcp.NewToolResultText(string(jsonData)), nil
	}
}
