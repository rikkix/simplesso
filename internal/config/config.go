package config

import (
	"os"

	"github.com/pelletier/go-toml/v2"
	"github.com/rikkix/simplesso/utils/errors"
)

const (
	// ErrWrongServerSection is the error message for wrong server section.
	ErrWrongServerSection = "Wrong server section."
	// ErrWrongGitHubSection is the error message for wrong github section.
	ErrWrongGitHubSection = "Wrong github section."
	// ErrWrongUsersSection is the error message for wrong users section.
	ErrWrongUsersSection = "Wrong users section."
	// ErrWrongTokensSection is the error message for wrong tokens section.
	ErrWrongTokensSection = "Wrong tokens section."
	// ErrWrongServicesSection is the error message for wrong services section.
	ErrWrongServicesSection = "Wrong services section."

	// ErrOpeningConfigFile is the error message for opening config file.
	ErrOpeningConfigFile = "Error opening config file."
	// ErrDecodingConfigFile is the error message for decoding config file.
	ErrDecodingConfigFile = "Error decoding config file."
	// ErrCleaningConfigData is the error message for cleaning config data.
	ErrCleaningConfigData = "Error cleaning config data."
)

type Config struct {
	Server   Server   `toml:"server"`
	GitHub  	GitHub   `toml:"github"`
	Users    []User   `toml:"users"`
	Tokens   []Token   `toml:"tokens"`
	Services []Service `toml:"services"`
	users	 map[string]*User
	users_github map[string]*User
	services map[string]*Service
}

// Clean and validate config data.
func (c *Config) Clean() errors.TraceableError {
	c.users = make(map[string]*User)
	c.users_github = make(map[string]*User)
	c.services = make(map[string]*Service)
	
	var err errors.TraceableError
	err = c.Server.Clean()
	if err != nil {
		return err.From(ErrWrongServerSection)
	}

	err = c.GitHub.Clean()
	if err != nil {
		return err.From(ErrWrongGitHubSection)
	}

	for i, u := range c.Users {
		err = u.Clean()
		if err != nil {
			return err.From(ErrWrongUsersSection)
		}
		c.Users[i] = u
		c.users[u.Name] = &c.Users[i]
		c.users_github[u.GitHub] = &c.Users[i]
	}

	for i, t := range c.Tokens {
		err = t.Clean()
		if err != nil {
			return err.From(ErrWrongTokensSection)
		}
		c.Tokens[i] = t
	}

	for i, s := range c.Services {
		err = s.Clean()
		if err != nil {
			return err.From(ErrWrongServicesSection)
		}
		c.Services[i] = s
		c.services[s.Host] = &c.Services[i]
	}

	return nil
}

// Parse config from toml file.
func FromFile(path string) (*Config, errors.TraceableError) {
	// Check file existence.
	f, e := os.Open(path)
	if e != nil {
		return nil, errors.New(e.Error()).From(ErrOpeningConfigFile)
	}

	config := Config{}
	// Decode toml file.
	dec := toml.NewDecoder(f)
	e = dec.Decode(&config)
	if e != nil {
		return nil, errors.New(e.Error()).From(ErrDecodingConfigFile)
	}

	// Clean and validate config data.
	err := config.Clean()
	if err != nil {
		return nil, err.From(ErrCleaningConfigData)
	}

	return &config, nil
}

func (c *Config) FindUser(name string) *User {
	return c.users[name]
}

func (c *Config) FindService(host string) *Service {
	return c.services[host]
}

func (c *Config) FindUserByGitHub(name string) *User {
	return c.users_github[name]
}

