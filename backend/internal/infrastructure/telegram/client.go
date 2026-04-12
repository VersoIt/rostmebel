package telegram

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
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

func NewClient(token, chatID string) *Client {
	return &Client{
		token:  token,
		chatID: chatID,
		http:   &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *Client) SendOrderNotification(order OrderNotification) error {
	if c.token == "" || c.chatID == "" {
		return nil // Service disabled
	}

	text := buildOrderText(order)

	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", c.token)

	resp, err := c.http.PostForm(apiURL, url.Values{
		"chat_id": {c.chatID},
		"text":    {text},
	})

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram api returned status %d", resp.StatusCode)
	}

	return nil
}

func buildOrderText(order OrderNotification) string {
	lines := []string{
		"Новая заявка РОСТ Мебель",
		"",
		"Клиент: " + order.Name,
		"Телефон: " + order.Phone,
		"Проект: " + order.Product,
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

func contactMethodLabel(method string) string {
	switch method {
	case "phone":
		return "Звонок"
	case "whatsapp":
		return "WhatsApp"
	case "telegram":
		return "Telegram"
	case "email":
		return "Email"
	default:
		return method
	}
}
