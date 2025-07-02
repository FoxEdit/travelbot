package telegram

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"travelWallet/internal/domain"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var ErrWalletNotFound = errors.New("кошелек не найден")

func (b *Bot) handleStart(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID

	isWalletAlreadyExists, err := b.app.Wallets.IsExists.IsExists(strconv.FormatInt(chatID, 10))
	if err != nil {
		b.sendMessage(chatID, fmt.Sprintf("Ошибка идентефикации пользователя: %s.", err.Error()), WalletSupportKeyboard()) // TODO handle this
		b.stateManager.SetState(chatID, StateAwaitingHelpRequest, nil)
		return
	}

	if isWalletAlreadyExists {
		text := "С возвращением!"
		b.sendMessage(chatID, text, MainMenuKeyboard())
		b.stateManager.ClearState(chatID)
		return
	}

	welcomeText := "Добро пожаловать в Travel Wallet Bot! Этот бот создан для ручного ведения счета в условиях отсуствия интернет банкинга, но при наличии пластиковой карты."
	guideText := "Для создания кошелька напишите вашу основную валюту для расплаты (базовую, которой будете пополнять ваш баланс) в формате международном стандарте <i>(RUB, KGS, KZT...)</i>"

	welcomeMsg := tgbotapi.NewMessage(chatID, welcomeText)
	welcomeMsg.ParseMode = tgbotapi.ModeHTML
	guideMsg := tgbotapi.NewMessage(chatID, guideText)
	guideMsg.ParseMode = tgbotapi.ModeHTML

	b.stateManager.SetState(chatID, StateAwaitingChooseBaseCurrency, make(map[string]string))
	b.api.Send(welcomeMsg)
	b.api.Send(guideMsg)
}

func (b *Bot) handleWalletMenu(update tgbotapi.Update) {
	b.sendMessage(update.Message.Chat.ID, "Выберите действие с кошельком:", WalletMenuKeyboard())
}

func (b *Bot) handleBackToMainMenu(update tgbotapi.Update) {
	b.sendMessage(update.Message.Chat.ID, "Главное меню:", MainMenuKeyboard())
}

func (b *Bot) handleAddForeignCurrency(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID

	formattedText := "Введите новую валюту трехсимвольного формата\n<blockquote>USD/EUR/RUB</blockquote>"
	msg := tgbotapi.NewMessage(chatID, formattedText)
	msg.ParseMode = tgbotapi.ModeHTML
	b.api.Send(msg)

	b.stateManager.SetState(chatID, StateAwaitingAddForeignCurrency, make(map[string]string))
}

func (b *Bot) handleGetBalance(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	strChatID := strconv.FormatInt(chatID, 10)

	errorProcess := func(chatID int64, err error) {
		b.sendMessage(chatID, "Ошибка получения баланса: "+err.Error(), WalletSupportKeyboard())
		b.stateManager.SetState(chatID, StateAwaitingHelpRequest, nil)
	}

	balance, err := b.app.Wallets.GetBalance.GetBalance(strChatID)

	if err != nil {
		errorProcess(chatID, err)
		return
	}

	balanceBaseCurrency, err := b.app.Wallets.GetBaseCurrency.GetBaseCurrency(strChatID)
	if err != nil {
		errorProcess(chatID, err)
		return
	}

	foreignCurrencies, err := b.app.Wallets.GetForeign.GetCurrencies(strChatID)
	if err != nil {
		errorProcess(chatID, err)
		return
	}

	baseCurrencyDomainStructure, _ := domain.NewCurrency(balance, balanceBaseCurrency)
	foreignCurrenciesBalance, err := b.app.Exchange.ConvertCurrency.ConvertFromBaseToMany(baseCurrencyDomainStructure, foreignCurrencies...) // TODO replace 1st arg "domain.Currency" to string?
	if err != nil {
		errorProcess(chatID, err)
		return
	}

	var responseText strings.Builder
	responseText.WriteString(fmt.Sprintf("Ваш баланс: %s %s", balance.String(), balanceBaseCurrency))

	if len(foreignCurrenciesBalance) != 0 {
		responseText.WriteString("\n\nБаланс в валюте:\n")
		for _, v := range foreignCurrenciesBalance {
			responseText.WriteString(fmt.Sprintf("💠 %s: %s\n", v.GetCurrencyCode(), v.GetAmount().Mul(balance).StringFixed(2)))
		}
	}

	b.sendMessage(chatID, responseText.String(), WalletMenuKeyboard())
}

func (b *Bot) handleRemoveForeignCurrency(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	strChatID := strconv.FormatInt(chatID, 10)

	foreignCurrencies, err := b.app.Wallets.GetForeign.GetCurrencies(strChatID)
	if err != nil {
		b.sendMessage(chatID, "Ошибка получения списка иностранных валют: "+err.Error(), WalletSupportKeyboard())
		b.stateManager.SetState(chatID, StateAwaitingHelpRequest, nil)
		return
	}

	if len(foreignCurrencies) == 0 {
		b.sendMessage(chatID, "У вас нет добавленных иностранных валют.", WalletMenuKeyboard())
		b.stateManager.ClearState(chatID)
		return
	}

	formattedText := "Введите валюту которую хотите удалить."
	b.sendMessage(chatID, formattedText, CurrencyKeyboard(foreignCurrencies))
	b.stateManager.SetState(chatID, StateAwaitingRemoveForeignCurrency, make(map[string]string))
}

func (b *Bot) handleChangeBaseCurrency(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID

	formattedText := "Введите новую базовую валюту трехсимвольного формата\n<blockquote>USD/EUR/RUB</blockquote>"
	msg := tgbotapi.NewMessage(chatID, formattedText)
	msg.ParseMode = tgbotapi.ModeHTML
	b.api.Send(msg)

	b.stateManager.SetState(chatID, StateAwaitingChangeBaseCurrency, make(map[string]string))
}

func (b *Bot) handleHelp(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID

	responseText := `
<b>Добро пожаловать в Справку!</b>

Это ваш гид по основным функциям <i>Travel Wallet Bot</i>.

<b>Основные операции</b>
• <code>Баланс</code> — покажет баланс в вашей основной валюте и иностранных валютах, если они были добавлены.
• <code>Внести средства</code> — запустит диалог для пополнения счета. Вам нужно будет указать сумму и валюту.
• <code>Снять средства</code> — аналогично, но для снятия денег со счета.
• <code>Добавить иностранную валюту</code> — добавит нужную вам валюту в список при просмотре баланса.
• <code>Удалить иностранную валюту</code> — удалит ненужную вам валюту из списка при просмотре баланса.
• <code>Изменить базовую валюту</code> — изменит вашу основную валюту на выбранную.

<b>Что в разработке?</b>
<i>Скоро появится:</i> 
• Обмен валют и калькулятор при вводе любых чисел (ваш запрос сможет быть в такой форме: 12.45+4.35+1+5.4, что удобно при коррекции баланса с магазинных чеков).
• Прерывание операций (например, если передумали пополнять баланс)
`

	msg := tgbotapi.NewMessage(chatID, responseText)
	msg.ParseMode = tgbotapi.ModeHTML
	msg.ReplyMarkup = MainMenuKeyboard()

	b.api.Send(msg)
	b.stateManager.ClearState(chatID)
}

func (b *Bot) handleDeposit(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	b.stateManager.SetState(chatID, StateAwaitingDepositAmount, make(map[string]string))
	b.sendMessage(chatID, "Введите сумму для пополнения:", tgbotapi.NewRemoveKeyboard(true))
}

func (b *Bot) handleWithdraw(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	b.stateManager.SetState(chatID, StateAwaitingWithdrawAmount, make(map[string]string))
	b.sendMessage(chatID, "Введите сумму для снятия:", tgbotapi.NewRemoveKeyboard(true))
}

func (b *Bot) handleNothing(update tgbotapi.Update) {
	chatID := update.Message.Chat.ID
	b.sendMessage(chatID, "Функция в разработке.", nil)
}
