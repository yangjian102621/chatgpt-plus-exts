package mj

import (
	"strings"

	discord "github.com/bwmarrin/discordgo"
)

type Scene string

const (
	// FirstTrigger /** 首次触发生成 */
	FirstTrigger Scene = "FirstTrigger"
	// GenerateEnd /** 生成图片结束 */
	GenerateEnd Scene = "GenerateEnd"
	// GenerateEditError /** 发送的指令midjourney生成过程中发现错误 */
	GenerateEditError Scene = "GenerateEditError"
	// RichText 富文本
	RichText Scene = "RichText"
)

func DiscordMsgCreate(s *discord.Session, m *discord.MessageCreate) {
	// 过滤频道

	// 过滤掉自己发送的消息
	//if m.Author.ID == s.State.User.ID {
	//	return
	//}

	/******** *********/
	/******** *********/
	logger.Info("content: ", m.Content)
	logger.Info("Attachments: ", m.Attachments)
	/******** *********/

	if strings.Contains(m.Content, "(Waiting to start)") && !strings.Contains(m.Content, "Rerolling **") {
		trigger(m.Content, FirstTrigger)
		return
	}
	for _, attachment := range m.Attachments {
		if attachment.Width > 0 && attachment.Height > 0 {
			replay(m)
			return
		}
	}
}

func DiscordMsgUpdate(s *discord.Session, m *discord.MessageUpdate) {
	// 过滤频道
	//if m.Author == nil {
	//	return
	//}

	// 过滤掉自己发送的消息
	//if m.Author.ID == s.State.User.ID {
	//	return
	//}

	/******** *********/
	logger.Info("content: ", m.Content)
	logger.Info("Attachments: ", m.Attachments)

	if strings.Contains(m.Content, "(Stopped)") {
		trigger(m.Content, GenerateEditError)
		return
	}
	if len(m.Embeds) > 0 {
		send(m.Embeds)
		return
	}
}

type ReqCb struct {
	Embeds  []*discord.MessageEmbed `json:"embeds,omitempty"`
	Discord *discord.MessageCreate  `json:"discord,omitempty"`
	Content string                  `json:"content,omitempty"`
	Type    Scene                   `json:"type"`
}

func replay(m *discord.MessageCreate) {
	body := ReqCb{
		Discord: m,
		Type:    GenerateEnd,
	}
	callback(body)
}

func send(embeds []*discord.MessageEmbed) {
	body := ReqCb{
		Embeds: embeds,
		Type:   RichText,
	}
	callback(body)
}

func trigger(content string, t Scene) {
	body := ReqCb{
		Content: content,
		Type:    t,
	}
	callback(body)
}

func callback(params interface{}) {
	logger.Infof("请求回调接口：%+v", params)
	// TODO 发送回调请求到应用
}
