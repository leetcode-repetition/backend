package main

import (
	"encoding/json"

	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	shared "github.com/jmurrah/leetcode-repetition-shared"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Raw request: %+v", request)
	log.Printf("Query parameters: %+v", request.QueryStringParameters)
	log.Printf("Path parameters: %+v", request.PathParameters)

	username := request.QueryStringParameters["username"]

	log.Printf("Username: %s", username)

	var problems = []map[string]interface{}{}
	for _, problem := range shared.GetProblemsFromDatabase(username) {
		problems = append(problems, map[string]interface{}{
			"link":               problem.Link,
			"titleSlug":          problem.TitleSlug,
			"repeatDate":         problem.RepeatDate,
			"lastCompletionDate": problem.LastCompletionDate,
		})
	}

	log.Printf("Problems: %+v", problems)

	responseBody, _ := json.Marshal(map[string]interface{}{
		"message": "Get table data processed",
		"table":   problems,
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
