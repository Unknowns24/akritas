package incident

import "time"

type GitHubIssueReferenceDTO struct {
	Number     int       `json:"number"`
	URL        string    `json:"url"`
	Repository string    `json:"repository"`
	CreatedAt  time.Time `json:"created_at"`
}
