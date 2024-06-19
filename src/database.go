package main

import (
	"fmt"
	"os"
	"time"

	"github.com/supabase-community/supabase-go"
)

type Problem struct {
	TitleSlug      string
	Difficulty     string
	Tags           []string
	CompletedDates []time.Time
	RepeatDate     time.Time
}

func create_supabase_client() *supabase.Client {
	client, err := supabase.NewClient(os.Getenv("SUPABASE_URL"), os.Getenv("SUPABASE_KEY"), &supabase.ClientOptions{})

	if err != nil {
		fmt.Println("cannot initalize client", err)
	}

	return client
}

func add_username_and_email_to_database(leetcode_username string, email string) {
	client := create_supabase_client()
	table := os.Getenv("SUPABASE_TABLE")

	found_username, _, _ := client.From(table).Select("username", "", false).Eq("username", leetcode_username).Execute()

	if string(found_username) == "[]" {
		// User not in db = ADD USERNAME to database
		data := map[string]interface{}{
			"username": leetcode_username,
			"email":    []string{email},
			"problems": []Problem{},
		}
		client.From(table).Insert(data, false, "Failure", "Success", "1")
	} else {
		// User IN db, update the email list
		data := map[string]interface{}{
			"email": []string{email},
		}
		client.From(table).Update(data, "append", "1").Eq("username", leetcode_username).Execute()
		fmt.Println("world")
	}
}
