package gemini

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/genai"
)

// GenerateContent generates content using Google Gemini with a system prompt and user prompt.
func GenerateContent(ctx context.Context, apiKey, systemPrompt, userPrompt string) (string, error) {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return "", fmt.Errorf("gemini client: %w", err)
	}

	result, err := client.Models.GenerateContent(ctx, "gemini-2.5-flash", genai.Text(userPrompt), &genai.GenerateContentConfig{
		SystemInstruction: &genai.Content{
			Parts: []*genai.Part{
				{Text: systemPrompt},
			},
		},
	})
	if err != nil {
		return "", fmt.Errorf("gemini generate: %w", err)
	}

	return strings.TrimSpace(result.Text()), nil
}
