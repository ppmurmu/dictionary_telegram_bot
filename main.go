package main

import (
	controllers "bot/models"
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// variables
var (
	bot                *tgbotapi.BotAPI
	failedMsg          = "Oops, we are not yet aware of this word!\nTry checking the spelling and try again."
	unexpectedErrorMsg = "Looks like the bot is unresponsive. Give us some time."
)

// main function
func main() {

	var err error
	//paste the bot key you got from botfather
	bot, err = tgbotapi.NewBotAPI("<YOUR_BOT_KEY_FROM_BOTFATHER>")
	if err != nil {
		// Abort if something is wrong
		log.Panic(err)
	}

	// Set this to true to log all interactions with telegram servers
	bot.Debug = false

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// Create a new cancellable background context. Calling `cancel()` leads to the cancellation of the context
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	// `updates` is a golang channel which receives telegram updates
	updates := bot.GetUpdatesChan(u)

	// Pass cancellable context to goroutine
	go receiveUpdates(ctx, updates)

	// Tell the user the bot is online
	log.Println("Start listening for updates. Press enter to stop")

	// Wait for a newline symbol, then cancel handling updates
	bufio.NewReader(os.Stdin).ReadBytes('\n')
	getRequest("stumble")
	cancel()

}

func receiveUpdates(ctx context.Context, updates tgbotapi.UpdatesChannel) {
	// `for {` means the loop is infinite until we manually stop it
	for {
		select {
		// stop looping if ctx is cancelled
		case <-ctx.Done():
			return
		// receive update from channel and then handle it
		case update := <-updates:
			handleMessage(update.Message)
		}
	}
}

func handleMessage(message *tgbotapi.Message) {
	user := message.From
	text := message.Text

	if user == nil {
		return
	}

	// Print to console
	log.Printf("%s wrote %s", user.FirstName, text)

	welcomeMsg := fmt.Sprintf("Welcome %s! I'm your dictionary bot.\nSimply send a word to get its meaning", user)
	var err error
	if strings.HasPrefix(text, "/start") {

		//err = handleCommand(message.Chat.ID, text)
		msg := tgbotapi.NewMessage(message.Chat.ID, welcomeMsg)
		// To preserve markdown, we attach entities (bold, italic..)
		msg.Entities = message.Entities
		_, err = bot.Send(msg)
	} else if len(text) > 0 {
		var response = getRequest(text)
		msg := tgbotapi.NewMessage(message.Chat.ID, response)
		// To preserve markdown, we attach entities (bold, italic..)
		msg.Entities = message.Entities
		_, err = bot.Send(msg)

	} else {
		msg := tgbotapi.NewMessage(message.Chat.ID, unexpectedErrorMsg)
		_, err = bot.Send(msg)
	}

	if err != nil {
		log.Printf("An error occured: %s", err.Error())
	}
}

// http client for persistent connection
func httpClient() *http.Client {
	client := &http.Client{Timeout: 10 * time.Second}
	return client
}

func getRequest(s string) string {
	var myUrl = fmt.Sprintf("https://api.dictionaryapi.dev/api/v2/entries/en/%s", s)

	client := httpClient()
	req, err := http.NewRequest(http.MethodGet, myUrl, nil)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Errored when sending request to the server")
		return unexpectedErrorMsg
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Errored when sending request to the server")
		return unexpectedErrorMsg
	}

	if resp.StatusCode == 200 {
		return onSuccess(responseBody, s)
	} else if resp.StatusCode == 404 {
		return failedMsg
	}

	return unexpectedErrorMsg

}

func onSuccess(data []byte, s string) string {
	var meaning controllers.Meaning
	err := json.Unmarshal(data, &meaning)
	if err != nil {
		log.Fatal(err)
	}
	var msg string = fmt.Sprintf("%s %s \n\n", s, meaning[0].Phonetic)
	for _, i := range meaning {
		for _, j := range i.Meanings {
			msg += fmt.Sprintf("- %s \n", j.PartOfSpeech)
			for l, m := range j.Definitions {
				msg += fmt.Sprintf("%d. %s \n", l+1, m.Definition)
			}
			if len(j.Synonyms) != 0 {
				msg += "synonyms: "
				for _, m := range j.Synonyms {
					msg += fmt.Sprintf("%s, ", m)
				}
				msg += "\n"
			}
			if len(j.Antonyms) != 0 {
				msg += "antonyms: "
				for _, m := range j.Antonyms {
					msg += fmt.Sprintf("%s, ", m)
				}
				msg += "\n"
			}
			msg += "\n"

		}
	}
	return msg

}
