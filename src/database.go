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

type EmailListResponse struct {
	EmailList []string
}

type Email struct {
	Emails []string
}

type Problem struct {
	Problems []map[string]string
}

func create_supabase_client() *supabase.Client {
	client, err := supabase.NewClient(os.Getenv("SUPABASE_URL"), os.Getenv("SUPABASE_KEY"), &supabase.ClientOptions{})
	if err != nil {
		fmt.Println("cannot initalize client", err)
	}
	return client
}

func add_username_and_email_to_database(leetcode_username string, new_email string) {
	client := create_supabase_client()
	table := os.Getenv("SUPABASE_TABLE")

	found_username, _, _ := client.From(table).Select("username", "", false).Eq("username", leetcode_username).Execute()

	if string(found_username) == "[]" {
		data := map[string]interface{}{
			"username":   leetcode_username,
			"email_list": []string{new_email},
			"problems":   []LeetCodeProblem{},
		}
		client.From(table).Insert(data, false, "Failure", "Success", "1").Execute()
		return
	}

	var old_emails []EmailListResponse
	var updatedEmails []string

	response, _, _ := client.From(table).Select("emails", "", false).Eq("username", leetcode_username).Execute()
	json.Unmarshal(response, &old_emails)

	for _, old_email := range old_emails {
		updatedEmails = append(updatedEmails, old_email.EmailList...)
	}
	updatedEmails = append(updatedEmails, new_email)

	data := map[string]interface{}{
		"email_list": updatedEmails,
	}
	client.From(table).Update(data, "", "").Eq("username", leetcode_username).Execute()
}

func get_emails_and_problems(leetcode_username string) ([]string, []map[string]string) {
	client := create_supabase_client()
	table := os.Getenv("SUPABASE_TABLE")

	emails_data, _, _ := client.From(table).Select("emails", "", false).Eq("username", leetcode_username).Execute()
	problems_data, _, _ := client.From(table).Select("problems", "", false).Eq("username", leetcode_username).Execute()

	var emailStruct []Email
	var problemStruct []Problem

	json.Unmarshal([]byte(string(emails_data)), &emailStruct)
	json.Unmarshal([]byte(string(problems_data)), &problemStruct)

	emails := emailStruct[0].Emails
	problems := problemStruct[0].Problems

	fmt.Println(emails)
	fmt.Println(problems)

	return emails, problems
}
