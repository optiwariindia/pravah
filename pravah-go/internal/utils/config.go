package utils

import (
	"os"
	"strconv"
	"time"
)

const DefaultRetentionHours = 24

func EventRetentionHours() int {
	raw := os.Getenv("EVENT_RETENTION_HOURS")
	hours, err := strconv.Atoi(raw)
	if raw == "" || err != nil || hours < 1 || hours > 8760 {
		return DefaultRetentionHours
	}
	return hours
}

func EventExpiration(createdAt time.Time, retentionHours int) time.Time {
	return createdAt.Add(time.Duration(retentionHours) * time.Hour)
}
