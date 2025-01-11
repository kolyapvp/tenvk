package taskBot

import (
	"context"
	"log"
	"taskbot/internal/delivery/tg"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var (
	BotToken   = "76402885797pBL3x2UoNros"
	WebhookURL = "3123.12312.3213"
)

func startTaskBot(ctx context.Context, httpListenAddr string) error {
	return tg.StartTgBot(httpListenAddr, WebhookURL, BotToken)
}

func main() {
	err := startTaskBot(context.Background(), ":8081")
	if err != nil {
		log.Fatalln(err)
	}
}

// это заглушка чтобы импорт сохранился
func __dummy() {
	tgbotapi.APIEndpoint = "_dummy"
}
