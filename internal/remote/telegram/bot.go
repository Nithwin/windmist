package telegram

import (
	"fmt"
	"strconv"

	"github.com/Nithwin/WindMist/internal/config"
	"github.com/Nithwin/WindMist/internal/remote"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramController struct {
	bot       *tgbotapi.BotAPI
	cfg       config.TelegramConfig
	allowedID int64
	hub       *remote.Hub
	stopChan  chan struct{}
}

func New(cfg config.TelegramConfig) (*TelegramController, error) {
	if cfg.BotToken == "" {
		return nil, fmt.Errorf("telegram bot token is required")
	}

	id, err := strconv.ParseInt(cfg.AllowedID, 10, 64)
	if err != nil {
		id = 0 // Allow username based validation later, or require numeric for now
	}

	return &TelegramController{
		cfg:       cfg,
		allowedID: id,
		stopChan:  make(chan struct{}),
	}, nil
}

func (t *TelegramController) Name() string {
	return "telegram"
}

func (t *TelegramController) Start(hub *remote.Hub) error {
	t.hub = hub
	bot, err := tgbotapi.NewBotAPI(t.cfg.BotToken)
	if err != nil {
		return fmt.Errorf("failed to connect to telegram: %w", err)
	}
	t.bot = bot

	go t.listen()
	return nil
}

func (t *TelegramController) Stop() error {
	close(t.stopChan)
	if t.bot != nil {
		t.bot.StopReceivingUpdates()
	}
	return nil
}

func (t *TelegramController) SendMessage(text string) error {
	if t.bot == nil || t.allowedID == 0 {
		return nil
	}

	// Telegram messages must not exceed 4096 characters
	if len(text) > 4000 {
		text = text[:4000] + "\n...[truncated]"
	}

	msg := tgbotapi.NewMessage(t.allowedID, text)
	// Optionally set Markdown parse mode, but WindMist markdown might break Telegram's strict MDV2 parser
	// msg.ParseMode = "Markdown"

	_, err := t.bot.Send(msg)
	return err
}

func (t *TelegramController) listen() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := t.bot.GetUpdatesChan(u)

	for {
		select {
		case update := <-updates:
			if update.Message == nil { // ignore any non-Message Updates
				continue
			}

			// Validate user
			if t.allowedID == 0 {
				t.allowedID = update.Message.From.ID
				t.SendMessage(fmt.Sprintf("✅ Telegram bot auto-bound to your user ID: %d", t.allowedID))
			} else if update.Message.From.ID != t.allowedID {
				msg := tgbotapi.NewMessage(update.Message.From.ID, "❌ Unauthorized access attempt.")
				_, _ = t.bot.Send(msg)
				continue
			}

			if update.Message.IsCommand() {
				t.handleCommand(update.Message)
			} else {
				// Treat normal text as a prompt
				t.hub.Incoming <- remote.Command{Type: "ask", Args: []string{update.Message.Text}}
			}
		case <-t.stopChan:
			return
		}
	}
}

func (t *TelegramController) handleCommand(message *tgbotapi.Message) {
	switch message.Command() {
	case "status":
		t.SendMessage("🌀 WindMist is running and connected!")
	case "providers":
		t.hub.Incoming <- remote.Command{Type: "list_providers"}
	case "models":
		t.hub.Incoming <- remote.Command{Type: "list_models"}
	case "provider":
		args := message.CommandArguments()
		if args == "" {
			t.SendMessage("Usage: /provider <name>")
			return
		}
		t.hub.Incoming <- remote.Command{Type: "provider", Args: []string{args}}
	case "model":
		args := message.CommandArguments()
		if args == "" {
			t.SendMessage("Usage: /model <name>")
			return
		}
		t.hub.Incoming <- remote.Command{Type: "model", Args: []string{args}}
	case "ask":
		args := message.CommandArguments()
		if args == "" {
			t.SendMessage("Usage: /ask <prompt>")
			return
		}
		t.hub.Incoming <- remote.Command{Type: "ask", Args: []string{args}}
	default:
		t.SendMessage("Unknown command. Supported:\n/status\n/providers\n/provider <name>\n/models\n/model <name>\n/ask <prompt>")
	}
}
