package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mickaelroger/mcp-freshrss/client"
)

func NewListCategoriesTool(feverClient *client.FeverClient) mcp.Tool {
	return mcp.NewTool(
		"freshrss_list_categories",
		mcp.WithDescription("List all categories (groups) from FreshRSS. Returns a list of category objects with id and title."),
	)
}

func HandleListCategories(feverClient *client.FeverClient) func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		response, err := feverClient.GetGroups(ctx)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to list categories: %v", err)), nil
		}

		result := map[string]interface{}{
			"categories": response.Groups,
			"count":      len(response.Groups),
		}

		jsonData, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to format response: %v", err)), nil
		}

		return mcp.NewToolResultText(string(jsonData)), nil
	}
}
