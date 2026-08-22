package domain

import "strings"

type DetectionRuleCode string

const (
	DetectionRuleErrorLevel       DetectionRuleCode = "error_level"
	DetectionRuleFatalLevel       DetectionRuleCode = "fatal_level"
	DetectionRulePanic            DetectionRuleCode = "panic"
	DetectionRuleStackTrace       DetectionRuleCode = "stack_trace"
	DetectionRuleHTTP5xx          DetectionRuleCode = "http_5xx"
	DetectionRuleProcessCrash     DetectionRuleCode = "process_crash"
	DetectionRuleContainerRestart DetectionRuleCode = "container_restart"
)

func (c DetectionRuleCode) Validate() error {
	switch c {
	case DetectionRuleErrorLevel, DetectionRuleFatalLevel, DetectionRulePanic, DetectionRuleStackTrace,
		DetectionRuleHTTP5xx, DetectionRuleProcessCrash, DetectionRuleContainerRestart:
		return nil
	default:
		return ErrInvalidDetectionRule.Wrap(validationCause("detection rule code"))
	}
}

type BuiltInDetectionRule struct {
	Code        DetectionRuleCode
	DisplayName string
	Enabled     bool
}

func NewBuiltInDetectionRule(code DetectionRuleCode, displayName string) (BuiltInDetectionRule, error) {
	rule := BuiltInDetectionRule{Code: code, DisplayName: strings.TrimSpace(displayName), Enabled: true}
	if err := rule.Validate(); err != nil {
		return BuiltInDetectionRule{}, err
	}
	return rule, nil
}

func (r BuiltInDetectionRule) Validate() error {
	if r.Code.Validate() != nil || !nonBlank(r.DisplayName) || !r.Enabled {
		return ErrInvalidDetectionRule.Wrap(validationCause("built-in detection rule"))
	}
	return nil
}

func AllBuiltInDetectionRules() []BuiltInDetectionRule {
	return []BuiltInDetectionRule{
		mustBuiltInDetectionRule(DetectionRuleErrorLevel, "Error level"),
		mustBuiltInDetectionRule(DetectionRuleFatalLevel, "Fatal level"),
		mustBuiltInDetectionRule(DetectionRulePanic, "Panic"),
		mustBuiltInDetectionRule(DetectionRuleStackTrace, "Stack trace"),
		mustBuiltInDetectionRule(DetectionRuleHTTP5xx, "HTTP 5xx"),
		mustBuiltInDetectionRule(DetectionRuleProcessCrash, "Process crash"),
		mustBuiltInDetectionRule(DetectionRuleContainerRestart, "Container restart"),
	}
}

func mustBuiltInDetectionRule(code DetectionRuleCode, displayName string) BuiltInDetectionRule {
	rule, err := NewBuiltInDetectionRule(code, displayName)
	if err != nil {
		panic(err)
	}
	return rule
}
