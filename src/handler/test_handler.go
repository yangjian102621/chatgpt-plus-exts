package handler

import (
	"chatgpt-plus-exts/core"
	"chatgpt-plus-exts/store"
	"github.com/gin-gonic/gin"
)

type TestHandler struct {
	BaseHandler
	mq *store.RedisQueue
}

func NewTestHandler(app *core.AppServer, mqs *store.RedisMQs) *TestHandler {
	h := TestHandler{mq: mqs.WeChat}
	h.App = app
	return &h
}

func (h *TestHandler) TestWechat(c *gin.Context) {

}
