package crypto

import (
	"crypto/rand"
)

// RandomBytes returns a random byte slice of the given length.
func RandomBytes(nbytes int) []byte {
	b := make([]byte, nbytes)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return b
}

const CHARS = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const LEN_CHARS = len(CHARS)
// RandomString returns a random string of the given length.
func RandomString(nchars int) string {
	b := RandomBytes(nchars)
	for i, c := range b {
		b[i] = CHARS[int(c)%LEN_CHARS]
	}
	return string(b)
}

const DIGIT_CHARS = "0123456789"
const LEN_DIGIT_CHARS = len(DIGIT_CHARS)

func RandomDigits(nchars int) string {
	b := RandomBytes(nchars)
	for i, c := range b {
		b[i] = DIGIT_CHARS[int(c)%LEN_DIGIT_CHARS]
	}
	return string(b)
}
