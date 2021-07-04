package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"sync"
	"syscall"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"

	"github.com/Machiel/telegrambot"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var (
	telegram        telegrambot.TelegramBot
	telegramBot     *tgbotapi.BotAPI
	db              *sqlx.DB
	emptyOpts       = telegrambot.SendMessageOptions{}
	commandHandlers = []commandHandler{
		commandDisablePictures,
		commandHelp,
		commandStart,
		commandStop,
		commandReport,
		commandSetup,
		commandMessage,
	}
	startJobs            = make(chan int64, 10000)
	updatesQueue         = make(chan *tgbotapi.Update, 10000)
	endConversationQueue = make(chan EndConversationEvent, 10000)
	updateMap            = NewIdRecorder()
	messageMap           = NewIdRecorder()
	stopped              = false
)

func main() {

	numCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPU)

	var err error

	log.Println("Starting...")
	mysqlUser := os.Getenv("MYSQL_USER")
	mysqlPassword := os.Getenv("MYSQL_PASSWORD")
	mysqlDatabaseName := os.Getenv("MYSQL_DATABASE_NAME")
	telegramBotKey := os.Getenv("TELEGRAM_BOT_KEY")

	dsn := fmt.Sprintf("%s:%s@(localhost:3306)/%s?parseTime=true", mysqlUser, mysqlPassword, mysqlDatabaseName)
	db, err = sqlx.Open("mysql", dsn)

	if err != nil {
		panic(err)
	}

	telegram = telegrambot.New(telegramBotKey)
	telegramBot, err = tgbotapi.NewBotAPI(telegramBotKey)
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func(jobs <-chan int64) {
		defer wg.Done()
		log.Println("Starting match user job")
		matchUsers(jobs)
	}(startJobs)

	for j := 0; j < 1; j++ {
		wg.Add(1)
		go func(jobs chan<- int64) {
			defer wg.Done()
			log.Println("Started load available user job")
			loadAvailableUsers(jobs)
		}(startJobs)
	}

	var workerWg sync.WaitGroup
	for i := 0; i < 3; i++ {
		workerWg.Add(1)
		go func(queue <-chan *tgbotapi.Update) {
			defer workerWg.Done()
			log.Println("Started a message worker...")
			updateWorker(queue)
		}(updatesQueue)
	}

	for x := 0; x < 1; x++ {
		wg.Add(1)

		go func(queue <-chan EndConversationEvent) {
			defer wg.Done()
			log.Println("Started end convo worker...")
			endConversationWorker(queue)
		}(endConversationQueue)
	}

	var receiverWg sync.WaitGroup
	receiverWg.Add(1)
	go func() {
		defer receiverWg.Done()
		log.Println("Started update worker")

			var offset int

		for {
			//log.Println("Requesting updates")
			offset = processUpdates(offset)
			//log.Println("Request completed")
			time.Sleep(1 * time.Second)

			if stopped {
				break
			}
		}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)
	signal.Notify(sigs, syscall.SIGTERM)
	done := make(chan bool, 1)

	go func() {
		<-sigs
		done <- true
	}()

	<-done

	log.Printf("Stopping...")

	stopped = true

	receiverWg.Wait()

	close(updatesQueue)

	workerWg.Wait()

	close(startJobs)
	close(endConversationQueue)

	log.Printf("Waiting for goroutines to stop...")

	wg.Wait()

	log.Printf("Closed...")
}

func loadAvailableUsers(startJobs chan<- int64) {

	for {

		u, err := retrieveAllAvailableUsers()

		if err != nil {
			log.Printf("Error retrieving everyone available: %s", err)
		} else {
			for _, x := range u {
				startJobs <- x.ChatID
			}
		}

		time.Sleep(10 * time.Second)
	}

}

// User holds user data
type User struct {
	ID            int64         `db:"id"`
	ChatID        int64         `db:"chat_id"`
	Available     bool          `db:"available"`
	LastActivity  time.Time     `db:"last_activity"`
	MatchChatID   sql.NullInt64 `db:"match_chat_id"`
	RegisterDate  time.Time     `db:"register_date"`
	PreviousMatch sql.NullInt64 `db:"previous_match"`
	AllowPictures bool          `db:"allow_pictures"`
	BannedUntil   NullTime      `db:"banned_until"`
	Gender        int           `db:"gender"`
	Tags          string        `db:"tags"`
	MatchMode     int           `db:"match_mode"`
}

func (u User) IsProfileFinish() bool {
	if u.Gender == 0 || u.MatchMode < 0 || u.Tags == "" {
		return false
	}
	return true
}

func (u User) GetNeedFinishProfile() string {

	field := make([]string, 0, 3)

	if u.Gender == 0 {
		field = append(field, "Gender")
	}
	if u.MatchMode < 0 {
		field = append(field, "Gender Preference")
	}
	if u.Tags == "" {
		field = append(field, "Goal")
	}

	return strings.Join(field, ",")
}

func retrieveUser(chatID int64) (User, error) {
	var u User
	err := db.Get(&u, "SELECT * FROM users WHERE chat_id = ?", chatID)
	return u, err
}

func retrieveOrCreateUser(chatID int64) (User, error) {
	var u User
	var count int
	err := db.Get(&count, "SELECT COUNT(*) FROM users WHERE chat_id = ?", chatID)

	if err != nil {
		return u, err
	}

	if count == 0 {
		_, err = db.Exec("INSERT INTO users(chat_id, available, last_activity, register_date, allow_pictures) VALUES (?, ?, ?, ?, 1)", chatID, false, time.Now(), time.Now())

		if err != nil {
			return u, err
		}

		telegram.SendMessage(chatID, `Welcome! This bot provides opportunities for interaction with your fellow Singaporeans! Especially in light of Covid-19 restrictions, perhaps such opportunities have reduced—and in some cases, have taken a toll on mental health. Whether you are here for a listening ear, to make a friend, somehow find a partner, find a lunch buddy or just have a quick wholesome chat, @sgchatterbot is for you. While conversations are anonymous, you're more than welcome to exchange contacts after a quick chat (do stay safe though)! Note that you're not 100% anonymous—if you break any rules or do anything illegal (especially noting there may be minors on here)—it will result in a permanent ban across all bots and a police report against you. This includes, but is not limited to impersonation, harassment, grooming, or sending of explicit images. If a user prefers to stay anonymous, this should be respected. Advertising or spam is not allowed as well.
		To configure your profile:
		/setup
		Then to start a conversation, enter
		/start
		If you feel like ending the conversation, type:
		/end
		If you want another chat partner, type /start again after typing /end!
		Enter /report (followed by reason) into the chat immediately to report abuse; if you are unable to do so because chat has ended, DO NOT start a new chat—immediately contact @aaldentnay . Note that users who abuse the platform will be banned permanently.
		E.g.
		/report creep
		(do not leave the reason blank! It must contain a reason for the report to send.)
		Enter /nopics to prevent others from sending pics to you!
		HEAD TO @singaporebotchannel for rules, announcements, etc. first before proceeding! Subscribe to this! :)
		Feel free to contact @aaldentnay if you need any assistance!
		Note that information of anyone who breaches rules will be tracked. Anything illegal will be reported to the police. Be careful with your personal information.
		The "friends/wholesome talk" segment should be strictly only for that; anyone "thirsty" should probably stick within the "tinder" option (even then, if you do anything out of line, you will still be banned/reported to authorities)
		Have fun!`, emptyOpts)
	}

	return retrieveUser(chatID)
}

func updateLastActivity(id int64) {
	db.Exec("UPDATE users SET last_activity = ? WHERE id = ?", time.Now(), id)
}

func updateGender(id int64, gender int) {
	_, err := db.Exec("UPDATE users SET gender = ? WHERE id = ?", gender, id)
	if err != nil {
		log.Println("update gender info error", err.Error())
	}
}

func updateMathMode(id int64, mathMode int) {
	_, err := db.Exec("UPDATE users SET match_mode = ? WHERE id = ?", mathMode, id)
	if err != nil {
		log.Println("update match mode error", err.Error())
	}
}

func updateTags(id int64, tags string) {
	_, err := db.Exec("UPDATE users SET tags = ? WHERE id = ?", tags, id)
	if err != nil {
		log.Println("update tags error", err.Error())
	}
}

func retrieveAllAvailableUsers() ([]User, error) {
	var u []User
	err := db.Select(&u, "SELECT * FROM users WHERE available = 1 AND match_chat_id IS NULL")
	return u, err
}

func retrieveAvailableUsers(c int64, user User) ([]User, error) {
	var u []User

	sql := `SELECT * FROM users WHERE (gender > 0 AND tags!="" AND match_mode > -1) AND chat_id != ? AND available = 1 AND match_chat_id IS NULL`

	switch user.MatchMode {
	case 1:
		sql = sql + fmt.Sprintf(" AND gender = 1 AND (match_mode = 0 OR match_mode = %d)", user.Gender)
	case 2:
		sql = sql + fmt.Sprintf(" AND gender = 2 AND (match_mode = 0 OR match_mode = %d)", user.Gender)
	default:
		sql = sql + fmt.Sprintf(" AND (match_mode = 0 OR match_mode = %d)", user.Gender)
	}

	if user.Tags != "" {
		sql = sql + ` AND tags = "` + user.Tags + `"`
	}

	err := db.Select(&u, sql, c)
	//err := db.Select(&u, "SELECT * FROM users WHERE chat_id != ? AND available = 1 AND match_chat_id IS NULL", c)
	return u, err
}

func shuffle(a []User) {
	for i := range a {
		j := rand.Intn(i + 1)
		a[i], a[j] = a[j], a[i]
	}
}

func handleMessage(message *tgbotapi.Message) {

	u, err := retrieveOrCreateUser(message.Chat.ID)

	if err != nil {
		log.Printf("retrieveOrCreateUser err: %s", err.Error())
		return
	}

	if u.BannedUntil.Valid && time.Now().Before(u.BannedUntil.Time) {
		date := u.BannedUntil.Time.Format("02 January 2006")
		response := fmt.Sprintf("You are banned until %s", date)
		_, err := telegram.SendMessage(message.Chat.ID, response, emptyOpts)
		if err != nil {
			log.Printf("handleMessage telegram.SendMessage err: %s", err.Error())
			return
		}
	}

	sendToHandler(u, message)

	updateLastActivity(u.ID)

}

func sendToHandler(u User, message *tgbotapi.Message) {

	log.Printf("msg_id: %d sendToHandler", message.MessageID)

	for _, handler := range commandHandlers {

		res := handler(u, message)
		if res {
			return
		}

	}

}

func processUpdates(offset int) int {

	log.Printf("Fetching with offset %d", offset)

	updates, err := telegramBot.GetUpdates(tgbotapi.UpdateConfig{
		Offset:  offset,
		Limit:   100,
		Timeout: 20,
	})

	if err != nil {
		log.Printf("GetUpdates err: %s", err.Error())
		return offset
	}

	return handleUpdates(updates, offset)

}

func handleUpdate(update *tgbotapi.Update) {

	if update.Message != nil {

		if !messageMap.IsSent(update.Message.MessageID) {
			if messageMap.SetSent(update.Message.MessageID) {
				handleMessage(update.Message)
			} else {
				log.Printf("message id:%d is handled", update.Message.MessageID)
			}

		} else {
			log.Printf("message id:%d is handled", update.Message.MessageID)
		}

	} else if update.CallbackQuery != nil {
		handleCallbackQuery(update.CallbackQuery)
	}

}

func updateWorker(updates <-chan *tgbotapi.Update) {

	for update := range updates {

		log.Printf("msg_id: %d updateWorker", update.Message.MessageID)

		if !updateMap.IsSent(update.UpdateID) {
			if updateMap.SetSent(update.UpdateID) {
				handleUpdate(update)
			} else {
				log.Printf("update id:%d is handled", update.UpdateID)
			}
		} else {
			log.Printf("update id:%d is handled", update.UpdateID)
		}
	}
}

func handleUpdates(updates []tgbotapi.Update, offset int) int {

	for i, update := range updates {

		log.Printf("offset: %d msg_id: %d msg_text: %s", update.UpdateID, update.Message.MessageID, update.Message.Text)

		if update.UpdateID >= offset {
			if update.UpdateID%1000 == 0 {
				log.Printf("Update ID: %d", update.UpdateID)
			}
			offset = update.UpdateID + 1
		}

		updatesQueue <- &updates[i]
	}

	return offset
}
