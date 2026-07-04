package commands

import (
	"fmt"
	"time"

	"geowar-bot/config"
	"github.com/bwmarrin/discordgo"
)

type WarCmd struct{}

func NewWarCmd() *WarCmd {
	return &WarCmd{}
}

func NewWarCommand() *Command {
	w := NewWarCmd()
	return &Command{
		Name:        "declarewar",
		Aliases:     []string{"warrequest", "war"},
		Description: "Submit a war declaration request",
		Usage:       "declarewar",
		Category:    "Utilities",
		Run: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
			w.SendPrompt(s, m.ChannelID)
		},
	}
}

func (w *WarCmd) SendPrompt(s *discordgo.Session, channelID string) {
	embed := &discordgo.MessageEmbed{
		Title:       "⚔️ War Declaration Request",
		Description: "Click the button below to submit a war declaration request.",
		Color:       0xed4245,
	}
	btn := discordgo.Button{
		Label:    "Declare War",
		Style:    discordgo.DangerButton,
		CustomID: "declarewar_prompt",
		Emoji:    &discordgo.ComponentEmoji{Name: "⚔️"},
	}
	s.ChannelMessageSendComplex(channelID, &discordgo.MessageSend{
		Embed: embed,
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{btn},
			},
		},
	})
}

func (w *WarCmd) HandleSlash(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: "war_modal",
			Title:    "Declare War",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    "war_declaring_nation",
							Label:       "Your Nation",
							Style:       discordgo.TextInputShort,
							Placeholder: "Enter your nation name",
							Required:    true,
							MaxLength:   100,
						},
					},
				},
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    "war_enemy_nation",
							Label:       "Enemy Nation",
							Style:       discordgo.TextInputShort,
							Placeholder: "Enter the enemy nation name",
							Required:    true,
							MaxLength:   100,
						},
					},
				},
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    "war_reason",
							Label:       "Reason For War",
							Style:       discordgo.TextInputParagraph,
							Placeholder: "Must be reasonable. Staff can deny weak reasons.",
							Required:    true,
							MaxLength:   1000,
						},
					},
				},
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    "war_demands",
							Label:       "What Do You Want From Them?",
							Style:       discordgo.TextInputParagraph,
							Placeholder: "Must be reasonable. Staff can deny/change demands.",
							Required:    true,
							MaxLength:   1000,
						},
					},
				},
			},
		},
	})
	if err != nil {
		fmt.Println("Error sending war modal:", err)
	}
}

func (w *WarCmd) HandleModalSubmit(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var declaringNation, enemyNation, reason, demands string
	for _, row := range i.ModalSubmitData().Components {
		actionRow, ok := row.(*discordgo.ActionsRow)
		if !ok {
			continue
		}
		for _, comp := range actionRow.Components {
			textInput, ok := comp.(*discordgo.TextInput)
			if !ok {
				continue
			}
			switch textInput.CustomID {
			case "war_declaring_nation":
				declaringNation = textInput.Value
			case "war_enemy_nation":
				enemyNation = textInput.Value
			case "war_reason":
				reason = textInput.Value
			case "war_demands":
				demands = textInput.Value
			}
		}
	}

	desc := fmt.Sprintf(
		"## 🏳️ Nation Declaring War\n**%s**\n\n"+
			"## 🎯 Nation Being Declared On\n**%s**\n\n"+
			"## 📜 Reason For War\n%s\n\n"+
			"> Must be reasonable. Staff can deny weak or invalid reasons.\n\n"+
			"## 💰 What Do You Want From Them?\n%s\n\n"+
			"> Must be reasonable. Staff can approve, change, or deny the demands.\n\n"+
			"---\n\n"+
			"## ℹ️ Important Info\n"+
			"- Both nations must agree on a war time before the war starts.\n"+
			"- Staff will decide if the reason is valid.\n"+
			"- Staff will decide if the demands are fair.\n"+
			"- Fake reasons, troll wars, or unfair demands can be denied.",
		declaringNation,
		enemyNation,
		reason,
		demands,
	)

	embed := &discordgo.MessageEmbed{
		Title:       "⚔️ WAR DECLARATION REQUEST",
		Description: desc,
		Color:       0xed4245,
		Timestamp:   time.Now().Format(time.RFC3339),
		Author: &discordgo.MessageEmbedAuthor{
			Name:    fmt.Sprintf("Submitted by %s", i.Member.User.Username),
			IconURL: i.Member.User.AvatarURL(""),
		},
	}

	targetChanID := config.WarDeclarationsChannel
	_, err := s.ChannelMessageSendEmbed(targetChanID, embed)
	if err != nil {
		targetChanID = i.ChannelID
		_, err = s.ChannelMessageSendEmbed(targetChanID, embed)
		if err != nil {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Failed to post war declaration request.",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			return
		}
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Your war declaration request has been submitted in <#%s>!", targetChanID),
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}
