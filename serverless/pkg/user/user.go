package user

import (
	"encoding/json"
	"errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/robinmuhia/GolangProjects/severless/pkg/validators"
)

var (
	ErrorFailedToFetchRecord="Failed to fetch record"
	ErrorFailedToUnmarshalRecord="Failed to unmarshal records into json"
	ErrorInvalidUserData="Invalid user data passed"
	ErrorInvalidEmail="Invalid Email passed"
	ErrorCouldNotMarshalItem="Could not marshal item"
	ErrorCouldNotDeleteItem="Could not delete item"
	ErrorCOuldNotPutItem="Could not update item"
	ErrorUserAlreadyExist="User already exists"
	ErrorUserDoesNotExist="User does not exist"
)

type User struct {
	Email     string `json:"email"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func FetchUser(email, tablename string, dynaClient dynamodbiface.DynamoDBAPI)(*User,error) {

	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"email":{
				S:aws.String(email),
			},
		},
		TableName: aws.String(tablename),
	}
	result,err := dynaClient.GetItem(input)
	if err != nil{
		return nil, errors.New(ErrorFailedToFetchRecord)
	}
	user := new(User)
	
	err = dynamodbattribute.UnmarshalMap(result.Item,user)
	if err != nil{
		return nil,errors.New(ErrorFailedToUnmarshalRecord)
	}
	return user,nil
}

func FetchUsers(tablename string, dynaClient dynamodbiface.DynamoDBAPI)(*[]User, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(tablename),
	}
	result,err := dynaClient.Scan(input)
	if err != nil{
		return nil,errors.New(ErrorFailedToFetchRecord)
	}
	users := new([]User)

	err = dynamodbattribute.UnmarshalListOfMaps(result.Items,users)

	if err != nil{
		return nil,errors.New(ErrorFailedToUnmarshalRecord)
	}

	return users,nil
}

func CreateUser(req events.APIGatewayProxyRequest, tablename string, dynaClient dynamodbiface.DynamoDBAPI)(*User,error){
	var user User
	if err := json.Unmarshal([]byte(req.Body),&user); err != nil{
		return nil, errors.New(ErrorInvalidUserData)
	}
	if !validators.IsEmailValid(user.Email){
		return nil,errors.New(ErrorInvalidEmail)
	}
	currentUser ,_ := FetchUser(user.Email,tablename,dynaClient)
	if currentUser != nil && len(currentUser.Email) != 0 {
		return nil, errors.New(ErrorUserAlreadyExist)
	}

	newUser, err := dynamodbattribute.MarshalMap(user)
	if err != nil{
		return nil,errors.New(ErrorCouldNotMarshalItem)
	}
	input := &dynamodb.PutItemInput{
		Item: newUser,
		TableName: aws.String(tablename),
	}
	_,err = dynaClient.PutItem(input)
	if err != nil{
		return nil, errors.New(ErrorCOuldNotPutItem)
	}
	return &user,nil
}

func UpdateUser(req events.APIGatewayProxyRequest, tablename string, dynaClient dynamodbiface.DynamoDBAPI)(*User ,error) {
	var user User
	if err := json.Unmarshal([]byte(req.Body),&user); err != nil{
		return nil, errors.New(ErrorInvalidUserData)
	}
	currentUser ,_ := FetchUser(user.Email,tablename,dynaClient)
	if currentUser != nil && len(currentUser.Email) != 0 {
		return nil, errors.New(ErrorUserAlreadyExist)
	}
	updateUser, err := dynamodbattribute.MarshalMap(user)
	if err != nil{
		return nil,errors.New(ErrorCouldNotMarshalItem)
	}
	input := &dynamodb.PutItemInput{
		Item: updateUser,
		TableName: aws.String(tablename),
	}
	_,err = dynaClient.PutItem(input)
	if err != nil{
		return nil, errors.New(ErrorCOuldNotPutItem)
	}
	return &user,nil
}

func DeleteUser(req events.APIGatewayProxyRequest, tablename string, dynaClient dynamodbiface.DynamoDBAPI) error {
	
	email := req.QueryStringParameters["email"]
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"email":{
				S: aws.String(email),
			},
		},
		TableName: aws.String(tablename),
	}
	_, err := dynaClient.DeleteItem(input)
	if err != nil {
		return errors.New(ErrorCouldNotDeleteItem)
	}

	return nil

}