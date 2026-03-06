package tools

import (
	"context"
	"strings"

	"github.com/kkdai/youtube/v2"
	"github.com/mark3labs/mcp-go/mcp"
)

func NewYoutubeTranscriptTool() mcp.Tool {
	return mcp.NewTool(
		"youtube_get_transcript",
		mcp.WithDescription("Get the transcript of a YouTube video as plain text. Returns only the transcript content without timestamps or metadata, optimized for LLM analysis. Supports only English ('en') and French ('fr') languages."),
		mcp.WithString("url",
			mcp.Required(),
			mcp.Description("The YouTube video URL or video ID (e.g., https://www.youtube.com/watch?v=VIDEO_ID or just VIDEO_ID)"),
		),
		mcp.WithString("language",
			mcp.Description("Language code for the transcript: 'en' for English or 'fr' for French. Defaults to 'en' if not specified."),
			mcp.Enum("en", "fr"),
		),
	)
}

func HandleYoutubeTranscript() func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	return func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		videoURL, err := request.RequireString("url")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		language := request.GetString("language", "en")

		if language != "en" && language != "fr" {
			return mcp.NewToolResultError("Invalid language: only 'en' (English) and 'fr' (French) are supported"), nil
		}

		client := youtube.Client{}

		video, err := client.GetVideo(videoURL)
		if err != nil {
			return mcp.NewToolResultError("Failed to fetch video: " + err.Error()), nil
		}

		transcript, err := client.GetTranscript(video, language)
		if err != nil {
			if err == youtube.ErrTranscriptDisabled {
				return mcp.NewToolResultError("No transcript available for this video"), nil
			}
			return mcp.NewToolResultError("Failed to fetch transcript: " + err.Error()), nil
		}

		var plainText strings.Builder
		for _, segment := range transcript {
			if segment.Text != "" {
				plainText.WriteString(segment.Text)
				plainText.WriteString(" ")
			}
		}

		transcriptText := strings.TrimSpace(plainText.String())

		if transcriptText == "" {
			return mcp.NewToolResultError("Transcript is empty"), nil
		}

		return mcp.NewToolResultText(transcriptText), nil
	}
}
