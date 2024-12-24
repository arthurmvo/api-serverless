package main

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Clan struct {
	ID      string   `json:"id"`
	Clan    string   `json:"clan"`
	Icon    string   `json:"icon"`
	Element string   `json:"element"`
	Lurer   string   `json:"lurer"`
	Hunts   []string `json:"hunts"`
}

func getClans(request events.APIGatewayProxyRequest) (events.LambdaFunctionURLResponse, error) {
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)

	input := &dynamodb.ScanInput{
		TableName: aws.String("clans"),
	}

	result, err := svc.Scan(input)
	if err != nil {
		return events.LambdaFunctionURLResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf("Error scanning DynamoDB: %s", err.Error()),
		}, nil
	}

	var clans []Clan
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &clans)
	if err != nil {
		return events.LambdaFunctionURLResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf("Error unmarshalling result: %s", err.Error()),
		}, nil
	}

	body, err := json.Marshal(clans)
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

func main() {
	lambda.Start(getClans)
}