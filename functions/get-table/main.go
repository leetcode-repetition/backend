package main

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	username := request.QueryStringParameters["username"]

	var problems = []map[string]interface{}{}
	for _, problem := range getProblemsFromDatabase(username) {
		problems = append(problems, map[string]interface{}{
			"link":               problem.Link,
			"titleSlug":          problem.TitleSlug,
			"repeatDate":         problem.RepeatDate,
			"lastCompletionDate": problem.LastCompletionDate,
		})
	}

	responseBody, _ := json.Marshal(map[string]interface{}{
		"message": "Get table data processed",
		"table":   problems,
	})

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(responseBody),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}

func main() {
	lambda.Start(handler)
}
