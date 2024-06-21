package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/supabase-community/supabase-go"
)

type LeetCodeProblem struct {
	Link           string
	TitleSlug      string
	Difficulty     string
	Tags           []string
	CompletedDates []time.Time
	RepeatDate     time.Time
}

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

func get_subscriber(email string) Subscriber {
	client := create_supabase_client()
	table := os.Getenv("SUPABASE_TABLE")

	raw_data, _, _ := client.From(table).Select("*", "", false).Eq("email", email).Execute()
	var infoStruct []Subscriber

	json.Unmarshal([]byte(string(raw_data)), &infoStruct)

	return infoStruct[0]
}
