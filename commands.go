package main

import (
	"log"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// CommandHandler supplies an interface for handling messages
type commandHandler func(u User, m *tgbotapi.Message) bool

func commandDisablePictures(u User, m *tgbotapi.Message) bool {
	if len(m.Text) < 7 || strings.ToLower(m.Text[0:7]) != "/nopics" {
		return false
	}

	if u.AllowPictures {
		db.Exec("UPDATE users SET allow_pictures = 0 WHERE id = ?", u.ID)
		telegram.SendMessage(u.ChatID, "Users won't be able to send you photos anymore!", emptyOpts)
		return true
	}

	db.Exec("UPDATE users SET allow_pictures = 1 WHERE id = ?", u.ID)
	telegram.SendMessage(u.ChatID, "Users can now send you photos!", emptyOpts)
	return true
}

func commandStart(u User, m *tgbotapi.Message) bool {

	if len(m.Text) < 6 {
		return false
	}

	if strings.ToLower(m.Text[0:6]) != "/start" {
		return false
	}

	if u.Available {
		return false
	}

	if u.MatchChatID.Valid {
		return false
	}

	db.Exec("UPDATE users SET available = 1 WHERE chat_id = ?", u.ChatID)

	telegram.SendMessage(u.ChatID, "Looking for another student to match you with... Hold on! (This may take a while! Keep your notifications on!) **NOTE: If you send anything illegal here, your data will be handed over to the police. Your User ID is anonymous only until you break the rules. A police report for harassment/defamation will be filed if you pass off another user's contact as if it is yours.** To report a user, enter **/report (followed by a reason; don't leave blank)** into the chat. If chat with a user you want to report has already ended, then **do not** start a new chat—immediately contact the admin @aaldentnay . **Also, a note to some guys here—pls stop being thirsty on here because that scares new users away; I'm taking a huge leapt of faith when I set up a platform like this. Those reported will be PERMANENTLY banned.** Misuse of /report , if not accidental, can also result in ban.", emptyOpts)
	startJobs <- u.ChatID

	return true
}

func commandStop(u User, m *tgbotapi.Message) bool {

	if len(m.Text) < 4 {
		return false
	}

	rightCommand := strings.ToLower(m.Text[0:4]) == "/bye" || strings.ToLower(m.Text[0:4]) == "/end"

	if !rightCommand {
		return false
	}

	if !u.Available {
		return false
	}

	telegram.SendMessage(u.ChatID, "We're ending the conversation...", emptyOpts)

	endConversationQueue <- EndConversationEvent{ChatID: u.ChatID}

	return true
}

func commandReport(u User, m *tgbotapi.Message) bool {

	if len(m.Text) < 7 || strings.ToLower(m.Text[0:7]) != "/report" {
		return false
	}

	if !u.Available || !u.MatchChatID.Valid {
		return false
	}

	report := m.Text[7:]
	report = strings.TrimSpace(report)

	if len(report) == 0 {
		telegram.SendMessage(u.ChatID, "Usage /report: /report <reason>", emptyOpts)
		return true
	}

	partner, err := retrieveUser(u.MatchChatID.Int64)

	if err != nil {
		log.Println("Error retrieving partner")
		return true
	}

	db.Exec("INSERT INTO reports (user_id, report, reporter_id, created_at) VALUES (?, ?, ?, ?)", partner.ID, report, u.ID, time.Now())

	telegram.SendMessage(u.ChatID, "User has been reported!", emptyOpts)

	return true
}

func commandMessage(u User, m *tgbotapi.Message) bool {

	if !u.Available {
		return false
	}

	if !u.MatchChatID.Valid {
		return false
	}

	chatID := u.MatchChatID.Int64
	partner, err := retrieveUser(chatID)

	if err != nil {
		log.Println("[ERROR] Could not retrieve partner %d", chatID)
		return false
	}

	if m.Photo != nil && len(*m.Photo) > 0 {

		if !partner.AllowPictures {
			telegram.SendMessage(chatID, "User tried to send you a photo, but you disabled this,  you can enable photos by using the /nopics command", emptyOpts)
			telegram.SendMessage(u.ChatID, "User disabled photos, and will not receive your photos", emptyOpts)
			return true
		}

		var toSend tgbotapi.PhotoSize

		for _, t := range *m.Photo {
			if t.FileSize > toSend.FileSize {
				toSend = t
			}
		}

		telegram.SendMessage(chatID, "User sends you a photo!", emptyOpts)
		_, err = telegram.SendPhoto(chatID, toSend.FileID, emptyOpts)

	} else if m.Sticker != nil {
		telegram.SendMessage(chatID, "User sends you a sticker!", emptyOpts)
		_, err = telegram.SendSticker(chatID, m.Sticker.FileID, emptyOpts)
	} else if m.Location != nil {
		telegram.SendMessage(chatID, "User sends you a location!", emptyOpts)
		_, err = telegram.SendLocation(chatID,
			m.Location.Latitude,
			m.Location.Longitude,
			emptyOpts,
		)
	} else if m.Document != nil {
		telegram.SendMessage(chatID, "User sends you a document!", emptyOpts)
		_, err = telegram.SendDocument(chatID, m.Document.FileID, emptyOpts)
	} else if m.Audio != nil {
		telegram.SendMessage(chatID, "User sends you an audio file!", emptyOpts)
		_, err = telegram.SendAudio(chatID, m.Audio.FileID, emptyOpts)
	} else if m.Video != nil {
		telegram.SendMessage(chatID, "User sends you a video file!", emptyOpts)
		_, err = telegram.SendVideo(chatID, m.Video.FileID, emptyOpts)
	} else {
		_, err = telegram.SendMessage(chatID, "User: "+m.Text, emptyOpts)
	}

	if err != nil {
		log.Printf("Forward error: %s", err)
	}

	return true

}

func commandHelp(u User, m *tgbotapi.Message) bool {

	if len(m.Text) < 5 {
		return false
	}

	if strings.ToLower(m.Text[0:5]) != "/help" {
		return false
	}

	telegram.SendMessage(m.Chat.ID, `Help:

Use /start to start looking for a conversational partner, once you're matched you can use /end to end the conversation.

Use /report to report a user, use it as follows:
/report <reason>

Use /nopics to disable receiving photos, and /nopics if you want to enable it again.

HEAD OVER to @unichatbotchannel for rules, updates, announcements or info on how you can support the bot!

Sending images and videos are a beta functionality, but appear to be working fine.

If you require any help, feel free to contact @aaldentnay !`, emptyOpts)

	return true
}

func commandSetup(u User, m *tgbotapi.Message) bool {

	if len(m.Text) < 6 {
		return false
	}

	if strings.ToLower(m.Text[0:6]) != "/setup" {
		return false
	}

	msg := tgbotapi.NewMessage(m.Chat.ID, `What is your gender?`)

	msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup([]tgbotapi.InlineKeyboardButton{
		{
			Text:         "Male",
			CallbackData: &GenderMale,
		},
		{
			Text:         "Female",
			CallbackData: &GenderFemale,
		},
	})

	_, err := telegramBot.Send(msg)
	if err != nil {
		log.Println(err.Error())
	}

	return true
}
