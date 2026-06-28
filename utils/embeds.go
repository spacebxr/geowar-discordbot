package utils

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

func NewEmbed(title, description string, color int) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       title,
		Description: description,
		Color:       color,
		Timestamp:   time.Now().Format(time.RFC3339),
	}
}

func NewErrorEmbed(description string) *discordgo.MessageEmbed {
	return NewEmbed("Error", description, 0xed4245)
}

func NewSuccessEmbed(description string) *discordgo.MessageEmbed {
	return NewEmbed("Success", description, 0x57f287)
}

func NewWarningEmbed(description string) *discordgo.MessageEmbed {
	return NewEmbed("Warning", description, 0xfee75c)
}

func AddField(embed *discordgo.MessageEmbed, name, value string, inline bool) *discordgo.MessageEmbed {
	embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
		Name:   name,
		Value:  value,
		Inline: inline,
	})
	return embed
}

func SetFooter(embed *discordgo.MessageEmbed, text, iconURL string) *discordgo.MessageEmbed {
	embed.Footer = &discordgo.MessageEmbedFooter{
		Text:    text,
		IconURL: iconURL,
	}
	return embed
}

func SetAuthor(embed *discordgo.MessageEmbed, name, iconURL string) *discordgo.MessageEmbed {
	embed.Author = &discordgo.MessageEmbedAuthor{
		Name:    name,
		IconURL: iconURL,
	}
	return embed
}

func SetThumbnail(embed *discordgo.MessageEmbed, url string) *discordgo.MessageEmbed {
	embed.Thumbnail = &discordgo.MessageEmbedThumbnail{
		URL: url,
	}
	return embed
}

func SetImage(embed *discordgo.MessageEmbed, url string) *discordgo.MessageEmbed {
	embed.Image = &discordgo.MessageEmbedImage{
		URL: url,
	}
	return embed
}
