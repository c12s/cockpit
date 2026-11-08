package utils

import (
	"encoding/hex"
	"fmt"
	"strings"
)

// ValidateDigest enforces the digest form the node and gateway accept:
// sha256:<64 lowercase hex>.
func ValidateDigest(digest string) error {
	const prefix = "sha256:"
	if !strings.HasPrefix(digest, prefix) || len(digest) != len(prefix)+64 {
		return fmt.Errorf("digest %q must be sha256:<64 hex>", digest)
	}
	if _, err := hex.DecodeString(digest[len(prefix):]); err != nil {
		return fmt.Errorf("digest %q is not valid hex", digest)
	}
	return nil
}
