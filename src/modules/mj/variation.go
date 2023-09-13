package mj

import (
	"fmt"
	"time"
)

type VariationRequest struct {
	Index       int32  `json:"index"`
	GuildID     string `json:"guild_id"`
	ChannelID   string `json:"channel_id"`
	MessageID   string `json:"message_id"`
	MessageHash string `json:"message_hash"`
}

// Variation  以指定的图片的视角进行变换再创作，注意需要在对应的频道中关闭 Remix 变换，否则 Variation 指令将不会生效
func (c *MidJourneyClient) Variation(variationReq *VariationRequest) error {
	flags := 0
	interactionsReq := &InteractionsRequest{
		Type:          3,
		ApplicationID: ApplicationID,
		GuildID:       variationReq.GuildID,
		ChannelID:     variationReq.ChannelID,
		MessageFlags:  &flags,
		MessageID:     &variationReq.MessageID,
		SessionID:     SessionID,
		Data: map[string]any{
			"component_type": 2,
			"custom_id":      fmt.Sprintf("MJ::JOB::variation::%d::%s", variationReq.Index, variationReq.MessageHash),
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
