package wexin

import (
	"github.com/eatmoreapple/openwechat"
	"github.com/skip2/go-qrcode"
	"gorm.io/gorm"
)

// MessageHandler 消息处理
func MessageHandler(msg *openwechat.Message, db *gorm.DB) {
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
			// TODO 推送消息到微信消息回调地址
		}
	}
}

// QrCodeCallBack 登录扫码回调，
func QrCodeCallBack(uuid string) {
	logger.Info("请使用微信扫描下面二维码登录")
	q, _ := qrcode.New("https://login.weixin.qq.com/l/"+uuid, qrcode.Medium)
	logger.Info(q.ToString(true))
}
