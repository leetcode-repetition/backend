package main

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	username := request.QueryStringParameters["username"]
	problemTitleSlug := request.QueryStringParameters["problemTitleSlug"]

	deleteProblemFromDatabase(username, problemTitleSlug)

	responseBody, _ := json.Marshal(map[string]interface{}{
		"message": "Delete row data processed",
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
