package wexin

import (
	"chatgpt-plus-exts/core"
	logger2 "chatgpt-plus-exts/logger"
	"chatgpt-plus-exts/store"
	"chatgpt-plus-exts/vo"
	"github.com/eatmoreapple/openwechat"
	"github.com/imroc/req/v3"
	"github.com/skip2/go-qrcode"
	"time"
)

// 微信收款机器人
var logger = logger2.GetLogger()

type WeChatBot struct {
	bot    *openwechat.Bot
	config *core.WeChatConfig
	mq     *store.RedisQueue
	token  string
}

func NewWeChatBot(config *core.Config, mqs *store.RedisMQs) *WeChatBot {
	bot := openwechat.DefaultBot(openwechat.Desktop)
	return &WeChatBot{
		bot:    bot,
		config: &config.WeChatConfig,
		token:  config.CallbackToken,
		mq:     mqs.WeChat,
	}
}

func (b *WeChatBot) Run() error {
	logger.Info("Starting WeChat Bot...")

	// set message handler
	b.bot.MessageHandler = func(msg *openwechat.Message) {
		b.messageHandler(msg)
	}
	// scan code login callback
	b.bot.UUIDCallback = b.qrCodeCallBack

	// create hot login storage
	reloadStorage := openwechat.NewJsonFileHotReloadStorage("storage.json")
	err := b.bot.HotLogin(reloadStorage, true)
	if err != nil {
		logger.Error("login error: %v", err)
		return b.bot.Login()
	}
	logger.Info("微信登录成功！")
	return nil
}

// ConsumeMessages consume messages
func (b *WeChatBot) ConsumeMessages() {
	logger.Info("Starting consume wechat messages...")
	client := req.C().SetTimeout(10 * time.Second)
	for {
		var message Transaction
		err := b.mq.Take(&message)
		if err != nil {
			logger.Errorf("taking message with error: %v", err)
			continue
		}
		var res vo.BizVo
		r, err := client.R().
			SetHeader("X-TOKEN", b.token).
			SetBody(message).
			SetSuccessResult(&res).
			Post(b.config.CallbackUrl)
		if err != nil || r.IsErrorState() || !res.Success() {
			logger.Errorf("消息推送失败：%v%v%v", err, r.Err, res.Message)
			b.mq.Push(message)
			continue
		}

		logger.Infof("推送微信转账消息成功： %+v", message)
	}
}

// message handler
func (b *WeChatBot) messageHandler(msg *openwechat.Message) {
	sender, err := msg.Sender()
	if err != nil {
		return
	}

	// 只处理微信支付的推送消息
	if sender.NickName == "微信支付" ||
		msg.MsgType == openwechat.MsgTypeApp ||
		msg.AppMsgType == openwechat.AppMsgTypeUrl {
		// 解析支付金额
		message, err := parseTransactionMessage(msg.Content)
		if err == nil {
			transaction := extractTransaction(message)
			logger.Infof("解析到收款信息：%+v", transaction)
			// push transaction to message queue
			b.mq.Push(transaction)
		}
	}
}

func (b *WeChatBot) qrCodeCallBack(uuid string) {
	logger.Info("请使用微信扫描下面二维码登录")
	q, _ := qrcode.New("https://login.weixin.qq.com/l/"+uuid, qrcode.Medium)
	logger.Info(q.ToString(true))
}
