package domain

import "strings"

type CodeChangeType string

const (
	CodeChangeAdded    CodeChangeType = "added"
	CodeChangeModified CodeChangeType = "modified"
	CodeChangeDeleted  CodeChangeType = "deleted"
)

func (t CodeChangeType) Validate() error {
	switch t {
	case CodeChangeAdded, CodeChangeModified, CodeChangeDeleted:
		return nil
	default:
		return ErrInvalidCodeChangeType.Wrap(validationCause("code change type"))
	}
}

type CodeChange struct {
	FilePath   string
	ChangeType CodeChangeType
	Patch      string
	Redacted   bool
}

func NewCodeChange(filePath string, changeType CodeChangeType, patch string) (CodeChange, error) {
	change := CodeChange{
		FilePath: strings.TrimSpace(filePath), ChangeType: changeType, Patch: patch, Redacted: true,
	}
	if err := change.Validate(); err != nil {
		return CodeChange{}, err
	}
	return change, nil
}

func (c CodeChange) Validate() error {
	if !nonBlank(c.FilePath) || len(c.FilePath) > 4096 || c.ChangeType.Validate() != nil ||
		!nonBlank(c.Patch) || len(c.Patch) > 200000 || !c.Redacted {
		return ErrInvalidCodeChange.Wrap(validationCause("code change"))
	}
	return nil
}
