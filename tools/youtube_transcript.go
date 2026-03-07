package tools

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
)

// IsYtDlpAvailable returns true if the yt-dlp command is found on PATH.
func IsYtDlpAvailable() bool {
	_, err := exec.LookPath("yt-dlp")
	return err == nil
}

func NewYoutubeTranscriptTool() mcp.Tool {
	return mcp.NewTool(
		"youtube_get_transcript",
		mcp.WithDescription("Get the transcript of a YouTube video as plain text using yt-dlp. Returns only the transcript content without timestamps or metadata, optimized for LLM analysis. Supports English ('en') and French ('fr') languages."),
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

		// Create a temporary directory that is cleaned up after use.
		tmpDir, err := os.MkdirTemp("", "mcp-freshrss-yt-*")
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to create temp directory: %v", err)), nil
		}
		defer os.RemoveAll(tmpDir)

		outputTemplate := filepath.Join(tmpDir, "transcript")

		// Run yt-dlp to download only the auto-generated subtitles, no video.
		cmd := exec.CommandContext(ctx,
			"yt-dlp",
			"--write-auto-sub",
			"--sub-lang", language,
			"--skip-download",
			"--no-playlist",
			"-o", outputTemplate,
			videoURL,
		)
		out, err := cmd.CombinedOutput()
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("yt-dlp failed: %v\n%s", err, string(out))), nil
		}

		// yt-dlp writes the file as <output>.<lang>.vtt
		vttPath := outputTemplate + "." + language + ".vtt"
		vttData, err := os.ReadFile(vttPath)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Subtitle file not found (the video may not have a transcript in '%s'): %v", language, err)), nil
		}

		text := parseVTT(string(vttData))
		if text == "" {
			return mcp.NewToolResultError("Transcript is empty"), nil
		}

		return mcp.NewToolResultText(text), nil
	}
}

// inlineTagRe matches VTT inline timing tags such as <00:00:01.234> and <c>.
var inlineTagRe = regexp.MustCompile(`<[^>]+>`)

// timestampLineRe matches VTT cue timing lines, e.g.:
// "00:00:00.400 --> 00:00:02.149 align:start position:0%"
var timestampLineRe = regexp.MustCompile(`^\d{2}:\d{2}:\d{2}\.\d{3} -->`)

// parseVTT converts a WebVTT file into plain text.
//
// YouTube's auto-generated VTT files have a peculiar structure: each cue
// block repeats the previous line as plain text on the first line, then adds
// new words with inline timing tags on the second line.  Collecting only the
// first plain-text line of each block (and de-duplicating consecutive equal
// lines) gives a clean, non-repetitive transcript.
func parseVTT(vtt string) string {
	scanner := bufio.NewScanner(strings.NewReader(vtt))

	var lines []string // candidate plain-text lines from cue blocks
	inCue := false     // true while we are inside a cue block
	firstInCue := true // true for the first content line of a cue

	for scanner.Scan() {
		line := scanner.Text()

		switch {
		case line == "WEBVTT" || strings.HasPrefix(line, "Kind:") || strings.HasPrefix(line, "Language:"):
			// Header lines – skip.
			inCue = false

		case timestampLineRe.MatchString(line):
			// Start of a new cue block.
			inCue = true
			firstInCue = true

		case line == "":
			// Blank line ends the current cue block.
			inCue = false

		case inCue && firstInCue:
			// First content line of a cue – this is the "clean" previous-sentence
			// line that YouTube inserts for context.  Strip any residual inline tags
			// and collect it.
			clean := inlineTagRe.ReplaceAllString(line, "")
			clean = strings.TrimSpace(clean)
			if clean != "" {
				lines = append(lines, clean)
			}
			firstInCue = false

		default:
			// Subsequent content lines within a cue – skip (they contain the
			// inline-tagged new words that will appear as the "clean" first line
			// of the very next cue block).
		}
	}

	// De-duplicate consecutive identical lines (YouTube repeats the same text
	// across consecutive cue blocks during short pauses).
	var deduped []string
	for _, l := range lines {
		if len(deduped) == 0 || deduped[len(deduped)-1] != l {
			deduped = append(deduped, l)
		}
	}

	return strings.Join(deduped, " ")
}
