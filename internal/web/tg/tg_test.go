package tg_test

import (
	"os"
	"strconv"
	"testing"

	"github.com/rikkix/simplesso/internal/web/tg"
)

func TestTG_SendConfirmaion(t *testing.T) {
	tg, err := tg.New(os.Getenv("TG_TOKEN"))
	if err != nil {
		t.Fatal(err)
	}
	chatid, err := strconv.Atoi(os.Getenv("TG_CHATID"))
	if err != nil {
		t.Fatal(err)
	}
	err = tg.SendConfirmaion(int64(chatid), "test-reqid")
	if err != nil {
		t.Fatal(err)
	}

	
}