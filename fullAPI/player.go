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
func createPlayer(ctx context.Context, request events.LambdaFunctionURLRequest, params lambdahandler.Params) (interface{}, error) {
	// Unmarshal the request body into a Player struct
	var newPlayer Player
	err := json.Unmarshal([]byte(request.Body), &newPlayer)
	if err != nil {
		return events.LambdaFunctionURLResponse{
			StatusCode: 400,
			Body:       fmt.Sprintf("Error unmarshalling request body: %s", err.Error()),
		}, nil
	}

	fmt.Println(newPlayer, "newPlayer")

	// Marshal the Player struct into an AttributeValue
	av, err := dynamodbattribute.MarshalMap(newPlayer)
	if err != nil {
		return events.LambdaFunctionURLResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf("Error marshalling player: %s", err.Error()),
		}, nil
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
		return events.LambdaFunctionURLResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf("Error putting item into DynamoDB: %s", err.Error()),
		}, nil
	}

	// Marshal the new player into JSON
	newPlayerJSON, err := json.Marshal(newPlayer)
	if err != nil {
		return events.LambdaFunctionURLResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf("Error marshalling player to JSON: %s", err.Error()),
		}, nil
	}

	return events.LambdaFunctionURLResponse{
		StatusCode: 201,
		Body:       string(newPlayerJSON),
	}, nil
}

// func updatePlayer(request events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
func updatePlayer(ctx context.Context, request events.LambdaFunctionURLRequest, params lambdahandler.Params) (interface{}, error) {

	uid, err := strconv.Atoi(params["uid"])
	if err != nil {
		return events.LambdaFunctionURLResponse{
			StatusCode: 400,
			Body:       fmt.Sprintf("Invalid ID format: %s", err.Error()),
		}, nil
	}

	// Unmarshal the request body into a Player struct
	var updatedPlayer Player
	err = json.Unmarshal([]byte(request.Body), &updatedPlayer)
	if err != nil {
		return events.LambdaFunctionURLResponse{
			StatusCode: 400,
			Body:       fmt.Sprintf("Error unmarshalling request body: %s", err.Error()),
		}, nil
	}

	// Ensure the ID in the request body matches the ID in the path
	if updatedPlayer.UID != uid {
		return events.LambdaFunctionURLResponse{
			StatusCode: 400,
			Body:       "ID in the request body does not match ID in the path",
		}, nil
	}

	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)

	av, err := dynamodbattribute.MarshalMap(updatedPlayer)
	if err != nil {
		return events.LambdaFunctionURLResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf("Error marshalling player: %s", err.Error()),
		}, nil
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String("players"), // Replace with your table name
		Item:      av,
	}

	_, err = svc.PutItem(input)
	if err != nil {
		return events.LambdaFunctionURLResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf("Error putting item into DynamoDB: %s", err.Error()),
		}, nil
	}

	jsonBody, _ := json.Marshal(updatedPlayer)

	return lambdahandler.SuccessResponse(jsonBody), nil
}

// func getPlayer(request events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
func getPlayer(ctx context.Context, req events.LambdaFunctionURLRequest, params lambdahandler.Params) (interface{}, error) {
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
		return events.LambdaFunctionURLResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf("Error getting item from DynamoDB: %s", err.Error()),
		}, nil
	}

	if result.Item == nil {
		return events.LambdaFunctionURLResponse{
			StatusCode: 404,
			Body:       "Player not found",
		}, nil
	}

	var player Player
	err = dynamodbattribute.UnmarshalMap(result.Item, &player)
	if err != nil {
		return events.LambdaFunctionURLResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf("Error unmarshalling DynamoDB item: %s", err.Error()),
		}, nil
	}

	playerJSON, err := json.Marshal(player)
	if err != nil {
		return events.LambdaFunctionURLResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf("Error marshalling player to JSON: %s", err.Error()),
		}, nil
	}

	return lambdahandler.SuccessResponse(playerJSON), nil

}

// func deletePlayer(request events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
func deletePlayer(ctx context.Context, request events.LambdaFunctionURLRequest, params lambdahandler.Params) (interface{}, error) {
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
		return events.LambdaFunctionURLResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf("Error deleting item from DynamoDB: %s", err.Error()),
		}, nil
	}

	return lambdahandler.SuccessResponse("Player deleted successfully"), nil
}
