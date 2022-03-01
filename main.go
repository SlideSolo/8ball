// Package handler contains an HTTP Cloud Function to handle update from Telegram whenever a users interacts with the
// bot.
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

// Pass token and sensible APIs through environment variables
const telegramApiBaseUrl string = "https://api.telegram.org/bot"
const telegramApiSendMessage string = "/sendMessage"
const telegramTokenEnv string = "BOT_TOKEN"

var telegramApi string = telegramApiBaseUrl + os.Getenv(telegramTokenEnv) + telegramApiSendMessage

// Update is a Telegram object that we receive every time an user interacts with the bot.
type Update struct {
	UpdateId int     `json:"update_id"`
	Message  Message `json:"message"`
}

// Implements the fmt.String interface to get the representation of an Update as a string.
func (u Update) String() string {
	return fmt.Sprintf("(update id: %d, message: %s)", u.UpdateId, u.Message)
}

// Message is a Telegram object that can be found in an update.
// Note that not all Update contains a Message. Update for an Inline Query doesn't.
type Message struct {
	Text string `json:"text"`
	Chat Chat   `json:"chat"`
}

// Implements the fmt.String interface to get the representation of a Message as a string.
func (m Message) String() string {
	return fmt.Sprintf("(text: %s, chat: %s, audio %s)", m.Text, m.Chat)
}

// A Chat indicates the conversation to which the Message belongs.
type Chat struct {
	Id int `json:"id"`
}

// Implements the fmt.String interface to get the representation of a Chat as a string.
func (c Chat) String() string {
	return fmt.Sprintf("(id: %d)", c.Id)
}

// HandleTelegramWebHook sends a message back to the chat with a punchline starting by the message provided by the user.
func HandleTelegramWebHook(w http.ResponseWriter, r *http.Request) {
	// Parse incoming request
	var update, err = parseTelegramRequest(r)
	if err != nil {
		log.Printf("error parsing update, %s", err.Error())
		return
	}

	answers := getAnswers()
	rand.Seed(time.Now().UnixNano())
	phraseIndex := rand.Intn(19)
	reply := answers[phraseIndex]

	// Send the punchline back to Telegram
	var telegramResponseBody, errTelegram = sendTextToTelegramChat(update.Message.Chat.Id, reply)
	if errTelegram != nil {
		log.Printf("got error %s from telegram, response body is %s", errTelegram.Error(), telegramResponseBody)
	} else {
		log.Printf("punchline %s successfully distributed to chat id %d", reply, update.Message.Chat.Id)
	}
}

// parseTelegramRequest handles incoming update from the Telegram web hook
func parseTelegramRequest(r *http.Request) (*Update, error) {
	var update Update
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		log.Printf("could not decode incoming update %s", err.Error())
		return nil, err
	}
	if update.UpdateId == 0 {
		log.Printf("invalid update id, got update id = 0")
		return nil, errors.New("invalid update id of 0 indicates failure to parse incoming update")
	}
	return &update, nil
}

// sendTextToTelegramChat sends a text message to the Telegram chat identified by its chat Id
func sendTextToTelegramChat(chatId int, text string) (string, error) {

	log.Printf("Sending %s to chat_id: %d", text, chatId)
	response, err := http.PostForm(
		telegramApi,
		url.Values{
			"chat_id": {strconv.Itoa(chatId)},
			"text":    {text},
		})

	if err != nil {
		log.Printf("error when posting text to the chat: %s", err.Error())
		return "", err
	}
	defer response.Body.Close()

	var bodyBytes, errRead = ioutil.ReadAll(response.Body)
	if errRead != nil {
		log.Printf("error in parsing telegram answer %s", errRead.Error())
		return "", err
	}
	bodyString := string(bodyBytes)
	log.Printf("Body of Telegram Response: %s", bodyString)

	return bodyString, nil
}

func getAnswers() [20]string {
	var ans [20]string

	ans[0] = "Бесспорно"
	ans[1] = "Предрешено"
	ans[2] = "Никаких сомнений"
	ans[3] = "Определённо да"
	ans[4] = "Можешь быть уверен в этом"

	ans[5] = "Мне кажется — «да»"
	ans[6] = "Вероятнее всего"
	ans[7] = "Хорошие перспективы"
	ans[8] = "Знаки говорят — «да»"
	ans[9] = "Да"

	ans[10] = "Пока не ясно, попробуй снова"
	ans[11] = "Спроси позже"
	ans[12] = "Лучше не рассказывать"
	ans[13] = "Сейчас нельзя предсказать"
	ans[14] = "Сконцентрируйся и спроси опять"

	ans[15] = "Даже не думай"
	ans[16] = "Мой ответ — «нет»"
	ans[17] = "По моим данным — «нет»"
	ans[18] = "Перспективы не очень хорошие"
	ans[19] = "Весьма сомнительно"

	return ans
}
