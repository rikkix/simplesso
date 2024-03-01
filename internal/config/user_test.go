package config_test

import (
	"testing"

	"github.com/pelletier/go-toml/v2"
	"github.com/rikkix/simplesso/internal/config"
)

func TestUserUnmarshalAndClean(t *testing.T) {
	doc := `
name = "sam"
email = "sam@smith.io"
telegram_id = 987654321
	`

	expected := config.User{
		Name:       "sam",
		Email:      "sam@smith.io",
		TelegramId: 987654321,
	}

	var u config.User

	err := toml.Unmarshal([]byte(doc), &u)
	if err != nil {
		t.Fatal(err)
	}

	if u != expected {
		t.Errorf("expected: %v, got: %v", expected, u)
	}

	err = u.Clean()
	if err != nil {
		t.Fatal(err)
	}
}

func TestUserUnmarshalAndClean_OmittedField(t *testing.T) {
	doc := `
	name = "sam"
	telegram_id = 987654321
	`

	expected := config.User{
		Name:       "sam",
		Email:      "",
		TelegramId: 987654321,
	}

	var u config.User

	err := toml.Unmarshal([]byte(doc), &u)
	if err != nil {
		t.Fatal(err)
	}

	if u != expected {
		t.Errorf("expected: %v, got: %v", expected, u)
	}

	terr := u.Clean()
	if terr.Last() != config.ErrEmptyEmail {
		t.Errorf("expected: %v, got: %v", config.ErrEmptyEmail, terr.Last())
	}
}