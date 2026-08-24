package qvac

import "errors"

var (
	ErrInvalidEndpoint    = errors.New("QVAC endpoint must be a local HTTP(S) URL")
	ErrUnavailable        = errors.New("QVAC runtime is unavailable")
	ErrModelUnavailable   = errors.New("QVAC model is unavailable")
	ErrContextOverflow    = errors.New("QVAC context window was exceeded")
	ErrInvalidModelOutput = errors.New("QVAC returned an invalid investigation result")
	ErrUnknownTool        = errors.New("QVAC requested an unknown tool")
	ErrToolLimitExceeded  = errors.New("QVAC tool-call limit exceeded")
)
