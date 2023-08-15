package mj

import (
	"chatgpt-plus-exts/core"
	logger2 "chatgpt-plus-exts/logger"
	"chatgpt-plus-exts/store"
	"chatgpt-plus-exts/vo"
	"github.com/bwmarrin/discordgo"
	"github.com/gorilla/websocket"
	"github.com/imroc/req/v3"
	"net/http"
	"net/url"
	"time"
)

// MidJourney 机器人

var logger = logger2.GetLogger()

type MidJourneyBot struct {
	config *core.MidJourneyConfig
	bot    *discordgo.Session
	mq     *store.RedisQueue
	token  string
}

func NewMidJourneyBot(config *core.Config, mqs *store.RedisMQs) (*MidJourneyBot, error) {
	discord, err := discordgo.New("Bot " + config.MidJourneyConfig.BotToken)
	if err != nil {
		return nil, err
	}

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
		mq:     mqs.MidJourney,
		token:  config.CallbackToken,
	}, nil
}

func (b *MidJourneyBot) Run() {
	b.bot.Identify.Intents = discordgo.IntentsAllWithoutPrivileged | discordgo.IntentsGuildMessages | discordgo.IntentMessageContent
	b.bot.AddHandler(b.messageCreate)
	b.bot.AddHandler(b.messageUpdate)

	logger.Info("Starting MidJourney Bot...")
	err := b.bot.Open()
	if err != nil {
		logger.Error("Error opening Discord connection:", err)
		return
	}
	logger.Info("Starting MidJourney Bot successfully!")
}

// ConsumeMessages consume messages
func (b *MidJourneyBot) ConsumeMessages() {
	logger.Info("Starting consume MidJourney messages...")
	client := req.C().SetTimeout(10 * time.Second)
	for {
		var message CBReq
		err := b.mq.LPop(&message)
		if err != nil {
			logger.Errorf("taking message with error: %v", err)
			continue
		}
		var res vo.BizVo
		r, err := client.R().
			SetHeader("Authorization", b.token).
			SetBody(message).
			SetSuccessResult(&res).
			Post(b.config.CallbackUrl)
		if err != nil || r.IsErrorState() || !res.Success() {
			logger.Errorf("消息推送失败，Network Err: %v, Http Err: %v, Resp Err: %v", err, r.Err, res.Message)
			b.mq.LPush(message)
			time.Sleep(time.Second)
			continue
		}
	}
}
