package handlers

import (
	"strconv"
	. "golang_telegram_bot/types"
	. "golang_telegram_bot/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func AllTextHandler(
	message *tgbotapi.Message, 
	nodeTree Node, 
	cartMap *map[string]Cart) tgbotapi.MessageConfig {

	// Переменная сообщения для отправки 
	var msg tgbotapi.MessageConfig

	switch message.Text {
	case "Посмотреть меню":
		msg = tgbotapi.NewMessage((*message).Chat.ID, "Выберите категорию продукта")
		msg.ReplyMarkup = VerticalDataInlineKeyboardMaker(
			[]string{ nodeTree.GetName() },
			[]string{ nodeTree.GetCallback() },
		)

	case "Моя корзина":
		cart, ok := (*cartMap)[strconv.FormatInt((*message).Chat.ID, 10)] 
		if ok && len(cart.PositionItems) != 0 {
			kb := CartInlineKeyboardMarkup(cart)
			msg = tgbotapi.NewMessage((*message).Chat.ID, "Список товаров в вашей корзине\nСуммарная стоимость: "+cart.SumPrice)
			msg.ReplyMarkup = kb
		} else {
			msg = tgbotapi.NewMessage((*message).Chat.ID, "Корзина пока пуста")
		}
	}
	return msg
}

func AllCommandsHandler(message *tgbotapi.Message) tgbotapi.MessageConfig {
	// Переменная сообщения для отправки 
	var msg tgbotapi.MessageConfig

	switch message.Command() {
	case "start":
		msg = tgbotapi.NewMessage(message.Chat.ID, "Привет! Я бот для продажи пиццы, в моем меню есть много вкусного)")
		msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("Посмотреть меню"),
				tgbotapi.NewKeyboardButton("Моя корзина"),
			),
		)
		return msg
	}

	return msg
}

func AllCallbackHandler( 
		cartMap *map[string]Cart,
		callbackQuery *tgbotapi.CallbackQuery,
		nodeTree Node) tgbotapi.MessageConfig {

	// Переменная сообщения для отправки 
	var msg tgbotapi.MessageConfig

	// Ищем нужный калбек в узлах дерева меню
	nameArr, callbackArr, priceArr := Navigation((*callbackQuery).Data, nodeTree)
	// Если он нашелся, отправляем пользователю нужный узел меню
	if len(nameArr) != 0 { 
		// Если список с ценами не нулевой, значит мы на узле, к которому прикреплены конечные узлы меню (позиции для заказа)
		if len(priceArr) != 0 {
			// Создаем клавиатуру для добавления позиции в заказ
			kb := PriceListInlineKeyboardMarkup(nameArr, priceArr, callbackArr)
			msg = tgbotapi.NewMessage((*callbackQuery).Message.Chat.ID, "Выберите товар")
			msg.ReplyMarkup = kb
		} else {
			// Создаем клавиатуру для дальнейшей навигации по меню
			kb := VerticalDataInlineKeyboardMaker(nameArr, callbackArr)
			msg = tgbotapi.NewMessage((*callbackQuery).Message.Chat.ID, "Выберите пункт меню")
			msg.ReplyMarkup = kb
		}
		return msg
	}

	err := AddInCartButtonCallbackHandler(
		cartMap,  
		callbackQuery,
		nodeTree,
	)

	if err == nil {
		msg = tgbotapi.NewMessage((*callbackQuery).Message.Chat.ID, "Позиция добавлена в корзину")
		return msg
	}

	err = RemovOutCartButtonCallbackHandler(
		cartMap,  
		callbackQuery,
		nodeTree,
	)

	if err == nil {
		msg = tgbotapi.NewMessage((*callbackQuery).Message.Chat.ID, "Позиция удалена из корзины")
		return msg
	}

	msg = tgbotapi.NewMessage((*callbackQuery).Message.Chat.ID, err.Error())
	return msg
	/*
	// Если калбек не найден в узлах, значит это мог быть калбек от кнопки добавления/удаления позиции
	isRemoving, isCallbackFound := AddOrRemovButtonCallbackHandler(
		cartMap,  
		callbackQuery,
		nodeTree,
	)
    // Если калбек действительно был от кнопки добавления/удаления то печатаем соответствующее сообщение
	if isCallbackFound {
		// Если калбек был удален, пишем сообщение об удалени позиции из корзины
		if isRemoving {
			msg = tgbotapi.NewMessage((*callbackQuery).Message.Chat.ID, "Позиция удалена из корзины")
		} 
		// Если калбек не был удален, пишем сообщение о добавлении позиции в корзину
		if !isRemoving {
			msg = tgbotapi.NewMessage((*callbackQuery).Message.Chat.ID, "Позиция добавлена в корзину")
		}
		return msg
	}
	return msg*/
}