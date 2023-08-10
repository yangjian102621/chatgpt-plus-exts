package mj

import (
	"chatgpt-plus-exts/core"
	logger2 "chatgpt-plus-exts/logger"
	"github.com/bwmarrin/discordgo"
	"github.com/gorilla/websocket"
	"net/http"
	"net/url"
)

// MidJourney 机器人

var logger = logger2.GetLogger()

type MidJourneyBot struct {
	config *core.MidJourneyConfig
	bot    *discordgo.Session
}

func NewMidJourneyBot(config *core.Config) (*MidJourneyBot, error) {
	discord, err := discordgo.New("Bot " + config.MidJourneyConfig.BotToken)
	if err != nil {
		return nil, err
	}

	discord.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentMessageContent
	discord.AddHandler(DiscordMsgCreate)
	discord.AddHandler(DiscordMsgUpdate)

	if config.ProxyURL != "" {
		proxy, _ := url.Parse(config.ProxyURL)
		discord.Client = &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxy),
			},
		}
		discord.Dialer = &websocket.Dialer{
			Proxy: http.ProxyURL(proxy),
		}
	}

	return &MidJourneyBot{
		config: &config.MidJourneyConfig,
		bot:    discord,
	}, nil
}

func (b *MidJourneyBot) Run() {
	logger.Info("Starting MidJourney Bot...")
	err := b.bot.Open()
	if err != nil {
		logger.Error("Error opening Discord connection:", err)
		return
	}
}
