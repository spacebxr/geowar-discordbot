package commands

import (
	"fmt"
	"sort"

	"geowar-bot/config"
	"github.com/bwmarrin/discordgo"
)

func NewHelpCommand(cm *CommandManager) *Command {
	return &Command{
		Name:        "help",
		Aliases:     []string{"h", "commands"},
		Description: "Shows all commands or info about a specific command",
		Usage:       "help [command]",
		Category:    "Utilities",
		Run: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
			if len(args) > 0 {
				cmd := cm.Get(args[0])
				if cmd == nil {
					s.ChannelMessageSendEmbed(m.ChannelID, ErrorEmbed("Command not found."))
					return
				}
				aliases := ""
				if len(cmd.Aliases) > 0 {
					aliases = "`" + joinStrings(cmd.Aliases, "`, `") + "`"
				}
				embed := InfoEmbed(fmt.Sprintf("Help: %s%s", config.Prefix, cmd.Name), cmd.Description)
				embed.Fields = []*discordgo.MessageEmbedField{
					{Name: "Usage", Value: fmt.Sprintf("`%s%s`", config.Prefix, cmd.Usage), Inline: false},
					{Name: "Category", Value: cmd.Category, Inline: true},
					{Name: "Aliases", Value: aliases, Inline: true},
				}
				s.ChannelMessageSendEmbed(m.ChannelID, embed)
				return
			}

			categories := make(map[string][]*Command)
			for _, cmd := range cm.GetAll() {
				categories[cmd.Category] = append(categories[cmd.Category], cmd)
			}

			var catNames []string
			for cat := range categories {
				catNames = append(catNames, cat)
			}
			sort.Strings(catNames)

			embed := InfoEmbed("GeoWar Bot Commands", fmt.Sprintf("Prefix: `%s`\nUse `%shelp <command>` for more info.", config.Prefix, config.Prefix))
			for _, cat := range catNames {
				cmds := categories[cat]
				var cmdList string
				for _, c := range cmds {
					cmdList += fmt.Sprintf("`%s%s` ", config.Prefix, c.Name)
				}
				embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
					Name:   cat,
					Value:  cmdList,
					Inline: false,
				})
			}
			embed.Footer = &discordgo.MessageEmbedFooter{Text: "GeoWar SMP"}
			s.ChannelMessageSendEmbed(m.ChannelID, embed)
		},
	}
}

func joinStrings(strs []string, sep string) string {
	result := ""
	for i, s := range strs {
		if i > 0 {
			result += sep
		}
		result += s
	}
	return result
}
