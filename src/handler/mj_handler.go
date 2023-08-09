package handler

import (
	"chatgpt-plus-exts/core"
	"github.com/gin-gonic/gin"
)

// MidJourney API handler implementations

type MidJourneyHandler struct {
	BaseHandler
}

func NewMidJourneyHandler(app *core.AppServer) *MidJourneyHandler {
	h := MidJourneyHandler{}
	h.App = app
	return &h
}

func (h *MidJourneyHandler) Image(c *gin.Context) {

}
