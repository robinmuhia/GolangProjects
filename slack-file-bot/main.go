package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/slack-go/slack"
)

func main(){
	envErr := godotenv.Load()
    if envErr != nil {
        log.Fatalf("err loading env: %v", envErr)
    }
	slackBotToken := os.Getenv("SLACK_BOT_TOKEN")
	slackChannelId := os.Getenv("SLACK_CHANNEL_ID")
	slackApi := slack.New(slackBotToken)
	channelArr := []string{slackChannelId}
	fileArr := []string{"test.pdf"}

	for i := 0; i < len(fileArr); i++{
		params := slack.FileUploadParameters{
			Channels: channelArr,
			File: fileArr[i],
		}
		file, err := slackApi.UploadFile(params)
		if err != nil{
			fmt.Printf("%s\n",err)
			return
		}
		fmt.Printf("Name: %s , URL: %s\n",file.Name,file.URL)
	}
}