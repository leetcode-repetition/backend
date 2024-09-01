package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"
)

type LeetCodeProblem struct {
	Link                string
	Title               string
	TitleSlug           string
	Difficulty          string
	CompletedTimestamps []int64
	RepeatTimestamp     int64
}

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

func get_completed_timestamps(title_slug string) []int64 {
	subscriber := get_subscriber("jacobhmurrah@gmail.com")
	problems := subscriber.Problems

	for _, problem := range problems {
		if problem.TitleSlug == title_slug {
			return problem.CompletedTimestamps
		}
	}

	return []int64{}
}

func get_offset(difficulty string, completed_dates []int64) int64 {
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

func create_leetcode_problem_object(title_slug string, latest_completion_timestamp int64) LeetCodeProblem {
	problem_data, _ := query_leetcode_api("http://localhost:3000/select?titleSlug=" + title_slug)

	difficulty, _ := problem_data["difficulty"].(string)

	completed_timestamps := get_completed_timestamps(title_slug)
	completed_timestamps = append(completed_timestamps, latest_completion_timestamp)

	return LeetCodeProblem{
		Link:                problem_data["link"].(string),
		Title:               problem_data["questionTitle"].(string),
		TitleSlug:           title_slug,
		Difficulty:          difficulty,
		CompletedTimestamps: completed_timestamps,
		RepeatTimestamp:     latest_completion_timestamp + get_offset(difficulty, completed_timestamps),
	}
}

func transform_into_leetcode_problems(subscriber Subscriber, unformatted_problems []interface{}) []LeetCodeProblem {
	var formatted_problems []LeetCodeProblem

	for _, item := range unformatted_problems {
		unformatted_problem, _ := item.(map[string]interface{})
		latest_completed_timestamp, _ := strconv.ParseInt(unformatted_problem["timestamp"].(string), 10, 64)

		if latest_completed_timestamp > (time.Now().Unix() - DAY_IN_SECONDS) {
			problem := create_leetcode_problem_object(unformatted_problem["titleSlug"].(string), latest_completed_timestamp)
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
// func main() {
// 	godotenv.Load("../.env")

// 	//alfa-leetcode-api.onrender.com/username/acsubmission
// 	for _, subscriber := range get_all_subscribers() {
// 		submission_data, _ := query_leetcode_api("http://localhost:3000/" + subscriber.Username + "/acsubmission")
// 		unformatted_problems, _ := submission_data["submission"].([]interface{})
// 		problems := transform_into_leetcode_problems(subscriber, unformatted_problems)
// 		send_spaced_repetition_email()
// 	}
// }

func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func main() {
	fmt.Println("program running!")
	http.HandleFunc("/hello", enableCORS(helloHandler))

	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hello, World!")

	response := map[string]string{"message": "Hello, World!"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// func getSubscribersHandler(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != http.MethodGet {
// 		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 		return
// 	}

// 	subscribers := get_all_subscribers_with_recent_activity()

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(subscribers)
// }
