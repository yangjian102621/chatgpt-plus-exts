package mj

import (
	"chatgpt-plus-exts/utils"
	"regexp"
	"strings"

	discord "github.com/bwmarrin/discordgo"
)

type TaskStatus string

const (
	Start    = TaskStatus("Started")
	Running  = TaskStatus("Running")
	Stopped  = TaskStatus("Stopped")
	Finished = TaskStatus("Finished")
)

type Image struct {
	URL      string `json:"url"`
	ProxyURL string `json:"proxy_url"`
	Filename string `json:"filename"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	Size     int    `json:"size"`
}

type CBReq struct {
	MessageId   string
	MessageHash string
	Image       Image      `json:"image"`
	Content     string     `json:"content"`
	Status      TaskStatus `json:"status"`
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
		req := CBReq{Content: extractPrompt(m.Content), Status: Start}
		b.mq.RPush(req)
		return
	}

	b.addAttachment(m.Content, m.Attachments)
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
		req := CBReq{Content: extractPrompt(m.Content), Status: Stopped}
		b.mq.RPush(req)
		return
	}

	b.addAttachment(m.Content, m.Attachments)

}

func (b *MidJourneyBot) addAttachment(content string, attachments []*discord.MessageAttachment) {
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
		}
		req := CBReq{Image: image, Content: extractPrompt(content), Status: status}
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
