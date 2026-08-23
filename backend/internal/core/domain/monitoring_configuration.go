package domain

import (
	"regexp"
	"strings"
	"time"
)

const (
	DefaultGroupingWindow = 30 * time.Minute
	DefaultContextRecords = 20
	MaxContextRecords     = 1000
	MaxMonitoringPatterns = 100
	MaxPatternLength      = 500
)

type MonitoringConfiguration struct {
	Enabled         bool          `gorm:"column:monitoring_enabled"`
	ErrorPatterns   []string      `gorm:"serializer:json;type:jsonb;column:error_patterns"`
	IgnoredPatterns []string      `gorm:"serializer:json;type:jsonb;column:ignored_patterns"`
	GroupingWindow  time.Duration `gorm:"column:grouping_window_ns"`
	ContextBefore   int           `gorm:"column:context_before"`
	ContextAfter    int           `gorm:"column:context_after"`
}

func (c MonitoringConfiguration) Clone() MonitoringConfiguration {
	c.ErrorPatterns = cloneStrings(c.ErrorPatterns)
	c.IgnoredPatterns = cloneStrings(c.IgnoredPatterns)
	return c
}

func DefaultMonitoringConfiguration() MonitoringConfiguration {
	return MonitoringConfiguration{
		Enabled: false, ErrorPatterns: []string{}, IgnoredPatterns: []string{},
		GroupingWindow: DefaultGroupingWindow, ContextBefore: DefaultContextRecords, ContextAfter: DefaultContextRecords,
	}
}

func NewMonitoringConfiguration(
	enabled bool,
	errorPatterns, ignoredPatterns []string,
	groupingWindow time.Duration,
	contextBefore, contextAfter int,
) (MonitoringConfiguration, error) {
	configuration := MonitoringConfiguration{
		Enabled: enabled, ErrorPatterns: cloneStrings(errorPatterns), IgnoredPatterns: cloneStrings(ignoredPatterns),
		GroupingWindow: groupingWindow, ContextBefore: contextBefore, ContextAfter: contextAfter,
	}
	if err := configuration.Validate(); err != nil {
		return MonitoringConfiguration{}, err
	}
	return configuration, nil
}

func (c MonitoringConfiguration) Validate() error {
	if c.GroupingWindow <= 0 || c.ContextBefore < 0 || c.ContextBefore > MaxContextRecords ||
		c.ContextAfter < 0 || c.ContextAfter > MaxContextRecords || len(c.ErrorPatterns) > MaxMonitoringPatterns ||
		len(c.IgnoredPatterns) > MaxMonitoringPatterns {
		return ErrInvalidMonitoringConfiguration.Wrap(validationCause("monitoring limits"))
	}
	if !validPatterns(c.ErrorPatterns) || !validPatterns(c.IgnoredPatterns) {
		return ErrInvalidMonitoringConfiguration.Wrap(validationCause("monitoring patterns"))
	}
	return nil
}

func validPatterns(patterns []string) bool {
	for _, pattern := range patterns {
		if !nonBlank(pattern) || len(pattern) > MaxPatternLength || strings.TrimSpace(pattern) != pattern {
			return false
		}
		if _, err := regexp.Compile(pattern); err != nil {
			return false
		}
	}
	return true
}
