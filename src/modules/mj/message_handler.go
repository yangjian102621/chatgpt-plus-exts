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
	Image   Image      `json:"image"`
	Content string     `json:"content"`
	Status  TaskStatus `json:"status"`
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

	if strings.Contains(m.Content, "(Waiting to start)") && !strings.Contains(m.Content, "Rerolling **") {
		// parse content
		req := CBReq{Content: extractPrompt(m.Content), Status: Start}
		b.mq.Push(req)
		return
	}
	for _, attachment := range m.Attachments {
		if attachment.Width == 0 || attachment.Height == 0 || attachment.ContentType != "image/png" {
			continue
		}
		var image Image
		err := utils.CopyObject(attachment, &image)
		if err != nil {
			logger.Error("Error with copy object: ", err)
			continue
		}
		req := CBReq{Image: image, Content: extractPrompt(m.Content), Status: Finished}
		b.mq.Push(req)
		break // only get one image
	}
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

	if strings.Contains(m.Content, "(Stopped)") {
		req := CBReq{Content: extractPrompt(m.Content), Status: Stopped}
		b.mq.Push(req)
		return
	}
	for _, attachment := range m.Attachments {
		if attachment.Width == 0 || attachment.Height == 0 || attachment.ContentType != "image/png" {
			continue
		}
		var image Image
		err := utils.CopyObject(attachment, &image)
		if err != nil {
			logger.Error("Error with copy object: ", err)
			continue
		}
		req := CBReq{Image: image, Content: extractPrompt(m.Content), Status: Running}
		b.mq.Push(req)
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
