package handlers

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type ReactionRoleHandler struct {
	reactionRoles []ReactionRoleConfig
}

type ReactionRoleConfig struct {
	MessageID string
	ChannelID string
	Emoji     string
	RoleID    string
	GuildID   string
}

func NewReactionRoleHandler(s *discordgo.Session) *ReactionRoleHandler {
	return &ReactionRoleHandler{}
}

func (h *ReactionRoleHandler) HandleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionMessageComponent {
		return
	}

	customID := i.MessageComponentData().CustomID
	if !strings.HasPrefix(customID, "rr_") {
		return
	}

	parts := strings.SplitN(customID, "_", 3)
	if len(parts) < 3 {
		return
	}

	roleID := parts[2]
	member := i.Member
	if member == nil {
		return
	}

	hasRole := false
	for _, r := range member.Roles {
		if r == roleID {
			hasRole = true
			break
		}
	}

	var response string
	if hasRole {
		err := s.GuildMemberRoleRemove(i.GuildID, member.User.ID, roleID)
		if err != nil {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{Content: "Failed to remove role.", Flags: discordgo.MessageFlagsEphemeral},
			})
			return
		}
		response = "Role removed!"
	} else {
		err := s.GuildMemberRoleAdd(i.GuildID, member.User.ID, roleID)
		if err != nil {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{Content: "Failed to add role.", Flags: discordgo.MessageFlagsEphemeral},
			})
			return
		}
		response = "Role added!"
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: response,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}

func (h *ReactionRoleHandler) HandleReactionAdd(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	for _, rr := range h.reactionRoles {
		if rr.MessageID == r.MessageID && rr.ChannelID == r.ChannelID && rr.Emoji == r.Emoji.Name {
			s.GuildMemberRoleAdd(r.GuildID, r.UserID, rr.RoleID)
			return
		}
	}
}

func (h *ReactionRoleHandler) AddReactionRole(guildID, channelID, messageID, emoji, roleID string) {
	h.reactionRoles = append(h.reactionRoles, ReactionRoleConfig{
		MessageID: messageID,
		ChannelID: channelID,
		Emoji:     emoji,
		RoleID:    roleID,
		GuildID:   guildID,
	})
}

func (h *ReactionRoleHandler) RemoveReactionRole(messageID, emoji string) {
	var filtered []ReactionRoleConfig
	for _, rr := range h.reactionRoles {
		if rr.MessageID == messageID && rr.Emoji == emoji {
			continue
		}
		filtered = append(filtered, rr)
	}
	h.reactionRoles = filtered
}

func (h *ReactionRoleHandler) CreateRolePanel(s *discordgo.Session, channelID string, roles []RolePanelOption) (*discordgo.Message, error) {
	var description string
	var components []discordgo.MessageComponent
	var row discordgo.ActionsRow

	for i, role := range roles {
		description += fmt.Sprintf("%s - <@&%s>\n", role.Emoji, role.RoleID)
		row.Components = append(row.Components, discordgo.Button{
			Label:    role.Label,
			Style:    discordgo.SecondaryButton,
			CustomID: fmt.Sprintf("rr_role_%s", role.RoleID),
			Emoji:    &discordgo.ComponentEmoji{Name: role.Emoji},
		})

		if (i+1)%5 == 0 || i == len(roles)-1 {
			components = append(components, row)
			row = discordgo.ActionsRow{}
		}
	}

	embed := &discordgo.MessageEmbed{
		Title:       "Reaction Roles",
		Description: "Click a button to get the corresponding role!\n\n" + description,
		Color:       0x2b2d31,
	}

	return s.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
		Embed:      embed,
		Components: components,
	})
}

type RolePanelOption struct {
	Label  string
	Emoji  string
	RoleID string
}
