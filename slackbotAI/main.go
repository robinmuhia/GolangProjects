package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/krognol/go-wolfram"
	"github.com/shomali11/slacker"
	"github.com/tidwall/gjson"
	witai "github.com/wit-ai/wit-go/v2"
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

func main() {
	err := godotenv.Load()
	if err != nil{
		log.Fatal("Error loading .env file")
	}

	slackBotToken := os.Getenv("SLACK_BOT_TOKEN")
	slackAppToken := os.Getenv("SLACK_SOCKET_TOKEN")
	witaiServerToken := os.Getenv("WITAI_SERVER_TOKEN")
	wolframAppId := os.Getenv("WOLFRAM_APP_ID")

	slackBot := slacker.NewClient(slackBotToken,slackAppToken)
	witaiClient := witai.NewClient(witaiServerToken)
	wolframClient := &wolfram.Client{
		AppID: wolframAppId,
	}

	go printCommandEvents(slackBot.CommandEvents())

	var example [1] string;
	example[0] = "My year of birth is 2023"

	slackBot.Command("Query for bot - <message>",&slacker.CommandDefinition{
		Description: "Send any question to wolfram",
		Examples: example[:],
		Handler: func(bc slacker.BotContext, r slacker.Request, w slacker.ResponseWriter) {
			query := r.Param("message")
			message,_ := witaiClient.Parse(&witai.MessageRequest{
				Query: query,
			})
			data,_ := json.MarshalIndent(message, "","    ")
			rough := string(data[:])
			value := gjson.Get(rough,"entities.wit$wolfram_search_query:wolfram_search_query.0.value")
			question := value.String()
			answer,err := wolframClient.GetSpokentAnswerQuery(question,wolfram.Metric,1000)
			if err != nil{
				fmt.Printf("Err:%v",err)
			}
			w.Reply(answer)
		},
	})

	ctx,cancel := context.WithCancel(context.Background())
	defer cancel()

	err = slackBot.Listen(ctx)
	
	if err != nil{
		log.Fatal(err)
	}
}