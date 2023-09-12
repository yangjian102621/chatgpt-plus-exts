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
	Hash     string `json:"hash"`
}

type CBReq struct {
	MessageId   string     `json:"message_id"`
	ReferenceId string     `json:"reference_id"`
	Image       Image      `json:"image"`
	Content     string     `json:"content"`
	Prompt      string     `json:"prompt"`
	Status      TaskStatus `json:"status"`
	Progress    int        `json:"progress"`
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
	var referenceId = ""
	if m.ReferencedMessage != nil {
		referenceId = m.ReferencedMessage.ID
	}
	if strings.Contains(m.Content, "(Waiting to start)") && !strings.Contains(m.Content, "Rerolling **") {
		// parse content
		req := CBReq{
			MessageId:   m.ID,
			ReferenceId: referenceId,
			Prompt:      extractPrompt(m.Content),
			Content:     m.Content,
			Progress:    0,
			Status:      Start}
		b.mq.RPush(req)
		return
	}

	b.addAttachment(m.ID, referenceId, m.Content, m.Attachments)
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

	var referenceId = ""
	if m.ReferencedMessage != nil {
		referenceId = m.ReferencedMessage.ID
	}
	if strings.Contains(m.Content, "(Stopped)") {
		req := CBReq{
			MessageId:   m.ID,
			ReferenceId: referenceId,
			Prompt:      extractPrompt(m.Content),
			Content:     m.Content,
			Progress:    extractProgress(m.Content),
			Status:      Stopped}
		b.mq.RPush(req)
		return
	}

	b.addAttachment(m.ID, referenceId, m.Content, m.Attachments)

}

func (b *MidJourneyBot) addAttachment(messageId string, referenceId string, content string, attachments []*discord.MessageAttachment) {
	progress := extractProgress(content)
	var status TaskStatus
	if progress == 100 {
		status = Finished
	} else {
		status = Running
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
			MessageId:   messageId,
			ReferenceId: referenceId,
			Image:       image,
			Prompt:      extractPrompt(content),
			Content:     content,
			Progress:    progress,
			Status:      status,
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
		return strings.TrimSpace(matches[1])
	}
	return ""
}

func extractProgress(input string) int {
	pattern := `\((\d+)\%\)`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(input)
	if len(matches) > 1 {
		return utils.IntValue(matches[1], 0)
	}
	return 100
}

func extractHashFromFilename(filename string) string {
	if !strings.HasSuffix(filename, ".png") {
		return ""
	}

	index := strings.LastIndex(filename, "_")
	if index != -1 {
		return filename[index+1 : len(filename)-4]
	}
	return ""
}
