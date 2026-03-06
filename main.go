package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mark3labs/mcp-go/server"
	"github.com/mickaelroger/mcp-freshrss/client"
	"github.com/mickaelroger/mcp-freshrss/tools"
)

func main() {
	freshrssURL := os.Getenv("FRESHRSS_URL")
	if freshrssURL == "" {
		log.Fatal("FRESHRSS_URL environment variable is required")
	}

	apiKey := os.Getenv("FRESHRSS_API_KEY")
	if apiKey == "" {
		log.Fatal("FRESHRSS_API_KEY environment variable is required")
	}

	feverClient := client.NewFeverClient(freshrssURL, apiKey)

	s := server.NewMCPServer(
		"FreshRSS MCP Server",
		"1.0.0",
		server.WithToolCapabilities(false),
	)

	s.AddTool(tools.NewListCategoriesTool(feverClient), tools.HandleListCategories(feverClient))
	s.AddTool(tools.NewListFeedsTool(feverClient), tools.HandleListFeeds(feverClient))
	s.AddTool(tools.NewListItemsTool(feverClient), tools.HandleListItems(feverClient))
	s.AddTool(tools.NewGetItemTool(feverClient), tools.HandleGetItem(feverClient))
	s.AddTool(tools.NewMarkItemReadTool(feverClient), tools.HandleMarkItemRead(feverClient))
	s.AddTool(tools.NewMarkItemUnreadTool(feverClient), tools.HandleMarkItemUnread(feverClient))
	s.AddTool(tools.NewYoutubeTranscriptTool(), tools.HandleYoutubeTranscript())

	fmt.Fprintf(os.Stderr, "Starting FreshRSS MCP Server...\n")
	if err := server.ServeStdio(s); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
