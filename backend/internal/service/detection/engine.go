package detection

import (
	"regexp"
	"strings"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/google/uuid"
)

type Engine struct {
	ignored []*regexp.Regexp
	custom  []*regexp.Regexp
}

type Result struct {
	Timestamp         time.Time
	Message           string
	NormalizedMessage string
	Fingerprint       string
	PrimaryFamily     string
	Severity          domain.Severity
	Rules             []string
}

func NewEngine(configuration domain.MonitoringConfiguration) (*Engine, error) {
	if err := configuration.Validate(); err != nil {
		return nil, err
	}
	engine := &Engine{}
	for _, value := range configuration.IgnoredPatterns {
		compiled, err := regexp.Compile(value)
		if err != nil {
			return nil, err
		}
		engine.ignored = append(engine.ignored, compiled)
	}
	for _, value := range configuration.ErrorPatterns {
		compiled, err := regexp.Compile(value)
		if err != nil {
			return nil, err
		}
		engine.custom = append(engine.custom, compiled)
	}
	return engine, nil
}

type ruleDefinition struct {
	code     string
	severity domain.Severity
	match    func(string) bool
}

var (
	errorLevelPattern = regexp.MustCompile(`(?i)(?:^|[\s\[{,(])(?:ERROR|ERR)(?:$|[\s\]}:,)])|["']?level["']?\s*[=:]\s*["']?(?:error|err)["']?`)
	fatalPattern      = regexp.MustCompile(`(?i)(?:^|[\s\[{,(])FATAL(?:$|[\s\]}:,)])|["']?level["']?\s*[=:]\s*["']?fatal["']?`)
	panicPattern      = regexp.MustCompile(`(?i)(?:^|\s)panic(?:ed)?(?:\s|:|$)`)
	http5xxPattern    = regexp.MustCompile(`(?i)(?:\bstatus(?:_code)?\s*[=:]\s*5\d\d\b|\"\s+5\d\d\s+\d+\b|\b(?:GET|POST|PUT|PATCH|DELETE|HEAD|OPTIONS)\s+\S+\s+5\d\d\b)`)
	crashPattern      = regexp.MustCompile(`(?i)(?:segmentation fault|segfault|core dumped|out of memory|oom(?:killed)?|process (?:exited|exit) (?:with )?(?:status|code) [1-9]\d*)`)
	restartPattern    = regexp.MustCompile(`(?i)(?:restart(?:ed|ing)? container|container restart|restart[- ]loop|crashloopbackoff)`)
)

func builtInRules() []ruleDefinition {
	return []ruleDefinition{
		{string(domain.DetectionRuleErrorLevel), domain.SeverityError, errorLevelPattern.MatchString},
		{string(domain.DetectionRuleFatalLevel), domain.SeverityCritical, fatalPattern.MatchString},
		{string(domain.DetectionRulePanic), domain.SeverityCritical, panicPattern.MatchString},
		{string(domain.DetectionRuleStackTrace), domain.SeverityError, func(value string) bool { return exceptionPattern.MatchString(value) && framePattern.MatchString(value) }},
		{string(domain.DetectionRuleHTTP5xx), domain.SeverityError, http5xxPattern.MatchString},
		{string(domain.DetectionRuleProcessCrash), domain.SeverityCritical, crashPattern.MatchString},
		{string(domain.DetectionRuleContainerRestart), domain.SeverityCritical, restartPattern.MatchString},
	}
}

func (e *Engine) Detect(projectID uuid.UUID, event LogicalEvent) *Result {
	message := event.Message()
	if message == "" {
		return nil
	}
	for _, ignored := range e.ignored {
		if ignored.MatchString(message) {
			return nil
		}
	}
	rules := make([]string, 0, len(builtInRules())+len(e.custom))
	severity := domain.SeverityError
	for _, rule := range builtInRules() {
		if rule.match(message) {
			rules = append(rules, rule.code)
			if rule.severity == domain.SeverityCritical {
				severity = domain.SeverityCritical
			}
		}
	}
	for index, custom := range e.custom {
		if custom.MatchString(message) {
			rules = append(rules, "custom:"+itoa(index))
		}
	}
	if len(rules) == 0 {
		return nil
	}
	family := primaryFamily(rules)
	normalized := Normalize(message)
	timestamp := time.Time{}
	if len(event.Records) > 0 {
		timestamp = event.Records[0].Timestamp
	}
	return &Result{Timestamp: timestamp, Message: message, NormalizedMessage: normalized, Fingerprint: Fingerprint(projectID, family, normalized), PrimaryFamily: family, Severity: severity, Rules: rules}
}

func primaryFamily(rules []string) string {
	priority := []string{"panic", "fatal_level", "process_crash", "container_restart", "stack_trace", "http_5xx", "error_level"}
	joined := "\x00" + strings.Join(rules, "\x00") + "\x00"
	for _, family := range priority {
		if strings.Contains(joined, "\x00"+family+"\x00") {
			return family
		}
	}
	return "custom"
}

func itoa(value int) string {
	if value == 0 {
		return "0"
	}
	var buffer [20]byte
	position := len(buffer)
	for value > 0 {
		position--
		buffer[position] = byte('0' + value%10)
		value /= 10
	}
	return string(buffer[position:])
}
