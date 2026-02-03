package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/adminlove520/vulnDb-Notifier/internal/errors"
	"github.com/adminlove520/vulnDb-Notifier/internal/util"
)

type DiscordEmbed struct {
	Title       string `json:"title"`
	URL         string `json:"url"`
	Description string `json:"description"`
	Color       int    `json:"color"`
	Fields      []struct {
		Name   string `json:"name"`
		Value  string `json:"value"`
		Inline bool   `json:"inline"`
	} `json:"fields"`
	Footer struct {
		Text string `json:"text"`
	} `json:"footer"`
}

type DiscordMessage struct {
	Embeds []DiscordEmbed `json:"embeds"`
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func getRandomColor() int {
	return rand.Intn(16777215) + 1
}

func NotifyDiscord(vulnTitle string, link string, published string, categories string, description string, discordWebhook string) error {
	description = util.RemoveHTMLTags(description)

	embed := DiscordEmbed{
		Title:       vulnTitle,
		URL:         link,
		Description: description,
		Color:       getRandomColor(),
		Fields: []struct {
			Name   string `json:"name"`
			Value  string `json:"value"`
			Inline bool   `json:"inline"`
		}{
			{
				Name:   "Published",
				Value:  published,
				Inline: true,
			},
			{
				Name:   "Categories",
				Value:  categories,
				Inline: true,
			},
		},
		Footer: struct {
			Text string `json:"text"`
		}{
			Text: "vulnDb Notifier",
		},
	}

	message := DiscordMessage{
		Embeds: []DiscordEmbed{embed},
	}

	// Encode message payload as JSON
	payload, err := json.Marshal(message)
	if err != nil {
		return &errors.SlackNotificationError{Message: "Failed to marshal message: " + err.Error()}
	}

	// Make POST request to Discord webhook
	resp, err := http.Post(discordWebhook, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return &errors.SlackNotificationError{Message: "Failed to send message, check if discord webhook is valid: " + err.Error()}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return &errors.SlackNotificationError{Message: fmt.Sprintf("Failed to send message, status code: %d", resp.StatusCode)}
	}

	return nil
}
