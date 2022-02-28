package main

import (
	_ "embed"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/Nv7-Github/bsharp/bot"
)

//go:embed token.txt
var token string

func main() {
	start := time.Now()
	fmt.Println("Loading Bot...")
	err := os.MkdirAll("data", os.ModePerm)
	if err != nil {
		panic(err)
	}
	bot, err := bot.NewBot("data", token, "947593278147162113", "903380812135825459")
	if err != nil {
		panic(err)
	}
	fmt.Println("Loaded bot in", time.Since(start))

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	fmt.Println("Press Ctrl+C to exit!")
	<-stop

	fmt.Println("Stopping...")
	bot.Close()

	/*
		TODO:
		- Comment/Image
		- Leaderboards
	*/
}
