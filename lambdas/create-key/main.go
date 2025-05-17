package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
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
	origin := request.Headers["Origin"]
	if origin == "" {
		origin = request.Headers["origin"]
	}
	if origin == "" {
		origin = "null"
	}
	cors := map[string]string{
		"Access-Control-Allow-Origin":      origin,
		"Access-Control-Allow-Credentials": "true",
		"Access-Control-Allow-Methods":     "OPTIONS,GET,POST,DELETE",
		"Access-Control-Allow-Headers":     "Content-Type,x-pkce-verifier,x-auth-code,x-csrf-token,x-api-key",
		"Vary":                             "Origin",
	}

	if request.HTTPMethod == http.MethodOptions {
		return events.APIGatewayProxyResponse{
			StatusCode: 204,
			Headers:    cors,
		}, nil
	}

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
			Headers:    cors,
			Body:       "Error creating API key",
		}, err
	}

	planInput := &apigateway.CreateUsagePlanKeyInput{
		KeyId:       keyResult.Id,
		KeyType:     aws.String("API_KEY"),
		UsagePlanId: aws.String(os.Getenv("USAGE_PLAN_ID")),
	}
	if _, err = apiGatewayClient.CreateUsagePlanKey(ctx, planInput); err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Headers:    cors,
			Body:       "Error creating usage plan key",
		}, err
	}

	responseBody, _ := json.Marshal(map[string]interface{}{
		"message": "Generated new API key!",
		"apiKey":  aws.ToString(keyResult.Value),
	})

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    cors,
		Body:       string(responseBody),
	}, nil
}

func main() { lambda.Start(handler) }
