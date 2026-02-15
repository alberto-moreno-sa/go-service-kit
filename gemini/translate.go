package gemini

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/genai"
)

// Translate translates text to the specified target language using Google Gemini.
func Translate(ctx context.Context, apiKey, text, targetLang string) (string, error) {
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return "", fmt.Errorf("gemini client: %w", err)
	}

	prompt := fmt.Sprintf("Translate the following text to %s. Return only the translated text, nothing else.", targetLang)

	result, err := client.Models.GenerateContent(ctx, "gemini-2.5-flash", genai.Text(text), &genai.GenerateContentConfig{
		SystemInstruction: &genai.Content{
			Parts: []*genai.Part{
				{Text: prompt},
			},
		},
	})
	if err != nil {
		return "", fmt.Errorf("gemini generate: %w", err)
	}

	return strings.TrimSpace(result.Text()), nil
}
