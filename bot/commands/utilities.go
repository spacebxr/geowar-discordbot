package commands

import (
	"fmt"
	"strings"
	"time"

	"geowar-bot/config"
	"github.com/bwmarrin/discordgo"
)

func NewUtilityCommands(s *discordgo.Session, cm *CommandManager) *Command {
	cmds := []*Command{
		{
			Name:        "ping",
			Description: "Check bot latency",
			Usage:       "ping",
			Category:    "Utilities",
			Run: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
				latency := s.HeartbeatLatency().Milliseconds()
				s.ChannelMessageSendEmbed(m.ChannelID, InfoEmbed("Pong! 🏓", fmt.Sprintf("Latency: **%dms**", latency)))
			},
		},
		{
			Name:        "userinfo",
			Aliases:     []string{"ui", "whois"},
			Description: "Get info about a user",
			Usage:       "userinfo [user]",
			Category:    "Utilities",
			Run: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
				targetID := m.Author.ID
				if len(args) > 0 {
					targetID = extractID(args[0])
				}
				member, err := s.GuildMember(m.GuildID, targetID)
				if err != nil {
					s.ChannelMessageSendEmbed(m.ChannelID, ErrorEmbed("User not found."))
					return
				}
				user := member.User
				joinedAt := member.JoinedAt.Format("Jan 2, 2006")
				createdAtTS, _ := discordgo.SnowflakeTimestamp(member.User.ID)
				createdAt := fmt.Sprintf("<t:%d:R>", createdAtTS.Unix())

				var roles []string
				for _, roleID := range member.Roles {
					role, err := s.State.Role(m.GuildID, roleID)
					if err == nil {
						roles = append(roles, role.Mention())
					}
				}

				embed := &discordgo.MessageEmbed{
					Title:       user.Username,
					Color:       config.EmbedColor,
					Thumbnail:   &discordgo.MessageEmbedThumbnail{URL: user.AvatarURL("")},
					Description: fmt.Sprintf("**ID:** %s\n**Mention:** <@%s>", user.ID, user.ID),
					Fields: []*discordgo.MessageEmbedField{
						{Name: "Joined Server", Value: joinedAt, Inline: true},
						{Name: "Account Created", Value: createdAt, Inline: true},
						{Name: "Roles", Value: strings.Join(roles, " "), Inline: false},
						{Name: "Is Bot", Value: fmt.Sprintf("%t", user.Bot), Inline: true},
					},
					Footer: &discordgo.MessageEmbedFooter{Text: fmt.Sprintf("ID: %s", user.ID)},
				}
				s.ChannelMessageSendEmbed(m.ChannelID, embed)
			},
		},
		{
			Name:        "serverinfo",
			Aliases:     []string{"si", "guildinfo"},
			Description: "Get info about the server",
			Usage:       "serverinfo",
			Category:    "Utilities",
			Run: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
				guild, err := s.State.Guild(m.GuildID)
				if err != nil {
					s.ChannelMessageSendEmbed(m.ChannelID, ErrorEmbed("Could not fetch server info."))
					return
				}
				memberCount := guild.MemberCount
				boostCount := guild.PremiumSubscriptionCount

				embed := &discordgo.MessageEmbed{
					Title:     guild.Name,
					Color:     config.EmbedColor,
					Thumbnail: &discordgo.MessageEmbedThumbnail{URL: guild.IconURL("")},
					Fields: []*discordgo.MessageEmbedField{
						{Name: "Owner", Value: fmt.Sprintf("<@%s>", guild.OwnerID), Inline: true},
						{Name: "Members", Value: fmt.Sprintf("%d", memberCount), Inline: true},
						{Name: "Boosts", Value: fmt.Sprintf("%d", boostCount), Inline: true},
						{Name: "Channels", Value: fmt.Sprintf("%d", len(guild.Channels)), Inline: true},
						{Name: "Roles", Value: fmt.Sprintf("%d", len(guild.Roles)), Inline: true},
					},
					Footer: &discordgo.MessageEmbedFooter{Text: fmt.Sprintf("ID: %s", guild.ID)},
				}
				s.ChannelMessageSendEmbed(m.ChannelID, embed)
			},
		},
		{
			Name:        "avatar",
			Aliases:     []string{"av", "pfp"},
			Description: "Get a user's avatar",
			Usage:       "avatar [user]",
			Category:    "Utilities",
			Run: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
				targetID := m.Author.ID
				if len(args) > 0 {
					targetID = extractID(args[0])
				}
				member, err := s.GuildMember(m.GuildID, targetID)
				if err != nil {
					s.ChannelMessageSendEmbed(m.ChannelID, ErrorEmbed("User not found."))
					return
				}
				embed := &discordgo.MessageEmbed{
					Title: fmt.Sprintf("%s's Avatar", member.User.Username),
					Color: config.EmbedColor,
					Image: &discordgo.MessageEmbedImage{URL: member.User.AvatarURL("4096")},
				}
				s.ChannelMessageSendEmbed(m.ChannelID, embed)
			},
		},
		{
			Name:        "announce",
			Description: "Make an announcement in a channel",
			Usage:       "announce <channel> <message>",
			Category:    "Moderation",
			StaffOnly:   true,
			Run: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
				if len(args) < 2 {
					s.ChannelMessageSendEmbed(m.ChannelID, ErrorEmbed("Usage: `_announce <#channel> <message>`"))
					return
				}
				channelID := extractID(args[0])
				message := strings.Join(args[1:], " ")
				embed := &discordgo.MessageEmbed{
					Title:       "📢 Announcement",
					Description: message,
					Color:       config.EmbedColor,
					Footer:      &discordgo.MessageEmbedFooter{Text: fmt.Sprintf("Announced by %s", m.Author.Username)},
				}
				_, err := s.ChannelMessageSendEmbed(channelID, embed)
				if err != nil {
					s.ChannelMessageSendEmbed(m.ChannelID, ErrorEmbed("Failed to send announcement. Invalid channel?"))
					return
				}
				s.ChannelMessageSendEmbed(m.ChannelID, SuccessEmbed("Announcement sent!"))
			},
		},
		{
			Name:        "say",
			Description: "Make the bot say something",
			Usage:       "say <message>",
			Category:    "Utilities",
			StaffOnly:   true,
			Run: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
				if len(args) < 1 {
					s.ChannelMessageSendEmbed(m.ChannelID, ErrorEmbed("Usage: `_say <message>`"))
					return
				}
				message := strings.Join(args, " ")
				s.ChannelMessageSend(m.ChannelID, message)
			},
		},
		{
			Name:        "report",
			Description: "Report a user to staff",
			Usage:       "report <user> <reason>",
			Category:    "Utilities",
			Run: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
				if len(args) < 2 {
					s.ChannelMessageSendEmbed(m.ChannelID, ErrorEmbed("Usage: `_report <user> <reason>`"))
					return
				}
				targetID := extractID(args[0])
				reason := strings.Join(args[1:], " ")

				var reportChanID string
				channels, _ := s.GuildChannels(m.GuildID)
				for _, ch := range channels {
					if ch.Name == config.ReportsChannel {
						reportChanID = ch.ID
						break
					}
				}
				if reportChanID == "" {
					s.ChannelMessageSendEmbed(m.ChannelID, SuccessEmbed("Report submitted. Staff will review it shortly."))
					return
				}

				embed := &discordgo.MessageEmbed{
					Title:       "New Report",
					Color:       config.ErrorColor,
					Fields: []*discordgo.MessageEmbedField{
						{Name: "Reporter", Value: fmt.Sprintf("<@%s>", m.Author.ID), Inline: true},
						{Name: "Reported", Value: fmt.Sprintf("<@%s>", targetID), Inline: true},
						{Name: "Reason", Value: reason, Inline: false},
					},
					Timestamp: m.Timestamp.Format(time.RFC3339),
					Footer:    &discordgo.MessageEmbedFooter{Text: fmt.Sprintf("Reporter ID: %s", m.Author.ID)},
				}
				s.ChannelMessageSendEmbed(reportChanID, embed)
				s.ChannelMessageSendEmbed(m.ChannelID, SuccessEmbed("Report submitted. Staff will review it shortly."))
			},
		},
	}

	base := &Command{Name: "utilities", Category: "Utilities"}
	for _, cmd := range cmds {
		cm.Register(cmd)
	}
	return base
}
