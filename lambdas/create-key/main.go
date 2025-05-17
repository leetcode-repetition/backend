package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

type LeetCodeUserStatus struct {
	UserId   json.Number `json:"userId"`
	Username string      `json:"username"`
}

type LeetCodeGraphQLResponse struct {
	Data struct {
		UserStatus LeetCodeUserStatus `json:"userStatus"`
	} `json:"data"`
}

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Printf("Failed to initialize API Gateway client: %v", err)
		return
	}
	apiGatewayClient = apigateway.NewFromConfig(cfg)
}

func fetchLeetCodeUserInfo(csrfToken, leetcodeSession string) (*LeetCodeUserStatus, error) {
	if csrfToken == "" || leetcodeSession == "" {
		return nil, fmt.Errorf("missing required authentication tokens")
	}

	query := map[string]interface{}{
		"query": `
            query {
                userStatus {
                    userId
                    username
                }
            }
        `,
	}

	jsonQuery, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "https://leetcode.com/graphql", bytes.NewBuffer(jsonQuery))
	if err != nil {
		return nil, err
	}

	// Match the exact headers used in the JavaScript implementation
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-csrftoken", csrfToken) // Note: lowercase 'x-csrftoken' instead of 'X-Csrf-Token'
	req.Header.Set("Cookie", fmt.Sprintf("csrftoken=%s; LEETCODE_SESSION=%s", csrfToken, leetcodeSession))

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	log.Printf("LeetCode API response: %s", string(body))

	var graphQLResponse LeetCodeGraphQLResponse
	if err := json.Unmarshal(body, &graphQLResponse); err != nil {
		return nil, err
	}

	return &graphQLResponse.Data.UserStatus, nil
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	cors := map[string]string{
		"Access-Control-Allow-Origin":      request.Headers["origin"],
		"Access-Control-Allow-Credentials": "true",
		"Access-Control-Allow-Methods":     "POST",
		"Access-Control-Allow-Headers":     "Content-Type,X-Api-Key,X-Pkce-Verifier,X-Auth-Code,X-Csrf-Token,X-Leetcode-Session",
		"Content-Type":                     "application/json",
	}
	log.Printf("Request headers: %+v", request.Headers)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	csrfToken := request.Headers["x-csrf-token"]
	leetcodeSession := request.Headers["x-leetcode-session"]

	userInfo, _ := fetchLeetCodeUserInfo(csrfToken, leetcodeSession)
	userId := userInfo.UserId.String()

	log.Printf("%+v", userInfo)

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
