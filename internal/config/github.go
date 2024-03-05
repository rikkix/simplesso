package config

import (
	"strings"

	"github.com/rikkix/simplesso/utils/errors"
)

type GitHub struct {
	ClientID	 string `toml:"client_id"`
	ClientSecret string `toml:"client_secret"`
}

func (g *GitHub) Clean() errors.TraceableError {
	g.ClientID = strings.TrimSpace(g.ClientID)

	g.ClientSecret = strings.TrimSpace(g.ClientSecret)
	
	return nil
}