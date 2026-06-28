package models

import "time"

type GuildConfig struct {
	GuildID          string            `json:"guild_id"`
	Prefix           string            `json:"prefix"`
	TicketCategoryID string            `json:"ticket_category_id"`
	TicketLogChanID  string            `json:"ticket_log_channel_id"`
	LogChannelID     string            `json:"log_channel_id"`
	WelcomeChannelID string            `json:"welcome_channel_id"`
	SuggestChanID    string            `json:"suggest_channel_id"`
	ReportsChanID    string            `json:"reports_channel_id"`
	RulesChannelID   string            `json:"rules_channel_id"`
	ModRoleID        string            `json:"mod_role_id"`
	AdminRoleID      string            `json:"admin_role_id"`
	StaffRoleIDs     []string          `json:"staff_role_ids"`
	WelcomeEnabled   bool              `json:"welcome_enabled"`
	WelcomeMessage   string            `json:"welcome_message"`
	LoggingEnabled   bool              `json:"logging_enabled"`
	AutoModEnabled   bool              `json:"automod_enabled"`
	BannedWords      []string          `json:"banned_words"`
	BannedDomains    []string          `json:"banned_domains"`
	TicketCounter    int               `json:"ticket_counter"`
	ReactionRoles    []ReactionRole    `json:"reaction_roles"`
	Warns            map[string][]Warn `json:"warns"`
	MutedRoles       map[string]string `json:"muted_roles"`
	Giveaways        []Giveaway        `json:"giveaways"`
}

type ReactionRole struct {
	MessageID string `json:"message_id"`
	ChannelID string `json:"channel_id"`
	Emoji     string `json:"emoji"`
	RoleID    string `json:"role_id"`
}

type Warn struct {
	Reason   string    `json:"reason"`
	ModID    string    `json:"mod_id"`
	Date     time.Time `json:"date"`
	WarnID   string    `json:"warn_id"`
}

type Ticket struct {
	TicketID  string    `json:"ticket_id"`
	ChannelID string    `json:"channel_id"`
	UserID    string    `json:"user_id"`
	Reason    string    `json:"reason"`
	Open      bool      `json:"open"`
	CreatedAt time.Time `json:"created_at"`
	ClaimedBy string    `json:"claimed_by,omitempty"`
}

type Giveaway struct {
	MessageID string    `json:"message_id"`
	ChannelID string    `json:"channel_id"`
	Prize     string    `json:"prize"`
	Winners   int       `json:"winners"`
	EndTime   time.Time `json:"end_time"`
	HostID    string    `json:"host_id"`
	Entries   []string  `json:"entries"`
	Ended     bool      `json:"ended"`
}

type Mute struct {
	UserID  string    `json:"user_id"`
	EndTime time.Time `json:"end_time"`
	Reason  string    `json:"reason"`
	ModID   string    `json:"mod_id"`
}
