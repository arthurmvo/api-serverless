package main

import (
	"context"
	"fmt"

	"github.com/arthurmvo/lambdahandler"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// func getPokemons() (events.LambdaFunctionURLResponse, error) {
func getPokemons(ctx context.Context, req events.LambdaFunctionURLRequest, params lambdahandler.Params) (interface{}, lambdahandler.LambdaError) {
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)

	input := &dynamodb.ScanInput{
		TableName: aws.String("pokemons"),
	}

	result, err := svc.Scan(input)
	if err != nil {
		return nil, lambdahandler.NewLambdaError(500, fmt.Sprintf("Error scanning DynamoDB: %s", err.Error()))
	}

	var pokemons []Pokemon
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &pokemons)
	if err != nil {
		return nil, lambdahandler.NewLambdaError(500, fmt.Sprintf("Error unmarshalling result: %s", err.Error()))
	}

	return pokemons, nil
}
