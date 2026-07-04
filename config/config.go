package config

import (
	"os"

	"github.com/bwmarrin/discordgo"
)

var (
	Prefix = "_"

	Token = getEnv("DISCORD_TOKEN", "YOUR_BOT_TOKEN_HERE")

	EmbedColor    = 0x2b2d31
	ErrorColor    = 0xed4245
	SuccessColor  = 0x57f287
	WarningColor  = 0xfee75c

	TicketCategoryName = "Tickets"
	TicketLogChannel   = "ticket-logs"
	WelcomeChannel     = "welcome"
	LogChannel         = "server-logs"
	SuggestionsChannel = "suggestions"
	ReportsChannel     = "reports"
	RulesChannel       = "rules"
	WarDeclarationsChannel = "1514653247074471936"

	ModRoleName   = "Mod"
	AdminRoleName = "Admin"

	StaffRoleNames = []string{"Admin", "Mod", "Staff"}

	MaxTicketsPerUser = 3
	TicketCloseDelay  = 5

	Intents = discordgo.IntentsAll
)

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
