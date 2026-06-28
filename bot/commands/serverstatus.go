package commands

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"geowar-bot/config"
	"github.com/bwmarrin/discordgo"
)

type MCStatusResponse struct {
	Online   bool   `json:"online"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Version  struct {
		NameRaw   string `json:"name_raw"`
		NameClean string `json:"name_clean"`
		Protocol  int    `json:"protocol"`
	} `json:"version"`
	Players struct {
		Online int      `json:"online"`
		Max    int      `json:"max"`
		List   []Player `json:"list"`
	} `json:"players"`
	MOTD struct {
		Raw   string `json:"raw"`
		Clean string `json:"clean"`
		HTML  string `json:"html"`
	} `json:"motd"`
	Icon string `json:"icon"`
}

type Player struct {
	Name string `json:"name"`
	UUID string `json:"uuid"`
}

const mcAPIURL = "https://api.mcstatus.io/v2/status/java/node.vividmcc.com:25593"

type ServerStatusCmd struct {
	lastFetch time.Time
	cached    *MCStatusResponse
}

func NewServerStatusCmd() *ServerStatusCmd {
	return &ServerStatusCmd{}
}

func NewServerStatusCommand() *Command {
	ss := NewServerStatusCmd()
	return &Command{
		Name:        "serverstatus",
		Aliases:     []string{"status", "mcstatus", "minecraft"},
		Description: "Shows the GeoWar Minecraft server status (updates every minute)",
		Usage:       "serverstatus",
		Category:    "Utilities",
		Run: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
			ss.sendStatus(s, m.ChannelID, m.Message)
		},
	}
}

func (ss *ServerStatusCmd) HandleSlash(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Fetching server status...",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})

	embed, err := ss.buildStatusEmbed()
	if err != nil {
		s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
			Content: fmt.Sprintf("Failed to fetch server status: %v", err),
			Flags:   discordgo.MessageFlagsEphemeral,
		})
		return
	}

	msg, err := s.FollowupMessageCreate(i.Interaction, true, &discordgo.WebhookParams{
		Embeds: []*discordgo.MessageEmbed{embed},
	})
	if err != nil {
		return
	}

	go ss.autoUpdateWebhook(s, i.Interaction, msg.ID)
}

func (ss *ServerStatusCmd) sendStatus(s *discordgo.Session, channelID string, replyTo *discordgo.Message) {
	embed, err := ss.buildStatusEmbed()
	if err != nil {
		s.ChannelMessageSendEmbed(channelID, ErrorEmbed(fmt.Sprintf("Failed to fetch server status: %v", err)))
		return
	}

	msg, err := s.ChannelMessageSendEmbed(channelID, embed)
	if err != nil {
		return
	}

	s.MessageReactionAdd(channelID, msg.ID, "🔄")

	go ss.autoUpdate(s, channelID, msg.ID)
}

func (ss *ServerStatusCmd) fetchStatus() (*MCStatusResponse, error) {
	if time.Since(ss.lastFetch) < 30*time.Second && ss.cached != nil {
		return ss.cached, nil
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(mcAPIURL)
	if err != nil {
		return nil, fmt.Errorf("connection failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var status MCStatusResponse
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	ss.cached = &status
	ss.lastFetch = time.Now()
	return &status, nil
}

func (ss *ServerStatusCmd) buildStatusEmbed() (*discordgo.MessageEmbed, error) {
	status, err := ss.fetchStatus()
	if err != nil {
		return nil, err
	}

	color := config.SuccessColor
	statusEmoji := "🟢"
	statusText := "Online"
	if !status.Online {
		color = config.ErrorColor
		statusEmoji = "🔴"
		statusText = "Offline"
	}

	embed := &discordgo.MessageEmbed{
		Title:     fmt.Sprintf("%s GeoWar Minecraft Server", statusEmoji),
		Color:     color,
		Timestamp: time.Now().Format(time.RFC3339),
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Updates every minute | GeoWar SMP",
		},
	}

	if !status.Online {
		embed.Description = "The server is currently **offline**."
		return embed, nil
	}

	version := status.Version.NameClean
	if version == "" {
		version = status.Version.NameRaw
	}

	embed.Description = fmt.Sprintf("**%s:%d**", status.Host, status.Port)

	playerList := "No players online"
	if status.Players.Online > 0 {
		var names []string
		for _, p := range status.Players.List {
			names = append(names, p.Name)
		}
		if len(names) > 0 {
			playerList = strings.Join(names, ", ")
			if len(playerList) > 900 {
				playerList = playerList[:900] + "..."
			}
		}
	}

	motd := status.MOTD.Clean
	if len(motd) > 200 {
		motd = motd[:200] + "..."
	}

	embed.Fields = []*discordgo.MessageEmbedField{
		{Name: "Status", Value: fmt.Sprintf("%s %s", statusEmoji, statusText), Inline: true},
		{Name: "Players", Value: fmt.Sprintf("**%d** / **%d**", status.Players.Online, status.Players.Max), Inline: true},
		{Name: "Version", Value: version, Inline: true},
		{Name: "Player List", Value: playerList, Inline: false},
		{Name: "MOTD", Value: motd, Inline: false},
	}

	if status.Icon != "" {
		embed.Thumbnail = &discordgo.MessageEmbedThumbnail{URL: status.Icon}
	}

	return embed, nil
}

func (ss *ServerStatusCmd) autoUpdate(s *discordgo.Session, channelID, messageID string) {
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		embed, err := ss.buildStatusEmbed()
		if err != nil {
			continue
		}

		_, err = s.ChannelMessageEditEmbed(channelID, messageID, embed)
		if err != nil {
			return
		}
	}
}

func (ss *ServerStatusCmd) autoUpdateWebhook(s *discordgo.Session, i *discordgo.Interaction, messageID string) {
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		embed, err := ss.buildStatusEmbed()
		if err != nil {
			continue
		}

		_, err = s.FollowupMessageEdit(i, messageID, &discordgo.WebhookEdit{
			Embeds: &[]*discordgo.MessageEmbed{embed},
		})
		if err != nil {
			return
		}
	}
}
