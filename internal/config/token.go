package config

import (
	"strings"

	"github.com/rikkix/simplesso/utils/errors"
)

const (
	// ErrEmptyTokenName is the error message for empty token name.
	ErrEmptyTokenName = "Token name cannot be empty."

	// ErrEmptyToken is the error message for empty token.
	ErrEmptyToken = "Token cannot be empty."
)

// Struct Token for [[tokens]] section in config file.
// example:
// ==============================
// name = "example-name"
// token = "example_token"
type Token struct {
	Name  string `toml:"name"`
	Token string `toml:"token"`
}

func (t *Token) Clean() errors.TraceableError {
	t.Name = strings.TrimSpace(t.Name)
	if t.Name == "" {
		return errors.New(ErrEmptyTokenName)
	}

	t.Token = strings.TrimSpace(t.Token)
	if t.Token == "" {
		return errors.New(ErrEmptyToken)
	}

	return nil
}