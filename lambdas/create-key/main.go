package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
)

var apiGatewayClient *apigateway.Client

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Printf("Failed to initialize API Gateway client: %v", err)
		return
	}
	apiGatewayClient = apigateway.NewFromConfig(cfg)
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	log.Printf("Query parameters: %+v", request.QueryStringParameters)
	userId := request.QueryStringParameters["userId"]
	log.Printf("UserId: %s", userId)

	keyName := "LRE_" + userId + "_" + time.Now().Format("20060102150405")
	log.Printf("Creating API key with name: %s", keyName)

	keyInput := &apigateway.CreateApiKeyInput{
		Name:    aws.String(keyName),
		Enabled: true,
	}

	keyResult, err := apiGatewayClient.CreateApiKey(ctx, keyInput)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Error creating API key",
		}, err
	}

	planInput := &apigateway.CreateUsagePlanKeyInput{
		KeyId:       keyResult.Id,
		KeyType:     aws.String("API_KEY"),
		UsagePlanId: aws.String(os.Getenv("USAGE_PLAN_ID")),
	}

	_, err = apiGatewayClient.CreateUsagePlanKey(ctx, planInput)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Error creating usage plan key",
		}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       aws.ToString(keyResult.Value),
		Headers: map[string]string{
			"Content-Type":                 "application/json",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "GET, POST, DELETE, OPTIONS",
			"Access-Control-Allow-Headers": "Content-Type",
		},
	}, nil
}

func main() {
	lambda.Start(handler)
}
