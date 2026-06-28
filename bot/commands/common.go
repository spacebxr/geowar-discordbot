package commands

import (
	"strings"

	"geowar-bot/config"
	"github.com/bwmarrin/discordgo"
)

func IsStaff(s *discordgo.Session, guildID, userID string) bool {
	member, err := s.GuildMember(guildID, userID)
	if err != nil {
		return false
	}
	for _, roleID := range member.Roles {
		role, err := s.State.Role(guildID, roleID)
		if err != nil {
			continue
		}
		for _, staffName := range config.StaffRoleNames {
			if strings.EqualFold(role.Name, staffName) {
				return true
			}
		}
	}
	return false
}

func IsAdmin(s *discordgo.Session, guildID, userID string) bool {
	member, err := s.GuildMember(guildID, userID)
	if err != nil {
		return false
	}
	for _, roleID := range member.Roles {
		role, err := s.State.Role(guildID, roleID)
		if err != nil {
			continue
		}
		if strings.EqualFold(role.Name, config.AdminRoleName) {
			return true
		}
	}
	return userID == guildID
}

func HasRoleByName(s *discordgo.Session, guildID, userID, roleName string) bool {
	member, err := s.GuildMember(guildID, userID)
	if err != nil {
		return false
	}
	for _, roleID := range member.Roles {
		role, err := s.State.Role(guildID, roleID)
		if err != nil {
			continue
		}
		if strings.EqualFold(role.Name, roleName) {
			return true
		}
	}
	return false
}

func FindRoleByName(s *discordgo.Session, guildID, name string) string {
	roles, err := s.GuildRoles(guildID)
	if err != nil {
		return ""
	}
	for _, role := range roles {
		if strings.EqualFold(role.Name, name) {
			return role.ID
		}
	}
	return ""
}

func ErrorEmbed(desc string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       "Error",
		Description: desc,
		Color:       config.ErrorColor,
	}
}

func SuccessEmbed(desc string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       "Success",
		Description: desc,
		Color:       config.SuccessColor,
	}
}

func InfoEmbed(title, desc string) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       title,
		Description: desc,
		Color:       config.EmbedColor,
	}
}
