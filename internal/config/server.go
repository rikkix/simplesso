package config

import (
	"strings"

	"github.com/rikkix/simplesso/utils/crypto"
	"github.com/rikkix/simplesso/utils/errors"
)

const (
	ErrEmptyServerTelegramToken = "Telegram token cannot be empty."
	ErrEmptySsoHost = "SSO host cannot be empty."

	DEFAULT_KEY_SIZE = 32
)

// Struct Server for [server] section in config file.
// example:
// ==============================
// listen_address = "localhost:5000"
// sso_jwt_secret = "example_secret"
// services_jwt_secret = "example_secret"
// telegram_token = "00000000:xxxxxxxx"
type Server struct {
	ListenAddress 	string `toml:"listen_address"`
	WebPath 		string `toml:"web_path"`
	SsoHost			string `toml:"sso_host"`
	SsoJwtSecret 	string `toml:"sso_jwt_secret"`
	ServicesJwt 	string `toml:"services_jwt_secret"`
	TelegramToken 	string `toml:"telegram_token"`
	CookieSecure	bool `toml:"cookie_secure"`

	ssoSecretBytes 	[]byte
	servicesSecretBytes []byte
}

// Clean and validate server data.
func (s *Server) Clean() errors.TraceableError {
	s.ListenAddress = strings.TrimSpace(s.ListenAddress)
	if s.ListenAddress == "" {
		s.ListenAddress = "127.0.0.1:5000"
	}

	s.WebPath = strings.TrimSpace(s.WebPath)
	if s.WebPath == "" {
		s.WebPath = "templates/"
	}

	s.SsoHost = strings.TrimSpace(s.SsoHost)
	if s.SsoHost == "" {
		return errors.New(ErrEmptySsoHost)
	}

	s.SsoJwtSecret = strings.TrimSpace(s.SsoJwtSecret)
	if s.SsoJwtSecret == "" {
		s.ssoSecretBytes = crypto.RandomBytes(DEFAULT_KEY_SIZE)
	} else {
		s.ssoSecretBytes = []byte(s.SsoJwtSecret)
	}

	s.ServicesJwt = strings.TrimSpace(s.ServicesJwt)
	if s.ServicesJwt == "" {
		s.servicesSecretBytes = crypto.RandomBytes(DEFAULT_KEY_SIZE)
	} else {
		s.servicesSecretBytes = []byte(s.ServicesJwt)
	}

	s.TelegramToken = strings.TrimSpace(s.TelegramToken)
	if s.TelegramToken == "" {
		return errors.New(ErrEmptyServerTelegramToken)
	}

	return nil
}

// GetSsoSecretBytes returns the secret bytes for the SSO JWT.
func (s *Server) GetSsoSecretBytes() []byte {
	return s.ssoSecretBytes
}

// GetServicesSecretBytes returns the secret bytes for the services JWT.
func (s *Server) GetServicesSecretBytes() []byte {
	return s.servicesSecretBytes
}
