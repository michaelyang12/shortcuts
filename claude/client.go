package claude

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/anthropics/anthropic-sdk-go"
)

func encodeBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

type Client struct {
	api anthropic.Client
}

func NewClient() *Client {
	return &Client{
		api: anthropic.NewClient(),
	}
}

// Text sends a text prompt to Claude and returns the response.
func (c *Client) Text(ctx context.Context, text, prompt string) (string, error) {
	msg, err := c.api.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.ModelClaudeSonnet4_5_20250929,
		MaxTokens: 4096,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(
				anthropic.NewTextBlock(fmt.Sprintf("Content:\n%s\n\nTask:\n%s", text, prompt)),
			),
		},
	})
	if err != nil {
		return "", fmt.Errorf("claude text request failed: %w", err)
	}
	return extractText(msg), nil
}

// Vision sends images + a prompt to Claude and returns the response.
func (c *Client) Vision(ctx context.Context, imagePaths []string, prompt string) (string, error) {
	var blocks []anthropic.ContentBlockParamUnion
	for _, path := range imagePaths {
		data, err := os.ReadFile(path)
		if err != nil {
			return "", fmt.Errorf("reading image %s: %w", path, err)
		}
		b64 := encodeBase64(data)
		blocks = append(blocks, anthropic.NewImageBlockBase64("image/jpeg", b64))
	}
	blocks = append(blocks, anthropic.NewTextBlock(prompt))

	msg, err := c.api.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.ModelClaudeSonnet4_5_20250929,
		MaxTokens: 4096,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(blocks...),
		},
	})
	if err != nil {
		return "", fmt.Errorf("claude vision request failed: %w", err)
	}
	return extractText(msg), nil
}

// VisionBase64 sends a single base64-encoded image + prompt to Claude.
func (c *Client) VisionBase64(ctx context.Context, b64Image, mediaType, prompt string) (string, error) {
	msg, err := c.api.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.ModelClaudeSonnet4_5_20250929,
		MaxTokens: 4096,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(
				anthropic.NewImageBlockBase64(mediaType, b64Image),
				anthropic.NewTextBlock(prompt),
			),
		},
	})
	if err != nil {
		return "", fmt.Errorf("claude vision request failed: %w", err)
	}
	return extractText(msg), nil
}

func extractText(msg *anthropic.Message) string {
	for _, block := range msg.Content {
		if block.Type == "text" {
			return block.Text
		}
	}
	return ""
}
