package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/joho/godotenv"
)

type LeetCodeProblem struct {
	Link                string
	Title               string
	Difficulty          string
	RepeatTimestamp     int64
	CompletedTimestamps []int64
}

// func query_leetcode_api(url string) (map[string]interface{}, error) {
// 	response, err := http.Get(url)
// 	if err != nil {
// 		return nil, fmt.Errorf("error sending request to API endpoint. %v", err)
// 	}

// 	defer response.Body.Close()

// 	response_body, err := io.ReadAll(response.Body)
// 	if err != nil {
// 		return nil, fmt.Errorf("error reading data. %v", err)
// 	}

// 	var formatted_response map[string]interface{}
// 	json.Unmarshal(response_body, &formatted_response)

// 	return formatted_response, nil
// }

// func get_completed_timestamps(title_slug string) []int64 {
// 	subscriber := get_subscriber("jacobhmurrah@gmail.com")
// 	problems := subscriber.Problems

// 	for _, problem := range problems {
// 		if problem.TitleSlug == title_slug {
// 			return problem.CompletedTimestamps
// 		}
// 	}

// 	return []int64{}
// }

// func get_offset(difficulty string, completed_dates []int64) int64 {
// 	length_completed := len(completed_dates)

// 	maxCompleted := map[string]map[string]int64{
// 		"Easy":   {"max_exponent": 2, "days": 14},
// 		"Medium": {"max_exponent": 3, "days": 7},
// 		"Hard":   {"max_exponent": 4, "days": 3},
// 	}

// 	if length_completed > int(maxCompleted[difficulty]["max_exponent"]) {
// 		length_completed = int(maxCompleted[difficulty]["max_exponent"])
// 	}

// 	return int64(math.Pow(2, float64(length_completed))) * maxCompleted["difficulty"]["days"] * DAY_IN_SECONDS
// }

// func get_leetcode_problems(username string) []LeetCodeProblem {
// 	raw_problems, _ := get_problems_from_database(username)

// 	difficulty, _ := problem_data["difficulty"].(string)

// 	completed_timestamps := get_completed_timestamps(title_slug)
// 	completed_timestamps = append(completed_timestamps, latest_completion_timestamp)

// 	return LeetCodeProblem{
// 		Link:                problem_data["link"].(string),
// 		Title:               problem_data["questionTitle"].(string),
// 		Difficulty:          difficulty,
// 		CompletedTimestamps: completed_timestamps,
// 		RepeatTimestamp:     latest_completion_timestamp + get_offset(difficulty, completed_timestamps),
// 	}
// }

// func format_problems(unformatted_problems []interface{}) []LeetCodeProblem {
// 	var formatted_problems []LeetCodeProblem

// 	for _, item := range unformatted_problems {
// 		unformatted_problem, _ := item.(map[string]interface{})
// 		latest_completed_timestamp, _ := strconv.ParseInt(unformatted_problem["timestamp"].(string), 10, 64)

// 		if latest_completed_timestamp > (time.Now().Unix() - DAY_IN_SECONDS) {
// 			problem := create_leetcode_problem_object(unformatted_problem["titleSlug"].(string), latest_completed_timestamp)
// 			formatted_problems = append(formatted_problems, problem)
// 		}
// 	}

// 	return formatted_problems
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

func genericHandler(specificHandler func(map[string]interface{}) map[string]interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestData map[string]interface{}
		json.NewDecoder(r.Body).Decode(&requestData)
		fmt.Println("Received data:", requestData)

		responseData := specificHandler(requestData)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(responseData)
	}
}

func getTableHandler(incoming_data map[string]interface{}) map[string]interface{} {
	fmt.Println("Processing get-table data:", incoming_data)
	var problems = []map[string]interface{}{}

	for _, problem := range get_problems_from_database(incoming_data["username"].(string)) {
		problems = append(problems, map[string]interface{}{
			"link":                problem.Link,
			"title":               problem.Title,
			"difficulty":          problem.Difficulty,
			"repeatDate":          time.Unix(problem.RepeatTimestamp, 0).Format("1/12/24"),
			"lastCompleted":       time.Unix(problem.CompletedTimestamps[len(problem.CompletedTimestamps)-1], 0).Format("1/12/24"),
			"completedTimesCount": len(problem.CompletedTimestamps),
		})
	}

	return map[string]interface{}{
		"message": "Get table data processed",
		"data":    problems,
	}
}

func deleteRowHandler(data map[string]interface{}) map[string]interface{} {
	fmt.Println("Processing delete-row data:", data)
	return map[string]interface{}{
		"message": "Delete row data processed",
		"data":    data,
	}
}

func main() {
	godotenv.Load("../.env")
	fmt.Println("program running!")

	http.HandleFunc("/get-table", enableCORS(genericHandler(getTableHandler)))
	http.HandleFunc("/delete-row", enableCORS(genericHandler(deleteRowHandler)))

	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
