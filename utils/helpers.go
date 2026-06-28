package utils

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

func ParseDuration(text string) (time.Duration, bool) {
	re := regexp.MustCompile(`^(\d+)\s*([smhdw])$`)
	matches := re.FindStringSubmatch(strings.TrimSpace(strings.ToLower(text)))
	if len(matches) < 3 {
		return 0, false
	}

	amount, err := strconv.Atoi(matches[1])
	if err != nil || amount <= 0 {
		return 0, false
	}

	var d time.Duration
	switch matches[2] {
	case "s":
		d = time.Duration(amount) * time.Second
	case "m":
		d = time.Duration(amount) * time.Minute
	case "h":
		d = time.Duration(amount) * time.Hour
	case "d":
		d = time.Duration(amount) * 24 * time.Hour
	case "w":
		d = time.Duration(amount) * 7 * 24 * time.Hour
	default:
		return 0, false
	}
	return d, true
}

func FormatDuration(d time.Duration) string {
	secs := int(d.Seconds())
	periods := []struct {
		name string
		secs int
	}{
		{"w", 604800},
		{"d", 86400},
		{"h", 3600},
		{"m", 60},
		{"s", 1},
	}

	var parts []string
	for _, p := range periods {
		if secs >= p.secs {
			count := secs / p.secs
			parts = append(parts, strconv.Itoa(count)+p.name)
			secs %= p.secs
		}
	}
	if len(parts) == 0 {
		return "0s"
	}
	return strings.Join(parts, " ")
}

func IsStaffRole(roleName string) bool {
	staffRoles := []string{"Admin", "Mod", "Staff"}
	for _, r := range staffRoles {
		if strings.EqualFold(roleName, r) {
			return true
		}
	}
	return false
}
