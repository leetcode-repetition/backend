package main

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	shared "github.com/jmurrah/leetcode-repetition-shared"
)

func init() {
	if err := shared.InitSupabaseClient(); err != nil {
		log.Printf("Failed to initialize Supabase client: %v", err)
	}
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	username := request.QueryStringParameters["username"]

	var requestData map[string]interface{}
	json.Unmarshal([]byte(request.Body), &requestData)

	problem := shared.LeetCodeProblem{
		Link:               requestData["link"].(string),
		TitleSlug:          requestData["titleSlug"].(string),
		RepeatDate:         requestData["repeatDate"].(string),
		LastCompletionDate: requestData["lastCompletionDate"].(string),
	}

	shared.UpsertProblemIntoDatabase(username, problem)

	responseBody, _ := json.Marshal(map[string]interface{}{
		"message": "Inserted row data processed",
		"data":    requestData,
	})

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(responseBody),
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
