package handlers

import (
	"fmt"
	"time"

	"geowar-bot/config"
	"github.com/bwmarrin/discordgo"
)

type LoggingHandler struct{}

func NewLoggingHandler(s *discordgo.Session) *LoggingHandler {
	return &LoggingHandler{}
}

func (h *LoggingHandler) getLogChannel(s *discordgo.Session, guildID string) string {
	channels, err := s.GuildChannels(guildID)
	if err != nil {
		return ""
	}
	for _, ch := range channels {
		if ch.Name == config.LogChannel {
			return ch.ID
		}
	}
	return ""
}

func (h *LoggingHandler) HandleMessageDelete(s *discordgo.Session, m *discordgo.MessageDelete) {
	logChanID := h.getLogChannel(s, m.GuildID)
	if logChanID == "" {
		return
	}

	embed := &discordgo.MessageEmbed{
		Title:     "Message Deleted",
		Color:     config.ErrorColor,
		Fields: []*discordgo.MessageEmbedField{
			{Name: "Channel", Value: fmt.Sprintf("<#%s>", m.ChannelID), Inline: true},
			{Name: "Message ID", Value: m.ID, Inline: true},
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	if m.Message != nil && m.Message.Author != nil {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name: "Author", Value: fmt.Sprintf("<@%s>", m.Message.Author.ID), Inline: true,
		})
		if m.Message.Content != "" {
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
				Name: "Content", Value: m.Message.Content, Inline: false,
			})
		}
	}

	s.ChannelMessageSendEmbed(logChanID, embed)
}

func (h *LoggingHandler) HandleMessageUpdate(s *discordgo.Session, m *discordgo.MessageUpdate) {
	logChanID := h.getLogChannel(s, m.GuildID)
	if logChanID == "" || m.Message == nil {
		return
	}

	if m.Message.Author == nil || m.Message.Author.Bot {
		return
	}

	if m.Message.Content == "" {
		return
	}

		msg, err := s.State.Message(m.ChannelID, m.ID)
	if err != nil {
		return
	}
	if msg.Content == m.Message.Content {
		return
	}

	embed := &discordgo.MessageEmbed{
		Title: "Message Edited",
		Color: config.WarningColor,
		Fields: []*discordgo.MessageEmbedField{
			{Name: "Author", Value: fmt.Sprintf("<@%s>", m.Message.Author.ID), Inline: true},
			{Name: "Channel", Value: fmt.Sprintf("<#%s>", m.ChannelID), Inline: true},
			{Name: "Before", Value: truncate(msg.Content, 900), Inline: false},
			{Name: "After", Value: truncate(m.Message.Content, 900), Inline: false},
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	s.ChannelMessageSendEmbed(logChanID, embed)
}

func (h *LoggingHandler) HandleMemberAdd(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
	logChanID := h.getLogChannel(s, m.GuildID)
	if logChanID == "" {
		return
	}

	embed := &discordgo.MessageEmbed{
		Title:     "Member Joined",
		Color:     config.SuccessColor,
		Thumbnail: &discordgo.MessageEmbedThumbnail{URL: m.User.AvatarURL("")},
		Fields: []*discordgo.MessageEmbedField{
			{Name: "User", Value: fmt.Sprintf("<@%s>", m.User.ID), Inline: true},
			{Name: "ID", Value: m.User.ID, Inline: true},
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	s.ChannelMessageSendEmbed(logChanID, embed)
}

func (h *LoggingHandler) HandleMemberRemove(s *discordgo.Session, m *discordgo.GuildMemberRemove) {
	logChanID := h.getLogChannel(s, m.GuildID)
	if logChanID == "" {
		return
	}

	embed := &discordgo.MessageEmbed{
		Title:     "Member Left",
		Color:     config.ErrorColor,
		Thumbnail: &discordgo.MessageEmbedThumbnail{URL: m.User.AvatarURL("")},
		Fields: []*discordgo.MessageEmbedField{
			{Name: "User", Value: fmt.Sprintf("<@%s>", m.User.ID), Inline: true},
			{Name: "ID", Value: m.User.ID, Inline: true},
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	s.ChannelMessageSendEmbed(logChanID, embed)
}

func (h *LoggingHandler) HandleVoiceStateUpdate(s *discordgo.Session, v *discordgo.VoiceStateUpdate) {
	logChanID := h.getLogChannel(s, v.GuildID)
	if logChanID == "" {
		return
	}

	if v.BeforeUpdate == nil && v.ChannelID != "" {
		embed := &discordgo.MessageEmbed{
			Title: "Voice: User Joined",
			Color: config.SuccessColor,
			Fields: []*discordgo.MessageEmbedField{
				{Name: "User", Value: fmt.Sprintf("<@%s>", v.UserID), Inline: true},
				{Name: "Channel", Value: fmt.Sprintf("<#%s>", v.ChannelID), Inline: true},
			},
			Timestamp: time.Now().Format(time.RFC3339),
		}
		s.ChannelMessageSendEmbed(logChanID, embed)
	} else if v.ChannelID == "" && v.BeforeUpdate != nil {
		embed := &discordgo.MessageEmbed{
			Title: "Voice: User Left",
			Color: config.ErrorColor,
			Fields: []*discordgo.MessageEmbedField{
				{Name: "User", Value: fmt.Sprintf("<@%s>", v.UserID), Inline: true},
				{Name: "Channel", Value: fmt.Sprintf("<#%s>", v.BeforeUpdate.ChannelID), Inline: true},
			},
			Timestamp: time.Now().Format(time.RFC3339),
		}
		s.ChannelMessageSendEmbed(logChanID, embed)
	} else if v.BeforeUpdate != nil && v.ChannelID != "" && v.BeforeUpdate.ChannelID != v.ChannelID {
		embed := &discordgo.MessageEmbed{
			Title: "Voice: User Moved",
			Color: config.WarningColor,
			Fields: []*discordgo.MessageEmbedField{
				{Name: "User", Value: fmt.Sprintf("<@%s>", v.UserID), Inline: true},
				{Name: "From", Value: fmt.Sprintf("<#%s>", v.BeforeUpdate.ChannelID), Inline: true},
				{Name: "To", Value: fmt.Sprintf("<#%s>", v.ChannelID), Inline: true},
			},
			Timestamp: time.Now().Format(time.RFC3339),
		}
		s.ChannelMessageSendEmbed(logChanID, embed)
	}
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
