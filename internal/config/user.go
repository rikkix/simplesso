package config

import (
	"strings"

	"github.com/rikkix/simplesso/utils/errors"
)

const (
	// ErrEmptyName is the error message for empty name.
	ErrEmptyUserName = "User name cannot be empty."

	// ErrEmptyEmail is the error message for empty email.
	ErrEmptyUserEmail = "User email cannot be empty."

	// ErrEmptyTelegramId is the error message for empty telegram id.
	ErrEmptyUserTelegramId = "User telegram id cannot be empty (0)."
)

// Struct User for [[users]] section in config file.
// example:
// ==============================
// name = "sam"
// email = "sam@smith.io"
// telegram_id = 987654321
type User struct {
	Name       string `toml:"name"`
	Email      string `toml:"email"`
	TelegramId int    `toml:"telegram_id"`
}

// Clean and validate user data.
func (u *User) Clean() errors.TraceableError {
	u.Name = strings.TrimSpace(u.Name)
	if u.Name == "" {
		return errors.New(ErrEmptyUserName)
	}

	u.Email = strings.TrimSpace(u.Email)
	u.Email = strings.ToLower(u.Email)
	if u.Email == "" {
		return errors.New(ErrEmptyUserEmail)
	}

	if u.TelegramId == 0 {
		return errors.New(ErrEmptyUserTelegramId)
	}

	return nil
}
