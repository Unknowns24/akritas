package commitcorrelation

import (
	"sort"
	"strings"
	"time"

	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
	"github.com/Unknowns24/akritas/backend/internal/service/evidencesafety"
)

const (
	DefaultLimit = 5
	beforeWindow = 24 * time.Hour
	afterWindow  = time.Hour
)

type IncidentWindow struct {
	FirstSeenAt time.Time
	LastSeenAt  time.Time
	Files       []string
}

type CorrelatedCommit struct {
	SHA       string
	Timestamp time.Time
	Author    string
	Message   string
	URL       string
	Reason    string
}

func Select(window IncidentWindow, commits []portsout.RepositoryCommitSummary, limit int) []CorrelatedCommit {
	if limit <= 0 || limit > DefaultLimit {
		limit = DefaultLimit
	}
	firstSeen := window.FirstSeenAt.UTC()
	lastSeen := window.LastSeenAt.UTC()
	if lastSeen.IsZero() {
		lastSeen = firstSeen
	}
	type scored struct {
		commit CorrelatedCommit
		score  int
		index  int
	}
	seen := map[string]struct{}{}
	selected := make([]scored, 0, len(commits))
	for index, commit := range commits {
		sha := strings.TrimSpace(commit.SHA)
		if sha == "" {
			continue
		}
		if _, duplicate := seen[sha]; duplicate {
			continue
		}
		seen[sha] = struct{}{}
		timestamp, _ := time.Parse(time.RFC3339, strings.TrimSpace(commit.Date))
		score, reason := scoreCommit(window, timestamp)
		if score < 0 {
			continue
		}
		selected = append(selected, scored{
			commit: CorrelatedCommit{
				SHA:       evidencesafety.RedactAndLimit(sha, 64),
				Timestamp: timestamp.UTC(),
				Author:    evidencesafety.RedactAndLimit(commit.Author, 200),
				Message:   evidencesafety.RedactAndLimit(commit.Message, 1000),
				URL:       evidencesafety.RedactAndLimit(commit.URL, 1000),
				Reason:    reason,
			},
			score: score,
			index: index,
		})
	}
	sort.SliceStable(selected, func(left, right int) bool {
		if selected[left].score != selected[right].score {
			return selected[left].score > selected[right].score
		}
		leftTime, rightTime := selected[left].commit.Timestamp, selected[right].commit.Timestamp
		if !leftTime.Equal(rightTime) {
			return leftTime.After(rightTime)
		}
		if selected[left].commit.SHA != selected[right].commit.SHA {
			return selected[left].commit.SHA < selected[right].commit.SHA
		}
		return selected[left].index < selected[right].index
	})
	if len(selected) > limit {
		selected = selected[:limit]
	}
	out := make([]CorrelatedCommit, 0, len(selected))
	for _, item := range selected {
		out = append(out, item.commit)
	}
	return out
}

func scoreCommit(window IncidentWindow, timestamp time.Time) (int, string) {
	if timestamp.IsZero() {
		return 1, "Commit reciente sin timestamp parseable; se conserva como contexto, no como causa confirmada."
	}
	firstSeen := window.FirstSeenAt.UTC()
	lastSeen := window.LastSeenAt.UTC()
	if lastSeen.IsZero() {
		lastSeen = firstSeen
	}
	if !firstSeen.IsZero() && timestamp.Before(firstSeen.Add(-beforeWindow)) {
		return -1, ""
	}
	if !lastSeen.IsZero() && timestamp.After(lastSeen.Add(afterWindow)) {
		return -1, ""
	}
	if !firstSeen.IsZero() && (timestamp.Equal(firstSeen) || timestamp.Before(firstSeen)) {
		return 3, "Commit anterior cercano al incidente; potencialmente relacionado, no causalidad confirmada."
	}
	return 2, "Commit dentro de la ventana temporal del incidente; potencialmente relacionado, no causalidad confirmada."
}
