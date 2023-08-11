package handler

import (
	"chatgpt-plus-exts/core"
	"chatgpt-plus-exts/modules/mj"
	"chatgpt-plus-exts/utils/resp"
	"chatgpt-plus-exts/vo"
	"github.com/gin-gonic/gin"
)

// MidJourney API handler implementations

type MidJourneyHandler struct {
	BaseHandler
	client *mj.MidJourneyClient
}

func NewMidJourneyHandler(app *core.AppServer, client *mj.MidJourneyClient) *MidJourneyHandler {
	h := MidJourneyHandler{client: client}
	h.App = app
	return &h
}

func (h *MidJourneyHandler) Image(c *gin.Context) {
	var data struct {
		Prompt string `json:"prompt"`
	}
	if err := c.ShouldBindJSON(&data); err != nil || data.Prompt == "" {
		resp.ERROR(c, vo.InvalidArgs)
		return
	}

	logger.Infof("收到绘画任务请求：%+v", data.Prompt)
	if err := h.client.Imagine(&mj.ImagineRequest{
		GuildID:   h.App.Config.MidJourneyConfig.GuildId,
		ChannelID: h.App.Config.MidJourneyConfig.ChanelId,
		Prompt:    data.Prompt,
	}); err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	resp.SUCCESS(c)
}
