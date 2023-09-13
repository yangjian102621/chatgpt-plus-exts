package mj

import (
	"fmt"
	"time"
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
		Nonce: fmt.Sprintf("%d", time.Now().UnixNano()),
	}

	url := "https://discord.com/api/v9/interactions"
	var res InteractionsResult
	r, err := c.client.R().SetHeader("Authorization", c.config.UserToken).
		SetHeader("Content-Type", "application/json").
		SetBody(interactionsReq).
		SetErrorResult(&res).
		Post(url)
	if err != nil || r.IsErrorState() {
		return fmt.Errorf("error with http request: %v%v%v", err, r.Err, res.Message)
	}

	return nil
}
