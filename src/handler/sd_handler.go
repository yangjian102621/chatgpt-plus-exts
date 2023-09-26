package handler

import (
	"chatgpt-plus-exts/core"
	"chatgpt-plus-exts/modules/sd"
	"chatgpt-plus-exts/utils"
	"chatgpt-plus-exts/utils/resp"
	"chatgpt-plus-exts/vo"
	"fmt"
	"github.com/gin-gonic/gin"
)

// StableDiffusionHandler API handler implementations

type StableDiffusionHandler struct {
	BaseHandler
	client *sd.StableDiffusionClient
}

func NewStableDiffusionHandler(app *core.AppServer, client *sd.StableDiffusionClient) *StableDiffusionHandler {
	h := StableDiffusionHandler{client: client}
	h.App = app
	return &h
}

// Tex2Img 文生图
func (h *StableDiffusionHandler) Tex2Img(c *gin.Context) {
	var data sd.Txt2ImgReq
	if err := c.ShouldBindJSON(&data); err != nil || data.Prompt == "" {
		resp.ERROR(c, vo.InvalidArgs)
		return
	}
	if !h.App.Config.SdConfig.Enabled {
		resp.ERROR(c, "Stable Diffusion service is disabled")
		return
	}

	if data.Width <= 0 {
		data.Width = 512
	}
	if data.Height <= 0 {
		data.Height = 512
	}
	if data.CfgScale <= 0 {
		data.CfgScale = 7
	}
	if data.Seed == 0 {
		data.Seed = -1
	}
	if data.Steps <= 0 {
		data.Steps = 20
	}
	if data.Sampler == "" {
		data.Sampler = "Euler a"
	}

	data.TaskId = fmt.Sprintf("task(%s)", utils.RandString(15))
	logger.Infof("收到 Image 任务：%+v", data)
	if err := h.client.Txt2Img(data); err != nil {
		resp.ERROR(c, err.Error())
		return
	}

	resp.SUCCESS(c)
}
