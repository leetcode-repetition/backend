package main

import (
	"fmt"
	"os"
	"time"
)

func get_recent_submissions() {
	data, _ := os.ReadFile("graphql/get_recent_submissions.graphql")
	query := string(data)

	fmt.Println(query)
}

func schedule(f func(), firstRun time.Duration) {
	time.Sleep(firstRun)
	ticker := time.NewTicker(24 * time.Hour)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			f()
		}
	}
}

func main() {
	// url := "https://leetcode.com/graphql"
	get_recent_submissions()
}
