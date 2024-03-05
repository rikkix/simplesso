package config

import (
	"strings"

	"github.com/rikkix/simplesso/utils/errors"
)

const (
	// ErrEmptyServiceName is the error message for empty service name.
	ErrEmptyServiceName = "Service name cannot be empty."
	// ErrEmptyServiceHost is the error message for empty service host.
	ErrEmptyServiceHost = "Service host cannot be empty."
)

// Struct Service for [[services]] section in config file.
// example:
// ==============================
// name = "code-server"
// host = "code.example.com"
// users = [ "john" ]
// tokens = [ "local-pc", "my-server" ]
// bypass = [ "/static/", "/some-path/" ]
type Service struct {
	Name    string   `toml:"name"`
	Host    string   `toml:"host"`
	Users   []string `toml:"users"`
	users_map map[string]bool
	Tokens  []string `toml:"tokens"`
	tokens_map map[string]bool
	Bypass  []string `toml:"bypass"`
}

// Clean and validate service data.
func (s *Service) Clean() errors.TraceableError {
	s.Name = strings.TrimSpace(s.Name)
	if s.Name == "" {
		return errors.New(ErrEmptyServiceName)
	}

	s.Host = strings.ToLower(strings.TrimSpace(s.Host))
	if s.Host == "" {
		return errors.New(ErrEmptyServiceHost)
	}

	for i, u := range s.Users {
		s.Users[i] = strings.TrimSpace(u)
	}
	s.users_map = make(map[string]bool)
	for _, u := range s.Users {
		s.users_map[u] = true
	}

	for i, t := range s.Tokens {
		s.Tokens[i] = strings.TrimSpace(t)
	}
	s.tokens_map = make(map[string]bool)
	for _, t := range s.Tokens {
		s.tokens_map[t] = true
	}

	for i, p := range s.Bypass {
		s.Bypass[i] = strings.ToLower(strings.TrimSpace(p))
	}

	return nil
}

// Check if user is allowed to access the service.
func (s *Service) IsUserAllowed(user string) bool {
	_, ok := s.users_map[user]
	return ok
}

// Check if token is allowed to access the service.
func (s *Service) IsTokenAllowed(token string) bool {
	_, ok := s.tokens_map[token]
	return ok
}

// Check if path is allowed to bypass the service.
func (s *Service) IsPathAllowed(path string) bool {
	if len(s.Bypass) == 0 {
		return false
	}
	path = strings.ToLower(path)
	for _, p := range s.Bypass {
		if strings.HasPrefix(path, p) {
			return true
		}
	}
	return false
}

func (s *Service) IsBypass(path string) bool {
	if len(s.Bypass) == 0 { return false }
	path = strings.ToLower(strings.TrimSpace(path))
	for _, p := range s.Bypass {
		if strings.HasPrefix(path, p) {
			return true
		}
	}
	return false
}