package commands

import (
	"fmt"
	"strings"
	"time"

	"geowar-bot/config"
	"github.com/bwmarrin/discordgo"
)

func NewPollCommands(s *discordgo.Session) *Command {
	return &Command{
		Name:        "poll",
		Aliases:     []string{"vote"},
		Description: "Create a poll. Separate question and options with |",
		Usage:       "poll <question> | <option1> | <option2> ...",
		Category:    "Utilities",
		Run: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
			if len(args) < 1 {
				s.ChannelMessageSendEmbed(m.ChannelID, ErrorEmbed("Usage: `_poll <question> | <option1> | <option2> ...`"))
				return
			}
			content := strings.Join(args, " ")
			parts := strings.Split(content, "|")
			if len(parts) < 2 {
				s.ChannelMessageSendEmbed(m.ChannelID, ErrorEmbed("Separate question and options with `|`"))
				return
			}

			question := strings.TrimSpace(parts[0])
			emojis := []string{"1️⃣", "2️⃣", "3️⃣", "4️⃣", "5️⃣", "6️⃣", "7️⃣", "8️⃣", "9️⃣", "🔟"}
			var options []string
			for i, p := range parts[1:] {
				if i >= len(emojis) {
					break
				}
				options = append(options, fmt.Sprintf("%s %s", emojis[i], strings.TrimSpace(p)))
			}
			if len(options) < 2 {
				s.ChannelMessageSendEmbed(m.ChannelID, ErrorEmbed("Need at least 2 options."))
				return
			}

			embed := &discordgo.MessageEmbed{
				Title:       fmt.Sprintf("📊 Poll: %s", question),
				Description: strings.Join(options, "\n"),
				Color:       config.EmbedColor,
				Footer:      &discordgo.MessageEmbedFooter{Text: fmt.Sprintf("Poll by %s", m.Author.Username)},
				Timestamp:   time.Now().Format(time.RFC3339),
			}

			msg, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
			if err != nil {
				s.ChannelMessageSendEmbed(m.ChannelID, ErrorEmbed("Failed to create poll."))
				return
			}

			for i := range options {
				s.MessageReactionAdd(m.ChannelID, msg.ID, emojis[i])
			}
		},
	}
}
