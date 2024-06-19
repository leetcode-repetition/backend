package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/supabase-community/supabase-go"
)

func create_supabase_client() {
	err := godotenv.Load("../.env")
	fmt.Println(err)
	supabase_url := os.Getenv("SUPABASE_URL")
	supabase_key := os.Getenv("SUPABASE_KEY")

	Options := &supabase.ClientOptions{}

	client, err := supabase.NewClient(supabase_url, supabase_key, Options)
	if err != nil {
		fmt.Println("cannot initalize client", err)
	}
	fmt.Println(client, err)

	data, count, err := client.From("leetcode_profiles").Select("*", "1", true).Execute()
	fmt.Println(string(data), count, err)
}

func add_username_and_email_to_database(leetcode_username string, email string) int {
	create_supabase_client()
	return 1
}
