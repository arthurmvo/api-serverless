package main

import (
	// "context"
	"context"
	"fmt"

	// "github.com/arthurmvo/lambdahandler"
	"github.com/arthurmvo/lambdahandler"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// func getClans() (events.LambdaFunctionURLResponse, error) {
func getClans(ctx context.Context, req events.LambdaFunctionURLRequest, params lambdahandler.Params) (interface{}, lambdahandler.LambdaError) {
	fmt.Print("Starting router")
	// fmt.Print(req)
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)

	input := &dynamodb.ScanInput{
		TableName: aws.String("clans"),
	}

	result, err := svc.Scan(input)
	if err != nil {
		return nil, lambdahandler.NewLambdaError(500, fmt.Sprintf("Error scanning DynamoDB: %s", err.Error()))
	}

	var clans []Clan
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &clans)
	if err != nil {
		return nil, lambdahandler.NewLambdaError(500, fmt.Sprintf("Error unmarshalling result: %s", err.Error()))
	}

	return clans, nil
}
