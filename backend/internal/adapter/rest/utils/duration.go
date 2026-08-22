package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

var isoDuration = regexp.MustCompile(`^PT(?:(\d+)H)?(?:(\d+)M)?(?:(\d+)S)?$`)

func ParseISODuration(value string) (time.Duration, error) {
	matches := isoDuration.FindStringSubmatch(value)
	if matches == nil || value == "PT" {
		return 0, fmt.Errorf("invalid ISO-8601 duration")
	}
	hours, minutes, seconds := atoi(matches[1]), atoi(matches[2]), atoi(matches[3])
	duration := time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute + time.Duration(seconds)*time.Second
	if duration <= 0 {
		return 0, fmt.Errorf("invalid ISO-8601 duration")
	}
	return duration, nil
}

func FormatISODuration(duration time.Duration) string {
	if duration%time.Hour == 0 {
		return fmt.Sprintf("PT%dH", int(duration/time.Hour))
	}
	if duration%time.Minute == 0 {
		return fmt.Sprintf("PT%dM", int(duration/time.Minute))
	}
	return fmt.Sprintf("PT%dS", int(duration/time.Second))
}

func atoi(value string) int {
	if value == "" {
		return 0
	}
	parsed, _ := strconv.Atoi(value)
	return parsed
}
