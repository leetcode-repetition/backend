package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type GraphQL struct {
	Query     string
	Variables map[string]interface{}
}

func query_leetcode_graphql_api(request *GraphQL) ([]byte, error) {
	body, _ := json.Marshal(request)
	response, err := http.Post("https://leetcode.com/graphql", "application/json", bytes.NewBuffer(body))

	if err != nil {
		return nil, fmt.Errorf("error sending request to API endpoint. %v", err)
	}
	defer response.Body.Close()

	response_body, _ := io.ReadAll(response.Body)
	// fmt.Println(string(response_body))
	return response_body, nil
}

// func schedule(f func(), firstRun time.Duration) {
// 	time.Sleep(firstRun)
// 	ticker := time.NewTicker(24 * time.Hour)
// 	defer ticker.Stop()
// 	for {
// 		select {
// 		case <-ticker.C:
// 			f()
// 		}
// 	}
// }

func main() {
	godotenv.Load("../.env")

	var leetcode_username string
	var email string

	fmt.Print("Enter LeetCode username: ")
	fmt.Scanln(&leetcode_username)

	fmt.Print("Enter email: ")
	fmt.Scanln(&email)

	add_username_and_email_to_database(leetcode_username, email)

	data, _ := os.ReadFile("graphql/get_recent_submissions.graphql")
	query := string(data)
	request := &GraphQL{
		Query: query,
		Variables: map[string]interface{}{
			"username": leetcode_username,
		},
	}

	query_leetcode_graphql_api(request)
}
