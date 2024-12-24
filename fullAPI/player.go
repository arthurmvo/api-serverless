package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

func createPlayer(request events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
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

func updatePlayer(request events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {

	fmt.Println(request, "request")
	fmt.Println(request.RawPath, "request.RawPath")
	fmt.Println(request.RequestContext, "request.RequestContext")
	uidStr := strings.Split(request.RequestContext.HTTP.Path, "/")[2]
	uid, _ := strconv.Atoi(uidStr)

	// Unmarshal the request body into a Player struct
	var updatedPlayer Player
	err := json.Unmarshal([]byte(request.Body), &updatedPlayer)
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

	return events.LambdaFunctionURLResponse{
		StatusCode: 200,
		Body:       string(jsonBody),
	}, nil
}

func getPlayer(request events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	uidStr := strings.Split(request.RequestContext.HTTP.Path, "/")[2]

	// Find player in DynamoDB
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)

	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("players"), // Replace with your table name
		Key: map[string]*dynamodb.AttributeValue{
			"uid": {
				N: aws.String(uidStr),
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

	return events.LambdaFunctionURLResponse{
		StatusCode: 200,
		Body:       string(playerJSON),
	}, nil

}

func deletePlayer(request events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	uidStr := strings.Split(request.RequestContext.HTTP.Path, "/")[2]

	// Find player in DynamoDB
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)

	_, err := svc.DeleteItem(&dynamodb.DeleteItemInput{
		TableName: aws.String("players"), // Replace with your table name
		Key: map[string]*dynamodb.AttributeValue{
			"uid": {
				N: aws.String(uidStr),
			},
		},
	})
	if err != nil {
		return events.LambdaFunctionURLResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf("Error deleting item from DynamoDB: %s", err.Error()),
		}, nil
	}

	return events.LambdaFunctionURLResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "OPTIONS, GET, PUT, POST",
			"Access-Control-Allow-Headers": "Content-Type",
		},
		Body: "Player deleted successfully",
	}, nil
}
