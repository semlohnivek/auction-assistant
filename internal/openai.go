// internal/openai.go
package internal

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"bidzauction/config"

	openai "github.com/sashabaranov/go-openai"
)

type LotDetails struct {
	Title           string `json:"title"`
	Description     string `json:"description"`
	Make            string `json:"make"`
	Model           string `json:"model"`
	Condition       string `json:"condition"`
	Year            string `json:"year"`
	CountryOfOrigin string `json:"country_of_origin"`
}

func AnalyzeImageURLs(imageURLs []string) (LotDetails, error) {
	cfg := config.Current.OpenAI

	client := openai.NewClient(cfg.Key)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var parts []openai.ChatMessagePart
	for _, url := range imageURLs {
		parts = append(parts, openai.ChatMessagePart{
			Type:     openai.ChatMessagePartTypeImageURL,
			ImageURL: &openai.ChatMessageImageURL{URL: url},
		})
	}

	req := openai.ChatCompletionRequest{
		Model:       cfg.Model,
		MaxTokens:   cfg.MaxTokens,
		Temperature: float32(cfg.Temperature),
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: cfg.SystemPrompt,
			},
			{
				Role:         openai.ChatMessageRoleUser,
				MultiContent: parts,
			},
		},
	}

	resp, err := client.CreateChatCompletion(ctx, req)
	if err != nil {
		return LotDetails{}, fmt.Errorf("OpenAI API error: %w", err)
	}
	if len(resp.Choices) == 0 {
		return LotDetails{}, errors.New("no response choices from OpenAI")
	}

	text := strings.TrimSpace(resp.Choices[0].Message.Content)
	var details LotDetails
	if err := json.Unmarshal([]byte(text), &details); err != nil {
		return LotDetails{}, fmt.Errorf("error parsing OpenAI JSON response: %w\nRaw: %s", err, text)
	}

	return details, nil
}
