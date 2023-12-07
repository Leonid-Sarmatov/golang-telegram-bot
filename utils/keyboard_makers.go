package utils

import (
	types "golang_telegram_bot/types"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Функция для создания горизонтальной клавиатуры
func HorizontalDataInlineKeyboardMaker(names, callbacks []string) tgbotapi.InlineKeyboardMarkup {
	var buttons []tgbotapi.InlineKeyboardButton

	for i := 0; i < len(names); i+=1 {
		buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(names[i], callbacks[i]))
	}

	return tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(buttons...))
}

// Функция для создания вертикальной клавиатуры
func VerticalDataInlineKeyboardMaker(names, callbacks []string) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton

	for i := 0; i < len(names); i+=1 {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(names[i], callbacks[i])))
	}

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

// Функция создания прайс-лсчта
func PriceListInlineKeyboardMarkup(names, prices, callbacks []string) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton

	for i := 0; i < len(names); i+=1 {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(names[i], "None"),
			tgbotapi.NewInlineKeyboardButtonData("В корзину за: "+prices[i]+"p", callbacks[i])),
		)
	}

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

// Функция создания отображения корзины
func CartInlineKeyboardMarkup(cart types.Cart) tgbotapi.InlineKeyboardMarkup {
	positionItems := cart.PositionItems
	var rows [][]tgbotapi.InlineKeyboardButton

	for k, v := range positionItems {
		for _, i := range v {
			rows = append(rows, tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(k+" "+i.Price+"p", "None"),
				tgbotapi.NewInlineKeyboardButtonData("Удалить из корзины", "delete"+i.Callback)),
			)
		}
	}

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}