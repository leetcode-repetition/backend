package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/supabase-community/supabase-go"
)

type LeetCodeProblem struct {
	TitleSlug      string
	Difficulty     string
	Tags           []string
	CompletedDates []time.Time
	RepeatDate     time.Time
}

type Row struct {
	Email    string              `json:"email"`
	Username string              `json:"username"`
	Problems []map[string]string `json:"problems"`
}

func create_supabase_client() *supabase.Client {
	client, err := supabase.NewClient(os.Getenv("SUPABASE_URL"), os.Getenv("SUPABASE_KEY"), &supabase.ClientOptions{})
	if err != nil {
		fmt.Println("cannot initalize client", err)
	}
	return client
}

func add_row_to_database(email string, username string) {
	row := Row{
		Email:    email,
		Username: username,
		Problems: []map[string]string{
			{
				"test": "solution1",
			},
			{
				"test": "solution2",
			},
		},
	}

	client := create_supabase_client()
	table := os.Getenv("SUPABASE_TABLE")
	client.From(table).Upsert(row, "email", "success", "").Execute()
}

func get_row(email string) Row {
	client := create_supabase_client()
	table := os.Getenv("SUPABASE_TABLE")

	raw_data, _, _ := client.From(table).Select("*", "", false).Eq("email", email).Execute()
	var infoStruct []Row

	json.Unmarshal([]byte(string(raw_data)), &infoStruct)

	return infoStruct[0]
}
