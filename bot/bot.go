package bot

import (
	"geowar-bot/bot/commands"
	"geowar-bot/bot/handlers"
	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	Session      *discordgo.Session
	Commands     *commands.CommandManager
	Tickets      *handlers.TicketHandler
	Logging      *handlers.LoggingHandler
	Welcome      *handlers.WelcomeHandler
	Reaction     *handlers.ReactionRoleHandler
	ServerStatus *commands.ServerStatusCmd
}

func New(s *discordgo.Session) *Bot {
	b := &Bot{
		Session:      s,
		Commands:     commands.New(s),
		Tickets:      handlers.NewTicketHandler(s),
		Logging:      handlers.NewLoggingHandler(s),
		Welcome:      handlers.NewWelcomeHandler(s),
		Reaction:     handlers.NewReactionRoleHandler(s),
		ServerStatus: commands.NewServerStatusCmd(),
	}
	return b
}

func (b *Bot) MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot || m.Author.ID == s.State.User.ID {
		return
	}
	b.Commands.HandleMessageCreate(s, m)
}

func (b *Bot) MessageDelete(s *discordgo.Session, m *discordgo.MessageDelete) {
	b.Logging.HandleMessageDelete(s, m)
}

func (b *Bot) MessageUpdate(s *discordgo.Session, m *discordgo.MessageUpdate) {
	b.Logging.HandleMessageUpdate(s, m)
}

func (b *Bot) GuildMemberAdd(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
	b.Welcome.HandleMemberAdd(s, m)
	b.Logging.HandleMemberAdd(s, m)
}

func (b *Bot) GuildMemberRemove(s *discordgo.Session, m *discordgo.GuildMemberRemove) {
	b.Logging.HandleMemberRemove(s, m)
}

func (b *Bot) VoiceStateUpdate(s *discordgo.Session, v *discordgo.VoiceStateUpdate) {
	b.Logging.HandleVoiceStateUpdate(s, v)
}

func (b *Bot) InteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionApplicationCommand {
		name := i.ApplicationCommandData().Name
		switch name {
		case "serverstatus":
			b.ServerStatus.HandleSlash(s, i)
		}
		return
	}
	b.Tickets.HandleInteraction(s, i)
	b.Reaction.HandleInteraction(s, i)
}

func (b *Bot) MessageReactionAdd(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	if r.UserID == s.State.User.ID {
		return
	}
	b.Reaction.HandleReactionAdd(s, r)
}
