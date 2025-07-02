package telegram

import (
	"fmt"
	"log"
	"slices"
	"strconv"
	"strings"
	"travelWallet/internal/app"
	"travelWallet/internal/domain"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/shopspring/decimal"
)

// Bot представляет нашего Telegram бота и его зависимости.
type Bot struct {
	api          *tgbotapi.BotAPI
	app          *app.Application
	stateManager *StateManager
}

// NewBot создает новый экземпляр бота.
func NewBot(api *tgbotapi.BotAPI, application *app.Application) *Bot {
	return &Bot{
		api:          api,
		app:          application,
		stateManager: NewStateManager(),
	}
}

// Start запускает бота для получения обновлений.
func (b *Bot) Start() error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.api.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		chatID := update.Message.Chat.ID

		if state, data, ok := b.stateManager.GetState(chatID); ok {
			b.handleState(update, state, data)
			continue
		}

		b.handleMessage(update)
	}
	return nil
}

// handleMessage распределяет сообщения без состояния по нужным обработчикам.
func (b *Bot) handleMessage(update tgbotapi.Update) {
	switch update.Message.Text {
	case "/start":
		b.handleStart(update)
	case "Мой кошелек":
		b.handleWalletMenu(update)
	case "Баланс":
		b.handleGetBalance(update)
	case "Внести средства":
		b.handleDeposit(update)
	case "Снять средства":
		b.handleWithdraw(update)
	case "Добавить иностранную валюту":
		b.handleAddForeignCurrency(update)
	case "Удалить иностранную валюту":
		b.handleRemoveForeignCurrency(update)
	case "Изменить базовую валюту":
		b.handleChangeBaseCurrency(update)
	case "Главное меню":
		b.handleBackToMainMenu(update)
	case "Как пользоваться":
		b.handleHelp(update)
	case "Конвертер валют":
		b.handleNothing(update) // TODO

	default:
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Неизвестная команда, давай попробуем по новой?")
		b.api.Send(msg)
		b.handleStart(update)
	}
}

func (b *Bot) handleState(update tgbotapi.Update, state UserState, data map[string]string) {
	chatID, strChatID := update.Message.Chat.ID, strconv.FormatInt(update.Message.Chat.ID, 10)
	text := update.Message.Text

	switch state {
	case StateAwaitingChooseBaseCurrency:
		if len(text) != 3 {
			b.sendMessage(chatID, "Введенная валюта не соотвествует ISO 4217, повторите попытку.", nil)
			return
		}
		currency := strings.ToUpper(text)
		err := b.app.Wallets.Create.Create(strChatID, currency)
		if err != nil {
			b.sendMessage(chatID, "Ошибка при создании кошелька.", WalletSupportKeyboard())
			b.stateManager.SetState(chatID, StateAwaitingHelpRequest, make(map[string]string))
			return
		}

		b.sendMessage(chatID, "Базовая валюта успешно установлена!", MainMenuKeyboard())
		b.stateManager.ClearState(chatID)

	case StateAwaitingDepositAmount:
		foreignCurrencies, err := b.app.Wallets.GetForeign.GetCurrencies(strChatID)
		if err != nil {
			b.sendMessage(chatID, "Не удалось получить сохраненные валюты.", WalletSupportKeyboard())
			b.stateManager.SetState(chatID, StateAwaitingHelpRequest, make(map[string]string))
			return
		}

		baseCurrency, err := b.app.Wallets.GetBaseCurrency.GetBaseCurrency(strChatID)
		if err != nil {
			b.sendMessage(chatID, "Не удалось получить базовую валюту.", WalletSupportKeyboard())
			b.stateManager.SetState(chatID, StateAwaitingHelpRequest, make(map[string]string))
			return
		}

		allCurrencies := append([]string{baseCurrency}, foreignCurrencies...)

		amount, err := strconv.ParseFloat(text, 64)
		if err != nil || amount <= 0 {
			b.sendMessage(chatID, "Неверная сумма. Пожалуйста, введите положительное число.", nil)
			return
		}

		data["amount"] = text

		b.stateManager.SetState(chatID, StateAwaitingDepositCurrency, data)
		b.sendMessage(chatID, "Теперь выберите или введите валюту, в которой будем пополнять", CurrencyKeyboard(allCurrencies))

	case StateAwaitingDepositCurrency:
		currency := strings.ToUpper(text)
		amountStr, ok := data["amount"]
		if !ok {
			b.sendMessage(chatID, "Внутренняя ошибка: не найден контекст запроса. Попробуйте снова.", MainMenuKeyboard())
			b.stateManager.ClearState(chatID)
			return
		}
		amount, _ := strconv.ParseFloat(amountStr, 64)
		amountForCurrency := decimal.NewFromFloat(amount)
		depositCurrency, _ := domain.NewCurrency(amountForCurrency, currency)

		err := b.app.Wallets.Deposit.Deposit(strconv.FormatInt(chatID, 10), depositCurrency)
		if err != nil {
			b.sendMessage(chatID, "Ошибка пополнения: "+err.Error(), WalletMenuKeyboard())
		} else {
			b.sendMessage(chatID, "Баланс успешно пополнен!", WalletMenuKeyboard())
		}
		b.stateManager.ClearState(chatID)

	case StateAwaitingWithdrawAmount:
		foreignCurrencies, err := b.app.Wallets.GetForeign.GetCurrencies(strChatID)
		if err != nil {
			b.sendMessage(chatID, "Не удалось получить сохраненные валюты.", WalletSupportKeyboard())
			b.stateManager.SetState(chatID, StateAwaitingHelpRequest, make(map[string]string))
			return
		}

		baseCurrency, err := b.app.Wallets.GetBaseCurrency.GetBaseCurrency(strChatID)
		if err != nil {
			b.sendMessage(chatID, "Не удалось получить базовую валюту.", WalletSupportKeyboard())
			b.stateManager.SetState(chatID, StateAwaitingHelpRequest, make(map[string]string))
			return
		}

		allCurrencies := append([]string{baseCurrency}, foreignCurrencies...)

		amount, err := strconv.ParseFloat(text, 64)
		if err != nil || amount <= 0 {
			b.sendMessage(chatID, "Неверная сумма. Пожалуйста, введите положительное число.", nil)
			return
		}
		data["amount"] = text
		b.stateManager.SetState(chatID, StateAwaitingWithdrawCurrency, data)
		b.sendMessage(chatID, "Теперь выберите или введите валюту, в которой будем снимать.", CurrencyKeyboard(allCurrencies))

	case StateAwaitingWithdrawCurrency:
		currency := strings.ToUpper(text)
		amountStr, ok := data["amount"]
		if !ok {
			b.sendMessage(chatID, "Произошла ошибка, не найдена сумма. Попробуйте снова.", MainMenuKeyboard())
			b.stateManager.ClearState(chatID)
			return
		}
		amount, _ := strconv.ParseFloat(amountStr, 64)
		amountForCurrency := decimal.NewFromFloat(amount)
		withdrawCurrency, _ := domain.NewCurrency(amountForCurrency, currency)

		err := b.app.Wallets.Withdraw.Withdraw(strconv.FormatInt(chatID, 10), withdrawCurrency)
		if err != nil {
			b.sendMessage(chatID, "Ошибка снятия: "+err.Error(), WalletMenuKeyboard())
		} else {
			b.sendMessage(chatID, "Средства успешно сняты!", WalletMenuKeyboard())
		}
		b.stateManager.ClearState(chatID)

	case StateAwaitingAddForeignCurrency:
		if len(text) != 3 {
			b.sendMessage(chatID, "Некорректный формат валюты - используйте трехзначный код (например, USD) и попробуйте еще раз.", nil)
			return
		}

		inputCurrency := strings.ToUpper(text)
		currencies, err := b.app.Wallets.GetForeign.GetCurrencies(strChatID)
		if err != nil {
			b.sendMessage(chatID, "Ошибка во время получения иностранных валют.", WalletSupportKeyboard())
			b.stateManager.SetState(chatID, StateAwaitingHelpRequest, nil)
			return
		}

		if slices.Contains(currencies, inputCurrency) {
			b.sendMessage(chatID, "Данная валюта уже есть в списке добавленных", WalletMenuKeyboard())
			b.stateManager.ClearState(chatID)
			return
		}

		err = b.app.Wallets.AddForeign.AddCurrency(strChatID, inputCurrency)
		if err != nil {
			b.sendMessage(chatID, fmt.Sprintf("Ошибка во время добавления новой валюты: %s", err.Error()), WalletSupportKeyboard())
			b.stateManager.SetState(chatID, StateAwaitingHelpRequest, nil)
			return
		}

		b.sendMessage(chatID, "Валюта успешно добавлена", WalletMenuKeyboard())
		b.stateManager.ClearState(chatID)

	case StateAwaitingRemoveForeignCurrency:
		if len(text) != 3 {
			b.sendMessage(chatID, "Некорректный формат валюты - используйте трехзначный код (например, USD) и попробуйте еще раз.", nil)
			return
		}

		inputCurrency := strings.ToUpper(text)
		currencies, err := b.app.Wallets.GetForeign.GetCurrencies(strChatID)
		if err != nil {
			b.sendMessage(chatID, "Ошибка во время получения иностранных валют.", WalletSupportKeyboard())
			b.stateManager.SetState(chatID, StateAwaitingHelpRequest, nil)
			return
		}

		if !slices.Contains(currencies, inputCurrency) {
			b.sendMessage(chatID, "Данная валюта отсуствует в списке иностранных валют", WalletMenuKeyboard())
			b.stateManager.ClearState(chatID)
			return
		}

		err = b.app.Wallets.RemoveForeign.RemoveCurrency(strChatID, inputCurrency)
		if err != nil {
			b.sendMessage(chatID, fmt.Sprintf("Ошибка во время удаления валюты: %s", err.Error()), WalletSupportKeyboard())
			b.stateManager.SetState(chatID, StateAwaitingHelpRequest, nil)
			return
		}

		b.sendMessage(chatID, "Валюта успешно удалена", WalletMenuKeyboard())
		b.stateManager.ClearState(chatID)

	case StateAwaitingChangeBaseCurrency:
		if len(text) != 3 {
			b.sendMessage(chatID, "Некорректный формат валюты - используйте трехзначный код (например, USD) и попробуйте еще раз.", nil)
			return
		}

		err := b.app.Wallets.ChangeBase.ChangeBaseCurrency(strChatID, text)
		if err != nil {
			b.sendMessage(chatID, fmt.Sprintf("Ошибка во установки новой базовой валюты: %s", err.Error()), WalletSupportKeyboard())
			b.stateManager.SetState(chatID, StateAwaitingHelpRequest, nil)
			return
		}

		b.sendMessage(chatID, "Новая базовая валюта успешно установлена", WalletMenuKeyboard())
		b.stateManager.ClearState(chatID)

	case StateAwaitingHelpRequest:
		b.sendMessage(chatID, "Помощь: @FoxEdit", tgbotapi.NewRemoveKeyboard(true))
		b.stateManager.ClearState(chatID)
	}
}

// sendMessage - это удобная обертка для отправки сообщений.
func (b *Bot) sendMessage(chatID int64, text string, keyboard interface{}) {
	msg := tgbotapi.NewMessage(chatID, text)
	if keyboard != nil {
		msg.ReplyMarkup = keyboard
	}
	_, err := b.api.Send(msg)
	if err != nil {
		log.Printf("Failed to send message to chat %d: %v", chatID, err)
	}
}
