package detection

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/google/uuid"
)

const fingerprintVersion = "akr-detection-v1"

func Fingerprint(projectID uuid.UUID, family, normalizedMessage string) string {
	payload := fingerprintVersion + "\x00" + projectID.String() + "\x00" + family + "\x00" + normalizedMessage
	sum := sha256.Sum256([]byte(payload))
	return "sha256:" + hex.EncodeToString(sum[:])
}
