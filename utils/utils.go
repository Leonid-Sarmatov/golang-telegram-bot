package utils

import (
	"strings"
	"strconv"
	. "golang_telegram_bot/types"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Функция-обработчик для поиска нужной позиции заказа по калбеку, и добавления/удаления этой позиции
func AddOrRemovButtonCallbackHandler(cartMap *map[string]Cart, callbackQuery *tgbotapi.CallbackQuery, node Node) (bool, bool) {
	// ID пользователя
	id := strconv.FormatInt((*callbackQuery).Message.Chat.ID, 10)
	// Строка каллбека
	callback := (*callbackQuery).Data
	// Флаг успешного удаления элемента
	isRemoving := false
	// Флаг успешного нахождения калбека
	isCallbackFound := false

	// Если найден нужный калбек, то добавляем позицию в карзину 
	if node.GetCallback() == callback {
		// Создаем позицию для заказа 
		pos := EndingMenu {
			Name: node.GetName(),
			Price: node.GetPrice(),
			Callback: node.GetCallback(),
		}
		cart, ok := (*cartMap)[id];
		// Если корзина для пользователя уже создана, то обновляем ее
		if ok {
			cart.AddPositionItem(pos)
			(*cartMap)[id] = cart
		// Иначе создаем новую корзину и добавляем в нее позицию 
		} else {
			newCart := Cart {
				SumPrice: "0",
				PositionItems: make(map[string][]EndingMenu),
			}
			newCart.AddPositionItem(pos)
			(*cartMap)[id] = newCart
		}
		return false, true
	}

	// Если у калбека есть префикс "delete" - значит нужно не добавить, а удалить позицию
	if strings.HasPrefix(callback, "delete") && strings.HasSuffix(callback, node.GetCallback()) {
		// Создаем позицию для удаления 
		pos := EndingMenu {
			Name: node.GetName(),
			Price: node.GetPrice(),
			Callback: node.GetCallback(),
		}
		cart, ok := (*cartMap)[id];
		// Удаляем позицию из корзины
		if ok {
			cart.RemovePositionItem(pos)
			(*cartMap)[id] = cart
			return true, true
		} 
		return false, true
	}

	// Рекурсивно ищем калбек в дочерних узлах этого узла
	for _, v := range node.GetSubmenus() {
		isr, isf := AddOrRemovButtonCallbackHandler(cartMap, callbackQuery, v)
		isRemoving = isRemoving || isr
		isCallbackFound = isCallbackFound || isf
	}

	return isRemoving, isCallbackFound
}



// Функция поиска узла по калбеку
func Navigation(callback string, node Node) ([]string, []string, []string) {
	calbackArr := make([]string, 0)
	nameArr := make([]string, 0)
	priceArr := make([]string, 0)

	// Если нашли нужный калбек, то добавляем в списки его параметры
	if node.GetCallback() == callback {
		for _, v := range node.GetSubmenus() {
			calbackArr = append(calbackArr, v.GetCallback())
			nameArr = append(nameArr, v.GetName())
			if v.IsEnding() {
				priceArr = append(priceArr, v.GetPrice())
			} 
		}
		return nameArr, calbackArr, priceArr
	}

	// Рекурсивно ищем калбек в дочерних узлах этого узла
	for _, v := range node.GetSubmenus() {
		n, c, p := Navigation(callback, v)
		priceArr = append(priceArr, p...)
		calbackArr = append(calbackArr, c...)
		nameArr = append(nameArr, n...)
	}
	return nameArr, calbackArr, priceArr
}