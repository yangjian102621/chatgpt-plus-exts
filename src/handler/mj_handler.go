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
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, vo.InvalidArgs)
		return
	}

	if err := h.client.Imagine(&mj.ImagineRequest{
		GuildID:   h.App.Config.MidJourneyConfig.GuildId,
		ChannelID: h.App.Config.MidJourneyConfig.ChanelId,
		Prompt:    "A chinese girl, long hair and shawl. looking at view, At the age of 15 or 16, her skin was better than snow and was beautiful. The appearance was extremely beautiful, the whole body was dressed in red, and the hair was tied with a gold band. When the snow reflected, it was even more brilliant.",
	}); err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	resp.SUCCESS(c)
}
