package gemini

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	defaultTimeout       = 30 * time.Second
	maxErrorBodyBytes    = 8 * 1024
	maxResponseBodyBytes = 2 * 1024 * 1024
)

type Client struct {
	apiKey string
	model  string
	client *http.Client
}

func NewClient(apiKey, model string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: defaultTimeout}
	}

	return &Client{
		apiKey: apiKey,
		model:  model,
		client: httpClient,
	}
}

type generateContentRequest struct {
	Contents []content `json:"contents"`
}

type content struct {
	Parts []part `json:"parts"`
}

type part struct {
	Text string `json:"text"`
}

type generateContentResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`

	Error *struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Status  string `json:"status"`
	} `json:"error"`
}

type searchProductsResult struct {
	IDs []int64 `json:"ids"`
}

func (c *Client) SearchProducts(ctx context.Context, userQuery string, projectsJSON string) ([]int64, error) {
	if strings.TrimSpace(c.apiKey) == "" {
		return []int64{}, nil
	}
	if strings.TrimSpace(c.model) == "" {
		return nil, fmt.Errorf("gemini model is empty")
	}

	prompt := buildPrompt(userQuery, projectsJSON)

	payload := generateContentRequest{
		Contents: []content{
			{
				Parts: []part{
					{Text: prompt},
				},
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal gemini request: %w", err)
	}

	endpoint := fmt.Sprintf(
		"https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s",
		c.model,
		c.apiKey,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create gemini request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("execute gemini request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		snippet, _ := io.ReadAll(io.LimitReader(resp.Body, maxErrorBodyBytes))
		return nil, fmt.Errorf("gemini api returned status %d: %s", resp.StatusCode, strings.TrimSpace(string(snippet)))
	}

	limitedBody := io.LimitReader(resp.Body, maxResponseBodyBytes)

	var result generateContentResponse
	if err := json.NewDecoder(limitedBody).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode gemini response: %w", err)
	}

	if result.Error != nil {
		return nil, fmt.Errorf(
			"gemini api error: code=%d status=%s message=%s",
			result.Error.Code,
			result.Error.Status,
			result.Error.Message,
		)
	}

	text, ok := firstCandidateText(result)
	if !ok {
		return []int64{}, nil
	}

	cleanJSON, err := extractJSON(text)
	if err != nil {
		return nil, fmt.Errorf("extract json from gemini response: %w", err)
	}

	var finalResult searchProductsResult
	if err := json.Unmarshal([]byte(cleanJSON), &finalResult); err != nil {
		return nil, fmt.Errorf("unmarshal gemini result json: %w; raw=%q", err, cleanJSON)
	}

	return finalResult.IDs, nil
}

func buildPrompt(userQuery string, projectsJSON string) string {
	return fmt.Sprintf(`
Ты — экспертный ИИ-консультант РОСТ Мебель.
Твоя задача: проанализировать запрос пользователя и подобрать из списка ниже наиболее подходящие реализованные проекты.

КРИТИЧЕСКИЕ ПРАВИЛА:
1. КАТЕГОРИЯ: Если в запросе указан тип мебели (кухня, шкаф, гардеробная), возвращай ТОЛЬКО проекты из этой категории. Если пользователь ищет КУХНЮ, не предлагай шкафы, даже если они подходят по бюджету.
2. БЮДЖЕТ: Если указана сумма (например, "до 100000"), исключай все проекты, бюджет которых превышает эту сумму.
3. ТЕГИ: Анализируй стиль (лофт, классика, сканди) и материалы.
4. СОРТИРОВКА: Первыми ставь проекты, максимально похожие на описание.

СПИСОК ПРОЕКТОВ (JSON):
%s

ЗАПРОС ПОЛЬЗОВАТЕЛЯ: %q

ОТВЕТЬ СТРОГО В ФОРМАТЕ JSON: {"ids": [id1, id2, ...]}
Если ничего не подходит, верни {"ids": []}.
`, projectsJSON, userQuery)
}

func firstCandidateText(resp generateContentResponse) (string, bool) {
	if len(resp.Candidates) == 0 {
		return "", false
	}
	if len(resp.Candidates[0].Content.Parts) == 0 {
		return "", false
	}

	text := strings.TrimSpace(resp.Candidates[0].Content.Parts[0].Text)
	if text == "" {
		return "", false
	}

	return text, true
}

func extractJSON(text string) (string, error) {
	text = strings.TrimSpace(text)
	if text == "" {
		return "", fmt.Errorf("empty response text")
	}

	if strings.HasPrefix(text, "```") {
		lines := strings.Split(text, "\n")
		if len(lines) >= 2 {
			lines = lines[1:]
		}
		if n := len(lines); n > 0 && strings.TrimSpace(lines[n-1]) == "```" {
			lines = lines[:n-1]
		}
		text = strings.TrimSpace(strings.Join(lines, "\n"))
	}

	start := strings.Index(text, "{")
	end := strings.LastIndex(text, "}")
	if start == -1 || end == -1 || start > end {
		return "", fmt.Errorf("json object not found")
	}

	return strings.TrimSpace(text[start : end+1]), nil
}
