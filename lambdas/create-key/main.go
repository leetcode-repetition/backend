package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/apigateway"
	"github.com/golang-jwt/jwt/v4"

	shared "github.com/leetcode-repetition/shared"
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

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
	IdToken      string `json:"id_token,omitempty"`
}

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Printf("Failed to initialize API Gateway client: %v", err)
		return
	}
	apiGatewayClient = apigateway.NewFromConfig(cfg)
	if err := shared.InitSupabaseClient(); err != nil {
		log.Printf("Failed to initialize Supabase client: %v", err)
	}
}

func extractSubjectFromIdToken(idToken string) (string, error) {
	parser := jwt.Parser{
		SkipClaimsValidation: true,
	}

	token, _, err := parser.ParseUnverified(idToken, jwt.MapClaims{})
	if err != nil {
		return "", fmt.Errorf("failed to parse ID token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("failed to extract claims from token")
	}

	sub, ok := claims["sub"].(string)
	if !ok || sub == "" {
		return "", fmt.Errorf("subject claim missing from token")
	}

	return sub, nil
}

func exchangeCodeForToken(authCode, pkceVerifier string, clientID string, redirectURI string, tokenEndpoint string) (*TokenResponse, error) {
	clientSecret := os.Getenv("OAUTH_CLIENT_SECRET")

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("code", authCode)
	data.Set("redirect_uri", redirectURI)
	data.Set("code_verifier", pkceVerifier)

	req, err := http.NewRequest("POST", tokenEndpoint, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("token request failed with status %d: %s", resp.StatusCode, string(body))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var tokenResp TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, err
	}

	return &tokenResp, nil
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

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-csrftoken", csrfToken)
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

func CreateNewApiKey(ctx context.Context, userId string) (string, int64, error) {
	keyCreationTime := time.Now().UnixMilli()
	keyName := "LRE_" + userId + "_" + strconv.FormatInt(keyCreationTime, 10)
	log.Printf("Creating API key with name: %s", keyName)

	keyInput := &apigateway.CreateApiKeyInput{
		Name:    aws.String(keyName),
		Enabled: true,
	}

	keyResult, err := apiGatewayClient.CreateApiKey(ctx, keyInput)
	if err != nil {
		return "e", 0, err
	}

	planInput := &apigateway.CreateUsagePlanKeyInput{
		KeyId:       keyResult.Id,
		KeyType:     aws.String("API_KEY"),
		UsagePlanId: aws.String(os.Getenv("USAGE_PLAN_ID")),
	}
	if _, err = apiGatewayClient.CreateUsagePlanKey(ctx, planInput); err != nil {
		return "e", 0, err
	}

	return *keyResult.Value, keyCreationTime, nil
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	cors := map[string]string{
		"Access-Control-Allow-Origin":      request.Headers["origin"],
		"Access-Control-Allow-Credentials": "true",
		"Access-Control-Allow-Methods":     "POST",
		"Access-Control-Allow-Headers":     "Content-Type,X-Api-Key,X-Pkce-Verifier,X-Auth-Code,X-Csrf-Token,X-Leetcode-Session,X-Client-ID,X-Token-Endpoint",
		"Content-Type":                     "application/json",
	}
	log.Printf("Request headers: %+v", request.Headers)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	authCode := request.Headers["x-auth-code"]
	pkceVerifier := request.Headers["x-pkce-verifier"]
	clientID := request.Headers["x-client-id"]
	tokenEndpoint := request.Headers["x-token-endpoint"]
	csrfToken := request.Headers["x-csrf-token"]
	leetcodeSession := request.Headers["x-leetcode-session"]

	var requestBody struct {
		RedirectURI string `json:"redirectUri"`
	}
	if err := json.Unmarshal([]byte(request.Body), &requestBody); err != nil {
		log.Printf("Error parsing request body: %v", err)
	}

	token, err := exchangeCodeForToken(authCode, pkceVerifier, clientID, requestBody.RedirectURI, tokenEndpoint)
	if err != nil {
		log.Printf("Token exchange error: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: 401,
			Headers:    cors,
			Body:       `{"error": "Authentication failed"}`,
		}, nil
	}

	userInfo, _ := fetchLeetCodeUserInfo(csrfToken, leetcodeSession)
	userId := userInfo.UserId.String()
	username := userInfo.Username

	userIdentifier, err := extractSubjectFromIdToken(token.IdToken)
	if err != nil {
		log.Printf("Failed to extract user ID from token: %v", err)
	}

	apiKey, apiKeyCreationTime := shared.GetApiKeyFromDatabase(userId, userIdentifier)
	if apiKey == "" || apiKeyCreationTime == 0 {
		log.Printf("creating new API key for user %s", userId)
		apiKey, apiKeyCreationTime, err = CreateNewApiKey(ctx, userId)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Headers:    cors,
				Body:       "Error creating API key",
			}, err
		}
		shared.UpsertApiKeyIntoDatabase(userId, userIdentifier, apiKey, apiKeyCreationTime)
	}

	log.Printf("%+v", userInfo)
	responseBody, _ := json.Marshal(map[string]interface{}{
		"message":            "Generated new API key!",
		"apiKey":             apiKey,
		"username":           username,
		"userId":             userId,
		"apiKeyCreationTime": apiKeyCreationTime,
	})

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    cors,
		Body:       string(responseBody),
	}, nil
}

func main() { lambda.Start(handler) }
