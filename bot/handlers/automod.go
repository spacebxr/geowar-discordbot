package handlers

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

type AutoModHandler struct {
	spamTracker map[string]*SpamData
	mu          sync.Mutex
}

type SpamData struct {
	Count    int
	FirstMsg time.Time
	LastWarn time.Time
}

func NewAutoModHandler(s *discordgo.Session) *AutoModHandler {
	return &AutoModHandler{
		spamTracker: make(map[string]*SpamData),
	}
}

func (h *AutoModHandler) HandleMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}

	// Check for mass mentions
	if len(m.Mentions) > 5 {
		h.deleteAndWarn(s, m, "Mass mention detected")
		return
	}

	// Check for invite links
	content := strings.ToLower(m.Content)
	if strings.Contains(content, "discord.gg/") || strings.Contains(content, "discord.com/invite/") {
		// Allow if in staff channel or user is staff
		if !isStaffByRoles(s, m.GuildID, m.Author.ID) {
			h.deleteAndWarn(s, m, "Discord invite links are not allowed")
			return
		}
	}

	// Check for excessive caps
	if len(content) > 20 {
		upperCount := 0
		for _, ch := range content {
			if ch >= 'A' && ch <= 'Z' {
				upperCount++
			}
		}
		if float64(upperCount)/float64(len(content)) > 0.7 {
			h.deleteAndWarn(s, m, "Excessive caps detected")
			return
		}
	}

	// Spam detection
	h.checkSpam(s, m)
}

func (h *AutoModHandler) deleteAndWarn(s *discordgo.Session, m *discordgo.MessageCreate, reason string) {
	s.ChannelMessageDelete(m.ChannelID, m.ID)

	warnEmbed := &discordgo.MessageEmbed{
		Title:       "Auto-Mod Warning",
		Description: fmt.Sprintf("<@%s>, %s.", m.Author.ID, reason),
		Color:       0xfee75c,
	}
	s.ChannelMessageSendEmbed(m.ChannelID, warnEmbed)
}

func (h *AutoModHandler) checkSpam(s *discordgo.Session, m *discordgo.MessageCreate) {
	h.mu.Lock()
	defer h.mu.Unlock()

	key := fmt.Sprintf("%s:%s", m.GuildID, m.Author.ID)
	data, exists := h.spamTracker[key]
	now := time.Now()

	if !exists {
		h.spamTracker[key] = &SpamData{
			Count:    1,
			FirstMsg: now,
		}
		return
	}

	// Reset if more than 5 seconds have passed
	if now.Sub(data.FirstMsg) > 5*time.Second {
		data.Count = 1
		data.FirstMsg = now
		return
	}

	data.Count++

	if data.Count > 5 && now.Sub(data.LastWarn) > 10*time.Second {
		data.LastWarn = now
		h.mu.Unlock()
		s.ChannelMessageDelete(m.ChannelID, m.ID)
		embed := &discordgo.MessageEmbed{
			Title:       "Auto-Mod: Spam Detected",
			Description: fmt.Sprintf("<@%s>, please stop spamming.", m.Author.ID),
			Color:       0xed4245,
		}
		s.ChannelMessageSendEmbed(m.ChannelID, embed)
		h.mu.Lock()
	}
}

func isStaffByRoles(s *discordgo.Session, guildID, userID string) bool {
	member, err := s.GuildMember(guildID, userID)
	if err != nil {
		return false
	}
	return member.Permissions&discordgo.PermissionAdministrator != 0 ||
		member.Permissions&discordgo.PermissionModerateMembers != 0
}

func (h *AutoModHandler) AddBannedWord(guildID, word string) {
	// Would be persisted to config
}

func (h *AutoModHandler) RemoveBannedWord(guildID, word string) {
	// Would be persisted from config
}
