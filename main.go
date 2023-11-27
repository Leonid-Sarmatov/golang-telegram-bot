package main

import (
	"log"
	//"os"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Node interface {
	isEnding() bool
	getName() string
	getCallback() string
	getSubmenus() []Node
}

type NodeMenu struct {
	Name string
	Callback string
	SubMenu []Node
}

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

func main() {
    e1 := EndingMenu { Name: "Пеперони", Price: "100p", Callback: "peperoni" }
	e2 := EndingMenu { Name: "Четыре сыра", Price: "100p", Callback: "4_chees" }
	e3 := EndingMenu { Name: "Маргарита", Price: "100p", Callback: "margarita" }

	e4 := EndingMenu { Name: "Капучино", Callback: "capuchino" }
	e5 := EndingMenu { Name: "Латте", Callback: "latte" }

	n1 := NodeMenu { Name: "Пицца", Callback: "pizza", SubMenu: []Node{ e1, e2, e3 }}
	n2 := NodeMenu { Name: "Напитки", Callback: "drink", SubMenu: []Node{ e4, e5 }}

	n0 := NodeMenu { Name: "Меню", Callback: "menu", SubMenu: []Node{ n1, n2 }}

	//token := os.Getenv("6058196438:AAH2svI0pJAcJ592nIojO1yuv43JwFwRlu4")
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
			nameArr, callbackArr := Navigation(update.CallbackQuery.Data, n0)
			if len(nameArr) != 0 {
				kb := VerticalDataInlineKeyboardMaker(nameArr, callbackArr)
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Выберите пункт меню")
				msg.ReplyMarkup = kb
				bot.Send(msg)
			}
			/*switch update.CallbackQuery.Data {
			case "1":
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Pizza 1")
				bot.Send(msg)
			case "2":
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Pizza 2")
				bot.Send(msg)
			case "3":
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Pizza 3")
				bot.Send(msg)
			}*/
		}

		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "start":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Hi, i'm pizza-maker bot. See my menu")
				msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
					tgbotapi.NewKeyboardButtonRow(
						tgbotapi.NewKeyboardButton("See menu"),
					),
				)
				bot.Send(msg)
			}
		}

		if update.Message.Text == "See menu" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "This is my menu.")
			msg.ReplyMarkup = VerticalDataInlineKeyboardMaker(
				[]string{ n0.Name },
				[]string{ n0.Callback },
			)
			/*msg.ReplyMarkup = VerticalDataInlineKeyboardMaker(
				[]string{"Pizza 1", "Pizza 2", "Pizza 3"},
				[]string{"1", "2", "3"},
			)*/
			bot.Send(msg)
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



func Navigation(callback string, node Node) ([]string, []string) {
	if node.isEnding() {
		return make([]string, 0), make([]string, 0)
	}

	if node.getCallback() == callback {
		calls := make([]string, 0)
		names := make([]string, 0)
		for _, v := range node.getSubmenus() {
			calls = append(calls, v.getCallback())
			names = append(names, v.getName())
		}
		return names, calls
	}

	calbackArr := make([]string, 0)
	nameArr := make([]string, 0)
	for _, v := range node.getSubmenus() {
		n, c := Navigation(callback, v)
		calbackArr = append(calbackArr, c...)
		nameArr = append(nameArr, n...)
	}
	return nameArr, calbackArr
}