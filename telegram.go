package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type TelegramNotifier struct {
	token  string
	chatID string
	client *http.Client
}

func NewTelegramNotifier(token, chatID string) *TelegramNotifier {
	if token == "" || chatID == "" {
		return nil
	}
	return &TelegramNotifier{
		token:  token,
		chatID: chatID,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func escapeHTML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	return s
}

func (t *TelegramNotifier) Send(ctx context.Context, msg Message) error {
	if t == nil {
		return nil
	}

	sender := msg.Sender
	if sender == "" {
		sender = "Unknown"
	}

	var text string
	if msg.HasLocation() {
		text = fmt.Sprintf("📍 <b>%s</b>\n%s\nhttps://www.google.com/maps?q=%f,%f",
			escapeHTML(sender), escapeHTML(msg.Text), *msg.Lat, *msg.Lon)
	} else {
		text = fmt.Sprintf("<b>%s</b>\n%s", escapeHTML(sender), escapeHTML(msg.Text))
	}

	if err := t.sendMessage(ctx, text); err != nil {
		return fmt.Errorf("sendMessage: %w", err)
	}

	if msg.HasLocation() {
		if err := t.sendLocation(ctx, *msg.Lat, *msg.Lon); err != nil {
			slog.Warn("failed to send location to Telegram", "error", err)
		}
	}

	return nil
}

func (t *TelegramNotifier) sendMessage(ctx context.Context, text string) error {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", t.token)
	resp, err := t.client.PostForm(apiURL, url.Values{
		"chat_id":    {t.chatID},
		"text":       {text},
		"parse_mode": {"HTML"},
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var result struct {
			Description string `json:"description"`
		}
		json.NewDecoder(resp.Body).Decode(&result)
		return fmt.Errorf("telegram API error %d: %s", resp.StatusCode, result.Description)
	}
	return nil
}

func (t *TelegramNotifier) sendLocation(ctx context.Context, lat, lon float64) error {
	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendLocation", t.token)
	resp, err := t.client.PostForm(apiURL, url.Values{
		"chat_id":   {t.chatID},
		"latitude":  {fmt.Sprintf("%f", lat)},
		"longitude": {fmt.Sprintf("%f", lon)},
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var result struct {
			Description string `json:"description"`
		}
		json.NewDecoder(resp.Body).Decode(&result)
		return fmt.Errorf("telegram API error %d: %s", resp.StatusCode, result.Description)
	}
	return nil
}
