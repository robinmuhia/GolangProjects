package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/shomali11/slacker"
)

func printCommandEvents(analyticsChannel <-chan *slacker.CommandEvent){
	for event := range analyticsChannel{
		fmt.Println("Command Events")
		fmt.Println(event.Timestamp)
		fmt.Println(event.Command)
		fmt.Println(event.Parameters)
		fmt.Println(event.Event)
		fmt.Println()
	}
}

func main(){
	envErr := godotenv.Load()
    if envErr != nil {
        log.Fatalf("err loading env: %v", envErr)
    }
	botToken := os.Getenv("SLACK_BOT_TOKEN")
	appToken := os.Getenv("SLACK_APP_TOKEN")
	bot := slacker.NewClient(botToken,appToken)

	go printCommandEvents(bot.CommandEvents())	
	var example [1] string;
	example[0] = "My year of birth is 2023"

	bot.Command("My year of birth is <years>", &slacker.CommandDefinition{
		Description:"Year of birth calculator",
		Examples: example[:],
		Handler: func(botCtx slacker.BotContext , request slacker.Request, response slacker.ResponseWriter){
			year := request.Param("years")
			yob,err := strconv.Atoi(year)
			if err != nil{
				fmt.Println("error")
			}
			age := 2023 - yob
			r := fmt.Sprintf("Age is %d",age)
			response.Reply(r)
		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err := bot.Listen(ctx)

	if err != nil{
		log.Fatal(err)
	}
}

