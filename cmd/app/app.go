package main

import (
	"fmt"
	"log"
	"os"
	hexaratepaikamaco "travelWallet/internal/adapters/gateways/hexarate.paikama.co"
	"travelWallet/internal/adapters/handlers/telegram"
	"travelWallet/internal/adapters/repository/plaintext"
	"travelWallet/internal/app"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println(err.Error())
	}

	plaintextWallet, err := plaintext.NewWalletRepository()
	if err != nil {
		return
	}

	currencyConverter := hexaratepaikamaco.NewConverter()

	core := app.NewApplication(plaintextWallet, currencyConverter)

	tgbotapi, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_API_TOKEN"))
	if err != nil {
		fmt.Printf("Failed to create bot: %v", err)
		return
	}

	fmt.Println("Bot is starting...")
	tgBot := telegram.NewBot(tgbotapi, core)
	tgBot.Start()
}
