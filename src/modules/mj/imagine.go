package mj

import (
	"fmt"
)

type ImagineRequest struct {
	GuildID   string `json:"guild_id"`
	ChannelID string `json:"channel_id"`
	Prompt    string `json:"prompt"`
}

func (c *MidJourneyClient) Imagine(imgReq *ImagineRequest) error {
	interactionsReq := &InteractionsRequest{
		Type:          2,
		ApplicationID: ApplicationID,
		GuildID:       imgReq.GuildID,
		ChannelID:     imgReq.ChannelID,
		SessionID:     SessionID,
		Data: map[string]any{
			"version": "1118961510123847772",
			"id":      "938956540159881230",
			"name":    "imagine",
			"type":    "1",
			"options": []map[string]any{
				{
					"type":  3,
					"name":  "prompt",
					"value": imgReq.Prompt,
				},
			},
			"application_command": map[string]any{
				"id":                         "938956540159881230",
				"application_id":             ApplicationID,
				"version":                    "1118961510123847772",
				"default_permission":         true,
				"default_member_permissions": nil,
				"type":                       1,
				"nsfw":                       false,
				"name":                       "imagine",
				"description":                "Create images with Midjourney",
				"dm_permission":              true,
				"options": []map[string]any{
					{
						"type":        3,
						"name":        "prompt",
						"description": "The prompt to imagine",
						"required":    true,
					},
				},
				"attachments": []any{},
			},
		},
	}

	url := "https://discord.com/api/v9/interactions"
	r, err := c.client.R().SetHeader("Authorization", c.config.UserToken).
		SetHeader("Content-Type", "application/json").
		SetBody(interactionsReq).
		Post(url)

	if err != nil || r.IsErrorState() {
		return fmt.Errorf("error with http request: %w%v", err, r.Err)
	}

	return nil
}
