package commands

import (
	"fmt"
	"strings"
	"time"

	"geowar-bot/config"
	"geowar-bot/utils"
	"github.com/bwmarrin/discordgo"
)

func NewModerationCommands(s *discordgo.Session, cm *CommandManager) *Command {
	base := &Command{
		Name:     "mod",
		Aliases:  []string{"moderation"},
		Category: "Moderation",
		Run: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {},
	}

	cm.Register(&Command{
		Name:        "kick",
		Description: "Kick a member from the server",
		Usage:       "kick <user> [reason]",
		Category:    "Moderation",
		StaffOnly:   true,
		Run: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
			if len(args) < 1 {
				s.ChannelMessageSendEmbed(m.ChannelID, ErrorEmbed("Usage: `_kick <user> [reason]`"))
				return
			}
			targetID := extractID(args[0])
			reason := "No reason provided"
			if len(args) > 1 {
				reason = strings.Join(args[1:], " ")
			}
			if err := s.GuildMemberDeleteWithReason(m.GuildID, targetID, reason); err != nil {
				s.ChannelMessageSendEmbed(m.ChannelID, ErrorEmbed(fmt.Sprintf("Failed to kick: %v", err)))
				return
			}
			s.ChannelMessageSendEmbed(m.ChannelID, SuccessEmbed(fmt.Sprintf("<@%s> has been kicked.\nReason: %s", targetID, reason)))
			logAction(s, m.GuildID, "Kick", fmt.Sprintf("<@%s> kicked <@%s>", m.Author.ID, targetID), reason)
		},
	})

	cm.Register(&Command{
		Name:        "ban",
		Description: "Ban a member from the server",
		Usage:       "ban <user> [reason]",
		Category:    "Moderation",
		StaffOnly:   true,
		Run: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
			if len(args) < 1 {
				s.ChannelMessageSendEmbed(m.ChannelID, ErrorEmbed("Usage: `_ban <user> [reason]`"))
				return
			}
			targetID := extractID(args[0])
			reason := "No reason provided"
			if len(args) > 1 {
				reason = strings.Join(args[1:], " ")
			}
			if err := s.GuildBanCreateWithReason(m.GuildID, targetID, reason, 0); err != nil {
				s.ChannelMessageSendEmbed(m.ChannelID, ErrorEmbed(fmt.Sprintf("Failed to ban: %v", err)))
				return
			}
			s.ChannelMessageSendEmbed(m.ChannelID, SuccessEmbed(fmt.Sprintf("<@%s> has been banned.\nReason: %s", targetID, reason)))
			logAction(s, m.GuildID, "Ban", fmt.Sprintf("<@%s> banned <@%s>", m.Author.ID, targetID), reason)
		},
	})

	cm.Register(&Command{
		Name:        "unban",
		Description: "Unban a user by ID",
		Usage:       "unban <user_id>",
		Category:    "Moderation",
		StaffOnly:   true,
		Run: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
			if len(args) < 1 {
				s.ChannelMessageSendEmbed(m.ChannelID, ErrorEmbed("Usage: `_unban <user_id>`"))
				return
			}
			if err := s.GuildBanDelete(m.GuildID, args[0]); err != nil {
				s.ChannelMessageSendEmbed(m.ChannelID, ErrorEmbed(fmt.Sprintf("Failed to unban: %v", err)))
				return
			}
			s.ChannelMessageSendEmbed(m.ChannelID, SuccessEmbed(fmt.Sprintf("Unbanned <@%s>", args[0])))
		},
	})

	cm.Register(&Command{
		Name:        "mute",
		Description: "Mute a member (give Muted role)",
		Usage:       "mute <user> [duration] [reason]",
		Category:    "Moderation",
		StaffOnly:   true,
		Run: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
			if len(args) < 1 {
				s.ChannelMessageSendEmbed(m.ChannelID, ErrorEmbed("Usage: `_mute <user> [duration] [reason]`"))
				return
			}
			targetID := extractID(args[0])
			duration := time.Duration(0)
			reasonStart := 1
			if len(args) > 1 {
				if d, ok := utils.ParseDuration(args[1]); ok {
					duration = d
					reasonStart = 2
				}
			}
			reason := "No reason provided"
			if len(args) > reasonStart {
				reason = strings.Join(args[reasonStart:], " ")
			}

			mutedRoleID := FindRoleByName(s, m.GuildID, "Muted")
			if mutedRoleID == "" {
				mutedRole, err := s.GuildRoleCreate(m.GuildID, &discordgo.RoleParams{
					Name:        "Muted",
					Color:       &i,
					Permissions: &zeroPerms,
				})
				if err != nil {
					s.ChannelMessageSendEmbed(m.ChannelID, ErrorEmbed("Failed to create Muted role. Create it manually."))
					return
				}
				mutedRoleID = mutedRole.ID
				for _, ch := range getAllChannels(s, m.GuildID) {
					s.ChannelPermissionSet(ch.ID, mutedRoleID, discordgo.PermissionOverwriteTypeRole, 0, discordgo.PermissionSendMessages)
				}
			}

			if err := s.GuildMemberRoleAdd(m.GuildID, targetID, mutedRoleID); err != nil {
				s.ChannelMessageSendEmbed(m.ChannelID, ErrorEmbed(fmt.Sprintf("Failed to mute: %v", err)))
				return
			}

			desc := fmt.Sprintf("<@%s> has been muted.", targetID)
			if duration > 0 {
				desc += fmt.Sprintf("\nDuration: %s", utils.FormatDuration(duration))
				go func() {
					time.Sleep(duration)
					s.GuildMemberRoleRemove(m.GuildID, targetID, mutedRoleID)
				}()
			}
			desc += fmt.Sprintf("\nReason: %s", reason)
			s.ChannelMessageSendEmbed(m.ChannelID, SuccessEmbed(desc))
			logAction(s, m.GuildID, "Mute", fmt.Sprintf("<@%s> muted <@%s>", m.Author.ID, targetID), reason)
		},
	})

	cm.Register(&Command{
		Name:        "unmute",
		Description: "Unmute a member",
		Usage:       "unmute <user>",
		Category:    "Moderation",
		StaffOnly:   true,
		Run: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
			if len(args) < 1 {
				s.ChannelMessageSendEmbed(m.ChannelID, ErrorEmbed("Usage: `_unmute <user>`"))
				return
			}
			targetID := extractID(args[0])
			mutedRoleID := FindRoleByName(s, m.GuildID, "Muted")
			if mutedRoleID == "" {
				s.ChannelMessageSendEmbed(m.ChannelID, ErrorEmbed("Muted role not found."))
				return
			}
			if err := s.GuildMemberRoleRemove(m.GuildID, targetID, mutedRoleID); err != nil {
				s.ChannelMessageSendEmbed(m.ChannelID, ErrorEmbed(fmt.Sprintf("Failed to unmute: %v", err)))
				return
			}
			s.ChannelMessageSendEmbed(m.ChannelID, SuccessEmbed(fmt.Sprintf("<@%s> has been unmuted.", targetID)))
		},
	})

	cm.Register(&Command{
		Name:        "warn",
		Description: "Warn a member",
		Usage:       "warn <user> [reason]",
		Category:    "Moderation",
		StaffOnly:   true,
		Run: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
			if len(args) < 1 {
				s.ChannelMessageSendEmbed(m.ChannelID, ErrorEmbed("Usage: `_warn <user> [reason]`"))
				return
			}
			targetID := extractID(args[0])
			reason := "No reason provided"
			if len(args) > 1 {
				reason = strings.Join(args[1:], " ")
			}
			embed := SuccessEmbed(fmt.Sprintf("<@%s> has been warned.\nReason: %s", targetID, reason))
			s.ChannelMessageSendEmbed(m.ChannelID, embed)
			logAction(s, m.GuildID, "Warn", fmt.Sprintf("<@%s> warned <@%s>", m.Author.ID, targetID), reason)
		},
	})

	cm.Register(&Command{
		Name:        "purge",
		Description: "Bulk delete messages",
		Usage:       "purge <amount>",
		Category:    "Moderation",
		StaffOnly:   true,
		Run: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
			if len(args) < 1 {
				s.ChannelMessageSendEmbed(m.ChannelID, ErrorEmbed("Usage: `_purge <amount>`"))
				return
			}
			amount := 0
			fmt.Sscanf(args[0], "%d", &amount)
			if amount < 1 || amount > 100 {
				s.ChannelMessageSendEmbed(m.ChannelID, ErrorEmbed("Amount must be between 1 and 100."))
				return
			}
			msgs, err := s.ChannelMessages(m.ChannelID, amount, "", "", "")
			if err != nil {
				s.ChannelMessageSendEmbed(m.ChannelID, ErrorEmbed("Failed to fetch messages."))
				return
			}
			var ids []string
			for _, msg := range msgs {
				ids = append(ids, msg.ID)
			}
			if err := s.ChannelMessagesBulkDelete(m.ChannelID, ids); err != nil {
				s.ChannelMessageSendEmbed(m.ChannelID, ErrorEmbed("Failed to delete messages."))
				return
			}
			confirm, _ := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Deleted %d messages.", len(ids)))
			time.Sleep(3 * time.Second)
			s.ChannelMessageDelete(m.ChannelID, confirm.ID)
		},
	})

	return base
}

func extractID(input string) string {
	if len(input) > 3 && input[:2] == "<@" {
		id := strings.TrimRight(input, ">")
		id = strings.TrimLeft(id, "<@!")
		id = strings.TrimLeft(id, "<@")
		return id
	}
	return input
}

var i = 0
var zeroPerms int64 = 0

func getAllChannels(s *discordgo.Session, guildID string) []*discordgo.Channel {
	channels, err := s.GuildChannels(guildID)
	if err != nil {
		return nil
	}
	return channels
}

func logAction(s *discordgo.Session, guildID, action, description, reason string) {
	embed := &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("Mod Log: %s", action),
		Description: description,
		Color:       config.EmbedColor,
		Fields: []*discordgo.MessageEmbedField{
			{Name: "Reason", Value: reason, Inline: false},
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}
	channels, _ := s.GuildChannels(guildID)
	for _, ch := range channels {
		if ch.Name == config.LogChannel {
			s.ChannelMessageSendEmbed(ch.ID, embed)
			return
		}
	}
}
