package main

import (
	"errors"
	"log"
	envhandler "mashinki/envHandler"
	"mashinki/logging"
	"mashinki/tgBot"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Checking if bot token is available
	botToken := envhandler.GetEnv("TG_TOKEN")
	if botToken == "" {
		logging.DefaultLogger.LogError(errors.New("TG_TOKEN not found in environment variables"))
		return
	}

	log.Println("Starting bot...")
	bot, err := tgBot.StartBot()
	if err != nil {
		logging.DefaultLogger.LogErrorF("Failed to start bot: %v", err)
		return
	}

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	// wait for signal for shutdown
	<-quit
	log.Println("Shutting down bot...")
	bot.Stop()
}
