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
			msg := tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, fmt.Sprintf("You selected %s as your Gender. %s", obj.GetOptionText(), obj.GetOptionNoteText()))
			telegramBot.Send(msg)
		}

		{
			msg := tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, `What gender do you want to match with?`)

			msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup([]tgbotapi.InlineKeyboardButton{
				{
					Text:         MatchModeOptionMaleText + MatchModeOptionMaleNoteText,
					CallbackData: &MatchModeMale,
				},
				{
					Text:         MatchModeOptionFemaleText + MatchModeOptionFemaleNoteText,
					CallbackData: &MatchModeFemale,
				},
				{
					Text:         MatchModeOptionAnythingText + MatchModeOptionAnythingNoteText,
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
			msg := tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, fmt.Sprintf("You selected %s as your Match. %s", obj.GetOptionText(), obj.GetOptionNoteText()))
			telegramBot.Send(msg)
		}

		{
			msg := tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, `What are you here for?`)

			msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup([]tgbotapi.InlineKeyboardButton{
				{
					Text:         GoalOptionDatingText + GoalOptionDatingNoteText,
					CallbackData: &GoalDating,
				},
				{
					Text:         GoalOptionFriendsText + GoalOptionFriendsNoteText,
					CallbackData: &GoalFriends,
				},
			})

			_, err := telegramBot.Send(msg)
			if err != nil {
				log.Println(err.Error())
			}
		}

	case GoalOptionType:

		// update tags
		updateTags(u.ID, obj.GetOptionText())

		// handle message
		{
			msg := tgbotapi.NewDeleteMessage(callbackQuery.Message.Chat.ID, callbackQuery.Message.MessageID)
			telegramBot.Send(msg)
		}

		{
			msg := tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, fmt.Sprintf("You selected %s as your Goal. %s", obj.GetOptionText(), obj.GetOptionNoteText()))
			telegramBot.Send(msg)
		}

	}

}
