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
	fmt.Println(string(response_body))
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

// rename this function to convey ADDING NEW EMAIL TO MAILING LIST
func main() {
	godotenv.Load("../.env")

	var username string
	var email string

	fmt.Print("Enter LeetCode username: ")
	fmt.Scanln(&username)

	fmt.Print("Enter email: ")
	fmt.Scanln(&email)

	// check if the username exists and check if email exists.

	query, _ := os.ReadFile("graphql/get_recent_submissions.graphql")
	payload := &GraphQL{
		Query: string(query),
		Variables: map[string]interface{}{
			"username": username,
		},
	}

	query_leetcode_graphql_api(payload)
	// fmt.Println(request)

	// add_row_to_database(email, username)
	// send_welcome_email(username)

	// send_daily_email(get_row(email))
}
