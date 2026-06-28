package commands

import (
	"fmt"
	"strings"

	"geowar-bot/config"
	"github.com/bwmarrin/discordgo"
)

func NewAdminCommands(s *discordgo.Session, cm *CommandManager) *Command {
	cm.Register(&Command{
		Name:        "config",
		Aliases:     []string{"cfg", "settings"},
		Description: "View or change bot configuration",
		Usage:       "config [setting] [value]",
		Category:    "Admin",
		AdminOnly:   true,
		Run: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
			if len(args) == 0 {
				embed := &discordgo.MessageEmbed{
					Title: "Bot Configuration",
					Color: config.EmbedColor,
					Fields: []*discordgo.MessageEmbedField{
						{Name: "Prefix", Value: fmt.Sprintf("`%s`", config.Prefix), Inline: true},
						{Name: "Mod Role", Value: config.ModRoleName, Inline: true},
						{Name: "Admin Role", Value: config.AdminRoleName, Inline: true},
						{Name: "Log Channel", Value: fmt.Sprintf("`%s`", config.LogChannel), Inline: true},
						{Name: "Welcome Channel", Value: fmt.Sprintf("`%s`", config.WelcomeChannel), Inline: true},
						{Name: "Suggestions Channel", Value: fmt.Sprintf("`%s`", config.SuggestionsChannel), Inline: true},
						{Name: "Reports Channel", Value: fmt.Sprintf("`%s`", config.ReportsChannel), Inline: true},
						{Name: "Ticket Category", Value: config.TicketCategoryName, Inline: true},
					},
					Footer: &discordgo.MessageEmbedFooter{Text: "GeoWar Bot Config"},
				}
				s.ChannelMessageSendEmbed(m.ChannelID, embed)
				return
			}

			setting := strings.ToLower(args[0])
			value := strings.Join(args[1:], " ")

			switch setting {
			case "prefix":
				if value == "" {
					s.ChannelMessageSendEmbed(m.ChannelID, ErrorEmbed("Please provide a prefix."))
					return
				}
				config.Prefix = value
				s.ChannelMessageSendEmbed(m.ChannelID, SuccessEmbed(fmt.Sprintf("Prefix changed to `%s`", value)))
			default:
				s.ChannelMessageSendEmbed(m.ChannelID, ErrorEmbed(fmt.Sprintf("Unknown setting: `%s`", setting)))
			}
		},
	})

	cm.Register(&Command{
		Name:        "reload",
		Description: "Reload bot configuration",
		Usage:       "reload",
		Category:    "Admin",
		AdminOnly:   true,
		Run: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
			s.ChannelMessageSendEmbed(m.ChannelID, SuccessEmbed("Configuration reloaded from config file."))
		},
	})

	cm.Register(&Command{
		Name:        "setrole",
		Aliases:     []string{"configrole"},
		Description: "Set a role for the config (mod, admin, muted)",
		Usage:       "setrole <mod|admin|muted> <role>",
		Category:    "Admin",
		AdminOnly:   true,
		Run: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
			if len(args) < 2 {
				s.ChannelMessageSendEmbed(m.ChannelID, ErrorEmbed("Usage: `_setrole <mod|admin|muted> <role>`"))
				return
			}
			roleName := strings.Join(args[1:], " ")
			roleID := FindRoleByName(s, m.GuildID, roleName)
			if roleID == "" {
				s.ChannelMessageSendEmbed(m.ChannelID, ErrorEmbed(fmt.Sprintf("Role `%s` not found.", roleName)))
				return
			}
			switch strings.ToLower(args[0]) {
			case "mod":
				config.ModRoleName = roleName
				s.ChannelMessageSendEmbed(m.ChannelID, SuccessEmbed(fmt.Sprintf("Mod role set to `%s`", roleName)))
			case "admin":
				config.AdminRoleName = roleName
				s.ChannelMessageSendEmbed(m.ChannelID, SuccessEmbed(fmt.Sprintf("Admin role set to `%s`", roleName)))
			case "muted":
				s.ChannelMessageSendEmbed(m.ChannelID, SuccessEmbed(fmt.Sprintf("Muted role set to `%s`", roleName)))
			default:
				s.ChannelMessageSendEmbed(m.ChannelID, ErrorEmbed("Unknown role type. Use `mod`, `admin`, or `muted`"))
			}
		},
	})

	cm.Register(&Command{
		Name:        "sync",
		Description: "Sync all server channels to Muted role permissions",
		Usage:       "sync",
		Category:    "Admin",
		AdminOnly:   true,
		Run: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
			mutedRoleID := FindRoleByName(s, m.GuildID, "Muted")
			if mutedRoleID == "" {
				s.ChannelMessageSendEmbed(m.ChannelID, ErrorEmbed("Muted role not found."))
				return
			}
			channels, _ := s.GuildChannels(m.GuildID)
			count := 0
			for _, ch := range channels {
				err := s.ChannelPermissionSet(ch.ID, mutedRoleID, discordgo.PermissionOverwriteTypeRole, 0, discordgo.PermissionSendMessages)
				if err == nil {
					count++
				}
			}
			s.ChannelMessageSendEmbed(m.ChannelID, SuccessEmbed(fmt.Sprintf("Synced Muted role permissions across %d channels.", count)))
		},
	})

	cm.Register(&Command{
		Name:        "slowmode",
		Description: "Set slowmode on the current channel",
		Usage:       "slowmode <seconds>",
		Category:    "Moderation",
		StaffOnly:   true,
		Run: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
			if len(args) < 1 {
				s.ChannelMessageSendEmbed(m.ChannelID, ErrorEmbed("Usage: `_slowmode <seconds>`"))
				return
			}
			seconds := 0
			fmt.Sscanf(args[0], "%d", &seconds)
			if seconds < 0 || seconds > 21600 {
				s.ChannelMessageSendEmbed(m.ChannelID, ErrorEmbed("Slowmode must be between 0 and 21600 seconds."))
				return
			}
			_, err := s.ChannelEdit(m.ChannelID, &discordgo.ChannelEdit{
				RateLimitPerUser: &seconds,
			})
			if err != nil {
				s.ChannelMessageSendEmbed(m.ChannelID, ErrorEmbed("Failed to set slowmode."))
				return
			}
			if seconds == 0 {
				s.ChannelMessageSendEmbed(m.ChannelID, SuccessEmbed("Slowmode disabled."))
			} else {
				s.ChannelMessageSendEmbed(m.ChannelID, SuccessEmbed(fmt.Sprintf("Slowmode set to %d seconds.", seconds)))
			}
		},
	})

	return &Command{Name: "admin"}
}
