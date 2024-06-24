package client

import (
	"bytes"
	"encoding/json"
	"net/http"

	"golang.org/x/xerrors"
)

type postBody struct {
	Content string `json:"content"`
}

func SendToDiscordWebhook(url string, text string) error {
	body, err := json.Marshal(postBody{
		Content: text,
	})
	if err != nil {
		return xerrors.Errorf(": %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		return xerrors.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	return nil
}

func SendToDiscordWebhookHandler(url string) func(string) error {
	return func(text string) error {
		err := SendToDiscordWebhook(url, text)
		if err != nil {
			return xerrors.Errorf(": %w", err)
		}
		return nil
	}
}
