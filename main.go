package main

import (
	"log"
	//"fmt"
	. "golang_telegram_bot/types"
	. "golang_telegram_bot/handlers"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var (
	token = "6058196438:AAH2svI0pJAcJ592nIojO1yuv43JwFwRlu4"
	cartMap = make(map[string]Cart)
)


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
			bot.Send(AllCallbackHandler(&cartMap, update.CallbackQuery, n0))
		}

		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			bot.Send(AllCommandsHandler(update.Message))
		}

		if update.Message.Text != "" {
			bot.Send(AllTextHandler(update.Message, &n0, &cartMap))
		}
	}
}