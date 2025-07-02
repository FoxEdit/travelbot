package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func MainMenuKeyboard() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Мой кошелек"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Конвертер валют"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Как пользоваться"),
		),
	)
}

func WalletMenuKeyboard() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Баланс"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Внести средства"),
			tgbotapi.NewKeyboardButton("Снять средства"),
		),
		tgbotapi.NewKeyboardButtonRow(

			tgbotapi.NewKeyboardButton("Добавить иностранную валюту"),
			tgbotapi.NewKeyboardButton("Удалить иностранную валюту"),
			tgbotapi.NewKeyboardButton("Изменить базовую валюту"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Главное меню"),
		),
	)
}

func WalletSupportKeyboard() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Обратится в поддержку"),
		),
	)
}

func CurrencyKeyboard(currencies []string) tgbotapi.ReplyKeyboardMarkup {
	var rowBuffer []tgbotapi.KeyboardButton
	var rows [][]tgbotapi.KeyboardButton
	rowElementsNumber := 2

	for _, text := range currencies {
		btn := tgbotapi.NewKeyboardButton(text)
		rowBuffer = append(rowBuffer, btn)

		// Добавляем, только когда ряд ПОЛНОСТЬЮ заполнился
		if len(rowBuffer) == rowElementsNumber {
			rows = append(rows, rowBuffer)
			rowBuffer = []tgbotapi.KeyboardButton{}
		}
	}

	if len(rowBuffer) > 0 {
		rows = append(rows, rowBuffer)
	}

	keyboard := tgbotapi.NewReplyKeyboard(rows...)
	keyboard.ResizeKeyboard = true
	return keyboard
}
