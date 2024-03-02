package tg

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gofiber/fiber/v2/log"
	"github.com/rikkix/simplesso/internal/web/loginreq"
)

type TG struct {
	bot    *tgbotapi.BotAPI
	reqdb  *loginreq.MemDB
	logger log.Logger
}

func New(tgtoken string, reqdb *loginreq.MemDB, logger log.Logger) (*TG, error) {
	bot, err := tgbotapi.NewBotAPI(tgtoken)
	if err != nil {
		return nil, err
	}
	return &TG{
		bot: bot,
		reqdb: reqdb,
		logger: logger,
	}, nil
}

const (
	ConfirmationTemplate = "*New Login Request*\nClick [Confirm] to continue login.\n\nreqid: `%s`"

	ConfirmedTemplate = "*✅ Login Request Confirmed*\nYour code is: `%s`.\n\nreqid: `%s`"

	DeniedTemplate = "*❌ Login Request Denied*\n\nreqid: `%s`"
)

func (t *TG) SendConfirmaion(chatID int64, reqid string) error {
	confirmBtn := tgbotapi.NewInlineKeyboardButtonData("✅ Confirm", fmt.Sprintf("confirm:%s", reqid))
	denyBtn := tgbotapi.NewInlineKeyboardButtonData("❌ Deny", fmt.Sprintf("deny:%s", reqid))
	mkup := tgbotapi.NewInlineKeyboardMarkup([]tgbotapi.InlineKeyboardButton{
		confirmBtn, denyBtn,
	})

	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf(ConfirmationTemplate, reqid))
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = mkup
	_, err := t.bot.Send(msg)
	return err
}

func (t *TG) handleCallbackQuery(update *tgbotapi.Update) {
	data := update.CallbackQuery.Data
	callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "❌ Request not found.")
	text := ""
	needupdate := false

	if strings.HasPrefix(data, "deny:") {
		reqid := strings.TrimPrefix(data, "deny:")
		succ := t.reqdb.RemoveReq(reqid)
		if succ {
			callback.Text = "✅ Request denied."
			text = fmt.Sprintf(DeniedTemplate, reqid)
			needupdate = true
		}
	} else if strings.HasPrefix(data, "confirm:") {
		reqid := strings.TrimPrefix(data, "confirm:")
		succ, code := t.reqdb.Confirm(reqid)
		if succ {
			callback.Text = "✅ Request confirmed."
			text = fmt.Sprintf(ConfirmedTemplate, code, reqid)
			needupdate = true
		}
	}

	if _, err := t.bot.Request(callback); err != nil {
		t.logger.Error("Failed to send callback response:", err)
	}

	if needupdate {
		chatid := update.CallbackQuery.Message.Chat.ID
		msgid := update.CallbackQuery.Message.MessageID
		msg := tgbotapi.NewEditMessageText(chatid, msgid, text)
		msg.ParseMode = "Markdown"
		msg.ReplyMarkup = nil
		if _, err := t.bot.Send(msg); err != nil {
			t.logger.Error("Failed to send message:", err)
		}
	}
}

func (t *TG) StartPolling() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := t.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.CallbackQuery != nil {
			t.handleCallbackQuery(&update)
		}
	}
}
