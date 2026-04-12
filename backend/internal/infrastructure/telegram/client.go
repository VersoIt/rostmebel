package telegram

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	defaultTimeout       = 10 * time.Second
	maxTelegramErrorBody = 8 * 1024
)

type Client struct {
	token  string
	chatID string
	http   *http.Client
}

type OrderNotification struct {
	Name          string
	Phone         string
	Product       string
	Comment       string
	ProjectType   string
	BudgetRange   string
	City          string
	ContactMethod string
}

func NewClient(token, chatID string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: defaultTimeout}
	}

	return &Client{
		token:  token,
		chatID: chatID,
		http:   httpClient,
	}
}

func (c *Client) SendOrderNotification(ctx context.Context, order OrderNotification) error {
	if strings.TrimSpace(c.token) == "" || strings.TrimSpace(c.chatID) == "" {
		return nil
	}

	text := buildOrderText(order)
	endpoint := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", c.token)

	form := url.Values{
		"chat_id": {c.chatID},
		"text":    {text},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, strings.NewReader(form.Encode()))
	if err != nil {
		return fmt.Errorf("create telegram request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("execute telegram request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, maxTelegramErrorBody))
		return fmt.Errorf("telegram api returned status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	return nil
}

func buildOrderText(order OrderNotification) string {
	lines := []string{
		"Новая заявка РОСТ Мебель",
		"",
		"Клиент: " + safeValue(order.Name),
		"Телефон: " + safeValue(order.Phone),
		"Проект: " + safeValue(order.Product),
	}

	appendOptional := func(label, value string) {
		value = strings.TrimSpace(value)
		if value != "" {
			lines = append(lines, label+": "+value)
		}
	}

	appendOptional("Тип", order.ProjectType)
	appendOptional("Бюджет", order.BudgetRange)
	appendOptional("Город", order.City)
	appendOptional("Связь", contactMethodLabel(order.ContactMethod))
	appendOptional("Комментарий", order.Comment)

	return strings.Join(lines, "\n")
}

func safeValue(v string) string {
	v = strings.TrimSpace(v)
	if v == "" {
		return "—"
	}
	return v
}

func contactMethodLabel(method string) string {
	switch strings.TrimSpace(strings.ToLower(method)) {
	case "phone":
		return "Звонок"
	case "whatsapp":
		return "WhatsApp"
	case "telegram":
		return "Telegram"
	case "email":
		return "Email"
	default:
		if method == "" {
			return "—"
		}
		return method
	}
}
