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

func create_supabase_client() (*supabase.Client, error) {
	client, err := supabase.NewClient(os.Getenv("SUPABASE_URL"), os.Getenv("SUPABASE_KEY"), &supabase.ClientOptions{})
	if err != nil {
		fmt.Println("Cannot initalize client", err)
	} else {
		fmt.Println("Initailized supabase client")
	}
	return client, err
}

func upsert_database(username string, problems []LeetCodeProblem) {
	user := User{
		Username: username,
		Problems: problems,
	}
	client, e := create_supabase_client()
	if e != nil {
		fmt.Println("Error creating supabase client:", e)
		return
	}
	table := os.Getenv("SUPABASE_TABLE")

	_, _, err := client.From(table).Upsert(user, "username", "success", "").Execute()
	if err != nil {
		fmt.Println("Error upserting database:", err)
	}
	fmt.Println("Successfully upserted database entry for user:", username)
}

func get_problems_from_database(username string) []LeetCodeProblem {
	var problems []LeetCodeProblem
	var raw_response []map[string]json.RawMessage

	client, e := create_supabase_client()
	if e != nil {
		fmt.Println("Error creating supabase client:", e)
		return []LeetCodeProblem{}
	}
	table := os.Getenv("SUPABASE_TABLE")

	fmt.Println("Getting problems from database for user:", username)
	raw_data, _, _ := client.From(table).Select("problems", "", false).Eq("username", username).Execute()
	json.Unmarshal(raw_data, &raw_response)

	if len(raw_response) == 0 {
		fmt.Println("User not found. Adding user to database:", username)
		upsert_database(username, []LeetCodeProblem{})
		return []LeetCodeProblem{}
	}
	json.Unmarshal(raw_response[0]["problems"], &problems)
	return problems
}
