package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/supabase-community/supabase-go"
)

type Subscriber struct {
	Email    string            `json:"email"`
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

func add_new_subscriber_to_database(email string, username string) {
	subscriber := Subscriber{
		Email:    email,
		Username: username,
		Problems: []LeetCodeProblem{},
	}

	client := create_supabase_client()
	table := os.Getenv("SUPABASE_TABLE")
	client.From(table).Upsert(subscriber, "email", "success", "").Execute()
}

func get_all_subscribers_with_recent_activity() []Subscriber {
	client := create_supabase_client()
	table := os.Getenv("SUPABASE_TABLE")

	raw_data, _, _ := client.From(table).Select("*", "", false).Execute()
	var subscribers []Subscriber

	json.Unmarshal([]byte(string(raw_data)), &subscribers)

	return subscribers
}

func get_subscriber(email string) Subscriber {
	client := create_supabase_client()
	table := os.Getenv("SUPABASE_TABLE")

	raw_data, _, _ := client.From(table).Select("*", "", false).Eq("email", email).Execute()
	var subscriber []Subscriber

	json.Unmarshal([]byte(string(raw_data)), &subscriber)

	return subscriber[0]
}
