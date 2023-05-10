package main

import (
	"PasswordBot/pkg/repository"
	"PasswordBot/pkg/service"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/robfig/cron"
	"log"
	"net/http"
	"os"
)

const (
	BotToken   = "6009569817:AAGp-5LQ3Wft436L_ybiEtaotPVM0RHnoGY"
	WebhookURL = "https://24c0-91-188-188-211.ngrok-free.app"
)

func main() {
	db := repository.NewRepo()
	cron := cron.New()
	bot, err := tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		log.Fatalf("New bot API failed: %s", err)
	}
	bot.Debug = true
	fmt.Printf("Authorized on account %s\n", bot.Self.UserName)
	//wh, err := tgbotapi.NewWebhook(WebhookURL)
	//if err != nil {
	//	log.Fatalf("NewWebhook failed: %s", err)
	//}
	//_, err = bot.Request(wh)
	//if err != nil {
	//	log.Fatalf("Setting Webhook failed: %s", err)
	//}
	updates := bot.ListenForWebhook("/")
	http.HandleFunc("/state", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("all is working"))
		if err != nil {
			return
		}
	})
	port := os.Getenv("PORT")
	fmt.Println(port)
	if port == "" {
		port = "8081"
	}
	go func() {
		log.Fatalln("http err:", http.ListenAndServe(":"+port, nil))
	}()
	fmt.Println("start listening :" + port)

	worker := service.NewHandler(db)
	cron.AddFunc("@daily", func() {
		worker.Repo.Clear()
	})
	cron.Start()
	for upd := range updates {
		log.Printf("upd: %#v\n", upd)
		if upd.Message == nil {
			log.Println("Change of message")
			continue
		}
		if upd.Message.From == nil {
			log.Println("nil User")
			continue
		}
		command := upd.Message.Command()
		if command == "" {
			msg := tgbotapi.NewMessage(upd.Message.Chat.ID, "That's not a command")
			_, errSend := bot.Send(msg)
			if errSend != nil {
				log.Printf("Sending failed: %s", errSend)
			}
			continue
		}
		worker.Commander(upd, command)
	}
}
