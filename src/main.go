package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func generic_handler(specific_handler func(*http.Request, map[string]interface{}) map[string]interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request_data map[string]interface{}
		json.NewDecoder(r.Body).Decode(&request_data)
		fmt.Println("Received data:", request_data)

		response_data := specific_handler(r, request_data)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response_data)
	}
}

func get_table_handler(r *http.Request, data map[string]interface{}) map[string]interface{} {
	username := r.URL.Query().Get("username")
	if username == "" {
		return map[string]interface{}{"error": "Username not provided"}
	}
	fmt.Println("Processing get-table data for user:", username)

	var problems = []map[string]interface{}{}

	for _, problem := range get_problems_from_database(username) {
		problems = append(problems, map[string]interface{}{
			"link":               problem.Link,
			"titleSlug":          problem.TitleSlug,
			"difficulty":         problem.Difficulty,
			"repeatDate":         problem.RepeatDate,
			"lastCompletionDate": problem.LastCompletionDate,
		})
	}
	fmt.Println("Problems for user", username, ":", problems)
	return map[string]interface{}{
		"message": "Get table data processed",
		"table":   problems,
	}
}

func delete_row_handler(r *http.Request, data map[string]interface{}) map[string]interface{} {
	fmt.Println("Processing delete-row data:", data)

	username := r.URL.Query().Get("username")
	problem_title_slug := r.URL.Query().Get("problemTitleSlug")
	if username == "" || problem_title_slug == "" {
		fmt.Println("Username or problem title slug not provided")
		return map[string]interface{}{"error": "Username or problem title slug not provided"}
	}

	delete_problem_from_database(username, problem_title_slug)

	return map[string]interface{}{
		"message": "Delete row data processed",
		"data":    data,
	}
}

func insert_row_handler(r *http.Request, data map[string]interface{}) map[string]interface{} {
	username := r.URL.Query().Get("username")
	if username == "" {
		return map[string]interface{}{"error": "Username not provided"}
	}

	problem := LeetCodeProblem{
		Link:               data["link"].(string),
		TitleSlug:          data["titleSlug"].(string),
		Difficulty:         data["difficulty"].(string),
		RepeatDate:         data["repeatDate"].(string),
		LastCompletionDate: data["lastCompletionDate"].(string),
	}
	upsert_problem_into_database(username, problem)

	return map[string]interface{}{
		"message": "Inserted row data processed",
		"data":    data,
	}
}

func main() {
	godotenv.Load()
	fmt.Println("program running!")

	http.HandleFunc("/get-table", enableCORS(generic_handler(get_table_handler)))
	http.HandleFunc("/delete-row", enableCORS(generic_handler(delete_row_handler)))
	http.HandleFunc("/insert-row", enableCORS(generic_handler(insert_row_handler)))

	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
