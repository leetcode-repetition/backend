package shared

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/supabase-community/supabase-go"
)

func createSupabaseClient() (*supabase.Client, error) {
	client, err := supabase.NewClient(os.Getenv("SUPABASE_URL"), os.Getenv("SUPABASE_KEY"), &supabase.ClientOptions{})
	if err != nil {
		fmt.Println("Cannot initalize client", err)
	} else {
		fmt.Println("Initailized supabase client")
	}
	return client, err
}

func upsertProblemIntoDatabase(username string, problem LeetCodeProblem) error {
	client, err := createSupabaseClient()
	if err != nil {
		return err
	}

	table := os.Getenv("SUPABASE_TABLE")
	_, _, err = client.From(table).
		Upsert(map[string]interface{}{
			"username":           username,
			"titleSlug":          problem.TitleSlug,
			"link":               problem.Link,
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

func deleteProblemFromDatabase(username string, problem_title_slug string) error {
	client, err := createSupabaseClient()
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

func getProblemsFromDatabase(username string) []LeetCodeProblem {
	var problems []LeetCodeProblem

	client, e := createSupabaseClient()
	if e != nil {
		fmt.Println("Error creating supabase client:", e)
		return []LeetCodeProblem{}
	}
	table := os.Getenv("SUPABASE_TABLE")

	fmt.Println("Getting problems from database for user:", username)
	rawData, _, err := client.From(table).Select("*", "", false).Eq("username", username).Execute()
	if err != nil {
		fmt.Println("Error fetching data:", err)
		return []LeetCodeProblem{}
	}

	fmt.Println("Raw data:", string(rawData))

	var rawProblems []map[string]interface{}
	err = json.Unmarshal(rawData, &rawProblems)
	if err != nil {
		fmt.Println("Error unmarshaling data:", err)
		return []LeetCodeProblem{}
	}
	for _, rawProblem := range rawProblems {
		problem := LeetCodeProblem{
			Link:               rawProblem["link"].(string),
			TitleSlug:          rawProblem["titleSlug"].(string),
			RepeatDate:         rawProblem["repeatDate"].(string),
			LastCompletionDate: rawProblem["lastCompletionDate"].(string),
		}
		problems = append(problems, problem)
	}

	fmt.Printf("Problems for user %s: %+v\n", username, problems)
	return problems
}
