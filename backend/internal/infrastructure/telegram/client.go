package telegram

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	token  string
	chatID string
	http   *http.Client
}

func NewClient(token, chatID string) *Client {
	return &Client{
		token:  token,
		chatID: chatID,
		http:   &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *Client) SendOrderNotification(name, phone, product, comment string) error {
	if c.token == "" || c.chatID == "" {
		return nil // Service disabled
	}

	text := fmt.Sprintf("🛋️ **Новая заявка РОСТ Мебель**\n\n"+
		"👤 Клиент: %s\n"+
		"📞 Телефон: %s\n"+
		"📦 Проект: %s\n"+
		"💬 Комментарий: %s",
		name, phone, product, comment)

	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", c.token)
	
	resp, err := c.http.PostForm(apiURL, url.Values{
		"chat_id":    {c.chatID},
		"text":       {text},
		"parse_mode": {"Markdown"},
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
