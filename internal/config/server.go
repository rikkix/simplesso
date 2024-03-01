package config

import (
	"strings"

	"github.com/rikkix/simplesso/utils/crypto"
	"github.com/rikkix/simplesso/utils/errors"
)

const (
	ErrEmptyServerTelegramToken = "Telegram token cannot be empty."

	DEFAULT_KEY_SIZE = 32
)

// Struct Server for [server] section in config file.
// example:
// ==============================
// sso_jwt_secret = "example_secret"
// services_jwt_secret = "example_secret"
// telegram_token = "00000000:xxxxxxxx"
type Server struct {
	SsoJwtSecret 	string `toml:"sso_jwt_secret"`
	ServicesJwt 	string `toml:"services_jwt_secret"`
	TelegramToken 	string `toml:"telegram_token"`
}

// Clean and validate server data.
func (s *Server) Clean() errors.TraceableError {
	s.SsoJwtSecret = strings.TrimSpace(s.SsoJwtSecret)
	if s.SsoJwtSecret == "" {
		s.SsoJwtSecret = crypto.HexString(DEFAULT_KEY_SIZE)
	}

	s.ServicesJwt = strings.TrimSpace(s.ServicesJwt)
	if s.ServicesJwt == "" {
		s.ServicesJwt = crypto.HexString(DEFAULT_KEY_SIZE)
	}

	s.TelegramToken = strings.TrimSpace(s.TelegramToken)
	if s.TelegramToken == "" {
		return errors.New(ErrEmptyServerTelegramToken)
	}
}