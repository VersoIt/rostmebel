package gemini

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	apiKey string
	model  string
	client *http.Client
}

func NewClient(apiKey, model string) *Client {
	return &Client{
		apiKey: apiKey,
		model:  model,
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

type GeminiRequest struct {
	Contents []Content `json:"contents"`
}

type Content struct {
	Parts []Part `json:"parts"`
}

type Part struct {
	Text string `json:"text"`
}

type GeminiResponse struct {
	Candidates []Candidate `json:"candidates"`
}

type Candidate struct {
	Content Content `json:"content"`
}

type AIResponse struct {
	IDs []int64 `json:"ids"`
}

func (c *Client) SearchProducts(ctx context.Context, userQuery string, productsJSON string) ([]int64, error) {
	if c.apiKey == "" {
		return nil, fmt.Errorf("api key is not set")
	}

	prompt := fmt.Sprintf(`
Ты — экспертный ИИ-дизайнер магазина РОСТ Мебель на базе Gemma 4. 
Твоя задача: проанализировать свободный запрос пользователя и подобрать из списка ниже наиболее подходящие товары.

ИНСТРУКЦИИ:
1. Анализируй стиль (скандинавский, лофт, классика), материал, цвет и бюджет, если они указаны.
2. Если в запросе есть цена (например, "до 50000"), исключай товары дороже.
3. Если запрос касается комнаты (например, "для спальни"), подбирай товары с соответствующими тегами.
4. Сортируй результат по степени соответствия запросу.

СПИСОК ТОВАРОВ (JSON):
%s

ЗАПРОС ПОЛЬЗОВАТЕЛЯ: "%s"

ОТВЕТЬ СТРОГО В ФОРМАТЕ JSON: {"ids": [id1, id2, ...]}
Если ничего не подходит, верни {"ids": []}.
`, productsJSON, userQuery)

	reqBody := GeminiRequest{
		Contents: []Content{
			{
				Parts: []Part{{Text: prompt}},
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", c.model, c.apiKey)
	
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("gemini api returned status %d", resp.StatusCode)
	}

	var geminiResp GeminiResponse
	if err := json.NewDecoder(resp.Body).Decode(&geminiResp); err != nil {
		return nil, err
	}

	if len(geminiResp.Candidates) == 0 || len(geminiResp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("gemini api returned no candidates")
	}

	responseText := geminiResp.Candidates[0].Content.Parts[0].Text
	
	// Try to parse the response text as JSON
	var aiResp AIResponse
	// Sometimes Gemini wraps JSON in markdown blocks
	cleanedText := cleanJSONResponse(responseText)
	if err := json.Unmarshal([]byte(cleanedText), &aiResp); err != nil {
		return nil, fmt.Errorf("failed to parse gemini response: %w, text: %s", err, responseText)
	}

	return aiResp.IDs, nil
}

func cleanJSONResponse(text string) string {
	// Simple cleanup for markdown code blocks
	if start := bytes.Index([]byte(text), []byte("```json")); start != -1 {
		text = text[start+7:]
	} else if start := bytes.Index([]byte(text), []byte("```")); start != -1 {
		text = text[start+3:]
	}
	if end := bytes.LastIndex([]byte(text), []byte("```")); end != -1 {
		text = text[:end]
	}
	return text
}
