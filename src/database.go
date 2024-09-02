package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/supabase-community/supabase-go"
)

type User struct {
	Username string            `json:"username"`
	Problems []LeetCodeProblem `json:"problems"`
}

func create_supabase_client() *supabase.Client {
	client, err := supabase.NewClient(os.Getenv("SUPABASE_URL"), os.Getenv("SUPABASE_KEY"), &supabase.ClientOptions{})
	if err != nil {
		fmt.Println("cannot initalize client", err)
	}
	return client
}

func upsert_database(username string, problems []LeetCodeProblem) {
	user := User{
		Username: username,
		Problems: problems,
	}
	client := create_supabase_client()
	table := os.Getenv("SUPABASE_TABLE")
	client.From(table).Upsert(user, "username", "success", "").Execute()
}

func get_problems_from_database(username string) []LeetCodeProblem {
	var problems []LeetCodeProblem
	var raw_response []map[string]json.RawMessage

	client := create_supabase_client()
	table := os.Getenv("SUPABASE_TABLE")

	raw_data, _, _ := client.From(table).Select("problems", "", false).Eq("username", username).Execute()
	json.Unmarshal(raw_data, &raw_response)

	if len(raw_response) == 0 {
		upsert_database(username, []LeetCodeProblem{})
	}
	json.Unmarshal(raw_response[0]["problems"], &problems)
	return problems
}
