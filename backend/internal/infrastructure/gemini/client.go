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

func (c *Client) SearchProducts(ctx context.Context, userQuery string, projectsJSON string) ([]int64, error) {
	if c.apiKey == "" {
		return []int64{}, nil
	}

	prompt := fmt.Sprintf(`
Ты — экспертный ИИ-консультант РОСТ Мебель. 
Твоя задача: проанализировать запрос пользователя и подобрать из списка ниже наиболее подходящие реализованные проекты.

КРИТИЧЕСКИЕ ПРАВИЛА:
1. КАТЕГОРИЯ: Если в запросе указан тип мебели (кухня, шкаф, гардеробная), возвращай ТОЛЬКО проекты из этой категории. Если пользователь ищет КУХНЮ, не предлагай шкафы, даже если они подходят по бюджету.
2. БЮДЖЕТ: Если указана сумма (например, "до 100000"), исключай все проекты, бюджет которых превышает эту сумму.
3. ТЕГИ: Анализируй стиль (лофт, классика, сканди) и материалы.
4. СОРТИРОВКА: Первыми ставь проекты, максимально похожие на описание.

СПИСОК ПРОЕКТОВ (JSON):
%s

ЗАПРОС ПОЛЬЗОВАТЕЛЯ: "%s"

ОТВЕТЬ СТРОГО В ФОРМАТЕ JSON: {"ids": [id1, id2, ...]}
Если ничего не подходит, верни {"ids": []}.
`, projectsJSON, userQuery)

	payload := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]interface{}{
					{"text": prompt},
				},
			},
		},
	}

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", c.model, c.apiKey)
	
	data, _ := json.Marshal(payload)
	req, _ := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if len(result.Candidates) == 0 || len(result.Candidates[0].Content.Parts) == 0 {
		return []int64{}, nil
	}

	cleanJSON := extractJSON(result.Candidates[0].Content.Parts[0].Text)
	
	var finalResult struct {
		IDs []int64 `json:"ids"`
	}
	if err := json.Unmarshal([]byte(cleanJSON), &finalResult); err != nil {
		return nil, nil
	}

	return finalResult.IDs, nil
}

func extractJSON(text string) string {
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
