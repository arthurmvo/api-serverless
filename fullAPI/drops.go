package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/arthurmvo/lambdahandler"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// func getDrops() (events.LambdaFunctionURLResponse, error) {
func getDrops(ctx context.Context, req events.LambdaFunctionURLRequest, params lambdahandler.Params) (interface{}, error) {
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)

	input := &dynamodb.ScanInput{
		TableName: aws.String("drops"),
	}

	result, err := svc.Scan(input)
	if err != nil {
		return events.LambdaFunctionURLResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf("Error scanning DynamoDB: %s", err.Error()),
		}, nil
	}

	var drops []Drop
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &drops)
	if err != nil {
		return events.LambdaFunctionURLResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf("Error unmarshalling result: %s", err.Error()),
		}, nil
	}

	body, err := json.Marshal(drops)
	if err != nil {
		return events.LambdaFunctionURLResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf("Error marshalling response: %s", err.Error()),
		}, nil
	}

	return events.LambdaFunctionURLResponse{
		StatusCode: 200,
		Body:       string(body),
	}, nil
}

// func createDrop(request events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
func createDrop(ctx context.Context, request events.LambdaFunctionURLRequest, params lambdahandler.Params) (interface{}, error) {
	var drop Drop
	err := json.Unmarshal([]byte(request.Body), &drop)
	if err != nil {
		return events.LambdaFunctionURLResponse{
			StatusCode: 400,
			Body:       fmt.Sprintf("Error unmarshalling request body: %s", err.Error()),
		}, nil
	}

	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)

	av, err := dynamodbattribute.MarshalMap(drop)
	if err != nil {
		return events.LambdaFunctionURLResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf("Error marshalling drop: %s", err.Error()),
		}, nil
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String("drops"),
		Item:      av,
	}

	_, err = svc.PutItem(input)
	if err != nil {
		return events.LambdaFunctionURLResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf("Error putting item into DynamoDB: %s", err.Error()),
		}, nil
	}

	return events.LambdaFunctionURLResponse{
		StatusCode: 201,
		Body:       "Drop created successfully",
	}, nil
}
