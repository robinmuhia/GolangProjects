package main

import (
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/joho/godotenv"
	"github.com/robinmuhia/GolangProjects/severless/pkg/handlers"
)


var(
	dynaClient dynamodbiface.DynamoDBAPI
)

func main() {
	err := godotenv.Load()
	if err != nil {
	  log.Fatal("Error loading .env file")
	}
	region := os.Getenv("AWS_REGION")
	awsSession,err := session.NewSession(&aws.Config{
		Region:aws.String(region),
	})
	if err != nil{
		return
	}
	dynaClient = dynamodb.New(awsSession)
	lambda.Start(handler)
}

const tablename = "LambdaInGoUser"

func handler(req events.APIGatewayProxyRequest)(*events.APIGatewayProxyResponse,error){
	switch req.HTTPMethod{
	case "GET":
		return handlers.GetUser(req,tablename,dynaClient)
	case "POST":
		return handlers.CreateUser(req,tablename,dynaClient)
	case "PUT":
		return handlers.UpdateUser(req,tablename,dynaClient)
	case "DELETE":
		return handlers.DeleteUser(req,tablename,dynaClient)
	default:
		return handlers.UnhandledMethod()
	}
}