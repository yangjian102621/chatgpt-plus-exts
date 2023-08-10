package mj

import (
	"chatgpt-plus-exts/core"
	"github.com/imroc/req/v3"
	"time"
)

// MidJourney client

type MidJourneyClient struct {
	client *req.Client
	config *core.MidJourneyConfig
}

func NewMjClient(config *core.Config) *MidJourneyClient {
	client := req.C().SetTimeout(10 * time.Second)
	// set proxy URL
	if config.ProxyURL != "" {
		client.SetProxyURL(config.ProxyURL)
	}
	return &MidJourneyClient{client: client, config: &config.MidJourneyConfig}
}
