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

const guild = "903380812135825459"

func main() {
	start := time.Now()
	fmt.Println("Loading Bot...")
	bot, err := bot.NewBot("data", token)
	if err != nil {
		panic(err)
	}
	fmt.Println("Loaded bot in", time.Since(start))

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	fmt.Println("Press Ctrl+C to exit!")
	<-stop

	bot.Close()
}
