package handler

import (
	"chatgpt-plus-exts/core"
	wexin "chatgpt-plus-exts/modules/weixin"
	"chatgpt-plus-exts/store"
	"chatgpt-plus-exts/utils/resp"
	"chatgpt-plus-exts/vo"
	"github.com/gin-gonic/gin"
)

// MidJourney API handler implementations

type MidJourneyHandler struct {
	BaseHandler
	mq *store.RedisQueue
}

func NewMidJourneyHandler(app *core.AppServer, mqs *store.RedisMQs) *MidJourneyHandler {
	h := MidJourneyHandler{mq: mqs.WeChat}
	h.App = app
	return &h
}

func (h *MidJourneyHandler) Image(c *gin.Context) {
	var data wexin.Transaction
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, vo.InvalidArgs)
		return
	}
	h.mq.Push(data)
	resp.SUCCESS(c)
}
