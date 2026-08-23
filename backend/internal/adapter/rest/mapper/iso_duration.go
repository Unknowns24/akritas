package mapper

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"
)

var fixedISODuration = regexp.MustCompile(`^P(?:(\d+)D)?(?:T(?:(\d+)H)?(?:(\d+)M)?(?:(\d+)S)?)?$`)

func parseFixedISODuration(value string) (time.Duration, error) {
	matches := fixedISODuration.FindStringSubmatch(value)
	if matches == nil || value == "P" || value == "PT" {
		return 0, errors.New("invalid fixed ISO-8601 duration")
	}
	parts := make([]int64, 4)
	for index := range parts {
		if matches[index+1] != "" {
			parsed, err := strconv.ParseInt(matches[index+1], 10, 64)
			if err != nil {
				return 0, err
			}
			parts[index] = parsed
		}
	}
	window := time.Duration(parts[0])*24*time.Hour + time.Duration(parts[1])*time.Hour + time.Duration(parts[2])*time.Minute + time.Duration(parts[3])*time.Second
	if window <= 0 {
		return 0, errors.New("duration must be positive")
	}
	return window, nil
}

func formatFixedISODuration(value time.Duration) string {
	if value%(24*time.Hour) == 0 {
		return fmt.Sprintf("P%dD", value/(24*time.Hour))
	}
	if value%time.Hour == 0 {
		return fmt.Sprintf("PT%dH", value/time.Hour)
	}
	if value%time.Minute == 0 {
		return fmt.Sprintf("PT%dM", value/time.Minute)
	}
	return fmt.Sprintf("PT%dS", value/time.Second)
}
