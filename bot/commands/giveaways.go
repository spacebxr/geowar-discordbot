package commands

import (
	"fmt"
	"strings"
	"time"

	"geowar-bot/config"
	"geowar-bot/utils"
	"github.com/bwmarrin/discordgo"
)

func NewGiveawayCommands(s *discordgo.Session) *Command {
	return &Command{
		Name:        "giveaway",
		Aliases:     []string{"ga"},
		Description: "Giveaway commands: start, reroll, end",
		Usage:       "giveaway <start|reroll|end> <duration> <winners> <prize>",
		Category:    "Giveaways",
		StaffOnly:   true,
		Run: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
			if len(args) < 1 {
				s.ChannelMessageSendEmbed(m.ChannelID, ErrorEmbed("Usage: `_giveaway start|reroll|end`"))
				return
			}
			switch args[0] {
			case "start":
				startGiveaway(s, m, args[1:])
			case "reroll":
				rerollGiveaway(s, m, args[1:])
			case "end":
				endGiveaway(s, m, args[1:])
			default:
				s.ChannelMessageSendEmbed(m.ChannelID, ErrorEmbed("Usage: `_giveaway start|reroll|end`"))
			}
		},
	}
}

func startGiveaway(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(args) < 3 {
		s.ChannelMessageSendEmbed(m.ChannelID, ErrorEmbed("Usage: `_giveaway start <duration> <winners> <prize>`"))
		return
	}

	duration, ok := utils.ParseDuration(args[0])
	if !ok {
		s.ChannelMessageSendEmbed(m.ChannelID, ErrorEmbed("Invalid duration. Use format like `10m`, `1h`, `7d`"))
		return
	}

	winners := 0
	fmt.Sscanf(args[1], "%d", &winners)
	if winners < 1 || winners > 10 {
		s.ChannelMessageSendEmbed(m.ChannelID, ErrorEmbed("Winners must be between 1 and 10."))
		return
	}

	prize := strings.Join(args[2:], " ")
	endTime := time.Now().Add(duration)

	embed := &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("🎉 Giveaway: %s", prize),
		Description: fmt.Sprintf("React with 🎉 to enter!\n**Winners:** %d\n**Host:** %s\n**Ends:** <t:%d:R>", winners, m.Author.Mention(), endTime.Unix()),
		Color:       config.EmbedColor,
		Footer:      &discordgo.MessageEmbedFooter{Text: fmt.Sprintf("Ends at")},
		Timestamp:   endTime.Format(time.RFC3339),
	}

	msg, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
	if err != nil {
		s.ChannelMessageSendEmbed(m.ChannelID, ErrorEmbed("Failed to create giveaway."))
		return
	}

	s.MessageReactionAdd(m.ChannelID, msg.ID, "🎉")

	go func() {
		time.Sleep(duration)
		endGiveawayByMessage(s, m.GuildID, m.ChannelID, msg.ID, prize, winners, m.Author.ID)
	}()
}

func rerollGiveaway(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(args) < 2 {
		s.ChannelMessageSendEmbed(m.ChannelID, ErrorEmbed("Usage: `_giveaway reroll <message_id> <winners>`"))
		return
	}
	msgID := args[0]
	winners := 0
	fmt.Sscanf(args[1], "%d", &winners)
	if winners < 1 {
		winners = 1
	}

	msg, err := s.ChannelMessage(m.ChannelID, msgID)
	if err != nil {
		s.ChannelMessageSendEmbed(m.ChannelID, ErrorEmbed("Message not found."))
		return
	}

	var entries []string
	for _, reaction := range msg.Reactions {
		if reaction.Emoji.Name == "🎉" {
			users, _ := s.MessageReactions(m.ChannelID, msgID, "🎉", 100, "", "")
			for _, u := range users {
				if !u.Bot {
					entries = append(entries, u.ID)
				}
			}
		}
	}

	if len(entries) == 0 {
		s.ChannelMessageSendEmbed(m.ChannelID, ErrorEmbed("No valid entries found."))
		return
	}

	var winnersList []string
	for i := 0; i < winners && i < len(entries); i++ {
				winnersList = append(winnersList, fmt.Sprintf("<@%s>", entries[i]))
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("🎉 New winners: %s! Congratulations on winning the giveaway!", strings.Join(winnersList, ", ")))
}

func endGiveaway(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	if len(args) < 2 {
		s.ChannelMessageSendEmbed(m.ChannelID, ErrorEmbed("Usage: `_giveaway end <message_id> <winners>`"))
		return
	}
	msgID := args[0]
	winners := 0
	fmt.Sscanf(args[1], "%d", &winners)
	if winners < 1 {
		winners = 1
	}
	endGiveawayByMessage(s, m.GuildID, m.ChannelID, msgID, "Giveaway", winners, m.Author.ID)
}

func endGiveawayByMessage(s *discordgo.Session, guildID, channelID, msgID, prize string, winners int, hostID string) {
	msg, err := s.ChannelMessage(channelID, msgID)
	if err != nil {
		return
	}

	var entries []string
	for _, reaction := range msg.Reactions {
		if reaction.Emoji.Name == "🎉" {
			users, _ := s.MessageReactions(channelID, msgID, "🎉", 100, "", "")
			for _, u := range users {
				if !u.Bot {
					entries = append(entries, u.ID)
				}
			}
		}
	}

	var winnersList []string
	for i := 0; i < winners && i < len(entries); i++ {
		winnersList = append(winnersList, fmt.Sprintf("<@%s>", entries[i]))
	}

	resultEmbed := &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("🎉 Giveaway Ended: %s", prize),
		Color:       config.EmbedColor,
		Timestamp:   time.Now().Format(time.RFC3339),
	}

	if len(winnersList) > 0 {
		resultEmbed.Description = fmt.Sprintf("**Winners:** %s\nCongratulations!", strings.Join(winnersList, ", "))
		s.ChannelMessageSend(channelID, fmt.Sprintf("Congratulations %s! You won **%s**!", strings.Join(winnersList, ", "), prize))
	} else {
		resultEmbed.Description = "Not enough participants."
		s.ChannelMessageSendEmbed(channelID, resultEmbed)
	}
}
