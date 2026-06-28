package commands

import (
	"fmt"
	"time"

	"geowar-bot/config"
	"github.com/bwmarrin/discordgo"
)

func NewTicketCommands(s *discordgo.Session, cm *CommandManager) *Command {
	cm.Register(&Command{
		Name:        "ticket",
		Aliases:     []string{"tickets"},
		Description: "Ticket system commands",
		Usage:       "ticket <setup|add|remove|close|claim|unclaim>",
		Category:    "Tickets",
		StaffOnly:   true,
		Run: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
			if len(args) < 1 {
				s.ChannelMessageSendEmbed(m.ChannelID, ErrorEmbed("Usage: `_ticket setup|add|remove|close|claim|unclaim`"))
				return
			}
			sub := args[0]
			switch sub {
			case "setup":
				setupTicketPanel(s, m)
			case "add":
				if len(args) < 2 {
					s.ChannelMessageSendEmbed(m.ChannelID, ErrorEmbed("Usage: `_ticket add <user>`"))
					return
				}
				addToTicket(s, m, args[1])
			case "remove":
				if len(args) < 2 {
					s.ChannelMessageSendEmbed(m.ChannelID, ErrorEmbed("Usage: `_ticket remove <user>`"))
					return
				}
				removeFromTicket(s, m, args[1])
			case "close":
				closeTicket(s, m)
			case "claim":
				claimTicket(s, m)
			case "unclaim":
				unclaimTicket(s, m)
			default:
				s.ChannelMessageSendEmbed(m.ChannelID, ErrorEmbed("Unknown subcommand. Use `setup|add|remove|close|claim|unclaim`"))
			}
		},
	})
	return &Command{Name: "tickets"}
}

func setupTicketPanel(s *discordgo.Session, m *discordgo.MessageCreate) {
	embed := &discordgo.MessageEmbed{
		Title:       "GeoWar Support Tickets",
		Description: "Need help from staff? Click the button below to create a ticket!\n\n**Rules:**\n- Do not create unnecessary tickets\n- Be respectful to staff\n- Use the form to describe your issue",
		Color:       config.EmbedColor,
		Thumbnail: &discordgo.MessageEmbedThumbnail{URL: s.State.User.AvatarURL("")},
		Footer:     &discordgo.MessageEmbedFooter{Text: "GeoWar SMP"},
	}

	view := &TicketPanelView{}
	msg, err := s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Embed: embed,
		Components: view.Components(),
	})
	if err != nil {
		s.ChannelMessageSendEmbed(m.ChannelID, ErrorEmbed("Failed to create ticket panel."))
		return
	}
	_, _ = msg, err
}

func addToTicket(s *discordgo.Session, m *discordgo.MessageCreate, userID string) {
	targetID := extractID(userID)
	perms := &discordgo.PermissionOverwrite{
		ID:    targetID,
		Type:  discordgo.PermissionOverwriteTypeMember,
		Allow: discordgo.PermissionViewChannel | discordgo.PermissionSendMessages | discordgo.PermissionReadMessageHistory,
	}
	_, err := s.ChannelEditComplex(m.ChannelID, &discordgo.ChannelEdit{
		PermissionOverwrites: append(getOverwrites(s, m.ChannelID), perms),
	})
	if err != nil {
		s.ChannelMessageSendEmbed(m.ChannelID, ErrorEmbed("Failed to add user."))
		return
	}
	s.ChannelMessageSendEmbed(m.ChannelID, SuccessEmbed(fmt.Sprintf("Added <@%s> to this ticket.", targetID)))
}

func removeFromTicket(s *discordgo.Session, m *discordgo.MessageCreate, userID string) {
	targetID := extractID(userID)
	perms := &discordgo.PermissionOverwrite{
		ID:    targetID,
		Type:  discordgo.PermissionOverwriteTypeMember,
		Deny:  discordgo.PermissionViewChannel,
	}
	_, err := s.ChannelEditComplex(m.ChannelID, &discordgo.ChannelEdit{
		PermissionOverwrites: append(getOverwrites(s, m.ChannelID), perms),
	})
	if err != nil {
		s.ChannelMessageSendEmbed(m.ChannelID, ErrorEmbed("Failed to remove user."))
		return
	}
	s.ChannelMessageSendEmbed(m.ChannelID, SuccessEmbed(fmt.Sprintf("Removed <@%s> from this ticket.", targetID)))
}

func closeTicket(s *discordgo.Session, m *discordgo.MessageCreate) {
	embed := &discordgo.MessageEmbed{
		Title:       "Ticket Closing",
		Description: fmt.Sprintf("This ticket will be closed in %d seconds.", config.TicketCloseDelay),
		Color:       config.WarningColor,
	}
	s.ChannelMessageSendEmbed(m.ChannelID, embed)
	time.Sleep(time.Duration(config.TicketCloseDelay) * time.Second)
	s.ChannelDelete(m.ChannelID)
}

func claimTicket(s *discordgo.Session, m *discordgo.MessageCreate) {
	embed := &discordgo.MessageEmbed{
		Title:       "Ticket Claimed",
		Description: fmt.Sprintf("<@%s> is now handling this ticket.", m.Author.ID),
		Color:       config.SuccessColor,
	}
	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}

func unclaimTicket(s *discordgo.Session, m *discordgo.MessageCreate) {
	embed := &discordgo.MessageEmbed{
		Title:       "Ticket Unclaimed",
		Description: "This ticket is no longer being handled.",
		Color:       config.WarningColor,
	}
	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}

func getOverwrites(s *discordgo.Session, channelID string) []*discordgo.PermissionOverwrite {
	ch, err := s.State.Channel(channelID)
	if err != nil {
		return nil
	}
	return ch.PermissionOverwrites
}

type TicketPanelView struct{}

func (v *TicketPanelView) Components() []discordgo.MessageComponent {
	return []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label:    "Create Ticket",
					Style:    discordgo.PrimaryButton,
					CustomID: "create_ticket",
					Emoji:    &discordgo.ComponentEmoji{Name: "🎫"},
				},
			},
		},
	}
}
