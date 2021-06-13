package main

import (
	"encoding/json"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func handleCallbackQuery(callbackQuery *tgbotapi.CallbackQuery) {

	if len(callbackQuery.Data) == 0 {
		return
	}

	u, err := retrieveOrCreateUser(callbackQuery.Message.Chat.ID)

	if err != nil {
		log.Println(err)
		return
	}

	obj := new(KeyboardCallbackData)
	if err := json.Unmarshal([]byte(callbackQuery.Data), &obj); err != nil {
		log.Println("json unamrshal error", err.Error())
	}

	switch obj.OptionType {
	case GenderOptionType:

		// update gender
		updateGender(u.ID, obj.OptionValue)

		// handle message
		{
			msg := tgbotapi.NewDeleteMessage(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID)
			telegramBot.Send(msg)
		}

		{
			msg := tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, fmt.Sprintf("You selected %s as your Gender", obj.GetOptionText()))
			telegramBot.Send(msg)
		}

		{
			msg := tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, `What gender do you want to match with?`)

			msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup([]tgbotapi.InlineKeyboardButton{
				{
					Text:         MatchModeOptionMaleText,
					CallbackData: &MatchModeMale,
				},
				{
					Text:         MatchModeOptionFemaleText,
					CallbackData: &MatchModeFemale,
				},
				{
					Text:         MatchModeOptionAnythingText,
					CallbackData: &MatchModeAnything,
				},
			})

			_, err := telegramBot.Send(msg)
			if err != nil {
				log.Println(err.Error())
			}
		}

	case MatchModeOptionType:

		// update gender
		updateMathMode(u.ID, obj.OptionValue)

		// handle message
		{
			msg := tgbotapi.NewDeleteMessage(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID)
			telegramBot.Send(msg)
		}

		{
			msg := tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, fmt.Sprintf("You selected %s as your Match", obj.GetOptionText()))
			telegramBot.Send(msg)
		}

		{
			msg := tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, `What are you here for?`)

			msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup([]tgbotapi.InlineKeyboardButton{
				{
					Text:         GoalOptionDatingText,
					CallbackData: &GoalDating,
				},
				{
					Text:         GoalOptionFriendsText,
					CallbackData: &GoalFriends,
				},
			})

			_, err := telegramBot.Send(msg)
			if err != nil {
				log.Println(err.Error())
			}
		}

	case GoalOptionType:

		// update gender
		updateTags(u.ID, obj.GetOptionText())

		// handle message
		{
			msg := tgbotapi.NewDeleteMessage(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID)
			telegramBot.Send(msg)
		}

		{
			msg := tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, fmt.Sprintf("You selected %s as your Goal", obj.GetOptionText()))
			telegramBot.Send(msg)
		}

	}

}
