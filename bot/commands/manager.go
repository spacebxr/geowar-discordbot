package commands

import (
	"strings"

	"geowar-bot/config"
	"github.com/bwmarrin/discordgo"
)

type Command struct {
	Name        string
	Aliases     []string
	Description string
	Usage       string
	Category    string
	StaffOnly   bool
	AdminOnly   bool
	Run         func(s *discordgo.Session, m *discordgo.MessageCreate, args []string)
}

type CommandManager struct {
	commands map[string]*Command
}

func New(s *discordgo.Session) *CommandManager {
	cm := &CommandManager{
		commands: make(map[string]*Command),
	}
	cm.Register(NewHelpCommand(cm))
	cm.Register(NewModerationCommands(s, cm))
	cm.Register(NewTicketCommands(s, cm))
	cm.Register(NewSuggestionCommands(s))
	cm.Register(NewPollCommands(s))
	cm.Register(NewUtilityCommands(s, cm))
	cm.Register(NewGiveawayCommands(s))
	cm.Register(NewAdminCommands(s, cm))
	cm.Register(NewServerStatusCommand())
	cm.Register(NewWarCommand())
	return cm
}

func (cm *CommandManager) Register(cmd *Command) {
	cm.commands[cmd.Name] = cmd
}

func (cm *CommandManager) Get(name string) *Command {
	if cmd, ok := cm.commands[name]; ok {
		return cmd
	}
	for _, cmd := range cm.commands {
		for _, alias := range cmd.Aliases {
			if strings.EqualFold(alias, name) {
				return cmd
			}
		}
	}
	return nil
}

func (cm *CommandManager) GetAll() map[string]*Command {
	return cm.commands
}

func (cm *CommandManager) HandleMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	prefix := config.Prefix
	content := strings.TrimSpace(m.Content)

	if !strings.HasPrefix(content, prefix) {
		return
	}

	args := strings.Fields(content[len(prefix):])
	if len(args) == 0 {
		return
	}

	cmdName := strings.ToLower(args[0])
	cmd := cm.Get(cmdName)
	if cmd == nil {
		return
	}

	if cmd.StaffOnly {
		if !IsStaff(s, m.GuildID, m.Author.ID) {
			s.ChannelMessageSendEmbed(m.ChannelID, ErrorEmbed("This command is staff-only."))
			return
		}
	}

	if cmd.AdminOnly {
		if !IsAdmin(s, m.GuildID, m.Author.ID) {
			s.ChannelMessageSendEmbed(m.ChannelID, ErrorEmbed("This command is admin-only."))
			return
		}
	}

	cmd.Run(s, m, args[1:])
}
