package main

import (
	"log"
	"strconv"
	"fmt"
	//"os"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)



// Структура описывающая корзину заказа
type Position struct {
	Name string
	Price string
}
type Cart struct {
	SumPrice int
	PositionItems map[string]Position
}

// Добавления в заказ позиции 
func (c *Cart) AddOrderItems(p Position) error {
	n, err := strconv.Atoi(p.Price)
	if err != nil {
		return fmt.Errorf("ERROR: Price value is not integer. Order: %s\n", p.Name)
	}

	c.SumPrice += n
	if c.PositionItems == nil {
		c.PositionItems = make(map[string]Position)
	}
	c.PositionItems[p.Name] = p
	return nil
}

// Удаление позиции из заказа
func (c *Cart) RemoveOrderItems(p Position) error {
	n, err := strconv.Atoi(p.Price)
	if err != nil {
		return fmt.Errorf("ERROR: Price value is not integer. Order: %s\n", p.Name)
	}

	c.SumPrice -= n
	delete(c.PositionItems, p.Name)
	return nil
}

// Интерфейс, описывающий обстрактный узел дерева меню
type Node interface {
	isEnding() bool
	getName() string
	getCallback() string
	getSubmenus() []Node
	getPrice() string
}

// Структура, описывающая узел дерева меню, не являющийся конечным
type NodeMenu struct {
	Name string
	Callback string
	SubMenu []Node
}

// Структура, описывающая конечные узлы дерева, то есть позиции для заказа
type EndingMenu struct {
	Name string
	Price string
	Callback string
}



func (m NodeMenu) isEnding() bool {
	return false
}

func (m NodeMenu) getName() string {
	return m.Name
}

func (m NodeMenu) getCallback() string {
	return m.Callback
}

func (m NodeMenu) getSubmenus() []Node {
	return m.SubMenu
}

func (m NodeMenu) getPrice() string {
	return "None"
}



func (e EndingMenu) isEnding() bool {
	return true
}

func (e EndingMenu) getName() string {
	return e.Name
}

func (e EndingMenu) getCallback() string {
	return e.Callback
}

func (e EndingMenu) getSubmenus() []Node {
	return nil
}

func (e EndingMenu) getPrice() string {
	return e.Price
}

func main() {
	// ***** Начало создания структуры меню *****
    e1 := EndingMenu { Name: "Пеперони", Price: "100", Callback: "peperoni" }
	e2 := EndingMenu { Name: "Четыре сыра", Price: "120", Callback: "4_chees" }
	e3 := EndingMenu { Name: "Маргарита", Price: "90", Callback: "margarita" }

	e4 := EndingMenu { Name: "Капучино", Price: "150", Callback: "capuchino" }
	e5 := EndingMenu { Name: "Латте", Price: "150", Callback: "latte" }

	e6 := EndingMenu { Name: "Сендвич с ветчиной", Price: "120", Callback: "vetchina" }
	e7 := EndingMenu { Name: "Сендвич с курицей", Price: "130", Callback: "kurica" }

	n1 := NodeMenu { Name: "Пицца", Callback: "pizza", SubMenu: []Node{ e1, e2, e3 }}
	n2 := NodeMenu { Name: "Напитки", Callback: "drink", SubMenu: []Node{ e4, e5 }}
	n3 := NodeMenu { Name: "Сендвичи", Callback: "sandvuch", SubMenu: []Node{ e6, e7 }}

	n0 := NodeMenu { Name: "Меню", Callback: "menu", SubMenu: []Node{ n1, n2, n3 }}
	// ***** Конец создания структуры меню *****

	cartMap := make(map[string]Cart)

	token := "6058196438:AAH2svI0pJAcJ592nIojO1yuv43JwFwRlu4"

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true

	log.Printf("Autorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {

		if update.CallbackQuery != nil {
		    // Ищем нужный калбек в дереве меню
			nameArr, callbackArr, priceArr := Navigation(update.CallbackQuery.Data, n0)
			// Если он нашелся, отправляем пользователю нужный узел меню
			if len(nameArr) != 0 { 
				// Если список с уенами не нулевой, значит мы на конце меню, то есть на позициях для заказа
				if len(priceArr) != 0 {
					// Создаем клавиатуру для добавления позиции в заказ
					kb := PriceListInlineKeyboardMarkup(nameArr, priceArr, callbackArr)
					msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Выберите товар")
					msg.ReplyMarkup = kb
					bot.Send(msg)
				} else {
					// Создаем клавиатуру для дальнейшей навигации по меню
					kb := VerticalDataInlineKeyboardMaker(nameArr, callbackArr)
					msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Выберите пункт меню")
					msg.ReplyMarkup = kb
					bot.Send(msg)
				}
			// Если калбек не найден, значит это был калбек от кнопки добавления позиции в заказ
			} else {
				AddOrderButtonCallbackHandler(
					&cartMap, 
					strconv.FormatInt(update.CallbackQuery.Message.Chat.ID, 10), 
					update.CallbackQuery.Data,
					n0,
				)
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Пункт добавлен в корзину")
				bot.Send(msg)
			}
			
		}

		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Привет! Я бот для продажи пиццы, в моем меню есть много вкусного)")
				msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
					tgbotapi.NewKeyboardButtonRow(
						tgbotapi.NewKeyboardButton("Посмотреть меню"),
						tgbotapi.NewKeyboardButton("Моя корзина"),
					),
				)
				bot.Send(msg)
			}
		}

		if update.Message.Text == "Посмотреть меню" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите категорию продукта")
			msg.ReplyMarkup = VerticalDataInlineKeyboardMaker(
				[]string{ n0.Name },
				[]string{ n0.Callback },
			)
			bot.Send(msg)
		}

		if update.Message.Text == "Моя корзина" {
			cart, ok := cartMap[strconv.FormatInt(update.Message.Chat.ID, 10)] 
			if ok {
				kb := CartInlineKeyboardMarkup(cart.PositionItems)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Список товаров в вашей корзине")
				msg.ReplyMarkup = kb
				bot.Send(msg)
			} else {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Корзина пока пуста")
				bot.Send(msg)
			}
		}
	}
}



func HorizontalDataInlineKeyboardMaker(names, callbacks []string) tgbotapi.InlineKeyboardMarkup {
	var buttons []tgbotapi.InlineKeyboardButton

	for i := 0; i < len(names); i+=1 {
		buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(names[i], callbacks[i]))
	}

	return tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(buttons...))
}

func VerticalDataInlineKeyboardMaker(names, callbacks []string) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton

	for i := 0; i < len(names); i+=1 {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(names[i], callbacks[i])))
	}

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

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

func CartInlineKeyboardMarkup(positionItems map[string]Position) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton

	for k, v := range positionItems {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(k+" "+v.Price+"p", "None"),
			tgbotapi.NewInlineKeyboardButtonData("Удалить из корзины", "None")),
		)
	}

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}



// Функция-обработчик для поиска нужной позиции заказа по калбеку, и добавления этой позиции в заказ
func AddOrderButtonCallbackHandler(cartMap *map[string]Cart, userId, callback string, node Node) error {
	if node.getCallback() == callback {
		// Создаем позицию для заказа
		p := Position {
			Name: node.getName(),
			Price: node.getPrice(),
		}
		// Если корзина для пользователя уже создана, то обновляем ее
		cart, ok := (*cartMap)[userId];
		if ok {
			cart.AddOrderItems(p)
			(*cartMap)[userId] = cart
		// Иначе создаем новую корзину и добавляем в нее позицию 
		} else {
			var newCart Cart
			newCart.AddOrderItems(p)
			(*cartMap)[userId] = newCart
		}
		return nil
	}

	// Рекурсивно ищем калбек в дочерних узлах этого узла
	for _, v := range node.getSubmenus() {
		AddOrderButtonCallbackHandler(cartMap, userId, callback, v)
	}

	if _, ok := (*cartMap)[userId]; ok {
		return nil
	}

	return fmt.Errorf("ERROR: Can not add position in CartMap. ID: %s, Node: %s\n", userId, node)
}



// Функция поиска узла по калбеку
func Navigation(callback string, node Node) ([]string, []string, []string) {

	// Если нашли нужный калбек, то добавляем в списки его параметры
	if node.getCallback() == callback {
		calls := make([]string, 0)
		names := make([]string, 0)
		prices := make([]string, 0)

		for _, v := range node.getSubmenus() {
			calls = append(calls, v.getCallback())
			names = append(names, v.getName())
			if v.isEnding() {
				prices = append(prices, v.getPrice())
			} 
		}
		fmt.Println("==========")
		fmt.Println("calls: ", calls, "  names: ", names, "  prices: ", prices, " node: ", node)
		fmt.Println("==========")
		return names, calls, prices
	}

	calbackArr := make([]string, 0)
	nameArr := make([]string, 0)
	priceArr := make([]string, 0)

	// Рекурсивно ищем калбек в дочерних узлах этого узла
	for _, v := range node.getSubmenus() {
		n, c, p := Navigation(callback, v)
		priceArr = append(priceArr, p...)
		calbackArr = append(calbackArr, c...)
		nameArr = append(nameArr, n...)
	}
	return nameArr, calbackArr, priceArr
}