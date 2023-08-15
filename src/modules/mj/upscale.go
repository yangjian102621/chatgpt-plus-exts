package mj

import (
	"fmt"
)

type UpscaleRequest struct {
	Index       int32  `json:"index"`
	GuildID     string `json:"guild_id"`
	ChannelID   string `json:"channel_id"`
	MessageID   string `json:"message_id"`
	MessageHash string `json:"message_hash"`
}

func (c *MidJourneyClient) Upscale(upscaleReq *UpscaleRequest) error {
	flags := 0
	interactionsReq := &InteractionsRequest{
		Type:          3,
		ApplicationID: ApplicationID,
		GuildID:       upscaleReq.GuildID,
		ChannelID:     upscaleReq.ChannelID,
		MessageFlags:  &flags,
		MessageID:     &upscaleReq.MessageID,
		SessionID:     SessionID,
		Data: map[string]any{
			"component_type": 2,
			"custom_id":      fmt.Sprintf("MJ::JOB::upsample::%d::%s", upscaleReq.Index, upscaleReq.MessageHash),
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
