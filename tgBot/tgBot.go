package tgBot

import (
	"context"
	"log"
	envhandler "mashinki/envHandler"
	"mashinki/logging"
	"mashinki/parser"
	"mashinki/taxes"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	cmdStart   = "start"
	btnFindCar = "🚗 Найти информацию о машине"
)

var (
	// Ading button
	mainKeyboard = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(btnFindCar),
		),
	)
)

// Состояния пользователя
type UserState struct {
	WaitingForURL bool
}

type Bot struct {
	api         *tgbotapi.BotAPI
	cancel      context.CancelFunc
	userStates  map[int64]*UserState
	statesMutex sync.RWMutex
}

func StartBot() (*Bot, error) {
	token := envhandler.GetEnv("TG_TOKEN")
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	log.Printf("Authorized on account %s", api.Self.UserName)

	ctx, cancel := context.WithCancel(context.Background())
	bot := &Bot{
		api:         api,
		cancel:      cancel,
		userStates:  make(map[int64]*UserState),
		statesMutex: sync.RWMutex{},
	}

	go bot.run(ctx)

	return bot, nil
}

func (b *Bot) Stop() {
	if b.cancel != nil {
		b.cancel()
	}
}

func (b *Bot) getUserState(chatID int64) *UserState {
	b.statesMutex.RLock()
	state, exists := b.userStates[chatID]
	b.statesMutex.RUnlock()

	if !exists {
		b.statesMutex.Lock()
		state = &UserState{}
		b.userStates[chatID] = state
		b.statesMutex.Unlock()
	}

	return state
}

func (b *Bot) setUserState(chatID int64, state *UserState) {
	b.statesMutex.Lock()
	b.userStates[chatID] = state
	b.statesMutex.Unlock()
}

func (b *Bot) handleMessage(ctx context.Context, update tgbotapi.Update) {
	if update.Message == nil {
		return
	}

	chatID := update.Message.Chat.ID
	state := b.getUserState(chatID)

	var msg tgbotapi.MessageConfig

	switch {
	case update.Message.Command() == cmdStart:
		msg = tgbotapi.NewMessage(chatID, "Привет! Я помогу найти информацию о машине и рассчитаю таможенные платежи. Нажми на кнопку ниже:")
		msg.ReplyMarkup = mainKeyboard

	case update.Message.Text == btnFindCar:
		state.WaitingForURL = true
		b.setUserState(chatID, state)
		msg = tgbotapi.NewMessage(chatID, "Отправь мне ссылку на машину с сайта che168.com")

	case state.WaitingForURL:
		state.WaitingForURL = false
		b.setUserState(chatID, state)

		processingMsg := tgbotapi.NewMessage(chatID, "🔄 Получаю информацию о машине и рассчитываю таможенные платежи...")
		if _, err := b.api.Send(processingMsg); err != nil {
			logging.DefaultLogger.LogErrorF("Error sending processing message: %v", err)
		}

		carInfo, err := parser.GetCarInfo(update.Message.Text)
		if err != nil {
			logging.DefaultLogger.LogErrorF("Error getting car info: %v", err)
			msg = tgbotapi.NewMessage(chatID, "❌ Ошибка при получении информации о машине")
		} else {
			fullInfo := taxes.NewFullCarInfo(&carInfo)
			msg = tgbotapi.NewMessage(chatID, "✅ "+fullInfo.String())
		}
		msg.ReplyMarkup = mainKeyboard
		msg.ParseMode = "Markdown"

	default:
		msg = tgbotapi.NewMessage(chatID, "Нажми на кнопку ниже, чтобы начать:")
		msg.ReplyMarkup = mainKeyboard
	}

	if _, err := b.api.Send(msg); err != nil {
		log.Printf("Error sending message: %v", err)
	}
}

func (b *Bot) run(ctx context.Context) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	// Making 100 chans for users
	const maxWorkers = 100
	sem := make(chan struct{}, maxWorkers)

	for {
		select {
		case <-ctx.Done():
			return
		case update := <-updates:
			// Start a new goroutine for each update
			sem <- struct{}{}
			go func(update tgbotapi.Update) {
				defer func() { <-sem }()
				b.handleMessage(ctx, update)
			}(update)
		}
	}
}
