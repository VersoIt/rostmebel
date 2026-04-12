package gemini

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	defaultTimeout       = 30 * time.Second
	maxErrorBodyBytes    = 8 * 1024
	maxResponseBodyBytes = 2 * 1024 * 1024

	DefaultModel = "gemini-2.5-flash"
)

var ErrDisabled = errors.New("gemini is disabled")

type Client struct {
	apiKey         string
	model          string
	fallbackModels []string
	baseURL        string
	client         *http.Client
}

type ClientOptions struct {
	APIKey         string
	Model          string
	FallbackModels []string
	BaseURL        string
	HTTPClient     *http.Client
}

func NewClient(apiKey, model string, httpClient *http.Client) *Client {
	return NewClientWithOptions(ClientOptions{
		APIKey:     apiKey,
		Model:      model,
		HTTPClient: httpClient,
	})
}

func NewClientWithOptions(opts ClientOptions) *Client {
	httpClient := opts.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{Timeout: defaultTimeout}
	}
	baseURL := strings.TrimRight(strings.TrimSpace(opts.BaseURL), "/")
	if baseURL == "" {
		baseURL = "https://generativelanguage.googleapis.com/v1beta"
	}

	return &Client{
		apiKey:         strings.TrimSpace(opts.APIKey),
		model:          normalizeModel(opts.Model),
		fallbackModels: normalizeModels(opts.FallbackModels),
		baseURL:        baseURL,
		client:         httpClient,
	}
}

type generateContentRequest struct {
	Contents         []content        `json:"contents"`
	GenerationConfig generationConfig `json:"generationConfig"`
}

type generationConfig struct {
	ResponseMimeType string `json:"responseMimeType,omitempty"`
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

type googleErrorResponse struct {
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
		return nil, ErrDisabled
	}
	models := c.models()
	if len(models) == 0 {
		return nil, fmt.Errorf("gemini model list is empty")
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
		GenerationConfig: generationConfig{
			ResponseMimeType: "application/json",
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal gemini request: %w", err)
	}

	var errs []error
	for _, model := range models {
		ids, err := c.generateProjectIDs(ctx, model, body)
		if err == nil {
			return ids, nil
		}
		errs = append(errs, err)
		if !isModelFallbackError(err) {
			break
		}
	}

	return nil, errors.Join(errs...)
}

func (c *Client) generateProjectIDs(ctx context.Context, model string, body []byte) ([]int64, error) {
	endpoint := fmt.Sprintf("%s/models/%s:generateContent?key=%s", c.baseURL, url.PathEscape(model), url.QueryEscape(c.apiKey))

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
		apiErr := &APIError{
			Model:      model,
			StatusCode: resp.StatusCode,
			Body:       strings.TrimSpace(string(snippet)),
		}
		var googleErr googleErrorResponse
		if err := json.Unmarshal(snippet, &googleErr); err == nil && googleErr.Error != nil {
			apiErr.ProviderCode = googleErr.Error.Code
			apiErr.ProviderStatus = googleErr.Error.Status
			apiErr.ProviderMessage = googleErr.Error.Message
		}
		return nil, apiErr
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

type APIError struct {
	Model           string
	StatusCode      int
	ProviderCode    int
	ProviderStatus  string
	ProviderMessage string
	Body            string
}

func (e *APIError) Error() string {
	if e.ProviderStatus != "" || e.ProviderMessage != "" {
		return fmt.Sprintf(
			"gemini model %q returned status %d provider_code=%d provider_status=%s provider_message=%q",
			e.Model,
			e.StatusCode,
			e.ProviderCode,
			e.ProviderStatus,
			e.ProviderMessage,
		)
	}
	return fmt.Sprintf("gemini model %q returned status %d: %s", e.Model, e.StatusCode, e.Body)
}

func isModelFallbackError(err error) bool {
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		return false
	}
	return apiErr.StatusCode == http.StatusBadRequest || apiErr.StatusCode == http.StatusNotFound
}

func (c *Client) models() []string {
	models := []string{c.model, DefaultModel}
	models = append(models, c.fallbackModels...)
	return uniqueModels(models)
}

func normalizeModels(models []string) []string {
	normalized := make([]string, 0, len(models))
	for _, model := range models {
		if strings.TrimSpace(model) == "" {
			continue
		}
		model = normalizeModel(model)
		if model != "" {
			normalized = append(normalized, model)
		}
	}
	return normalized
}

func normalizeModel(model string) string {
	model = strings.TrimSpace(model)
	model = strings.TrimPrefix(model, "models/")
	if model == "" {
		return DefaultModel
	}
	return model
}

func uniqueModels(models []string) []string {
	seen := make(map[string]struct{}, len(models))
	unique := make([]string, 0, len(models))
	for _, model := range models {
		model = strings.TrimSpace(model)
		if model == "" {
			continue
		}
		if _, ok := seen[model]; ok {
			continue
		}
		seen[model] = struct{}{}
		unique = append(unique, model)
	}
	return unique
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
