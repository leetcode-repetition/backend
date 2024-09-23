package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

type LeetCodeProblem struct {
	link               string
	titleSlug          string
	difficulty         string
	repeatDate         string
	lastCompletionDate string
}

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

func genericHandler(specificHandler func(*http.Request, map[string]interface{}) map[string]interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requestData map[string]interface{}
		json.NewDecoder(r.Body).Decode(&requestData)
		fmt.Println("Received data:", requestData)

		responseData := specificHandler(r, requestData)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(responseData)
	}
}

func getTableHandler(r *http.Request, data map[string]interface{}) map[string]interface{} {
	username := r.URL.Query().Get("username")
	if username == "" {
		return map[string]interface{}{"error": "Username not provided"}
	}
	fmt.Println("Processing get-table data for user:", username)

	var problems = []map[string]interface{}{}

	for _, problem := range get_problems_from_database(username) {
		problems = append(problems, map[string]interface{}{
			"link":               problem.link,
			"titleSlug":          problem.titleSlug,
			"difficulty":         problem.difficulty,
			"repeatDate":         problem.repeatDate,
			"lastCompletionDate": problem.lastCompletionDate,
		})
	}
	fmt.Println("Problems for user", username, ":", problems)
	return map[string]interface{}{
		"message": "Get table data processed",
		"table":   problems,
	}
}

func deleteRowHandler(r *http.Request, data map[string]interface{}) map[string]interface{} {
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

func insertRowHandler(r *http.Request, data map[string]interface{}) map[string]interface{} {
	username := r.URL.Query().Get("username")
	if username == "" {
		return map[string]interface{}{"error": "Username not provided"}
	}

	problem := LeetCodeProblem{
		link:               data["link"].(string),
		titleSlug:          data["titleSlug"].(string),
		difficulty:         data["difficulty"].(string),
		repeatDate:         data["repeatDate"].(string),
		lastCompletionDate: data["lastCompletionDate"].(string),
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

	http.HandleFunc("/get-table", enableCORS(genericHandler(getTableHandler)))
	http.HandleFunc("/delete-row", enableCORS(genericHandler(deleteRowHandler)))
	http.HandleFunc("/insert-row", enableCORS(genericHandler(insertRowHandler)))

	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
