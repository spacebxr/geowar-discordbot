package handlers

import (
	"fmt"
	"time"

	"geowar-bot/config"
	"github.com/bwmarrin/discordgo"
)

type WelcomeHandler struct{}

func NewWelcomeHandler(s *discordgo.Session) *WelcomeHandler {
	return &WelcomeHandler{}
}

func (h *WelcomeHandler) getWelcomeChannel(s *discordgo.Session, guildID string) string {
	channels, err := s.GuildChannels(guildID)
	if err != nil {
		return ""
	}
	for _, ch := range channels {
		if ch.Name == config.WelcomeChannel {
			return ch.ID
		}
	}
	return ""
}

func (h *WelcomeHandler) HandleMemberAdd(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
	welcomeChanID := h.getWelcomeChannel(s, m.GuildID)
	if welcomeChanID == "" {
		return
	}

	guild, err := s.State.Guild(m.GuildID)
	if err != nil {
		return
	}

	memberCount := guild.MemberCount

	embed := &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("Welcome to %s, %s!", guild.Name, m.User.Username),
		Description: fmt.Sprintf("Hey <@%s>, welcome to **%s**!\n\nMake sure to check out <#%s> and enjoy your stay!", m.User.ID, guild.Name, config.RulesChannel),
		Color:       config.EmbedColor,
		Thumbnail:   &discordgo.MessageEmbedThumbnail{URL: m.User.AvatarURL("")},
		Fields: []*discordgo.MessageEmbedField{
			{Name: "Member #", Value: fmt.Sprintf("%d", memberCount), Inline: true},
		},
		Timestamp: time.Now().Format(time.RFC3339),
		Footer:    &discordgo.MessageEmbedFooter{Text: "GeoWar SMP"},
	}

	s.ChannelMessageSendEmbed(welcomeChanID, embed)
}
