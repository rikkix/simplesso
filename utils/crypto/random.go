package crypto

import (
	"crypto/rand"
	"fmt"
)

// HexString returns a hex string of nbytes random bytes.
// This function is cryptographically secure.
func HexString(nbytes int) string {
	b := make([]byte, nbytes)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%x", b)
}