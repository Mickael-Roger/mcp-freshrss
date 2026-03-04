package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mickaelroger/mcp-freshrss/client"
	"github.com/mickaelroger/mcp-freshrss/models"
)

func NewListFeedsTool(feverClient *client.FeverClient) mcp.Tool {
	return mcp.NewTool(
		"freshrss_list_feeds",
		mcp.WithDescription("List feeds for one or all categories with read/unread filtering. Returns feed objects with id, title, url, and unread count."),
		mcp.WithNumber("category_id",
			mcp.Description("Filter by category ID. If not provided, lists all feeds."),
		),
		mcp.WithString("read_status",
			mcp.Description("Filter by read status: 'all' (default), 'read', or 'unread'"),
			mcp.Enum("all", "read", "unread"),
		),
	)
}

func HandleListFeeds(feverClient *client.FeverClient) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		categoryID := request.GetInt("category_id", 0)
		readStatus := request.GetString("read_status", "all")

		feedsResponse, err := feverClient.GetFeeds(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to list feeds: %v", err)), nil
		}

		groupsResponse, err := feverClient.GetGroups(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to get groups: %v", err)), nil
		}

		unreadIDs := make(map[int]bool)
		if readStatus == "unread" || readStatus == "read" {
			ids, err := feverClient.GetUnreadItemIDs(ctx)
			if err != nil {
				return mcp.NewToolResultError(fmt.Sprintf("Failed to get unread items: %v", err)), nil
			}
			for _, id := range ids {
				unreadIDs[id] = true
			}
		}

		var filteredFeeds []map[string]interface{}
		for _, feed := range feedsResponse.Feeds {
			if categoryID > 0 {
				found := false
				for _, fg := range feedsResponse.FeedsGroups {
					if fg.GroupID == categoryID {
						if containsFeedID(fg.FeedIDs, feed.ID) {
							found = true
							break
						}
					}
				}
				if !found {
					continue
				}
			}

			feedMap := map[string]interface{}{
				"id":       feed.ID,
				"title":    feed.Title,
				"url":      feed.URL,
				"site_url": feed.SiteURL,
				"is_spark": feed.IsSpark,
			}

			if readStatus != "all" {
				feedMap["has_unread"] = hasUnreadItems(feed.ID, groupsResponse.FeedsGroups, unreadIDs)
				if readStatus == "unread" && !feedMap["has_unread"].(bool) {
					continue
				}
				if readStatus == "read" && feedMap["has_unread"].(bool) {
					continue
				}
			}

			filteredFeeds = append(filteredFeeds, feedMap)
		}

		result := map[string]interface{}{
			"feeds": filteredFeeds,
			"count": len(filteredFeeds),
		}

		if categoryID > 0 {
			result["category_id"] = categoryID
		}

		jsonData, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to format response: %v", err)), nil
		}

		return mcp.NewToolResultText(string(jsonData)), nil
	}
}

func parseFeedIDs(idsStr string) []int {
	if idsStr == "" {
		return []int{}
	}

	parts := strings.Split(idsStr, ",")
	ids := make([]int, 0, len(parts))

	for _, part := range parts {
		var id int
		if _, err := fmt.Sscanf(strings.TrimSpace(part), "%d", &id); err == nil {
			ids = append(ids, id)
		}
	}

	return ids
}

func containsFeedID(idsStr string, feedID int) bool {
	ids := parseFeedIDs(idsStr)
	for _, id := range ids {
		if id == feedID {
			return true
		}
	}
	return false
}

func hasUnreadItems(feedID int, feedsGroups []models.FeedsGroup, unreadIDs map[int]bool) bool {
	return len(unreadIDs) > 0
}
