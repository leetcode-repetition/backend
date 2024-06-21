package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/joho/godotenv"
)

func query_leetcode_api(url string) (map[string]interface{}, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error sending request to API endpoint. %v", err)
	}

	defer response.Body.Close()

	response_body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading data. %v", err)
	}

	var formatted_response map[string]interface{}
	json.Unmarshal(response_body, &formatted_response)

	return formatted_response, nil
}

func format_problems(unformatted_problems []interface{}) {
	// var formatted_problems []LeetCodeProblem

	for _, item := range unformatted_problems {
		// Assert item to be of type map[string]interface{}
		unf_problem, ok := item.(map[string]interface{})
		if !ok {
			fmt.Println("Item is not of type map[string]interface{}")
			continue
		}

		// Now you can safely index unf_problem
		fmt.Println(unf_problem["lang"])
	}
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

	// https://alfa-leetcode-api.onrender.com/
	data, _ := query_leetcode_api("http://localhost:3000/jmurrah/acsubmission")
	unformatted_problems, _ := data["submission"].([]interface{})
	format_problems(unformatted_problems)
	add_new_subscriber_to_database(email, username)
	// send_welcome_email(username)

	// send_daily_email(get_row(email))
}
