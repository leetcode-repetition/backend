package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"time"
)

const DAY_IN_SECONDS = 60 * 60 * 24

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

func get_completed_dates(title_slug string) []int64 {
	subscriber := get_subscriber("jacobhmurrah@gmail.com")
	problems := subscriber.Problems

	for _, problem := range problems {
		if problem.TitleSlug == title_slug {
			return problem.CompletedTimestamps
		}
	}

	return []int64{}
}

func get_offset(difficulty string, title_slug string) int64 {
	completed_dates := get_completed_dates(title_slug)
	length_completed := len(completed_dates)

	maxCompleted := map[string]map[string]int64{
		"Easy":   {"max_exponent": 2, "days": 14},
		"Medium": {"max_exponent": 3, "days": 7},
		"Hard":   {"max_exponent": 4, "days": 3},
	}

	if length_completed > int(maxCompleted[difficulty]["max_exponent"]) {
		length_completed = int(maxCompleted[difficulty]["max_exponent"])
	}

	return int64(math.Pow(2, float64(length_completed))) * maxCompleted["difficulty"]["days"] * DAY_IN_SECONDS
}

func create_leetcode_problem_object(title_slug string) LeetCodeProblem {
	completed_dates := get_completed_dates(title_slug)

	// get previos completed
	data, _ := query_leetcode_api("http://localhost:3000/select?titleSlug=" + title_slug)

	difficulty, _ := data["difficulty"].(string)
	repeat_on := time.Now().Unix() + get_offset(difficulty, title_slug)

	return LeetCodeProblem{
		Link:                data["link"].(string),
		Title:               data["questionTitle"].(string),
		TitleSlug:           title_slug,
		Difficulty:          data["difficulty"].(string),
		CompletedTimestamps: completed_dates,
		RepeatTimestamp:     repeat_on,
	}
}

func format_problems(unformatted_problems []interface{}) []LeetCodeProblem {
	var formatted_problems []LeetCodeProblem

	for _, item := range unformatted_problems {
		unformatted_problem, _ := item.(map[string]interface{})
		timestamp, _ := strconv.ParseInt(unformatted_problem["timestamp"].(string), 10, 64)

		if timestamp > (time.Now().Unix() - DAY_IN_SECONDS) {
			problem := create_leetcode_problem_object(unformatted_problem["titleSlug"].(string))
			formatted_problems = append(formatted_problems, problem)
		}
	}

	return formatted_problems
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

func handle_new_subscriptions() {
	// connect this to frontend!!!
	var username string
	var email string

	fmt.Print("Enter LeetCode username: ")
	fmt.Scanln(&username)

	fmt.Print("Enter email: ")
	fmt.Scanln(&email)

	add_new_subscriber_to_database(email, username)
	send_welcome_email(username)
}

// func get_problems_accepted_yesterday() int {
func main() {
	// https://alfa-leetcode-api.onrender.com/username/acsubmission
	data, _ := query_leetcode_api("http://localhost:3000/jmurrah/acsubmission")
	unformatted_problems, _ := data["submission"].([]interface{})
	format_problems(unformatted_problems)
}
