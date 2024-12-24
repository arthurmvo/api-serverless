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

type Drop struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Icon        string `json:"icon"`
	Link        string `json:"link"`
	IsNightmare bool   `json:"isNightmare"`
	NpcPrice    int    `json:"npcPrice"`
	Chance      int    `json:"chance"`
	IsRare      bool   `json:"isRare"`
}

func getDrops(request events.APIGatewayProxyRequest) (events.LambdaFunctionURLResponse, error) {
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

func main() {
	lambda.Start(getDrops)
}
