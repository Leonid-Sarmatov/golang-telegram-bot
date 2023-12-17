package utils

import (
	"strings"
	"strconv"
	"fmt"
	. "golang_telegram_bot/types"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Функция-обработчик для поиска нужной позиции заказа по калбеку, и добавления этой позиции
func AddInCartButtonCallbackHandler(cartMap *map[string]Cart, callbackQuery *tgbotapi.CallbackQuery, node Node) error {
	// ID пользователя
	id := strconv.FormatInt((*callbackQuery).Message.Chat.ID, 10)
	// Строка каллбека
	callback := (*callbackQuery).Data
	// Флаг ошибки, поднимается когда найден каллбек. если False, значит калбек найден не был
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
		return nil
	}

	// Рекурсивно ищем калбек в дочерних узлах этого узла
	for _, v := range node.GetSubmenus() {
		err := AddInCartButtonCallbackHandler(cartMap, callbackQuery, v)
		if err == nil {
			isCallbackFound = isCallbackFound || true
		}
	}

	// Возвращаем ошибку, если калбек не найден
	if isCallbackFound == false {
		return fmt.Errorf("Callback not found")
	} else {
		return nil
	}
}

// Функция-обработчик для поиска нужной позиции заказа по калбеку, и удаления этой позиции
func RemovOutCartButtonCallbackHandler(cartMap *map[string]Cart, callbackQuery *tgbotapi.CallbackQuery, node Node) error {
	// ID пользователя
	id := strconv.FormatInt((*callbackQuery).Message.Chat.ID, 10)
	// Строка каллбека
	callback := (*callbackQuery).Data
	// Флаг успешного удаления элемента
	isRemoving := false
	// Флаг успешного нахождения калбека
	isCallbackFound := false

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
			return nil
		}
		 
		return fmt.Errorf("Cart not found for id")
	}

	// Рекурсивно ищем калбек в дочерних узлах этого узла
	for _, v := range node.GetSubmenus() {
		err := RemovOutCartButtonCallbackHandler(cartMap, callbackQuery, v)
		if err != nil {

		}
	}

	return fmt.Errorf
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