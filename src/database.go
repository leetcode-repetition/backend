package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/supabase-community/supabase-go"
)

func create_supabase_client() (*supabase.Client, error) {
	client, err := supabase.NewClient(os.Getenv("SUPABASE_URL"), os.Getenv("SUPABASE_KEY"), &supabase.ClientOptions{})
	if err != nil {
		fmt.Println("Cannot initalize client", err)
	} else {
		fmt.Println("Initailized supabase client")
	}
	return client, err
}

func upsert_problem_into_database(username string, problem LeetCodeProblem) error {
	client, err := create_supabase_client()
	if err != nil {
		return err
	}

	table := os.Getenv("SUPABASE_TABLE")
	_, _, err = client.From(table).
		Upsert(map[string]interface{}{
			"username":           username,
			"titleSlug":          problem.TitleSlug,
			"link":               problem.Link,
			"difficulty":         problem.Difficulty,
			"repeatDate":         problem.RepeatDate,
			"lastCompletionDate": problem.LastCompletionDate,
		}, "username,titleSlug", "", "").
		Execute()

	if err != nil {
		fmt.Println("Error upserting database:", err)
	}
	fmt.Println("Successfully upserted database entry for user:", username)
	return err
}

func delete_problem_from_database(username string, problem_title_slug string) error {
	client, err := create_supabase_client()
	if err != nil {
		return err
	}

	table := os.Getenv("SUPABASE_TABLE")
	_, _, err = client.From(table).
		Delete("", "").
		Eq("username", username).
		Eq("titleSlug", problem_title_slug).
		Execute()

	if err != nil {
		fmt.Println("Error deleting database entry:", err)
	}
	fmt.Println("Successfully deleted database entry for user:", username)
	return err
}

func get_problems_from_database(username string) []LeetCodeProblem {
	var problems []LeetCodeProblem

	client, e := create_supabase_client()
	if e != nil {
		fmt.Println("Error creating supabase client:", e)
		return []LeetCodeProblem{}
	}
	table := os.Getenv("SUPABASE_TABLE")

	fmt.Println("Getting problems from database for user:", username)
	raw_data, _, err := client.From(table).Select("*", "", false).Eq("username", username).Execute()
	if err != nil {
		fmt.Println("Error fetching data:", err)
		return []LeetCodeProblem{}
	}

	fmt.Println("Raw data:", string(raw_data))

	var rawProblems []map[string]interface{}
	err = json.Unmarshal(raw_data, &rawProblems)
	if err != nil {
		fmt.Println("Error unmarshaling data:", err)
		return []LeetCodeProblem{}
	}
	for _, rawProblem := range rawProblems {
		problem := LeetCodeProblem{
			Link:               rawProblem["link"].(string),
			TitleSlug:          rawProblem["titleSlug"].(string),
			Difficulty:         rawProblem["difficulty"].(string),
			RepeatDate:         rawProblem["repeatDate"].(string),
			LastCompletionDate: rawProblem["lastCompletionDate"].(string),
		}
		problems = append(problems, problem)
	}

	now := time.Now()
	sort.Slice(problems, func(i, j int) bool {
		dateI, _ := time.Parse("1/2/06", problems[i].RepeatDate)
		dateJ, _ := time.Parse("1/2/06", problems[j].RepeatDate)
		diffI := dateI.Sub(now).Abs()
		diffJ := dateJ.Sub(now).Abs()
		return diffI < diffJ
	})

	fmt.Printf("Problems for user %s: %+v\n", username, problems)
	return problems
}
