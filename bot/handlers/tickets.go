package handlers

import (
	"fmt"
	"strings"
	"time"

	"geowar-bot/config"
	"github.com/bwmarrin/discordgo"
)

type TicketHandler struct {
	openTickets map[string]bool
}

func NewTicketHandler(s *discordgo.Session) *TicketHandler {
	return &TicketHandler{
		openTickets: make(map[string]bool),
	}
}

func (h *TicketHandler) HandleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionMessageComponent:
		customID := i.MessageComponentData().CustomID
		switch customID {
		case "create_ticket":
			h.createTicket(s, i)
		case "close_ticket":
			h.closeTicket(s, i)
		case "claim_ticket":
			h.claimTicket(s, i)
		}
	case discordgo.InteractionModalSubmit:
		h.handleModalSubmit(s, i)
	}
}

type TicketModal struct{}

func (h *TicketHandler) createTicket(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: "ticket_modal",
			Title:    "Create Support Ticket",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    "ticket_reason",
							Label:       "Reason for ticket",
							Style:       discordgo.TextInputParagraph,
							Placeholder: "Describe your issue...",
							Required:    true,
							MaxLength:   500,
						},
					},
				},
			},
		},
	})
	if err != nil {
		fmt.Println("Error sending modal:", err)
	}
}

func (h *TicketHandler) handleModalSubmit(s *discordgo.Session, i *discordgo.InteractionCreate) {
	customID := i.ModalSubmitData().CustomID
	if customID != "ticket_modal" {
		return
	}

	reason := i.ModalSubmitData().Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput).Value
	userID := i.Member.User.ID
	guildID := i.GuildID

	var ticketChan *discordgo.Channel
	var categoryID string

	channels, _ := s.GuildChannels(guildID)
	for _, ch := range channels {
		if ch.Name == config.TicketCategoryName && ch.Type == discordgo.ChannelTypeGuildCategory {
			categoryID = ch.ID
			break
		}
	}
	if categoryID == "" {
		cat, err := s.GuildChannelCreateComplex(guildID, discordgo.GuildChannelCreateData{
			Name: config.TicketCategoryName,
			Type: discordgo.ChannelTypeGuildCategory,
		})
		if err == nil {
			categoryID = cat.ID
		}
	}

	ticketName := fmt.Sprintf("ticket-%s", strings.ToLower(strings.Split(i.Member.User.Username, "#")[0]))
	if len(ticketName) > 32 {
		ticketName = ticketName[:32]
	}

	ticketChan, err := s.GuildChannelCreateComplex(guildID, discordgo.GuildChannelCreateData{
		Name:     ticketName,
		Type:     discordgo.ChannelTypeGuildText,
		ParentID: categoryID,
		PermissionOverwrites: []*discordgo.PermissionOverwrite{
			{
				ID:    guildID,
				Type:  discordgo.PermissionOverwriteTypeRole,
				Deny:  discordgo.PermissionViewChannel,
			},
			{
				ID:    userID,
				Type:  discordgo.PermissionOverwriteTypeMember,
				Allow: discordgo.PermissionViewChannel | discordgo.PermissionSendMessages | discordgo.PermissionReadMessageHistory,
			},
		},
	})
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to create ticket channel.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	for _, roleName := range config.StaffRoleNames {
		for _, ch := range channels {
			if strings.EqualFold(ch.Name, roleName) || strings.EqualFold(ch.Name, roleName) {
							}
		}
		roleID := findRoleIDByName(s, guildID, roleName)
		if roleID != "" {
			s.ChannelPermissionSet(ticketChan.ID, roleID, discordgo.PermissionOverwriteTypeRole, discordgo.PermissionViewChannel|discordgo.PermissionSendMessages|discordgo.PermissionReadMessageHistory, 0)
		}
	}

	ticketEmbed := &discordgo.MessageEmbed{
		Title:       "New Ticket",
		Description: fmt.Sprintf("**User:** <@%s>\n**Reason:** %s", userID, reason),
		Color:       config.EmbedColor,
		Timestamp:   time.Now().Format(time.RFC3339),
	}

	closeBtn := discordgo.Button{
		Label:    "Close Ticket",
		Style:    discordgo.DangerButton,
		CustomID: "close_ticket",
		Emoji:    &discordgo.ComponentEmoji{Name: "🔒"},
	}

	claimBtn := discordgo.Button{
		Label:    "Claim Ticket",
		Style:    discordgo.SecondaryButton,
		CustomID: "claim_ticket",
		Emoji:    &discordgo.ComponentEmoji{Name: "✋"},
	}

	s.ChannelMessageSendComplex(ticketChan.ID, &discordgo.MessageSend{
		Embed: ticketEmbed,
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{closeBtn, claimBtn},
			},
		},
	})

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Ticket created in <#%s>", ticketChan.ID),
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}

func (h *TicketHandler) closeTicket(s *discordgo.Session, i *discordgo.InteractionCreate) {
	embed := &discordgo.MessageEmbed{
		Title:       "Closing Ticket",
		Description: fmt.Sprintf("This ticket will be closed in %d seconds.", config.TicketCloseDelay),
		Color:       config.WarningColor,
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})

	time.Sleep(time.Duration(config.TicketCloseDelay) * time.Second)
	s.ChannelDelete(i.ChannelID)
}

func (h *TicketHandler) claimTicket(s *discordgo.Session, i *discordgo.InteractionCreate) {
	embed := &discordgo.MessageEmbed{
		Title:       "Ticket Claimed",
		Description: fmt.Sprintf("<@%s> is now handling this ticket.", i.Member.User.ID),
		Color:       config.SuccessColor,
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
}

func findRoleIDByName(s *discordgo.Session, guildID, name string) string {
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
