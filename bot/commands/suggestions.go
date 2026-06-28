package commands

import (
	"fmt"
	"strings"

	"geowar-bot/config"
	"github.com/bwmarrin/discordgo"
)

func NewSuggestionCommands(s *discordgo.Session) *Command {
	return &Command{
		Name:        "suggest",
		Aliases:     []string{"suggestion"},
		Description: "Submit a suggestion for the server",
		Usage:       "suggest <suggestion>",
		Category:    "Utilities",
		Run: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
			if len(args) < 1 {
				s.ChannelMessageSendEmbed(m.ChannelID, ErrorEmbed("Usage: `_suggest <suggestion>`"))
				return
			}
			suggestion := strings.Join(args, " ")

			var targetChanID string
			channels, _ := s.GuildChannels(m.GuildID)
			for _, ch := range channels {
				if ch.Name == config.SuggestionsChannel {
					targetChanID = ch.ID
					break
				}
			}
			if targetChanID == "" {
				targetChanID = m.ChannelID
			}

			embed := &discordgo.MessageEmbed{
				Title:       "New Suggestion",
				Description: suggestion,
				Color:       config.EmbedColor,
				Author: &discordgo.MessageEmbedAuthor{
					Name:    m.Author.Username,
					IconURL: m.Author.AvatarURL(""),
				},
				Footer: &discordgo.MessageEmbedFooter{Text: fmt.Sprintf("ID: %s", m.Author.ID)},
			}

			msg, err := s.ChannelMessageSendEmbed(targetChanID, embed)
			if err != nil {
				s.ChannelMessageSendEmbed(m.ChannelID, ErrorEmbed("Failed to post suggestion."))
				return
			}

			s.MessageReactionAdd(targetChanID, msg.ID, "✅")
			s.MessageReactionAdd(targetChanID, msg.ID, "❌")
			s.ChannelMessageSendEmbed(m.ChannelID, SuccessEmbed("Your suggestion has been posted!"))
		},
	}
}
