package mj

import (
	"chatgpt-plus-exts/utils"
	"regexp"
	"strings"

	discord "github.com/bwmarrin/discordgo"
)

type TaskStatus string
type TaskType string

const (
	Start    = TaskStatus("Started")
	Running  = TaskStatus("Running")
	Stopped  = TaskStatus("Stopped")
	Finished = TaskStatus("Finished")

	TaskImage   = TaskType("Image")   // 创建
	TaskUpScale = TaskType("Upscale") //放大
)

type Image struct {
	URL      string `json:"url"`
	ProxyURL string `json:"proxy_url"`
	Filename string `json:"filename"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	Size     int    `json:"size"`
	Hash     string `json:"hash"`
}

type CBReq struct {
	Type      TaskType   `json:"type"`
	MessageId string     `json:"message_id"`
	Image     Image      `json:"image"`
	Content   string     `json:"content"`
	Prompt    string     `json:"prompt"`
	Status    TaskStatus `json:"status"`
}

func (b *MidJourneyBot) messageCreate(s *discord.Session, m *discord.MessageCreate) {
	// ignore messages for other channels
	if m.GuildID != b.config.GuildId || m.ChannelID != b.config.ChanelId {
		return
	}
	// ignore messages for self
	if m.Author.ID == s.State.User.ID {
		return
	}

	logger.Infof("CREATE: %s", utils.JsonEncode(m))

	if strings.Contains(m.Content, "(Waiting to start)") && !strings.Contains(m.Content, "Rerolling **") {
		// parse content
		req := CBReq{
			Type:      TaskImage,
			MessageId: m.ID,
			Prompt:    extractPrompt(m.Content),
			Content:   m.Content,
			Status:    Start}
		b.mq.RPush(req)
		return
	}

	b.addAttachment(TaskImage, m.ID, m.Content, m.Attachments)
}

func (b *MidJourneyBot) messageUpdate(s *discord.Session, m *discord.MessageUpdate) {
	// ignore messages for other channels
	if m.GuildID != b.config.GuildId || m.ChannelID != b.config.ChanelId {
		return
	}
	// ignore messages for self
	if m.Author.ID == s.State.User.ID {
		return
	}

	logger.Infof("UPDATE: %s", utils.JsonEncode(m))

	if strings.Contains(m.Content, "(Stopped)") {
		req := CBReq{
			Type:      TaskImage,
			MessageId: m.ID,
			Prompt:    extractPrompt(m.Content),
			Content:   m.Content,
			Status:    Stopped}
		b.mq.RPush(req)
		return
	}

	b.addAttachment(TaskImage, m.ID, m.Content, m.Attachments)

}

func (b *MidJourneyBot) addAttachment(t TaskType, messageId string, content string, attachments []*discord.MessageAttachment) {
	pattern := `\(\d+\%\)`
	re := regexp.MustCompile(pattern)
	match := re.FindStringSubmatch(content)
	var status TaskStatus
	if len(match) > 0 {
		status = Running
	} else {
		status = Finished
	}
	for _, attachment := range attachments {
		if attachment.Width == 0 || attachment.Height == 0 {
			continue
		}
		image := Image{
			URL:      attachment.URL,
			Height:   attachment.Height,
			ProxyURL: attachment.ProxyURL,
			Width:    attachment.Width,
			Size:     attachment.Size,
			Filename: attachment.Filename,
			Hash:     extractHashFromFilename(attachment.Filename),
		}
		req := CBReq{
			Type:      t,
			MessageId: messageId,
			Image:     image,
			Prompt:    extractPrompt(content),
			Content:   content,
			Status:    status,
		}
		b.mq.RPush(req)
		break // only get one image
	}
}

// extract prompt from string
func extractPrompt(input string) string {
	pattern := `\*\*(.*?)\*\*`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(input)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func extractHashFromFilename(filename string) string {
	index := strings.LastIndex(filename, "_")
	if index != -1 {
		return filename[index+1 : len(filename)-4]
	}
	return ""
}
