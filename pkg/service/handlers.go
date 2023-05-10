package service

import (
	"database/sql"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"math/rand"
	"strings"
	"time"
)

const Alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ_"
const AlphLen = 63

type Handler struct {
	Repo Repo
	Bot  *tgbotapi.BotAPI
}

func NewHandler(repo Repo) *Handler {
	return &Handler{Repo: repo}
}

type Repo interface {
	SetLogin(userID int64, serviceName string, pwd string) error
	GetLogin(userID int64, serviceName string) (string, error)
	Delete(userID int64, serviceName string) error
	Clear()
}

func (h *Handler) Commander(upd tgbotapi.Update, comm string) {
	switch comm {
	case "info":
		h.Start(upd)
	case "set":
		h.SaveLogin(upd)
	case "get":
		h.ShowPassword(upd)
	case "del":
		h.Clear(upd)
	default:
		h.SendMsg(upd, "Неизвестная команда\nНажмите /info")
		return
	}
}

func (h *Handler) SendMsg(upd tgbotapi.Update, msg string) {
	ans := tgbotapi.NewMessage(upd.Message.Chat.ID, msg)
	_, err := h.Bot.Send(ans)
	if err != nil {
		log.Print("Не получилось послать сообщение пользователю " + upd.Message.From.UserName)
	}
	return
}

func (h *Handler) SaveLogin(upd tgbotapi.Update) {
	userID := upd.Message.From.ID
	serviceName := upd.Message.CommandArguments()
	if serviceName == "" {
		msg := "Нужно было ввести название сервиса"
		h.SendMsg(upd, msg)
		return
	}
	rand.Seed(time.Now().UTC().UnixNano())
	pwd := CreatePassword()
	err := h.Repo.SetLogin(userID, serviceName, pwd)
	if err == nil {
		msg := fmt.Sprintf("Пароль для сервиса %s был сгенерирован ранее", serviceName)
		h.SendMsg(upd, msg)
		return

	}
	if strings.Contains(err.Error(), "ok") {
		msg := fmt.Sprintf("Созданный пароль для сервиса %s: %s", serviceName, pwd)
		h.SendMsg(upd, msg)
		return
	}
	msg := fmt.Sprintf("Fatal error\n press /info to proceed")
	h.SendMsg(upd, msg)
}

func (h *Handler) ShowPassword(upd tgbotapi.Update) {
	userID := upd.Message.From.ID
	serviceName := upd.Message.CommandArguments()
	if serviceName == "" {
		msg := "Нужно было ввести название сервиса"
		h.SendMsg(upd, msg)
		return
	}
	pwd, err := h.Repo.GetLogin(userID, serviceName)
	if err == nil {
		msg := fmt.Sprintf("Пароль для сервиса %s: %s", serviceName, pwd)
		h.SendMsg(upd, msg)
		return
	}
	if err == sql.ErrNoRows {
		msg := fmt.Sprintf("Пароль для сервиса %s: не генерировался или был удалён", serviceName)
		h.SendMsg(upd, msg)
		return
	}
	msg := fmt.Sprintf("Fatal error\n press /info to proceed")
	h.SendMsg(upd, msg)
}

func (h *Handler) Clear(upd tgbotapi.Update) {
	userID := upd.Message.From.ID
	serviceName := upd.Message.CommandArguments()
	if serviceName == "" {
		msg := "Нужно было ввести название сервиса"
		h.SendMsg(upd, msg)
		return
	}
	err := h.Repo.Delete(userID, serviceName)
	if err == nil {
		msg := fmt.Sprintf("Пароль для сервиса %s успешно удалён из базы", serviceName)
		h.SendMsg(upd, msg)
	}
	msg := fmt.Sprintf("Fatal error\n press /info to proceed")
	h.SendMsg(upd, msg)
}

func CreatePassword() string {
	pwd := make([]byte, 10)
	for i := range pwd {
		pwd[i] = Alphabet[rand.Intn(AlphLen)]
	}
	return string(pwd)
}

func (h *Handler) Start(upd tgbotapi.Update) {
	msg := "/set - добавляет логин и пароль к сервису (в качестве логина используется название сервиса)\n" +
		"/get - получает логин и пароль по названию сервиса\n" +
		"del - удаляет значения для сервиса\n" +
		"/info - вернуться к этому экрану"
	h.SendMsg(upd, msg)
}
