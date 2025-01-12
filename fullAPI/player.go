package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/arthurmvo/lambdahandler"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// func createPlayer(request events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
func createPlayer(ctx context.Context, request events.LambdaFunctionURLRequest, params lambdahandler.Params) (interface{}, lambdahandler.LambdaError) {
	// Unmarshal the request body into a Player struct
	var newPlayer Player
	err := json.Unmarshal([]byte(request.Body), &newPlayer)
	if err != nil {
		return nil, lambdahandler.NewLambdaError(400, fmt.Sprintf("Error unmarshalling request body: %s", err.Error()))
	}

	fmt.Println(newPlayer, "newPlayer")

	// Marshal the Player struct into an AttributeValue
	av, err := dynamodbattribute.MarshalMap(newPlayer)
	if err != nil {
		return nil, lambdahandler.NewLambdaError(500, fmt.Sprintf("Error marshalling player: %s", err.Error()))
	}

	// Put the new player into DynamoDB
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)

	input := &dynamodb.PutItemInput{
		TableName: aws.String("players"), // Replace with your table name
		Item:      av,
	}

	_, err = svc.PutItem(input)
	if err != nil {
		return nil, lambdahandler.NewLambdaError(500, fmt.Sprintf("Error putting item into DynamoDB: %s", err.Error()))
	}

	return newPlayer, nil
}

// func updatePlayer(request events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
func updatePlayer(ctx context.Context, request events.LambdaFunctionURLRequest, params lambdahandler.Params) (interface{}, lambdahandler.LambdaError) {

	uid, err := strconv.Atoi(params["uid"])
	if err != nil {
		return nil, lambdahandler.NewLambdaError(400, fmt.Sprintf("Invalid ID format: %s", err.Error()))
	}

	// Unmarshal the request body into a Player struct
	var updatedPlayer Player
	err = json.Unmarshal([]byte(request.Body), &updatedPlayer)
	if err != nil {
		return nil, lambdahandler.NewLambdaError(400, fmt.Sprintf("Error unmarshalling request body: %s", err.Error()))
	}

	// Ensure the ID in the request body matches the ID in the path
	if updatedPlayer.UID != uid {
		return nil, lambdahandler.NewLambdaError(400, "ID in the request body does not match ID in the path")
	}

	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)

	av, err := dynamodbattribute.MarshalMap(updatedPlayer)
	if err != nil {
		return nil, lambdahandler.NewLambdaError(500, fmt.Sprintf("Error marshalling player: %s", err.Error()))
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String("players"), // Replace with your table name
		Item:      av,
	}

	_, err = svc.PutItem(input)
	if err != nil {
		return nil, lambdahandler.NewLambdaError(500, fmt.Sprintf("Error putting item into DynamoDB: %s", err.Error()))
	}

	return updatedPlayer, nil
}

// func getPlayer(request events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
func getPlayer(ctx context.Context, req events.LambdaFunctionURLRequest, params lambdahandler.Params) (interface{}, lambdahandler.LambdaError) {
	uid := params["uid"]

	// Find player in DynamoDB
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)

	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("players"), // Replace with your table name
		Key: map[string]*dynamodb.AttributeValue{
			"uid": {
				N: aws.String(uid),
			},
		},
	})
	if err != nil {
		return nil, lambdahandler.NewLambdaError(500, fmt.Sprintf("Error getting item from DynamoDB: %s", err.Error()))
	}

	if result.Item == nil {
		return nil, lambdahandler.NewLambdaError(404, "Player not found")
	}

	var player Player
	err = dynamodbattribute.UnmarshalMap(result.Item, &player)
	if err != nil {
		return nil, lambdahandler.NewLambdaError(500, fmt.Sprintf("Error unmarshalling DynamoDB item: %s", err.Error()))
	}

	return player, nil

}

// func deletePlayer(request events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
func deletePlayer(ctx context.Context, request events.LambdaFunctionURLRequest, params lambdahandler.Params) (interface{}, lambdahandler.LambdaError) {
	uid := params["uid"]

	// Find player in DynamoDB
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)

	_, err := svc.DeleteItem(&dynamodb.DeleteItemInput{
		TableName: aws.String("players"), // Replace with your table name
		Key: map[string]*dynamodb.AttributeValue{
			"uid": {
				N: aws.String(uid),
			},
		},
	})
	if err != nil {
		return nil, lambdahandler.NewLambdaError(500, fmt.Sprintf("Error deleting item from DynamoDB: %s", err.Error()))
	}

	return "Player deleted successfully", nil
}
