package main

import (
	"log"
	//"os"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type rootMenu struct {
	
}

func main() {
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

			switch update.CallbackQuery.Data {
			case "1":
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Pizza 1")
				bot.Send(msg)
			case "2":
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Pizza 2")
				bot.Send(msg)
			case "3":
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, "Pizza 3")
				bot.Send(msg)
			}

			/*callbackMsg := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
			_, err := bot.Request(callbackMsg)
			if err != nil {
				log.Fatal(err)
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
			msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("Pizza 1", "1"),
					tgbotapi.NewInlineKeyboardButtonData("Pizza 2", "2"),
					tgbotapi.NewInlineKeyboardButtonData("Pizza 3", "3"),
					tgbotapi.NewInlineKeyboardButtonURL("More pizzes:", "https://example.com"),
				),
			)
			bot.Send(msg)
		}
	}
}