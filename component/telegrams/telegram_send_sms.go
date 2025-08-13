package telegrams

import (
	"errors"
	"flag"
	"strconv"
	"strings"

	"github.com/teoit/gosctx"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	ErrTelegramTokenEmpty = errors.New("error Token Telegram Empate")
	ErrTelegramToken      = errors.New("error Token Telegram Auth")
	ErrTelegramDev        = errors.New("dont set env telegram dev true")
)

type SendMessageTelegramSVC interface {
	SendMessageTelegramDev(msg string) error
	SendMessageTelegram(groupId int64, msg string) error
}

type telegramClient struct {
	id          string
	token       string
	dev         bool
	groupdev    string
	groupIdsDev []int64
	logger      gosctx.Logger
	bot         *tgbotapi.BotAPI
}

func NewTelegramClient(id string) *telegramClient {
	return &telegramClient{id: id}
}

func (t *telegramClient) ID() string {
	return t.id
}

func (t *telegramClient) InitFlags() {
	flag.StringVar(&t.token, "telegram-token", "", "token telegram")
	flag.StringVar(&t.groupdev, "telegram-group-dev", "", "group dev telegram")
	flag.BoolVar(&t.dev, "telegram-dev", false, "send dev full message")
}

func (t *telegramClient) Activate(sc gosctx.ServiceContext) error {
	t.logger = sc.Logger(t.id)

	bot, err := tgbotapi.NewBotAPI(t.token)
	if err != nil {
		return ErrTelegramToken
	}
	bot.Debug = false
	t.bot = bot

	if !t.dev {
		return nil
	}

	groupDev := strings.Split(t.groupdev, ",")
	if len(groupDev) == 0 {
		t.logger.Error("error group dev id env not exist")
		return nil
	}
	groupIds := []int64{}
	for _, val := range groupDev {
		idstr := strings.TrimSpace(val)
		if idstr != "" {
			id, err := strconv.ParseInt(idstr, 10, 64)
			if err != nil {
				t.logger.Errorf("error parse group_id - %s", idstr)
				continue
			}
			groupIds = append(groupIds, id)
		}
	}
	t.groupIdsDev = groupIds

	return nil
}

func (t *telegramClient) Stop() error {
	return nil
}

func (t *telegramClient) SendMessageTelegramDev(msg string) error {
	if !t.dev {
		return ErrTelegramDev
	}
	if t.groupIdsDev == nil {
		return nil
	}
	for _, groupId := range t.groupIdsDev {
		if err := t.sendSMS(groupId, msg); err != nil {
			return err
		}
	}
	return nil
}

func (t *telegramClient) SendMessageTelegram(groupId int64, msg string) error {
	if t.dev {
		go t.SendMessageTelegramDev(msg)
	}
	return t.sendSMS(groupId, msg)
}

func (t *telegramClient) sendSMS(groupId int64, msg string) error {
	mes := tgbotapi.NewMessage(groupId, msg)
	mes.ParseMode = tgbotapi.ModeHTML
	_, err := t.bot.Send(mes)
	if err != nil {
		t.logger.Error("error send message dev ", err)
		return err
	}
	return nil
}
